package secretstores

// ArkSecHubSecretStoresFilters
type ArkSecHubSecretStoresFilters struct {
	Behavior string   `json:"behavior,omitempty" mapstructure:"behavior,omitempty" desc:"Secret store behavior. Allowed Values: SECRETS_TARGET, SECRETS_SOURCE. Default Value: SECRETS_TARGET" default:"SECRETS_TARGET"`
	Filters  []string `json:"filters,omitempty" mapstructure:"filters,omitempty" desc:"Secret store filters. Example: --Filter 'type EQ AWS_ASM' --Filter 'data.accountId EQ 123412341234'"`
}
