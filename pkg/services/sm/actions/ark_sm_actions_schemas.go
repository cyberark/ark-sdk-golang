package actions

import smmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sm/models"

// ActionToSchemaMap is a map that defines the mapping between SM action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"list-sessions":               nil,
	"count-sessions":              nil,
	"list-sessions-by":            &smmodels.ArkSMSessionsFilter{},
	"count-sessions-by":           &smmodels.ArkSMSessionsFilter{},
	"session":                     &smmodels.ArkSIASMGetSession{},
	"list-session-activities":     &smmodels.ArkSIASMGetSessionActivities{},
	"count-session-activities":    &smmodels.ArkSIASMGetSessionActivities{},
	"list-session-activities-by":  &smmodels.ArkSMSessionActivitiesFilter{},
	"count-session-activities-by": &smmodels.ArkSMSessionActivitiesFilter{},
	"sessions-stats":              nil,
}
