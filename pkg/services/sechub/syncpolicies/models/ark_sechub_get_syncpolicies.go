package models

// ArkSecHubGetSyncPolicies contains the query option for retrieving sync policies.
type ArkSecHubGetSyncPolicies struct {
	Projection string `json:"projection,omitempty" mapstructure:"projection,omitempty" desc:"Data representation method (EXTEND, REGULAR)" flag:"projection" default:"REGULAR" choices:"EXTEND,REGULAR"`
}
