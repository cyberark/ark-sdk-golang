package syncpolicies

// ArkSecHubSyncPoliciesFilters contains the policy id for the policy to retrieve
type ArkSecHubSyncPoliciesFilters struct {
	Filters    string `json:"filters,omitempty" mapstructure:"filters,omitempty" desc:"Sync Policy filters. Example: --Filter 'target.id EQ store-cfd25162-f8a9-4d94-8d36-f46c4b60d65'" flag:"filters"`
	Projection string `json:"projection,omitempty" mapstructure:"projection,omitempty" description:"Data representation method (EXTEND,REGULAR)" default:"REGULAR" flag:"projection" choices:"EXTEND,REGULAR"`
}
