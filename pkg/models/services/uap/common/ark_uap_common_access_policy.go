package common

// Possible values for DelegationClassification
const (
	DelegationClassificationUnrestricted = "Unrestricted"
	DelegationClassificationRestricted   = "Restricted"
)

// ArkUAPCommonAccessPolicy represents the access policy in UAP.
type ArkUAPCommonAccessPolicy struct {
	Metadata                 ArkUAPMetadata    `json:"metadata,omitempty" mapstructure:"metadata,omitempty" flag:"metadata" desc:"Policy metadata id name and extra information"`
	Principals               []ArkUAPPrincipal `json:"principals,omitempty" mapstructure:"principals,omitempty" flag:"principals" desc:"List of users, groups and roles that the policy applies to"`
	DelegationClassification string            `json:"delegation_classification" mapstructure:"delegation_classification" flag:"delegation-classification" desc:"Indicates the user rights for the current policy" choices:"Restricted,Unrestricted" default:"Unrestricted"`
}
