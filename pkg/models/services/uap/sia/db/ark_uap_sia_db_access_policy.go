package db

import sia "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/sia/common"

// ArkUAPSIADBAccessPolicy represents a DB access policy for SIA.
type ArkUAPSIADBAccessPolicy struct {
	sia.ArkUAPSIACommonAccessPolicy `mapstructure:",squash"`
	Targets                         map[string]ArkUAPSIADBTargets `json:"targets,omitempty" mapstructure:"targets,omitempty" flag:"targets" desc:"The targets of the db access policy" choices:"FQDN/IP"`
}
