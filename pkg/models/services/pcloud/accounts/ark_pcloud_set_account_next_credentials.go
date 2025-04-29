package accounts

// ArkPCloudSetAccountNextCredentials represents the details required to set the next credentials for an account.
type ArkPCloudSetAccountNextCredentials struct {
	AccountID      string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to change the password for" flag:"account-id" validate:"required"`
	NewCredentials string `json:"new_credentials" mapstructure:"new_credentials" desc:"Next credentials to set" flag:"new-credentials" validate:"required"`
}
