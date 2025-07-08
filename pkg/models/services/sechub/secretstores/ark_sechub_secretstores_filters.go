package secretstores

type ArkSecHubSecretStoresFilters struct {
	Behavior string `json:"behavior,omitempty" mapstructure:"behavior,omitempty" desc:"The type of secret store (SECRETS_TARGET,SECRETS_SOURCE)" default:"SECRETS_TARGET" choices:"SECRETS_TARGET,SECRETS_SOURCE"`
	Filters  string `json:"filters,omitempty" mapstructure:"filters,omitempty" desc:"Secret store filters. Example: --Filter 'type EQ AWS_ASM' --Filter 'data.accountId EQ 123412341234'"`
}
