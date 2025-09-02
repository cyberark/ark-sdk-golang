package models

// ArkIdentityUpdateUser represents the schema for updating a user's details.
type ArkIdentityUpdateUser struct {
	UserID       string `json:"user_id,omitempty" mapstructure:"user_id" flag:"user-id" desc:"Users id that we change the details for"`
	Username     string `json:"username,omitempty" mapstructure:"username" flag:"username" desc:"Username that we change the details for"`
	NewUsername  string `json:"new_username,omitempty" mapstructure:"new_username" flag:"new-username" desc:"Name of the user to change"`
	DisplayName  string `json:"display_name,omitempty" mapstructure:"display_name" flag:"display-name" desc:"Display name of the user to change"`
	Email        string `json:"email,omitempty" mapstructure:"email" flag:"email" desc:"Email of the user to change"`
	MobileNumber string `json:"mobile_number,omitempty" mapstructure:"mobile_number" flag:"mobile-number" desc:"Mobile number of the user to change"`
}
