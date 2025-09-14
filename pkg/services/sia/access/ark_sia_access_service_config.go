package access

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siaaccessactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/access/actions"
)

// ServiceConfig is the configuration for the ArkSIAAccessService.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-access",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			siaaccessactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			siaaccessactions.TerraformActionAccessConnectorResource,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SIA Access service.
var ServiceGenerator = NewArkSIAAccessService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
