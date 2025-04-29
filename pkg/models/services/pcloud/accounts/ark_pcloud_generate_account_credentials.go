package accounts

// ArkPCloudGenerateAccountCredentials represents the details required to generate account credentials.
type ArkPCloudGenerateAccountCredentials struct {
	AccountID string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to generate the password for" flag:"account-id" validate:"required"`
}
