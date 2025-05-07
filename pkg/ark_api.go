package api

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	"github.com/cyberark/ark-sdk-golang/pkg/services/cmgr"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity/directories"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity/roles"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity/users"
	"github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/accounts"
	"github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/safes"
	siaaccess "github.com/cyberark/ark-sdk-golang/pkg/services/sia/access"
	siak8s "github.com/cyberark/ark-sdk-golang/pkg/services/sia/k8s"
	siasecretsdb "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/db"
	siasecretsvm "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/vm"
	siasso "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso"
	siaworkspacesdb "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db"
	siatargetsets "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/targetsets"
)

// ArkAPI Wraps different API functionality of Ark Services.
type ArkAPI struct {
	authenticators []auth.ArkAuth
	services       map[string]*services.ArkService
	profile        *models.ArkProfile
}

// NewArkAPI creates a new ArkAPI instance with the provided authenticators and profile.
func NewArkAPI(authenticators []auth.ArkAuth, profile *models.ArkProfile) (*ArkAPI, error) {
	var err error
	if profile == nil {
		profile, err = (*profiles.DefaultProfilesLoader()).LoadDefaultProfile()
		if err != nil {
			return nil, err
		}
	}
	return &ArkAPI{
		authenticators: authenticators,
		services:       make(map[string]*services.ArkService),
		profile:        profile,
	}, nil
}

func (api *ArkAPI) loadServiceAuthenticators(config services.ArkServiceConfig) []auth.ArkAuth {
	var authenticators []auth.ArkAuth
	for _, authenticator := range api.authenticators {
		for _, name := range config.RequiredAuthenticatorNames {
			if authenticator.AuthenticatorName() == name {
				authenticators = append(authenticators, authenticator)
			}
		}
	}
	for _, authenticator := range api.authenticators {
		for _, name := range config.OptionalAuthenticatorNames {
			if authenticator.AuthenticatorName() == name {
				authenticators = append(authenticators, authenticator)
			}
		}
	}
	return authenticators
}

// Authenticator returns the authenticator with the specified name from the ArkAPI instance.
func (api *ArkAPI) Authenticator(authenticatorName string) (auth.ArkAuth, error) {
	for _, authenticator := range api.authenticators {
		if authenticator.AuthenticatorName() == authenticatorName {
			return authenticator, nil
		}
	}
	return nil, fmt.Errorf("%s is not supported or not found", authenticatorName)
}

// Profile returns the profile associated with the ArkAPI instance.
func (api *ArkAPI) Profile() *models.ArkProfile {
	return api.profile
}

// SiaSso returns the SiaSSO service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) SiaSso() (*siasso.ArkSIASSOService, error) {
	if ssoServiceInterface, ok := api.services[siasso.SIASSOServiceConfig.ServiceName]; ok {
		return (*ssoServiceInterface).(*siasso.ArkSIASSOService), nil
	}
	ssoService, err := siasso.NewArkSIASSOService(api.loadServiceAuthenticators(siasso.SIASSOServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var ssoBaseService services.ArkService = ssoService
	api.services[siasso.SIASSOServiceConfig.ServiceName] = &ssoBaseService
	return ssoService, nil
}

// SiaK8s returns the SiaK8s service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) SiaK8s() (*siak8s.ArkSIAK8SService, error) {
	if k8sServiceInterface, ok := api.services[siak8s.SIAK8SServiceConfig.ServiceName]; ok {
		return (*k8sServiceInterface).(*siak8s.ArkSIAK8SService), nil
	}
	k8sService, err := siak8s.NewArkSIAK8SService(api.loadServiceAuthenticators(siak8s.SIAK8SServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var k8sBaseService services.ArkService = k8sService
	api.services[siak8s.SIAK8SServiceConfig.ServiceName] = &k8sBaseService
	return k8sService, nil
}

// SiaWorkspacesTargetSets returns the SiaTargetSets service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) SiaWorkspacesTargetSets() (*siatargetsets.ArkSIATargetSetsWorkspaceService, error) {
	if targetSetsServiceInterface, ok := api.services[siatargetsets.SIATargetSetsWorkspaceServiceConfig.ServiceName]; ok {
		return (*targetSetsServiceInterface).(*siatargetsets.ArkSIATargetSetsWorkspaceService), nil
	}
	targetSetsService, err := siatargetsets.NewArkSIATargetSetsWorkspaceService(api.loadServiceAuthenticators(siatargetsets.SIATargetSetsWorkspaceServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var targetSetsBaseService services.ArkService = targetSetsService
	api.services[siatargetsets.SIATargetSetsWorkspaceServiceConfig.ServiceName] = &targetSetsBaseService
	return targetSetsService, nil
}

// SiaWorkspacesDB returns the Workspaces DB service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) SiaWorkspacesDB() (*siaworkspacesdb.ArkSIADBWorkspaceService, error) {
	if workspacesDBServiceInterface, ok := api.services[siaworkspacesdb.SIADBWorkspaceServiceConfig.ServiceName]; ok {
		return (*workspacesDBServiceInterface).(*siaworkspacesdb.ArkSIADBWorkspaceService), nil
	}
	workspacesDBService, err := siaworkspacesdb.NewArkSIADBWorkspaceService(api.loadServiceAuthenticators(siaworkspacesdb.SIADBWorkspaceServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var workspacesDBBaseService services.ArkService = workspacesDBService
	api.services[siaworkspacesdb.SIADBWorkspaceServiceConfig.ServiceName] = &workspacesDBBaseService
	return workspacesDBService, nil
}

// SiaSecretsVM returns the SiaSecretsVM service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) SiaSecretsVM() (*siasecretsvm.ArkSIASecretsVMService, error) {
	if secretsVMServiceInterface, ok := api.services[siasecretsvm.SIASecretsVMServiceConfig.ServiceName]; ok {
		return (*secretsVMServiceInterface).(*siasecretsvm.ArkSIASecretsVMService), nil
	}
	secretsVMService, err := siasecretsvm.NewArkSIASecretsVMService(api.loadServiceAuthenticators(siasecretsvm.SIASecretsVMServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var secretsVMBaseService services.ArkService = secretsVMService
	api.services[siasecretsvm.SIASecretsVMServiceConfig.ServiceName] = &secretsVMBaseService
	return secretsVMService, nil
}

// SiaSecretsDB returns the SiaSecretsDB service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) SiaSecretsDB() (*siasecretsdb.ArkSIASecretsDBService, error) {
	if secretsDBServiceInterface, ok := api.services[siasecretsdb.SIASecretsDBServiceConfig.ServiceName]; ok {
		return (*secretsDBServiceInterface).(*siasecretsdb.ArkSIASecretsDBService), nil
	}
	secretsDBService, err := siasecretsdb.NewArkSIASecretsDBService(api.loadServiceAuthenticators(siasecretsdb.SIASecretsDBServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var secretsDBBaseService services.ArkService = secretsDBService
	api.services[siasecretsdb.SIASecretsDBServiceConfig.ServiceName] = &secretsDBBaseService
	return secretsDBService, nil
}

// SiaAccess returns the SiaAccess service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) SiaAccess() (*siaaccess.ArkSIAAccessService, error) {
	if accessServiceInterface, ok := api.services[siaaccess.SIAAccessServiceConfig.ServiceName]; ok {
		return (*accessServiceInterface).(*siaaccess.ArkSIAAccessService), nil
	}
	accessService, err := siaaccess.NewArkSIAAccessService(api.loadServiceAuthenticators(siaaccess.SIAAccessServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var accessBaseService services.ArkService = accessService
	api.services[siaaccess.SIAAccessServiceConfig.ServiceName] = &accessBaseService
	return accessService, nil
}

// Cmgr returns the Cmgr service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) Cmgr() (*cmgr.ArkCmgrService, error) {
	if cmgrServiceInterface, ok := api.services[cmgr.CmgrServiceConfig.ServiceName]; ok {
		return (*cmgrServiceInterface).(*cmgr.ArkCmgrService), nil
	}
	cmgrService, err := cmgr.NewArkCmgrService(api.loadServiceAuthenticators(cmgr.CmgrServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var cmgrBaseService services.ArkService = cmgrService
	api.services[cmgr.CmgrServiceConfig.ServiceName] = &cmgrBaseService
	return cmgrService, nil
}

// PCloudSafes returns the PCloudSafes service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) PCloudSafes() (*safes.ArkPCloudSafesService, error) {
	if pcloudSafesServiceInterface, ok := api.services[safes.PCloudSafesServiceConfig.ServiceName]; ok {
		return (*pcloudSafesServiceInterface).(*safes.ArkPCloudSafesService), nil
	}
	pcloudSafesService, err := safes.NewArkPCloudSafesService(api.loadServiceAuthenticators(safes.PCloudSafesServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var pcloudSafesBaseService services.ArkService = pcloudSafesService
	api.services[safes.PCloudSafesServiceConfig.ServiceName] = &pcloudSafesBaseService
	return pcloudSafesService, nil
}

// PCloudAccounts returns the PCloudAccounts service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) PCloudAccounts() (*accounts.ArkPCloudAccountsService, error) {
	if pcloudAccountsServiceInterface, ok := api.services[accounts.PCloudAccountsServiceConfig.ServiceName]; ok {
		return (*pcloudAccountsServiceInterface).(*accounts.ArkPCloudAccountsService), nil
	}
	pcloudAccountsService, err := accounts.NewArkPCloudAccountsService(api.loadServiceAuthenticators(accounts.PCloudAccountsServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var pcloudAccountsBaseService services.ArkService = pcloudAccountsService
	api.services[accounts.PCloudAccountsServiceConfig.ServiceName] = &pcloudAccountsBaseService
	return pcloudAccountsService, nil
}

// IdentityDirectories returns the IdentityDirectories service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) IdentityDirectories() (*directories.ArkIdentityDirectoriesService, error) {
	if directoriesServiceInterface, ok := api.services[directories.IdentityDirectoriesServiceConfig.ServiceName]; ok {
		return (*directoriesServiceInterface).(*directories.ArkIdentityDirectoriesService), nil
	}
	directoriesService, err := directories.NewArkIdentityDirectoriesService(api.loadServiceAuthenticators(directories.IdentityDirectoriesServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var directoriesBaseService services.ArkService = directoriesService
	api.services[directories.IdentityDirectoriesServiceConfig.ServiceName] = &directoriesBaseService
	return directoriesService, nil
}

// IdentityRoles returns the IdentityRoles service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) IdentityRoles() (*roles.ArkIdentityRolesService, error) {
	if rolesServiceInterface, ok := api.services[roles.IdentityRolesServiceConfig.ServiceName]; ok {
		return (*rolesServiceInterface).(*roles.ArkIdentityRolesService), nil
	}
	rolesService, err := roles.NewArkIdentityRolesService(api.loadServiceAuthenticators(roles.IdentityRolesServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var rolesBaseService services.ArkService = rolesService
	api.services[roles.IdentityRolesServiceConfig.ServiceName] = &rolesBaseService
	return rolesService, nil
}

// IdentityUsers returns the IdentityUsers service from the ArkAPI instance. If the service is not already created, it creates a new one.
func (api *ArkAPI) IdentityUsers() (*users.ArkIdentityUsersService, error) {
	if usersServiceInterface, ok := api.services[users.IdentityUsersServiceConfig.ServiceName]; ok {
		return (*usersServiceInterface).(*users.ArkIdentityUsersService), nil
	}
	usersService, err := users.NewArkIdentityUsersService(api.loadServiceAuthenticators(users.IdentityUsersServiceConfig)...)
	if err != nil {
		return nil, err
	}
	var usersBaseService services.ArkService = usersService
	api.services[users.IdentityUsersServiceConfig.ServiceName] = &usersBaseService
	return usersService, nil
}
