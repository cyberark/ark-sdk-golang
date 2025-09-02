package models

// ArkSIATargetSetsStats represents the statistics of target sets in a workspace.
type ArkSIATargetSetsStats struct {
	TargetSetsCount              int            `json:"target_sets_count" mapstructure:"target_sets_count" flag:"target-sets-count" desc:"Total target sets count" validate:"required"`
	TargetSetsCountPerSecretType map[string]int `json:"target_sets_count_per_secret_type" mapstructure:"target_sets_count_per_secret_type" flag:"target-sets-count-per-secret-type" desc:"Target sets count per secret type" validate:"required"`
}
