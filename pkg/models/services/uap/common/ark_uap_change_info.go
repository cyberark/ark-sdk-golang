package common

// ArkUAPChangeInfo represents the change information in UAP.
type ArkUAPChangeInfo struct {
	User string `json:"user,omitempty" mapstructure:"user,omitempty" flag:"user" desc:"Username of the user who made the change" validate:"omitempty,min=1,max=512"`
	Time string `json:"time,omitempty" mapstructure:"time,omitempty" flag:"time" desc:"Time of the change"`
}
