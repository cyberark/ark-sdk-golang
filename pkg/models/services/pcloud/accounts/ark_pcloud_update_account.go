package accounts

// ArkPCloudUpdateAccount represents the details required to update an account.
type ArkPCloudUpdateAccount struct {
	ArkPCloudAccountSecretManagement
	ArkPCloudAccountRemoteMachinesAccess
	AccountID                 string                 `json:"account_id" mapstructure:"account_id" desc:"The account id to update" flag:"account-id" validate:"required"`
	Name                      string                 `json:"name,omitempty" mapstructure:"name,omitempty" desc:"Name of the account to update" flag:"name"`
	Address                   string                 `json:"address,omitempty" mapstructure:"address,omitempty" desc:"Address of the account to update" flag:"address"`
	Username                  string                 `json:"username,omitempty" mapstructure:"username,omitempty" desc:"Username of the account to update" flag:"username"`
	PlatformID                string                 `json:"platform_id,omitempty" mapstructure:"platform_id,omitempty" desc:"Platform id to relate the account to to update" flag:"platform-id"`
	SafeName                  string                 `json:"safe_name,omitempty" mapstructure:"safe_name,omitempty" desc:"Safe name to store the account in to update" flag:"safe-name"`
	SecretType                string                 `json:"secret_type,omitempty" mapstructure:"secret_type,omitempty" desc:"Type of the secret of the account to update" flag:"secret-type"`
	PlatformAccountProperties map[string]interface{} `json:"platform_account_properties,omitempty" mapstructure:"platform_account_properties,omitempty" desc:"Different properties related to the platform the account is related to to update" flag:"platform-account-properties"`
}
