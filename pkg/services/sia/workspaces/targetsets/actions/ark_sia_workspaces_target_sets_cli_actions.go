package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the SIA Workspace Target Sets action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "target-sets",
		ActionDescription: "SIA Workspaces Target Sets Actions.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
