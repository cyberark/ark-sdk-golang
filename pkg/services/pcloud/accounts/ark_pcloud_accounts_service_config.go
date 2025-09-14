package accounts

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	pcloudaccountsactions "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/accounts/actions"
)

// ServiceConfig is the configuration for the pcloud accounts service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "pcloud-accounts",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			pcloudaccountsactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			pcloudaccountsactions.TerraformActionAccountResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			pcloudaccountsactions.TerraformActionAccountDataSource,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the ArkPCloudAccountsService.
var ServiceGenerator = NewArkPCloudAccountsService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
