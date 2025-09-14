package workspaces

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siaworkspacesdbactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db/actions"
	siaworkspacestargetsetsactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/targetsets/actions"
)

// CLIAction is a struct that defines the SIA Workspaces action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "workspaces",
		ActionDescription: "SIA Workspaces Actions",
		ActionVersion:     1,
	},
	Subactions: []*actions.ArkServiceCLIActionDefinition{
		siaworkspacestargetsetsactions.CLIAction,
		siaworkspacesdbactions.CLIAction,
	},
}

// ServiceConfig is the configuration for the sia workspaces services.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-workspaces",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			CLIAction,
		},
	},
}

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
