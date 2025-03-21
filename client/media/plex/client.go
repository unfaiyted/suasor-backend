// plex/client.go
package plex

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"suasor/client/media/interfaces"
	"suasor/models"
	"time"

	"github.com/LukeHagar/plexgo"
	"github.com/LukeHagar/plexgo/models/operations"
)

// init is automatically called when package is imported
func init() {
	interfaces.RegisterProvider(models.MediaClientTypePlex, NewPlexClient)
}

// PlexClient implements MediaContentProvider for Plex
type PlexClient struct {
	interfaces.BaseMediaClient
	config     models.PlexConfig
	httpClient *http.Client
	baseURL    string
	plexAPI    *plexgo.PlexAPI
}

// NewPlexClient creates a new Plex client
func NewPlexClient(ctx context.Context, clientID uint64, config any) (interfaces.MediaContentProvider, error) {
	plexConfig, ok := config.(models.PlexConfig)
	if !ok {
		return nil, fmt.Errorf("invalid Plex configuration")
	}

	// Initialize the Plex API client
	plexAPI := plexgo.New(
		plexgo.WithSecurity(plexConfig.Token),
	)

	return &PlexClient{
		BaseMediaClient: interfaces.BaseMediaClient{
			ClientID:   clientID,
			ClientType: models.MediaClientTypePlex,
		},
		config:  plexConfig,
		plexAPI: plexAPI,
		baseURL: plexConfig.Host,
	}, nil
}

// Capability methods
func (c *PlexClient) SupportsMovies() bool      { return true }
func (c *PlexClient) SupportsTVShows() bool     { return true }
func (c *PlexClient) SupportsMusic() bool       { return true }
func (c *PlexClient) SupportsPlaylists() bool   { return true }
func (c *PlexClient) SupportsCollections() bool { return true }

// Helper functions for common operations

// makeFullURL creates a complete URL from a resource path
func (c *PlexClient) makeFullURL(resourcePath string) string {
	if resourcePath == "" {
		return ""
	}

	if strings.HasPrefix(resourcePath, "http") {
		return resourcePath
	}

	return fmt.Sprintf("%s%s", c.baseURL, resourcePath)
}

// findLibrarySectionByType returns the section key for the specified type
func (c *PlexClient) findLibrarySectionByType(ctx context.Context, sectionType string) (string, error) {
	libraries, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get libraries: %w", err)
	}

	for _, dir := range libraries.Object.MediaContainer.Directory {
		if dir.Type == sectionType {
			return dir.Key, nil
		}
	}

	return "", nil
}

func (c *PlexClient) createChildMetadataFromPlexItem(item *operations.GetMetadataChildrenMetadata) interfaces.MediaMetadata {
	metadata := interfaces.MediaMetadata{
		Title:       *item.Title,
		Description: *item.Summary,
		Artwork: interfaces.Artwork{
			Thumbnail: c.makeFullURL(*item.Thumb),
		},
		ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
			Source: "plex",
			ID:     *item.RatingKey,
		}},
	}

	// Add optional fields if present
	if item.UpdatedAt != nil {
		metadata.UpdatedAt = time.Unix(int64(*item.UpdatedAt), 0)
	}
	if *item.AddedAt != 0 {
		metadata.AddedAt = time.Unix(int64(*item.AddedAt), 0)
	}
	if item.ParentYear != nil {
		metadata.ReleaseYear = *item.ParentYear
	}

	return metadata
}

// createMetadataFromPlexItem creates a MediaMetadata from a Plex item
func (c *PlexClient) createMetadataFromPlexItem(item *operations.GetLibraryItemsMetadata) interfaces.MediaMetadata {
	metadata := interfaces.MediaMetadata{
		Title:       item.Title,
		Description: item.Summary,
		Artwork: interfaces.Artwork{
			Thumbnail: c.makeFullURL(*item.Thumb),
		},
		ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
			Source: "plex",
			ID:     item.RatingKey,
		}},
	}

	// Add optional fields if present
	if item.UpdatedAt != nil {
		metadata.UpdatedAt = time.Unix(int64(*item.UpdatedAt), 0)
	}
	if item.AddedAt != 0 {
		metadata.AddedAt = time.Unix(int64(item.AddedAt), 0)
	}
	if item.Year != nil {
		metadata.ReleaseYear = *item.Year
	}
	if item.Rating != nil {
		metadata.Ratings = interfaces.Ratings{
			interfaces.Rating{
				Source: "plex",
				Value:  float32(*item.Rating),
			},
		}
	}
	if item.Duration != nil {
		metadata.Duration = time.Duration(*item.Duration) * time.Millisecond
	}
	if item.Studio != nil {
		metadata.Studios = []string{*item.Studio}
	}
	if item.ContentRating != nil {
		metadata.ContentRating = *item.ContentRating
	}

	// Add genres if present
	if item.Genre != nil {
		metadata.Genres = make([]string, 0, len(item.Genre))
		for _, genre := range item.Genre {
			if genre.Tag != nil {
				metadata.Genres = append(metadata.Genres, *genre.Tag)
			}
		}
	}

	return metadata
}

// createMetadataFromPlexItem creates a MediaMetadata from a Plex item
func (c *PlexClient) createMediaMetadataFromPlexItem(item *operations.GetMediaMetaDataMetadata) interfaces.MediaMetadata {
	metadata := interfaces.MediaMetadata{
		Title:       item.Title,
		Description: item.Summary,
		Artwork:     interfaces.Artwork{
			// Thumbnail: c.makeFullURL(*item.Thumb),
		},
		ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
			Source: "plex",
			ID:     item.RatingKey,
		}},
	}

	// Add optional fields if present
	if item.AddedAt != 0 {
		metadata.AddedAt = time.Unix(int64(item.AddedAt), 0)
	}

	metadata.UpdatedAt = time.Unix(int64(item.UpdatedAt), 0)
	metadata.ReleaseYear = item.Year

	if item.Rating != nil {
		metadata.Ratings = interfaces.Ratings{
			interfaces.Rating{
				Source: "plex",
				Value:  float32(*item.Rating),
			},
		}
	} // if item.Duration != nil {
	// 	metadata.Duration = *item.Duration
	// }
	if item.Studio != nil {
		metadata.Studios = []string{*item.Studio}
	}
	if item.ContentRating != nil {
		metadata.ContentRating = *item.ContentRating
	}

	// Add genres if present
	if item.Genre != nil {
		metadata.Genres = make([]string, 0, len(item.Genre))
		for _, genre := range item.Genre {
			metadata.Genres = append(metadata.Genres, genre.Tag)
		}
	}

	return metadata
}

// GetPlaylists retrieves playlists from Plex
func (c *PlexClient) GetPlaylists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Playlist, error) {
	res, err := c.plexAPI.Playlists.GetPlaylists(ctx, operations.PlaylistTypeAudio.ToPointer(), operations.QueryParamSmartOne.ToPointer())
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		return []interfaces.Playlist{}, nil
	}

	playlists := make([]interfaces.Playlist, 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		playlist := interfaces.Playlist{
			MediaItem: interfaces.MediaItem{
				ExternalID: *item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Description: *item.Summary,
					Title:       *item.Title,
					Artwork:     interfaces.Artwork{
						// Thumbnail: c.makeFullURL(*item.Thumb),
					},
					ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
						Source: "plex",
						ID:     *item.RatingKey,
					}},
					UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
					AddedAt:   time.Unix(int64(*item.AddedAt), 0),
				},
			},
		}
		c.BaseMediaClient.AddClientInfo(&playlist.MediaItem)
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

// GetCollections retrieves collections from Plex
func (c *PlexClient) GetCollections(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Collection, error) {
	res, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	directories := res.Object.MediaContainer.GetDirectory()
	collections := make([]interfaces.Collection, 0, len(directories))

	for _, dir := range directories {
		collection := interfaces.Collection{
			MediaItem: interfaces.MediaItem{
				ExternalID: dir.Key,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Title: dir.Title,
					Artwork: interfaces.Artwork{
						Thumbnail: c.makeFullURL(dir.Thumb),
					},
					ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
						Source: "plex",
						ID:     dir.Key,
					}},
				},
			},
		}

		// // Only add these fields if available
		// if dir.LeafCount != nil {
		// 	collection.ItemCount = *dir.LeafCount
		// }

		c.BaseMediaClient.AddClientInfo(&collection.MediaItem)
		collections = append(collections, collection)
	}

	return collections, nil
}

// GetMovies retrieves movies from Plex
func (c *PlexClient) GetMovies(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Movie, error) {
	// First, find the movie library section
	movieSectionKey, err := c.findLibrarySectionByType(ctx, "movie")
	if err != nil {
		return nil, err
	}

	if movieSectionKey == "" {
		return []interfaces.Movie{}, nil
	}

	// Get movies from the movie section
	sectionKey, _ := strconv.Atoi(movieSectionKey)
	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
		SectionKey: sectionKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		return []interfaces.Movie{}, nil
	}

	movies := make([]interfaces.Movie, 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		if item.Type != "movie" {
			continue
		}

		movie := interfaces.Movie{
			MediaItem: interfaces.MediaItem{
				ExternalID: item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata:   c.createMetadataFromPlexItem(&item),
			},
		}

		c.BaseMediaClient.AddClientInfo(&movie.MediaItem)
		movies = append(movies, movie)
	}

	return movies, nil
}

// GetTVShows retrieves TV shows from Plex
func (c *PlexClient) GetTVShows(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.TVShow, error) {
	// First, find the TV show library section
	tvSectionKey, err := c.findLibrarySectionByType(ctx, "show")
	if err != nil {
		return nil, err
	}

	if tvSectionKey == "" {
		return []interfaces.TVShow{}, nil
	}

	// Get TV shows from the TV section
	sectionKey, _ := strconv.Atoi(tvSectionKey)
	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
		SectionKey: sectionKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get TV shows: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		return []interfaces.TVShow{}, nil
	}

	shows := make([]interfaces.TVShow, 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		if item.Type != "show" {
			continue
		}

		show := interfaces.TVShow{
			MediaItem: interfaces.MediaItem{
				ExternalID: item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata:   c.createMetadataFromPlexItem(&item),
			},
		}

		if item.Rating != nil {
			show.Rating = float64(*item.Rating)
		}
		if item.Year != nil {
			show.ReleaseYear = *item.Year
		}
		if item.ContentRating != nil {
			show.ContentRating = *item.ContentRating
		}
		if item.ChildCount != nil {
			show.SeasonCount = *item.ChildCount
		}
		if item.LeafCount != nil {
			show.EpisodeCount = int(*item.LeafCount)
		}

		if item.Genre != nil {
			show.Genres = make([]string, 0, len(item.Genre))
			for _, genre := range item.Genre {
				if genre.Tag != nil {
					show.Genres = append(show.Genres, *genre.Tag)
				}
			}
		}

		c.BaseMediaClient.AddClientInfo(&show.MediaItem)
		shows = append(shows, show)
	}

	return shows, nil
}

// GetTVShowSeasons retrieves seasons for a specific TV show
func (c *PlexClient) GetTVShowSeasons(ctx context.Context, showID string) ([]interfaces.Season, error) {
	ratingKey, _ := strconv.Atoi(showID)
	float64RatingKey := float64(ratingKey)

	childRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64RatingKey, plexgo.String("Stream"))
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show seasons: %w", err)
	}

	if childRes.Object.MediaContainer == nil || childRes.Object.MediaContainer.Metadata == nil {
		return []interfaces.Season{}, nil
	}

	seasons := make([]interfaces.Season, 0, len(childRes.Object.MediaContainer.Metadata))
	for _, item := range childRes.Object.MediaContainer.Metadata {
		if *item.Type != "season" {
			continue
		}

		season := interfaces.Season{
			MediaItem: interfaces.MediaItem{
				ExternalID: *item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Description: *item.Summary,
					Title:       *item.Title,
					Artwork: interfaces.Artwork{
						Thumbnail: c.makeFullURL(*item.Thumb),
					},
					ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
						Source: "plex",
						ID:     *item.RatingKey,
					}},
					UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
					AddedAt:   time.Unix(int64(*item.AddedAt), 0),
				},
			},
			EpisodeCount: *item.LeafCount,
			Number:       *item.Index,
		}

		c.BaseMediaClient.AddClientInfo(&season.MediaItem)
		seasons = append(seasons, season)
	}

	return seasons, nil
}

// GetTVShowEpisodes retrieves episodes for a specific season of a TV show
func (c *PlexClient) GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]interfaces.Episode, error) {
	// First get all seasons
	seasons, err := c.GetTVShowSeasons(ctx, showID)
	if err != nil {
		return nil, err
	}

	var seasonID string
	for _, season := range seasons {
		if season.Number == seasonNumber {
			for _, externalID := range season.MediaItem.Metadata.ExternalIDs {
				if externalID.Source == "plex" {
					seasonID = externalID.ID
					break
				}
			}
			break
		}
	}

	if seasonID == "" {
		return []interfaces.Episode{}, nil
	}

	ratingKey, _ := strconv.Atoi(seasonID)
	float64RatingKey := float64(ratingKey)
	childRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64RatingKey, plexgo.String("Stream"))
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show episodes: %w", err)
	}

	if childRes.Object.MediaContainer == nil || childRes.Object.MediaContainer.Metadata == nil {
		return []interfaces.Episode{}, nil
	}

	episodes := make([]interfaces.Episode, 0, len(childRes.Object.MediaContainer.Metadata))
	for _, item := range childRes.Object.MediaContainer.Metadata {
		if *item.Type != "episode" {
			continue
		}

		episode := interfaces.Episode{
			MediaItem: interfaces.MediaItem{
				ExternalID: *item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Description: *item.Summary,
					Title:       *item.Title,
					Artwork: interfaces.Artwork{
						Thumbnail: c.makeFullURL(*item.Thumb),
					},
					// ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
					// 	Source: "plex",
					// 	ID:     *item.RatingKey,
					// }},
					UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
					AddedAt:   time.Unix(int64(*item.AddedAt), 0),
				},
			},
			Number:       int64(*item.Index),
			SeasonNumber: int(*item.ParentIndex),
		}

		// Add studio if available
		if item.ParentStudio != nil {
			episode.MediaItem.Metadata.Studios = []string{*item.ParentStudio}
		}

		c.BaseMediaClient.AddClientInfo(&episode.MediaItem)
		episodes = append(episodes, episode)
	}

	return episodes, nil
}

// GetMusic retrieves music tracks from Plex
func (c *PlexClient) GetMusic(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicTrack, error) {
	// Find the music library section
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		return nil, err
	}

	if musicSectionKey == "" {
		return []interfaces.MusicTrack{}, nil
	}

	// For tracks, we need to traverse the hierarchy: artists > albums > tracks
	sectionKey, _ := strconv.Atoi(musicSectionKey)
	float64SectionKey := float64(sectionKey)
	res, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64SectionKey, plexgo.String("Stream"))
	if err != nil {
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	var tracks []interfaces.MusicTrack

	// Loop through artists
	if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
		for _, artist := range res.Object.MediaContainer.Metadata {
			artistKey, _ := strconv.Atoi(*artist.RatingKey)
			float64ArtistKey := float64(artistKey)
			albumsRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64ArtistKey, plexgo.String("Stream"))
			if err != nil {
				continue
			}

			// Loop through albums
			if albumsRes.Object.MediaContainer != nil && albumsRes.Object.MediaContainer.Metadata != nil {
				for _, album := range albumsRes.Object.MediaContainer.Metadata {
					albumKey, _ := strconv.Atoi(*album.RatingKey)

					float64AlbumKey := float64(albumKey)
					tracksRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64AlbumKey, plexgo.String("Stream"))
					if err != nil {
						continue
					}

					// Loop through tracks
					if tracksRes.Object.MediaContainer != nil && tracksRes.Object.MediaContainer.Metadata != nil {
						for _, item := range tracksRes.Object.MediaContainer.Metadata {
							if *item.Type != "track" {
								continue
							}

							track := interfaces.MusicTrack{
								MediaItem: interfaces.MediaItem{
									ExternalID: *item.RatingKey,
									ClientID:   c.ClientID,
									ClientType: string(c.ClientType),
									Metadata:   c.createChildMetadataFromPlexItem(&item),
								},
								Number:     *item.Index,
								ArtistID:   *artist.RatingKey,
								ArtistName: *artist.Title,
								AlbumID:    *album.RatingKey,
								AlbumName:  *album.Title,
							}
							c.BaseMediaClient.AddClientInfo(&track.MediaItem)
							tracks = append(tracks, track)

							// Limit number of tracks to avoid too large responses
							if len(tracks) >= 100 {
								return tracks, nil
							}
						}
					}
				}
			}
		}
	}

	return tracks, nil
}

// GetMusicArtists retrieves music artists from Plex
func (c *PlexClient) GetMusicArtists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicArtist, error) {
	// Find the music library section
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		return nil, err
	}

	if musicSectionKey == "" {
		return []interfaces.MusicArtist{}, nil
	}

	sectionKey, _ := strconv.Atoi(musicSectionKey)
	float64SectionKey := float64(sectionKey)
	res, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64SectionKey, plexgo.String("Stream"))
	if err != nil {
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		return []interfaces.MusicArtist{}, nil
	}

	artists := make([]interfaces.MusicArtist, 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		if *item.Type != "artist" {
			continue
		}

		artist := interfaces.MusicArtist{
			MediaItem: interfaces.MediaItem{
				ExternalID: *item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata:   c.createChildMetadataFromPlexItem(&item),
			},
		}

		c.BaseMediaClient.AddClientInfo(&artist.MediaItem)
		artists = append(artists, artist)
	}

	return artists, nil
}

// GetMusicAlbums retrieves music albums from Plex
func (c *PlexClient) GetMusicAlbums(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicAlbum, error) {
	// Find the music library section
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		return nil, err
	}

	if musicSectionKey == "" {
		return []interfaces.MusicAlbum{}, nil
	}

	// For albums, we need to traverse artists first
	sectionKey, _ := strconv.Atoi(musicSectionKey)
	float64SectionKey := float64(sectionKey)
	res, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64SectionKey, plexgo.String("Stream"))
	if err != nil {
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	var albums []interfaces.MusicAlbum

	// Loop through artists to get their albums
	if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
		for _, artist := range res.Object.MediaContainer.Metadata {
			artistKey, _ := strconv.Atoi(*artist.RatingKey)
			float64ArtistKey := float64(artistKey)
			albumsRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64ArtistKey, plexgo.String("Stream"))
			if err != nil {
				continue
			}

			if albumsRes.Object.MediaContainer != nil && albumsRes.Object.MediaContainer.Metadata != nil {
				for _, item := range albumsRes.Object.MediaContainer.Metadata {
					if *item.Type != "album" {
						continue
					}

					album := interfaces.MusicAlbum{
						MediaItem: interfaces.MediaItem{
							ExternalID: *item.RatingKey,
							ClientID:   c.ClientID,
							ClientType: string(c.ClientType),
							Metadata:   c.createChildMetadataFromPlexItem(&item),
						},
						ArtistID:   *artist.RatingKey,
						ArtistName: *artist.Title,
						TrackCount: *item.LeafCount,
					}

					c.BaseMediaClient.AddClientInfo(&album.MediaItem)
					albums = append(albums, album)
				}
			}
		}
	}

	return albums, nil
}

// GetWatchHistory retrieves watch history from Plex
func (c *PlexClient) GetWatchHistory(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.WatchHistoryItem, error) {
	// This would require querying Plex for watch history
	return []interfaces.WatchHistoryItem{}, fmt.Errorf("Watch history retrieval not yet implemented for Plex")
}

// GetMovieByID retrieves a specific movie by ID
func (c *PlexClient) GetMovieByID(ctx context.Context, id string) (interfaces.Movie, error) {
	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)
	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})
	if err != nil {
		return interfaces.Movie{}, fmt.Errorf("failed to get movie: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		return interfaces.Movie{}, fmt.Errorf("movie not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "movie" {
		return interfaces.Movie{}, fmt.Errorf("item is not a movie")
	}

	movie := interfaces.Movie{
		MediaItem: interfaces.MediaItem{
			ExternalID: item.RatingKey,
			ClientID:   c.ClientID,
			ClientType: string(c.ClientType),
			Metadata:   c.createMediaMetadataFromPlexItem(&item),
		},
	}

	c.BaseMediaClient.AddClientInfo(&movie.MediaItem)
	return movie, nil
}

// GetTVShowByID retrieves a specific TV show by ID
func (c *PlexClient) GetTVShowByID(ctx context.Context, id string) (interfaces.TVShow, error) {
	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)
	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})

	if err != nil {
		return interfaces.TVShow{}, fmt.Errorf("failed to get TV show: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		return interfaces.TVShow{}, fmt.Errorf("TV show not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "show" {
		return interfaces.TVShow{}, fmt.Errorf("item is not a TV show")
	}

	show := interfaces.TVShow{
		MediaItem: interfaces.MediaItem{
			ExternalID: item.RatingKey,
			ClientID:   c.ClientID,
			ClientType: string(c.ClientType),
			Metadata:   c.createMediaMetadataFromPlexItem(&item),
		},
	}

	if item.Rating != nil {
		show.Rating = float64(*item.Rating)
	}
	if item.ContentRating != nil {
		show.ContentRating = *item.ContentRating
	}
	if item.ChildCount != nil {
		show.SeasonCount = *item.ChildCount
	}
	if item.LeafCount != nil {
		show.EpisodeCount = int(*item.LeafCount)
	}

	if item.Genre != nil {
		show.Genres = make([]string, 0, len(item.Genre))
		for _, genre := range item.Genre {
			show.Genres = append(show.Genres, genre.Tag)
		}
	}

	c.BaseMediaClient.AddClientInfo(&show.MediaItem)
	return show, nil
}

// GetEpisodeByID retrieves a specific episode by ID
func (c *PlexClient) GetEpisodeByID(ctx context.Context, id string) (interfaces.Episode, error) {
	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)
	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})
	if err != nil {
		return interfaces.Episode{}, fmt.Errorf("failed to get episode: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		return interfaces.Episode{}, fmt.Errorf("episode not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "episode" {
		return interfaces.Episode{}, fmt.Errorf("item is not an episode")
	}

	// Get parent season info
	var seasonNumber int
	var showID string

	if item.ParentRatingKey != nil && item.ParentIndex != nil {
		seasonNumber = int(*item.ParentIndex)

		// Get show ID from parent season
		seasonKey, _ := strconv.Atoi(*item.ParentRatingKey)
		int64SeasonKey := int64(seasonKey)
		seasonRes, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
			RatingKey: int64SeasonKey,
		})

		if err == nil && seasonRes.Object.MediaContainer != nil &&
			seasonRes.Object.MediaContainer.Metadata != nil &&
			len(seasonRes.Object.MediaContainer.Metadata) > 0 {

			if seasonRes.Object.MediaContainer.Metadata[0].ParentRatingKey != nil {
				showID = *seasonRes.Object.MediaContainer.Metadata[0].ParentRatingKey
			}
		}
	}

	episode := interfaces.Episode{
		MediaItem: interfaces.MediaItem{
			ExternalID: item.RatingKey,
			ClientID:   c.ClientID,
			ClientType: string(c.ClientType),
			Metadata:   c.createMediaMetadataFromPlexItem(&item),
		},
		ShowID:       showID,
		SeasonNumber: seasonNumber,
		Number:       *item.Index,
	}

	c.BaseMediaClient.AddClientInfo(&episode.MediaItem)
	return episode, nil
}

// GetMusicTrackByID retrieves a specific music track by ID
func (c *PlexClient) GetMusicTrackByID(ctx context.Context, id string) (interfaces.MusicTrack, error) {
	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)
	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{RatingKey: int64RatingKey})
	if err != nil {
		return interfaces.MusicTrack{}, fmt.Errorf("failed to get music track: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		return interfaces.MusicTrack{}, fmt.Errorf("music track not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "track" {
		return interfaces.MusicTrack{}, fmt.Errorf("item is not a music track")
	}

	// Get album and artist info
	var albumName string
	var artistID string
	var artistName string

	if item.ParentRatingKey != nil {
		// Get album info
		albumKey, _ := strconv.Atoi(*item.ParentRatingKey)
		int64AlbumKey := int64(albumKey)
		albumRes, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
			RatingKey: int64AlbumKey,
		})

		if err == nil && albumRes.Object.MediaContainer != nil &&
			albumRes.Object.MediaContainer.Metadata != nil &&
			len(albumRes.Object.MediaContainer.Metadata) > 0 {

			albumName = albumRes.Object.MediaContainer.Metadata[0].Title

			// Get artist info if available
			if albumRes.Object.MediaContainer.Metadata[0].ParentRatingKey != nil {
				artistID = *albumRes.Object.MediaContainer.Metadata[0].ParentRatingKey

				artistKey, _ := strconv.Atoi(artistID)
				int64ArtistKey := int64(artistKey)
				artistRes, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
					RatingKey: int64ArtistKey,
				})

				if err == nil && artistRes.Object.MediaContainer != nil &&
					artistRes.Object.MediaContainer.Metadata != nil &&
					len(artistRes.Object.MediaContainer.Metadata) > 0 {

					artistName = artistRes.Object.MediaContainer.Metadata[0].Title
				}
			}
		}
	}

	track := interfaces.MusicTrack{
		MediaItem: interfaces.MediaItem{
			ExternalID: item.RatingKey,
			ClientID:   c.ClientID,
			ClientType: string(c.ClientType),
			Metadata:   c.createMediaMetadataFromPlexItem(&item),
		},
		AlbumName:  albumName,
		ArtistName: artistName,
		ArtistID:   artistID,
		AlbumID:    *item.ParentRatingKey,
		Number:     int(*item.Index),
	}

	c.BaseMediaClient.AddClientInfo(&track.MediaItem)
	return track, nil
}

// GetMusicGenres retrieves music genres from Plex
func (c *PlexClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	// Find the music library section
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		return nil, err
	}

	if musicSectionKey == "" {
		return []string{}, nil
	}

	// Get genres from the library items
	sectionKey, _ := strconv.Atoi(musicSectionKey)
	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{SectionKey: sectionKey})
	if err != nil {
		return nil, fmt.Errorf("failed to get music content: %w", err)
	}

	genreMap := make(map[string]bool)
	if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
		for _, item := range res.Object.MediaContainer.Metadata {
			if item.Genre != nil {
				for _, genre := range item.Genre {
					if genre.Tag != nil {
						genreMap[*genre.Tag] = true
					}
				}
			}
		}
	}

	genres := make([]string, 0, len(genreMap))
	for genre := range genreMap {
		genres = append(genres, genre)
	}

	return genres, nil
}

// GetMovieGenres retrieves movie genres from Plex
func (c *PlexClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	// Find the movie library section
	movieSectionKey, err := c.findLibrarySectionByType(ctx, "movie")
	if err != nil {
		return nil, err
	}

	if movieSectionKey == "" {
		return []string{}, nil
	}

	// Get genres directly from the genre endpoint
	sectionKey, _ := strconv.Atoi(movieSectionKey)
	res, err := c.plexAPI.Library.GetGenresLibrary(ctx, sectionKey, operations.GetGenresLibraryQueryParamTypeMovie)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie genres: %w", err)
	}

	genreMap := make(map[string]bool)
	if res.Object.MediaContainer != nil {
		for _, item := range res.Object.MediaContainer.GetDirectory() {
			genreMap[item.Title] = true
		}
	}

	genres := make([]string, 0, len(genreMap))
	for genre := range genreMap {
		genres = append(genres, genre)
	}

	return genres, nil
}
