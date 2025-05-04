package roles

// ArkIdentityDeleteRole represents the schema for deleting a role.
type ArkIdentityDeleteRole struct {
	RoleName string `json:"role_name,omitempty" mapstructure:"role_name" flag:"role-name" desc:"Role name to delete"`
	RoleID   string `json:"role_id,omitempty" mapstructure:"role_id" flag:"role-id" desc:"Role id to delete"`
}
