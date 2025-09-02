package models

// ArkIdentityRemoveRoleFromRole represents the schema for removing a role from another role.
type ArkIdentityRemoveRoleFromRole struct {
	RoleNameToRemove string `json:"role_name_to_remove" mapstructure:"role_name_to_remove" flag:"role-name-to-remove" desc:"Role name to remove from the role" required:"true"`
	RoleName         string `json:"role_name" mapstructure:"role_name" flag:"role-name" desc:"Name of the role to remove the role from" required:"true"`
}
