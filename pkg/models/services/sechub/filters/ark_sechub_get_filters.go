package filters

type ArkSecHubGetFilters struct {
	StoreID string `json:"store_id" mapstructure:"store_id" desc:"Secrets Store Id for Secrets Hub" flag:"store-id" validate:"required"`
}

// ArkSecHubGetFilterInfo holds the StoreId for the request.
// Get all the secrets filters related to a secret store source, by the secret store unique identifier.
type ArkSecHubGetFilter struct {
	StoreID  string `json:"store_id" mapstructure:"store_id" desc:"Secrets Store Id for Secrets Hub" flag:"store-id" validate:"required"`
	FilterID string `json:"filter_id" mapstructure:"filter_id" desc:"Filter ID for Secrets Hub" flag:"filter-id" validate:"required"`
}

// ArkSecHubFilterData represents the data field in the filter.
type ArkSecHubFilterData struct {
	SafeName string `json:"safe_name" mapstructure:"safe_name"`
}

// ArkSecHubFilter represents a single filter used to retrieve secrets from Ark Secrets Hub.
type ArkSecHubFilter struct {
	ID        string              `json:"id" mapstructure:"id"`
	Type      string              `json:"type" mapstructure:"type"`
	Data      ArkSecHubFilterData `json:"data" mapstructure:"data"`
	CreatedAt string              `json:"created_at" mapstructure:"created_at"`
	UpdatedAt string              `json:"updated_at" mapstructure:"updated_at"`
	CreatedBy string              `json:"created_by" mapstructure:"created_by"`
	UpdatedBy string              `json:"updated_by" mapstructure:"updated_by"`
}
