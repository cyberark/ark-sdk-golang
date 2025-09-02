package models

// ArkSecHubDeleteSecretStore contains the secret store ID to delete
type ArkSecHubDeleteSecretStore struct {
	SecretStoreID string `json:"secret_store_id" mapstructure:"secret_store_id" desc:"Secret store id to delete" flag:"secret-store-id" validate:"required"`
}
