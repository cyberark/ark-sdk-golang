package users

// ArkIdentityUserByName represents the schema for finding a user by their username.
type ArkIdentityUserByName struct {
	Username string `json:"username" mapstructure:"username" flag:"username" desc:"User name to find the id for" required:"true"`
}
