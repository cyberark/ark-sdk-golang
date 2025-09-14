package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the SIA K8S action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "db",
		ActionDescription: "SIA K8S Actions.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
