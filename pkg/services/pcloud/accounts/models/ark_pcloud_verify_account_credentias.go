package models

// ArkPCloudVerifyAccountCredentials represents the details required to verify account credentials.
type ArkPCloudVerifyAccountCredentials struct {
	AccountID string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to mark for validation" flag:"account-id" validate:"required"`
}
