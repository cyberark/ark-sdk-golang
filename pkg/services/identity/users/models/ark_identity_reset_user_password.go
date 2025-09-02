package models

// ArkIdentityResetUserPassword represents the schema for resetting a user's password.
type ArkIdentityResetUserPassword struct {
	Username    string `json:"username" mapstructure:"username" flag:"username" desc:"Username to reset the password for" required:"true"`
	NewPassword string `json:"new_password" mapstructure:"new_password" flag:"new-password" desc:"New password to reset to" required:"true"`
}
