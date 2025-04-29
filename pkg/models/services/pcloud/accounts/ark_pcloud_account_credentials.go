package accounts

// ArkPCloudAccountCredentials represents the credentials of an account.
type ArkPCloudAccountCredentials struct {
	AccountID string `json:"account_id" mapstructure:"account_id" desc:"The id of the account" flag:"account-id" validate:"required"`
	Password  string `json:"password" mapstructure:"password" desc:"The credentials" flag:"password" validate:"required"`
}
