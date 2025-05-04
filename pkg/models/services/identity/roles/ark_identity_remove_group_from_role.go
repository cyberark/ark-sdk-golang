package roles

// ArkIdentityRemoveGroupFromRole represents the schema for removing a group from a role.
type ArkIdentityRemoveGroupFromRole struct {
	GroupName string `json:"group_name" mapstructure:"group_name" flag:"group-name" desc:"Group name to remove from the role" required:"true"`
	RoleName  string `json:"role_name" mapstructure:"role_name" flag:"role-name" desc:"Name of the role to remove the group from" required:"true"`
}
