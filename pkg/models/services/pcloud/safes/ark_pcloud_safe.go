package safes

// ArkPCloudSafeCreator represents the creator of a safe.
type ArkPCloudSafeCreator struct {
	ID   string `json:"id" mapstructure:"id" desc:"ID of the safe creator" flag:"id"`
	Name string `json:"name" mapstructure:"name" desc:"Name of the safe creator" flag:"name"`
}

// ArkPCloudSafe represents a safe with additional details.
type ArkPCloudSafe struct {
	SafeName                  string               `json:"safe_name,omitempty" mapstructure:"safe_name" desc:"Name of the safe" flag:"safe-name"`
	Description               string               `json:"description,omitempty" mapstructure:"description" desc:"Description about the safe" flag:"description"`
	Location                  string               `json:"location,omitempty" mapstructure:"location" desc:"Location of the safe in the vault" flag:"location" default:"\\"`
	NumberOfDaysRetention     int                  `json:"number_of_days_retention,omitempty" mapstructure:"number_of_days_retention" desc:"Number of retention days on the safe objects" flag:"number-of-days-retention" default:"0"`
	NumberOfVersionsRetention int                  `json:"number_of_versions_retention,omitempty" mapstructure:"number_of_versions_retention" desc:"Number of retention versions on the safe objects" flag:"number-of-versions-retention"`
	AutoPurgeEnabled          bool                 `json:"auto_purge_enabled,omitempty" mapstructure:"auto_purge_enabled" desc:"Whether auto purge is enabled on the safe" flag:"auto-purge-enabled" default:"false"`
	OlacEnabled               bool                 `json:"olac_enabled,omitempty" mapstructure:"olac_enabled" desc:"Whether object level access control is enabled" flag:"olac-enabled" default:"false"`
	ManagingCPM               string               `json:"managing_cpm,omitempty" mapstructure:"managing_cpm" desc:"Managing CPM of the safe" flag:"managing-cpm"`
	Creator                   ArkPCloudSafeCreator `json:"creator,omitempty" mapstructure:"creator" desc:"Creator of the safe" flag:"creator"`
	CreationTime              int                  `json:"creation_time,omitempty" mapstructure:"creation_time" desc:"Creation time of the safe" flag:"creation-time"`
	LastModificationTime      int                  `json:"last_modification_time,omitempty" mapstructure:"last_modification_time" desc:"Last time the safe was modified" flag:"last-modification-time"`
	SafeID                    string               `json:"safe_id" mapstructure:"safe_id" desc:"Safe url to access with as an id" flag:"safe-id"`
	SafeNumber                int                  `json:"safe_number,omitempty" mapstructure:"safe_number" desc:"ID number of the safe" flag:"safe-number"`
	IsExpiredMember           bool                 `json:"is_expired_member,omitempty" mapstructure:"is_expired_member" desc:"Whether any member is expired" flag:"is-expired-member" default:"false"`
}
