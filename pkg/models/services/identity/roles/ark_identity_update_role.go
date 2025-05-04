package roles

// ArkIdentityUpdateRole represents the schema for updating a role.
type ArkIdentityUpdateRole struct {
	RoleName    string `json:"role_name,omitempty" mapstructure:"role_name" flag:"role-name" desc:"Role name to update"`
	RoleID      string `json:"role_id,omitempty" mapstructure:"role_id" flag:"role-id" desc:"Role id to update"`
	NewRoleName string `json:"new_role_name,omitempty" mapstructure:"new_role_name" flag:"new-role-name" desc:"New role name to update to"`
	Description string `json:"description,omitempty" mapstructure:"description" flag:"description" desc:"New description of the role"`
}
