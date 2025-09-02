package models

// DatabasesFamiliesToDefaultAuthMethod is a map of database family types to their default authentication methods.
var DatabasesFamiliesToDefaultAuthMethod = map[string]string{
	FamilyTypePostgres: LocalEphemeralUser,
	FamilyTypeOracle:   LocalEphemeralUser,
	FamilyTypeMSSQL:    ADEphemeralUser,
	FamilyTypeMySQL:    LocalEphemeralUser,
	FamilyTypeMariaDB:  LocalEphemeralUser,
	FamilyTypeDB2:      ADEphemeralUser,
	FamilyTypeMongo:    LocalEphemeralUser,
}

// Possible database authentication methods
const (
	ADEphemeralUser      string = "ad_ephemeral_user"
	LocalEphemeralUser   string = "local_ephemeral_user"
	RDSIAMAuthentication string = "rds_iam_authentication"
	AtlasEphemeralUser   string = "atlas_ephemeral_user"
)

// ArkSIADBAuthMethod represents the authentication method used for a database target.
type ArkSIADBAuthMethod struct {
	ID             int      `json:"id" mapstructure:"id" flag:"id" desc:"ID of the authentication method on the database"`
	AuthMethodType string   `json:"auth_method_type" mapstructure:"auth_method_type" flag:"auth-method-type" desc:"Type / name of the authentication method" choices:"ad_ephemeral_user,local_ephemeral_user,rds_iam_authentication,atlas_ephemeral_user"`
	Description    string   `json:"description" mapstructure:"description" flag:"description" desc:"Description about the authentication method"`
	Workspaces     []string `json:"workspaces" mapstructure:"workspaces" flag:"workspaces" desc:"Workspaces this authentication method is used in"`
}

// ArkSIADBDatabaseAuthMethod represents the relation between a database type and an authentication method.
type ArkSIADBDatabaseAuthMethod struct {
	ID             int                `json:"id" mapstructure:"id" flag:"id" desc:"ID of the relation between the authentication method and the database type"`
	ProviderFamily string             `json:"provider_family" mapstructure:"provider_family" flag:"provider-family" desc:"Name of the database family this authentication method is used for" choices:"Postgres,Oracle,MSSQL,MySQL,MariaDB,DB2,Mongo,Unknown"`
	AuthMethod     ArkSIADBAuthMethod `json:"auth_method" mapstructure:"auth_method" flag:"auth-method" desc:"The actual authentication method"`
	MethodEnabled  bool               `json:"method_enabled" mapstructure:"method_enabled" flag:"method-enabled" desc:"Whether this authentication method is enabled or not"`
}

// ArkSIADBDatabaseTargetConfiguredAuthMethod represents the relation between a database target and a configured authentication method.
type ArkSIADBDatabaseTargetConfiguredAuthMethod struct {
	DatabaseAuthMethod     ArkSIADBDatabaseAuthMethod `json:"database_auth_method" mapstructure:"database_auth_method" flag:"database-auth-method" desc:"Identifier for the configured auth method"`
	DatabaseTargetID       int                        `json:"database_target_id" mapstructure:"database_target_id" flag:"database-target-id" desc:"Database target identifier"`
	ConfiguredAuthMethodID int                        `json:"configured_auth_method_id,omitempty" mapstructure:"configured_auth_method_id,omitempty" flag:"configured-auth-method-id" desc:"The configured auth method id for the target"`
}
