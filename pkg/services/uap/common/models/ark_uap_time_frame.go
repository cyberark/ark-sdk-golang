package models

// ArkUAPTimeFrame represents the time frame for a policy.
type ArkUAPTimeFrame struct {
	FromTime string `json:"from_time,omitempty" mapstructure:"from_time,omitempty" flag:"from-time" desc:"Time from which the policy is effective"`
	ToTime   string `json:"to_time,omitempty" mapstructure:"to_time,omitempty" flag:"to-time" desc:"Time to which the policy is expired"`
}
