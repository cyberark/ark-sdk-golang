package models

// ArkIdentityDeleteUser represents the schema for deleting a user.
type ArkIdentityDeleteUser struct {
	UserID   string `json:"user_id,omitempty" mapstructure:"user_id" flag:"user-id" desc:"User ID to delete"`
	Username string `json:"username,omitempty" mapstructure:"username" flag:"username" desc:"Username to delete"`
}
