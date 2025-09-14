package configuration

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	sechubconfigurationactions "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/configuration/actions"
)

// ServiceConfig is the configuration for the Secrets Hub Configuration service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sechub-configuration",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			sechubconfigurationactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SecHub Configuration service.
var ServiceGenerator = NewArkSecHubConfigurationService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
