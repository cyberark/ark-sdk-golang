package models

// ArkIdentityAddAdminRightsToRole represents the schema for adding admin rights to a role.
type ArkIdentityAddAdminRightsToRole struct {
	RoleID      string   `json:"role_id,omitempty" mapstructure:"role_id" flag:"role-id" desc:"Role id to add admin rights to"`
	RoleName    string   `json:"role_name,omitempty" mapstructure:"role_name" flag:"role-name" desc:"Role name to add admin rights to"`
	AdminRights []string `json:"admin_rights" mapstructure:"admin_rights" flag:"admin-rights" desc:"Admin rights to add to the role" required:"true"`
}
