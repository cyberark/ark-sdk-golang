package models

// ArkPCloudAccountsFilter represents the filter options for accounts.
type ArkPCloudAccountsFilter struct {
	Search     string `json:"search,omitempty" mapstructure:"search,omitempty" desc:"Search by string" flag:"search"`
	SearchType string `json:"search_type,omitempty" mapstructure:"search_type,omitempty" desc:"Search type to filter with (contains or startswith)" flag:"search-type"`
	Sort       string `json:"sort,omitempty" mapstructure:"sort,omitempty" desc:"Sort results by given key" flag:"sort"`
	SafeName   string `json:"safe_name,omitempty" mapstructure:"safe_name,omitempty" desc:"Safe name to filter by" flag:"safe-name"`
	Offset     int    `json:"offset,omitempty" mapstructure:"offset,omitempty" desc:"Offset to the accounts list" flag:"offset"`
	Limit      int    `json:"limit,omitempty" mapstructure:"limit,omitempty" desc:"Limit of results" flag:"limit"`
}
