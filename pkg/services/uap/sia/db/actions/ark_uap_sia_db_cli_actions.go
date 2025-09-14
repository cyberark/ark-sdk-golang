package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the uap sia db action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "db",
		ActionDescription: "UAP SIA DB Policies Management.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
