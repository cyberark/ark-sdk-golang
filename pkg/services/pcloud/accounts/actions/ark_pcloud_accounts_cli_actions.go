package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the accounts action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "accounts",
		ActionDescription: "PCloud Accounts Management.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
