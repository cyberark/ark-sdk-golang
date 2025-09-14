package directories

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	identitydirectoriesactions "github.com/cyberark/ark-sdk-golang/pkg/services/identity/directories/actions"
)

// ServiceConfig is the configuration for the identity users service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "identity-directories",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			identitydirectoriesactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the ArkIdentityDirectoriesService.
var ServiceGenerator = NewArkIdentityDirectoriesService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
