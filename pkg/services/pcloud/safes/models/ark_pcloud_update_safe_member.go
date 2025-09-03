package models

// ArkPCloudUpdateSafeMember represents the details required to update a safe member.
type ArkPCloudUpdateSafeMember struct {
	SafeID                   string                          `json:"safe_id" mapstructure:"safe_id" desc:"Safe url id to update the member on" flag:"safe-id" validate:"required"`
	MemberName               string                          `json:"member_name" mapstructure:"member_name" desc:"Name of the member to update" flag:"member-name" validate:"required"`
	MembershipExpirationDate int                             `json:"membership_expiration_date,omitempty" mapstructure:"membership_expiration_date,omitempty" desc:"What is the member expiration date to update" flag:"membership-expiration-date"`
	Permissions              *ArkPCloudSafeMemberPermissions `json:"permissions,omitempty" mapstructure:"permissions,omitempty" desc:"Permissions of the safe member on the safe to update" flag:"permissions"`
	PermissionSet            string                          `json:"permission_set,omitempty" mapstructure:"permission_set,omitempty" desc:"Predefined permission set to update to (connect_only,read_only,approver,accounts_manager,full,custom)" flag:"permission-set" choices:"connect_only,read_only,approver,accounts_manager,full,custom"`
}
