package actions

import dbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/db/models"

// ActionToSchemaMap is a map that defines the mapping between db action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"psql":                      &dbmodels.ArkSIADBPsqlExecution{},
	"mysql":                     &dbmodels.ArkSIADBMysqlExecution{},
	"sqlcmd":                    &dbmodels.ArkSIADBSqlcmdExecution{},
	"generate-oracle-tnsnames":  &dbmodels.ArkSIADBOracleGenerateAssets{},
	"generate-proxy-full-chain": &dbmodels.ArkSIADBProxyFullChainGenerateAssets{},
}
