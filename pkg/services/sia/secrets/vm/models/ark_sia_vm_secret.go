package models

// Possible Secret Types
const (
	ProvisionerUser = "ProvisionerUser"
	PCloudAccount   = "PCloudAccount"
)

// ArkSIAVMDataMessage represents a data message in the Ark SIA VM.
type ArkSIAVMDataMessage struct {
	MessageID string `json:"message_id" mapstructure:"message_id" flag:"message-id" desc:"Data Message ID"`
	Data      string `json:"data" mapstructure:"data" flag:"data" desc:"Actual data"`
}

// ArkSIAVMSecretData represents the secret data in the Ark SIA VM.
type ArkSIAVMSecretData struct {
	SecretData      interface{} `json:"secret_data" mapstructure:"secret_data" flag:"secret-data" desc:"Actual secret data, can be of different types, and is base64 encoded if of SecretBytes, Otherwise Stored in the jit data message as a string Or as a dict of secret data to be encrypted"`
	TenantEncrypted bool        `json:"tenant_encrypted" mapstructure:"tenant_encrypted" flag:"tenant-encrypted" desc:"Whether this secret is encrypted by the tenant key or not"`
}

// ArkSIAVMSecret represents a secret in the Ark SIA VM.
type ArkSIAVMSecret struct {
	SecretID      string                 `json:"secret_id" mapstructure:"secret_id" flag:"secret-id" desc:"ID of the secret"`
	TenantID      string                 `json:"tenant_id,omitempty" mapstructure:"tenant_id,omitempty" flag:"tenant-id" desc:"Tenant ID of the secret"`
	Secret        ArkSIAVMSecretData     `json:"secret,omitempty" mapstructure:"secret,omitempty" flag:"secret" desc:"Secret itself"`
	SecretType    string                 `json:"secret_type" mapstructure:"secret_type" flag:"secret-type" desc:"Type of the secret" choices:"ProvisionerUser,PCloudAccount"`
	SecretDetails map[string]interface{} `json:"secret_details" mapstructure:"secret_details" flag:"secret-details" desc:"Secret extra details"`
	IsActive      bool                   `json:"is_active" mapstructure:"is_active" flag:"is-active" desc:"Whether this secret is active or not and can be retrieved or modified"`
	IsRotatable   bool                   `json:"is_rotatable" mapstructure:"is_rotatable" flag:"is-rotatable" desc:"Whether this secret can be rotated"`
	CreationTime  string                 `json:"creation_time" mapstructure:"creation_time" flag:"creation-time" desc:"Creation time of the secret"`
	LastModified  string                 `json:"last_modified" mapstructure:"last_modified" flag:"last-modified" desc:"Last time the secret was modified"`
	SecretName    string                 `json:"secret_name,omitempty" mapstructure:"secret_name,omitempty" flag:"secret-name" desc:"A friendly name label"`
}
