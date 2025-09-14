package targetsets

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siaworkspacestargetsetsactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/targetsets/actions"
)

// ServiceConfig is the configuration for the SIA target sets workspace service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-workspaces-target-sets",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			siaworkspacestargetsetsactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			siaworkspacestargetsetsactions.TerraformActionWorkspacesTargetSetsResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			siaworkspacestargetsetsactions.TerraformActionWorkspacesTargetSetsDataSource,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the ArkSIAWorkspacesTargetSetsService.
var ServiceGenerator = NewArkSIAWorkspacesTargetSetsService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
