package models

// ArkSIADeleteTargetSet represents the request to delete a target set in a workspace.
type ArkSIADeleteTargetSet struct {
	ID string `json:"id" mapstructure:"id" flag:"id" desc:"ID of the target set to delete" validate:"required"`
}
