package models

// ArkUAPPolicyStatus defines the possible status types for a policy in UAP.
const (
	StatusTypeActive     = "Active"
	StatusTypeSuspended  = "Suspended"
	StatusTypeExpired    = "Expired"
	StatusTypeValidating = "Validating"
	StatusTypeError      = "Error"
)

// ArkUAPPolicyStatus represents the status details of a policy.
type ArkUAPPolicyStatus struct {
	Status            string `json:"status" mapstructure:"status" flag:"status" desc:"The status type of the policy" choices:"Active,Suspended,Expired,Validating,Error,Warning" default:"Active"`
	StatusCode        string `json:"status_code,omitempty" mapstructure:"status_code,omitempty" flag:"status-code" desc:"The status code of the policy" validate:"max=100"`
	StatusDescription string `json:"status_description,omitempty" mapstructure:"status_description,omitempty" flag:"status-description" desc:"The status description of the policy" validate:"max=1000"`
	Link              string `json:"link,omitempty" mapstructure:"link,omitempty" flag:"link" desc:"A documentation link for the policy status" validate:"max=255"`
}
