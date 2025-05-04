package roles

// ArkIdentityRoleIDByName represents the schema for finding the ID of a role by its name.
type ArkIdentityRoleIDByName struct {
	RoleName string `json:"role_name" mapstructure:"role_name" flag:"role-name" desc:"Role name to find the id for" required:"true"`
}
