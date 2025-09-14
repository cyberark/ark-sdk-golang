package secrets

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siasecretsdbactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/db/actions"
	siasecretsvmactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/vm/actions"
)

// CLIAction is a struct that defines the SIA Secrets action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "secrets",
		ActionDescription: "SIA Secrets Actions",
		ActionVersion:     1,
	},
	Subactions: []*actions.ArkServiceCLIActionDefinition{
		siasecretsvmactions.CLIAction,
		siasecretsdbactions.CLIAction,
	},
}

// ServiceConfig is the configuration for the sia secrets services.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-secrets",
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
