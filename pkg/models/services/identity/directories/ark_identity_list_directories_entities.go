package directories

// ArkIdentityListDirectoriesEntities represents the schema for listing directory entities.
type ArkIdentityListDirectoriesEntities struct {
	Directories  []string `json:"directories,omitempty" mapstructure:"directories" flag:"directories" desc:"Directories to search on"`
	EntityTypes  []string `json:"entity_types,omitempty" mapstructure:"entity_types" flag:"entity-types" desc:"Member types to search"`
	Search       string   `json:"search,omitempty" mapstructure:"search" flag:"search" desc:"Search string to use"`
	PageSize     int      `json:"page_size" mapstructure:"page_size" flag:"page-size" desc:"Page size to emit" default:"10000"`
	Limit        int      `json:"limit" mapstructure:"limit" flag:"limit" desc:"Limit amount to list" default:"10000"`
	MaxPageCount int      `json:"max_page_count" mapstructure:"max_page_count" flag:"max-page-size" desc:"Max page count to reach to" default:"-1"`
}
