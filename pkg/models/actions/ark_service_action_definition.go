package actions

// ArkServiceActionType defines the type of action for an Ark service, such as CLI or Terraform.
type ArkServiceActionType string

// ArkServiceActionOperation defines the operation type for an Ark service action, such as create, read, update, delete, or state.
type ArkServiceActionOperation string

// Constants for ArkServiceActionType
const (
	ArkServiceActionTypeCLI                 ArkServiceActionType = "cli"
	ArkServiceActionTypeTerraformResource   ArkServiceActionType = "terraform-resource"
	ArkServiceActionTypeTerraformDataSource ArkServiceActionType = "terraform-data-source"
)

// Constants for ArkServiceActionOperation
const (
	CreateOperation ArkServiceActionOperation = "create"
	ReadOperation   ArkServiceActionOperation = "read"
	UpdateOperation ArkServiceActionOperation = "update"
	DeleteOperation ArkServiceActionOperation = "delete"
	StateOperation  ArkServiceActionOperation = "state"
)

// ArkServiceActionDefinition is an interface that defines the structure of an action for an Ark service.
type ArkServiceActionDefinition interface {
	ActionType() ArkServiceActionType
	ActionDefinitionName() string
}

// ArkServiceBaseActionDefinition is a struct that defines the base structure of an action using Ark SDK
type ArkServiceBaseActionDefinition struct {
	ActionName        string                 `mapstructure:"action_name" json:"action_name" desc:"Action name to be used in the cli commands"`
	ActionDescription string                 `mapstructure:"action_description,omitempty" json:"action_description,omitempty" desc:"Action description to be used in the cli commands"`
	ActionVersion     int64                  `mapstructure:"action_version,omitempty" json:"action_version,omitempty" desc:"Action version to be used in the cli commands"`
	Schemas           map[string]interface{} `mapstructure:"schemas,omitempty" json:"schemas,omitempty" desc:"Schemas for different cli actions for the definition"`
}

// ActionDefinitionName returns the name of the action definition.
func (a *ArkServiceBaseActionDefinition) ActionDefinitionName() string {
	return a.ActionName
}

// ArkServiceCLIActionDefinition is a struct that defines the structure of an action in the Ark CLI.
type ArkServiceCLIActionDefinition struct {
	ArkServiceBaseActionDefinition `mapstructure:",squash"`
	ActionAliases                  []string                          `mapstructure:"action_aliases,omitempty" json:"action_aliases,omitempty" desc:"Action aliases to be used in the cli commands"`
	Defaults                       map[string]map[string]interface{} `mapstructure:"defaults,omitempty" json:"defaults,omitempty" desc:"Defaults for the action schemas parameters"`
	AsyncActions                   []string                          `mapstructure:"async_actions,omitempty" json:"async_actions,omitempty" desc:"List of async actions as part of the schemas"`
	Subactions                     []*ArkServiceCLIActionDefinition  `mapstructure:"subactions,omitempty" json:"subactions,omitempty" desc:"Subactions to this action"`
}

// ActionType returns the type of action, which is CLI in this case.
func (a *ArkServiceCLIActionDefinition) ActionType() ArkServiceActionType {
	return ArkServiceActionTypeCLI
}

// ArkServiceBaseTerraformActionDefinition is a struct that defines the structure of an action in the Ark Terraform provider.
type ArkServiceBaseTerraformActionDefinition struct {
	ArkServiceBaseActionDefinition `mapstructure:",squash"`
	StateSchema                    interface{} `mapstructure:"state_schema,omitempty" json:"state_schema,omitempty" desc:"Schema for the state of the resource"`
	SensitiveAttributes            []string    `mapstructure:"sensitive_attributes,omitempty" json:"sensitive_attributes,omitempty" desc:"Used to set attributes as sensitive in the schema"`
	ExtraRequiredAttributes        []string    `mapstructure:"extra_required_attributes,omitempty" json:"extra_required_attributes,omitempty" desc:"Used to set attributes as required in the schema if not configured as validate required tag"`
	ComputedAsSetAttributes        []string    `mapstructure:"computed_as_set_attributes,omitempty" json:"computed_as_set_attributes,omitempty" desc:"Used to define list attributes as set attributes in the schema for non ordering unique collections"`
}

// ArkServiceTerraformResourceActionDefinition is a struct that defines the structure of a resource action in the Ark Terraform provider.
type ArkServiceTerraformResourceActionDefinition struct {
	ArkServiceBaseTerraformActionDefinition `mapstructure:",squash"`
	RawStateInference                       bool                                 `mapstructure:"raw_state_inference,omitempty" json:"raw_state_inference,omitempty" desc:"Used for cases where the schema is not a definitive one for the state"`
	ReadSchemaPath                          string                               `mapstructure:"read_schema_path,omitempty" json:"read_schema_path,omitempty" desc:"Used to find the inner schema within the state schema for the read schema input"`
	DeleteSchemaPath                        string                               `mapstructure:"delete_schema_path,omitempty" json:"delete_schema_path,omitempty" desc:"Used to find the inner schema within the state schema for the delete schema input"`
	SupportedOperations                     []ArkServiceActionOperation          `mapstructure:"supported_operations,omitempty" json:"supported_operations,omitempty" desc:"Defines the operations that this resource supports"`
	ActionsMappings                         map[ArkServiceActionOperation]string `mapstructure:"actions_mappings,omitempty" json:"actions_mappings,omitempty" desc:"Defines the mappings of operations to their action names"`
}

// ActionType returns the type of action, which is Terraform Resource in this case.
func (a *ArkServiceTerraformResourceActionDefinition) ActionType() ArkServiceActionType {
	return ArkServiceActionTypeTerraformResource
}

// ArkServiceTerraformDataSourceActionDefinition is a struct that defines the structure of a data source action in the Ark Terraform provider.
type ArkServiceTerraformDataSourceActionDefinition struct {
	ArkServiceBaseTerraformActionDefinition `mapstructure:",squash"`
	DataSourceAction                        string `mapstructure:"data_source_action,omitempty" json:"data_source_action,omitempty" desc:"Action name to be used for the data source"`
}

// ActionType returns the type of action, which is Terraform Data Source in this case.
func (a *ArkServiceTerraformDataSourceActionDefinition) ActionType() ArkServiceActionType {
	return ArkServiceActionTypeTerraformDataSource
}
