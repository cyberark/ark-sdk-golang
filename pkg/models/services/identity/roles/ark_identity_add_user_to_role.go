package roles

// ArkIdentityAddUserToRole represents the schema for adding a user to a role.
type ArkIdentityAddUserToRole struct {
	Username string `json:"username" mapstructure:"username" flag:"username" desc:"Username to add to the role" required:"true"`
	RoleName string `json:"role_name" mapstructure:"role_name" flag:"role-name" desc:"Name of the role to add the user to" required:"true"`
}
