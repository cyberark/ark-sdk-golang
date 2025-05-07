package db

// DatabasesEnginesToFamily is a map of database engine types to their family types.
var DatabasesEnginesToFamily = map[string]string{
	EngineTypeAuroraMysql:          FamilyTypeMySQL,
	EngineTypeAuroraPostgres:       FamilyTypePostgres,
	EngineTypeCustomSQLServerEE:    FamilyTypeMSSQL,
	EngineTypeCustomSQLServerSE:    FamilyTypeMSSQL,
	EngineTypeCustomSQLServerWeb:   FamilyTypeMSSQL,
	EngineTypeMariaDB:              FamilyTypeMariaDB,
	EngineTypeMariaDBSH:            FamilyTypeMariaDB,
	EngineTypeMSSQL:                FamilyTypeMSSQL,
	EngineTypeMSSQLSH:              FamilyTypeMSSQL,
	EngineTypeMySQL:                FamilyTypeMySQL,
	EngineTypeMySQLSH:              FamilyTypeMySQL,
	EngineTypeOracle:               FamilyTypeOracle,
	EngineTypeOracleEE:             FamilyTypeOracle,
	EngineTypeOracleSH:             FamilyTypeOracle,
	EngineTypeOracleEECDB:          FamilyTypeOracle,
	EngineTypeOracleSE2CDB:         FamilyTypeOracle,
	EngineTypeOracleSE2:            FamilyTypeOracle,
	EngineTypePostgres:             FamilyTypePostgres,
	EngineTypePostgresSH:           FamilyTypePostgres,
	EngineTypeSQLServer:            FamilyTypeMSSQL,
	EngineTypeSQLServerSH:          FamilyTypeMSSQL,
	EngineTypeDB2:                  FamilyTypeDB2,
	EngineTypeDB2SH:                FamilyTypeDB2,
	EngineTypeMongo:                FamilyTypeMongo,
	EngineTypeMongoSH:              FamilyTypeMongo,
	EngineTypeMSSQLSHVM:            FamilyTypeMSSQL,
	EngineTypeMSSQLAzureManaged:    FamilyTypeMSSQL,
	EngineTypeMSSQLAzureVM:         FamilyTypeMSSQL,
	EngineTypeMSSQLAWSEC2:          FamilyTypeMSSQL,
	EngineTypeMSSQLAWSRDS:          FamilyTypeMSSQL,
	EngineTypeDB2AWSRDS:            FamilyTypeDB2,
	EngineTypeDB2SHVM:              FamilyTypeDB2,
	EngineTypeOracleAWSRDS:         FamilyTypeOracle,
	EngineTypeOracleAWSVM:          FamilyTypeOracle,
	EngineTypeOracleSHVM:           FamilyTypeOracle,
	EngineTypeMariaDBSHVM:          FamilyTypeMariaDB,
	EngineTypeMariaDBAzureManaged:  FamilyTypeMariaDB,
	EngineTypeMariaDBAzureVM:       FamilyTypeMariaDB,
	EngineTypeMariaDBAWSVM:         FamilyTypeMariaDB,
	EngineTypeMariaDBAWSRDS:        FamilyTypeMariaDB,
	EngineTypeMariaDBAWSAurora:     FamilyTypeMariaDB,
	EngineTypeMySQLSHVM:            FamilyTypeMySQL,
	EngineTypeMySQLAzureManaged:    FamilyTypeMySQL,
	EngineTypeMySQLAzureVM:         FamilyTypeMySQL,
	EngineTypeMySQLAWSVM:           FamilyTypeMySQL,
	EngineTypeMySQLAWSRDS:          FamilyTypeMySQL,
	EngineTypeMySQLAWSAurora:       FamilyTypeMySQL,
	EngineTypePostgresSHVM:         FamilyTypePostgres,
	EngineTypePostgresAzureManaged: FamilyTypePostgres,
	EngineTypePostgresAzureVM:      FamilyTypePostgres,
	EngineTypePostgresAWSVM:        FamilyTypePostgres,
	EngineTypePostgresAWSRDS:       FamilyTypePostgres,
	EngineTypePostgresAWSAurora:    FamilyTypePostgres,
	EngineTypeMongoSHVM:            FamilyTypeMongo,
	EngineTypeMongoAWSDocDB:        FamilyTypeMongo,
	EngineTypeMongoAtlasManaged:    FamilyTypeMongo,
}

// DatabaseFamiliesDefaultPorts is a map of database family types to their default ports.
var DatabaseFamiliesDefaultPorts = map[string]int{
	FamilyTypePostgres: 5432,
	FamilyTypeOracle:   2484,
	FamilyTypeMSSQL:    1433,
	FamilyTypeMySQL:    3306,
	FamilyTypeMariaDB:  3306,
	FamilyTypeDB2:      50002,
	FamilyTypeMongo:    27017,
}

// DatabaseEngineTypes is a list of all possible database engine types.
var DatabaseEngineTypes = []string{
	EngineTypeAuroraPostgres,
	EngineTypeAuroraMysql,
	EngineTypeCustomSQLServerEE,
	EngineTypeCustomSQLServerSE,
	EngineTypeCustomSQLServerWeb,
	EngineTypeOracleEE,
	EngineTypeOracleEECDB,
	EngineTypeOracleSE2,
	EngineTypeOracleSE2CDB,
	EngineTypeSQLServer,
	EngineTypeOracle,
	EngineTypeMSSQL,
	EngineTypeMariaDB,
	EngineTypeMySQL,
	EngineTypePostgres,
	EngineTypeSQLServerSH,
	EngineTypeMSSQLSH,
	EngineTypeMySQLSH,
	EngineTypeMariaDBSH,
	EngineTypePostgresSH,
	EngineTypeOracleSH,
	EngineTypeDB2,
	EngineTypeDB2SH,
	EngineTypeMongo,
	EngineTypeMongoSH,
	EngineTypeMSSQLSHVM,
	EngineTypeMSSQLAzureManaged,
	EngineTypeMSSQLAzureVM,
	EngineTypeMSSQLAWSEC2,
	EngineTypeMSSQLAWSRDS,
	EngineTypeDB2AWSRDS,
	EngineTypeDB2SHVM,
	EngineTypeOracleAWSRDS,
	EngineTypeOracleAWSVM,
	EngineTypeOracleSHVM,
	EngineTypeMariaDBSHVM,
	EngineTypeMariaDBAzureManaged,
	EngineTypeMariaDBAzureVM,
	EngineTypeMariaDBAWSVM,
	EngineTypeMariaDBAWSRDS,
	EngineTypeMariaDBAWSAurora,
	EngineTypeMySQLSHVM,
	EngineTypeMySQLAzureManaged,
	EngineTypeMySQLAzureVM,
	EngineTypeMySQLAWSVM,
	EngineTypeMySQLAWSRDS,
	EngineTypeMySQLAWSAurora,
	EngineTypePostgresSHVM,
	EngineTypePostgresAzureManaged,
	EngineTypePostgresAzureVM,
	EngineTypePostgresAWSVM,
	EngineTypePostgresAWSRDS,
	EngineTypePostgresAWSAurora,
	EngineTypeMongoSHVM,
	EngineTypeMongoAWSDocDB,
	EngineTypeMongoAtlasManaged,
}

// DatabaseFamilyTypes is a list of all possible database family types.
var DatabaseFamilyTypes = []string{
	FamilyTypePostgres,
	FamilyTypeOracle,
	FamilyTypeMSSQL,
	FamilyTypeMySQL,
	FamilyTypeMariaDB,
	FamilyTypeDB2,
	FamilyTypeMongo,
	FamilyTypeUnknown,
}

// Possible database engine types.
const (
	EngineTypeAuroraPostgres       string = "aurora-postgresql"
	EngineTypeAuroraMysql          string = "aurora-mysql"
	EngineTypeCustomSQLServerEE    string = "custom-sqlserver-ee"
	EngineTypeCustomSQLServerSE    string = "custom-sqlserver-se"
	EngineTypeCustomSQLServerWeb   string = "custom-sqlserver-web"
	EngineTypeOracleEE             string = "oracle-ee"
	EngineTypeOracleEECDB          string = "oracle-ee-cdb"
	EngineTypeOracleSE2            string = "oracle-se2"
	EngineTypeOracleSE2CDB         string = "oracle-se2-cdb"
	EngineTypeSQLServer            string = "sqlserver"
	EngineTypeOracle               string = "oracle"
	EngineTypeMSSQL                string = "mssql"
	EngineTypeMariaDB              string = "mariadb"
	EngineTypeMySQL                string = "mysql"
	EngineTypePostgres             string = "postgres"
	EngineTypeSQLServerSH          string = "sqlserver-sh"
	EngineTypeMSSQLSH              string = "mssql-sh"
	EngineTypeMySQLSH              string = "mysql-sh"
	EngineTypeMariaDBSH            string = "mariadb-sh"
	EngineTypePostgresSH           string = "postgres-sh"
	EngineTypeOracleSH             string = "oracle-sh"
	EngineTypeDB2                  string = "db2"
	EngineTypeDB2SH                string = "db2-sh"
	EngineTypeMongo                string = "mongo"
	EngineTypeMongoSH              string = "mongo-sh"
	EngineTypeMSSQLSHVM            string = "mssql-sh-vm"
	EngineTypeMSSQLAzureManaged    string = "mssql-azure-managed"
	EngineTypeMSSQLAzureVM         string = "mssql-azure-vm"
	EngineTypeMSSQLAWSEC2          string = "mssql-aws-ec2"
	EngineTypeMSSQLAWSRDS          string = "mssql-aws-rds"
	EngineTypeDB2AWSRDS            string = "db2-aws-rds"
	EngineTypeDB2SHVM              string = "db2-sh-vm"
	EngineTypeOracleAWSRDS         string = "oracle-aws-rds"
	EngineTypeOracleAWSVM          string = "oracle-aws-vm"
	EngineTypeOracleSHVM           string = "oracle-sh-vm"
	EngineTypeMariaDBSHVM          string = "mariadb-sh-vm"
	EngineTypeMariaDBAzureManaged  string = "mariadb-azure-managed"
	EngineTypeMariaDBAzureVM       string = "mariadb-azure-vm"
	EngineTypeMariaDBAWSVM         string = "mariadb-aws-vm"
	EngineTypeMariaDBAWSRDS        string = "mariadb-aws-rds"
	EngineTypeMariaDBAWSAurora     string = "mariadb-aws-aurora"
	EngineTypeMySQLSHVM            string = "mysql-sh-vm"
	EngineTypeMySQLAzureManaged    string = "mysql-azure-managed"
	EngineTypeMySQLAzureVM         string = "mysql-azure-vm"
	EngineTypeMySQLAWSVM           string = "mysql-aws-vm"
	EngineTypeMySQLAWSRDS          string = "mysql-aws-rds"
	EngineTypeMySQLAWSAurora       string = "mysql-aws-aurora"
	EngineTypePostgresSHVM         string = "postgres-sh-vm"
	EngineTypePostgresAzureManaged string = "postgres-azure-managed"
	EngineTypePostgresAzureVM      string = "postgres-azure-vm"
	EngineTypePostgresAWSVM        string = "postgres-aws-vm"
	EngineTypePostgresAWSRDS       string = "postgres-aws-rds"
	EngineTypePostgresAWSAurora    string = "postgres-aws-aurora"
	EngineTypeMongoSHVM            string = "mongo-sh-vm"
	EngineTypeMongoAWSDocDB        string = "mongo-aws-docdb"
	EngineTypeMongoAtlasManaged    string = "mongo-atlas-managed"
)

// Possible database family types.
const (
	FamilyTypePostgres string = "Postgres"
	FamilyTypeOracle   string = "Oracle"
	FamilyTypeMSSQL    string = "MSSQL"
	FamilyTypeMySQL    string = "MySQL"
	FamilyTypeMariaDB  string = "MariaDB"
	FamilyTypeDB2      string = "DB2"
	FamilyTypeMongo    string = "Mongo"
	FamilyTypeUnknown  string = "Unknown"
)

// Possible database workspace types.
const (
	Cloud      string = "cloud"
	SelfHosted string = "self-hosted"
)

// ArkSIADBDatabaseProvider represents the structure for a database provider in the SIA workspace.
type ArkSIADBDatabaseProvider struct {
	ID        int    `json:"id" mapstructure:"id" flag:"id" desc:"ID of the provider"`
	Engine    string `json:"engine" mapstructure:"engine" flag:"engine" desc:"Engine type of the database provider"`
	Workspace string `json:"workspace" mapstructure:"workspace" flag:"workspace" desc:"Workspace of the database provider" choices:"cloud,self-hosted"`
	Family    string `json:"family" mapstructure:"family" flag:"family" desc:"Family of the database provider" choices:"Postgres,Oracle,MSSQL,MySQL,MariaDB,DB2,Mongo,Unknown"`
}
