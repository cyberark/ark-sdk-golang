package actions

import secretsdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/db/models"

// ActionToSchemaMap is a map that defines the mapping between Secrets DB action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"add-secret":      &secretsdbmodels.ArkSIADBAddSecret{},
	"update-secret":   &secretsdbmodels.ArkSIADBUpdateSecret{},
	"delete-secret":   &secretsdbmodels.ArkSIADBDeleteSecret{},
	"list-secrets":    nil,
	"list-secrets-by": &secretsdbmodels.ArkSIADBSecretsFilter{},
	"enable-secret":   &secretsdbmodels.ArkSIADBEnableSecret{},
	"disable-secret":  &secretsdbmodels.ArkSIADBDisableSecret{},
	"secret":          &secretsdbmodels.ArkSIADBGetSecret{},
	"secrets-stats":   nil,
}
