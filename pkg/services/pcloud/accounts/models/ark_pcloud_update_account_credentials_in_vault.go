package models

// ArkPCloudUpdateAccountCredentialsInVault represents the details required to update account credentials in the vault.
type ArkPCloudUpdateAccountCredentialsInVault struct {
	AccountID      string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to change the password for" flag:"account-id" validate:"required"`
	NewCredentials string `json:"new_credentials" mapstructure:"new_credentials" desc:"New credentials to set in vault" flag:"new-credentials" validate:"required"`
}
