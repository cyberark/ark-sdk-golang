package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the Directories action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "directories",
		ActionDescription: "Identity management of directories.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
