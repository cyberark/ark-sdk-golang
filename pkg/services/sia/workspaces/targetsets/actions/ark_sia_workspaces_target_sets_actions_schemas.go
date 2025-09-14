package actions

import targetsetsmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/targetsets/models"

// ActionToSchemaMap is a map that defines the mapping between TargetSets action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"add-target-set":          &targetsetsmodels.ArkSIAAddTargetSet{},
	"bulk-add-target-sets":    &targetsetsmodels.ArkSIABulkAddTargetSets{},
	"delete-target-set":       &targetsetsmodels.ArkSIADeleteTargetSet{},
	"bulk-delete-target-sets": &targetsetsmodels.ArkSIABulkDeleteTargetSets{},
	"update-target-set":       &targetsetsmodels.ArkSIAUpdateTargetSet{},
	"list-target-sets":        nil,
	"list-target-sets-by":     &targetsetsmodels.ArkSIATargetSetsFilter{},
	"target-set":              &targetsetsmodels.ArkSIAGetTargetSet{},
	"target-sets-stats":       nil,
}
