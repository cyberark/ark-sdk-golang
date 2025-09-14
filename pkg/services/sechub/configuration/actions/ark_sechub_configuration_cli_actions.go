package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the configuration action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "configuration",
		ActionDescription: "Sechub Configuration.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
