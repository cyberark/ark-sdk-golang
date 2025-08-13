package actions

// ArkServiceActionDefinition is a struct that defines the structure of an action in the Ark CLI.
type ArkServiceActionDefinition struct {
	ActionName   string                            `mapstructure:"action_name" json:"action_name" desc:"Action name to be used in the cli commands"`
	Schemas      map[string]interface{}            `mapstructure:"schemas,omitempty" json:"schemas,omitempty" desc:"Schemas for different cli actions for the definition"`
	Defaults     map[string]map[string]interface{} `mapstructure:"defaults,omitempty" json:"defaults,omitempty" desc:"Defaults for the action schemas parameters"`
	AsyncActions []string                          `mapstructure:"async_actions,omitempty" json:"async_actions,omitempty" desc:"List of async actions as part of the schemas"`
	Subactions   []*ArkServiceActionDefinition     `mapstructure:"subactions,omitempty" json:"subactions,omitempty" desc:"Subactions to this action"`
}
