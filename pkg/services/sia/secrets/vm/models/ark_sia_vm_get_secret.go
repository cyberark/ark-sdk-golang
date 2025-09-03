package models

// ArkSIAVMGetSecret represents the request to get a secret in a VM.
type ArkSIAVMGetSecret struct {
	SecretID string `json:"secret_id" mapstructure:"secret_id" flag:"secret-id" desc:"The secret id to get" validate:"required"`
}
