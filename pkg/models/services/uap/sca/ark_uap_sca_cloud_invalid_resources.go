package sca

// Possible values for ArkUAPSCACloudInvalidWorkspace status.
const (
	StatusRemoved   string = "REMOVED"
	StatusSuspended string = "SUSPENDED"
)

// ArkUAPSCACloudInvalidWorkspace represents an invalid workspace.
type ArkUAPSCACloudInvalidWorkspace struct {
	ID     string `json:"id" mapstructure:"id" flag:"id" desc:"Resource ID"`
	Status string `json:"status" mapstructure:"status" flag:"status" desc:"Workspace status" choices:"REMOVED,SUSPENDED"`
}

// ArkUAPSCACloudInvalidRole represents an invalid role.
type ArkUAPSCACloudInvalidRole struct {
	ID string `json:"id" mapstructure:"id" flag:"id" desc:"Invalid role ID"`
}

// ArkUAPSCACloudInvalidWebapp represents an invalid webapp.
type ArkUAPSCACloudInvalidWebapp struct {
	ID string `json:"id" mapstructure:"id" flag:"id" desc:"Invalid webapp ID"`
}

// ArkUAPSCACloudInvalidResources represents a collection of invalid resources.
type ArkUAPSCACloudInvalidResources struct {
	Workspaces []ArkUAPSCACloudInvalidWorkspace `json:"workspaces,omitempty" mapstructure:"workspaces,omitempty" flag:"workspaces" desc:"List of invalid workspaces"`
	Roles      []ArkUAPSCACloudInvalidRole      `json:"roles,omitempty" mapstructure:"roles,omitempty" flag:"roles" desc:"List of invalid roles"`
	Webapps    []ArkUAPSCACloudInvalidWebapp    `json:"webapps,omitempty" mapstructure:"webapps,omitempty" flag:"webapps" desc:"List of invalid webapps"`
}
