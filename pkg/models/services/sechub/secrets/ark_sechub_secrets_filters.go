package secrets

// ArkSecHubSecretsFilter represents the filter options for accounts.
type ArkSecHubSecretsFilter struct {
	Projection string `json:"projection,omitempty" mapstructure:"projection,omitempty" desc:"Whether to use extended projection or not - Allowed Values: 'EXTEND','REGULAR'" flag:"projection" validate:"oneof=EXTEND REGULAR"`
	Filter     string `json:"filter,omitempty" mapstructure:"filter,omitempty" desc:"Filter to apply" flag:"filter"`
	Sort       string `json:"sort,omitempty" mapstructure:"sort,omitempty" desc:"Sort results by given key" flag:"sort"`
	Offset     int    `json:"offset,omitempty" mapstructure:"offset,omitempty" desc:"Offset to the accounts list (min=0,max=1500000)" flag:"offset" default:"0" validate:"min=0,max=1500000"`
	Limit      int    `json:"limit,omitempty" mapstructure:"limit,omitempty" desc:"Limit of results (min=1,max=1000)" flag:"limit" validate:"min=1,max=1000"`
}
