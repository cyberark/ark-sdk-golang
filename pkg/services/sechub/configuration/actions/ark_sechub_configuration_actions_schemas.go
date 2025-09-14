package actions

import sechubconfiguration "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/configuration/models"

// ActionToSchemaMap is a map that defines the mapping between Sec Hub configuration action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"configuration":     nil,
	"set-configuration": &sechubconfiguration.ArkSecHubSetConfiguration{},
}
