package common

// ArkUAPDeletePolicyRequest represents the request to delete a policy in UAP.
type ArkUAPDeletePolicyRequest struct {
	PolicyID string `json:"policy_id" mapstructure:"policy_id" flag:"policy-id" desc:"Policy id to be deleted"`
}
