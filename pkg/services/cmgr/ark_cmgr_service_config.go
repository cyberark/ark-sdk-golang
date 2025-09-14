package cmgr

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	cmgractions "github.com/cyberark/ark-sdk-golang/pkg/services/cmgr/actions"
)

// ServiceConfig is the configuration for the connector management service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "cmgr",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			cmgractions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			cmgractions.TerraformActionNetworkResource,
			cmgractions.TerraformActionPoolResource,
			cmgractions.TerraformActionPoolIdentifierResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			cmgractions.TerraformActionNetworkDataSource,
			cmgractions.TerraformActionPoolDataSource,
			cmgractions.TerraformActionPoolIdentifierDataSource,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the ArkCmgrService.
var ServiceGenerator = NewArkCmgrService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, true)
	if err != nil {
		panic(err)
	}
}
