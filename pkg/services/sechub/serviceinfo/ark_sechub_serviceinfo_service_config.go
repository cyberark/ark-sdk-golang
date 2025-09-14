package serviceinfo

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	sechubserviceinfoactions "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/serviceinfo/actions"
)

// ServiceConfig is the configuration for the Secrets Hub Service Info service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sechub-serviceinfo",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			sechubserviceinfoactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SecHub Service Info service.
var ServiceGenerator = NewArkSecHubServiceInfoService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
