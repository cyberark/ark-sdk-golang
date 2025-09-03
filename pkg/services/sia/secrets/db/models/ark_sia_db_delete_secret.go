package models

// ArkSIADBDeleteSecret is the struct for deleting a secret from the Ark SIA DB.
type ArkSIADBDeleteSecret struct {
	SecretID   string `json:"secret_id,omitempty" mapstructure:"secret_id" flag:"secret-id" desc:"ID of the secret to delete"`
	SecretName string `json:"secret_name,omitempty" mapstructure:"secret_name" flag:"secret-name" desc:"Name of the secret to delete"`
}
