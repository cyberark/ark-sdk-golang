package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the uap sia vm action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "vm",
		ActionDescription: "UAP SIA VM Policies Management.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
