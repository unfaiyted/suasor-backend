// helpers.go
package emby

import (
	"github.com/antihax/optional"
	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
)

// applyQueryOptions applies the common query options to Emby API parameters
func applyQueryOptions(queryParams *embyclient.ItemsServiceApiGetItemsOpts, options *types.QueryOptions) {
	if options == nil {
		return
	}

	if options.Limit > 0 {
		queryParams.Limit = optional.NewInt32(int32(options.Limit))
	}

	if options.Offset > 0 {
		queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
	}

	if options.Sort != "" {
		queryParams.SortBy = optional.NewString(options.Sort)
		if options.SortOrder == "desc" {
			queryParams.SortOrder = optional.NewString("Descending")
		} else {
			queryParams.SortOrder = optional.NewString("Ascending")
		}
	}
}
