package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	accessmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/access/models"
	dbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/db/models"
	siak8s "github.com/cyberark/ark-sdk-golang/pkg/services/sia/k8s/models"
	secretsdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/db/models"
	secretsvmmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/vm/models"
	siasshca "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sshca/models"
	ssomodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso/models"
	workspacesdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db/models"
	targetsetsmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/targetsets/models"
)

// SSOActionToSchemaMap is a map that defines the mapping between SSO action names and their corresponding schema types.
var SSOActionToSchemaMap = map[string]interface{}{
	"short-lived-password":           &ssomodels.ArkSIASSOGetShortLivedPassword{},
	"short-lived-client-certificate": &ssomodels.ArkSIASSOGetShortLivedClientCertificate{},
	"short-lived-oracle-wallet":      &ssomodels.ArkSIASSOGetShortLivedOracleWallet{},
	"short-lived-rdp-file":           &ssomodels.ArkSIASSOGetShortLivedRDPFile{},
	"short-lived-token-info":         &ssomodels.ArkSIASSOGetTokenInfo{},
	"short-lived-ssh-key":            &ssomodels.ArkSIASSOGetSSHKey{},
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
	"add-target-set":          &targetsetsmodels.ArkSIAAddTargetSet{},
	"bulk-add-target-sets":    &targetsetsmodels.ArkSIABulkAddTargetSets{},
	"delete-target-set":       &targetsetsmodels.ArkSIADeleteTargetSet{},
	"bulk-delete-target-sets": &targetsetsmodels.ArkSIABulkDeleteTargetSets{},
	"update-target-set":       &targetsetsmodels.ArkSIAUpdateTargetSet{},
	"list-target-sets":        nil,
	"list-target-sets-by":     &targetsetsmodels.ArkSIATargetSetsFilter{},
	"target-set":              &targetsetsmodels.ArkSIAGetTargetSet{},
	"target-sets-stats":       nil,
}

// DBWorkspaceActionToSchemaMap is a map that defines the mapping between DB workspace action names and their corresponding schema types.
var DBWorkspaceActionToSchemaMap = map[string]interface{}{
	"add-database":      &workspacesdbmodels.ArkSIADBAddDatabase{},
	"delete-database":   &workspacesdbmodels.ArkSIADBDeleteDatabase{},
	"update-database":   &workspacesdbmodels.ArkSIADBUpdateDatabase{},
	"database":          &workspacesdbmodels.ArkSIADBGetDatabase{},
	"list-databases":    nil,
	"list-databases-by": &workspacesdbmodels.ArkSIADBDatabasesFilter{},
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
	"add-secret":      &secretsvmmodels.ArkSIAVMAddSecret{},
	"change-secret":   &secretsvmmodels.ArkSIAVMChangeSecret{},
	"delete-secret":   &secretsvmmodels.ArkSIAVMDeleteSecret{},
	"list-secrets":    nil,
	"list-secrets-by": &secretsvmmodels.ArkSIAVMSecretsFilter{},
	"secret":          &secretsvmmodels.ArkSIAVMGetSecret{},
	"secrets-stats":   nil,
}

// SecretsVMAction is a struct that defines the Secrets VM action for the Ark service.
var SecretsVMAction = actions.ArkServiceActionDefinition{
	ActionName: "vm",
	Schemas:    SecretsVMActionToSchemaMap,
}

// SecretsDBActionToSchemaMap is a map that defines the mapping between Secrets DB action names and their corresponding schema types.
var SecretsDBActionToSchemaMap = map[string]interface{}{
	"add-secret":      &secretsdbmodels.ArkSIADBAddSecret{},
	"update-secret":   &secretsdbmodels.ArkSIADBUpdateSecret{},
	"delete-secret":   &secretsdbmodels.ArkSIADBDeleteSecret{},
	"list-secrets":    nil,
	"list-secrets-by": &secretsdbmodels.ArkSIADBSecretsFilter{},
	"enable-secret":   &secretsdbmodels.ArkSIADBEnableSecret{},
	"disable-secret":  &secretsdbmodels.ArkSIADBDisableSecret{},
	"secret":          &secretsdbmodels.ArkSIADBGetSecret{},
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
	"connector-setup-script":      &accessmodels.ArkSIAGetConnectorSetupScript{},
	"install-connector":           &accessmodels.ArkSIAInstallConnector{},
	"uninstall-connector":         &accessmodels.ArkSIAUninstallConnector{},
	"test-connector-reachability": &accessmodels.ArkSIATestConnectorReachability{},
	"delete-connector":            &accessmodels.ArkSIADeleteConnector{},
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
	"public-key":             &siasshca.ArkSIAGetSSHPublicKey{},
	"public-key-script":      &siasshca.ArkSIAGetSSHPublicKey{},
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
