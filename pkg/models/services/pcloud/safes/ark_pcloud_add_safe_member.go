package safes

// ArkPCloudAddSafeMember represents the details required to add a safe member.
type ArkPCloudAddSafeMember struct {
	SafeID                   string                          `json:"safe_id" mapstructure:"safe_id" desc:"Safe url id to add the member to" flag:"safe-id" validate:"required"`
	MemberName               string                          `json:"member_name" mapstructure:"member_name" desc:"Name of the member to add" flag:"member-name" validate:"required"`
	MemberType               string                          `json:"member_type" mapstructure:"member_type" desc:"Type of the member (User,Group,Role)" flag:"member-type" validate:"required" choices:"User,Group,Role"`
	SearchIn                 string                          `json:"search_in,omitempty" mapstructure:"search_in,omitempty" desc:"Where to search the member in, vault or a domain" flag:"search-in"`
	MembershipExpirationDate int                             `json:"membership_expiration_date,omitempty" mapstructure:"membership_expiration_date,omitempty" desc:"What is the member expiration date" flag:"membership-expiration-date"`
	Permissions              *ArkPCloudSafeMemberPermissions `json:"permissions,omitempty" mapstructure:"permissions,omitempty" desc:"Permissions of the safe member on the safe"`
	PermissionSet            string                          `json:"permission_set" mapstructure:"permission_set,omitempty" desc:"Predefined permission set to use (connect_only,read_only,approver,accounts_manager,full,custom)" flag:"permission-set" default:"read_only" choices:"connect_only,read_only,approver,accounts_manager,full,custom"`
}
