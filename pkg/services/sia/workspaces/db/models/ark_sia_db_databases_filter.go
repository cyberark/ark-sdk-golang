package models

// ArkSIADBDatabasesFilter represents the filter criteria for retrieving databases in a workspace.
type ArkSIADBDatabasesFilter struct {
	Name              string        `json:"name,omitempty" mapstructure:"name,omitempty" flag:"name" desc:"Name of the database to filter on"`
	ProviderFamily    string        `json:"provider_family,omitempty" mapstructure:"provider_family,omitempty" flag:"provider-family" desc:"List filter by family" choices:"Postgres,Oracle,MSSQL,MySQL,MariaDB,DB2,Mongo,Unknown"`
	ProviderEngine    string        `json:"provider_engine,omitempty" mapstructure:"provider_engine,omitempty" flag:"provider-engine" desc:"List filter by engine"`
	ProviderWorkspace string        `json:"provider_workspace,omitempty" mapstructure:"provider_workspace,omitempty" flag:"provider-workspace" desc:"List filter by workspace" choices:"AWS,AZURE,GCP,ON-PREMISE,ATLAS"`
	AuthMethods       []string      `json:"auth_methods,omitempty" mapstructure:"auth_methods,omitempty" flag:"auth-methods" desc:"Auth method types to filter on" choices:"ad_ephemeral_user,local_ephemeral_user,rds_iam_authentication,atlas_ephemeral_user"`
	Tags              []ArkSIADBTag `json:"tags,omitempty" mapstructure:"tags,omitempty" flag:"tags" desc:"List filter by tags"`
	DBWarningsFilter  string        `json:"db_warnings_filter,omitempty" mapstructure:"db_warnings_filter,omitempty" flag:"db-warnings-filter" desc:"Filter by databases who are with warnings / incomplete" choices:"no_certificates,no_secrets,any_error"`
}
