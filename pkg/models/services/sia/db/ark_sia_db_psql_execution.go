package db

// ArkSIADBPsqlExecution defines the structure for executing PostgreSQL commands in the ArkDBA context.
type ArkSIADBPsqlExecution struct {
	ArkSIADBBaseExecution `mapstructure:",squash"`
	PsqlPath              string `json:"psql_path" mapstructure:"psql_path" flag:"psql-path" desc:"Path to the psql executable" default:"psql"`
}
