package roles

// ArkIdentityRole represents the schema for a role.
type ArkIdentityRole struct {
	RoleID   string `json:"role_id" mapstructure:"role_id" flag:"role-id" desc:"Identifier of the role" required:"true"`
	RoleName string `json:"role_name" mapstructure:"role_name" flag:"role-name" desc:"Name of the role" required:"true"`
}
