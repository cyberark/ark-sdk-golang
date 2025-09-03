package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	accountsmodels "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/accounts/models"
	safesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/safes/models"
)

// AccountsActionToSchemaMap is a map that defines the mapping between Ark PCloud action names and their corresponding schema types.
var AccountsActionToSchemaMap = map[string]interface{}{
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

// AccountsAction is a struct that defines the Ark PCloud action for the Ark service.
var AccountsAction = actions.ArkServiceActionDefinition{
	ActionName: "accounts",
	Schemas:    AccountsActionToSchemaMap,
}

// SafesActionToSchemaMap is a map that defines the mapping between Ark PCloud action names and their corresponding schema types.
var SafesActionToSchemaMap = map[string]interface{}{
	"add-safe":             &safesmodels.ArkPCloudAddSafe{},
	"update-safe":          &safesmodels.ArkPCloudUpdateSafe{},
	"delete-safe":          &safesmodels.ArkPCloudDeleteSafe{},
	"safe":                 &safesmodels.ArkPCloudGetSafe{},
	"list-safes":           nil,
	"list-safes-by":        &safesmodels.ArkPCloudSafesFilters{},
	"safes-stats":          nil,
	"add-safe-member":      &safesmodels.ArkPCloudAddSafeMember{},
	"update-safe-member":   &safesmodels.ArkPCloudUpdateSafeMember{},
	"delete-safe-member":   &safesmodels.ArkPCloudDeleteSafeMember{},
	"safe-member":          &safesmodels.ArkPCloudGetSafeMember{},
	"list-safe-members":    &safesmodels.ArkPCloudListSafeMembers{},
	"list-safe-members-by": &safesmodels.ArkPCloudSafeMembersFilters{},
	"safe-members-stats":   &safesmodels.ArkPCloudGetSafeMembersStats{},
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
