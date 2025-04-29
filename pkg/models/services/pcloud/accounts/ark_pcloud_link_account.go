package accounts

// ArkPCloudLinkAccount represents the details required to link an account.
type ArkPCloudLinkAccount struct {
	AccountID          string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to link" flag:"account-id" validate:"required"`
	Safe               string `json:"safe" mapstructure:"safe" desc:"The safe in which the linked account is stored" flag:"safe" validate:"required"`
	ExtraPasswordIndex int    `json:"extra_password_index" mapstructure:"extra_password_index" desc:"The linked account extra password index" flag:"extra-password-index" validate:"required"`
	Folder             string `json:"folder" mapstructure:"folder" desc:"Folder of the linked account" flag:"folder" validate:"required"`
	Name               string `json:"name" mapstructure:"name" desc:"The linked account name" flag:"name" validate:"required"`
}
