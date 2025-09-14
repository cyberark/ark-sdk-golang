package identity

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	identitydirectoriesactions "github.com/cyberark/ark-sdk-golang/pkg/services/identity/directories/actions"
	identityrolesactions "github.com/cyberark/ark-sdk-golang/pkg/services/identity/roles/actions"
	identityusersactions "github.com/cyberark/ark-sdk-golang/pkg/services/identity/users/actions"
)

// CLIAction is a struct that defines the identity action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "identity",
		ActionDescription: "Identity provides a single centralized interface for provisioning users and setting up the authentication for users of the Shared Services platform.",
		ActionVersion:     1,
	},
	ActionAliases: []string{"idaptive", "id"},
	Subactions: []*actions.ArkServiceCLIActionDefinition{
		identitydirectoriesactions.CLIAction,
		identityrolesactions.CLIAction,
		identityusersactions.CLIAction,
	},
}

// ServiceConfig is the configuration for the identity service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "identity",
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
	err := services.Register(ServiceConfig, true)
	if err != nil {
		panic(err)
	}
}
