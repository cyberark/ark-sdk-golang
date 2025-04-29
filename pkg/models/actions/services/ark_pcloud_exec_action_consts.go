package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	pcloudaccounts "github.com/cyberark/ark-sdk-golang/pkg/models/services/pcloud/accounts"
	pcloudsafes "github.com/cyberark/ark-sdk-golang/pkg/models/services/pcloud/safes"
)

// AccountsActionToSchemaMap is a map that defines the mapping between Ark PCloud action names and their corresponding schema types.
var AccountsActionToSchemaMap = map[string]interface{}{
	"add-account":                         &pcloudaccounts.ArkPCloudAddAccount{},
	"update-account":                      &pcloudaccounts.ArkPCloudUpdateAccount{},
	"delete-account":                      &pcloudaccounts.ArkPCloudDeleteAccount{},
	"account":                             &pcloudaccounts.ArkPCloudGetAccount{},
	"account-credentials":                 &pcloudaccounts.ArkPCloudGetAccountCredentials{},
	"list-accounts":                       nil,
	"list-accounts-by":                    &pcloudaccounts.ArkPCloudAccountsFilter{},
	"list-account-secret-versions":        &pcloudaccounts.ArkPCloudListAccountSecretVersions{},
	"generate-account-credentials":        &pcloudaccounts.ArkPCloudGenerateAccountCredentials{},
	"verify-account-credentials":          &pcloudaccounts.ArkPCloudVerifyAccountCredentials{},
	"change-account-credentials":          &pcloudaccounts.ArkPCloudChangeAccountCredentials{},
	"set-account-next-credentials":        &pcloudaccounts.ArkPCloudSetAccountNextCredentials{},
	"update-account-credentials-in-vault": &pcloudaccounts.ArkPCloudUpdateAccountCredentialsInVault{},
	"reconcile-account-credentials":       &pcloudaccounts.ArkPCloudReconcileAccountCredentials{},
	"link-account":                        &pcloudaccounts.ArkPCloudLinkAccount{},
	"unlink-account":                      &pcloudaccounts.ArkPCloudUnlinkAccount{},
	"accounts-stats":                      nil,
}

// AccountsAction is a struct that defines the Ark PCloud action for the Ark service.
var AccountsAction = actions.ArkServiceActionDefinition{
	ActionName: "accounts",
	Schemas:    AccountsActionToSchemaMap,
}

// SafesActionToSchemaMap is a map that defines the mapping between Ark PCloud action names and their corresponding schema types.
var SafesActionToSchemaMap = map[string]interface{}{
	"add-safe":             &pcloudsafes.ArkPCloudAddSafe{},
	"update-safe":          &pcloudsafes.ArkPCloudUpdateSafe{},
	"delete-safe":          &pcloudsafes.ArkPCloudDeleteSafe{},
	"safe":                 &pcloudsafes.ArkPCloudGetSafe{},
	"list-safes":           nil,
	"list-safes-by":        &pcloudsafes.ArkPCloudSafesFilters{},
	"safes-stats":          nil,
	"add-safe-member":      &pcloudsafes.ArkPCloudAddSafeMember{},
	"update-safe-member":   &pcloudsafes.ArkPCloudUpdateSafeMember{},
	"delete-safe-member":   &pcloudsafes.ArkPCloudDeleteSafeMember{},
	"safe-member":          &pcloudsafes.ArkPCloudGetSafeMember{},
	"list-safe-members":    &pcloudsafes.ArkPCloudListSafeMembers{},
	"list-safe-members-by": &pcloudsafes.ArkPCloudSafeMembersFilters{},
	"safe-members-stats":   &pcloudsafes.ArkPCloudGetSafeMembersStats{},
	"safes-members-stats":  nil,
}

// SafesAction is a struct that defines the Ark PCloud action for the Ark service.
var SafesAction = actions.ArkServiceActionDefinition{
	ActionName: "safes",
	Schemas:    SafesActionToSchemaMap,
}

// PCloudActions is a struct that defines the Ark PCloud action for the Ark service.
var PCloudActions = &actions.ArkServiceActionDefinition{
	ActionName: "pcloud",
	Subactions: []*actions.ArkServiceActionDefinition{
		&AccountsAction,
		&SafesAction,
	},
}
