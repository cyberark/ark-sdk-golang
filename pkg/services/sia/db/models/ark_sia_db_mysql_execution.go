package models

// ArkSIADBMysqlExecution defines the structure for executing MySQL commands in the ArkDBA context.
type ArkSIADBMysqlExecution struct {
	ArkSIADBBaseExecution `mapstructure:",squash"`
	MysqlPath             string `json:"mysql_path" mapstructure:"mysql_path" flag:"mysql-path" desc:"Path to the mysql executable" default:"mysql"`
}
