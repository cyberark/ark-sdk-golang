package targetsets

// ArkSIABulkDeleteTargetSets represents the request to delete multiple target sets in a workspace.
type ArkSIABulkDeleteTargetSets struct {
	TargetSets []string `json:"target_sets" mapstructure:"target_sets" flag:"target-sets" desc:"List of target sets names to delete" validate:"required,dive,required"`
}
