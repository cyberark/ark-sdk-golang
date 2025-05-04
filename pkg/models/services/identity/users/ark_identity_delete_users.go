package users

// ArkIdentityDeleteUsers represents the schema for deleting multiple users.
type ArkIdentityDeleteUsers struct {
	UserIDs []string `json:"user_ids" mapstructure:"user_ids" flag:"user-ids" desc:"User IDs to delete" required:"true"`
}
