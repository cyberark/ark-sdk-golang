package common

// ArkUAPResponse represents the response containing a policy ID.
type ArkUAPResponse struct {
	PolicyID string `json:"policy_id" mapstructure:"policy_id" flag:"policy-id" desc:"Policy id"`
}
