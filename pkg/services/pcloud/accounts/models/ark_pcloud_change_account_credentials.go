package models

// ArkPCloudChangeAccountCredentials represents the details required to change account credentials.
type ArkPCloudChangeAccountCredentials struct {
	AccountID string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to change the password for" flag:"account-id" validate:"required"`
}
