package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the uap sca action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "sca",
		ActionDescription: "UAP SCA Policies Management.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
