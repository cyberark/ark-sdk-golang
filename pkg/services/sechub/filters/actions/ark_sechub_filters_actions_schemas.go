package actions

import filtersmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/filters/models"

// ActionToSchemaMap is a map that defines the mapping between Sec Hub filters action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"filter":        &filtersmodels.ArkSecHubGetFilter{},
	"list-filters":  &filtersmodels.ArkSecHubGetFilters{},
	"add-filter":    &filtersmodels.ArkSecHubAddFilter{},
	"delete-filter": &filtersmodels.ArkSecHubDeleteFilter{},
}
