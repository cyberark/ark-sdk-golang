package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the filters action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "filters",
		ActionDescription: "Sechub Filters.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
