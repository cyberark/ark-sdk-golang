package actions

import accountsmodels "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/accounts/models"

// ActionToSchemaMap is a map that defines the mapping between Ark PCloud action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"add-account":                         &accountsmodels.ArkPCloudAddAccount{},
	"update-account":                      &accountsmodels.ArkPCloudUpdateAccount{},
	"delete-account":                      &accountsmodels.ArkPCloudDeleteAccount{},
	"account":                             &accountsmodels.ArkPCloudGetAccount{},
	"account-credentials":                 &accountsmodels.ArkPCloudGetAccountCredentials{},
	"list-accounts":                       nil,
	"list-accounts-by":                    &accountsmodels.ArkPCloudAccountsFilter{},
	"list-account-secret-versions":        &accountsmodels.ArkPCloudListAccountSecretVersions{},
	"generate-account-credentials":        &accountsmodels.ArkPCloudGenerateAccountCredentials{},
	"verify-account-credentials":          &accountsmodels.ArkPCloudVerifyAccountCredentials{},
	"change-account-credentials":          &accountsmodels.ArkPCloudChangeAccountCredentials{},
	"set-account-next-credentials":        &accountsmodels.ArkPCloudSetAccountNextCredentials{},
	"update-account-credentials-in-vault": &accountsmodels.ArkPCloudUpdateAccountCredentialsInVault{},
	"reconcile-account-credentials":       &accountsmodels.ArkPCloudReconcileAccountCredentials{},
	"link-account":                        &accountsmodels.ArkPCloudLinkAccount{},
	"unlink-account":                      &accountsmodels.ArkPCloudUnlinkAccount{},
	"accounts-stats":                      nil,
}
