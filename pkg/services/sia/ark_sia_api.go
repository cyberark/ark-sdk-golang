package sia

import (
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia/access"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia/k8s"
	dbsecrets "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/db"
	vmsecrets "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/vm"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/targetsets"
)

// ArkSIAAPI is a struct that provides access to the Ark SIA API as a wrapped set of services.
type ArkSIAAPI struct {
	ssoService          *sso.ArkSIASSOService
	k8sService          *k8s.ArkSIAK8SService
	targetSetsService   *targetsets.ArkSIAWorkspacesTargetSetsService
	workspacesDBService *db.ArkSIAWorkspacesDBService
	vmSecretsService    *vmsecrets.ArkSIASecretsVMService
	dbSecretsService    *dbsecrets.ArkSIASecretsDBService
	accessService       *access.ArkSIAAccessService
}

// NewArkSIAAPI creates a new instance of ArkSIAAPI with the provided ArkISPAuth.
func NewArkSIAAPI(ispAuth *auth.ArkISPAuth) (*ArkSIAAPI, error) {
	var baseIspAuth auth.ArkAuth = ispAuth
	ssoService, err := sso.NewArkSIASSOService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	k8sService, err := k8s.NewArkSIAK8SService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	targetSetsService, err := targetsets.NewArkSIAWorkspacesTargetSetsService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	workspaceDBService, err := db.NewArkSIAWorkspacesDBService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	vmSecretsService, err := vmsecrets.NewArkSIASecretsVMService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	dbSecretsService, err := dbsecrets.NewArkSIASecretsDBService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	accessService, err := access.NewArkSIAAccessService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	return &ArkSIAAPI{
		ssoService:          ssoService,
		k8sService:          k8sService,
		targetSetsService:   targetSetsService,
		workspacesDBService: workspaceDBService,
		vmSecretsService:    vmSecretsService,
		dbSecretsService:    dbSecretsService,
		accessService:       accessService,
	}, nil
}

// Sso returns the SSO service of the ArkSIAAPI instance.
func (api *ArkSIAAPI) Sso() *sso.ArkSIASSOService {
	return api.ssoService
}

// K8s returns the K8S service of the ArkSIAAPI instance.
func (api *ArkSIAAPI) K8s() *k8s.ArkSIAK8SService {
	return api.k8sService
}

// WorkspacesTargetSets returns the TargetSets service of the ArkSIAAPI instance.
func (api *ArkSIAAPI) WorkspacesTargetSets() *targetsets.ArkSIAWorkspacesTargetSetsService {
	return api.targetSetsService
}

// WorkspacesDB returns the workspace DB service of the ArkSIAAPI instance.
func (api *ArkSIAAPI) WorkspacesDB() *db.ArkSIAWorkspacesDBService {
	return api.workspacesDBService
}

// SecretsVM returns the VM Secrets service of the ArkSIAAPI instance.
func (api *ArkSIAAPI) SecretsVM() *vmsecrets.ArkSIASecretsVMService {
	return api.vmSecretsService
}

// SecretsDB returns the DB Secrets service of the ArkSIAAPI instance.
func (api *ArkSIAAPI) SecretsDB() *dbsecrets.ArkSIASecretsDBService {
	return api.dbSecretsService
}

// Access returns the access service of the ArkSIAAPI instance.
func (api *ArkSIAAPI) Access() *access.ArkSIAAccessService {
	return api.accessService
}
