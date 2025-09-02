package models

// ArkIdentityUserByID represents the schema for finding a user by ID.
type ArkIdentityUserByID struct {
	UserID string `json:"user_id" mapstructure:"user_id" flag:"user-id" desc:"Id to find the user for" required:"true"`
}
