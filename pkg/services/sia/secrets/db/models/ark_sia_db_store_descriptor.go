package models

// ArkSIADBStoreDescriptor represents the descriptor of a store in the Ark SIA DB.
type ArkSIADBStoreDescriptor struct {
	StoreID   string `json:"store_id,omitempty" mapstructure:"store_id" desc:"ID of the store"`
	StoreType string `json:"store_type,omitempty" mapstructure:"store_type" desc:"Type of the store"`
}
