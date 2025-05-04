package roles

// ArkIdentityRemoveUserFromRole represents the schema for removing a user from a role.
type ArkIdentityRemoveUserFromRole struct {
	Username string `json:"username" mapstructure:"username" flag:"username" desc:"Username to remove from the role" required:"true"`
	RoleName string `json:"role_name" mapstructure:"role_name" flag:"role-name" desc:"Name of the role to remove the user from" required:"true"`
}
