package common

// ArkUAPGetPolicyByNameRequest represents the request to get a policy by its name.
type ArkUAPGetPolicyByNameRequest struct {
	PolicyName string `json:"policy_name" mapstructure:"policy_name" flag:"policy-name" desc:"Policy name to be retrieved"`
}
