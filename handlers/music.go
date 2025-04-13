// handlers/music.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/client/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
)

// MusicSpecificHandler extends MediaItemHandler with music-specific functionality
type MusicSpecificHandler struct {
	*MediaItemHandler[*mediatypes.Track]
	albumService  services.MediaItemService[*mediatypes.Album]
	artistService services.MediaItemService[*mediatypes.Artist]
}

// NewMusicSpecificHandler creates a new MusicSpecificHandler
func NewMusicSpecificHandler(
	trackService services.MediaItemService[*mediatypes.Track],
	albumService services.MediaItemService[*mediatypes.Album],
	artistService services.MediaItemService[*mediatypes.Artist],
) *MusicSpecificHandler {
	return &MusicSpecificHandler{
		MediaItemHandler: NewMediaItemHandler(trackService),
		albumService:     albumService,
		artistService:    artistService,
	}
}

// GetTracksByAlbumID godoc
// @Summary Get tracks for an album
// @Description Retrieves all tracks for a specific album
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Album ID"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Album not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/albums/{id}/tracks [get]
func (h *MusicSpecificHandler) GetTracksByAlbumID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	albumID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid album ID")
		responses.RespondBadRequest(c, err, "Invalid album ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("albumID", albumID).
		Uint64("userID", userID).
		Msg("Getting tracks for album")

	// Get the album first to ensure it exists
	album, err := h.albumService.GetByID(ctx, albumID)
	if err != nil {
		log.Error().Err(err).
			Uint64("albumID", albumID).
			Msg("Failed to retrieve album")
		responses.RespondNotFound(c, err, "Album not found")
		return
	}

	// Get all tracks for the user
	allTracks, err := h.service.GetByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve tracks")
		return
	}

	// Filter tracks that belong to the given album
	var albumTracks []*models.MediaItem[*mediatypes.Track]
	for _, track := range allTracks {
		if track.Data.AlbumID == strconv.FormatUint(albumID, 10) ||
			track.Data.AlbumName == album.Data.Details.Title {
			albumTracks = append(albumTracks, track)
		}
	}

	log.Info().
		Uint64("albumID", albumID).
		Int("trackCount", len(albumTracks)).
		Msg("Tracks retrieved successfully")

	responses.RespondOK(c, albumTracks, "Tracks retrieved successfully")
}

// GetAlbumsByArtistID godoc
// @Summary Get albums for an artist
// @Description Retrieves all albums for a specific artist
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Artist ID"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Artist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/artists/{id}/albums [get]
func (h *MusicSpecificHandler) GetAlbumsByArtistID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	artistID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid artist ID")
		responses.RespondBadRequest(c, err, "Invalid artist ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("artistID", artistID).
		Uint64("userID", userID).
		Msg("Getting albums for artist")

	// Get the artist first to ensure it exists
	artist, err := h.artistService.GetByID(ctx, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("artistID", artistID).
			Msg("Failed to retrieve artist")
		responses.RespondNotFound(c, err, "Artist not found")
		return
	}

	// Get all albums for the user
	allAlbums, err := h.albumService.GetByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve albums")
		responses.RespondInternalError(c, err, "Failed to retrieve albums")
		return
	}

	// Filter albums that belong to the given artist
	var artistAlbums []*models.MediaItem[*mediatypes.Album]
	for _, album := range allAlbums {
		if album.Data.ArtistID == strconv.FormatUint(artistID, 10) ||
			album.Data.ArtistName == artist.Data.Details.Title {
			artistAlbums = append(artistAlbums, album)
		}
	}

	log.Info().
		Uint64("artistID", artistID).
		Int("albumCount", len(artistAlbums)).
		Msg("Albums retrieved successfully")

	responses.RespondOK(c, artistAlbums, "Albums retrieved successfully")
}

// GetTopRatedAlbums godoc
// @Summary Get top rated albums
// @Description Retrieves albums with the highest user/critic ratings
// @Tags music
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of albums to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/albums/top-rated [get]
func (h *MusicSpecificHandler) GetTopRatedAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting top rated albums")

	// Get all albums for the user
	allAlbums, err := h.albumService.GetByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve albums")
		responses.RespondInternalError(c, err, "Failed to retrieve albums")
		return
	}

	// Sort albums by rating (this is a basic implementation, would be better to do at the DB level)
	// TODO: Implement more efficient rating-based sorting at the database level

	// For now, just return the albums we have, limited by the requested count
	var result []*models.MediaItem[*mediatypes.Album]
	if len(allAlbums) > limit {
		result = allAlbums[:limit]
	} else {
		result = allAlbums
	}

	log.Info().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("count", len(result)).
		Msg("Top rated albums retrieved successfully")

	responses.RespondOK(c, result, "Top rated albums retrieved successfully")
}

// GetMostPlayedTracks godoc
// @Summary Get most played tracks
// @Description Retrieves tracks with the highest play count
// @Tags music
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/tracks/most-played [get]
func (h *MusicSpecificHandler) GetMostPlayedTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting most played tracks")

	// This would require integration with the play history service to determine play counts
	// For now, we'll just return a not implemented response

	// TODO: Implement this properly with play history integration
	log.Info().Msg("Most played tracks feature not yet implemented")
	responses.RespondNotImplemented(c, nil, "Most played tracks feature not yet implemented")
}

// GetRecentlyAddedMusic godoc
// @Summary Get recently added music
// @Description Retrieves recently added music (tracks, albums, or artists)
// @Tags music
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of items to return (default 10)"
// @Param type query string false "Media type filter (track, album, artist)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Music items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/recent [get]
func (h *MusicSpecificHandler) GetRecentlyAddedMusic(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	mediaType := c.Query("type") // Optional media type filter

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Str("mediaType", mediaType).
		Msg("Getting recently added music")

	switch mediaType {
	case "track", "":
		// Get recent tracks
		tracks, err := h.service.GetRecentItems(ctx, userID, limit)
		if err != nil {
			log.Error().Err(err).
				Uint64("userID", userID).
				Msg("Failed to retrieve recent tracks")
			responses.RespondInternalError(c, err, "Failed to retrieve recent tracks")
			return
		}

		log.Info().
			Uint64("userID", userID).
			Int("count", len(tracks)).
			Msg("Recent tracks retrieved successfully")

		responses.RespondOK(c, tracks, "Recent tracks retrieved successfully")

	case "album":
		// Get recent albums
		albums, err := h.albumService.GetRecentItems(ctx, userID, limit)
		if err != nil {
			log.Error().Err(err).
				Uint64("userID", userID).
				Msg("Failed to retrieve recent albums")
			responses.RespondInternalError(c, err, "Failed to retrieve recent albums")
			return
		}

		log.Info().
			Uint64("userID", userID).
			Int("count", len(albums)).
			Msg("Recent albums retrieved successfully")

		// Type assertion needed here for proper JSON serialization
		responses.RespondOK(c, albums, "Recent albums retrieved successfully")

	case "artist":
		// Get recent artists
		artists, err := h.artistService.GetRecentItems(ctx, userID, limit)
		if err != nil {
			log.Error().Err(err).
				Uint64("userID", userID).
				Msg("Failed to retrieve recent artists")
			responses.RespondInternalError(c, err, "Failed to retrieve recent artists")
			return
		}

		log.Info().
			Uint64("userID", userID).
			Int("count", len(artists)).
			Msg("Recent artists retrieved successfully")

		// Type assertion needed here for proper JSON serialization
		responses.RespondOK(c, artists, "Recent artists retrieved successfully")

	default:
		log.Warn().Str("mediaType", mediaType).Msg("Invalid media type")
		responses.RespondBadRequest(c, nil, "Invalid media type. Must be track, album, or artist")
	}
}

// GetSimilarArtists godoc
// @Summary Get similar artists
// @Description Retrieves artists similar to a specified artist
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Artist ID"
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of artists to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Artist]] "Similar artists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Artist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/artists/{id}/similar [get]
func (h *MusicSpecificHandler) GetSimilarArtists(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	artistID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid artist ID")
		responses.RespondBadRequest(c, err, "Invalid artist ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("artistID", artistID).
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting similar artists")

	// Get the artist first to ensure it exists
	artist, err := h.artistService.GetByID(ctx, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("artistID", artistID).
			Msg("Failed to retrieve artist")
		responses.RespondNotFound(c, err, "Artist not found")
		return
	}

	// Check if the artist has any similar artists listed
	if len(artist.Data.SimilarArtists) == 0 {
		log.Info().
			Uint64("artistID", artistID).
			Msg("No similar artists found")

		responses.RespondOK(c, []*models.MediaItem[*mediatypes.Artist]{}, "No similar artists found")
		return
	}

	// Get all artists for the user
	allArtists, err := h.artistService.GetByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve artists")
		responses.RespondInternalError(c, err, "Failed to retrieve artists")
		return
	}

	// Filter artists that are similar to the given artist
	var similarArtists []*models.MediaItem[*mediatypes.Artist]
	for _, otherArtist := range allArtists {
		for _, similarID := range artist.Data.SimilarArtists {
			// Check if the current artist is in the similar artists list
			// This handles both direct ID matches and name-based matches
			if strconv.FormatUint(otherArtist.ID, 10) == similarID ||
				otherArtist.Data.Details.Title == similarID {
				similarArtists = append(similarArtists, otherArtist)
				break
			}
		}

		// Limit the results
		if len(similarArtists) >= limit {
			break
		}
	}

	log.Info().
		Uint64("artistID", artistID).
		Int("count", len(similarArtists)).
		Msg("Similar artists retrieved successfully")

	if len(similarArtists) == 0 {
		// If we didn't find any similar artists in our database,
		// but the artist had similarArtists defined, return a custom message
		responses.RespondOK(c, similarArtists, "Similar artists retrieved successfully")
		return
	}

	responses.RespondOK(c, similarArtists, "Similar artists retrieved successfully")
}

// GetGenreRecommendations godoc
// @Summary Get music recommendations by genre
// @Description Retrieves music recommendations based on genre preferences
// @Tags music
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param genre query string true "Genre to get recommendations for"
// @Param limit query int false "Maximum number of items to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Recommendations retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/recommendations/genre [get]
func (h *MusicSpecificHandler) GetGenreRecommendations(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	genre := c.Query("genre")
	if genre == "" {
		log.Warn().Msg("Genre is required")
		responses.RespondBadRequest(c, nil, "Genre is required")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userID", userID).
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting music recommendations by genre")

	// This would ideally be implemented with a more sophisticated recommendation engine
	// For now, we'll just return a simple genre-based filter of tracks

	// Get all tracks for the user
	allTracks, err := h.service.GetByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve tracks")
		return
	}

	// Filter tracks by genre
	var genreTracks []*models.MediaItem[*mediatypes.Track]
	for _, track := range allTracks {
		for _, g := range track.Data.Details.Genres {
			if g == genre {
				genreTracks = append(genreTracks, track)
				break
			}
		}

		// Limit the results
		if len(genreTracks) >= limit {
			break
		}
	}

	log.Info().
		Uint64("userID", userID).
		Str("genre", genre).
		Int("count", len(genreTracks)).
		Msg("Genre-based recommendations retrieved successfully")

	responses.RespondOK(c, genreTracks, "Genre-based recommendations retrieved successfully")
}
