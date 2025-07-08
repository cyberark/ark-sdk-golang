package secretstores

// ArkSecHubSetSecretStoresState represents the body of sent when setting multiple secret store states
type ArkSecHubSetSecretStoresState struct {
	Action         string   `json:"action" mapstructure:"action" flag:"action" desc:"State to set secret stores to (enable,disable)" choices:"enable,disable"`
	SecretStoreIDs []string `json:"secret_store_ids" mapstructure:"secret_store_ids" desc:"List of Secret Store ids to set state for" flag:"secret-store-ids" validate:"required"`
}

// ArkSecHubSetSecretStoresStateResults represents the individual object for each secret store for which
// the secret store state was set
type ArkSecHubSetSecretStoresStateResults struct {
	SecretStoreID string `json:"secret_store_id" mapstructure:"secret_store_id"`
	Result        string `json:"result" mapstructure:"result"`
	ErrorMessage  string `json:"error_message" mapstructure:"error_message"`
}

// ArkSecHubSetSecretStoresStateResponse is the outer object which contains the indvidual secret store state
// response objects
type ArkSecHubSetSecretStoresStateResponse struct {
	Results []ArkSecHubSetSecretStoresStateResults `json:"results" mapstructure:"results"`
}
