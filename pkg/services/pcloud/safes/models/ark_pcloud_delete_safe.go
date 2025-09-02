package models

// ArkPCloudDeleteSafe represents the details required to delete a safe.
type ArkPCloudDeleteSafe struct {
	SafeID string `json:"safe_id" mapstructure:"safe_id" desc:"Safe id to delete" flag:"safe-id" validate:"required"`
}
