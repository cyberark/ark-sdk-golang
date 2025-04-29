package safes

// ArkPCloudGetSafeMember represents the details required to get a safe member's details.
type ArkPCloudGetSafeMember struct {
	SafeID     string `json:"safe_id" mapstructure:"safe_id" desc:"Safe url id to get the member from" flag:"safe-id" validate:"required"`
	MemberName string `json:"member_name" mapstructure:"member_name" desc:"Name of the member to get" flag:"member-name" validate:"required"`
}
