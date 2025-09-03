package models

// ArkPCloudReconcileAccountCredentials represents the details required to mark an account for reconciliation.
type ArkPCloudReconcileAccountCredentials struct {
	AccountID string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to mark for reconciliation" flag:"account-id" validate:"required"`
}
