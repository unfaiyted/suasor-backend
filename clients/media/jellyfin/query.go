package jellyfin

import (
	jellyfin "github.com/sj14/jellyfin-go/api"
	"strings"
	"suasor/clients/media/types"
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
func NewJellyfinQueryOptions(options *types.QueryOptions) *JellyfinQueryOptions {
	var jellyfinOptions *JellyfinQueryOptions

	// Create a new instance if none was provided
	if options == nil {
		jellyfinOptions = &JellyfinQueryOptions{}
	}

	jellyfinOptions.FromQueryOptions(options)

	return jellyfinOptions

}

// ToItemsRequest converts JellyfinQueryOptions to an API request object for GetItems
func (j *JellyfinQueryOptions) SetItemsRequest(req *jellyfin.ApiGetItemsRequest) {

	// Apply all parameters using method-based approach
	if j.UserId != nil {
		req.UserId(*j.UserId)
	}

	if j.MaxOfficialRating != nil {
		req.MaxOfficialRating(*j.MaxOfficialRating)
	}

	if j.HasThemeSong != nil {
		req.HasThemeSong(*j.HasThemeSong)
	}

	if j.HasThemeVideo != nil {
		req.HasThemeVideo(*j.HasThemeVideo)
	}

	if j.HasSubtitles != nil {
		req.HasSubtitles(*j.HasSubtitles)
	}

	if j.HasSpecialFeature != nil {
		req.HasSpecialFeature(*j.HasSpecialFeature)
	}

	if j.HasTrailer != nil {
		req.HasTrailer(*j.HasTrailer)
	}

	if j.AdjacentTo != nil {
		req.AdjacentTo(*j.AdjacentTo)
	}

	if j.IndexNumber != nil {
		req.IndexNumber(*j.IndexNumber)
	}

	if j.ParentIndexNumber != nil {
		req.ParentIndexNumber(*j.ParentIndexNumber)
	}

	if j.HasParentalRating != nil {
		req.HasParentalRating(*j.HasParentalRating)
	}

	if j.IsHd != nil {
		req.IsHd(*j.IsHd)
	}

	if j.Is4K != nil {
		req.Is4K(*j.Is4K)
	}

	if j.LocationTypes != nil {
		req.LocationTypes(*j.LocationTypes)
	}

	if j.ExcludeLocationTypes != nil {
		req.ExcludeLocationTypes(*j.ExcludeLocationTypes)
	}

	if j.IsMissing != nil {
		req.IsMissing(*j.IsMissing)
	}

	if j.IsUnaired != nil {
		req.IsUnaired(*j.IsUnaired)
	}

	if j.MinCommunityRating != nil {
		req.MinCommunityRating(*j.MinCommunityRating)
	}

	if j.MinCriticRating != nil {
		req.MinCriticRating(*j.MinCriticRating)
	}

	if j.MinPremiereDate != nil {
		req.MinPremiereDate(*j.MinPremiereDate)
	}

	if j.MinDateLastSaved != nil {
		req.MinDateLastSaved(*j.MinDateLastSaved)
	}

	if j.MinDateLastSavedForUser != nil {
		req.MinDateLastSavedForUser(*j.MinDateLastSavedForUser)
	}

	if j.MaxPremiereDate != nil {
		req.MaxPremiereDate(*j.MaxPremiereDate)
	}

	if j.HasOverview != nil {
		req.HasOverview(*j.HasOverview)
	}

	if j.HasImdbId != nil {
		req.HasImdbId(*j.HasImdbId)
	}

	if j.HasTmdbId != nil {
		req.HasTmdbId(*j.HasTmdbId)
	}

	if j.HasTvdbId != nil {
		req.HasTvdbId(*j.HasTvdbId)
	}

	if j.IsMovie != nil {
		req.IsMovie(*j.IsMovie)
	}

	if j.IsSeries != nil {
		req.IsSeries(*j.IsSeries)
	}

	if j.IsNews != nil {
		req.IsNews(*j.IsNews)
	}

	if j.IsKids != nil {
		req.IsKids(*j.IsKids)
	}

	if j.IsSports != nil {
		req.IsSports(*j.IsSports)
	}

	if j.ExcludeItemIds != nil && len(*j.ExcludeItemIds) > 0 {
		req.ExcludeItemIds(*j.ExcludeItemIds)
	}

	if j.StartIndex != nil {
		req.StartIndex(*j.StartIndex)
	}

	if j.Limit != nil {
		req.Limit(*j.Limit)
	}

	if j.Recursive != nil {
		req.Recursive(*j.Recursive)
	}

	if j.SearchTerm != nil {
		req.SearchTerm(*j.SearchTerm)
	}

	if j.SortOrder != nil {
		req.SortOrder(*j.SortOrder)
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

	if j.ImageTypes != nil {
		req.ImageTypes(*j.ImageTypes)
	}

	if j.SortBy != nil {
		req.SortBy(*j.SortBy)
	}

	if j.IsPlayed != nil {
		req.IsPlayed(*j.IsPlayed)
	}

	if j.Genres != nil {
		req.Genres(*j.Genres)
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
func (j *JellyfinQueryOptions) FromQueryOptions(options *types.QueryOptions) *JellyfinQueryOptions {
	if options == nil {
		return nil
	}

	// Create a new instance if none was provided
	if j == nil {
		j = &JellyfinQueryOptions{}
	}

	// ItemIDs
	if options.ItemIDs != "" {
		ids := strings.Split(options.ItemIDs, ",")
		j.Ids = &ids
	}

	// Limit
	if options.Limit > 0 {
		limit := int32(options.Limit)
		j.Limit = &limit
	}

	// Offset/StartIndex
	if options.Offset > 0 {
		startIndex := int32(options.Offset)
		j.StartIndex = &startIndex
	}

	// Sort
	if options.Sort != "" {
		sortBy := []jellyfin.ItemSortBy{jellyfin.ItemSortBy(options.Sort)}
		j.SortBy = &sortBy

		// SortOrder
		if options.SortOrder == "desc" {
			j.SortOrder = &[]jellyfin.SortOrder{jellyfin.SORTORDER_DESCENDING}
		} else {
			j.SortOrder = &[]jellyfin.SortOrder{jellyfin.SORTORDER_ASCENDING}
		}
	}

	// Search term
	if options.Query != "" {
		j.SearchTerm = &options.Query

		// Enable recursive search when searching
		recursive := true
		j.Recursive = &recursive

		// Increase limit for search results if not explicitly set
		if options.Limit <= 0 && j.Limit == nil {
			defaultLimit := int32(50)
			j.Limit = &defaultLimit
		}
	}

	// MediaType filter
	if options.MediaType != "" {
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
		}
	}

	// Genre filter
	if options.Genre != "" {
		genres := []string{options.Genre}
		j.Genres = &genres
	}

	// Favorite filter
	if options.Favorites {
		favorite := true
		j.IsFavorite = &favorite
	}

	// Year filter
	if options.Year > 0 {
		year := int32(options.Year)
		years := []int32{year}
		j.Years = &years
	}

	// Person filters
	if options.Actor != "" || options.Director != "" || options.Creator != "" {
		// Use the first non-empty person value
		var person string
		if options.Actor != "" {
			person = options.Actor
		} else if options.Director != "" {
			person = options.Director
		} else if options.Creator != "" {
			person = options.Creator
		}

		if person != "" {
			j.Person = &person
		}
	}

	// Content rating filter
	if options.ContentRating != "" {
		ratings := []string{options.ContentRating}
		j.OfficialRatings = &ratings
	}

	// Tags filter
	if len(options.Tags) > 0 {
		j.Tags = &options.Tags
	}

	// Recently added filter
	if options.RecentlyAdded {
		sortBy := []jellyfin.ItemSortBy{jellyfin.ITEMSORTBY_DATE_CREATED, jellyfin.ITEMSORTBY_SORT_NAME}
		j.SortBy = &sortBy

		sortOrder := []jellyfin.SortOrder{jellyfin.SORTORDER_DESCENDING}
		j.SortOrder = &sortOrder
	}

	// Recently played filter
	if options.RecentlyPlayed {
		sortBy := []jellyfin.ItemSortBy{jellyfin.ITEMSORTBY_DATE_PLAYED, jellyfin.ITEMSORTBY_SORT_NAME}
		j.SortBy = &sortBy

		sortOrder := []jellyfin.SortOrder{jellyfin.SORTORDER_DESCENDING}
		j.SortOrder = &sortOrder
	}

	// Watched filter
	if options.Watched {
		watched := true
		j.IsPlayed = &watched
	}

	// Date filters
	if !options.DateAddedAfter.IsZero() {
		j.MinDateLastSaved = options.DateAddedAfter
	}

	if !options.DateAddedBefore.IsZero() {
		j.MaxPremiereDate = options.DateAddedBefore
	}

	if !options.ReleasedAfter.IsZero() {
		j.MinPremiereDate = options.ReleasedAfter
	}

	if !options.ReleasedBefore.IsZero() {
		j.MaxPremiereDate = options.ReleasedBefore
	}

	// Rating filter
	if options.MinimumRating > 0 {
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
