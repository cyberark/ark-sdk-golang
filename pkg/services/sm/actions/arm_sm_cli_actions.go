package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the SM action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "sm",
		ActionDescription: "The CyberArk Audit space centralizes session monitoring across all CyberArk services on the Shared Services platform to provide a comprehensive display of all sessions as a unified view. This enables enhanced auditing as well as incident-response investigation.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
	ActionAliases: []string{"sessionmonitoring"},
}
