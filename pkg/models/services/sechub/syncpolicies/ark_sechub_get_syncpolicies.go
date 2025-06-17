package syncpolicies

// ArkSecHubGetSyncPolicies contains the query option for retrieving sync policies.
type ArkSecHubGetSyncPolicies struct {
	Projection string `json:"projection,omitempty" mapstructure:"projection,omitempty" description:"Data representation method. Allowed values: 'EXTEND, REGULAR'" default:"REGULAR"`
}
