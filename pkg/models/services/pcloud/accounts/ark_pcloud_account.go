package accounts

// Possible Secret Types
const (
	Password = "password"
	Key      = "key"
)

// ArkPCloudAccountSecretManagement represents the secret management properties of an account.
type ArkPCloudAccountSecretManagement struct {
	AutomaticManagementEnabled bool   `json:"automatic_management_enabled,omitempty" mapstructure:"automatic_management_enabled,omitempty" desc:"Whether automatic management of the account is enabled or not" flag:"automatic-management-enabled"`
	ManualManagementReason     string `json:"manual_management_reason,omitempty" mapstructure:"manual_management_reason,omitempty" desc:"The reason for disabling automatic management" flag:"manual-management-reason"`
	LastModifiedTime           int    `json:"last_modified_time,omitempty" mapstructure:"last_modified_time,omitempty" desc:"Last time the management properties were modified" flag:"last-modified-time"`
}

// ArkPCloudAccountRemoteMachinesAccess represents the remote machine access properties of an account.
type ArkPCloudAccountRemoteMachinesAccess struct {
	RemoteMachines                   []string `json:"remote_machines,omitempty" mapstructure:"remote_machines,omitempty" desc:"Remote machines the access of this account is allowed" flag:"remote-machines"`
	AccessRestrictedToRemoteMachines bool     `json:"access_restricted_to_remote_machines,omitempty" mapstructure:"access_restricted_to_remote_machines,omitempty" desc:"Whether the access is only restricted to those remote machines" flag:"access-restricted-to-remote-machines"`
}

// ArkPCloudAccount represents the full properties of an account.
type ArkPCloudAccount struct {
	AccountID                 string                               `json:"account_id" mapstructure:"account_id" desc:"ID of the account" flag:"account-id" validate:"required"`
	Status                    string                               `json:"status,omitempty" mapstructure:"status,omitempty" desc:"Status of the account" flag:"status"`
	CreatedTime               int                                  `json:"created_time,omitempty" mapstructure:"created_time,omitempty" desc:"Creation time of the account" flag:"created-time"`
	CategoryModificationTime  int                                  `json:"category_modification_time,omitempty" mapstructure:"category_modification_time,omitempty" desc:"Category modification time of the account" flag:"category-modification-time"`
	Name                      string                               `json:"name" mapstructure:"name" desc:"Name of the account" flag:"name" validate:"required"`
	SafeName                  string                               `json:"safe_name" mapstructure:"safe_name" desc:"Safe name to store the account in" flag:"safe-name" validate:"required"`
	PlatformID                string                               `json:"platform_id,omitempty" mapstructure:"platform_id,omitempty" desc:"Platform id to relate the account to" flag:"platform-id"`
	UserName                  string                               `json:"user_name,omitempty" mapstructure:"user_name,omitempty" desc:"Username of the account" flag:"user-name"`
	Address                   string                               `json:"address,omitempty" mapstructure:"address,omitempty" desc:"Address of the account" flag:"address"`
	SecretType                string                               `json:"secret_type,omitempty" mapstructure:"secret_type,omitempty" desc:"Type of the secret of the account (password,key)" flag:"secret-type" choices:"password,key"`
	PlatformAccountProperties map[string]interface{}               `json:"platform_account_properties,omitempty" mapstructure:"platform_account_properties,omitempty" desc:"Different properties related to the platform the account is related to" flag:"platform-account-properties"`
	SecretManagement          ArkPCloudAccountSecretManagement     `json:"secret_management,omitempty" mapstructure:"secret_management,omitempty" desc:"Secret mgmt related properties" flag:"secret-management"`
	RemoteMachinesAccess      ArkPCloudAccountRemoteMachinesAccess `json:"remote_machines_access,omitempty" mapstructure:"remote_machines_access,omitempty" desc:"Remote machines access related properties" flag:"remote-machines-access"`
}
