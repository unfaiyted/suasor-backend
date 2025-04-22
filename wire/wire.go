//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"fmt"
	"time"

	"github.com/google/wire"
	"gorm.io/gorm"
	"suasor/client"
	"suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/database"
	"suasor/handlers"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
	"suasor/services"
	"suasor/services/jobs"
	"suasor/services/jobs/recommendation"
	apptypes "suasor/types"
	"suasor/types/models"
)

// SystemHandlers contains core system handlers
type SystemHandlers struct {
	AuthHandler   *handlers.AuthHandler
	UserHandler   *handlers.UserHandler
	ConfigHandler *handlers.ConfigHandler
	JobHandler    *handlers.JobHandler
	HealthHandler *handlers.HealthHandler
	SearchHandler *handlers.SearchHandler

	UserConfigHandler *handlers.UserConfigHandler
}

// MediaHandlers contains media-related handlers
type MediaHandlers struct {
	// Core Handlers
	MovieHandler   handlers.UserMediaItemHandler[*types.Movie]
	SeriesHandler  handlers.UserMediaItemHandler[*types.Series]
	SeasonHandler  handlers.UserMediaItemHandler[*types.Season]
	EpisodeHandler handlers.UserMediaItemHandler[*types.Episode]
	TrackHandler   handlers.UserMediaItemHandler[*types.Track]
	AlbumHandler   handlers.UserMediaItemHandler[*types.Album]
	ArtistHandler  handlers.UserMediaItemHandler[*types.Artist]

	// List Handlers
	PlaylistHandler   handlers.UserListHandler[*types.Playlist]
	CollectionHandler handlers.UserListHandler[*types.Collection]

	// Special Media Handlers
	PeopleHandler *handlers.PeopleHandler
	CreditHandler *handlers.CreditHandler
}

// MediaDataHandlers contains handlers for media item data
type MediaDataHandlers struct {
	MovieDataHandler      handlers.UserMediaItemDataHandler[*types.Movie]
	SeriesDataHandler     handlers.UserMediaItemDataHandler[*types.Series]
	SeasonDataHandler     handlers.UserMediaItemDataHandler[*types.Season]
	EpisodeDataHandler    handlers.UserMediaItemDataHandler[*types.Episode]
	TrackDataHandler      handlers.UserMediaItemDataHandler[*types.Track]
	AlbumDataHandler      handlers.UserMediaItemDataHandler[*types.Album]
	ArtistDataHandler     handlers.UserMediaItemDataHandler[*types.Artist]
	PlaylistDataHandler   handlers.UserMediaItemDataHandler[*types.Playlist]
	CollectionDataHandler handlers.UserMediaItemDataHandler[*types.Collection]
}

// ClientHandlers contains client-related handlers
type ClientHandlers struct {
	// Master Handler
	ClientsHandler *handlers.ClientsHandler

	// Media Clients
	EmbyHandler     *handlers.ClientHandler[*clienttypes.EmbyConfig]
	JellyfinHandler *handlers.ClientHandler[*clienttypes.JellyfinConfig]
	PlexHandler     *handlers.ClientHandler[*clienttypes.PlexConfig]
	SubsonicHandler *handlers.ClientHandler[*clienttypes.SubsonicConfig]

	// Automation Clients
	RadarrHandler *handlers.ClientHandler[*clienttypes.RadarrConfig]
	SonarrHandler *handlers.ClientHandler[*clienttypes.SonarrConfig]
	LidarrHandler *handlers.ClientHandler[*clienttypes.LidarrConfig]

	// AI Clients - Generic handlers need to be instantiated with specific types
	AIHandler     *handlers.AIHandler[*clienttypes.ClaudeConfig] // We're choosing Claude as a default
	ClaudeHandler *handlers.ClientHandler[*clienttypes.ClaudeConfig]
	OpenAIHandler *handlers.ClientHandler[*clienttypes.OpenAIConfig]
	OllamaHandler *handlers.ClientHandler[*clienttypes.OllamaConfig]

	// Metadata - Generic handler needs to be instantiated with a specific type
	MetadataHandler *handlers.MetadataClientHandler[*clienttypes.TMDBConfig]
}

// SpecializedHandlers contains specialized handlers
type SpecializedHandlers struct {
	RecommendationHandler *handlers.RecommendationHandler

	CalendarHandler *handlers.CalendarHandler
}

// ApplicationHandlers contains all handlers organized by category
type ApplicationHandlers struct {
	System      SystemHandlers
	Media       MediaHandlers
	MediaData   MediaDataHandlers
	Clients     ClientHandlers
	Specialized SpecializedHandlers
}

// ProvideDB provides the database connection
func ProvideDB() (*gorm.DB, error) {
	// In a real implementation, this would load proper configuration from environment or config file
	// This is a simplified version for demonstration purposes
	dbConfig := apptypes.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Name:     "suasor",
	}

	// Initialize database connection
	return database.Initialize(context.Background(), dbConfig)
}

// ----- Repository Providers -----

// Media Item Repositories
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

// User Media Item Repositories
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

// Core Media Item Data Repositories
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

// User Media Item Data Repositories
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

// ----- Service Providers -----

// Core Media Item Services
func ProvideCoreMovieService(
	repo repository.MediaItemRepository[*types.Movie],
) services.CoreMediaItemService[*types.Movie] {
	return services.NewCoreMediaItemService[*types.Movie](repo)
}

func ProvideCoreSeriesService(
	repo repository.MediaItemRepository[*types.Series],
) services.CoreMediaItemService[*types.Series] {
	return services.NewCoreMediaItemService[*types.Series](repo)
}

func ProvideCoreSeasonService(
	repo repository.MediaItemRepository[*types.Season],
) services.CoreMediaItemService[*types.Season] {
	return services.NewCoreMediaItemService[*types.Season](repo)
}

func ProvideCoreEpisodeService(
	repo repository.MediaItemRepository[*types.Episode],
) services.CoreMediaItemService[*types.Episode] {
	return services.NewCoreMediaItemService[*types.Episode](repo)
}

func ProvideCoreTrackService(
	repo repository.MediaItemRepository[*types.Track],
) services.CoreMediaItemService[*types.Track] {
	return services.NewCoreMediaItemService[*types.Track](repo)
}

func ProvideCoreAlbumService(
	repo repository.MediaItemRepository[*types.Album],
) services.CoreMediaItemService[*types.Album] {
	return services.NewCoreMediaItemService[*types.Album](repo)
}

func ProvideCoreArtistService(
	repo repository.MediaItemRepository[*types.Artist],
) services.CoreMediaItemService[*types.Artist] {
	return services.NewCoreMediaItemService[*types.Artist](repo)
}

func ProvideCorePlaylistService(
	repo repository.MediaItemRepository[*types.Playlist],
) services.CoreMediaItemService[*types.Playlist] {
	return services.NewCoreMediaItemService[*types.Playlist](repo)
}

func ProvideCoreCollectionService(
	repo repository.MediaItemRepository[*types.Collection],
) services.CoreMediaItemService[*types.Collection] {
	return services.NewCoreMediaItemService[*types.Collection](repo)
}

// Core List Services
func ProvideCorePlaylistListService(
	repo repository.MediaItemRepository[*types.Playlist],
) services.CoreListService[*types.Playlist] {
	return services.NewCoreListService[*types.Playlist](repo)
}

func ProvideCoreCollectionListService(
	repo repository.MediaItemRepository[*types.Collection],
) services.CoreListService[*types.Collection] {
	return services.NewCoreListService[*types.Collection](repo)
}

// User Media Item Services
func ProvideUserMovieService(
	coreService services.CoreMediaItemService[*types.Movie],
	userRepo repository.UserMediaItemRepository[*types.Movie],
) services.UserMediaItemService[*types.Movie] {
	return services.NewUserMediaItemService[*types.Movie](coreService, userRepo)
}

func ProvideUserSeriesService(
	coreService services.CoreMediaItemService[*types.Series],
	userRepo repository.UserMediaItemRepository[*types.Series],
) services.UserMediaItemService[*types.Series] {
	return services.NewUserMediaItemService[*types.Series](coreService, userRepo)
}

func ProvideUserSeasonService(
	coreService services.CoreMediaItemService[*types.Season],
	userRepo repository.UserMediaItemRepository[*types.Season],
) services.UserMediaItemService[*types.Season] {
	return services.NewUserMediaItemService[*types.Season](coreService, userRepo)
}

func ProvideUserEpisodeService(
	coreService services.CoreMediaItemService[*types.Episode],
	userRepo repository.UserMediaItemRepository[*types.Episode],
) services.UserMediaItemService[*types.Episode] {
	return services.NewUserMediaItemService[*types.Episode](coreService, userRepo)
}

func ProvideUserTrackService(
	coreService services.CoreMediaItemService[*types.Track],
	userRepo repository.UserMediaItemRepository[*types.Track],
) services.UserMediaItemService[*types.Track] {
	return services.NewUserMediaItemService[*types.Track](coreService, userRepo)
}

func ProvideUserAlbumService(
	coreService services.CoreMediaItemService[*types.Album],
	userRepo repository.UserMediaItemRepository[*types.Album],
) services.UserMediaItemService[*types.Album] {
	return services.NewUserMediaItemService[*types.Album](coreService, userRepo)
}

func ProvideUserArtistService(
	coreService services.CoreMediaItemService[*types.Artist],
	userRepo repository.UserMediaItemRepository[*types.Artist],
) services.UserMediaItemService[*types.Artist] {
	return services.NewUserMediaItemService[*types.Artist](coreService, userRepo)
}

func ProvideUserPlaylistService(
	coreService services.CoreMediaItemService[*types.Playlist],
	userRepo repository.UserMediaItemRepository[*types.Playlist],
) services.UserMediaItemService[*types.Playlist] {
	return services.NewUserMediaItemService[*types.Playlist](coreService, userRepo)
}

func ProvideUserCollectionService(
	coreService services.CoreMediaItemService[*types.Collection],
	userRepo repository.UserMediaItemRepository[*types.Collection],
) services.UserMediaItemService[*types.Collection] {
	return services.NewUserMediaItemService[*types.Collection](coreService, userRepo)
}

// User List Services
func ProvideUserPlaylistListService(
	coreService services.CoreListService[*types.Playlist],
	userRepo repository.UserMediaItemRepository[*types.Playlist],
	userDataRepo repository.UserMediaItemDataRepository[*types.Playlist],
) services.UserListService[*types.Playlist] {
	return services.NewUserListService[*types.Playlist](coreService, userRepo, userDataRepo)
}

func ProvideUserCollectionListService(
	coreService services.CoreListService[*types.Collection],
	userRepo repository.UserMediaItemRepository[*types.Collection],
	userDataRepo repository.UserMediaItemDataRepository[*types.Collection],
) services.UserListService[*types.Collection] {
	return services.NewUserListService[*types.Collection](coreService, userRepo, userDataRepo)
}

// Core User Media Item Data Services
func ProvideCoreMovieDataService(
	coreService services.CoreMediaItemService[*types.Movie],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Movie],
) services.CoreUserMediaItemDataService[*types.Movie] {
	return services.NewCoreUserMediaItemDataService[*types.Movie](coreService, coreRepo)
}

func ProvideCoreSeriesDataService(
	coreService services.CoreMediaItemService[*types.Series],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Series],
) services.CoreUserMediaItemDataService[*types.Series] {
	return services.NewCoreUserMediaItemDataService[*types.Series](coreService, coreRepo)
}

func ProvideCoreEpisodeDataService(
	coreService services.CoreMediaItemService[*types.Episode],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Episode],
) services.CoreUserMediaItemDataService[*types.Episode] {
	return services.NewCoreUserMediaItemDataService[*types.Episode](coreService, coreRepo)
}

func ProvideCoreSeasonDataService(
	coreService services.CoreMediaItemService[*types.Season],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Season],
) services.CoreUserMediaItemDataService[*types.Season] {
	return services.NewCoreUserMediaItemDataService[*types.Season](coreService, coreRepo)
}

func ProvideCoreTrackDataService(
	coreService services.CoreMediaItemService[*types.Track],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Track],
) services.CoreUserMediaItemDataService[*types.Track] {
	return services.NewCoreUserMediaItemDataService[*types.Track](coreService, coreRepo)
}

func ProvideCoreAlbumDataService(
	coreService services.CoreMediaItemService[*types.Album],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Album],
) services.CoreUserMediaItemDataService[*types.Album] {
	return services.NewCoreUserMediaItemDataService[*types.Album](coreService, coreRepo)
}

func ProvideCoreArtistDataService(
	coreService services.CoreMediaItemService[*types.Artist],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Artist],
) services.CoreUserMediaItemDataService[*types.Artist] {
	return services.NewCoreUserMediaItemDataService[*types.Artist](coreService, coreRepo)
}

func ProvideCorePlaylistDataService(
	coreService services.CoreMediaItemService[*types.Playlist],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Playlist],
) services.CoreUserMediaItemDataService[*types.Playlist] {
	return services.NewCoreUserMediaItemDataService[*types.Playlist](coreService, coreRepo)
}

func ProvideCoreCollectionDataService(
	coreService services.CoreMediaItemService[*types.Collection],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Collection],
) services.CoreUserMediaItemDataService[*types.Collection] {
	return services.NewCoreUserMediaItemDataService[*types.Collection](coreService, coreRepo)
}

// User Media Item Data Services
func ProvideUserMovieDataService(
	coreService services.CoreUserMediaItemDataService[*types.Movie],
	userRepo repository.UserMediaItemDataRepository[*types.Movie],
) services.UserMediaItemDataService[*types.Movie] {
	return services.NewUserMediaItemDataService[*types.Movie](coreService, userRepo)
}

func ProvideUserSeriesDataService(
	coreService services.CoreUserMediaItemDataService[*types.Series],
	userRepo repository.UserMediaItemDataRepository[*types.Series],
) services.UserMediaItemDataService[*types.Series] {
	return services.NewUserMediaItemDataService[*types.Series](coreService, userRepo)
}

func ProvideUserSeasonDataService(
	coreService services.CoreUserMediaItemDataService[*types.Season],
	userRepo repository.UserMediaItemDataRepository[*types.Season],
) services.UserMediaItemDataService[*types.Season] {
	return services.NewUserMediaItemDataService[*types.Season](coreService, userRepo)
}

func ProvideUserEpisodeDataService(
	coreService services.CoreUserMediaItemDataService[*types.Episode],
	userRepo repository.UserMediaItemDataRepository[*types.Episode],
) services.UserMediaItemDataService[*types.Episode] {
	return services.NewUserMediaItemDataService[*types.Episode](coreService, userRepo)
}

func ProvideUserTrackDataService(
	coreService services.CoreUserMediaItemDataService[*types.Track],
	userRepo repository.UserMediaItemDataRepository[*types.Track],
) services.UserMediaItemDataService[*types.Track] {
	return services.NewUserMediaItemDataService[*types.Track](coreService, userRepo)
}

func ProvideUserAlbumDataService(
	coreService services.CoreUserMediaItemDataService[*types.Album],
	userRepo repository.UserMediaItemDataRepository[*types.Album],
) services.UserMediaItemDataService[*types.Album] {
	return services.NewUserMediaItemDataService[*types.Album](coreService, userRepo)
}

func ProvideUserArtistDataService(
	coreService services.CoreUserMediaItemDataService[*types.Artist],
	userRepo repository.UserMediaItemDataRepository[*types.Artist],
) services.UserMediaItemDataService[*types.Artist] {
	return services.NewUserMediaItemDataService[*types.Artist](coreService, userRepo)
}

func ProvideUserPlaylistDataService(
	coreService services.CoreUserMediaItemDataService[*types.Playlist],
	userRepo repository.UserMediaItemDataRepository[*types.Playlist],
) services.UserMediaItemDataService[*types.Playlist] {
	return services.NewUserMediaItemDataService[*types.Playlist](coreService, userRepo)
}

func ProvideUserCollectionDataService(
	coreService services.CoreUserMediaItemDataService[*types.Collection],
	userRepo repository.UserMediaItemDataRepository[*types.Collection],
) services.UserMediaItemDataService[*types.Collection] {
	return services.NewUserMediaItemDataService[*types.Collection](coreService, userRepo)
}

// ----- Handler Providers -----

// Core Media Item Handlers
func ProvideCoreMovieHandler(
	service services.CoreMediaItemService[*types.Movie],
) handlers.CoreMediaItemHandler[*types.Movie] {
	return handlers.NewCoreMediaItemHandler[*types.Movie](service)
}

func ProvideCoreSeriesHandler(
	service services.CoreMediaItemService[*types.Series],
) handlers.CoreMediaItemHandler[*types.Series] {
	return handlers.NewCoreMediaItemHandler[*types.Series](service)
}

func ProvideCoreSeasonHandler(
	service services.CoreMediaItemService[*types.Season],
) handlers.CoreMediaItemHandler[*types.Season] {
	return handlers.NewCoreMediaItemHandler[*types.Season](service)
}

func ProvideCoreEpisodeHandler(
	service services.CoreMediaItemService[*types.Episode],
) handlers.CoreMediaItemHandler[*types.Episode] {
	return handlers.NewCoreMediaItemHandler[*types.Episode](service)
}

func ProvideCoreTrackHandler(
	service services.CoreMediaItemService[*types.Track],
) handlers.CoreMediaItemHandler[*types.Track] {
	return handlers.NewCoreMediaItemHandler[*types.Track](service)
}

func ProvideCoreAlbumHandler(
	service services.CoreMediaItemService[*types.Album],
) handlers.CoreMediaItemHandler[*types.Album] {
	return handlers.NewCoreMediaItemHandler[*types.Album](service)
}

func ProvideCoreArtistHandler(
	service services.CoreMediaItemService[*types.Artist],
) handlers.CoreMediaItemHandler[*types.Artist] {
	return handlers.NewCoreMediaItemHandler[*types.Artist](service)
}

func ProvideCorePlaylistHandler(
	service services.CoreMediaItemService[*types.Playlist],
) handlers.CoreMediaItemHandler[*types.Playlist] {
	return handlers.NewCoreMediaItemHandler[*types.Playlist](service)
}

func ProvideCoreCollectionHandler(
	service services.CoreMediaItemService[*types.Collection],
) handlers.CoreMediaItemHandler[*types.Collection] {
	return handlers.NewCoreMediaItemHandler[*types.Collection](service)
}

// Core List Handlers
func ProvideCorePlaylistListHandler(
	coreHandler handlers.CoreMediaItemHandler[*types.Playlist],
	listService services.CoreListService[*types.Playlist],
) handlers.CoreListHandler[*types.Playlist] {
	return handlers.NewCoreListHandler[*types.Playlist](coreHandler, listService)
}

func ProvideCoreCollectionListHandler(
	coreHandler handlers.CoreMediaItemHandler[*types.Collection],
	listService services.CoreListService[*types.Collection],
) handlers.CoreListHandler[*types.Collection] {
	return handlers.NewCoreListHandler[*types.Collection](coreHandler, listService)
}

// User Media Item Handlers
func ProvideUserMovieHandler(
	service services.UserMediaItemService[*types.Movie],
) handlers.UserMediaItemHandler[*types.Movie] {
	return handlers.NewUserMediaItemHandler[*types.Movie](service)
}

func ProvideUserSeriesHandler(
	service services.UserMediaItemService[*types.Series],
) handlers.UserMediaItemHandler[*types.Series] {
	return handlers.NewUserMediaItemHandler[*types.Series](service)
}

func ProvideUserSeasonHandler(
	service services.UserMediaItemService[*types.Season],
) handlers.UserMediaItemHandler[*types.Season] {
	return handlers.NewUserMediaItemHandler[*types.Season](service)
}

func ProvideUserEpisodeHandler(
	service services.UserMediaItemService[*types.Episode],
) handlers.UserMediaItemHandler[*types.Episode] {
	return handlers.NewUserMediaItemHandler[*types.Episode](service)
}

func ProvideUserTrackHandler(
	service services.UserMediaItemService[*types.Track],
) handlers.UserMediaItemHandler[*types.Track] {
	return handlers.NewUserMediaItemHandler[*types.Track](service)
}

func ProvideUserAlbumHandler(
	service services.UserMediaItemService[*types.Album],
) handlers.UserMediaItemHandler[*types.Album] {
	return handlers.NewUserMediaItemHandler[*types.Album](service)
}

func ProvideUserArtistHandler(
	service services.UserMediaItemService[*types.Artist],
) handlers.UserMediaItemHandler[*types.Artist] {
	return handlers.NewUserMediaItemHandler[*types.Artist](service)
}

func ProvideUserPlaylistHandler(
	service services.UserMediaItemService[*types.Playlist],
) handlers.UserMediaItemHandler[*types.Playlist] {
	return handlers.NewUserMediaItemHandler[*types.Playlist](service)
}

func ProvideUserCollectionHandler(
	service services.UserMediaItemService[*types.Collection],
) handlers.UserMediaItemHandler[*types.Collection] {
	return handlers.NewUserMediaItemHandler[*types.Collection](service)
}

// User List Handlers
func ProvideUserPlaylistListHandler(
	coreHandler handlers.CoreListHandler[*types.Playlist],
	itemService services.UserMediaItemService[*types.Playlist],
	listService services.UserListService[*types.Playlist],
) handlers.UserListHandler[*types.Playlist] {
	return handlers.NewUserListHandler[*types.Playlist](coreHandler, itemService, listService)
}

func ProvideUserCollectionListHandler(
	coreHandler handlers.CoreListHandler[*types.Collection],
	itemService services.UserMediaItemService[*types.Collection],
	listService services.UserListService[*types.Collection],
) handlers.UserListHandler[*types.Collection] {
	return handlers.NewUserListHandler[*types.Collection](coreHandler, itemService, listService)
}

// Core User Media Item Data Handlers
func ProvideCoreMovieDataHandler(
	service services.CoreUserMediaItemDataService[*types.Movie],
) handlers.CoreUserMediaItemDataHandler[*types.Movie] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Movie](service)
}

func ProvideCoreSeriesDataHandler(
	service services.CoreUserMediaItemDataService[*types.Series],
) handlers.CoreUserMediaItemDataHandler[*types.Series] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Series](service)
}

func ProvideCoreEpisodeDataHandler(
	service services.CoreUserMediaItemDataService[*types.Episode],
) handlers.CoreUserMediaItemDataHandler[*types.Episode] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Episode](service)
}

func ProvideCoreTrackDataHandler(
	service services.CoreUserMediaItemDataService[*types.Track],
) handlers.CoreUserMediaItemDataHandler[*types.Track] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Track](service)
}

func ProvideCoreSeasonDataHandler(
	service services.CoreUserMediaItemDataService[*types.Season],
) handlers.CoreUserMediaItemDataHandler[*types.Season] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Season](service)
}

func ProvideCoreAlbumDataHandler(
	service services.CoreUserMediaItemDataService[*types.Album],
) handlers.CoreUserMediaItemDataHandler[*types.Album] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Album](service)
}

func ProvideCoreArtistDataHandler(
	service services.CoreUserMediaItemDataService[*types.Artist],
) handlers.CoreUserMediaItemDataHandler[*types.Artist] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Artist](service)
}

func ProvideCorePlaylistDataHandler(
	service services.CoreUserMediaItemDataService[*types.Playlist],
) handlers.CoreUserMediaItemDataHandler[*types.Playlist] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Playlist](service)
}

func ProvideCoreCollectionDataHandler(
	service services.CoreUserMediaItemDataService[*types.Collection],
) handlers.CoreUserMediaItemDataHandler[*types.Collection] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Collection](service)
}

// User Media Item Data Handlers
func ProvideUserMovieDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Movie],
	service services.UserMediaItemDataService[*types.Movie],
) handlers.UserMediaItemDataHandler[*types.Movie] {
	return handlers.NewUserMediaItemDataHandler[*types.Movie](coreHandler, service)
}

func ProvideUserSeriesDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Series],
	service services.UserMediaItemDataService[*types.Series],
) handlers.UserMediaItemDataHandler[*types.Series] {
	return handlers.NewUserMediaItemDataHandler[*types.Series](coreHandler, service)
}

func ProvideUserSeasonDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Season],
	service services.UserMediaItemDataService[*types.Season],
) handlers.UserMediaItemDataHandler[*types.Season] {
	return handlers.NewUserMediaItemDataHandler[*types.Season](coreHandler, service)
}

func ProvideUserEpisodeDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Episode],
	service services.UserMediaItemDataService[*types.Episode],
) handlers.UserMediaItemDataHandler[*types.Episode] {
	return handlers.NewUserMediaItemDataHandler[*types.Episode](coreHandler, service)
}

func ProvideUserTrackDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Track],
	service services.UserMediaItemDataService[*types.Track],
) handlers.UserMediaItemDataHandler[*types.Track] {
	return handlers.NewUserMediaItemDataHandler[*types.Track](coreHandler, service)
}

func ProvideUserAlbumDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Album],
	service services.UserMediaItemDataService[*types.Album],
) handlers.UserMediaItemDataHandler[*types.Album] {
	return handlers.NewUserMediaItemDataHandler[*types.Album](coreHandler, service)
}

func ProvideUserArtistDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Artist],
	service services.UserMediaItemDataService[*types.Artist],
) handlers.UserMediaItemDataHandler[*types.Artist] {
	return handlers.NewUserMediaItemDataHandler[*types.Artist](coreHandler, service)
}

func ProvideUserPlaylistDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Playlist],
	service services.UserMediaItemDataService[*types.Playlist],
) handlers.UserMediaItemDataHandler[*types.Playlist] {
	return handlers.NewUserMediaItemDataHandler[*types.Playlist](coreHandler, service)
}

func ProvideUserCollectionDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Collection],
	service services.UserMediaItemDataService[*types.Collection],
) handlers.UserMediaItemDataHandler[*types.Collection] {
	return handlers.NewUserMediaItemDataHandler[*types.Collection](coreHandler, service)
}

// ----- System Repository Providers -----

// ProvideUserRepository provides a new UserRepository
func ProvideUserRepository(db *gorm.DB) repository.UserRepository {
	return repository.NewUserRepository(db)
}

// ProvideSessionRepository provides a new SessionRepository
func ProvideSessionRepository(db *gorm.DB) repository.SessionRepository {
	return repository.NewSessionRepository(db)
}

// ProvideConfigRepository provides a new ConfigRepository
func ProvideConfigRepository(db *gorm.DB) repository.ConfigRepository {
	return repository.NewConfigRepository()
}

// ProvideJobRepository provides a new JobRepository
func ProvideJobRepository(db *gorm.DB) repository.JobRepository {
	return repository.NewJobRepository(db)
}

// ProvideUserConfigRepository provides a new UserConfigRepository
func ProvideUserConfigRepository(db *gorm.DB) repository.UserConfigRepository {
	return repository.NewUserConfigRepository(db)
}

// ----- System Service Providers -----

// ProvideAuthService provides a new AuthService
func ProvideAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
) services.AuthService {
	// In a real implementation, these would come from configuration
	return services.NewAuthService(
		userRepo,
		sessionRepo,
		"your-jwt-secret",    // JWT secret
		24*time.Hour,         // Access token expiry (24 hours)
		7*24*time.Hour,       // Refresh token expiry (7 days)
		"suasor",             // Token issuer
		"suasor-application", // Token audience
	)
}

// ProvideUserService provides a new UserService
func ProvideUserService(
	userRepo repository.UserRepository,
) services.UserService {
	return services.NewUserService(userRepo)
}

// ProvideConfigService provides a new ConfigService
func ProvideConfigService(
	configRepo repository.ConfigRepository,
) services.ConfigService {
	return services.NewConfigService(configRepo)
}

// ProvideHealthService provides a new HealthService
func ProvideHealthService(db *gorm.DB) services.HealthService {
	return services.NewHealthService(db)
}

// ProvideSearchRepository provides a new SearchRepository
func ProvideSearchRepository(db *gorm.DB) repository.SearchRepository {
	return repository.NewSearchRepository(db)
}

// ProvideEmbyClientRepository provides a ClientRepository for EmbyConfig
func ProvideEmbyClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.EmbyConfig] {
	return repository.NewClientRepository[*clienttypes.EmbyConfig](db)
}

// ProvideJellyfinClientRepository provides a ClientRepository for JellyfinConfig
func ProvideJellyfinClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.JellyfinConfig] {
	return repository.NewClientRepository[*clienttypes.JellyfinConfig](db)
}

// ProvidePlexClientRepository provides a ClientRepository for PlexConfig
func ProvidePlexClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.PlexConfig] {
	return repository.NewClientRepository[*clienttypes.PlexConfig](db)
}

// ProvideSubsonicClientRepository provides a ClientRepository for SubsonicConfig
func ProvideSubsonicClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.SubsonicConfig] {
	return repository.NewClientRepository[*clienttypes.SubsonicConfig](db)
}

// ProvideSonarrClientRepository provides a ClientRepository for SonarrConfig
func ProvideSonarrClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.SonarrConfig] {
	return repository.NewClientRepository[*clienttypes.SonarrConfig](db)
}

// ProvideRadarrClientRepository provides a ClientRepository for RadarrConfig
func ProvideRadarrClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.RadarrConfig] {
	return repository.NewClientRepository[*clienttypes.RadarrConfig](db)
}

// ProvideLidarrClientRepository provides a ClientRepository for LidarrConfig
func ProvideLidarrClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.LidarrConfig] {
	return repository.NewClientRepository[*clienttypes.LidarrConfig](db)
}

// ProvideClaudeClientRepository provides a ClientRepository for ClaudeConfig
func ProvideClaudeClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.ClaudeConfig] {
	return repository.NewClientRepository[*clienttypes.ClaudeConfig](db)
}

// ProvideOpenAIClientRepository provides a ClientRepository for OpenAIConfig
func ProvideOpenAIClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.OpenAIConfig] {
	return repository.NewClientRepository[*clienttypes.OpenAIConfig](db)
}

// ProvideOllamaClientRepository provides a ClientRepository for OllamaConfig
func ProvideOllamaClientRepository(db *gorm.DB) repository.ClientRepository[*clienttypes.OllamaConfig] {
	return repository.NewClientRepository[*clienttypes.OllamaConfig](db)
}

// ProvideMetadataClientService provides a MetadataClientService for TMDBConfig
func ProvideMetadataClientService(factory *client.ClientFactoryService, db *gorm.DB) *services.MetadataClientService[*clienttypes.TMDBConfig] {
	repo := repository.NewClientRepository[*clienttypes.TMDBConfig](db)
	return services.NewMetadataClientService(factory, repo)
}

// ProvideClientRepositories creates a simplified ClientRepositories instance with essential repositories
func ProvideClientRepositories(
	embyRepo repository.ClientRepository[*clienttypes.EmbyConfig],
	jellyfinRepo repository.ClientRepository[*clienttypes.JellyfinConfig],
	plexRepo repository.ClientRepository[*clienttypes.PlexConfig],
	claudeRepo repository.ClientRepository[*clienttypes.ClaudeConfig],
) repobundles.ClientRepositories {
	// Creating null repositories for the ones we're not providing
	db, _ := gorm.Open(nil, &gorm.Config{}) // This creates a null DB for the empty repositories
	subsonicRepo := repository.NewClientRepository[*clienttypes.SubsonicConfig](db)
	sonarrRepo := repository.NewClientRepository[*clienttypes.SonarrConfig](db)
	radarrRepo := repository.NewClientRepository[*clienttypes.RadarrConfig](db)
	lidarrRepo := repository.NewClientRepository[*clienttypes.LidarrConfig](db)
	openaiRepo := repository.NewClientRepository[*clienttypes.OpenAIConfig](db)
	ollamaRepo := repository.NewClientRepository[*clienttypes.OllamaConfig](db)
	
	return repobundles.NewClientRepositories(
		embyRepo,
		jellyfinRepo,
		plexRepo,
		subsonicRepo,
		sonarrRepo,
		radarrRepo,
		lidarrRepo,
		claudeRepo,
		openaiRepo,
		ollamaRepo,
	)
}

// ProvideCoreMediaItemRepositories creates a CoreMediaItemRepositories instance
func ProvideCoreMediaItemRepositories(db *gorm.DB) repobundles.CoreMediaItemRepositories {
	return repobundles.NewCoreMediaItemRepositories(
		repository.NewMediaItemRepository[*types.Movie](db),
		repository.NewMediaItemRepository[*types.Series](db),
		repository.NewMediaItemRepository[*types.Season](db),
		repository.NewMediaItemRepository[*types.Episode](db),
		repository.NewMediaItemRepository[*types.Track](db),
		repository.NewMediaItemRepository[*types.Album](db),
		repository.NewMediaItemRepository[*types.Artist](db),
		repository.NewMediaItemRepository[*types.Collection](db),
		repository.NewMediaItemRepository[*types.Playlist](db),
	)
}

// ProvideClientFactoryService provides a ClientFactoryService
func ProvideClientFactoryService() *client.ClientFactoryService {
	return client.GetClientFactoryService()
}

// ProvideSearchService provides a SearchService with all required dependencies
func ProvideSearchService(
	searchRepo repository.SearchRepository,
	personRepo repository.PersonRepository,
	clientRepos repobundles.ClientRepositories,
	itemRepos repobundles.CoreMediaItemRepositories,
	clientFactoryService *client.ClientFactoryService,
) services.SearchService {
	// Create a proper SearchService implementation
	return services.NewSearchService(
		searchRepo,
		clientRepos,
		itemRepos,
		personRepo,
		clientFactoryService,
	)
}

// ProvidePersonRepository provides a new PersonRepository
func ProvidePersonRepository(db *gorm.DB) repository.PersonRepository {
	return repository.NewPersonRepository(db)
}

// ProvideCreditRepository provides a new CreditRepository
func ProvideCreditRepository(db *gorm.DB) repository.CreditRepository {
	return repository.NewCreditRepository(db)
}

// ProvidePersonService provides a new PersonService
func ProvidePersonService(
	personRepo repository.PersonRepository,
	creditRepo repository.CreditRepository,
) *services.PersonService {
	return services.NewPersonService(personRepo, creditRepo)
}

// ProvideRecommendationJob provides a RecommendationJob instance
func ProvideRecommendationJob() *recommendation.RecommendationJob {
	// This is a simplified version for demonstration purposes
	return &recommendation.RecommendationJob{}
}

// ProvideMediaSyncJob provides a MediaSyncJob instance
func ProvideMediaSyncJob() *jobs.MediaSyncJob {
	// This is a simplified version for demonstration purposes
	return &jobs.MediaSyncJob{}
}

// ProvideWatchHistorySyncJob provides a WatchHistorySyncJob instance
func ProvideWatchHistorySyncJob() *jobs.WatchHistorySyncJob {
	// This is a simplified version for demonstration purposes
	return &jobs.WatchHistorySyncJob{}
}

// No longer needed in our simplified version

// ProvideJobService provides a simplified JobService
func ProvideJobService(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	userConfigRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*types.Movie],
	seriesRepo repository.MediaItemRepository[*types.Series],
	userMovieDataRepo repository.UserMediaItemDataRepository[*types.Movie],
	userSeriesDataRepo repository.UserMediaItemDataRepository[*types.Series],
	recommendationJob *recommendation.RecommendationJob,
	mediaSyncJob *jobs.MediaSyncJob,
	watchHistorySyncJob *jobs.WatchHistorySyncJob,
) services.JobService {
	// Create null repositories and dependencies for the ones we removed
	db, _ := gorm.Open(nil, &gorm.Config{})
	trackRepo := repository.NewMediaItemRepository[*types.Track](db)
	userTrackDataRepo := repository.NewUserMediaItemDataRepository[*types.Track](
		db,
		repository.NewCoreUserMediaItemDataRepository[*types.Track](db),
	)
	favoritesSyncJob := &jobs.FavoritesSyncJob{} // Empty stub
	
	return services.NewJobService(
		jobRepo,             // jobRepo
		userRepo,            // userRepo
		userConfigRepo,      // configRepo
		movieRepo,           // movieRepo
		seriesRepo,          // seriesRepo
		trackRepo,           // musicRepo
		userMovieDataRepo,   // userMovieDataRepo
		userSeriesDataRepo,  // userSeriesDataRepo
		userTrackDataRepo,   // userMusicDataRepo
		recommendationJob,   // recommendationJob
		mediaSyncJob,        // mediaSyncJob
		watchHistorySyncJob, // watchHistorySyncJob
		favoritesSyncJob,    // favoritesSyncJob
	)
}

// ProvideUserConfigService provides a simplified UserConfigService
func ProvideUserConfigService(
	repo repository.UserConfigRepository,
	jobService services.JobService,
	recommendationJob *recommendation.RecommendationJob,
) services.UserConfigService {
	return services.NewUserConfigService(
		repo,       // userConfigRepo
		jobService, // jobService
		recommendationJob,
	)
}

// ----- Client Media Item Services -----

// SimpleClientMediaItemService is a simplified implementation of ClientMediaItemService
// that doesn't actually perform persistence operations but delegates to a client adaptor
type simpleClientMediaItemService[T clienttypes.ClientMediaConfig, U types.MediaData] struct {
	services.CoreMediaItemService[U] // Embed the core service
}

// NewSimpleClientMediaItemService creates a new simple client media item service
func NewSimpleClientMediaItemService[T clienttypes.ClientMediaConfig, U types.MediaData](
	coreService services.CoreMediaItemService[U],
) services.ClientMediaItemService[T, U] {
	return &simpleClientMediaItemService[T, U]{
		CoreMediaItemService: coreService,
	}
}

// GetByClientID retrieves all media items associated with a specific client
func (s *simpleClientMediaItemService[T, U]) GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[U], error) {
	// This is a simplified implementation that would be replaced with actual client API calls
	return []*models.MediaItem[U]{}, nil
}

// GetByClientItemID retrieves a media item by its client-specific ID
func (s *simpleClientMediaItemService[T, U]) GetByClientItemID(ctx context.Context, itemID string, clientID uint64) (*models.MediaItem[U], error) {
	// This is a simplified implementation that would be replaced with actual client API calls
	return nil, fmt.Errorf("not implemented")
}

// GetByMultipleClients retrieves all media items associated with any of the specified clients
func (s *simpleClientMediaItemService[T, U]) GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[U], error) {
	// This is a simplified implementation that would be replaced with actual client API calls
	return []*models.MediaItem[U]{}, nil
}

// SearchAcrossClients searches for media items across multiple clients
func (s *simpleClientMediaItemService[T, U]) SearchAcrossClients(ctx context.Context, query types.QueryOptions, clientIDs []uint64) (map[uint64][]*models.MediaItem[U], error) {
	// This is a simplified implementation that would be replaced with actual client API calls
	return map[uint64][]*models.MediaItem[U]{}, nil
}

// SyncItemBetweenClients creates or updates a mapping between a media item and a target client
func (s *simpleClientMediaItemService[T, U]) SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error {
	// This is a simplified implementation that would be replaced with actual client API calls
	return nil
}

// DeleteClientItem deletes a client item
func (s *simpleClientMediaItemService[T, U]) DeleteClientItem(ctx context.Context, clientID uint64, itemID string) error {
	// This is a simplified implementation that would be replaced with actual client API calls
	return nil
}

// ProvideEmbyMovieClientService provides a client media item service for Emby movies
func ProvideEmbyMovieClientService(
	coreService services.CoreMediaItemService[*types.Movie],
) services.ClientMediaItemService[*clienttypes.EmbyConfig, *types.Movie] {
	return NewSimpleClientMediaItemService[*clienttypes.EmbyConfig, *types.Movie](coreService)
}

// ProvideEmbySeriesClientService provides a client media item service for Emby series
func ProvideEmbySeriesClientService(
	coreService services.CoreMediaItemService[*types.Series],
) services.ClientMediaItemService[*clienttypes.EmbyConfig, *types.Series] {
	return NewSimpleClientMediaItemService[*clienttypes.EmbyConfig, *types.Series](coreService)
}

// ProvideJellyfinMovieClientService provides a client media item service for Jellyfin movies
func ProvideJellyfinMovieClientService(
	coreService services.CoreMediaItemService[*types.Movie],
) services.ClientMediaItemService[*clienttypes.JellyfinConfig, *types.Movie] {
	return NewSimpleClientMediaItemService[*clienttypes.JellyfinConfig, *types.Movie](coreService)
}

// ProvideJellyfinSeriesClientService provides a client media item service for Jellyfin series
func ProvideJellyfinSeriesClientService(
	coreService services.CoreMediaItemService[*types.Series],
) services.ClientMediaItemService[*clienttypes.JellyfinConfig, *types.Series] {
	return NewSimpleClientMediaItemService[*clienttypes.JellyfinConfig, *types.Series](coreService)
}

// ProvidePlexMovieClientService provides a client media item service for Plex movies
func ProvidePlexMovieClientService(
	coreService services.CoreMediaItemService[*types.Movie],
) services.ClientMediaItemService[*clienttypes.PlexConfig, *types.Movie] {
	return NewSimpleClientMediaItemService[*clienttypes.PlexConfig, *types.Movie](coreService)
}

// ProvidePlexSeriesClientService provides a client media item service for Plex series
func ProvidePlexSeriesClientService(
	coreService services.CoreMediaItemService[*types.Series],
) services.ClientMediaItemService[*clienttypes.PlexConfig, *types.Series] {
	return NewSimpleClientMediaItemService[*clienttypes.PlexConfig, *types.Series](coreService)
}

// ----- System Handler Providers -----

func ProvideAuthHandler(service services.AuthService) *handlers.AuthHandler {
	return handlers.NewAuthHandler(service)
}

func ProvideUserHandler(
	service services.UserService,
	configService services.ConfigService,
) *handlers.UserHandler {
	return handlers.NewUserHandler(service, configService)
}

// ProvideConfigHandler provides a new ConfigHandler
func ProvideConfigHandler(service services.ConfigService) *handlers.ConfigHandler {
	return handlers.NewConfigHandler(service)
}

// ProvideJobHandler provides a new JobHandler
func ProvideJobHandler(service services.JobService) *handlers.JobHandler {
	return handlers.NewJobHandler(service)
}

// ProvideHealthHandler provides a new HealthHandler
func ProvideHealthHandler(service services.HealthService) *handlers.HealthHandler {
	return handlers.NewHealthHandler(service)
}

// ProvideSearchHandler provides a new SearchHandler
func ProvideSearchHandler(service services.SearchService) *handlers.SearchHandler {
	return handlers.NewSearchHandler(service)
}

// ProvidePeopleHandler provides a new PeopleHandler
func ProvidePeopleHandler(service *services.PersonService) *handlers.PeopleHandler {
	return handlers.NewPeopleHandler(service)
}

// ProvideUserConfigHandler provides a new UserConfigHandler
func ProvideUserConfigHandler(service services.UserConfigService) *handlers.UserConfigHandler {
	return handlers.NewUserConfigHandler(service)
}

// ----- Specialized Handler Providers -----

// ProvideRecommendationService provides a RecommendationService
func ProvideRecommendationService(jobRepo repository.JobRepository) services.RecommendationService {
	// Simplified placeholder implementation
	return &recommendationServiceImpl{
		jobRepo: jobRepo,
	}
}

// Simple implementation of RecommendationService
type recommendationServiceImpl struct{
    jobRepo repository.JobRepository
}

// GetRecommendations retrieves recommendations for a user with optional filtering
func (s *recommendationServiceImpl) GetRecommendations(ctx context.Context, userID uint64, mediaType string, limit, offset int) ([]models.Recommendation, error) {
    return []models.Recommendation{}, nil
}

// GetRecommendationByID retrieves a specific recommendation by ID
func (s *recommendationServiceImpl) GetRecommendationByID(ctx context.Context, id uint64) (*models.Recommendation, error) {
    return nil, nil
}

// MarkRecommendationAsViewed marks a recommendation as viewed
func (s *recommendationServiceImpl) MarkRecommendationAsViewed(ctx context.Context, id uint64, userID uint64) error {
    return nil
}

// RateRecommendation sets a user rating for a recommendation
func (s *recommendationServiceImpl) RateRecommendation(ctx context.Context, id uint64, userID uint64, rating float32) error {
    return nil
}

// StoreRecommendations stores recommendations for a user
func (s *recommendationServiceImpl) StoreRecommendations(ctx context.Context, recommendations []*models.Recommendation) error {
    return nil
}

// GetRecentRecommendations retrieves recently created recommendations for a user
func (s *recommendationServiceImpl) GetRecentRecommendations(ctx context.Context, userID uint64, days int, limit int) ([]models.Recommendation, error) {
    return []models.Recommendation{}, nil
}

// GetTopRecommendations retrieves top-scored recommendations for a user
func (s *recommendationServiceImpl) GetTopRecommendations(ctx context.Context, userID uint64, minScore float32, limit int) ([]models.Recommendation, error) {
    return []models.Recommendation{}, nil
}

// GetRecommendationCount returns the count of recommendations for a user
func (s *recommendationServiceImpl) GetRecommendationCount(ctx context.Context, userID uint64, mediaType string) (int, error) {
    return 0, nil
}

// DeleteExpiredRecommendations deletes all expired recommendations
func (s *recommendationServiceImpl) DeleteExpiredRecommendations(ctx context.Context) error {
    return nil
}

// ProvideRecommendationHandler provides a RecommendationHandler
func ProvideRecommendationHandler(service services.RecommendationService) *handlers.RecommendationHandler {
	return handlers.NewRecommendationHandler(service)
}

// ProvideCreditService provides a CreditService
func ProvideCreditService(db *gorm.DB) *services.CreditService {
	// Simplified placeholder implementation
	return &services.CreditService{}
}

// ProvideCreditHandler provides a CreditHandler
func ProvideCreditHandler(service *services.CreditService) *handlers.CreditHandler {
	return handlers.NewCreditHandler(service)
}

// ProvideCalendarHandler provides a CalendarHandler
func ProvideCalendarHandler() *handlers.CalendarHandler {
	// Simplified placeholder implementation
	return &handlers.CalendarHandler{}
}

// ----- Client Handler Providers -----

// ProvideClientsHandler provides a simplified ClientsHandler
func ProvideClientsHandler(
	embyService services.ClientService[*clienttypes.EmbyConfig],
	jellyfinService services.ClientService[*clienttypes.JellyfinConfig],
	plexService services.ClientService[*clienttypes.PlexConfig],
	claudeService services.ClientService[*clienttypes.ClaudeConfig],
) *handlers.ClientsHandler {
	// Create stub services for the ones we're not providing
	clientFactory := client.GetClientFactoryService()
	db, _ := gorm.Open(nil, &gorm.Config{}) // This creates a null DB for the empty repositories
	
	subsonicRepo := repository.NewClientRepository[*clienttypes.SubsonicConfig](db)
	sonarrRepo := repository.NewClientRepository[*clienttypes.SonarrConfig](db)
	radarrRepo := repository.NewClientRepository[*clienttypes.RadarrConfig](db)
	lidarrRepo := repository.NewClientRepository[*clienttypes.LidarrConfig](db)
	openaiRepo := repository.NewClientRepository[*clienttypes.OpenAIConfig](db)
	ollamaRepo := repository.NewClientRepository[*clienttypes.OllamaConfig](db)
	
	subsonicService := services.NewClientService[*clienttypes.SubsonicConfig](clientFactory, subsonicRepo)
	sonarrService := services.NewClientService[*clienttypes.SonarrConfig](clientFactory, sonarrRepo)
	radarrService := services.NewClientService[*clienttypes.RadarrConfig](clientFactory, radarrRepo)
	lidarrService := services.NewClientService[*clienttypes.LidarrConfig](clientFactory, lidarrRepo)
	openaiService := services.NewClientService[*clienttypes.OpenAIConfig](clientFactory, openaiRepo)
	ollamaService := services.NewClientService[*clienttypes.OllamaConfig](clientFactory, ollamaRepo)
	
	return handlers.NewClientsHandler(
		embyService,
		jellyfinService,
		plexService,
		subsonicService,
		sonarrService,
		radarrService,
		lidarrService,
		claudeService,
		openaiService,
		ollamaService,
	)
}

// ProvideClientService provides a ClientService
func ProvideEmbyClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.EmbyConfig]) services.ClientService[*clienttypes.EmbyConfig] {
	return services.NewClientService[*clienttypes.EmbyConfig](factory, repo)
}

func ProvideJellyfinClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.JellyfinConfig]) services.ClientService[*clienttypes.JellyfinConfig] {
	return services.NewClientService[*clienttypes.JellyfinConfig](factory, repo)
}

func ProvidePlexClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.PlexConfig]) services.ClientService[*clienttypes.PlexConfig] {
	return services.NewClientService[*clienttypes.PlexConfig](factory, repo)
}

func ProvideSubsonicClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.SubsonicConfig]) services.ClientService[*clienttypes.SubsonicConfig] {
	return services.NewClientService[*clienttypes.SubsonicConfig](factory, repo)
}

func ProvideSonarrClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.SonarrConfig]) services.ClientService[*clienttypes.SonarrConfig] {
	return services.NewClientService[*clienttypes.SonarrConfig](factory, repo)
}

func ProvideRadarrClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.RadarrConfig]) services.ClientService[*clienttypes.RadarrConfig] {
	return services.NewClientService[*clienttypes.RadarrConfig](factory, repo)
}

func ProvideLidarrClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.LidarrConfig]) services.ClientService[*clienttypes.LidarrConfig] {
	return services.NewClientService[*clienttypes.LidarrConfig](factory, repo)
}

func ProvideClaudeClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.ClaudeConfig]) services.ClientService[*clienttypes.ClaudeConfig] {
	return services.NewClientService[*clienttypes.ClaudeConfig](factory, repo)
}

func ProvideOpenAIClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.OpenAIConfig]) services.ClientService[*clienttypes.OpenAIConfig] {
	return services.NewClientService[*clienttypes.OpenAIConfig](factory, repo)
}

func ProvideOllamaClientService(factory *client.ClientFactoryService, repo repository.ClientRepository[*clienttypes.OllamaConfig]) services.ClientService[*clienttypes.OllamaConfig] {
	return services.NewClientService[*clienttypes.OllamaConfig](factory, repo)
}

// ProvideClientHandler provides a ClientHandler
func ProvideEmbyClientHandler(service services.ClientService[*clienttypes.EmbyConfig]) *handlers.ClientHandler[*clienttypes.EmbyConfig] {
	return handlers.NewClientHandler[*clienttypes.EmbyConfig](service)
}

func ProvideJellyfinClientHandler(service services.ClientService[*clienttypes.JellyfinConfig]) *handlers.ClientHandler[*clienttypes.JellyfinConfig] {
	return handlers.NewClientHandler[*clienttypes.JellyfinConfig](service)
}

func ProvidePlexClientHandler(service services.ClientService[*clienttypes.PlexConfig]) *handlers.ClientHandler[*clienttypes.PlexConfig] {
	return handlers.NewClientHandler[*clienttypes.PlexConfig](service)
}

func ProvideSubsonicClientHandler(service services.ClientService[*clienttypes.SubsonicConfig]) *handlers.ClientHandler[*clienttypes.SubsonicConfig] {
	return handlers.NewClientHandler[*clienttypes.SubsonicConfig](service)
}

func ProvideSonarrClientHandler(service services.ClientService[*clienttypes.SonarrConfig]) *handlers.ClientHandler[*clienttypes.SonarrConfig] {
	return handlers.NewClientHandler[*clienttypes.SonarrConfig](service)
}

func ProvideRadarrClientHandler(service services.ClientService[*clienttypes.RadarrConfig]) *handlers.ClientHandler[*clienttypes.RadarrConfig] {
	return handlers.NewClientHandler[*clienttypes.RadarrConfig](service)
}

func ProvideLidarrClientHandler(service services.ClientService[*clienttypes.LidarrConfig]) *handlers.ClientHandler[*clienttypes.LidarrConfig] {
	return handlers.NewClientHandler[*clienttypes.LidarrConfig](service)
}

func ProvideClaudeClientHandler(service services.ClientService[*clienttypes.ClaudeConfig]) *handlers.ClientHandler[*clienttypes.ClaudeConfig] {
	return handlers.NewClientHandler[*clienttypes.ClaudeConfig](service)
}

func ProvideOpenAIClientHandler(service services.ClientService[*clienttypes.OpenAIConfig]) *handlers.ClientHandler[*clienttypes.OpenAIConfig] {
	return handlers.NewClientHandler[*clienttypes.OpenAIConfig](service)
}

func ProvideOllamaClientHandler(service services.ClientService[*clienttypes.OllamaConfig]) *handlers.ClientHandler[*clienttypes.OllamaConfig] {
	return handlers.NewClientHandler[*clienttypes.OllamaConfig](service)
}

// ProvideAIHandler provides an AIHandler
func ProvideAIHandler(
	factory *client.ClientFactoryService,
	claudeService services.ClientService[*clienttypes.ClaudeConfig],
) *handlers.AIHandler[*clienttypes.ClaudeConfig] {
	// We're using Claude as our default AI client type
	return handlers.NewAIHandler(factory, claudeService)
}

// ProvideMetadataClientHandler provides a MetadataClientHandler
func ProvideMetadataClientHandler(
	tmdbService *services.MetadataClientService[*clienttypes.TMDBConfig],
) *handlers.MetadataClientHandler[*clienttypes.TMDBConfig] {
	return handlers.NewMetadataClientHandler(tmdbService)
}

// ----- Client Media Item Handler Providers -----

// ProvideEmbyMovieHandler provides an EmbyMovieHandler
func ProvideEmbyMovieHandler(
	userHandler handlers.UserMediaItemHandler[*types.Movie],
	clientService services.ClientMediaItemService[*clienttypes.EmbyConfig, *types.Movie],
) handlers.ClientMediaItemHandler[*clienttypes.EmbyConfig, *types.Movie] {
	return handlers.NewClientMediaItemHandler[*clienttypes.EmbyConfig, *types.Movie](
		userHandler,
		clientService,
	)
}

// ProvideEmbySeriesHandler provides an EmbySeriesHandler
func ProvideEmbySeriesHandler(
	userHandler handlers.UserMediaItemHandler[*types.Series],
	clientService services.ClientMediaItemService[*clienttypes.EmbyConfig, *types.Series],
) handlers.ClientMediaItemHandler[*clienttypes.EmbyConfig, *types.Series] {
	return handlers.NewClientMediaItemHandler[*clienttypes.EmbyConfig, *types.Series](
		userHandler,
		clientService,
	)
}

// ProvideJellyfinMovieHandler provides a JellyfinMovieHandler
func ProvideJellyfinMovieHandler(
	userHandler handlers.UserMediaItemHandler[*types.Movie],
	clientService services.ClientMediaItemService[*clienttypes.JellyfinConfig, *types.Movie],
) handlers.ClientMediaItemHandler[*clienttypes.JellyfinConfig, *types.Movie] {
	return handlers.NewClientMediaItemHandler[*clienttypes.JellyfinConfig, *types.Movie](
		userHandler,
		clientService,
	)
}

// ProvideJellyfinSeriesHandler provides a JellyfinSeriesHandler
func ProvideJellyfinSeriesHandler(
	userHandler handlers.UserMediaItemHandler[*types.Series],
	clientService services.ClientMediaItemService[*clienttypes.JellyfinConfig, *types.Series],
) handlers.ClientMediaItemHandler[*clienttypes.JellyfinConfig, *types.Series] {
	return handlers.NewClientMediaItemHandler[*clienttypes.JellyfinConfig, *types.Series](
		userHandler,
		clientService,
	)
}

// ProvidePlexMovieHandler provides a PlexMovieHandler
func ProvidePlexMovieHandler(
	userHandler handlers.UserMediaItemHandler[*types.Movie],
	clientService services.ClientMediaItemService[*clienttypes.PlexConfig, *types.Movie],
) handlers.ClientMediaItemHandler[*clienttypes.PlexConfig, *types.Movie] {
	return handlers.NewClientMediaItemHandler[*clienttypes.PlexConfig, *types.Movie](
		userHandler,
		clientService,
	)
}

// ProvidePlexSeriesHandler provides a PlexSeriesHandler
func ProvidePlexSeriesHandler(
	userHandler handlers.UserMediaItemHandler[*types.Series],
	clientService services.ClientMediaItemService[*clienttypes.PlexConfig, *types.Series],
) handlers.ClientMediaItemHandler[*clienttypes.PlexConfig, *types.Series] {
	return handlers.NewClientMediaItemHandler[*clienttypes.PlexConfig, *types.Series](
		userHandler,
		clientService,
	)
}

// ----- Handler Group Providers -----

// ProvideSystemHandlers creates the system handlers group struct
func ProvideSystemHandlers(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	configHandler *handlers.ConfigHandler,
	jobHandler *handlers.JobHandler,
	healthHandler *handlers.HealthHandler,
	searchHandler *handlers.SearchHandler,
	userConfigHandler *handlers.UserConfigHandler,
) SystemHandlers {
	return SystemHandlers{
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		ConfigHandler:    configHandler,
		JobHandler:       jobHandler,
		HealthHandler:    healthHandler,
		SearchHandler:    searchHandler,
		UserConfigHandler: userConfigHandler,
	}
}

// ProvideMediaHandlers creates a simplified media handlers group struct
func ProvideMediaHandlers(
	movieHandler handlers.UserMediaItemHandler[*types.Movie],
	seriesHandler handlers.UserMediaItemHandler[*types.Series],
	playlistHandler handlers.UserListHandler[*types.Playlist],
	peopleHandler *handlers.PeopleHandler,
	creditHandler *handlers.CreditHandler,
) MediaHandlers {
	// Create default empty handlers for the ones we're not providing
	db, _ := gorm.Open(nil, &gorm.Config{})
	
	// Repositories
	seasonRepo := repository.NewMediaItemRepository[*types.Season](db)
	episodeRepo := repository.NewMediaItemRepository[*types.Episode](db)
	trackRepo := repository.NewMediaItemRepository[*types.Track](db)
	albumRepo := repository.NewMediaItemRepository[*types.Album](db)
	artistRepo := repository.NewMediaItemRepository[*types.Artist](db)
	collectionRepo := repository.NewMediaItemRepository[*types.Collection](db)
	
	// Core Services
	seasonCoreService := services.NewCoreMediaItemService[*types.Season](seasonRepo)
	episodeCoreService := services.NewCoreMediaItemService[*types.Episode](episodeRepo)
	trackCoreService := services.NewCoreMediaItemService[*types.Track](trackRepo)
	albumCoreService := services.NewCoreMediaItemService[*types.Album](albumRepo)
	artistCoreService := services.NewCoreMediaItemService[*types.Artist](artistRepo)
	collectionCoreService := services.NewCoreMediaItemService[*types.Collection](collectionRepo)
	
	// User Repositories
	userSeasonRepo := repository.NewUserMediaItemRepository[*types.Season](db)
	userEpisodeRepo := repository.NewUserMediaItemRepository[*types.Episode](db)
	userTrackRepo := repository.NewUserMediaItemRepository[*types.Track](db)
	userAlbumRepo := repository.NewUserMediaItemRepository[*types.Album](db)
	userArtistRepo := repository.NewUserMediaItemRepository[*types.Artist](db)
	userCollectionRepo := repository.NewUserMediaItemRepository[*types.Collection](db)
	
	// User Services
	userSeasonService := services.NewUserMediaItemService[*types.Season](seasonCoreService, userSeasonRepo)
	userEpisodeService := services.NewUserMediaItemService[*types.Episode](episodeCoreService, userEpisodeRepo)
	userTrackService := services.NewUserMediaItemService[*types.Track](trackCoreService, userTrackRepo)
	userAlbumService := services.NewUserMediaItemService[*types.Album](albumCoreService, userAlbumRepo)
	userArtistService := services.NewUserMediaItemService[*types.Artist](artistCoreService, userArtistRepo)
	
	// Collection List Service
	userCollectionDataRepo := repository.NewUserMediaItemDataRepository[*types.Collection](db, repository.NewCoreUserMediaItemDataRepository[*types.Collection](db))
	collectionListService := services.NewUserListService[*types.Collection](
		services.NewCoreListService[*types.Collection](collectionRepo),
		userCollectionRepo,
		userCollectionDataRepo,
	)
	
	// Handlers
	seasonHandler := handlers.NewUserMediaItemHandler[*types.Season](userSeasonService)
	episodeHandler := handlers.NewUserMediaItemHandler[*types.Episode](userEpisodeService)
	trackHandler := handlers.NewUserMediaItemHandler[*types.Track](userTrackService)
	albumHandler := handlers.NewUserMediaItemHandler[*types.Album](userAlbumService)
	artistHandler := handlers.NewUserMediaItemHandler[*types.Artist](userArtistService)
	
	// Collection Handler
	collectionCoreHandler := handlers.NewCoreMediaItemHandler[*types.Collection](collectionCoreService)
	collectionListCoreHandler := handlers.NewCoreListHandler[*types.Collection](
		collectionCoreHandler,
		services.NewCoreListService[*types.Collection](collectionRepo),
	)
	collectionHandler := handlers.NewUserListHandler[*types.Collection](
		collectionListCoreHandler,
		services.NewUserMediaItemService[*types.Collection](collectionCoreService, userCollectionRepo),
		collectionListService,
	)
	
	return MediaHandlers{
		MovieHandler:      movieHandler,
		SeriesHandler:     seriesHandler,
		SeasonHandler:     seasonHandler,
		EpisodeHandler:    episodeHandler,
		TrackHandler:      trackHandler,
		AlbumHandler:      albumHandler,
		ArtistHandler:     artistHandler,
		PlaylistHandler:   playlistHandler,
		CollectionHandler: collectionHandler,
		PeopleHandler:     peopleHandler,
		CreditHandler:     creditHandler,
	}
}

// ProvideMediaDataHandlers creates a simplified media data handlers group struct
func ProvideMediaDataHandlers(
	movieDataHandler handlers.UserMediaItemDataHandler[*types.Movie],
	seriesDataHandler handlers.UserMediaItemDataHandler[*types.Series],
	playlistDataHandler handlers.UserMediaItemDataHandler[*types.Playlist],
) MediaDataHandlers {
	// Create empty handlers for the ones we're not providing
	db, _ := gorm.Open(nil, &gorm.Config{})
	
	// Repositories and Core Services
	seasonRepo := repository.NewMediaItemRepository[*types.Season](db)
	episodeRepo := repository.NewMediaItemRepository[*types.Episode](db)
	trackRepo := repository.NewMediaItemRepository[*types.Track](db)
	albumRepo := repository.NewMediaItemRepository[*types.Album](db)
	artistRepo := repository.NewMediaItemRepository[*types.Artist](db)
	collectionRepo := repository.NewMediaItemRepository[*types.Collection](db)
	
	seasonCoreService := services.NewCoreMediaItemService[*types.Season](seasonRepo)
	episodeCoreService := services.NewCoreMediaItemService[*types.Episode](episodeRepo)
	trackCoreService := services.NewCoreMediaItemService[*types.Track](trackRepo)
	albumCoreService := services.NewCoreMediaItemService[*types.Album](albumRepo)
	artistCoreService := services.NewCoreMediaItemService[*types.Artist](artistRepo)
	collectionCoreService := services.NewCoreMediaItemService[*types.Collection](collectionRepo)
	
	// Core Data Repositories
	coreSeasonDataRepo := repository.NewCoreUserMediaItemDataRepository[*types.Season](db)
	coreEpisodeDataRepo := repository.NewCoreUserMediaItemDataRepository[*types.Episode](db)
	coreTrackDataRepo := repository.NewCoreUserMediaItemDataRepository[*types.Track](db)
	coreAlbumDataRepo := repository.NewCoreUserMediaItemDataRepository[*types.Album](db)
	coreArtistDataRepo := repository.NewCoreUserMediaItemDataRepository[*types.Artist](db)
	coreCollectionDataRepo := repository.NewCoreUserMediaItemDataRepository[*types.Collection](db)
	
	// Core Data Services
	coreSeasonDataService := services.NewCoreUserMediaItemDataService[*types.Season](seasonCoreService, coreSeasonDataRepo)
	coreEpisodeDataService := services.NewCoreUserMediaItemDataService[*types.Episode](episodeCoreService, coreEpisodeDataRepo)
	coreTrackDataService := services.NewCoreUserMediaItemDataService[*types.Track](trackCoreService, coreTrackDataRepo)
	coreAlbumDataService := services.NewCoreUserMediaItemDataService[*types.Album](albumCoreService, coreAlbumDataRepo)
	coreArtistDataService := services.NewCoreUserMediaItemDataService[*types.Artist](artistCoreService, coreArtistDataRepo)
	coreCollectionDataService := services.NewCoreUserMediaItemDataService[*types.Collection](collectionCoreService, coreCollectionDataRepo)
	
	// Core Data Handlers
	coreSeasonDataHandler := handlers.NewCoreUserMediaItemDataHandler[*types.Season](coreSeasonDataService)
	coreEpisodeDataHandler := handlers.NewCoreUserMediaItemDataHandler[*types.Episode](coreEpisodeDataService)
	coreTrackDataHandler := handlers.NewCoreUserMediaItemDataHandler[*types.Track](coreTrackDataService)
	coreAlbumDataHandler := handlers.NewCoreUserMediaItemDataHandler[*types.Album](coreAlbumDataService)
	coreArtistDataHandler := handlers.NewCoreUserMediaItemDataHandler[*types.Artist](coreArtistDataService)
	coreCollectionDataHandler := handlers.NewCoreUserMediaItemDataHandler[*types.Collection](coreCollectionDataService)
	
	// User Data Repositories
	userSeasonDataRepo := repository.NewUserMediaItemDataRepository[*types.Season](db, coreSeasonDataRepo)
	userEpisodeDataRepo := repository.NewUserMediaItemDataRepository[*types.Episode](db, coreEpisodeDataRepo)
	userTrackDataRepo := repository.NewUserMediaItemDataRepository[*types.Track](db, coreTrackDataRepo)
	userAlbumDataRepo := repository.NewUserMediaItemDataRepository[*types.Album](db, coreAlbumDataRepo)
	userArtistDataRepo := repository.NewUserMediaItemDataRepository[*types.Artist](db, coreArtistDataRepo)
	userCollectionDataRepo := repository.NewUserMediaItemDataRepository[*types.Collection](db, coreCollectionDataRepo)
	
	// User Data Services
	userSeasonDataService := services.NewUserMediaItemDataService[*types.Season](coreSeasonDataService, userSeasonDataRepo)
	userEpisodeDataService := services.NewUserMediaItemDataService[*types.Episode](coreEpisodeDataService, userEpisodeDataRepo)
	userTrackDataService := services.NewUserMediaItemDataService[*types.Track](coreTrackDataService, userTrackDataRepo)
	userAlbumDataService := services.NewUserMediaItemDataService[*types.Album](coreAlbumDataService, userAlbumDataRepo)
	userArtistDataService := services.NewUserMediaItemDataService[*types.Artist](coreArtistDataService, userArtistDataRepo)
	userCollectionDataService := services.NewUserMediaItemDataService[*types.Collection](coreCollectionDataService, userCollectionDataRepo)
	
	// User Data Handlers
	seasonDataHandler := handlers.NewUserMediaItemDataHandler[*types.Season](coreSeasonDataHandler, userSeasonDataService)
	episodeDataHandler := handlers.NewUserMediaItemDataHandler[*types.Episode](coreEpisodeDataHandler, userEpisodeDataService)
	trackDataHandler := handlers.NewUserMediaItemDataHandler[*types.Track](coreTrackDataHandler, userTrackDataService)
	albumDataHandler := handlers.NewUserMediaItemDataHandler[*types.Album](coreAlbumDataHandler, userAlbumDataService)
	artistDataHandler := handlers.NewUserMediaItemDataHandler[*types.Artist](coreArtistDataHandler, userArtistDataService)
	collectionDataHandler := handlers.NewUserMediaItemDataHandler[*types.Collection](coreCollectionDataHandler, userCollectionDataService)
	
	return MediaDataHandlers{
		MovieDataHandler:      movieDataHandler,
		SeriesDataHandler:     seriesDataHandler,
		SeasonDataHandler:     seasonDataHandler,
		EpisodeDataHandler:    episodeDataHandler,
		TrackDataHandler:      trackDataHandler,
		AlbumDataHandler:      albumDataHandler,
		ArtistDataHandler:     artistDataHandler,
		PlaylistDataHandler:   playlistDataHandler,
		CollectionDataHandler: collectionDataHandler,
	}
}

// ProvideClientHandlers creates a simplified client handlers group struct
func ProvideClientHandlers(
	clientsHandler *handlers.ClientsHandler,
	embyHandler *handlers.ClientHandler[*clienttypes.EmbyConfig],
	jellyfinHandler *handlers.ClientHandler[*clienttypes.JellyfinConfig],
	plexHandler *handlers.ClientHandler[*clienttypes.PlexConfig],
	aiHandler *handlers.AIHandler[*clienttypes.ClaudeConfig],
	claudeHandler *handlers.ClientHandler[*clienttypes.ClaudeConfig],
	metadataHandler *handlers.MetadataClientHandler[*clienttypes.TMDBConfig],
) ClientHandlers {
	// Create stub handlers for the ones we're not providing
	clientFactory := client.GetClientFactoryService()
	db, _ := gorm.Open(nil, &gorm.Config{}) 
	
	// Create repositories
	subsonicRepo := repository.NewClientRepository[*clienttypes.SubsonicConfig](db)
	sonarrRepo := repository.NewClientRepository[*clienttypes.SonarrConfig](db)
	radarrRepo := repository.NewClientRepository[*clienttypes.RadarrConfig](db)
	lidarrRepo := repository.NewClientRepository[*clienttypes.LidarrConfig](db)
	openaiRepo := repository.NewClientRepository[*clienttypes.OpenAIConfig](db)
	ollamaRepo := repository.NewClientRepository[*clienttypes.OllamaConfig](db)
	
	// Create services
	subsonicService := services.NewClientService[*clienttypes.SubsonicConfig](clientFactory, subsonicRepo)
	sonarrService := services.NewClientService[*clienttypes.SonarrConfig](clientFactory, sonarrRepo)
	radarrService := services.NewClientService[*clienttypes.RadarrConfig](clientFactory, radarrRepo)
	lidarrService := services.NewClientService[*clienttypes.LidarrConfig](clientFactory, lidarrRepo)
	openAIService := services.NewClientService[*clienttypes.OpenAIConfig](clientFactory, openaiRepo)
	ollamaService := services.NewClientService[*clienttypes.OllamaConfig](clientFactory, ollamaRepo)
	
	// Create handlers
	subsonicHandler := handlers.NewClientHandler[*clienttypes.SubsonicConfig](subsonicService)
	radarrHandler := handlers.NewClientHandler[*clienttypes.RadarrConfig](radarrService)
	sonarrHandler := handlers.NewClientHandler[*clienttypes.SonarrConfig](sonarrService)
	lidarrHandler := handlers.NewClientHandler[*clienttypes.LidarrConfig](lidarrService)
	openAIHandler := handlers.NewClientHandler[*clienttypes.OpenAIConfig](openAIService)
	ollamaHandler := handlers.NewClientHandler[*clienttypes.OllamaConfig](ollamaService)
	
	return ClientHandlers{
		ClientsHandler:  clientsHandler,
		EmbyHandler:     embyHandler,
		JellyfinHandler: jellyfinHandler,
		PlexHandler:     plexHandler,
		SubsonicHandler: subsonicHandler,
		RadarrHandler:   radarrHandler,
		SonarrHandler:   sonarrHandler,
		LidarrHandler:   lidarrHandler,
		AIHandler:       aiHandler,
		ClaudeHandler:   claudeHandler,
		OpenAIHandler:   openAIHandler,
		OllamaHandler:   ollamaHandler,
		MetadataHandler: metadataHandler,
	}
}

// ProvideSpecializedHandlers creates the specialized handlers group struct
func ProvideSpecializedHandlers(
	recommendationHandler *handlers.RecommendationHandler,
	calendarHandler *handlers.CalendarHandler,
) SpecializedHandlers {
	return SpecializedHandlers{
		RecommendationHandler: recommendationHandler,
		CalendarHandler:       calendarHandler,
	}
}

// InitializeAuthHandler initializes only the auth handler (which has no generics)
func InitializeAuthHandler(ctx context.Context) (*handlers.AuthHandler, error) {
	wire.Build(
		// Database
		ProvideDB,

		// Core Repositories
		ProvideUserRepository,
		ProvideSessionRepository,
		
		// Auth Service
		ProvideAuthService,
		
		// Handler
		ProvideAuthHandler,
	)

	return nil, nil
}

// InitializeConfigHandler initializes only the config handler
func InitializeConfigHandler(ctx context.Context) (*handlers.ConfigHandler, error) {
	wire.Build(
		ProvideDB,
		ProvideConfigRepository,
		ProvideConfigService,
		ProvideConfigHandler,
	)

	return nil, nil
}

// InitializeHealthHandler initializes only the health handler
func InitializeHealthHandler(ctx context.Context) (*handlers.HealthHandler, error) {
	wire.Build(
		ProvideDB,
		ProvideHealthService,
		ProvideHealthHandler,
	)

	return nil, nil
}

// InitializeUserHandler initializes only the user handler
func InitializeUserHandler(ctx context.Context) (*handlers.UserHandler, error) {
	wire.Build(
		ProvideDB,
		ProvideUserRepository,
		ProvideUserService,
		ProvideConfigRepository,
		ProvideConfigService,
		ProvideUserHandler,
	)

	return nil, nil
}

// InitializeAllHandlers initializes all handlers by combining individual Wire-generated handlers
// with manually created instances of the other handler groups to work around Wire's 
// generic type limitations
func InitializeAllHandlers(ctx context.Context) (ApplicationHandlers, error) {
	// Initialize non-generic handlers individually using Wire
	authHandler, err := InitializeAuthHandler(ctx)
	if err != nil {
		return ApplicationHandlers{}, fmt.Errorf("failed to initialize auth handler: %w", err)
	}
	
	configHandler, err := InitializeConfigHandler(ctx)
	if err != nil {
		return ApplicationHandlers{}, fmt.Errorf("failed to initialize config handler: %w", err)
	}
	
	healthHandler, err := InitializeHealthHandler(ctx)
	if err != nil {
		return ApplicationHandlers{}, fmt.Errorf("failed to initialize health handler: %w", err)
	}
	
	userHandler, err := InitializeUserHandler(ctx)
	if err != nil {
		return ApplicationHandlers{}, fmt.Errorf("failed to initialize user handler: %w", err)
	}
	
	// Get database connection for manual initialization of other handlers
	db, err := ProvideDB()
	if err != nil {
		return ApplicationHandlers{}, fmt.Errorf("failed to provide database connection: %w", err)
	}
	
	// Setup core services for people and credits
	personRepo := ProvidePersonRepository(db)
	creditRepo := ProvideCreditRepository(db)
	personService := ProvidePersonService(personRepo, creditRepo)
	creditService := ProvideCreditService(db)
	
	// Initialize recommendation service and handlers
	jobRepo := ProvideJobRepository(db)
	recommendationService := ProvideRecommendationService(jobRepo)
	recommendationHandler := ProvideRecommendationHandler(recommendationService)

	// Initialize handlers for people and credits
	peopleHandler := ProvidePeopleHandler(personService)
	creditHandler := ProvideCreditHandler(creditService)
	calendarHandler := ProvideCalendarHandler()
	
	// Set up client factory service
	clientFactoryForHandlers := ProvideClientFactoryService()
	
	// Initialize client repositories
	embyRepo := ProvideEmbyClientRepository(db)
	jellyfinRepo := ProvideJellyfinClientRepository(db)
	plexRepo := ProvidePlexClientRepository(db)
	claudeRepo := ProvideClaudeClientRepository(db)
	
	// Initialize client services
	embyService := ProvideEmbyClientService(clientFactoryForHandlers, embyRepo)
	jellyfinService := ProvideJellyfinClientService(clientFactoryForHandlers, jellyfinRepo)
	plexService := ProvidePlexClientService(clientFactoryForHandlers, plexRepo)
	claudeService := ProvideClaudeClientService(clientFactoryForHandlers, claudeRepo)
	
	// Initialize client handlers
	embyHandler := ProvideEmbyClientHandler(embyService)
	jellyfinHandler := ProvideJellyfinClientHandler(jellyfinService)
	plexHandler := ProvidePlexClientHandler(plexService)
	claudeHandler := ProvideClaudeClientHandler(claudeService)
	clientsHandler := ProvideClientsHandler(embyService, jellyfinService, plexService, claudeService)
	aiHandler := ProvideAIHandler(clientFactoryForHandlers, claudeService)
	
	// Initialize metadata client handler
	metadataService := ProvideMetadataClientService(clientFactoryForHandlers, db)
	metadataHandler := ProvideMetadataClientHandler(metadataService)
	
	// Initialize movie repositories and services
	movieRepo := ProvideMovieRepository(db)
	userMovieRepo := ProvideUserMovieRepository(db)
	coreMovieService := ProvideCoreMovieService(movieRepo)
	userMovieService := ProvideUserMovieService(coreMovieService, userMovieRepo)
	
	// Initialize movie data repositories and services
	coreMovieDataRepo := ProvideCoreMovieDataRepository(db)
	userMovieDataRepo := ProvideUserMovieDataRepository(db, coreMovieDataRepo)
	coreMovieDataService := ProvideCoreMovieDataService(coreMovieService, coreMovieDataRepo)
	userMovieDataService := ProvideUserMovieDataService(coreMovieDataService, userMovieDataRepo)
	
	// Initialize series repositories and services
	seriesRepo := ProvideSeriesRepository(db)
	userSeriesRepo := ProvideUserSeriesRepository(db)
	coreSeriesService := ProvideCoreSeriesService(seriesRepo)
	userSeriesService := ProvideUserSeriesService(coreSeriesService, userSeriesRepo)
	
	// Initialize series data repositories and services
	coreSeriesDataRepo := ProvideCoreSeriesDataRepository(db)
	userSeriesDataRepo := ProvideUserSeriesDataRepository(db, coreSeriesDataRepo)
	coreSeriesDataService := ProvideCoreSeriesDataService(coreSeriesService, coreSeriesDataRepo)
	userSeriesDataService := ProvideUserSeriesDataService(coreSeriesDataService, userSeriesDataRepo)
	
	// Initialize playlist repositories and services
	playlistRepo := ProvidePlaylistRepository(db)
	userPlaylistRepo := ProvideUserPlaylistRepository(db)
	corePlaylistService := ProvideCorePlaylistService(playlistRepo)
	corePlaylistListService := ProvideCorePlaylistListService(playlistRepo)
	userPlaylistService := ProvideUserPlaylistService(corePlaylistService, userPlaylistRepo)
	
	// Initialize playlist data repositories and services
	corePlaylistDataRepo := ProvideCorePlaylistDataRepository(db)
	userPlaylistDataRepo := ProvideUserPlaylistDataRepository(db, corePlaylistDataRepo)
	corePlaylistDataService := ProvideCorePlaylistDataService(corePlaylistService, corePlaylistDataRepo)
	userPlaylistDataService := ProvideUserPlaylistDataService(corePlaylistDataService, userPlaylistDataRepo)
	
	// Initialize list services
	userPlaylistListService := ProvideUserPlaylistListService(corePlaylistListService, userPlaylistRepo, userPlaylistDataRepo)
	
	// Initialize handlers
	userMovieHandler := ProvideUserMovieHandler(userMovieService)
	userSeriesHandler := ProvideUserSeriesHandler(userSeriesService)
	corePlaylistHandler := ProvideCorePlaylistHandler(corePlaylistService)
	corePlaylistListHandler := ProvideCorePlaylistListHandler(corePlaylistHandler, corePlaylistListService)
	userPlaylistListHandler := ProvideUserPlaylistListHandler(corePlaylistListHandler, userPlaylistService, userPlaylistListService)
	
	// Initialize media data handlers
	coreMovieDataHandler := ProvideCoreMovieDataHandler(coreMovieDataService)
	coreSeriesDataHandler := ProvideCoreSeriesDataHandler(coreSeriesDataService)
	corePlaylistDataHandler := ProvideCorePlaylistDataHandler(corePlaylistDataService)
	
	userMovieDataHandler := ProvideUserMovieDataHandler(coreMovieDataHandler, userMovieDataService)
	userSeriesDataHandler := ProvideUserSeriesDataHandler(coreSeriesDataHandler, userSeriesDataService)
	userPlaylistDataHandler := ProvideUserPlaylistDataHandler(corePlaylistDataHandler, userPlaylistDataService)
	
	// Create handler groups
	mediaHandlers := ProvideMediaHandlers(
		userMovieHandler,
		userSeriesHandler, 
		userPlaylistListHandler,
		peopleHandler,
		creditHandler,
	)
	
	mediaDataHandlers := ProvideMediaDataHandlers(
		userMovieDataHandler,
		userSeriesDataHandler,
		userPlaylistDataHandler,
	)
	
	clientHandlers := ProvideClientHandlers(
		clientsHandler,
		embyHandler,
		jellyfinHandler,
		plexHandler,
		aiHandler,
		claudeHandler,
		metadataHandler,
	)
	
	specializedHandlers := ProvideSpecializedHandlers(
		recommendationHandler,
		calendarHandler,
	)
	
	// Create the job handler and search handler manually since they depend on generic types
	jobRepository := repository.NewJobRepository(db)
	searchRepository := repository.NewSearchRepository(db)
	userConfigRepository := repository.NewUserConfigRepository(db)
	
	// Manually create job service and handler
	recJob := &recommendation.RecommendationJob{}
	mSyncJob := &jobs.MediaSyncJob{}
	wHistorySyncJob := &jobs.WatchHistorySyncJob{}
	favSyncJob := &jobs.FavoritesSyncJob{}
	
	// Create repositories for job service
	movieRepoForJob := repository.NewMediaItemRepository[*types.Movie](db)
	seriesRepoForJob := repository.NewMediaItemRepository[*types.Series](db)
	trackRepoForJob := repository.NewMediaItemRepository[*types.Track](db)
	
	// Create data repositories for job service
	coreMovieDataRepoForJob := repository.NewCoreUserMediaItemDataRepository[*types.Movie](db)
	userMovieDataRepoForJob := repository.NewUserMediaItemDataRepository[*types.Movie](db, coreMovieDataRepoForJob)
	
	coreSeriesDataRepoForJob := repository.NewCoreUserMediaItemDataRepository[*types.Series](db)
	userSeriesDataRepoForJob := repository.NewUserMediaItemDataRepository[*types.Series](db, coreSeriesDataRepoForJob)
	
	coreTrackDataRepoForJob := repository.NewCoreUserMediaItemDataRepository[*types.Track](db)
	userTrackDataRepoForJob := repository.NewUserMediaItemDataRepository[*types.Track](db, coreTrackDataRepoForJob)
	
	// Create user repository for job service
	userRepoForJob := repository.NewUserRepository(db)
	
	// Create job service
	jobService := services.NewJobService(
		jobRepository,
		userRepoForJob,
		userConfigRepository,
		movieRepoForJob,
		seriesRepoForJob,
		trackRepoForJob,
		userMovieDataRepoForJob,
		userSeriesDataRepoForJob,
		userTrackDataRepoForJob,
		recJob,
		mSyncJob,
		wHistorySyncJob,
		favSyncJob,
	)
	
	// Create job handler
	jobHandler := handlers.NewJobHandler(jobService)
	
	// Create user config service and handler
	userConfigService := services.NewUserConfigService(
		userConfigRepository,
		jobService,
		recJob,
	)
	userConfigHandler := handlers.NewUserConfigHandler(userConfigService)
	
	// Create search service and handler
	personRepoForSearch := repository.NewPersonRepository(db)
	
	// Create client repositories for search service
	embyRepoForSearch := repository.NewClientRepository[*clienttypes.EmbyConfig](db)
	jellyfinRepoForSearch := repository.NewClientRepository[*clienttypes.JellyfinConfig](db)
	plexRepoForSearch := repository.NewClientRepository[*clienttypes.PlexConfig](db)
	claudeRepoForSearch := repository.NewClientRepository[*clienttypes.ClaudeConfig](db)
	
	// Create client repository bundle
	clientRepos := repobundles.NewClientRepositories(
		embyRepoForSearch,
		jellyfinRepoForSearch,
		plexRepoForSearch,
		repository.NewClientRepository[*clienttypes.SubsonicConfig](db),
		repository.NewClientRepository[*clienttypes.SonarrConfig](db),
		repository.NewClientRepository[*clienttypes.RadarrConfig](db),
		repository.NewClientRepository[*clienttypes.LidarrConfig](db),
		claudeRepoForSearch,
		repository.NewClientRepository[*clienttypes.OpenAIConfig](db),
		repository.NewClientRepository[*clienttypes.OllamaConfig](db),
	)
	
	// Create core media item repositories bundle
	itemRepos := repobundles.NewCoreMediaItemRepositories(
		movieRepoForJob,
		seriesRepoForJob,
		repository.NewMediaItemRepository[*types.Season](db),
		repository.NewMediaItemRepository[*types.Episode](db),
		trackRepoForJob,
		repository.NewMediaItemRepository[*types.Album](db),
		repository.NewMediaItemRepository[*types.Artist](db),
		repository.NewMediaItemRepository[*types.Collection](db),
		repository.NewMediaItemRepository[*types.Playlist](db),
	)
	
	// Create client factory service
	clientFactoryServiceForSearch := client.GetClientFactoryService()
	
	// Create search service
	searchService := services.NewSearchService(
		searchRepository,
		clientRepos,
		itemRepos,
		personRepoForSearch,
		clientFactoryServiceForSearch,
	)
	
	// Create search handler
	searchHandler := handlers.NewSearchHandler(searchService)
	
	// Create system handlers struct
	systemHandlers := SystemHandlers{
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		ConfigHandler:    configHandler,
		JobHandler:       jobHandler,
		HealthHandler:    healthHandler,
		SearchHandler:    searchHandler,
		UserConfigHandler: userConfigHandler,
	}
	
	// Build the complete ApplicationHandlers
	return ApplicationHandlers{
		System:      systemHandlers,
		Media:       mediaHandlers,
		MediaData:   mediaDataHandlers,
		Clients:     clientHandlers,
		Specialized: specializedHandlers,
	}, nil
}
