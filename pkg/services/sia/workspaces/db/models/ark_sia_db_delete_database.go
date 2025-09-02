package models

// ArkSIADBDeleteDatabase represents the request to delete a database.
type ArkSIADBDeleteDatabase struct {
	ID   int    `json:"id,omitempty" mapstructure:"id,omitempty" flag:"id" desc:"Database id to delete"`
	Name string `json:"name,omitempty" mapstructure:"name,omitempty" flag:"name" desc:"Database name to delete"`
}
