package models

// ArkSIABulkAddTargetSetsItem represents the request to add multiple target sets to a strong account in a workspace.
type ArkSIABulkAddTargetSetsItem struct {
	StrongAccountID string               `json:"strong_account_id" mapstructure:"strong_account_id" flag:"strong-account-id" desc:"Secret ID of the strong account related to this set" validate:"required"`
	TargetSets      []ArkSIAAddTargetSet `json:"target_sets" mapstructure:"target_sets" flag:"target-sets" desc:"The target sets to associate with the strong account" validate:"required,dive"`
}

// ArkSIABulkAddTargetSets represents the request to add multiple target sets to a strong account in a workspace.
type ArkSIABulkAddTargetSets struct {
	TargetSetsMapping []ArkSIABulkAddTargetSetsItem `json:"target_sets_mapping" mapstructure:"target_sets_mapping" flag:"target-sets-mapping" desc:"Bulk of target set mappings to add" validate:"required,dive"`
}
