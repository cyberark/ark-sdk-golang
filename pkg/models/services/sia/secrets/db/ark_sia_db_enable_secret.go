package db

// ArkSIADBEnableSecret is the struct for enabling a secret in the Ark SIA DB.
type ArkSIADBEnableSecret struct {
	SecretID   string `json:"secret_id,omitempty" mapstructure:"secret_id" flag:"secret-id" desc:"ID of the secret to enable"`
	SecretName string `json:"secret_name,omitempty" mapstructure:"secret_name" flag:"secret-name" desc:"Name of the secret to enable"`
}
