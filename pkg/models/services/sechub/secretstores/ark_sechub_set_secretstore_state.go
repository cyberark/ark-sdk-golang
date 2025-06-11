package secretstores

type ArkSecHubSetSecretStoreState struct {
	SecretStoreID string `json:"secret_store_id" mapstructure:"secret_store_id" desc:"Secret Store id to get details for" flag:"secret-store-id" validate:"required,oneof=enable disable"`
	Action        string `json:"action" mapstructure:"action" flag:"action" desc:"State to set secret store to, Allowed Values: 'enable' and 'Disable'"`
}
