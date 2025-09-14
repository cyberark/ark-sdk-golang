package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the SIA Secrets DB action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "db",
		ActionDescription: "SIA Secrets DB Actions.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
