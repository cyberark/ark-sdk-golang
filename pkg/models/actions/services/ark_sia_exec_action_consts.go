package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	siaaccess "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/access"
	siak8s "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/k8s"
	siasecretsvm "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/secrets/vm"
	siasso "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/sso"
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

// TargetSetsAction is a struct that defines the TargetSets action for the Ark service.
var TargetSetsAction = actions.ArkServiceActionDefinition{
	ActionName: "target-sets",
	Schemas:    TargetSetsActionToSchemaMap,
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

// SecretsAction is a struct that defines the Secrets action for the Ark service.
var SecretsAction = actions.ArkServiceActionDefinition{
	ActionName: "secrets",
	Subactions: []*actions.ArkServiceActionDefinition{
		&SecretsVMAction,
	},
}

// AccessActionToSchemaMap is a map that defines the mapping between Access action names and their corresponding schema types.
var AccessActionToSchemaMap = map[string]interface{}{
	"connector-setup-script":      &siaaccess.ArkSIAGetConnectorSetupScript{},
	"install-connector":           &siaaccess.ArkSIAInstallConnector{},
	"uninstall-connector":         &siaaccess.ArkSIAUninstallConnector{},
	"test-connector-reachability": &siaaccess.ArkSiaTestConnectorReachability{},
}

// AccessAction is a struct that defines the Access action for the Ark service.
var AccessAction = actions.ArkServiceActionDefinition{
	ActionName: "access",
	Schemas:    AccessActionToSchemaMap,
}

// SIAActions is a struct that defines the SIA actions for the Ark service.
var SIAActions = &actions.ArkServiceActionDefinition{
	ActionName: "sia",
	Subactions: []*actions.ArkServiceActionDefinition{
		&SSOAction,
		&K8SAction,
		&TargetSetsAction,
		&SecretsAction,
		&AccessAction,
	},
}
