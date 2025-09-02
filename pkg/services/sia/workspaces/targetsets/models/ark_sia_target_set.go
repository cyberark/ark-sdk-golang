package models

// Possible target set types
const (
	Domain = "Domain"
	Suffix = "Suffix"
	Target = "Target"
)

// ArkSIATargetSet represents the structure for a target set in the SIA workspace.
type ArkSIATargetSet struct {
	ID                          string `json:"id" mapstructure:"id" flag:"id" desc:"The target set id" validate:"required"`
	Name                        string `json:"name" mapstructure:"name" flag:"name" desc:"The actual target set name / url" validate:"required"`
	Description                 string `json:"description,omitempty" mapstructure:"description,omitempty" flag:"description" desc:"Description about the target set"`
	ProvisionFormat             string `json:"provision_format,omitempty" mapstructure:"provision_format,omitempty" flag:"provision-format" desc:"Provisioning format for the target set ephemeral users"`
	EnableCertificateValidation bool   `json:"enable_certificate_validation,omitempty" mapstructure:"enable_certificate_validation,omitempty" flag:"enable-certificate-validation" desc:"Whether to enable certificate validation for the target set"`
	SecretType                  string `json:"secret_type,omitempty" mapstructure:"secret_type,omitempty" flag:"secret-type" desc:"Secret type of the target set" choices:"ProvisionerUser,PCloudAccount"`
	SecretID                    string `json:"secret_id,omitempty" mapstructure:"secret_id,omitempty" flag:"secret-id" desc:"Secret id of the target set"`
	Type                        string `json:"type" mapstructure:"type" flag:"type" desc:"Type of the target set" validate:"required" choices:"Domain,Suffix,Target"`
}
