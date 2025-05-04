package jellyfin

import (
	"context"
	jellyfin "github.com/sj14/jellyfin-go/api"
	"strings"
	"suasor/clients/media/types"
	"suasor/utils/logger"
	"time"
)

// JellyfinQueryOptions is a shadow struct that mirrors the Jellyfin API's query parameters
// but provides direct access to all fields. This helps overcome the private field access
// limitations in the official Jellyfin API client.
type JellyfinQueryOptions struct {
	UserId                  *string
	MaxOfficialRating       *string
	HasThemeSong            *bool
	HasThemeVideo           *bool
	HasSubtitles            *bool
	HasSpecialFeature       *bool
	HasTrailer              *bool
	AdjacentTo              *string
	IndexNumber             *int32
	ParentIndexNumber       *int32
	HasParentalRating       *bool
	IsHd                    *bool
	Is4K                    *bool
	LocationTypes           *[]jellyfin.LocationType
	ExcludeLocationTypes    *[]jellyfin.LocationType
	IsMissing               *bool
	IsUnaired               *bool
	MinCommunityRating      *float64
	MinCriticRating         *float64
	MinPremiereDate         *time.Time
	MinDateLastSaved        *time.Time
	MinDateLastSavedForUser *time.Time
	MaxPremiereDate         *time.Time
	HasOverview             *bool
	HasImdbId               *bool
	HasTmdbId               *bool
	HasTvdbId               *bool
	IsMovie                 *bool
	IsSeries                *bool
	IsNews                  *bool
	IsKids                  *bool
	IsSports                *bool
	ExcludeItemIds          *[]string
	StartIndex              *int32
	Limit                   *int32
	Recursive               *bool
	SearchTerm              *string
	SortOrder               *[]jellyfin.SortOrder
	ParentId                *string
	Fields                  *[]jellyfin.ItemFields
	ExcludeItemTypes        *[]jellyfin.BaseItemKind
	IncludeItemTypes        *[]jellyfin.BaseItemKind
	Filters                 *[]jellyfin.ItemFilter
	IsFavorite              *bool
	MediaTypes              *[]jellyfin.MediaType
	ImageTypes              *[]jellyfin.ImageType
	SortBy                  *[]jellyfin.ItemSortBy
	IsPlayed                *bool
	Genres                  *[]string
	OfficialRatings         *[]string
	Tags                    *[]string
	Years                   *[]int32
	EnableUserData          *bool
	ImageTypeLimit          *int32
	EnableImageTypes        *[]jellyfin.ImageType
	Person                  *string
	PersonIds               *[]string
	PersonTypes             *[]string
	Studios                 *[]string
	StudioIds               *[]string
	Artists                 *[]string
	ExcludeArtistIds        *[]string
	ArtistIds               *[]string
	AlbumArtistIds          *[]string
	ContributingArtistIds   *[]string
	Albums                  *[]string
	AlbumIds                *[]string
	Ids                     *[]string
	VideoTypes              *[]jellyfin.VideoType
	MinOfficialRating       *string
	IsLocked                *bool
	IsPlaceHolder           *bool
	HasOfficialRating       *bool
	CollapseBoxSetItems     *bool
	MinWidth                *int32
	MinHeight               *int32
	MaxWidth                *int32
	MaxHeight               *int32
	Is3D                    *bool
	SeriesStatus            *[]jellyfin.SeriesStatus
	NameStartsWithOrGreater *string
	NameStartsWith          *string
	NameLessThan            *string
	GenreIds                *[]string
	EnableTotalRecordCount  *bool
	EnableImages            *bool
}

// NewJellyfinQueryOptions creates a new instance of JellyfinQueryOptions
func NewJellyfinQueryOptions(ctx context.Context, options *types.QueryOptions) *JellyfinQueryOptions {
	// Always create a new instance
	jellyfinOptions := &JellyfinQueryOptions{}

	// Only process options if they are provided
	if options != nil {
		jellyfinOptions.FromQueryOptions(ctx, options)
	}

	return jellyfinOptions
}

// ToItemsRequest converts JellyfinQueryOptions to an API request object for GetItems
func (j *JellyfinQueryOptions) SetItemsRequest(ctx context.Context, req *jellyfin.ApiGetItemsRequest) {
	log := logger.LoggerFromContext(ctx)
	// Defensive programming: check for nil pointers
	if j == nil || req == nil {
		log.Debug().Msg("No options provided, skipping query options")
		return
	}

	// Apply all parameters using method-based approach
	if j.UserId != nil {
		log.Debug().Str("userID", *j.UserId).Msg("Applying user ID filter")
		req.UserId(*j.UserId)
	}

	if j.MaxOfficialRating != nil {
		log.Debug().Str("maxOfficialRating", *j.MaxOfficialRating).Msg("Applying max official rating filter")
		req.MaxOfficialRating(*j.MaxOfficialRating)
	}

	if j.HasThemeSong != nil {
		log.Debug().Bool("hasThemeSong", *j.HasThemeSong).Msg("Applying has theme song filter")
		req.HasThemeSong(*j.HasThemeSong)
	}

	if j.HasThemeVideo != nil {
		log.Debug().Bool("hasThemeVideo", *j.HasThemeVideo).Msg("Applying has theme video filter")
		req.HasThemeVideo(*j.HasThemeVideo)
	}

	if j.HasSubtitles != nil {
		log.Debug().Bool("hasSubtitles", *j.HasSubtitles).Msg("Applying has subtitles filter")
		req.HasSubtitles(*j.HasSubtitles)
	}

	if j.HasSpecialFeature != nil {
		log.Debug().Bool("hasSpecialFeature", *j.HasSpecialFeature).Msg("Applying has special feature filter")
		req.HasSpecialFeature(*j.HasSpecialFeature)
	}

	if j.HasTrailer != nil {
		log.Debug().Bool("hasTrailer", *j.HasTrailer).Msg("Applying has trailer filter")
		req.HasTrailer(*j.HasTrailer)
	}

	if j.AdjacentTo != nil {
		log.Debug().Str("adjacentTo", *j.AdjacentTo).Msg("Applying adjacent to filter")
		req.AdjacentTo(*j.AdjacentTo)
	}

	if j.IndexNumber != nil {
		log.Debug().Int32("indexNumber", *j.IndexNumber).Msg("Applying index number filter")
		req.IndexNumber(*j.IndexNumber)
	}

	if j.ParentIndexNumber != nil {
		log.Debug().Int32("parentIndexNumber", *j.ParentIndexNumber).Msg("Applying parent index number filter")
		req.ParentIndexNumber(*j.ParentIndexNumber)
	}

	if j.HasParentalRating != nil {
		log.Debug().Bool("hasParentalRating", *j.HasParentalRating).Msg("Applying has parental rating filter")
		req.HasParentalRating(*j.HasParentalRating)
	}

	if j.IsHd != nil {
		log.Debug().Bool("isHd", *j.IsHd).Msg("Applying is HD filter")
		req.IsHd(*j.IsHd)
	}

	if j.Is4K != nil {
		log.Debug().Bool("is4K", *j.Is4K).Msg("Applying is 4K filter")
		req.Is4K(*j.Is4K)
	}

	if j.LocationTypes != nil {
		log.Debug().Interface("locationTypes", *j.LocationTypes).Msg("Applying location types filter")
		req.LocationTypes(*j.LocationTypes)
	}

	if j.ExcludeLocationTypes != nil {
		log.Debug().Interface("excludeLocationTypes", *j.ExcludeLocationTypes).Msg("Applying exclude location types filter")
		req.ExcludeLocationTypes(*j.ExcludeLocationTypes)
	}

	if j.IsMissing != nil {
		log.Debug().Bool("isMissing", *j.IsMissing).Msg("Applying is missing filter")
		req.IsMissing(*j.IsMissing)
	}

	if j.IsUnaired != nil {
		log.Debug().Bool("isUnaired", *j.IsUnaired).Msg("Applying is unaired filter")
		req.IsUnaired(*j.IsUnaired)
	}

	if j.MinCommunityRating != nil {
		log.Debug().Float64("minCommunityRating", *j.MinCommunityRating).Msg("Applying min community rating filter")
		req.MinCommunityRating(*j.MinCommunityRating)
	}

	if j.MinCriticRating != nil {
		log.Debug().Float64("minCriticRating", *j.MinCriticRating).Msg("Applying min critic rating filter")
		req.MinCriticRating(*j.MinCriticRating)
	}

	if j.MinPremiereDate != nil {
		log.Debug().Time("minPremiereDate", *j.MinPremiereDate).Msg("Applying min premiere")
		req.MinPremiereDate(*j.MinPremiereDate)
	}

	if j.MinDateLastSaved != nil {
		log.Debug().Time("minDateLastSaved", *j.MinDateLastSaved).Msg("Applying min date last saved filter")
		req.MinDateLastSaved(*j.MinDateLastSaved)
	}

	if j.MinDateLastSavedForUser != nil {
		log.Debug().Time("minDateLastSavedForUser", *j.MinDateLastSavedForUser).Msg("Applying min date last saved for user filter")
		req.MinDateLastSavedForUser(*j.MinDateLastSavedForUser)
	}

	if j.MaxPremiereDate != nil {
		log.Debug().Time("maxPremiereDate", *j.MaxPremiereDate).Msg("Applying max premiere date filter")
		req.MaxPremiereDate(*j.MaxPremiereDate)
	}

	if j.HasOverview != nil {
		log.Debug().Bool("hasOverview", *j.HasOverview).Msg("Applying has overview filter")
		req.HasOverview(*j.HasOverview)

	}

	if j.HasImdbId != nil {
		log.Debug().Bool("hasImdbId", *j.HasImdbId).Msg("Applying has IMDB ID filter")
		req.HasImdbId(*j.HasImdbId)
	}

	if j.HasTmdbId != nil {
		log.Debug().Bool("hasTmdbId", *j.HasTmdbId).Msg("Applying has TMDB ID filter")
		req.HasTmdbId(*j.HasTmdbId)
	}

	if j.HasTvdbId != nil {
		log.Debug().Bool("hasTvdbId", *j.HasTvdbId).Msg("Applying has TVDB ID filter")
		req.HasTvdbId(*j.HasTvdbId)
	}

	if j.IsMovie != nil {
		log.Debug().Bool("isMovie", *j.IsMovie).Msg("Applying is movie filter")
		req.IsMovie(*j.IsMovie)
	}

	if j.IsSeries != nil {
		log.Debug().Bool("isSeries", *j.IsSeries).Msg("Applying is series filter")
		req.IsSeries(*j.IsSeries)
	}

	if j.IsNews != nil {
		log.Debug().Bool("isNews", *j.IsNews).Msg("Applying is news filter")
		req.IsNews(*j.IsNews)
	}

	if j.IsKids != nil {
		log.Debug().Bool("isKids", *j.IsKids).Msg("Applying is kids filter")
		req.IsKids(*j.IsKids)
	}

	if j.IsSports != nil {
		log.Debug().Bool("isSports", *j.IsSports).Msg("Applying is sports filter")
		req.IsSports(*j.IsSports)
	}

	if j.ExcludeItemIds != nil && len(*j.ExcludeItemIds) > 0 {
		log.Debug().Interface("excludeItemIds", *j.ExcludeItemIds).Msg("Applying exclude item IDs filter")
		req.ExcludeItemIds(*j.ExcludeItemIds)
	}

	if j.StartIndex != nil {
		log.Debug().Int32("startIndex", *j.StartIndex).Msg("Applying start index filter")
		req.StartIndex(*j.StartIndex)

	}

	if j.Limit != nil {
		log.Debug().Int32("limit", *j.Limit).Msg("Applying limit filter")
		req.Limit(*j.Limit)
	}

	if j.Recursive != nil {
		log.Debug().Bool("recursive", *j.Recursive).Msg("Applying recursive filter")
		req.Recursive(*j.Recursive)
	}

	if j.SearchTerm != nil {
		log.Debug().Str("searchTerm", *j.SearchTerm).Msg("Applying search term filter")
		req.SearchTerm(*j.SearchTerm)
	}

	if j.SortOrder != nil {
		log.Debug().Interface("sortOrder", *j.SortOrder).Msg("Applying sort order filter")
		req.SortOrder(*j.SortOrder)
	}

	if j.ParentId != nil {
		log.Debug().Str("parentId", *j.ParentId).Msg("Applying parent ID filter")
		req.ParentId(*j.ParentId)
	}

	if j.Fields != nil {
		req.Fields(*j.Fields)
		log.Debug().Interface("fields", *j.Fields).Msg("Applying fields filter")
	}

	if j.ExcludeItemTypes != nil {
		log.Debug().Interface("excludeItemTypes", *j.ExcludeItemTypes).Msg("Applying exclude item types filter")
		req.ExcludeItemTypes(*j.ExcludeItemTypes)
	}

	if j.IncludeItemTypes != nil {
		log.Debug().Interface("includeItemTypes", *j.IncludeItemTypes).Msg("Applying include item types filter")
		req.IncludeItemTypes(*j.IncludeItemTypes)
	}

	if j.Filters != nil {
		log.Debug().Interface("filters", *j.Filters).Msg("Applying filters filter")
		req.Filters(*j.Filters)
	}

	if j.IsFavorite != nil {
		log.Debug().Bool("isFavorite", *j.IsFavorite).Msg("Applying is favorite filter")
		req.IsFavorite(*j.IsFavorite)
	}

	if j.MediaTypes != nil {
		log.Debug().Interface("mediaTypes", *j.MediaTypes).Msg("Applying media types filter")
		req.MediaTypes(*j.MediaTypes)
	}

	if j.ImageTypes != nil {
		log.Debug().Interface("imageTypes", *j.ImageTypes).Msg("Applying image types filter")
		req.ImageTypes(*j.ImageTypes)
	}

	if j.SortBy != nil {
		log.Debug().Interface("sortBy", *j.SortBy).Msg("Applying sort by filter")
		req.SortBy(*j.SortBy)
	}

	if j.IsPlayed != nil {
		log.Debug().Bool("isPlayed", *j.IsPlayed).Msg("Applying is played filter")
		req.IsPlayed(*j.IsPlayed)
	}

	if j.Genres != nil {
		log.Debug().Interface("genres", *j.Genres).Msg("Applying genres filter")
		req.Genres(*j.Genres)
	}

	if j.OfficialRatings != nil {
		log.Debug().Interface("officialRatings", *j.OfficialRatings).Msg("Applying official ratings filter")
		req.OfficialRatings(*j.OfficialRatings)
	}

	if j.Tags != nil {
		log.Debug().Interface("tags", *j.Tags).Msg("Applying tags filter")
		req.Tags(*j.Tags)
	}

	if j.Years != nil {
		log.Debug().Interface("years", *j.Years).Msg("Applying years filter")
		req.Years(*j.Years)
	}

	if j.EnableUserData != nil {
		req.EnableUserData(*j.EnableUserData)
	}

	if j.ImageTypeLimit != nil {
		req.ImageTypeLimit(*j.ImageTypeLimit)
	}

	if j.EnableImageTypes != nil {
		req.EnableImageTypes(*j.EnableImageTypes)
	}

	if j.Person != nil {
		req.Person(*j.Person)
	}

	if j.PersonIds != nil && len(*j.PersonIds) > 0 {
		req.PersonIds(*j.PersonIds)
	}

	if j.PersonTypes != nil && len(*j.PersonTypes) > 0 {
		req.PersonTypes(*j.PersonTypes)
	}

	if j.Studios != nil && len(*j.Studios) > 0 {
		req.Studios(*j.Studios)
	}

	if j.StudioIds != nil && len(*j.StudioIds) > 0 {
		req.StudioIds(*j.StudioIds)
	}

	if j.Artists != nil && len(*j.Artists) > 0 {
		req.Artists(*j.Artists)
	}

	if j.ExcludeArtistIds != nil && len(*j.ExcludeArtistIds) > 0 {
		req.ExcludeArtistIds(*j.ExcludeArtistIds)
	}

	if j.ArtistIds != nil && len(*j.ArtistIds) > 0 {
		req.ArtistIds(*j.ArtistIds)
	}

	if j.AlbumArtistIds != nil && len(*j.AlbumArtistIds) > 0 {
		req.AlbumArtistIds(*j.AlbumArtistIds)
	}

	if j.ContributingArtistIds != nil && len(*j.ContributingArtistIds) > 0 {
		req.ContributingArtistIds(*j.ContributingArtistIds)
	}

	if j.Albums != nil && len(*j.Albums) > 0 {
		req.Albums(*j.Albums)
	}

	if j.AlbumIds != nil && len(*j.AlbumIds) > 0 {
		req.AlbumIds(*j.AlbumIds)
	}

	if j.Ids != nil && len(*j.Ids) > 0 {
		req.Ids(*j.Ids)
	}

	if j.VideoTypes != nil {
		req.VideoTypes(*j.VideoTypes)
	}

	if j.MinOfficialRating != nil {
		req.MinOfficialRating(*j.MinOfficialRating)
	}

	if j.IsLocked != nil {
		req.IsLocked(*j.IsLocked)
	}

	if j.IsPlaceHolder != nil {
		req.IsPlaceHolder(*j.IsPlaceHolder)
	}

	if j.HasOfficialRating != nil {
		req.HasOfficialRating(*j.HasOfficialRating)
	}

	if j.CollapseBoxSetItems != nil {
		req.CollapseBoxSetItems(*j.CollapseBoxSetItems)
	}

	if j.MinWidth != nil {
		req.MinWidth(*j.MinWidth)
	}

	if j.MinHeight != nil {
		req.MinHeight(*j.MinHeight)
	}

	if j.MaxWidth != nil {
		req.MaxWidth(*j.MaxWidth)
	}

	if j.MaxHeight != nil {
		req.MaxHeight(*j.MaxHeight)
	}

	if j.Is3D != nil {
		req.Is3D(*j.Is3D)
	}

	if j.SeriesStatus != nil {
		req.SeriesStatus(*j.SeriesStatus)
	}

	if j.NameStartsWithOrGreater != nil {
		req.NameStartsWithOrGreater(*j.NameStartsWithOrGreater)
	}

	if j.NameStartsWith != nil {
		req.NameStartsWith(*j.NameStartsWith)
	}

	if j.NameLessThan != nil {
		req.NameLessThan(*j.NameLessThan)
	}

	if j.GenreIds != nil && len(*j.GenreIds) > 0 {
		req.GenreIds(*j.GenreIds)
	}

	if j.EnableTotalRecordCount != nil {
		req.EnableTotalRecordCount(*j.EnableTotalRecordCount)
	}

	if j.EnableImages != nil {
		req.EnableImages(*j.EnableImages)
	}

}

// ToArtistsRequest converts JellyfinQueryOptions to an API request object for GetArtists
func (j *JellyfinQueryOptions) SetArtistsRequest(req *jellyfin.ApiGetArtistsRequest) {
	// Defensive programming: check for nil pointers
	if j == nil || req == nil {
		return
	}

	if j.UserId != nil {
		req.UserId(*j.UserId)
	}

	if j.MinCommunityRating != nil {
		req.MinCommunityRating(*j.MinCommunityRating)
	}

	if j.StartIndex != nil {
		req.StartIndex(*j.StartIndex)
	}

	if j.Limit != nil {
		req.Limit(*j.Limit)
	}

	if j.SearchTerm != nil {
		req.SearchTerm(*j.SearchTerm)
	}

	if j.ParentId != nil {
		req.ParentId(*j.ParentId)
	}

	if j.Fields != nil {
		req.Fields(*j.Fields)
	}

	if j.ExcludeItemTypes != nil {
		req.ExcludeItemTypes(*j.ExcludeItemTypes)
	}

	if j.IncludeItemTypes != nil {
		req.IncludeItemTypes(*j.IncludeItemTypes)
	}

	if j.Filters != nil {
		req.Filters(*j.Filters)
	}

	if j.IsFavorite != nil {
		req.IsFavorite(*j.IsFavorite)
	}

	if j.MediaTypes != nil {
		req.MediaTypes(*j.MediaTypes)
	}

	if j.Genres != nil {
		req.Genres(*j.Genres)
	}

	if j.GenreIds != nil && len(*j.GenreIds) > 0 {
		req.GenreIds(*j.GenreIds)
	}

	if j.OfficialRatings != nil {
		req.OfficialRatings(*j.OfficialRatings)
	}

	if j.Tags != nil {
		req.Tags(*j.Tags)
	}

	if j.Years != nil {
		req.Years(*j.Years)
	}

	if j.EnableUserData != nil {
		req.EnableUserData(*j.EnableUserData)
	}

	if j.ImageTypeLimit != nil {
		req.ImageTypeLimit(*j.ImageTypeLimit)
	}

	if j.EnableImageTypes != nil {
		req.EnableImageTypes(*j.EnableImageTypes)
	}

	if j.Person != nil {
		req.Person(*j.Person)
	}

	if j.PersonIds != nil && len(*j.PersonIds) > 0 {
		req.PersonIds(*j.PersonIds)
	}

	if j.PersonTypes != nil && len(*j.PersonTypes) > 0 {
		req.PersonTypes(*j.PersonTypes)
	}

	if j.Studios != nil && len(*j.Studios) > 0 {
		req.Studios(*j.Studios)
	}

	if j.StudioIds != nil && len(*j.StudioIds) > 0 {
		req.StudioIds(*j.StudioIds)
	}

	if j.NameStartsWithOrGreater != nil {
		req.NameStartsWithOrGreater(*j.NameStartsWithOrGreater)
	}

	if j.NameStartsWith != nil {
		req.NameStartsWith(*j.NameStartsWith)
	}

	if j.NameLessThan != nil {
		req.NameLessThan(*j.NameLessThan)
	}

	if j.SortBy != nil {
		req.SortBy(*j.SortBy)
	}

	if j.SortOrder != nil {
		req.SortOrder(*j.SortOrder)
	}

	if j.EnableImages != nil {
		req.EnableImages(*j.EnableImages)
	}

	if j.EnableTotalRecordCount != nil {
		req.EnableTotalRecordCount(*j.EnableTotalRecordCount)
	}

}

// FromQueryOptions converts Suasor's QueryOptions to JellyfinQueryOptions
func (j *JellyfinQueryOptions) FromQueryOptions(ctx context.Context, options *types.QueryOptions) *JellyfinQueryOptions {
	log := logger.LoggerFromContext(ctx)
	if options == nil || j == nil {
		log.Debug().Msg("No options provided, skipping query options")
		return j
	}
	
	// Enable deeper debug logging to diagnose query option issues
	log.Debug().
		Interface("options", options).
		Msg("Converting query options to Jellyfin options")

	// ItemIDs
	if options.ItemIDs != "" {
		log.Debug().Str("itemIDs", options.ItemIDs).Msg("Applying item IDs filter")
		ids := strings.Split(options.ItemIDs, ",")
		j.Ids = &ids
	}

	// Limit
	if options.Limit > 0 {
		log.Debug().Int("limit", options.Limit).Msg("Applying limit filter")
		limit := int32(options.Limit)
		j.Limit = &limit
	}

	// Offset/StartIndex
	if options.Offset > 0 {
		log.Debug().Int("offset", options.Offset).Msg("Applying offset filter")
		startIndex := int32(options.Offset)
		j.StartIndex = &startIndex
	}

	// Sort
	if options.Sort != "" {
		log.Debug().Str("sort", string(options.Sort)).Msg("Applying sort filter")
		sortBy := []jellyfin.ItemSortBy{jellyfin.ItemSortBy(options.Sort)}
		j.SortBy = &sortBy

		// SortOrder
		if options.SortOrder == "desc" {
			log.Debug().Msg("Applying descending sort order")
			j.SortOrder = &[]jellyfin.SortOrder{jellyfin.SORTORDER_DESCENDING}
		} else {
			log.Debug().Msg("Applying ascending sort order")
			j.SortOrder = &[]jellyfin.SortOrder{jellyfin.SORTORDER_ASCENDING}
		}
	}

	// Search term
	if options.Query != "" {
		log.Debug().Str("query", options.Query).Msg("Applying search term filter")
		j.SearchTerm = &options.Query

		// Enable recursive search when searching
		recursive := true
		j.Recursive = &recursive

		// Increase limit for search results if not explicitly set
		if options.Limit <= 0 && j.Limit == nil {
			defaultLimit := int32(50)
			j.Limit = &defaultLimit
			log.Debug().Int32("defaultLimit", defaultLimit).Msg("Applying default search limit")
		}
	}

	// MediaType filter
	if options.MediaType != "" {
		log.Debug().Str("mediaType", string(options.MediaType)).Msg("Applying media type filter")
		var includeItemTypes []jellyfin.BaseItemKind

		switch options.MediaType {
		case types.MediaTypeCollection:
			includeItemTypes = append(includeItemTypes, jellyfin.BASEITEMKIND_COLLECTION_FOLDER, jellyfin.BASEITEMKIND_BOX_SET)
		case types.MediaTypeMovie:
			includeItemTypes = append(includeItemTypes, jellyfin.BASEITEMKIND_MOVIE)
		case types.MediaTypeSeries:
			includeItemTypes = append(includeItemTypes, jellyfin.BASEITEMKIND_SERIES)
		case types.MediaTypeSeason:
			includeItemTypes = append(includeItemTypes, jellyfin.BASEITEMKIND_SEASON)
		case types.MediaTypeEpisode:
			includeItemTypes = append(includeItemTypes, jellyfin.BASEITEMKIND_EPISODE)
		case types.MediaTypeArtist:
			includeItemTypes = append(includeItemTypes, jellyfin.BASEITEMKIND_MUSIC_ARTIST)
		case types.MediaTypeAlbum:
			includeItemTypes = append(includeItemTypes, jellyfin.BASEITEMKIND_MUSIC_ALBUM)
		case types.MediaTypeTrack:
			includeItemTypes = append(includeItemTypes, jellyfin.BASEITEMKIND_AUDIO)
		case types.MediaTypePlaylist:
			includeItemTypes = append(includeItemTypes, jellyfin.BASEITEMKIND_PLAYLIST)
		}

		if len(includeItemTypes) > 0 {
			j.IncludeItemTypes = &includeItemTypes
			log.Debug().Interface("includeItemTypes", includeItemTypes).Msg("Applying media type filter")
		}
	}

	// Genre filter
	if options.Genre != "" {
		log.Debug().Str("genre", options.Genre).Msg("Applying genre filter")
		genres := []string{options.Genre}
		j.Genres = &genres
	}

	// Favorite filter
	if options.Favorites {
		favorite := true
		log.Debug().Bool("favorite", favorite).Msg("Applying favorite filter")
		j.IsFavorite = &favorite
	}

	// Year filter
	if options.Year > 0 {
		year := int32(options.Year)
		years := []int32{year}
		j.Years = &years
		log.Debug().Int32("year", year).Msg("Applying year filter")
	}

	// Person filters
	if options.Actor != "" || options.Director != "" || options.Creator != "" {
		// Use the first non-empty person value
		log.Debug().
			Str("actor", options.Actor).
			Str("director", options.Director).
			Str("creator", options.Creator).
			Msg("Applying person filters")
		var person string
		if options.Actor != "" {
			person = options.Actor
		} else if options.Director != "" {
			person = options.Director
		} else if options.Creator != "" {
			person = options.Creator
		}

		if person != "" {
			log.Debug().Str("person", person).Msg("Applying person filter")
			j.Person = &person
		}
	}

	// Content rating filter
	if options.ContentRating != "" {
		log.Debug().Str("contentRating", options.ContentRating).Msg("Applying content rating filter")
		ratings := []string{options.ContentRating}
		j.OfficialRatings = &ratings
	}

	// Tags filter
	if len(options.Tags) > 0 {
		log.Debug().Strs("tags", options.Tags).Msg("Applying tags filter")
		j.Tags = &options.Tags
	}

	// Recently added filter
	if options.RecentlyAdded {
		log.Debug().Msg("Applying recently added filter")
		sortBy := []jellyfin.ItemSortBy{jellyfin.ITEMSORTBY_DATE_CREATED, jellyfin.ITEMSORTBY_SORT_NAME}
		j.SortBy = &sortBy

		sortOrder := []jellyfin.SortOrder{jellyfin.SORTORDER_DESCENDING}
		j.SortOrder = &sortOrder
	}

	// Recently played filter
	if options.RecentlyPlayed {
		log.Debug().Msg("Applying recently played filter")
		sortBy := []jellyfin.ItemSortBy{jellyfin.ITEMSORTBY_DATE_PLAYED, jellyfin.ITEMSORTBY_SORT_NAME}
		j.SortBy = &sortBy

		sortOrder := []jellyfin.SortOrder{jellyfin.SORTORDER_DESCENDING}
		j.SortOrder = &sortOrder
	}

	// Watched filter
	if options.Watched {
		log.Debug().Msg("Applying watched filter")
		watched := true
		j.IsPlayed = &watched
	}

	// Date filters
	if options.DateAddedAfter != nil && !options.DateAddedAfter.IsZero() {
		log.Debug().Time("dateAddedAfter", *options.DateAddedAfter).Msg("Applying date added after filter")
		j.MinDateLastSaved = options.DateAddedAfter
	}

	if options.DateAddedBefore != nil && !options.DateAddedBefore.IsZero() {
		log.Debug().Time("dateAddedBefore", *options.DateAddedBefore).Msg("Applying date added before filter")
		j.MaxPremiereDate = options.DateAddedBefore
	}

	if options.ReleasedAfter != nil && !options.ReleasedAfter.IsZero() {
		log.Debug().Time("releasedAfter", *options.ReleasedAfter).Msg("Applying released after filter")
		j.MinPremiereDate = options.ReleasedAfter
	}

	if options.ReleasedBefore != nil && !options.ReleasedBefore.IsZero() {
		log.Debug().Time("releasedBefore", *options.ReleasedBefore).Msg("Applying released before filter")
		j.MaxPremiereDate = options.ReleasedBefore
	}

	// Rating filter
	if options.MinimumRating > 0 {
		log.Debug().Float32("minimumRating", options.MinimumRating).Msg("Applying minimum rating filter")
		minRating := float64(options.MinimumRating)
		j.MinCommunityRating = &minRating
	}

	// Enable user data and images by default
	enableUserData := true
	j.EnableUserData = &enableUserData

	enableImages := true
	j.EnableImages = &enableImages

	enableTotalCount := true
	j.EnableTotalRecordCount = &enableTotalCount

	log.Debug().Msg("Successfully applied query options")

	return j
}

// ConvertItemParamsToGetItems creates a GetItems request from a GetItemsParams object
// func ConvertItemParamsToGetItems(ctx context.Context, params *jellyfin.ApiGetItemsRequest) *jellyfin.ApiGetItemsRequest {
// 	// Create new options and populate each field
// 	options := &JellyfinQueryOptions{}
//
// 	// Copy the values from params to options
// 	options.Limit = params.Limit
// 	options.StartIndex = params.StartIndex
// 	options.SortOrder = &params.SortOrder
// 	options.SortBy = &params.SortBy
// 	options.IncludeItemTypes = &params.IncludeItemTypes
// 	options.IsPlayed = params.IsPlayed
// 	options.IsFavorite = params.IsFavorite
// 	options.Recursive = params.Recursive
// 	options.MinCommunityRating = params.MinCommunityRating
//
// 	// Handle optional fields using IsSet() and Value() methods
// 	if params.SearchTerm.IsSet() {
// 		val := params.SearchTerm.Value()
// 		options.SearchTerm = &val
// 	}
//
// 	if params.Ids.IsSet() {
// 		ids := strings.Split(params.Ids.Value(), ",")
// 		options.Ids = &ids
// 	}
//
// 	if params.Person.IsSet() {
// 		val := params.Person.Value()
// 		options.Person = &val
// 	}
//
// 	// Convert to GetItemsRequest using the ToItemsRequest method
// 	return options.ToItemsRequest()
// }

// ApplyClientQueryOptions converts query options to Jellyfin API parameters
// func ApplyClientQueryOptions(queryParams *jellyfin.GetItemsParams, options *types.QueryOptions) {
// 	// Create a JellyfinQueryOptions from our standard options
// 	jellyfinOptions := NewJellyfinQueryOptions().FromQueryOptions(options)
//
// 	// ItemIDs
// 	if jellyfinOptions.Ids != nil && len(*jellyfinOptions.Ids) > 0 {
// 		queryParams.Ids(strings.Join(*jellyfinOptions.Ids, ","))
// 	}
//
// 	// Limit
// 	if jellyfinOptions.Limit != nil {
// 		queryParams.Limit(*jellyfinOptions.Limit)
// 	}
//
// 	// Offset/StartIndex
// 	if jellyfinOptions.StartIndex != nil {
// 		queryParams.StartIndex(*jellyfinOptions.StartIndex)
// 	}
//
// 	// Sort
// 	if jellyfinOptions.SortBy != nil {
// 		queryParams.SortBy(*jellyfinOptions.SortBy)
// 	}
//
// 	// SortOrder
// 	if jellyfinOptions.SortOrder != nil {
// 		queryParams.SortOrder(*jellyfinOptions.SortOrder)
// 	}
//
// 	// Search term
// 	if jellyfinOptions.SearchTerm != nil {
// 		queryParams.SearchTerm(*jellyfinOptions.SearchTerm)
// 	}
//
// 	// Recursive search
// 	if jellyfinOptions.Recursive != nil {
// 		queryParams.Recursive(*jellyfinOptions.Recursive)
// 	}
//
// 	// IncludeItemTypes
// 	if jellyfinOptions.IncludeItemTypes != nil {
// 		queryParams.IncludeItemTypes(*jellyfinOptions.IncludeItemTypes)
// 	}
//
// 	// Genre filter
// 	if jellyfinOptions.Genres != nil {
// 		queryParams.Genres(*jellyfinOptions.Genres)
// 	}
//
// 	// Favorite filter
// 	if jellyfinOptions.IsFavorite != nil {
// 		queryParams.IsFavorite(*jellyfinOptions.IsFavorite)
// 	}
//
// 	// Year filter
// 	if jellyfinOptions.Years != nil {
// 		queryParams.Years(*jellyfinOptions.Years)
// 	}
//
// 	// Person filter
// 	if jellyfinOptions.Person != nil {
// 		queryParams.Person(*jellyfinOptions.Person)
// 	}
//
// 	// Content rating filter
// 	if jellyfinOptions.OfficialRatings != nil {
// 		queryParams.OfficialRatings(*jellyfinOptions.OfficialRatings)
// 	}
//
// 	// Tags filter
// 	if jellyfinOptions.Tags != nil {
// 		queryParams.Tags(*jellyfinOptions.Tags)
// 	}
//
// 	// Watched filter
// 	if jellyfinOptions.IsPlayed != nil {
// 		queryParams.IsPlayed(*jellyfinOptions.IsPlayed)
// 	}
//
// 	// Date filters
// 	if jellyfinOptions.MinDateLastSaved != nil {
// 		queryParams.MinDateLastSaved(jellyfinOptions.MinDateLastSaved.Format(time.RFC3339))
// 	}
//
// 	if jellyfinOptions.MaxPremiereDate != nil {
// 		queryParams.MaxPremiereDate(jellyfinOptions.MaxPremiereDate.Format(time.RFC3339))
// 	}
//
// 	if jellyfinOptions.MinPremiereDate != nil {
// 		queryParams.MinPremiereDate(jellyfinOptions.MinPremiereDate.Format(time.RFC3339))
// 	}
//
// 	// Rating filter
// 	if jellyfinOptions.MinCommunityRating != nil {
// 		queryParams.MinCommunityRating(*jellyfinOptions.MinCommunityRating)
// 	}
//
// 	// Enable user data, images, record count
// 	if jellyfinOptions.EnableUserData != nil {
// 		queryParams.EnableUserData(*jellyfinOptions.EnableUserData)
// 	}
//
// 	if jellyfinOptions.EnableImages != nil {
// 		queryParams.EnableImages(*jellyfinOptions.EnableImages)
// 	}
//
// 	if jellyfinOptions.EnableTotalRecordCount != nil {
// 		queryParams.EnableTotalRecordCount(*jellyfinOptions.EnableTotalRecordCount)
// 	}
// }

// ApplyClientArtistOptions applies query options to an ArtistsRequest
// func ApplyClientArtistOptions(queryParams *jellyfin.GetArtistsRequest, options *types.QueryOptions) {
// 	// Create a new JellyfinQueryOptions from our standard options
// 	jellyfinOptions := NewJellyfinQueryOptions().FromQueryOptions(options)
//
// 	// Apply all parameters using method-based approach
// 	if jellyfinOptions.IncludeItemTypes != nil {
// 		queryParams.IncludeItemTypes(*jellyfinOptions.IncludeItemTypes)
// 	}
//
// 	if jellyfinOptions.UserId != nil {
// 		queryParams.UserId(*jellyfinOptions.UserId)
// 	}
//
// 	if jellyfinOptions.SortBy != nil {
// 		queryParams.SortBy(*jellyfinOptions.SortBy)
// 	}
//
// 	if jellyfinOptions.SortOrder != nil {
// 		queryParams.SortOrder(*jellyfinOptions.SortOrder)
// 	}
//
// 	if jellyfinOptions.Limit != nil {
// 		queryParams.Limit(*jellyfinOptions.Limit)
// 	}
//
// 	if jellyfinOptions.StartIndex != nil {
// 		queryParams.StartIndex(*jellyfinOptions.StartIndex)
// 	}
//
// 	if jellyfinOptions.SearchTerm != nil {
// 		queryParams.SearchTerm(*jellyfinOptions.SearchTerm)
// 	}
//
// 	if jellyfinOptions.Filters != nil {
// 		queryParams.Filters(*jellyfinOptions.Filters)
// 	}
//
// 	if jellyfinOptions.Fields != nil {
// 		queryParams.Fields(*jellyfinOptions.Fields)
// 	}
//
// 	if jellyfinOptions.ExcludeItemTypes != nil {
// 		queryParams.ExcludeItemTypes(*jellyfinOptions.ExcludeItemTypes)
// 	}
//
// 	if jellyfinOptions.Genres != nil {
// 		queryParams.Genres(*jellyfinOptions.Genres)
// 	}
//
// 	if jellyfinOptions.GenreIds != nil && len(*jellyfinOptions.GenreIds) > 0 {
// 		queryParams.GenreIds(strings.Join(*jellyfinOptions.GenreIds, ","))
// 	}
//
// 	if jellyfinOptions.OfficialRatings != nil {
// 		queryParams.OfficialRatings(*jellyfinOptions.OfficialRatings)
// 	}
//
// 	if jellyfinOptions.Tags != nil {
// 		queryParams.Tags(*jellyfinOptions.Tags)
// 	}
//
// 	if jellyfinOptions.Years != nil {
// 		queryParams.Years(*jellyfinOptions.Years)
// 	}
//
// 	if jellyfinOptions.EnableUserData != nil {
// 		queryParams.EnableUserData(*jellyfinOptions.EnableUserData)
// 	}
//
// 	if jellyfinOptions.IsFavorite != nil {
// 		queryParams.IsFavorite(*jellyfinOptions.IsFavorite)
// 	}
//
// 	if jellyfinOptions.MinCommunityRating != nil {
// 		queryParams.MinCommunityRating(*jellyfinOptions.MinCommunityRating)
// 	}
// }
