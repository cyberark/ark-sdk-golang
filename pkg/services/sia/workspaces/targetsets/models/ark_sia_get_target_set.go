package models

// ArkSIAGetTargetSet represents the request to retrieve a target set in a workspace.
type ArkSIAGetTargetSet struct {
	ID string `json:"id" mapstructure:"id" flag:"id" desc:"ID of the target set to retrieve"`
}
