package safes

import (
	"context"
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	safesmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/pcloud/safes"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"
)

// Constants for safes URLs
const (
	safesURL       = "/api/safes"
	safeURL        = "/api/safes/%s/"
	safeMembersURL = "/api/safes/%s/members"
	safeMemberURL  = "/api/safes/%s/members/%s/"
)

// SafeMembersPermissionsSets maps permission sets to their corresponding permissions
var SafeMembersPermissionsSets = map[string]safesmodels.ArkPCloudSafeMemberPermissions{
	safesmodels.ConnectOnly: {
		ListAccounts: true,
		UseAccounts:  true,
	},
	safesmodels.ReadOnly: {
		ListAccounts:     true,
		UseAccounts:      true,
		RetrieveAccounts: true,
	},
	safesmodels.Approver: {
		ListAccounts:                true,
		ViewSafeMembers:             true,
		ManageSafeMembers:           true,
		RequestsAuthorizationLevel1: true,
	},
	safesmodels.AccountsManager: {
		ListAccounts:                           true,
		UseAccounts:                            true,
		RetrieveAccounts:                       true,
		AddAccounts:                            true,
		UpdateAccountProperties:                true,
		UpdateAccountContent:                   true,
		InitiateCPMAccountManagementOperations: true,
		SpecifyNextAccountContent:              true,
		RenameAccounts:                         true,
		DeleteAccounts:                         true,
		UnlockAccounts:                         true,
		ViewSafeMembers:                        true,
		ManageSafeMembers:                      true,
		ViewAuditLog:                           true,
		AccessWithoutConfirmation:              true,
	},
	safesmodels.Full: {
		ListAccounts:                           true,
		UseAccounts:                            true,
		RetrieveAccounts:                       true,
		AddAccounts:                            true,
		UpdateAccountProperties:                true,
		UpdateAccountContent:                   true,
		InitiateCPMAccountManagementOperations: true,
		SpecifyNextAccountContent:              true,
		RenameAccounts:                         true,
		DeleteAccounts:                         true,
		UnlockAccounts:                         true,
		ViewSafeMembers:                        true,
		ManageSafeMembers:                      true,
		ViewAuditLog:                           true,
		AccessWithoutConfirmation:              true,
		RequestsAuthorizationLevel1:            true,
		ManageSafe:                             true,
		BackupSafe:                             true,
		MoveAccountsAndFolders:                 true,
		CreateFolders:                          true,
		DeleteFolders:                          true,
	},
}

// ArkPCloudSafesPage is a page of ArkPCloudSafe items.
type ArkPCloudSafesPage = common.ArkPage[safesmodels.ArkPCloudSafe]

// ArkPCloudSafeMembersPage is a page of ArkPCloudSafeMember items.
type ArkPCloudSafeMembersPage = common.ArkPage[safesmodels.ArkPCloudSafeMember]

// PCloudSafesServiceConfig is the configuration for the SIA pCloud Safes service.
var PCloudSafesServiceConfig = services.ArkServiceConfig{
	ServiceName:                "pcloud-safes",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkPCloudSafesService is the service for managing pCloud Safes.
type ArkPCloudSafesService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkPCloudSafesService creates a new instance of ArkPCloudSafesService.
func NewArkPCloudSafesService(authenticators ...auth.ArkAuth) (*ArkPCloudSafesService, error) {
	pcloudSafesService := &ArkPCloudSafesService{}
	var pcloudSafesServiceInterface services.ArkService = pcloudSafesService
	baseService, err := services.NewArkBaseService(pcloudSafesServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "privilegecloud", ".", "passwordvault", pcloudSafesService.refreshPCloudSafesAuth)
	if err != nil {
		return nil, err
	}
	pcloudSafesService.client = client
	pcloudSafesService.ispAuth = ispAuth
	pcloudSafesService.ArkBaseService = baseService
	return pcloudSafesService, nil
}

func (s *ArkPCloudSafesService) refreshPCloudSafesAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArkPCloudSafesService) listSafesWithFilters(
	search string,
	sort string,
	offset int,
	limit int,
) (<-chan *ArkPCloudSafesPage, error) {
	query := map[string]string{}
	if search != "" {
		query["search"] = search
	}
	if sort != "" {
		query["sort"] = sort
	}
	if offset > 0 {
		query["offset"] = fmt.Sprintf("%d", offset)
	}
	if limit > 0 {
		query["limit"] = fmt.Sprintf("%d", limit)
	}
	results := make(chan *ArkPCloudSafesPage)
	go func() {
		defer close(results)
		for {
			response, err := s.client.Get(context.Background(), safesURL, query)
			if err != nil {
				s.Logger.Error("Failed to list safes: %v", err)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					common.GlobalLogger.Warning("Error closing response body")
				}
			}(response.Body)
			if response.StatusCode != http.StatusOK {
				s.Logger.Error("Failed to list safes - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
				return
			}
			result, err := common.DeserializeJSONSnake(response.Body)
			if err != nil {
				s.Logger.Error("Failed to decode response: %v", err)
				return
			}
			resultMap := result.(map[string]interface{})
			var safesJSON []interface{}
			if value, ok := resultMap["value"]; ok {
				safesJSON = value.([]interface{})
			} else if safesData, ok := resultMap["Safes"]; ok {
				safesJSON = safesData.([]interface{})
			} else {
				s.Logger.Error("Failed to list safes, unexpected result")
				return
			}
			var safes []*safesmodels.ArkPCloudSafe
			if err := mapstructure.Decode(safesJSON, &safes); err != nil {
				s.Logger.Error("Failed to validate safes: %v", err)
				return
			}
			results <- &ArkPCloudSafesPage{Items: safes}
			if nextLink, ok := resultMap["nextLink"].(string); ok {
				nextQuery, _ := url.Parse(nextLink)
				queryValues := nextQuery.Query()
				query = make(map[string]string)
				for key, values := range queryValues {
					if len(values) > 0 {
						query[key] = values[0]
					}
				}
			} else {
				break
			}
		}
	}()
	return results, nil
}

func (s *ArkPCloudSafesService) listSafeMembersWithFilters(
	safeID string,
	search string,
	sort string,
	offset int,
	limit int,
	memberType string,
) (<-chan *ArkPCloudSafeMembersPage, error) {
	query := map[string]string{}
	if search != "" {
		query["search"] = search
	}
	if sort != "" {
		query["sort"] = sort
	}
	if offset > 0 {
		query["offset"] = fmt.Sprintf("%d", offset)
	}
	if limit > 0 {
		query["limit"] = fmt.Sprintf("%d", limit)
	}
	if memberType != "" {
		query["filter"] = fmt.Sprintf("memberType eq %s", memberType)
	}
	results := make(chan *ArkPCloudSafeMembersPage)
	go func() {
		defer close(results)
		for {
			response, err := s.client.Get(context.Background(), fmt.Sprintf(safeMembersURL, safeID), query)
			if err != nil {
				s.Logger.Error("Failed to list safe members: %v", err)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					common.GlobalLogger.Warning("Error closing response body")
				}
			}(response.Body)
			if response.StatusCode != http.StatusOK {
				s.Logger.Error("Failed to list safe members - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
				return
			}
			result, err := common.DeserializeJSONSnake(response.Body)
			if err != nil {
				s.Logger.Error("Failed to decode response: %v", err)
				return
			}
			resultMap := result.(map[string]interface{})
			var membersJSON []interface{}
			if value, ok := resultMap["value"]; ok {
				membersJSON = value.([]interface{})
			} else {
				s.Logger.Error("Failed to list safe members, unexpected result")
				return
			}
			var members []*safesmodels.ArkPCloudSafeMember
			if err := mapstructure.Decode(membersJSON, &members); err != nil {
				s.Logger.Error("Failed to validate safe members: %v", err)
				return
			}
			for _, member := range members {
				member.PermissionSet = safesmodels.Custom
				for permissionSet, permissions := range SafeMembersPermissionsSets {
					if reflect.DeepEqual(member.Permissions, permissions) {
						member.PermissionSet = permissionSet
						break
					}
				}
			}
			results <- &ArkPCloudSafeMembersPage{Items: members}
			if nextLink, ok := resultMap["nextLink"].(string); ok {
				nextQuery, _ := url.Parse(nextLink)
				queryValues := nextQuery.Query()
				query = make(map[string]string)
				for key, values := range queryValues {
					if len(values) > 0 {
						query[key] = values[0]
					}
				}
			} else {
				break
			}
		}
	}()
	return results, nil
}

// ListSafes returns a channel of ArkPCloudSafesPage containing all safes.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/SDK/Safes%20Web%20Services%20-%20List%20Safes.htm?
func (s *ArkPCloudSafesService) ListSafes() (<-chan *ArkPCloudSafesPage, error) {
	return s.listSafesWithFilters(
		"",
		"",
		0,
		0,
	)
}

// ListSafesBy returns a channel of ArkPCloudSafesPage containing safes filtered by the given filters.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/SDK/Safes%20Web%20Services%20-%20List%20Safes.htm?
func (s *ArkPCloudSafesService) ListSafesBy(safesFilters *safesmodels.ArkPCloudSafesFilters) (<-chan *ArkPCloudSafesPage, error) {
	return s.listSafesWithFilters(
		safesFilters.Search,
		safesFilters.Sort,
		safesFilters.Offset,
		safesFilters.Limit,
	)
}

// ListSafeMembers returns a channel of ArkPCloudSafeMembersPage containing all safe members.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/SDK/Safe%20Members%20WS%20-%20List%20Safe%20Members.htm
func (s *ArkPCloudSafesService) ListSafeMembers(listSafeMembers *safesmodels.ArkPCloudListSafeMembers) (<-chan *ArkPCloudSafeMembersPage, error) {
	return s.listSafeMembersWithFilters(
		listSafeMembers.SafeID,
		"",
		"",
		0,
		0,
		"",
	)
}

// ListSafeMembersBy returns a channel of ArkPCloudSafeMembersPage containing safe members filtered by the given filters.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/SDK/Safe%20Members%20WS%20-%20List%20Safe%20Members.htm
func (s *ArkPCloudSafesService) ListSafeMembersBy(safeMembersFilters *safesmodels.ArkPCloudSafeMembersFilters) (<-chan *ArkPCloudSafeMembersPage, error) {
	return s.listSafeMembersWithFilters(
		safeMembersFilters.SafeID,
		safeMembersFilters.Search,
		safeMembersFilters.Sort,
		safeMembersFilters.Offset,
		safeMembersFilters.Limit,
		safeMembersFilters.MemberType,
	)
}

// Safe retrieves a safe by its ID.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/SDK/Safes%20Web%20Services%20-%20Get%20Safes%20Details.htm
func (s *ArkPCloudSafesService) Safe(getSafe *safesmodels.ArkPCloudGetSafe) (*safesmodels.ArkPCloudSafe, error) {
	s.Logger.Info("Retrieving safe [%s]", getSafe.SafeID)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(safeURL, getSafe.SafeID), nil)
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
		return nil, fmt.Errorf("failed to retrieve safe - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	safeJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var safe safesmodels.ArkPCloudSafe
	err = mapstructure.Decode(safeJSON, &safe)
	if err != nil {
		return nil, err
	}
	return &safe, nil
}

// SafeMember retrieves a safe member by its safe ID and member name.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/SDK/Safe%20Members%20WS%20-%20List%20Safe%20Member.htm
func (s *ArkPCloudSafesService) SafeMember(getSafeMember *safesmodels.ArkPCloudGetSafeMember) (*safesmodels.ArkPCloudSafeMember, error) {
	s.Logger.Info("Retrieving safe member [%s] [%s]", getSafeMember.SafeID, getSafeMember.MemberName)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(safeMemberURL, getSafeMember.SafeID, getSafeMember.MemberName), nil)
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
		return nil, fmt.Errorf("failed to retrieve safe member - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	safeMemberJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var safeMember safesmodels.ArkPCloudSafeMember
	err = mapstructure.Decode(safeMemberJSON, &safeMember)
	if err != nil {
		return nil, err
	}
	safeMember.PermissionSet = safesmodels.Custom
	for permissionSet, permissions := range SafeMembersPermissionsSets {
		if reflect.DeepEqual(safeMember.Permissions, permissions) {
			safeMember.PermissionSet = permissionSet
			break
		}
	}
	return &safeMember, nil
}

// AddSafe adds a new safe.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/WebServices/Add%20Safe.htm
func (s *ArkPCloudSafesService) AddSafe(addSafe *safesmodels.ArkPCloudAddSafe) (*safesmodels.ArkPCloudSafe, error) {
	s.Logger.Info("Adding safe [%s]", addSafe.SafeName)
	addSafeJSON, err := common.SerializeJSONCamel(addSafe)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Post(context.Background(), safesURL, addSafeJSON)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to add safe - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	safeJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var safe safesmodels.ArkPCloudSafe
	err = mapstructure.Decode(safeJSON, &safe)
	if err != nil {
		return nil, err
	}
	return &safe, nil
}

// AddSafeMember adds a new member to a safe.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/WebServices/Add%20Safe%20Member.htm
func (s *ArkPCloudSafesService) AddSafeMember(addSafeMember *safesmodels.ArkPCloudAddSafeMember) (*safesmodels.ArkPCloudSafeMember, error) {
	s.Logger.Info("Adding safe member [%s] [%s]", addSafeMember.SafeID, addSafeMember.MemberName)
	if addSafeMember.PermissionSet == safesmodels.Custom && addSafeMember.Permissions == nil {
		return nil, fmt.Errorf("permission set is custom but permissions are not set")
	}
	if addSafeMember.PermissionSet != safesmodels.Custom {
		if permissions, ok := SafeMembersPermissionsSets[addSafeMember.PermissionSet]; ok {
			addSafeMember.Permissions = &permissions
		} else {
			return nil, fmt.Errorf("invalid permission set: %s", addSafeMember.PermissionSet)
		}
	}
	addSafeMemberJSON, err := common.SerializeJSONCamel(addSafeMember)
	if err != nil {
		return nil, err
	}
	delete(addSafeMemberJSON, "permissionSet")
	delete(addSafeMemberJSON, "safeId")
	response, err := s.client.Post(context.Background(), fmt.Sprintf(safeMembersURL, addSafeMember.SafeID), addSafeMemberJSON)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to add safe member - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	safeMemberJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var safeMember safesmodels.ArkPCloudSafeMember
	err = mapstructure.Decode(safeMemberJSON, &safeMember)
	if err != nil {
		return nil, err
	}
	safeMember.PermissionSet = safesmodels.Custom
	for permissionSet, permissions := range SafeMembersPermissionsSets {
		if reflect.DeepEqual(safeMember.Permissions, permissions) {
			safeMember.PermissionSet = permissionSet
			break
		}
	}
	return &safeMember, nil
}

// DeleteSafe deletes a safe by its ID.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/WebServices/Delete%20Safe.htm
func (s *ArkPCloudSafesService) DeleteSafe(deleteSafe *safesmodels.ArkPCloudDeleteSafe) error {
	s.Logger.Info("Deleting safe [%s]", deleteSafe.SafeID)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(safeURL, deleteSafe.SafeID), nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete safe - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// DeleteSafeMember deletes a member from a safe by its safe ID and member name.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/WebServices/Delete%20Safe%20Member.htm
func (s *ArkPCloudSafesService) DeleteSafeMember(deleteSafeMember *safesmodels.ArkPCloudDeleteSafeMember) error {
	s.Logger.Info("Deleting safe member [%s] [%s]", deleteSafeMember.SafeID, deleteSafeMember.MemberName)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(safeMemberURL, deleteSafeMember.SafeID, deleteSafeMember.MemberName), nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete safe member - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// UpdateSafe updates a safe by its ID.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/WebServices/Update%20Safe.htm
func (s *ArkPCloudSafesService) UpdateSafe(updateSafe *safesmodels.ArkPCloudUpdateSafe) (*safesmodels.ArkPCloudSafe, error) {
	s.Logger.Info("Updating safe [%s]", updateSafe.SafeID)
	updateSafeJSON, err := common.SerializeJSONCamel(updateSafe)
	if err != nil {
		return nil, err
	}
	delete(updateSafeJSON, "safeId")
	if len(updateSafeJSON) == 0 {
		return s.Safe(&safesmodels.ArkPCloudGetSafe{SafeID: updateSafe.SafeID})
	}
	response, err := s.client.Put(context.Background(), fmt.Sprintf(safeURL, updateSafe.SafeID), updateSafeJSON)
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
		return nil, fmt.Errorf("failed to update safe - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	safeJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var safe safesmodels.ArkPCloudSafe
	err = mapstructure.Decode(safeJSON, &safe)
	if err != nil {
		return nil, err
	}
	return &safe, nil
}

// UpdateSafeMember updates a member of a safe by its safe ID and member name.
// https://docs.cyberark.com/Product-Doc/OnlineHelp/PAS/Latest/en/Content/WebServices/Update%20Safe%20Member.htm
func (s *ArkPCloudSafesService) UpdateSafeMember(updateSafeMember *safesmodels.ArkPCloudUpdateSafeMember) (*safesmodels.ArkPCloudSafeMember, error) {
	s.Logger.Info("Updating safe member [%s] [%s]", updateSafeMember.SafeID, updateSafeMember.MemberName)
	if updateSafeMember.PermissionSet == safesmodels.Custom && updateSafeMember.Permissions == nil {
		return nil, fmt.Errorf("permission set is custom but permissions are not set")
	}
	if updateSafeMember.PermissionSet != safesmodels.Custom {
		if permissions, ok := SafeMembersPermissionsSets[updateSafeMember.PermissionSet]; ok {
			updateSafeMember.Permissions = &permissions
		} else {
			return nil, fmt.Errorf("invalid permission set: %s", updateSafeMember.PermissionSet)
		}
	}
	updateSafeMemberJSON, err := common.SerializeJSONCamel(updateSafeMember)
	if err != nil {
		return nil, err
	}
	delete(updateSafeMemberJSON, "safeId")
	delete(updateSafeMemberJSON, "memberName")
	delete(updateSafeMemberJSON, "permissionSet")
	if len(updateSafeMemberJSON) == 0 {
		return s.SafeMember(&safesmodels.ArkPCloudGetSafeMember{SafeID: updateSafeMember.SafeID, MemberName: updateSafeMember.MemberName})
	}
	response, err := s.client.Put(context.Background(), fmt.Sprintf(safeMemberURL, updateSafeMember.SafeID, updateSafeMember.MemberName), updateSafeMemberJSON)
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
		return nil, fmt.Errorf("failed to update safe member - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	safeMemberJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var safeMember safesmodels.ArkPCloudSafeMember
	err = mapstructure.Decode(safeMemberJSON, &safeMember)
	if err != nil {
		return nil, err
	}
	safeMember.PermissionSet = safesmodels.Custom
	for permissionSet, permissions := range SafeMembersPermissionsSets {
		if reflect.DeepEqual(safeMember.Permissions, permissions) {
			safeMember.PermissionSet = permissionSet
			break
		}
	}
	return &safeMember, nil
}

// SafesStats retrieves statistics about safes.
func (s *ArkPCloudSafesService) SafesStats() (*safesmodels.ArkPCloudSafesStats, error) {
	s.Logger.Info("Retrieving safes stats")
	safesChan, err := s.ListSafes()
	if err != nil {
		return nil, err
	}
	safes := make([]*safesmodels.ArkPCloudSafe, 0)
	for page := range safesChan {
		for _, safe := range page.Items {
			safes = append(safes, safe)
		}
	}
	var safesStats safesmodels.ArkPCloudSafesStats
	safesStats.SafesCount = len(safes)
	safesStats.SafesCountByLocation = make(map[string]int)
	safesStats.SafesCountByCreator = make(map[string]int)
	for _, safe := range safes {
		if _, ok := safesStats.SafesCountByLocation[safe.Location]; !ok {
			safesStats.SafesCountByLocation[safe.Location] = 0
		}
		if _, ok := safesStats.SafesCountByCreator[safe.Creator.Name]; !ok {
			safesStats.SafesCountByCreator[safe.Creator.Name] = 0
		}
		safesStats.SafesCountByLocation[safe.Location]++
		safesStats.SafesCountByCreator[safe.Creator.Name]++
	}
	return &safesStats, nil
}

// SafeMembersStats retrieves statistics about safe members for a specific safe.
func (s *ArkPCloudSafesService) SafeMembersStats(getSafeMembersStats *safesmodels.ArkPCloudGetSafeMembersStats) (*safesmodels.ArkPCloudSafeMembersStats, error) {
	s.Logger.Info("Retrieving safe members stats [%s]", getSafeMembersStats.SafeID)
	safeMembersChan, err := s.ListSafeMembers(&safesmodels.ArkPCloudListSafeMembers{SafeID: getSafeMembersStats.SafeID})
	if err != nil {
		return nil, err
	}
	safeMembers := make([]*safesmodels.ArkPCloudSafeMember, 0)
	for page := range safeMembersChan {
		for _, safeMember := range page.Items {
			safeMembers = append(safeMembers, safeMember)
		}
	}
	var safeMembersStats safesmodels.ArkPCloudSafeMembersStats
	safeMembersStats.SafeMembersCount = len(safeMembers)
	safeMembersStats.SafeMembersPermissionSets = make(map[string]int)
	safeMembersStats.SafeMembersTypesCount = make(map[string]int)
	for _, safeMember := range safeMembers {
		if safeMember.PermissionSet == "" {
			safeMember.PermissionSet = safesmodels.Custom
		}
		if _, ok := safeMembersStats.SafeMembersPermissionSets[safeMember.PermissionSet]; !ok {
			safeMembersStats.SafeMembersPermissionSets[safeMember.PermissionSet] = 0
		}
		if _, ok := safeMembersStats.SafeMembersTypesCount[safeMember.MemberType]; !ok {
			safeMembersStats.SafeMembersTypesCount[safeMember.MemberType] = 0
		}
		safeMembersStats.SafeMembersPermissionSets[safeMember.PermissionSet]++
		safeMembersStats.SafeMembersTypesCount[safeMember.MemberType]++
	}
	return &safeMembersStats, nil
}

// SafesMembersStats retrieves statistics about safe members for all safes.
func (s *ArkPCloudSafesService) SafesMembersStats() (*safesmodels.ArkPCloudSafesMembersStats, error) {
	s.Logger.Info("Retrieving safes members stats")
	safesChan, err := s.ListSafes()
	if err != nil {
		return nil, err
	}
	safesMembersStats := make(map[string]safesmodels.ArkPCloudSafeMembersStats)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error
	var once sync.Once

	for page := range safesChan {
		for _, safe := range page.Items {
			wg.Add(1)
			go func(safe *safesmodels.ArkPCloudSafe) {
				defer wg.Done()
				safeMembersStats, err := s.SafeMembersStats(&safesmodels.ArkPCloudGetSafeMembersStats{SafeID: safe.SafeURLID})
				if err != nil {
					once.Do(func() {
						firstErr = err
					})
					return
				}
				mu.Lock()
				safesMembersStats[safe.SafeName] = *safeMembersStats
				mu.Unlock()
			}(safe)
		}
	}
	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return &safesmodels.ArkPCloudSafesMembersStats{SafeMembersStats: safesMembersStats}, nil
}

// ServiceConfig returns the service configuration for the ArkPCloudSafesService.
func (s *ArkPCloudSafesService) ServiceConfig() services.ArkServiceConfig {
	return PCloudSafesServiceConfig
}
