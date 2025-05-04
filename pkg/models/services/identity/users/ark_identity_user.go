package users

import "time"

// ArkIdentityUser represents the schema for a user.
type ArkIdentityUser struct {
	UserID       string     `json:"user_id" mapstructure:"user_id" flag:"user-id" desc:"User identifier"`
	Username     string     `json:"username" mapstructure:"username" flag:"username" desc:"Name of the user"`
	DisplayName  string     `json:"display_name,omitempty" mapstructure:"display_name" flag:"display-name" desc:"Display name of the user"`
	Email        string     `json:"email,omitempty" mapstructure:"email" flag:"email" desc:"Email of the user"`
	MobileNumber string     `json:"mobile_number,omitempty" mapstructure:"mobile_number" flag:"mobile-number" desc:"Mobile number of the user"`
	Roles        []string   `json:"roles,omitempty" mapstructure:"roles" flag:"roles" desc:"Roles of the user"`
	LastLogin    *time.Time `json:"last_login,omitempty" mapstructure:"last_login" flag:"last-login" desc:"Last login of the user"`
}
