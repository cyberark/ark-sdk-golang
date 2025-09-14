package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the secret stores action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "secret-stores",
		ActionDescription: "Sechub Secret Stores.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
