package providers

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"suasor/client/media/types"
	"suasor/repository"
)

// DB provides the database connection
func ProvideDB() (*gorm.DB, error) {
	// In a real implementation, this would initialize an actual database connection
	// For this example, we're just returning nil and no error
	return nil, nil
}

// --- Core Media Item Repositories ---

func ProvideMediaItemRepository[T types.MediaData](db *gorm.DB) repository.MediaItemRepository[T] {
	return repository.NewMediaItemRepository[T](db)
}

func ProvideMovieRepository(db *gorm.DB) repository.MediaItemRepository[*types.Movie] {
	return repository.NewMediaItemRepository[*types.Movie](db)
}

func ProvideSeriesRepository(db *gorm.DB) repository.MediaItemRepository[*types.Series] {
	return repository.NewMediaItemRepository[*types.Series](db)
}

func ProvideSeasonRepository(db *gorm.DB) repository.MediaItemRepository[*types.Season] {
	return repository.NewMediaItemRepository[*types.Season](db)
}

func ProvideEpisodeRepository(db *gorm.DB) repository.MediaItemRepository[*types.Episode] {
	return repository.NewMediaItemRepository[*types.Episode](db)
}

func ProvideTrackRepository(db *gorm.DB) repository.MediaItemRepository[*types.Track] {
	return repository.NewMediaItemRepository[*types.Track](db)
}

func ProvideAlbumRepository(db *gorm.DB) repository.MediaItemRepository[*types.Album] {
	return repository.NewMediaItemRepository[*types.Album](db)
}

func ProvideArtistRepository(db *gorm.DB) repository.MediaItemRepository[*types.Artist] {
	return repository.NewMediaItemRepository[*types.Artist](db)
}

func ProvidePlaylistRepository(db *gorm.DB) repository.MediaItemRepository[*types.Playlist] {
	return repository.NewMediaItemRepository[*types.Playlist](db)
}

func ProvideCollectionRepository(db *gorm.DB) repository.MediaItemRepository[*types.Collection] {
	return repository.NewMediaItemRepository[*types.Collection](db)
}

// --- User Media Item Repositories ---

func ProvideUserMediaItemRepository[T types.MediaData](db *gorm.DB) repository.UserMediaItemRepository[T] {
	return repository.NewUserMediaItemRepository[T](db)
}

func ProvideUserMovieRepository(db *gorm.DB) repository.UserMediaItemRepository[*types.Movie] {
	return repository.NewUserMediaItemRepository[*types.Movie](db)
}

func ProvideUserSeriesRepository(db *gorm.DB) repository.UserMediaItemRepository[*types.Series] {
	return repository.NewUserMediaItemRepository[*types.Series](db)
}

func ProvideUserSeasonRepository(db *gorm.DB) repository.UserMediaItemRepository[*types.Season] {
	return repository.NewUserMediaItemRepository[*types.Season](db)
}

func ProvideUserEpisodeRepository(db *gorm.DB) repository.UserMediaItemRepository[*types.Episode] {
	return repository.NewUserMediaItemRepository[*types.Episode](db)
}

func ProvideUserTrackRepository(db *gorm.DB) repository.UserMediaItemRepository[*types.Track] {
	return repository.NewUserMediaItemRepository[*types.Track](db)
}

func ProvideUserAlbumRepository(db *gorm.DB) repository.UserMediaItemRepository[*types.Album] {
	return repository.NewUserMediaItemRepository[*types.Album](db)
}

func ProvideUserArtistRepository(db *gorm.DB) repository.UserMediaItemRepository[*types.Artist] {
	return repository.NewUserMediaItemRepository[*types.Artist](db)
}

func ProvideUserPlaylistRepository(db *gorm.DB) repository.UserMediaItemRepository[*types.Playlist] {
	return repository.NewUserMediaItemRepository[*types.Playlist](db)
}

func ProvideUserCollectionRepository(db *gorm.DB) repository.UserMediaItemRepository[*types.Collection] {
	return repository.NewUserMediaItemRepository[*types.Collection](db)
}

// --- Core User Media Item Data Repositories ---

func ProvideCoreUserMediaItemDataRepository[T types.MediaData](db *gorm.DB) repository.CoreUserMediaItemDataRepository[T] {
	return repository.NewCoreUserMediaItemDataRepository[T](db)
}

func ProvideCoreMovieDataRepository(db *gorm.DB) repository.CoreUserMediaItemDataRepository[*types.Movie] {
	return repository.NewCoreUserMediaItemDataRepository[*types.Movie](db)
}

func ProvideCoreSeriesDataRepository(db *gorm.DB) repository.CoreUserMediaItemDataRepository[*types.Series] {
	return repository.NewCoreUserMediaItemDataRepository[*types.Series](db)
}

func ProvideCoreSeasonDataRepository(db *gorm.DB) repository.CoreUserMediaItemDataRepository[*types.Season] {
	return repository.NewCoreUserMediaItemDataRepository[*types.Season](db)
}

func ProvideCoreEpisodeDataRepository(db *gorm.DB) repository.CoreUserMediaItemDataRepository[*types.Episode] {
	return repository.NewCoreUserMediaItemDataRepository[*types.Episode](db)
}

func ProvideCoreTrackDataRepository(db *gorm.DB) repository.CoreUserMediaItemDataRepository[*types.Track] {
	return repository.NewCoreUserMediaItemDataRepository[*types.Track](db)
}

func ProvideCoreAlbumDataRepository(db *gorm.DB) repository.CoreUserMediaItemDataRepository[*types.Album] {
	return repository.NewCoreUserMediaItemDataRepository[*types.Album](db)
}

func ProvideCoreArtistDataRepository(db *gorm.DB) repository.CoreUserMediaItemDataRepository[*types.Artist] {
	return repository.NewCoreUserMediaItemDataRepository[*types.Artist](db)
}

func ProvideCorePlaylistDataRepository(db *gorm.DB) repository.CoreUserMediaItemDataRepository[*types.Playlist] {
	return repository.NewCoreUserMediaItemDataRepository[*types.Playlist](db)
}

func ProvideCoreCollectionDataRepository(db *gorm.DB) repository.CoreUserMediaItemDataRepository[*types.Collection] {
	return repository.NewCoreUserMediaItemDataRepository[*types.Collection](db)
}

// --- User Media Item Data Repositories ---

func ProvideUserMediaItemDataRepository[T types.MediaData](
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[T],
) repository.UserMediaItemDataRepository[T] {
	return repository.NewUserMediaItemDataRepository[T](db, coreRepo)
}

func ProvideUserMovieDataRepository(
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Movie],
) repository.UserMediaItemDataRepository[*types.Movie] {
	return repository.NewUserMediaItemDataRepository[*types.Movie](db, coreRepo)
}

func ProvideUserSeriesDataRepository(
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Series],
) repository.UserMediaItemDataRepository[*types.Series] {
	return repository.NewUserMediaItemDataRepository[*types.Series](db, coreRepo)
}

func ProvideUserSeasonDataRepository(
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Season],
) repository.UserMediaItemDataRepository[*types.Season] {
	return repository.NewUserMediaItemDataRepository[*types.Season](db, coreRepo)
}

func ProvideUserEpisodeDataRepository(
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Episode],
) repository.UserMediaItemDataRepository[*types.Episode] {
	return repository.NewUserMediaItemDataRepository[*types.Episode](db, coreRepo)
}

func ProvideUserTrackDataRepository(
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Track],
) repository.UserMediaItemDataRepository[*types.Track] {
	return repository.NewUserMediaItemDataRepository[*types.Track](db, coreRepo)
}

func ProvideUserAlbumDataRepository(
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Album],
) repository.UserMediaItemDataRepository[*types.Album] {
	return repository.NewUserMediaItemDataRepository[*types.Album](db, coreRepo)
}

func ProvideUserArtistDataRepository(
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Artist],
) repository.UserMediaItemDataRepository[*types.Artist] {
	return repository.NewUserMediaItemDataRepository[*types.Artist](db, coreRepo)
}

func ProvideUserPlaylistDataRepository(
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Playlist],
) repository.UserMediaItemDataRepository[*types.Playlist] {
	return repository.NewUserMediaItemDataRepository[*types.Playlist](db, coreRepo)
}

func ProvideUserCollectionDataRepository(
	db *gorm.DB,
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Collection],
) repository.UserMediaItemDataRepository[*types.Collection] {
	return repository.NewUserMediaItemDataRepository[*types.Collection](db, coreRepo)
}

// --- Repository Provider Sets ---

// CoreRepositorySet is a provider set for core repositories
var CoreRepositorySet = wire.NewSet(
	ProvideMovieRepository,
	ProvideSeriesRepository,
	ProvideSeasonRepository,
	ProvideEpisodeRepository,
	ProvideTrackRepository,
	ProvideAlbumRepository,
	ProvideArtistRepository,
	ProvidePlaylistRepository,
	ProvideCollectionRepository,
)

// UserRepositorySet is a provider set for user repositories
var UserRepositorySet = wire.NewSet(
	ProvideUserMovieRepository,
	ProvideUserSeriesRepository,
	ProvideUserSeasonRepository,
	ProvideUserEpisodeRepository,
	ProvideUserTrackRepository,
	ProvideUserAlbumRepository,
	ProvideUserArtistRepository,
	ProvideUserPlaylistRepository,
	ProvideUserCollectionRepository,
)

// CoreDataRepositorySet is a provider set for core data repositories
var CoreDataRepositorySet = wire.NewSet(
	ProvideCoreMovieDataRepository,
	ProvideCoreSeriesDataRepository,
	ProvideCoreSeasonDataRepository,
	ProvideCoreEpisodeDataRepository,
	ProvideCoreTrackDataRepository,
	ProvideCoreAlbumDataRepository,
	ProvideCoreArtistDataRepository,
	ProvideCorePlaylistDataRepository,
	ProvideCoreCollectionDataRepository,
)

// UserDataRepositorySet is a provider set for user data repositories
var UserDataRepositorySet = wire.NewSet(
	ProvideUserMovieDataRepository,
	ProvideUserSeriesDataRepository,
	ProvideUserSeasonDataRepository,
	ProvideUserEpisodeDataRepository,
	ProvideUserTrackDataRepository,
	ProvideUserAlbumDataRepository,
	ProvideUserArtistDataRepository,
	ProvideUserPlaylistDataRepository,
	ProvideUserCollectionDataRepository,
)

// RepositorySet combines all repository provider sets
var RepositorySet = wire.NewSet(
	ProvideDB,
	CoreRepositorySet,
	UserRepositorySet,
	CoreDataRepositorySet,
	UserDataRepositorySet,
)

