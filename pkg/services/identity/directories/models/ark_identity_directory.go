package models

// ArkIdentityDirectory represents the schema for an identity directory.
type ArkIdentityDirectory struct {
	Directory            string `json:"directory" mapstructure:"directory" flag:"directory" desc:"Name of the directory" required:"true"`
	DirectoryServiceUUID string `json:"directory_service_uuid" mapstructure:"directory_service_uuid" flag:"directory-service-uuid" desc:"ID of the directory" required:"true"`
}
