package actions

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// CLIAction is a struct that defines the service info action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "service-info",
		ActionDescription: "Sechub Service Info.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
}
