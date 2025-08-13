package common

// ArkUAPTimeCondition represents the time conditions for a policy.
type ArkUAPTimeCondition struct {
	DaysOfTheWeek []int  `json:"days_of_the_week" mapstructure:"days_of_the_week" flag:"days-of-the-week" desc:"The days that the policy will be active" validate:"min=0,max=6" default:"0,1,2,3,4,5,6"`
	FromHour      string `json:"from_hour,omitempty" mapstructure:"from_hour,omitempty" flag:"from-hour" desc:"The policy will be active from hour" validate:"regexp=^\\w+$"`
	ToHour        string `json:"to_hour,omitempty" mapstructure:"to_hour,omitempty" flag:"to-hour" desc:"The policy will be active to time" validate:"regexp=^\\w+$"`
}
