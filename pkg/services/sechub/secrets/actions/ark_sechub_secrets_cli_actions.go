package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the secrets action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "secrets",
		ActionDescription: "Sechub Secrets.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
