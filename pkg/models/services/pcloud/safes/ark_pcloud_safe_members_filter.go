package safes

// ArkPCloudSafeMembersFilters represents the details required to filter the members of a safe.
type ArkPCloudSafeMembersFilters struct {
	SafeID     string `json:"safe_id" mapstructure:"safe_id" desc:"Which safe id to filter the members on" flag:"safe-id" validate:"required"`
	Search     string `json:"search,omitempty" mapstructure:"search" desc:"Search by string" flag:"search"`
	Sort       string `json:"sort,omitempty" mapstructure:"sort" desc:"Sort results by given key" flag:"sort"`
	Offset     int    `json:"offset,omitempty" mapstructure:"offset" desc:"Offset to the safe members list" flag:"offset"`
	Limit      int    `json:"limit,omitempty" mapstructure:"limit" desc:"Limit of results" flag:"limit"`
	MemberType string `json:"member_type,omitempty" mapstructure:"member_type" desc:"Filter by type of safe member" flag:"member-type"`
}
