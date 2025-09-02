package models

// ArkIdentityListDirectories represents the schema for listing directory types.
type ArkIdentityListDirectories struct {
	Directories []string `json:"directories,omitempty" mapstructure:"directories" flag:"directories" desc:"Directories types to list" required:"true"`
}
