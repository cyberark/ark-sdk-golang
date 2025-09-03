package models

// ArkSIADBGetSecret is the struct for retrieving a secret from the Ark SIA DB.
type ArkSIADBGetSecret struct {
	SecretID   string `json:"secret_id,omitempty" mapstructure:"secret_id" flag:"secret-id" desc:"ID of the secret to get"`
	SecretName string `json:"secret_name,omitempty" mapstructure:"secret_name" flag:"secret-name" desc:"Name of the secret to get"`
}
