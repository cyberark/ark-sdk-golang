package models

// ArkPCloudGetAccount represents the details required to retrieve an account.
type ArkPCloudGetAccount struct {
	AccountID string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to retrieve" flag:"account-id" validate:"required"`
}
