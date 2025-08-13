package common

// Possible values for PolicyType
const (
	PolicyTypeRecurring = "Recurring"
	PolicyTypeOnDemand  = "OnDemand"
)

// ArkUAPPolicyEntitlement represents the entitlement details of a policy.
type ArkUAPPolicyEntitlement struct {
	TargetCategory string `json:"target_category" mapstructure:"target_category" flag:"target-category" desc:"The target category of the policy" choices:"Cloud console,VM,DB"`
	LocationType   string `json:"location_type" mapstructure:"location_type" flag:"location-type" desc:"The location type of the policy" choices:"AWS,Azure,GCP,FQDN/IP"`
	PolicyType     string `json:"policy_type" mapstructure:"policy_type" flag:"policy-type" desc:"Policy type" choices:"Recurring,OnDemand" default:"Recurring"`
}
