package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the CMGR action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "cmgr",
		ActionDescription: "Connector Management mediates ISPSS services and is used by IT administrators to manage CyberArk components, communication tunnels and manage networks and pools.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
	ActionAliases: []string{"connectormanager", "cm"},
}
