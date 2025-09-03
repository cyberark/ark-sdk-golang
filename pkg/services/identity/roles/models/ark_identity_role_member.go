package models

// ArkIdentityRoleMember represents the schema for a role member.
type ArkIdentityRoleMember struct {
	MemberID   string `json:"member_id" mapstructure:"member_id" flag:"member-id" desc:"ID of the member" required:"true"`
	MemberName string `json:"member_name" mapstructure:"member_name" flag:"member-name" desc:"Name of the member" required:"true"`
	MemberType string `json:"member_type" mapstructure:"member_type" flag:"member-type" desc:"Type of the member" required:"true"`
}
