package models

// ArkPCloudListAccountSecretVersions represents the details required to list account secret versions.
type ArkPCloudListAccountSecretVersions struct {
	AccountID     string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to retrieve the secret versions for" flag:"account-id" validate:"required"`
	ShowTemporary bool   `json:"show_temporary" mapstructure:"show_temporary" desc:"Show temporary secrets as well" flag:"show-temporary" default:"false"`
}
