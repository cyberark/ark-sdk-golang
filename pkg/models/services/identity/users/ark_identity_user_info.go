package users

// ArkIdentityUserInfo represents the schema for user information.
type ArkIdentityUserInfo struct {
	FirstName  string                   `json:"first_name" mapstructure:"first_name" flag:"first-name" desc:"First name of the user"`
	LastName   string                   `json:"last_name" mapstructure:"last_name" flag:"last-name" desc:"Last name of the user"`
	HomeNumber string                   `json:"home_number" mapstructure:"home_number" flag:"home-number" desc:"Home number of the user"`
	Manager    string                   `json:"manager" mapstructure:"manager" flag:"manager" desc:"Manager of the user"`
	Username   string                   `json:"username" mapstructure:"username" flag:"username" desc:"Username info"`
	Groups     []map[string]interface{} `json:"groups" mapstructure:"groups" flag:"groups" desc:"AD groups of the user"`
	Rights     []string                 `json:"rights" mapstructure:"rights" flag:"rights" desc:"Administrative rights of the user"`
	Roles      []map[string]interface{} `json:"roles" mapstructure:"roles" flag:"roles" desc:"Roles of the user"`
}
