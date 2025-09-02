package models

// ArkPCloudGetSafe represents the details required to get a safe's details.
type ArkPCloudGetSafe struct {
	SafeID string `json:"safe_id" mapstructure:"safe_id" desc:"Safe id to get details for" flag:"safe-id" validate:"required"`
}
