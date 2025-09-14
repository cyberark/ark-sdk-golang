package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the SIA Secrets VM action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "vm",
		ActionDescription: "SIA Secrets VM Actions.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
