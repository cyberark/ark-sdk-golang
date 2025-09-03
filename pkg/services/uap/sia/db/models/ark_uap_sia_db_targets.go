package models

// ArkUAPSIADBTargets represents a collection of database instance targets in the UAP SIA DB.
type ArkUAPSIADBTargets struct {
	Instances []ArkUAPSIADBInstanceTarget `json:"instances" mapstructure:"instances" flag:"instances" desc:"List of database instance targets" validate:"min=1,max=1000"`
}
