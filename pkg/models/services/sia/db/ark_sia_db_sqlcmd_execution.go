package db

// ArkSIADBSqlcmdExecution defines the structure for executing SQLCMD commands in the ArkDBA context.
type ArkSIADBSqlcmdExecution struct {
	ArkSIADBBaseExecution `mapstructure:",squash"`
	SqlcmdPath            string `json:"sqlcmd_path" mapstructure:"sqlcmd_path" flag:"sqlcmd-path" desc:"Path to the sqlcmd executable" default:"sqlcmd"`
}
