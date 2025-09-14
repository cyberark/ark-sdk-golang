package actions

import safesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/safes/models"

// ActionToSchemaMap is a map that defines the mapping between Ark PCloud action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"add-safe":             &safesmodels.ArkPCloudAddSafe{},
	"update-safe":          &safesmodels.ArkPCloudUpdateSafe{},
	"delete-safe":          &safesmodels.ArkPCloudDeleteSafe{},
	"safe":                 &safesmodels.ArkPCloudGetSafe{},
	"list-safes":           nil,
	"list-safes-by":        &safesmodels.ArkPCloudSafesFilters{},
	"safes-stats":          nil,
	"add-safe-member":      &safesmodels.ArkPCloudAddSafeMember{},
	"update-safe-member":   &safesmodels.ArkPCloudUpdateSafeMember{},
	"delete-safe-member":   &safesmodels.ArkPCloudDeleteSafeMember{},
	"safe-member":          &safesmodels.ArkPCloudGetSafeMember{},
	"list-safe-members":    &safesmodels.ArkPCloudListSafeMembers{},
	"list-safe-members-by": &safesmodels.ArkPCloudSafeMembersFilters{},
	"safe-members-stats":   &safesmodels.ArkPCloudGetSafeMembersStats{},
	"safes-members-stats":  nil,
}
