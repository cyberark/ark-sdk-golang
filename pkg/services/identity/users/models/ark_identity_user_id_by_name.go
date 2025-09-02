package models

// ArkIdentityUserIDByName represents the schema for finding a user ID by their username.
type ArkIdentityUserIDByName struct {
	Username string `json:"username" mapstructure:"username" flag:"username" desc:"User name to find the id for" required:"true"`
}
