package vm

// ArkSIAVMChangeSecret represents the request to change a secret in a VM.
type ArkSIAVMChangeSecret struct {
	SecretID            string                 `json:"secret_id" mapstructure:"secret_id" flag:"secret-id" desc:"The secret id to change" validate:"required"`
	SecretName          string                 `json:"secret_name,omitempty" mapstructure:"secret_name,omitempty" flag:"secret-name" desc:"The new name of the secret"`
	SecretDetails       map[string]interface{} `json:"secret_details,omitempty" mapstructure:"secret_details,omitempty" flag:"secret-details" desc:"New secret details to add / change"`
	IsDisabled          bool                   `json:"is_disabled,omitempty" mapstructure:"is_disabled,omitempty" flag:"is-disabled" desc:"Whether to disable the secret" default:"false"`
	ProvisionerUsername string                 `json:"provisioner_username,omitempty" mapstructure:"provisioner_username,omitempty" flag:"provisioner-username" desc:"If provisioner user type secret, the new username"`
	ProvisionerPassword string                 `json:"provisioner_password,omitempty" mapstructure:"provisioner_password,omitempty" flag:"provisioner-password" desc:"If provisioner user type secret, the new password"`
	PCloudAccountSafe   string                 `json:"pcloud_account_safe,omitempty" mapstructure:"pcloud_account_safe,omitempty" flag:"pcloud-account-safe" desc:"If pcloud account type secret, the new account safe"`
	PCloudAccountName   string                 `json:"pcloud_account_name,omitempty" mapstructure:"pcloud_account_name,omitempty" flag:"pcloud-account-name" desc:"If pcloud account type secret, the new account name"`
}
