package vm

// ArkSIAVMAddSecret represents the request to add a secret in a VM.
type ArkSIAVMAddSecret struct {
	SecretName          string                 `json:"secret_name,omitempty" mapstructure:"secret_name,omitempty" flag:"secret-name" desc:"Optional name of the secret"`
	SecretDetails       map[string]interface{} `json:"secret_details,omitempty" mapstructure:"secret_details,omitempty" flag:"secret-details" desc:"Optional extra details about the secret"`
	SecretType          string                 `json:"secret_type" mapstructure:"secret_type" flag:"secret-type" desc:"Type of the secret to add, data is picked according to the chosen type (ProvisionerUser,PCloudAccount)" validate:"required" choices:"ProvisionerUser,PCloudAccount"`
	IsDisabled          bool                   `json:"is_disabled" mapstructure:"is_disabled" flag:"is-disabled" desc:"Whether the secret should be disabled or not" default:"false"`
	ProvisionerUsername string                 `json:"provisioner_username,omitempty" mapstructure:"provisioner_username,omitempty" flag:"provisioner-username" desc:"If provisioner user type is picked, the username"`
	ProvisionerPassword string                 `json:"provisioner_password,omitempty" mapstructure:"provisioner_password,omitempty" flag:"provisioner-password" desc:"If provisioner user type is picked, the password"`
	PCloudAccountSafe   string                 `json:"pcloud_account_safe,omitempty" mapstructure:"pcloud_account_safe,omitempty" flag:"pcloud-account-safe" desc:"If pcloud account type is picked, the account safe"`
	PCloudAccountName   string                 `json:"pcloud_account_name,omitempty" mapstructure:"pcloud_account_name,omitempty" flag:"pcloud-account-name" desc:"If pcloud account type is picked, the account name"`
}
