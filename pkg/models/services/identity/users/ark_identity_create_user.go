package users

// DefaultAdminRoles defines the default roles for an admin user.
var DefaultAdminRoles = []string{"DpaAdmin", "global auditor", "System Administrator"}

// ArkIdentityCreateUser represents the schema for creating a user.
// Note that username,display_name,email,mobile_number, and password will be auto generated if not given
type ArkIdentityCreateUser struct {
	Username     string   `json:"username" mapstructure:"username" flag:"username" desc:"Name of the user to create"`
	DisplayName  string   `json:"display_name,omitempty" mapstructure:"display_name" flag:"display-name" desc:"Display name of the user"`
	Email        string   `json:"email,omitempty" mapstructure:"email" flag:"email" desc:"Email of the user"`
	MobileNumber string   `json:"mobile_number,omitempty" mapstructure:"mobile_number" flag:"mobile-number" desc:"Mobile number of the user"`
	Suffix       string   `json:"suffix,omitempty" mapstructure:"suffix" flag:"suffix" desc:"Suffix to use for the username"`
	Password     string   `json:"password" mapstructure:"password" flag:"password" desc:"Password of the user"`
	Roles        []string `json:"roles" mapstructure:"roles" flag:"roles" desc:"Roles to add the user to, defaulted to DpaAdmin,global auditor,System Administrator"`
}
