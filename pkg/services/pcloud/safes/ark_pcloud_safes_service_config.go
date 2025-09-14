package safes

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	pcloudsafesactions "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/safes/actions"
)

// ServiceConfig is the configuration for the pcloud safes service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "pcloud-safes",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			pcloudsafesactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			pcloudsafesactions.TerraformActionSafeResource,
			pcloudsafesactions.TerraformActionSafeMemberResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			pcloudsafesactions.TerraformActionSafeDataSource,
			pcloudsafesactions.TerraformActionSafeMemberDataSource,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the ArkPCloudSafesService.
var ServiceGenerator = NewArkPCloudSafesService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
