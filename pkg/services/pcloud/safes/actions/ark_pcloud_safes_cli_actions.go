package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the safes action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "safes",
		ActionDescription: "PCloud Safes Management.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
