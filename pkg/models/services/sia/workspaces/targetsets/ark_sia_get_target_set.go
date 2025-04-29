package targetsets

// ArkSIAGetTargetSet represents the request to retrieve a target set in a workspace.
type ArkSIAGetTargetSet struct {
	Name string `json:"name" mapstructure:"name" flag:"name" desc:"Name of the target set to retrieve" validate:"required"`
}
