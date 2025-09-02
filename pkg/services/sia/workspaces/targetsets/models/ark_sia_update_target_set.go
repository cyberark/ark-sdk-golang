package models

// ArkSIAUpdateTargetSet represents the structure for updating a target set in the SIA workspace.
type ArkSIAUpdateTargetSet struct {
	ID                          string `json:"id" mapstructure:"id" flag:"id" desc:"The target set id" validate:"required"`
	Name                        string `json:"name,omitempty" mapstructure:"name" flag:"name" desc:"New name of the target set to update"`
	Description                 string `json:"description,omitempty" mapstructure:"description,omitempty" flag:"description" desc:"Updated description of the target set"`
	ProvisionFormat             string `json:"provision_format,omitempty" mapstructure:"provision_format,omitempty" flag:"provision-format" desc:"New provisioning format for the target set"`
	EnableCertificateValidation bool   `json:"enable_certificate_validation,omitempty" mapstructure:"enable_certificate_validation,omitempty" flag:"enable-certificate-validation" desc:"Updated enabling certificate validation"`
	SecretType                  string `json:"secret_type,omitempty" mapstructure:"secret_type,omitempty" flag:"secret-type" desc:"Secret type to update (ProvisionerUser,PCloudAccount)" choices:"ProvisionerUser,PCloudAccount"`
	SecretID                    string `json:"secret_id,omitempty" mapstructure:"secret_id,omitempty" flag:"secret-id" desc:"Secret id to update"`
	Type                        string `json:"type,omitempty" mapstructure:"type,omitempty" flag:"type" desc:"Type of the target set" choices:"Domain,Suffix,Target"`
}
