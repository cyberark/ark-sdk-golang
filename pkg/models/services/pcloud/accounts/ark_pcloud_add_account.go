package accounts

// ArkPCloudAddAccount represents the details required to add an account.
type ArkPCloudAddAccount struct {
	// Using inheritance on those for easier translation for CLI params
	ArkPCloudAccountSecretManagement     `mapstructure:",squash"`
	ArkPCloudAccountRemoteMachinesAccess `mapstructure:",squash"`
	Secret                               string                 `json:"secret" mapstructure:"secret" desc:"The secret of the account" flag:"secret" validate:"required"`
	Name                                 string                 `json:"name" mapstructure:"name,omitempty" desc:"Name of the account" flag:"name"`
	SafeName                             string                 `json:"safe_name" mapstructure:"safe_name" desc:"Safe name to store the account in" flag:"safe-name" validate:"required"`
	PlatformID                           string                 `json:"platform_id,omitempty" mapstructure:"platform_id,omitempty" desc:"Platform id to relate the account to" flag:"platform-id"`
	Username                             string                 `json:"username,omitempty" mapstructure:"username,omitempty" desc:"Username of the account" flag:"username"`
	Address                              string                 `json:"address,omitempty" mapstructure:"address,omitempty" desc:"Address of the account" flag:"address"`
	SecretType                           string                 `json:"secret_type,omitempty" mapstructure:"secret_type,omitempty" desc:"Type of the secret of the account (password,key)" flag:"secret-type" choices:"password,key"`
	PlatformAccountProperties            map[string]interface{} `json:"platform_account_properties,omitempty" mapstructure:"platform_account_properties,omitempty" desc:"Different properties related to the platform the account is related to" flag:"platform-account-properties"`
}
