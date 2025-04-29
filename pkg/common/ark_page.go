package common

// ArkPage is a generic struct representing a paginated response from the Ark service.
type ArkPage[T any] struct {
	Items []*T `json:"items" mapstructure:"items"`
}
