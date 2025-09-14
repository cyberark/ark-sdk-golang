package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the SIA Access action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "access",
		ActionDescription: "SIA Access Connectors.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
