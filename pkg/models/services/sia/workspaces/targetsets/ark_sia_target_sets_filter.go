package targetsets

// ArkSIATargetSetsFilter represents the filter criteria for retrieving target sets in a workspace.
type ArkSIATargetSetsFilter struct {
	Name       string `json:"name,omitempty" mapstructure:"name,omitempty" flag:"name" desc:"Name filter wildcard"`
	SecretType string `json:"secret_type,omitempty" mapstructure:"secret_type,omitempty" flag:"secret-type" desc:"Secret type filter" choices:"ProvisionerUser,PCloudAccount"`
}
