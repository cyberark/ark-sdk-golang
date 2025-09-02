package models

// ArkIdentityAddRoleToRole represents the schema for adding a role to another role.
type ArkIdentityAddRoleToRole struct {
	RoleNameToAdd string `json:"role_name_to_add" mapstructure:"role_name_to_add" flag:"role-name-to-add" desc:"Role name to add to the role" required:"true"`
	RoleName      string `json:"role_name" mapstructure:"role_name" flag:"role-name" desc:"Name of the role to add the role to" required:"true"`
}
