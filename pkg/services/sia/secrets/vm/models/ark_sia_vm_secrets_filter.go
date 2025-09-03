package models

// ArkSIAVMSecretsFilter represents the request to filter secrets in a VM.
type ArkSIAVMSecretsFilter struct {
	SecretTypes   []string               `json:"secret_types,omitempty" mapstructure:"secret_types,omitempty" flag:"secret-types" desc:"Type of secrets to filter"`
	Name          string                 `json:"name,omitempty" mapstructure:"name,omitempty" flag:"name" desc:"Name wildcard to filter with"`
	SecretDetails map[string]interface{} `json:"secret_details,omitempty" mapstructure:"secret_details,omitempty" flag:"secret-details" desc:"Secret details to filter with"`
	IsActive      bool                   `json:"is_active,omitempty" mapstructure:"is_active,omitempty" flag:"is-active" desc:"Filter only active / inactive secrets"`
}
