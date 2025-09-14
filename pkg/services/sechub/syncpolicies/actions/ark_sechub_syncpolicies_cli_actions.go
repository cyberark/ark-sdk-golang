package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the sync policies action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "sync-policies",
		ActionDescription: "Sechub Sync Policies.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
