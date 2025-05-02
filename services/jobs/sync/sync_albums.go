package sync

import (
	"context"
	"log"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/types/models"
)

// processAlbumBatchForArtist processes a batch of albums for a specific artist
func (j *MediaSyncJob) processAlbumBatchForArtist(
	ctx context.Context,
	albums []*models.MediaItem[*mediatypes.Album],
	clientArtistID string,
	clientID uint64,
	clientType clienttypes.ClientType) ([]*models.MediaItem[*mediatypes.Album], error) {

	clientMedia, _, err := j.getClientMedia(ctx, clientID)
	if err != nil {
		// Just log the error but continue processing with what we have
		log.Printf("Failed to get media client for album details: %v", err)
	}

	// Cast to music provider if possible
	var musicProvider providers.MusicProvider
	if clientMedia != nil {
		if mp, ok := clientMedia.(providers.MusicProvider); ok {
			musicProvider = mp
		}
	}

	processedAlbums := make([]*models.MediaItem[*mediatypes.Album], 0, len(albums))

	for _, album := range albums {
		albumClientItemID, exists := album.GetClientItemID(clientID)
		if !exists {
			log.Printf("Skipping album with no client IDs: %s", album.GetData().Details.Title)
			continue
		}

		// Check if the album already exists in the database
		existingAlbum, err := j.itemRepos.AlbumUserRepo().GetByClientItemID(ctx, clientID, albumClientItemID)
		if err != nil || existingAlbum == nil {
			existingAlbum, err = j.itemRepos.AlbumUserRepo().GetByExternalIDs(ctx, album.GetData().Details.ExternalIDs)
		}
		if err != nil || existingAlbum == nil {
			musicService, err := j.getMusicService(ctx, clientID, clientType)
			if err != nil {
				log.Printf("Error getting music service for client %d: %v", clientID, err)
				continue
			}
			// Try to find by title and artist
			existingAlbum, err = musicService.GetAlbumByTitleAndArtistName(ctx, album.GetData().Details.Title, album.GetData().ArtistName)
		}

		// Get tracks for this album - note that we get a single MediaItem here
		singleTrack, err := musicProvider.GetMusicTrackByID(ctx, albumClientItemID)
		if err != nil {
			log.Printf("Error getting tracks for album %s: %v", album.GetData().Details.Title, err)
		}

		// Create an empty TrackEntries collection
		trackEntries := make(mediatypes.TrackEntries, 0)

		// Process track if we have one
		if singleTrack != nil {
			// Get the client ID for the track
			trackClientID, hasID := singleTrack.GetClientItemID(clientID)
			if hasID {
				// Check if the track exists in the database
				existingTrack, err := j.itemRepos.TrackUserRepo().GetByClientItemID(ctx, clientID, trackClientID)
				if err != nil || existingTrack == nil {
					// Create a new track
					singleTrack.Title = singleTrack.GetData().Details.Title
					singleTrack.ReleaseDate = singleTrack.GetData().Details.ReleaseDate
					singleTrack.ReleaseYear = singleTrack.GetData().Details.ReleaseYear

					existingTrack, err = j.itemRepos.TrackUserRepo().Create(ctx, singleTrack)
					if err != nil {
						log.Printf("Error creating track: %v", err)
					} else {
						// Add to our track entries
						trackEntries = append(trackEntries, &mediatypes.TrackEntry{
							Number:  existingTrack.GetData().Number,
							Title:   existingTrack.GetData().Details.Title,
							TrackID: existingTrack.ID,
						})
					}
				} else {
					// Add existing track to our entries
					trackEntries = append(trackEntries, &mediatypes.TrackEntry{
						Number:  existingTrack.GetData().Number,
						Title:   existingTrack.GetData().Details.Title,
						TrackID: existingTrack.ID,
					})
				}
			}
		}

		// Add tracks to album
		album.GetData().Tracks = &trackEntries
		album.GetData().TrackCount = len(trackEntries)

		// Update existing album or save new one
		if existingAlbum != nil {
			existingAlbum.Merge(album)

			// Ensure album has tracks
			if existingAlbum.GetData().Tracks == nil || len(*existingAlbum.GetData().Tracks) == 0 {
				existingAlbum.GetData().Tracks = &trackEntries
				existingAlbum.GetData().TrackCount = len(trackEntries)
			}

			updatedAlbum, err := j.itemRepos.AlbumUserRepo().Update(ctx, existingAlbum)
			if err != nil {
				log.Printf("Error updating album: %v", err)
				continue
			}
			processedAlbums = append(processedAlbums, updatedAlbum)
		} else {
			// Set top level title and release fields
			album.Title = album.GetData().Details.Title
			album.ReleaseDate = album.GetData().Details.ReleaseDate
			album.ReleaseYear = album.GetData().Details.ReleaseYear

			newAlbum, err := j.itemRepos.AlbumUserRepo().Create(ctx, album)
			if err != nil {
				log.Printf("Error creating album: %v", err)
				continue
			}
			processedAlbums = append(processedAlbums, newAlbum)
		}
	}

	return processedAlbums, nil
}

// processArtistBatch processes a batch of music artists and saves them to the database
func (j *MediaSyncJob) processArtistBatch(ctx context.Context, artists []*models.MediaItem[*mediatypes.Artist], clientID uint64, clientType clienttypes.ClientType) error {
	// Try to get a music provider for this client to fetch album details
	clientMedia, _, err := j.getClientMedia(ctx, clientID)
	if err != nil {
		// Just log the error but continue processing with what we have
		log.Printf("Failed to get media client for album details: %v", err)
	}

	// Cast to music provider if possible
	var musicProvider providers.MusicProvider
	if clientMedia != nil {
		if mp, ok := clientMedia.(providers.MusicProvider); ok {
			musicProvider = mp
		}
	}

	for _, artist := range artists {
		clientArtistID, exists := artist.GetClientItemID(clientID)
		if !exists {
			log.Printf("Skipping artist with no client IDs: %s", artist.GetData().Details.Title)
			continue
		}

		options := &mediatypes.QueryOptions{
			ItemIDs: clientArtistID,
		}

		// Get Albums and tracks for this artist from the provider
		albums, err := musicProvider.GetMusicAlbums(ctx, options)
		if err != nil {
			log.Printf("Error getting albums for artist %s: %v", artist.GetData().Details.Title, err)
		}

		// Process albums
		processedAlbums, err := j.processAlbumBatchForArtist(ctx, albums, clientArtistID, clientID, clientType)
		if err != nil {
			log.Printf("Error processing album batch: %v", err)
		}

		for _, album := range processedAlbums {
			log.Printf("Album: %s", album.GetData().Details.Title)
			albumID := album.ID
			albumName := album.GetData().Details.Title

			// Get track IDs for this album
			trackIDs := album.Data.GetTrackIDs()

			// Add album tracks to artist
			artist.Data.AddAlbumTrackIDs(albumID, albumName, trackIDs)
		}

		// Check if the artist already exists in the database
		existingArtist, err := j.itemRepos.ArtistUserRepo().GetByClientItemID(ctx, clientID, clientArtistID)

		// Check by external services IDs
		if err != nil || existingArtist == nil {
			existingArtist, err = j.itemRepos.ArtistUserRepo().GetByExternalIDs(ctx, artist.GetData().Details.ExternalIDs)
		}

		// Check by title
		if err != nil || existingArtist == nil {
			log.Printf("Could not find artist '%s' by client ID or external IDs", artist.GetData().Details.Title)
		}

		if err == nil && existingArtist != nil {
			// Artist exists, update it
			existingArtist.Merge(artist)

			// Makes sure Data is merged properly
			if existingArtist.Data != nil && artist.Data != nil {
				existingArtist.Data.Merge(artist.Data)
			} else if existingArtist.Data == nil && artist.Data != nil {
				existingArtist.Data = artist.Data
			}

			updatedArtist, err := j.itemRepos.ArtistUserRepo().Update(ctx, existingArtist)
			j.updateAlbumsTracksArtistIDs(ctx, updatedArtist, clientID, clientType)
			if err != nil {
				log.Printf("Error updating artist: %v", err)
				continue
			}
		} else {
			// Artist doesn't exist, create it
			// Set top level title field
			artist.Title = artist.Data.Details.Title

			// Make sure data is properly initialized
			if artist.Data.Genres == nil {
				artist.Data.Genres = []string{}
			}

			// Create the artist
			createdArtist, err := j.itemRepos.ArtistUserRepo().Create(ctx, artist)
			if err != nil {
				log.Printf("Error creating artist: %v", err)
				continue
			}

			// Update album/track artist IDs
			j.updateAlbumsTracksArtistIDs(ctx, createdArtist, clientID, clientType)
		}
	}

	return nil
}

// processIndependentAlbumBatch processes a batch of albums that are not tied to a specific artist
func (j *MediaSyncJob) processIndependentAlbumBatch(
	ctx context.Context,
	albums []*models.MediaItem[*mediatypes.Album],
	clientID uint64,
	clientType clienttypes.ClientType) ([]*models.MediaItem[*mediatypes.Album], error) {

	clientMedia, _, err := j.getClientMedia(ctx, clientID)
	if err != nil {
		// Just log the error but continue processing with what we have
		log.Printf("Failed to get media client for album details: %v", err)
	}

	// Cast to music provider if possible
	var musicProvider providers.MusicProvider
	if clientMedia != nil {
		if mp, ok := clientMedia.(providers.MusicProvider); ok {
			musicProvider = mp
		}
	}

	processedAlbums := make([]*models.MediaItem[*mediatypes.Album], 0, len(albums))

	for _, album := range albums {
		albumClientItemID, exists := album.GetClientItemID(clientID)
		if !exists {
			log.Printf("Skipping album with no client IDs: %s", album.GetData().Details.Title)
			continue
		}

		// Check if the album already exists in the database
		existingAlbum, err := j.itemRepos.AlbumUserRepo().GetByClientItemID(ctx, clientID, albumClientItemID)
		if err != nil || existingAlbum == nil {
			existingAlbum, err = j.itemRepos.AlbumUserRepo().GetByExternalIDs(ctx, album.GetData().Details.ExternalIDs)
		}
		if err != nil || existingAlbum == nil {
			// Try to find by title and year
			existingAlbum, err = j.itemRepos.AlbumUserRepo().GetByTitleAndYear(ctx, clientID, album.GetData().Details.Title, album.GetData().Details.ReleaseYear)
		}

		// Get all album tracks through the client directly
		// Need to use GetMusic and filter by album ID since GetAlbumTracks isn't available
		var albumTracks []*models.MediaItem[*mediatypes.Track]
		if musicProvider != nil {
			allTracks, err := musicProvider.GetMusicTracks(ctx, &mediatypes.QueryOptions{})
			if err == nil {
				// Filter tracks to only those belonging to this album
				for _, track := range allTracks {
					_, hasID := track.GetClientItemID(clientID)
					if hasID && track.GetData().AlbumID > 0 {
						// If the track has album information, store it
						albumTracks = append(albumTracks, track)
					}
				}
			} else {
				log.Printf("Error getting tracks for album %s: %v", album.GetData().Details.Title, err)
			}
		}

		// Process tracks and create TrackEntries
		trackEntries := make(mediatypes.TrackEntries, 0, len(albumTracks))
		for _, track := range albumTracks {
			// Process track and save to database
			trackClientID, hasID := track.GetClientItemID(clientID)
			if !hasID {
				continue
			}

			savedTrack, err := j.itemRepos.TrackUserRepo().GetByClientItemID(ctx, clientID, trackClientID)
			if err != nil || savedTrack == nil {
				// Set top level title and release fields
				track.Title = track.GetData().Details.Title
				track.ReleaseDate = track.GetData().Details.ReleaseDate
				track.ReleaseYear = track.GetData().Details.ReleaseYear

				// Create the track
				savedTrack, err = j.itemRepos.TrackUserRepo().Create(ctx, track)
				if err != nil {
					log.Printf("Error creating track: %v", err)
					continue
				}
			}

			// Add to track entries
			trackEntries = append(trackEntries, &mediatypes.TrackEntry{
				Number:  savedTrack.GetData().Number,
				Title:   savedTrack.GetData().Details.Title,
				TrackID: savedTrack.ID,
			})
		}

		// Add tracks to album
		album.GetData().Tracks = &trackEntries
		album.GetData().TrackCount = len(trackEntries)

		// Update existing album or save new one
		if existingAlbum != nil {
			existingAlbum.Merge(album)

			// Ensure album has tracks
			if existingAlbum.GetData().Tracks == nil || len(*existingAlbum.GetData().Tracks) == 0 {
				existingAlbum.GetData().Tracks = &trackEntries
				existingAlbum.GetData().TrackCount = len(trackEntries)
			}

			updatedAlbum, err := j.itemRepos.AlbumUserRepo().Update(ctx, existingAlbum)
			if err != nil {
				log.Printf("Error updating album: %v", err)
				continue
			}

			// Update tracks with album references
			j.updateTracksAlbumID(ctx, updatedAlbum, clientID, clientType)

			processedAlbums = append(processedAlbums, updatedAlbum)
		} else {
			// Set top level title and release fields
			album.Title = album.GetData().Details.Title
			album.ReleaseDate = album.GetData().Details.ReleaseDate
			album.ReleaseYear = album.GetData().Details.ReleaseYear

			newAlbum, err := j.itemRepos.AlbumUserRepo().Create(ctx, album)
			if err != nil {
				log.Printf("Error creating album: %v", err)
				continue
			}

			// Update tracks with album references
			j.updateTracksAlbumID(ctx, newAlbum, clientID, clientType)

			processedAlbums = append(processedAlbums, newAlbum)
		}
	}

	return processedAlbums, nil
}
