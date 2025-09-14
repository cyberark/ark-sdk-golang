package actions

import secretsvmmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/vm/models"

// ActionToSchemaMap is a map that defines the mapping between Secrets VM action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"add-secret":      &secretsvmmodels.ArkSIAVMAddSecret{},
	"change-secret":   &secretsvmmodels.ArkSIAVMChangeSecret{},
	"delete-secret":   &secretsvmmodels.ArkSIAVMDeleteSecret{},
	"list-secrets":    nil,
	"list-secrets-by": &secretsvmmodels.ArkSIAVMSecretsFilter{},
	"secret":          &secretsvmmodels.ArkSIAVMGetSecret{},
	"secrets-stats":   nil,
}
