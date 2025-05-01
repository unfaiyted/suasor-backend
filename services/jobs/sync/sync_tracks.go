package sync

import (
	"context"
	"fmt"
	"log"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/utils/logger"
)

// processTrackBatchForAlbum processes a batch of tracks for a specific album
func (j *MediaSyncJob) processTrackBatchForAlbum(
	ctx context.Context,
	track *models.MediaItem[*mediatypes.Track],
	clientID uint64,
	clientType clienttypes.ClientType) (*mediatypes.TrackEntries, error) {
	
	// Create a new TrackEntries collection
	trackEntries := make(mediatypes.TrackEntries, 0)
	
	// Check if we have a valid track
	if track == nil {
		return &trackEntries, nil
	}
	
	// Get the track's client ID
	trackClientItemID, exists := track.GetClientItemID(clientID)
	if !exists {
		return &trackEntries, fmt.Errorf("track has no client ID")
	}
	
	// Check if track already exists in database
	existingTrack, err := j.itemRepos.TrackUserRepo().GetByClientItemID(ctx, clientID, trackClientItemID)
	if err != nil || existingTrack == nil {
		// Create a new track
		track.Title = track.GetData().Details.Title
		track.ReleaseDate = track.GetData().Details.ReleaseDate
		track.ReleaseYear = track.GetData().Details.ReleaseYear
		
		existingTrack, err = j.itemRepos.TrackUserRepo().Create(ctx, track)
		if err != nil {
			return &trackEntries, fmt.Errorf("failed to create track: %v", err)
		}
	}
	
	// Add the track to our entries
	trackEntries = append(trackEntries, &mediatypes.TrackEntry{
		Number: track.GetData().Number,
		Title: track.GetData().Details.Title,
		TrackID: existingTrack.ID,
	})
	
	return &trackEntries, nil
}

// processIndependentTrackBatch processes a batch of tracks that are not tied to a specific album or artist
func (j *MediaSyncJob) processIndependentTrackBatch(
	ctx context.Context,
	tracks []*models.MediaItem[*mediatypes.Track],
	clientID uint64,
	clientType clienttypes.ClientType) ([]*models.MediaItem[*mediatypes.Track], error) {

	musicService, err := j.getMusicService(ctx, clientID, clientType)
	if err != nil {
		return nil, err
	}
	processedTracks := make([]*models.MediaItem[*mediatypes.Track], 0, len(tracks))

	// Process each track
	for _, track := range tracks {
		trackClientItemID, exists := track.GetClientItemID(clientID)
		if !exists {
			log.Printf("Skipping track with no client IDs: %s", track.GetData().Details.Title)
			continue
		}

		// Check if the track already exists in the database
		existingTrack, err := j.itemRepos.TrackUserRepo().GetByClientItemID(ctx, clientID, trackClientItemID)
		if err != nil || existingTrack == nil {
			existingTrack, err = j.itemRepos.TrackUserRepo().GetByExternalIDs(ctx, track.GetData().Details.ExternalIDs)
		}
		if err != nil || existingTrack == nil {
			title := track.GetData().Details.Title
			artistName := track.GetData().ArtistName
			existingTrack, err = musicService.GetTrackByTitleAndArtistName(ctx, title, artistName)
		}

		// Update existing track or save new one
		if existingTrack != nil {
			existingTrack.Merge(track)

			updatedTrack, err := j.itemRepos.TrackUserRepo().Update(ctx, existingTrack)
			if err != nil {
				log.Printf("Error updating track: %v", err)
				continue
			}
			processedTracks = append(processedTracks, updatedTrack)
		} else {
			// Set top level title and release fields
			track.Title = track.GetData().Details.Title
			track.ReleaseDate = track.GetData().Details.ReleaseDate
			track.ReleaseYear = track.GetData().Details.ReleaseYear

			newTrack, err := j.itemRepos.TrackUserRepo().Create(ctx, track)
			if err != nil {
				log.Printf("Error creating track: %v", err)
				continue
			}
			processedTracks = append(processedTracks, newTrack)
		}
	}

	return processedTracks, nil
}

// updateTracksAlbumID updates tracks with the album ID
func (j *MediaSyncJob) updateTracksAlbumID(
	ctx context.Context,
	album *models.MediaItem[*mediatypes.Album],
	clientID uint64,
	clientType clienttypes.ClientType) {

	log := logger.LoggerFromContext(ctx)

	if album.GetData().Tracks == nil || len(*album.GetData().Tracks) == 0 {
		log.Warn().Msg("Album has no tracks")
		return
	}

	// Update each track to point to this album
	for _, track := range *album.GetData().Tracks {
		if track == nil {
			continue
		}

		// Get the full track record
		trackRecord, err := j.itemRepos.TrackUserRepo().GetByID(ctx, track.TrackID)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("trackID", track.TrackID).
				Uint64("albumID", album.ID).
				Msg("Failed to get track by ID")
			continue
		}

		// Update the album reference
		trackRecord.GetData().AlbumID = album.ID
		trackRecord.GetData().AlbumName = album.GetData().Details.Title

		// Set artist information if available
		if album.GetData().ArtistID != 0 && trackRecord.GetData().ArtistID == 0 {
			trackRecord.GetData().ArtistID = album.GetData().ArtistID
			trackRecord.GetData().ArtistName = album.GetData().ArtistName
		}

		// Save the updated track
		_, err = j.itemRepos.TrackUserRepo().Update(ctx, trackRecord)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("trackID", track.TrackID).
				Uint64("albumID", album.ID).
				Msg("Failed to update track")
		}
	}
}

// updateAlbumsTracksArtistIDs updates album and track records with the correct artist ID
func (j *MediaSyncJob) updateAlbumsTracksArtistIDs(
	ctx context.Context,
	artist *models.MediaItem[*mediatypes.Artist],
	clientID uint64, clientType clienttypes.ClientType) {

	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(clientType)).
		Msg("Updating albums/tracks artistIDs")

	// Guard against no albums
	if len(artist.Data.Albums) == 0 {
		log.Warn().Msg("Artist has no albums")
		return
	}

	// Debug logging of all albums and their IDs
	log.Debug().Msg("Before processing - Album information:")
	for i, album := range artist.Data.Albums {
		log.Debug().
			Int("index", i).
			Uint64("albumID", album.AlbumID).
			Str("albumName", album.AlbumName).
			Int("trackCount", len(album.TrackIDs)).
			Msg("Album info")
	}

	// Loop through each album entry
	for _, album := range artist.Data.Albums {
		// Skip if album ID is not set
		if album.AlbumID == 0 {
			log.Warn().
				Str("albumName", album.AlbumName).
				Msg("Album has no ID, skipping")
			continue
		}

		// Update album record
		albumRecord, err := j.itemRepos.AlbumUserRepo().GetByID(ctx, album.AlbumID)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("albumID", album.AlbumID).
				Str("albumName", album.AlbumName).
				Uint64("clientID", clientID).
				Str("clientType", string(clientType)).
				Msg("Failed to get album by ID")
			continue
		}

		// Set the artist ID and name on the album
		albumRecord.GetData().ArtistID = artist.ID
		albumRecord.GetData().ArtistName = artist.GetData().Details.Title
		_, err = j.itemRepos.AlbumUserRepo().Update(ctx, albumRecord)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("albumID", album.AlbumID).
				Str("albumName", album.AlbumName).
				Uint64("clientID", clientID).
				Str("clientType", string(clientType)).
				Msg("Failed to update album")
			continue
		}

		// Update all tracks for this album
		for _, trackID := range album.TrackIDs {
			track, err := j.itemRepos.TrackUserRepo().GetByID(ctx, trackID)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("trackID", trackID).
					Uint64("albumID", album.AlbumID).
					Uint64("clientID", clientID).
					Str("clientType", string(clientType)).
					Msg("Failed to get track by ID")
				continue
			}

			// Set album and artist IDs on the track
			track.GetData().AlbumID = album.AlbumID
			track.GetData().ArtistID = artist.ID
			track.GetData().ArtistName = artist.GetData().Details.Title
			track.GetData().AlbumName = album.AlbumName

			log.Debug().
				Uint64("trackID", trackID).
				Uint64("albumID", album.AlbumID).
				Uint64("artistID", artist.ID).
				Msg("Setting album and artist IDs for track")

			_, err = j.itemRepos.TrackUserRepo().Update(ctx, track)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("trackID", trackID).
					Uint64("clientID", clientID).
					Str("clientType", string(clientType)).
					Msg("Failed to update track")
				continue
			}
		}
	}
}

func (j *MediaSyncJob) getMusicService(ctx context.Context, clientID uint64, clientType clienttypes.ClientType) (services.ClientMusicService[clienttypes.ClientConfig], error) {
	// Get the client
	client, err := j.getClientMedia(ctx, clientID, clientType)
	if err != nil {
		return nil, err
	}

	// Get the music service for this client
	musicService, ok := client.(services.ClientMusicService[clienttypes.ClientMediaConfig])
	if !ok {
		return nil, fmt.Errorf("client does not implement music service interface")
	}

	return musicService, nil
}
