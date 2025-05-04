package roles

// ArkIdentityAddGroupToRole represents the schema for adding a group to a role.
type ArkIdentityAddGroupToRole struct {
	GroupName string `json:"group_name" mapstructure:"group_name" flag:"group-name" desc:"Group name to add to the role" required:"true"`
	RoleName  string `json:"role_name" mapstructure:"role_name" flag:"role-name" desc:"Name of the role to add the group to" required:"true"`
}
