package actions

// ArkServiceActionDefinition is a struct that defines the structure of an action in the Ark CLI.
type ArkServiceActionDefinition struct {
	ActionName   string                            `mapstructure:"action_name" json:"action_name" description:"Action name to be used in the cli commands"`
	Schemas      map[string]interface{}            `mapstructure:"schemas,omitempty" json:"schemas,omitempty" description:"Schemas for different cli actions for the definition"`
	Defaults     map[string]map[string]interface{} `mapstructure:"defaults,omitempty" json:"defaults,omitempty" description:"Defaults for the action schemas parameters"`
	AsyncActions []string                          `mapstructure:"async_actions,omitempty" json:"async_actions,omitempty" description:"List of async actions as part of the schemas"`
	Subactions   []*ArkServiceActionDefinition     `mapstructure:"subactions,omitempty" json:"subactions,omitempty" description:"Subactions to this action"`
}
