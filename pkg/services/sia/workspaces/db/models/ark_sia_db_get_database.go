package models

// ArkSIADBGetDatabase represents the request to retrieve a database in a workspace.
type ArkSIADBGetDatabase struct {
	ID   int    `json:"id,omitempty" mapstructure:"id,omitempty" flag:"id" desc:"Database id to get"`
	Name string `json:"name,omitempty" mapstructure:"name,omitempty" flag:"name" desc:"Database name to get"`
}
