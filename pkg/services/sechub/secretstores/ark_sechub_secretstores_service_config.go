package secretstores

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	sechubsecretstoresactions "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/secretstores/actions"
)

// ServiceConfig is the configuration for the Secrets Hub Secret Stores service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sechub-secretstores",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			sechubsecretstoresactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SecHub Secret Stores service.
var ServiceGenerator = NewArkSecHubSecretStoresService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
