package subsonic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"suasor/client/media/interfaces"
	"suasor/models"
	"suasor/utils"
	"time"

	gosonic "github.com/supersonic-app/go-subsonic/subsonic"
)

// SubsonicClient implements MediaContentProvider for Subsonic
type SubsonicClient struct {
	interfaces.BaseMediaClient
	config     models.SubsonicConfig
	httpClient *http.Client
	client     *gosonic.Client
	baseURL    string
}

// NewSubsonicClient creates a new Subsonic client
func NewSubsonicClient(clientID uint64, config models.SubsonicConfig) *SubsonicClient {
	// Get logger from context
	log := utils.LoggerFromContext(context.Background())

	log.Info().
		Uint64("clientID", clientID).
		Str("host", config.Host).
		Int("port", config.Port).
		Bool("ssl", config.SSL).
		Msg("Creating new Subsonic client")

	protocol := "http"
	if config.SSL {
		protocol = "https"
	}
	baseURL := fmt.Sprintf("%s://%s:%d", protocol, config.Host, config.Port)

	httpClient := &http.Client{Timeout: 30 * time.Second}

	// Create the go-subsonic client
	client := &gosonic.Client{
		Client:       httpClient,
		BaseUrl:      baseURL,
		User:         config.Username,
		ClientName:   "suasor",
		UserAgent:    "Suasor/1.0",
		PasswordAuth: true, // Using plain password auth for simplicity
	}

	// Authenticate with the Subsonic server
	err := client.Authenticate(config.Password)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to authenticate with Subsonic server")
	} else {
		log.Info().Msg("Successfully authenticated with Subsonic server")
	}

	return &SubsonicClient{
		BaseMediaClient: interfaces.BaseMediaClient{
			ClientID:   clientID,
			ClientType: models.MediaClientTypeSubsonic,
		},
		config:     config,
		httpClient: httpClient,
		client:     client,
		baseURL:    baseURL,
	}
}

// Register the provider factory
func init() {
	interfaces.RegisterProvider(models.MediaClientTypeSubsonic, func(ctx context.Context, clientID uint64, config interface{}) (interfaces.MediaContentProvider, error) {
		cfg, ok := config.(models.SubsonicConfig)
		if !ok {
			return nil, fmt.Errorf("invalid configuration for Subsonic client")
		}
		return NewSubsonicClient(clientID, cfg), nil
	})
}

// Capability methods - Subsonic only supports music
func (c *SubsonicClient) SupportsMusic() bool       { return true }
func (c *SubsonicClient) SupportsPlaylists() bool   { return true }
func (c *SubsonicClient) SupportsMovies() bool      { return false }
func (c *SubsonicClient) SupportsTVShows() bool     { return false }
func (c *SubsonicClient) SupportsBooks() bool       { return false }
func (c *SubsonicClient) SupportsCollections() bool { return false }

// GetMusic retrieves music tracks
func (c *SubsonicClient) GetMusic(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicTrack, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music tracks from Subsonic server")

	var tracks []interfaces.MusicTrack
	var err error

	// If filter/query provided, use search3
	if options != nil && len(options.Filters) > 0 {

		queryString := ""
		for _, value := range options.Filters {
			if queryString != "" {
				queryString += " "
			}
			queryString += value
		}
		tracks, err = c.searchMusic(ctx, queryString, options.Limit)
	} else {
		// Otherwise get random songs
		tracks, err = c.getRandomSongs(ctx, options)
	}

	if err != nil {
		return nil, err
	}

	// Add client info to each track
	for i := range tracks {
		c.AddClientInfo(&tracks[i].MediaItem)
	}

	log.Info().
		Int("trackCount", len(tracks)).
		Msg("Successfully retrieved music from Subsonic")

	return tracks, nil
}

// getRandomSongs retrieves random songs from the server
func (c *SubsonicClient) getRandomSongs(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicTrack, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Debug().Msg("Fetching random songs from Subsonic")

	params := make(map[string]string)

	if options != nil && options.Limit > 0 {
		params["size"] = strconv.Itoa(options.Limit)
	}

	resp, err := c.client.Get("getRandomSongs", params)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch random songs from Subsonic")
		return nil, err
	}

	if resp.RandomSongs == nil || len(resp.RandomSongs.Song) == 0 {
		log.Info().Msg("No songs returned from Subsonic")
		return []interfaces.MusicTrack{}, nil
	}

	tracks := make([]interfaces.MusicTrack, 0, len(resp.RandomSongs.Song))

	for _, song := range resp.RandomSongs.Song {
		track := convertChildToTrack(*song)
		track.ClientID = c.ClientID
		track.ClientType = string(c.ClientType)
		tracks = append(tracks, track)
	}

	return tracks, nil
}

// searchMusic searches for music by name
func (c *SubsonicClient) searchMusic(ctx context.Context, query string, limit int) ([]interfaces.MusicTrack, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Str("query", query).
		Int("limit", limit).
		Msg("Searching for music in Subsonic")

	params := make(map[string]string)
	params["query"] = query

	if limit > 0 {
		params["songCount"] = strconv.Itoa(limit)
	}

	resp, err := c.client.Get("search3", params)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", query).
			Msg("Failed to search for music in Subsonic")
		return nil, err
	}

	if resp.SearchResult3 == nil || len(resp.SearchResult3.Song) == 0 {
		log.Info().
			Str("query", query).
			Msg("No songs found matching query")
		return []interfaces.MusicTrack{}, nil
	}

	tracks := make([]interfaces.MusicTrack, 0, len(resp.SearchResult3.Song))

	for _, song := range resp.SearchResult3.Song {
		track := convertChildToTrack(*song)
		track.ClientID = c.ClientID
		track.ClientType = string(c.ClientType)
		tracks = append(tracks, track)
	}

	return tracks, nil
}

// Helper function to convert gosonic.Child to interfaces.MusicTrack
func convertChildToTrack(song gosonic.Child) interfaces.MusicTrack {
	// Convert duration from seconds to time.Duration
	duration := time.Duration(song.Duration) * time.Second

	track := interfaces.MusicTrack{
		MediaItem: interfaces.MediaItem{
			ExternalID: song.ID,
			Type:       "music",
			Metadata: interfaces.MediaMetadata{
				Title:       song.Title,
				Duration:    duration,
				ReleaseYear: song.Year, // Use ReleaseYear instead of Year
				Genres:      []string{song.Genre},
				Artwork: interfaces.Artwork{
					Poster: song.CoverArt, // Will be replaced with full URL later
				},
			},
		},
		ArtistName: song.Artist,
		AlbumTitle: song.Album,
		Number:     song.Track, // Number field instead of TrackNumber
	}

	return track
}

// GetMusicGenres retrieves available music genres
func (c *SubsonicClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music genres from Subsonic")

	resp, err := c.client.Get("getGenres", nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch genres from Subsonic")
		return nil, err
	}

	if resp.Genres == nil {
		log.Info().Msg("No genres returned from Subsonic")
		return []string{}, nil
	}

	genres := make([]string, 0, len(resp.Genres.Genre))
	for _, genre := range resp.Genres.Genre {
		genres = append(genres, genre.Name)
	}

	log.Info().
		Int("genreCount", len(genres)).
		Msg("Successfully retrieved music genres from Subsonic")

	return genres, nil
}

// GetPlaylists retrieves playlists from the Subsonic server
func (c *SubsonicClient) GetPlaylists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Playlist, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving playlists from Subsonic")

	resp, err := c.client.Get("getPlaylists", nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch playlists from Subsonic")
		return nil, err
	}

	if resp.Playlists == nil || len(resp.Playlists.Playlist) == 0 {
		log.Info().Msg("No playlists returned from Subsonic")
		return []interfaces.Playlist{}, nil
	}

	playlists := make([]interfaces.Playlist, 0, len(resp.Playlists.Playlist))

	for _, pl := range resp.Playlists.Playlist {
		playlist := interfaces.Playlist{
			MediaItem: interfaces.MediaItem{
				ExternalID: pl.ID,
				Type:       "playlist",
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Title:       pl.Name,
					Description: pl.Comment,
					Duration:    time.Duration(pl.Duration) * time.Second,
				},
			},
			ItemCount: pl.SongCount,
			Owner:     pl.Owner,
			IsPublic:  pl.Public,
		}

		// Add cover art if available
		if pl.CoverArt != "" {
			coverURL := c.GetCoverArtURL(pl.CoverArt)
			playlist.Metadata.Artwork.Poster = coverURL
		}

		playlists = append(playlists, playlist)
	}

	log.Info().
		Int("playlistCount", len(playlists)).
		Msg("Successfully retrieved playlists from Subsonic")

	return playlists, nil
}

// GetPlaylistItems retrieves tracks in a playlist
func (c *SubsonicClient) GetPlaylistItems(ctx context.Context, playlistID string) ([]interfaces.MusicTrack, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("playlistID", playlistID).
		Msg("Retrieving playlist items from Subsonic")

	params := make(map[string]string)
	params["id"] = playlistID

	resp, err := c.client.Get("getPlaylist", params)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to fetch playlist items from Subsonic")
		return nil, err
	}

	if resp.Playlist == nil || len(resp.Playlist.Entry) == 0 {
		log.Info().
			Str("playlistID", playlistID).
			Msg("No tracks found in playlist")
		return []interfaces.MusicTrack{}, nil
	}

	tracks := make([]interfaces.MusicTrack, 0, len(resp.Playlist.Entry))

	for _, song := range resp.Playlist.Entry {
		track := convertChildToTrack(*song)
		track.ClientID = c.ClientID
		track.ClientType = string(c.ClientType)
		tracks = append(tracks, track)
	}

	log.Info().
		Int("trackCount", len(tracks)).
		Str("playlistID", playlistID).
		Msg("Successfully retrieved playlist items from Subsonic")

	return tracks, nil
}

// GetMusicAlbums retrieves music albums from the Subsonic server
func (c *SubsonicClient) GetMusicAlbums(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicAlbum, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music albums from Subsonic")

	params := make(map[string]string)
	params["type"] = "alphabeticalByName"

	if options != nil {
		if options.Limit > 0 {
			params["size"] = strconv.Itoa(options.Limit)
		}
		if options.Offset > 0 {
			params["offset"] = strconv.Itoa(options.Offset)
		}
	}

	// Use getAlbumList2 which is tag-based instead of folder-based
	resp, err := c.client.Get("getAlbumList2", params)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch albums from Subsonic")
		return nil, err
	}

	if resp.AlbumList2 == nil || len(resp.AlbumList2.Album) == 0 {
		log.Info().Msg("No albums returned from Subsonic")
		return []interfaces.MusicAlbum{}, nil
	}

	albums := make([]interfaces.MusicAlbum, 0, len(resp.AlbumList2.Album))

	for _, album := range resp.AlbumList2.Album {
		musicAlbum := interfaces.MusicAlbum{
			MediaItem: interfaces.MediaItem{
				ExternalID: album.ID,
				Type:       "album",
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Title:       album.Name,
					ReleaseYear: album.Year,
					Duration:    time.Duration(album.Duration) * time.Second,
					Genres:      []string{album.Genre},
					Artwork: interfaces.Artwork{
						Poster: c.GetCoverArtURL(album.CoverArt),
					},
				},
			},
			ArtistName: album.Artist,
			TrackCount: album.SongCount,
		}
		albums = append(albums, musicAlbum)
	}

	log.Info().
		Int("albumCount", len(albums)).
		Msg("Successfully retrieved music albums from Subsonic")

	return albums, nil
}

// GetMusicArtists retrieves music artists from the Subsonic server
func (c *SubsonicClient) GetMusicArtists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicArtist, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music artists from Subsonic")

	resp, err := c.client.Get("getArtists", nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch artists from Subsonic")
		return nil, err
	}

	if resp.Artists == nil || len(resp.Artists.Index) == 0 {
		log.Info().Msg("No artists returned from Subsonic")
		return []interfaces.MusicArtist{}, nil
	}

	var artists []interfaces.MusicArtist

	// Flatten all artists from all indexes
	for _, index := range resp.Artists.Index {
		for _, artist := range index.Artist {
			// Apply pagination if needed
			if options != nil && options.Limit > 0 && len(artists) >= options.Limit {
				break
			}

			musicArtist := interfaces.MusicArtist{
				MediaItem: interfaces.MediaItem{
					ExternalID: artist.ID,
					Type:       "artist",
					ClientID:   c.ClientID,
					ClientType: string(c.ClientType),
					Metadata: interfaces.MediaMetadata{
						Title: artist.Name,
					},
				},
			}

			// Add cover art if available
			if artist.CoverArt != "" {
				musicArtist.Metadata.Artwork.Poster = c.GetCoverArtURL(artist.CoverArt)
			}

			artists = append(artists, musicArtist)
		}
	}

	log.Info().
		Int("artistCount", len(artists)).
		Msg("Successfully retrieved music artists from Subsonic")

	return artists, nil
}

// GetMusicTrackByID retrieves a specific music track by ID
func (c *SubsonicClient) GetMusicTrackByID(ctx context.Context, id string) (interfaces.MusicTrack, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", id).
		Msg("Retrieving specific music track from Subsonic")

	params := make(map[string]string)
	params["id"] = id

	resp, err := c.client.Get("getSong", params)
	if err != nil {
		log.Error().
			Err(err).
			Str("trackID", id).
			Msg("Failed to fetch track from Subsonic")
		return interfaces.MusicTrack{}, err
	}

	if resp.Song == nil {
		log.Error().
			Str("trackID", id).
			Msg("No track found with the specified ID")
		return interfaces.MusicTrack{}, fmt.Errorf("track with ID %s not found", id)
	}

	track := convertChildToTrack(*resp.Song)
	track.ClientID = c.ClientID
	track.ClientType = string(c.ClientType)

	log.Info().
		Str("trackID", id).
		Str("title", track.Metadata.Title).
		Str("artist", track.ArtistName).
		Msg("Successfully retrieved music track from Subsonic")

	return track, nil
}

// GetStreamURL returns the URL to stream a music track
func (c *SubsonicClient) GetStreamURL(ctx context.Context, trackID string) (string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", trackID).
		Msg("Generating stream URL for track")

	// We can't access the unexported setupRequest method, so build the URL manually
	protocol := "http"
	if c.config.SSL {
		protocol = "https"
	}

	// Create query parameters
	params := url.Values{}
	params.Add("id", trackID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", c.config.Username)
	params.Add("p", c.config.Password)

	streamURL := fmt.Sprintf("%s://%s:%d/rest/stream.view?%s",
		protocol, c.config.Host, c.config.Port, params.Encode())

	log.Debug().
		Str("trackID", trackID).
		Str("streamURL", streamURL).
		Msg("Generated stream URL for track")

	return streamURL, nil
}

// GetCoverArtURL returns the URL to download cover art
func (c *SubsonicClient) GetCoverArtURL(coverArtID string) string {
	if coverArtID == "" {
		return ""
	}

	// We can't access the unexported setupRequest method, so build the URL manually
	protocol := "http"
	if c.config.SSL {
		protocol = "https"
	}

	// Create query parameters
	params := url.Values{}
	params.Add("id", coverArtID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", c.config.Username)
	params.Add("p", c.config.Password)

	return fmt.Sprintf("%s://%s:%d/rest/getCoverArt.view?%s",
		protocol, c.config.Host, c.config.Port, params.Encode())
}

// The following methods are unsupported by Subsonic (music server only)

// Unsupported methods just return ErrFeatureNotSupported
var ErrFeatureNotSupported = fmt.Errorf("feature not supported by Subsonic")

func (c *SubsonicClient) GetMovies(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Movie, error) {
	return nil, ErrFeatureNotSupported
}

func (c *SubsonicClient) GetMovieByID(ctx context.Context, id string) (interfaces.Movie, error) {
	return interfaces.Movie{}, ErrFeatureNotSupported
}

func (c *SubsonicClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	return nil, ErrFeatureNotSupported
}

func (c *SubsonicClient) GetTVShows(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.TVShow, error) {
	return nil, ErrFeatureNotSupported
}

func (c *SubsonicClient) GetTVShowByID(ctx context.Context, id string) (interfaces.TVShow, error) {
	return interfaces.TVShow{}, ErrFeatureNotSupported
}

func (c *SubsonicClient) GetTVShowSeasons(ctx context.Context, showID string) ([]interfaces.Season, error) {
	return nil, ErrFeatureNotSupported
}

func (c *SubsonicClient) GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]interfaces.Episode, error) {
	return nil, ErrFeatureNotSupported
}

func (c *SubsonicClient) GetEpisodeByID(ctx context.Context, id string) (interfaces.Episode, error) {
	return interfaces.Episode{}, ErrFeatureNotSupported
}

func (c *SubsonicClient) GetCollections(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Collection, error) {
	return nil, ErrFeatureNotSupported
}

func (c *SubsonicClient) GetWatchHistory(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.WatchHistoryItem, error) {
	return nil, ErrFeatureNotSupported
}
