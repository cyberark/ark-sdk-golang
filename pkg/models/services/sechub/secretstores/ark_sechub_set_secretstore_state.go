package secretstores

type ArkSecHubSetSecretStoreState struct {
	SecretStoreID string `json:"secret_store_id" mapstructure:"secret_store_id" desc:"Secret Store id to get details for" flag:"secret-store-id" validate:"required"`
	Action        string `json:"action" mapstructure:"action" flag:"action" desc:"State to set secret store to (enable,disable)" default:"enable" choices:"enable,disable"`
}
