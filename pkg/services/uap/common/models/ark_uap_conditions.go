package models

// ArkUAPConditions represents the conditions for UAP policies.
type ArkUAPConditions struct {
	AccessWindow       ArkUAPTimeCondition `json:"access_window" mapstructure:"access_window" flag:"access-window" desc:"Indicate the time frame that the policy will be active"`
	MaxSessionDuration int                 `json:"max_session_duration" mapstructure:"max_session_duration" flag:"max-session-duration" desc:"Session length" validate:"min=1,max=24" default:"1"`
}
