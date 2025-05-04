package roles

// ArkIdentityCreateRole represents the schema for creating a role.
type ArkIdentityCreateRole struct {
	RoleName    string   `json:"role_name" mapstructure:"role_name" flag:"role-name" desc:"Role name to create" required:"true"`
	Description string   `json:"description,omitempty" mapstructure:"description" flag:"description" desc:"Description of the role"`
	AdminRights []string `json:"admin_rights" mapstructure:"admin_rights" flag:"admin-rights" desc:"Admin rights to add to the role"`
}
