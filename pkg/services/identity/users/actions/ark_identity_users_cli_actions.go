package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the roles action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "users",
		ActionDescription: "Identity management of users.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMapIdentityUsers,
	},
}
