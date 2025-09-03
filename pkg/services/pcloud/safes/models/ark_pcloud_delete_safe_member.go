package models

// ArkPCloudDeleteSafeMember represents the details required to remove a safe member.
type ArkPCloudDeleteSafeMember struct {
	SafeID     string `json:"safe_id" mapstructure:"safe_id" desc:"Safe url id to remove the member from" flag:"safe-id" validate:"required"`
	MemberName string `json:"member_name" mapstructure:"member_name" desc:"Name of the member to remove" flag:"member-name" validate:"required"`
}
