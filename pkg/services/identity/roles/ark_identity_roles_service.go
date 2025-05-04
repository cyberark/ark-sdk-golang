package roles

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/models/common/identity"
	directoriesmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/identity/directories"
	rolesmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/identity/roles"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity/directories"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"strings"
)

const (
	addUserToRoleURL         = "SaasManage/AddUsersAndGroupsToRole"
	createRoleURL            = "Roles/StoreRole"
	updateRoleURL            = "Roles/UpdateRole"
	roleMembersURL           = "Roles/GetRoleMembers"
	addAdminRightsToRoleURL  = "SaasManage/AssignSuperRights"
	removeUserFromRoleURL    = "SaasManage/RemoveUsersAndGroupsFromRole"
	deleteRoleURL            = "SaasManage/DeleteRole"
	directoryServiceQueryURL = "UserMgmt/DirectoryServiceQuery"
)

// IdentityRolesServiceConfig is the configuration for the identity roles service.
var IdentityRolesServiceConfig = services.ArkServiceConfig{
	ServiceName:                "identity-roles",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkIdentityRolesService is the service for managing identity roles.
type ArkIdentityRolesService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkIdentityRolesService creates a new instance of ArkIdentityRolesService.
func NewArkIdentityRolesService(authenticators ...auth.ArkAuth) (*ArkIdentityRolesService, error) {
	identityRolesService := &ArkIdentityRolesService{}
	var identityRolesServiceInterface services.ArkService = identityRolesService
	baseService, err := services.NewArkBaseService(identityRolesServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "", "", "api/idadmin", identityRolesService.refreshIdentityRolesAuth)
	if err != nil {
		return nil, err
	}
	client.UpdateHeaders(map[string]string{
		"X-IDAP-NATIVE-CLIENT": "true",
	})
	identityRolesService.client = client
	identityRolesService.ispAuth = ispAuth
	identityRolesService.ArkBaseService = baseService
	return identityRolesService, nil
}

func (s *ArkIdentityRolesService) refreshIdentityRolesAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// CreateRole creates a new role in the identity service.
func (s *ArkIdentityRolesService) CreateRole(createRole *rolesmodels.ArkIdentityCreateRole) (*rolesmodels.ArkIdentityRole, error) {
	s.Logger.Info("Trying to create role [%s]", createRole.RoleName)
	roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{
		RoleName: createRole.RoleName,
	})
	if err == nil && roleID != "" {
		s.Logger.Info("Role already exists with id [%s]", roleID)
		return &rolesmodels.ArkIdentityRole{
			RoleID:   roleID,
			RoleName: createRole.RoleName,
		}, nil
	}
	createRoleRequest := map[string]interface{}{
		"Name": createRole.RoleName,
	}
	if createRole.Description != "" {
		createRoleRequest["Description"] = createRole.Description
	}
	response, err := s.client.Post(context.Background(), createRoleURL, createRoleRequest)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result["success"].(bool) == false {
		return nil, fmt.Errorf("failed to create role - [%v]", result)
	}
	roleID = result["Result"].(map[string]interface{})["_RowKey"].(string)
	roleDetails := &rolesmodels.ArkIdentityRole{
		RoleName: createRole.RoleName,
		RoleID:   roleID,
	}
	s.Logger.Info(fmt.Sprintf("Role created with id [%s]", roleID))
	if len(createRole.AdminRights) > 0 {
		err = s.AddAdminRightsToRole(&rolesmodels.ArkIdentityAddAdminRightsToRole{
			RoleID:      roleDetails.RoleID,
			AdminRights: createRole.AdminRights,
		})
		if err != nil {
			return nil, err
		}
	}
	return roleDetails, nil
}

// UpdateRole updates an existing role in the identity service.
func (s *ArkIdentityRolesService) UpdateRole(updateRole *rolesmodels.ArkIdentityUpdateRole) error {
	if updateRole.RoleName != "" && updateRole.RoleID == "" {
		roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: updateRole.RoleName})
		if err != nil {
			return fmt.Errorf("failed to retrieve role ID by name: %v", err)
		}
		updateRole.RoleID = roleID
	}
	s.Logger.Info(fmt.Sprintf("Updating identity role [%s]", updateRole.RoleID))
	updateDict := map[string]interface{}{
		"Name": updateRole.RoleID,
	}
	if updateRole.NewRoleName != "" {
		updateDict["NewName"] = updateRole.NewRoleName
	}
	if updateRole.Description != "" {
		updateDict["Description"] = updateRole.Description
	}
	response, err := s.client.Post(context.Background(), updateRoleURL, updateDict)
	if err != nil {
		return fmt.Errorf("failed to update role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if result["success"].(bool) == false {
		return fmt.Errorf("failed to update role - [%v]", result)
	}
	s.Logger.Info("Role updated successfully")
	return nil
}

// ListRoleMembers retrieves the members of a role in the identity service.
func (s *ArkIdentityRolesService) ListRoleMembers(listRoleMembers *rolesmodels.ArkIdentityListRoleMembers) ([]*rolesmodels.ArkIdentityRoleMember, error) {
	if listRoleMembers.RoleName != "" && listRoleMembers.RoleID == "" {
		roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: listRoleMembers.RoleName})
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve role ID by name: %v", err)
		}
		listRoleMembers.RoleID = roleID
	}
	s.Logger.Info(fmt.Sprintf("Listing identity role [%s] members", listRoleMembers.RoleID))
	requestBody := map[string]interface{}{
		"Name": listRoleMembers.RoleID,
	}
	response, err := s.client.Post(context.Background(), roleMembersURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to list role members: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list role members - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if !result["success"].(bool) {
		return nil, fmt.Errorf("failed to list role members - [%v]", result)
	}
	var members []*rolesmodels.ArkIdentityRoleMember
	if resultMap, ok := result["Result"].(map[string]interface{}); ok {
		if results, ok := resultMap["Results"].([]interface{}); ok && len(results) > 0 {
			for _, r := range results {
				row := r.(map[string]interface{})["Row"].(map[string]interface{})
				members = append(members, &rolesmodels.ArkIdentityRoleMember{
					MemberID:   row["Guid"].(string),
					MemberName: row["Name"].(string),
					MemberType: strings.ToUpper(row["Type"].(string)),
				})
			}
		}
	}
	s.Logger.Info("Listed role members successfully")
	return members, nil
}

// AddAdminRightsToRole adds admin rights to a role in the identity service.
func (s *ArkIdentityRolesService) AddAdminRightsToRole(addAdminRightsToRole *rolesmodels.ArkIdentityAddAdminRightsToRole) error {
	s.Logger.Info(fmt.Sprintf("Adding admin rights [%v] to role [%s]", addAdminRightsToRole.AdminRights, addAdminRightsToRole.RoleName))

	if addAdminRightsToRole.RoleID == "" && addAdminRightsToRole.RoleName == "" {
		return fmt.Errorf("either role ID or role name must be given")
	}
	var roleID string
	if addAdminRightsToRole.RoleID != "" {
		roleID = addAdminRightsToRole.RoleID
	} else {
		var err error
		roleID, err = s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: addAdminRightsToRole.RoleName})
		if err != nil {
			return fmt.Errorf("failed to retrieve role ID by name: %v", err)
		}
	}
	requestBody := make([]map[string]interface{}, len(addAdminRightsToRole.AdminRights))
	for i, adminRight := range addAdminRightsToRole.AdminRights {
		requestBody[i] = map[string]interface{}{
			"Role": roleID,
			"Path": adminRight,
		}
	}
	response, err := s.client.Post(context.Background(), addAdminRightsToRoleURL, requestBody)
	if err != nil {
		return fmt.Errorf("failed to add admin rights to role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add admin rights to role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	if !result["success"].(bool) {
		return fmt.Errorf("failed to add admin rights to role - [%v]", result)
	}
	s.Logger.Info("Admin rights added to role successfully")
	return nil
}

// RoleIDByName retrieves the role ID by its name.
func (s *ArkIdentityRolesService) RoleIDByName(roleIDByName *rolesmodels.ArkIdentityRoleIDByName) (string, error) {
	s.Logger.Info(fmt.Sprintf("Retrieving role ID for name [%s]", roleIDByName.RoleName))
	directoriesService, err := directories.NewArkIdentityDirectoriesService(s.ispAuth)
	if err != nil {
		return "", fmt.Errorf("failed to initialize directories service: %v", err)
	}
	foundDirectories, err := directoriesService.ListDirectories(&directoriesmodels.ArkIdentityListDirectories{
		Directories: []string{identity.Identity},
	})
	if err != nil {
		return "", fmt.Errorf("failed to list directories: %v", err)
	}
	var directoryUUIDs []string
	for _, d := range foundDirectories {
		directoryUUIDs = append(directoryUUIDs, d.DirectoryServiceUUID)
	}
	specificRoleRequest := identity.NewDirectoryServiceQuerySpecificRoleRequest(roleIDByName.RoleName)
	specificRoleRequest.DirectoryServices = directoryUUIDs
	specificRoleRequest.Args = identity.DirectorySearchArgs{Limit: 1}
	var specificRoleRequestBody map[string]interface{}
	err = mapstructure.Decode(specificRoleRequest, &specificRoleRequestBody)
	if err != nil {
		return "", fmt.Errorf("failed to decode specific role request: %v", err)
	}
	response, err := s.client.Post(context.Background(), directoryServiceQueryURL, specificRoleRequestBody)
	if err != nil {
		return "", fmt.Errorf("failed to query directory services role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to query for directory services role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	if !result["success"].(bool) {
		return "", fmt.Errorf("failed to query for directory services role - [%v]", result)
	}
	var queryResponse identity.DirectoryServiceQueryResponse
	err = mapstructure.Decode(result, &queryResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}
	allRoles := queryResponse.Result.Roles.Results
	if len(allRoles) == 0 {
		return "", fmt.Errorf("no role found for given name")
	}
	return allRoles[0].Row.ID, nil
}

// AddUserToRole adds a user to a role in the identity service.
func (s *ArkIdentityRolesService) AddUserToRole(addUserToRole *rolesmodels.ArkIdentityAddUserToRole) error {
	s.Logger.Info(fmt.Sprintf("Adding user [%s] to role [%s]", addUserToRole.Username, addUserToRole.RoleName))
	roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: addUserToRole.RoleName})
	if err != nil {
		return fmt.Errorf("failed to retrieve role ID by name: %v", err)
	}
	requestBody := map[string]interface{}{
		"Name":  roleID,
		"Users": []string{addUserToRole.Username},
	}
	fmt.Printf("Request body: %v\n", requestBody)
	response, err := s.client.Post(context.Background(), addUserToRoleURL, requestBody)
	if err != nil {
		return fmt.Errorf("failed to add user to role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add user to role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	if !result["success"].(bool) {
		return fmt.Errorf("failed to add user to role - [%v]", result)
	}
	s.Logger.Info("User added to role successfully")
	return nil
}

// AddGroupToRole adds a group to a role in the identity service.
func (s *ArkIdentityRolesService) AddGroupToRole(addGroupToRole *rolesmodels.ArkIdentityAddGroupToRole) error {
	s.Logger.Info(fmt.Sprintf("Adding group [%s] to role [%s]", addGroupToRole.GroupName, addGroupToRole.RoleName))
	roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: addGroupToRole.RoleName})
	if err != nil {
		return fmt.Errorf("failed to retrieve role ID by name: %v", err)
	}
	requestBody := map[string]interface{}{
		"Name":   roleID,
		"Groups": []string{addGroupToRole.GroupName},
	}
	response, err := s.client.Post(context.Background(), addUserToRoleURL, requestBody)
	if err != nil {
		return fmt.Errorf("failed to add group to role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add group to role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	if !result["success"].(bool) {
		return fmt.Errorf("failed to add group to role - [%v]", result)
	}
	s.Logger.Info("Group added to role successfully")
	return nil
}

// AddRoleToRole adds a role to another role in the identity service.
func (s *ArkIdentityRolesService) AddRoleToRole(addRoleToRole *rolesmodels.ArkIdentityAddRoleToRole) error {
	s.Logger.Info(fmt.Sprintf("Adding role [%s] to role [%s]", addRoleToRole.RoleNameToAdd, addRoleToRole.RoleName))
	roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: addRoleToRole.RoleName})
	if err != nil {
		return fmt.Errorf("failed to retrieve role ID by name: %v", err)
	}
	requestBody := map[string]interface{}{
		"Name":  roleID,
		"Roles": []string{addRoleToRole.RoleNameToAdd},
	}
	response, err := s.client.Post(context.Background(), addUserToRoleURL, requestBody)
	if err != nil {
		return fmt.Errorf("failed to add role to role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add role to role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	if !result["success"].(bool) {
		return fmt.Errorf("failed to add role to role - [%v]", result)
	}
	s.Logger.Info("Role added to role successfully")
	return nil
}

// RemoveUserFromRole removes a user from a role in the identity service.
func (s *ArkIdentityRolesService) RemoveUserFromRole(removeUserFromRole *rolesmodels.ArkIdentityRemoveUserFromRole) error {
	s.Logger.Info(fmt.Sprintf("Removing user [%s] from role [%s]", removeUserFromRole.Username, removeUserFromRole.RoleName))
	roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: removeUserFromRole.RoleName})
	if err != nil {
		return fmt.Errorf("failed to retrieve role ID by name: %v", err)
	}
	requestBody := map[string]interface{}{
		"Name":  roleID,
		"Users": []string{removeUserFromRole.Username},
	}
	response, err := s.client.Post(context.Background(), removeUserFromRoleURL, requestBody)
	if err != nil {
		return fmt.Errorf("failed to remove user from role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to remove user from role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	if !result["success"].(bool) {
		return fmt.Errorf("failed to remove user from role - [%v]", result)
	}
	s.Logger.Info("User removed from role successfully")
	return nil
}

// RemoveGroupFromRole removes a group from a role in the identity service.
func (s *ArkIdentityRolesService) RemoveGroupFromRole(removeGroupFromRole *rolesmodels.ArkIdentityRemoveGroupFromRole) error {
	s.Logger.Info(fmt.Sprintf("Removing group [%s] from role [%s]", removeGroupFromRole.GroupName, removeGroupFromRole.RoleName))
	roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: removeGroupFromRole.RoleName})
	if err != nil {
		return fmt.Errorf("failed to retrieve role ID by name: %v", err)
	}
	requestBody := map[string]interface{}{
		"Name":   roleID,
		"Groups": []string{removeGroupFromRole.GroupName},
	}
	response, err := s.client.Post(context.Background(), removeUserFromRoleURL, requestBody)
	if err != nil {
		return fmt.Errorf("failed to remove group from role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to remove group from role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	if !result["success"].(bool) {
		return fmt.Errorf("failed to remove group from role - [%v]", result)
	}
	s.Logger.Info("Group removed from role successfully")
	return nil
}

// RemoveRoleFromRole removes a role from another role in the identity service.
func (s *ArkIdentityRolesService) RemoveRoleFromRole(removeRoleFromRole *rolesmodels.ArkIdentityRemoveRoleFromRole) error {
	s.Logger.Info(fmt.Sprintf("Removing role [%s] from role [%s]", removeRoleFromRole.RoleNameToRemove, removeRoleFromRole.RoleName))
	roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: removeRoleFromRole.RoleName})
	if err != nil {
		return fmt.Errorf("failed to retrieve role ID by name: %v", err)
	}
	requestBody := map[string]interface{}{
		"Name":  roleID,
		"Roles": []string{removeRoleFromRole.RoleNameToRemove},
	}
	response, err := s.client.Post(context.Background(), removeUserFromRoleURL, requestBody)
	if err != nil {
		return fmt.Errorf("failed to remove role from role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to remove role from role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	if !result["success"].(bool) {
		return fmt.Errorf("failed to remove role from role - [%v]", result)
	}
	s.Logger.Info("Role removed from role successfully")
	return nil
}

// DeleteRole deletes a role in the identity service.
func (s *ArkIdentityRolesService) DeleteRole(deleteRole *rolesmodels.ArkIdentityDeleteRole) error {
	s.Logger.Info(fmt.Sprintf("Deleting role [%s]", deleteRole.RoleName))
	if deleteRole.RoleName != "" && deleteRole.RoleID == "" {
		roleID, err := s.RoleIDByName(&rolesmodels.ArkIdentityRoleIDByName{RoleName: deleteRole.RoleName})
		if err != nil {
			return fmt.Errorf("failed to retrieve role ID by name: %v", err)
		}
		deleteRole.RoleID = roleID
	}
	requestBody := map[string]interface{}{
		"Name": deleteRole.RoleID,
	}
	response, err := s.client.Post(context.Background(), deleteRoleURL, requestBody)
	if err != nil {
		return fmt.Errorf("failed to delete role: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete role - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	if !result["success"].(bool) {
		return fmt.Errorf("failed to delete role - [%v]", result)
	}
	s.Logger.Info("Role deleted successfully")
	return nil
}

// ServiceConfig returns the service configuration for the ArkIdentityRolesService.
func (s *ArkIdentityRolesService) ServiceConfig() services.ArkServiceConfig {
	return IdentityRolesServiceConfig
}
