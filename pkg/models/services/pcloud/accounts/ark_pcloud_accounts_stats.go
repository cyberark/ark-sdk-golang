package accounts

// ArkPCloudAccountsStats represents the statistics of accounts.
type ArkPCloudAccountsStats struct {
	AccountsCount             int            `json:"accounts_count" mapstructure:"accounts_count" desc:"Overall accounts count" flag:"accounts-count"`
	AccountsCountByPlatformID map[string]int `json:"accounts_count_by_platform_id" mapstructure:"accounts_count_by_platform_id" desc:"Accounts count by platform id" flag:"accounts-count-by-platform-id"`
	AccountsCountBySafeName   map[string]int `json:"accounts_count_by_safe_name" mapstructure:"accounts_count_by_safe_name" desc:"Accounts count by safe name" flag:"accounts-count-by-safe-name"`
}
