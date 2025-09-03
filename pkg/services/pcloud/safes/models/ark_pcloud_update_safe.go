package models

// ArkPCloudUpdateSafe represents the request to update a safe in the PCloud vault.
type ArkPCloudUpdateSafe struct {
	SafeID                    string `json:"safe_id,omitempty" desc:"ID of the safe" flag:"safe-id" validate:"required"`
	SafeName                  string `json:"safe_name,omitempty" mapstructure:"safe_name,omitempty" desc:"Name of the safe" flag:"safe-name"`
	Description               string `json:"description,omitempty" mapstructure:"description,omitempty" desc:"Description about the safe" flag:"description"`
	Location                  string `json:"location,omitempty" mapstructure:"location,omitempty" desc:"Location of the safe in the vault" flag:"location" default:""`
	NumberOfDaysRetention     int    `json:"number_of_days_retention,omitempty" mapstructure:"number_of_days_retention,omitempty" desc:"Number of retention days on the safe objects" flag:"number-of-days-retention" default:"0"`
	NumberOfVersionsRetention int    `json:"number_of_versions_retention,omitempty" mapstructure:"number_of_versions_retention,omitempty" desc:"Number of retention versions on the safe objects" flag:"number-of-versions-retention"`
	AutoPurgeEnabled          bool   `json:"auto_purge_enabled,omitempty" mapstructure:"auto_purge_enabled,omitempty" desc:"Whether auto purge is enabled on the safe" flag:"auto-purge-enabled" default:"false"`
	OlacEnabled               bool   `json:"olac_enabled,omitempty" mapstructure:"olac_enabled,omitempty" desc:"Whether object level access control is enabled" flag:"olac-enabled" default:"false"`
	ManagingCPM               string `json:"managing_cpm,omitempty" mapstructure:"managing_cpm,omitempty" desc:"Managing CPM of the safe" flag:"managing-cpm"`
}
