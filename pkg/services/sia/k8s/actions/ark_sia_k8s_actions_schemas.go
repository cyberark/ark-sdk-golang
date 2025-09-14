package actions

import siak8s "github.com/cyberark/ark-sdk-golang/pkg/services/sia/k8s/models"

// ActionToSchemaMap is a map that defines the mapping between K8S action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"generate-kubeconfig": &siak8s.ArkSIAK8SGenerateKubeconfig{},
}
