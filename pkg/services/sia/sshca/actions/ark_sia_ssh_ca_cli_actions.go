package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the SIA SSH CA action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "ssh-ca",
		ActionDescription: "SIA SSH CA Actions.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
