package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	smmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sm/models"
)

// SMActionToSchemaMap is a map that defines the mapping between SM action names and their corresponding schema types.
var SMActionToSchemaMap = map[string]interface{}{
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

// SMActions is a struct that defines the SM actions for the Ark service.
var SMActions = &actions.ArkServiceActionDefinition{
	ActionName:        "sm",
	ActionAliases:     []string{"sessionmonitoring"},
	ActionDescription: "The CyberArk Audit space centralizes session monitoring across all CyberArk services on the Shared Services platform to provide a comprehensive display of all sessions as a unified view. This enables enhanced auditing as well as incident-response investigation.",
	Schemas:           SMActionToSchemaMap,
}
