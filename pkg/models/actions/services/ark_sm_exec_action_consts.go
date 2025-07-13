package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/models/services/sm"
)

// SMActionToSchemaMap is a map that defines the mapping between SM action names and their corresponding schema types.
var SMActionToSchemaMap = map[string]interface{}{
	"list-sessions":               nil,
	"count-sessions":              nil,
	"list-sessions-by":            &sm.ArkSMSessionsFilter{},
	"count-sessions-by":           &sm.ArkSMSessionsFilter{},
	"session":                     &sm.ArkSIASMGetSession{},
	"list-session-activities":     &sm.ArkSIASMGetSessionActivities{},
	"count-session-activities":    &sm.ArkSIASMGetSessionActivities{},
	"list-session-activities-by":  &sm.ArkSMSessionActivitiesFilter{},
	"count-session-activities-by": &sm.ArkSMSessionActivitiesFilter{},
	"sessions-stats":              nil,
}

// SMActions is a struct that defines the SM actions for the Ark service.
var SMActions = &actions.ArkServiceActionDefinition{
	ActionName: "sm",
	Schemas:    SMActionToSchemaMap,
}
