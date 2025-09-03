package models

// ArkPCloudUnlinkAccount represents the details required to unlink an account.
type ArkPCloudUnlinkAccount struct {
	AccountID          string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to unlink" flag:"account-id" validate:"required"`
	ExtraPasswordIndex string `json:"extra_password_index" mapstructure:"extra_password_index" desc:"The linked account extra password index" flag:"extra-password-index" validate:"required"`
}
