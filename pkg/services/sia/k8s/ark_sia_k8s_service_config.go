package k8s

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siak8sactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/k8s/actions"
)

// ServiceConfig is the configuration for the ArkSIAK8SService.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-k8s",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			siak8sactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SIA K8S service.
var ServiceGenerator = NewArkSIAK8SService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
