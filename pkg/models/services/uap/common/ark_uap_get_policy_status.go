package common

// ArkUAPGetPolicyStatus represents the request to get the status of a policy by its ID or name.
type ArkUAPGetPolicyStatus struct {
	PolicyID   string `json:"policy_id,omitempty" mapstructure:"policy_id,omitempty" flag:"policy-id" desc:"Policy id to get the status for"`
	PolicyName string `json:"policy_name,omitempty" mapstructure:"policy_name,omitempty" flag:"policy-name" desc:"Policy name to get the status for"`
}
