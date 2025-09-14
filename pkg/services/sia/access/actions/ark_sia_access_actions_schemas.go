package actions

import accessmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/access/models"

// ActionToSchemaMap is a map that defines the mapping between Access action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"connector-setup-script":      &accessmodels.ArkSIAGetConnectorSetupScript{},
	"install-connector":           &accessmodels.ArkSIAInstallConnector{},
	"uninstall-connector":         &accessmodels.ArkSIAUninstallConnector{},
	"test-connector-reachability": &accessmodels.ArkSIATestConnectorReachability{},
	"delete-connector":            &accessmodels.ArkSIADeleteConnector{},
}
