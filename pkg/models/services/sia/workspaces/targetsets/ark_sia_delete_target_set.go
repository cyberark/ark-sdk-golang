package targetsets

// ArkSIADeleteTargetSet represents the request to delete a target set in a workspace.
type ArkSIADeleteTargetSet struct {
	Name string `json:"name" mapstructure:"name" flag:"name" desc:"Name of the target set to delete" validate:"required"`
}
