package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	siaaccess "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/access"
	siak8s "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/k8s"
	siasecretsdb "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/secrets/db"
	siasecretsvm "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/secrets/vm"
	siasso "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/sso"
	siaworkspacesdb "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/workspaces/db"
	siatargetsets "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/workspaces/targetsets"
)

// SSOActionToSchemaMap is a map that defines the mapping between SSO action names and their corresponding schema types.
var SSOActionToSchemaMap = map[string]interface{}{
	"short-lived-password":           &siasso.ArkSIASSOGetShortLivedPassword{},
	"short-lived-client-certificate": &siasso.ArkSIASSOGetShortLivedClientCertificate{},
	"short-lived-oracle-wallet":      &siasso.ArkSIASSOGetShortLivedOracleWallet{},
	"short-lived-rdp-file":           &siasso.ArkSIASSOGetShortLivedRDPFile{},
	"short-lived-token-info":         &siasso.ArkSIASSOGetTokenInfo{},
	"short-lived-ssh-key":            &siasso.ArkSIASSOGetSSHKey{},
}

// SSOAction is a struct that defines the SSO action for the Ark service.
var SSOAction = actions.ArkServiceActionDefinition{
	ActionName: "sso",
	Schemas:    SSOActionToSchemaMap,
}

// K8SActionToSchemaMap is a map that defines the mapping between K8S action names and their corresponding schema types.
var K8SActionToSchemaMap = map[string]interface{}{
	"generate-kubeconfig": &siak8s.ArkSIAK8SGenerateKubeconfig{},
}

// K8SAction is a struct that defines the K8S action for the Ark service.
var K8SAction = actions.ArkServiceActionDefinition{
	ActionName: "k8s",
	Schemas:    K8SActionToSchemaMap,
}

// TargetSetsActionToSchemaMap is a map that defines the mapping between TargetSets action names and their corresponding schema types.
var TargetSetsActionToSchemaMap = map[string]interface{}{
	"add-target-set":          &siatargetsets.ArkSIAAddTargetSet{},
	"bulk-add-target-sets":    &siatargetsets.ArkSIABulkAddTargetSets{},
	"delete-target-set":       &siatargetsets.ArkSIADeleteTargetSet{},
	"bulk-delete-target-sets": &siatargetsets.ArkSIABulkDeleteTargetSets{},
	"update-target-set":       &siatargetsets.ArkSIAUpdateTargetSet{},
	"list-target-sets":        nil,
	"list-target-sets-by":     &siatargetsets.ArkSIATargetSetsFilter{},
	"target-set":              &siatargetsets.ArkSIAGetTargetSet{},
	"target-sets-stats":       nil,
}

// DBWorkspaceActionToSchemaMap is a map that defines the mapping between DB workspace action names and their corresponding schema types.
var DBWorkspaceActionToSchemaMap = map[string]interface{}{
	"add-database":      &siaworkspacesdb.ArkSIADBAddDatabase{},
	"delete-database":   &siaworkspacesdb.ArkSIADBDeleteDatabase{},
	"update-database":   &siaworkspacesdb.ArkSIADBUpdateDatabase{},
	"database":          &siaworkspacesdb.ArkSIADBGetDatabase{},
	"list-databases":    nil,
	"list-databases-by": &siaworkspacesdb.ArkSIADBDatabasesFilter{},
	"databases-stats":   nil,
	"list-engine-types": nil,
	"list-family-types": nil,
}

// TargetSetsAction is a struct that defines the TargetSets action for the Ark service.
var TargetSetsAction = actions.ArkServiceActionDefinition{
	ActionName: "target-sets",
	Schemas:    TargetSetsActionToSchemaMap,
}

// DBWorkspaceAction is a struct that defines the DB workspace action for the Ark service.
var DBWorkspaceAction = actions.ArkServiceActionDefinition{
	ActionName: "db",
	Schemas:    DBWorkspaceActionToSchemaMap,
}

// WorkspacesAction is a struct that defines the Workspaces action for the Ark service.
var WorkspacesAction = actions.ArkServiceActionDefinition{
	ActionName: "workspaces",
	Subactions: []*actions.ArkServiceActionDefinition{
		&TargetSetsAction,
		&DBWorkspaceAction,
	},
}

// SecretsVMActionToSchemaMap is a map that defines the mapping between Secrets VM action names and their corresponding schema types.
var SecretsVMActionToSchemaMap = map[string]interface{}{
	"add-secret":      &siasecretsvm.ArkSIAVMAddSecret{},
	"change-secret":   &siasecretsvm.ArkSIAVMChangeSecret{},
	"delete-secret":   &siasecretsvm.ArkSIAVMDeleteSecret{},
	"list-secrets":    nil,
	"list-secrets-by": &siasecretsvm.ArkSIAVMSecretsFilter{},
	"secret":          &siasecretsvm.ArkSIAVMGetSecret{},
	"secrets-stats":   nil,
}

// SecretsVMAction is a struct that defines the Secrets VM action for the Ark service.
var SecretsVMAction = actions.ArkServiceActionDefinition{
	ActionName: "vm",
	Schemas:    SecretsVMActionToSchemaMap,
}

// SecretsDBActionToSchemaMap is a map that defines the mapping between Secrets DB action names and their corresponding schema types.
var SecretsDBActionToSchemaMap = map[string]interface{}{
	"add-secret":      &siasecretsdb.ArkSIADBAddSecret{},
	"update-secret":   &siasecretsdb.ArkSIADBUpdateSecret{},
	"delete-secret":   &siasecretsdb.ArkSIADBDeleteSecret{},
	"list-secrets":    nil,
	"list-secrets-by": &siasecretsdb.ArkSIADBSecretsFilter{},
	"enable-secret":   &siasecretsdb.ArkSIADBEnableSecret{},
	"disable-secret":  &siasecretsdb.ArkSIADBDisableSecret{},
	"secret":          &siasecretsdb.ArkSIADBGetSecret{},
	"secrets-stats":   nil,
}

// SecretsDBAction is a struct that defines the Secrets DB action for the Ark service.
var SecretsDBAction = actions.ArkServiceActionDefinition{
	ActionName: "db",
	Schemas:    SecretsDBActionToSchemaMap,
}

// SecretsAction is a struct that defines the Secrets action for the Ark service.
var SecretsAction = actions.ArkServiceActionDefinition{
	ActionName: "secrets",
	Subactions: []*actions.ArkServiceActionDefinition{
		&SecretsVMAction,
		&SecretsDBAction,
	},
}

// AccessActionToSchemaMap is a map that defines the mapping between Access action names and their corresponding schema types.
var AccessActionToSchemaMap = map[string]interface{}{
	"connector-setup-script":      &siaaccess.ArkSIAGetConnectorSetupScript{},
	"install-connector":           &siaaccess.ArkSIAInstallConnector{},
	"uninstall-connector":         &siaaccess.ArkSIAUninstallConnector{},
	"test-connector-reachability": &siaaccess.ArkSIATestConnectorReachability{},
	"delete-connector":            &siaaccess.ArkSIADeleteConnector{},
}

// AccessAction is a struct that defines the Access action for the Ark service.
var AccessAction = actions.ArkServiceActionDefinition{
	ActionName: "access",
	Schemas:    AccessActionToSchemaMap,
}

// SSHCaActionToSchemaMap is a map that defines the mapping between ssh-ca action names and their corresponding schema types.
var SSHCaActionToSchemaMap = map[string]interface{}{
	"generate-new-ca":        nil,
	"deactivate-previous-ca": nil,
	"reactivate-previous-ca": nil,
}

// SSHCaAction is a struct that defines the Access action for the Ark service.
var SSHCaAction = actions.ArkServiceActionDefinition{
	ActionName: "ssh-ca",
	Schemas:    SSHCaActionToSchemaMap,
}

// DbActionToSchemaMap is a map that defines the mapping between ssh-ca action names and their corresponding schema types.
var DbActionToSchemaMap = map[string]interface{}{
	"psql":                      &dbmodels.ArkSIADBPsqlExecution{},
	"mysql":                     &dbmodels.ArkSIADBMysqlExecution{},
	"sqlcmd":                    &dbmodels.ArkSIADBSqlcmdExecution{},
	"generate-oracle-tnsnames":  &dbmodels.ArkSIADBOracleGenerateAssets{},
	"generate-proxy-full-chain": &dbmodels.ArkSIADBProxyFullChainGenerateAssets{},
}

// DbAction is a struct that defines the Access action for the Ark service.
var DbAction = actions.ArkServiceActionDefinition{
	ActionName: "db",
	Schemas:    DbActionToSchemaMap,
}

// SIAActions is a struct that defines the SIA actions for the Ark service.
var SIAActions = &actions.ArkServiceActionDefinition{
	ActionName: "sia",
	Subactions: []*actions.ArkServiceActionDefinition{
		&SSOAction,
		&K8SAction,
		&WorkspacesAction,
		&SecretsAction,
		&AccessAction,
		&SSHCaAction,
		&DbAction,
	},
}
