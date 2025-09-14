package scans

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	sechubsscansactions "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/scans/actions"
)

// ServiceConfig is the configuration for the Secrets Hub scans service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sechub-scans",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			sechubsscansactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SecHub scans service.
var ServiceGenerator = NewArkSecHubScansService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
