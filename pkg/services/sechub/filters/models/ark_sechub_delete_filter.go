package models

// ArkSecHubDeleteFilter holds the StoreID and FilterID for the request.
type ArkSecHubDeleteFilter struct {
	StoreID  string `json:"store_id" mapstructure:"store_id" desc:"Secrets Store Id for Secrets Hub" flag:"store-id" validate:"required"`
	FilterID string `json:"filter_id,omitempty" mapstructure:"filter_id,omitempty" desc:"Filter ID for Secrets Hub" flag:"filter-id" validate:"required"`
}
