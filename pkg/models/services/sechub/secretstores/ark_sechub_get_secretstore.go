package secretstores

// ArkSecHubGetSecretStore represents the details required to get a Secret Stores details.
type ArkSecHubGetSecretStore struct {
	SecretStoreID string `json:"secret_store_id" mapstructure:"secret_store_id" desc:"Secret Store id to get details for" flag:"secret-store-id" validate:"required"`
}
