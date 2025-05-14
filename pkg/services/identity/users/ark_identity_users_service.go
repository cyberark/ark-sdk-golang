package users

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	rolesmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/identity/roles"
	usersmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/identity/users"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity/directories"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity/roles"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	createUserURL        = "CDirectoryService/CreateUser"
	deleteUserURL        = "CDirectoryService/DeleteUser"
	updateUserURL        = "CDirectoryService/ChangeUser"
	removeUsersURL       = "UserMgmt/RemoveUsers"
	resetUserPasswordURL = "UserMgmt/ResetUserPassword"
	redrockQueryURL      = "Redrock/query"
	userInfoURL          = "OAuth2/UserInfo/__idaptive_cybr_user_oidc"
)

// IdentityUsersServiceConfig is the configuration for the identity users service.
var IdentityUsersServiceConfig = services.ArkServiceConfig{
	ServiceName:                "identity-users",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkIdentityUsersService is the service for managing identity users.
type ArkIdentityUsersService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkIdentityUsersService creates a new instance of ArkIdentityUsersService.
func NewArkIdentityUsersService(authenticators ...auth.ArkAuth) (*ArkIdentityUsersService, error) {
	identityUsersService := &ArkIdentityUsersService{}
	var identityUsersServiceInterface services.ArkService = identityUsersService
	baseService, err := services.NewArkBaseService(identityUsersServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "", "", "api/idadmin", identityUsersService.refreshIdentityUsersAuth)
	if err != nil {
		return nil, err
	}
	client.UpdateHeaders(map[string]string{
		"X-IDAP-NATIVE-CLIENT": "true",
	})
	identityUsersService.client = client
	identityUsersService.ispAuth = ispAuth
	identityUsersService.ArkBaseService = baseService
	return identityUsersService, nil
}

func (s *ArkIdentityUsersService) refreshIdentityUsersAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// CreateUser creates a new user in the identity service.
func (s *ArkIdentityUsersService) CreateUser(createUser *usersmodels.ArkIdentityCreateUser) (*usersmodels.ArkIdentityUser, error) {
	if createUser.Username == "" {
		createUser.Username = fmt.Sprintf("ark_user_%s", common.RandomString(10))
	}
	if createUser.DisplayName == "" {
		createUser.DisplayName = fmt.Sprintf("%s %s", strings.Title(common.RandomString(5)), strings.Title(common.RandomString(7)))
	}
	if createUser.Email == "" {
		createUser.Email = fmt.Sprintf("%s@email.com", strings.ToLower(common.RandomString(6)))
	}
	if createUser.MobileNumber == "" {
		createUser.MobileNumber = fmt.Sprintf("+44-987-654-%s", common.RandomNumberString(4))
	}
	if createUser.Password == "" {
		createUser.Password = common.RandomPassword(25)
	}
	if createUser.Roles == nil {
		createUser.Roles = usersmodels.DefaultAdminRoles
	}
	s.Logger.Info("Creating identity user [%s]", createUser.Username)
	if createUser.Suffix == "" {
		directoriesService, err := directories.NewArkIdentityDirectoriesService(s.ispAuth)
		if err != nil {
			return nil, err
		}
		createUser.Suffix, err = directoriesService.TenantDefaultSuffix()
		if err != nil {
			return nil, err
		}
	}
	createUserRequest := map[string]interface{}{
		"DisplayName":             createUser.DisplayName,
		"Name":                    fmt.Sprintf("%s@%s", createUser.Username, createUser.Suffix),
		"Mail":                    createUser.Email,
		"Password":                createUser.Password,
		"MobileNumber":            createUser.MobileNumber,
		"InEverybodyRole":         "true",
		"InSysAdminRole":          "false",
		"ForcePasswordChangeNext": "false",
		"SendEmailInvite":         "false",
		"SendSmsInvite":           "false",
	}
	response, err := s.client.Post(context.Background(), createUserURL, createUserRequest)
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
		return nil, fmt.Errorf("failed to create user - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result["success"].(bool) == false {
		return nil, fmt.Errorf("failed to create user - [%v]", result)
	}
	if createUser.Roles != nil {
		rolesService, err := roles.NewArkIdentityRolesService(s.ispAuth)
		if err != nil {
			return nil, err
		}
		for _, role := range createUser.Roles {
			err := rolesService.AddUserToRole(&rolesmodels.ArkIdentityAddUserToRole{
				Username: fmt.Sprintf("%s@%s", createUser.Username, createUser.Suffix),
				RoleName: role,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	userID := result["Result"].(string)
	s.Logger.Info("User created successfully with id [%s]", userID)
	return &usersmodels.ArkIdentityUser{
		UserID:       userID,
		Username:     fmt.Sprintf("%s@%s", createUser.Username, createUser.Suffix),
		DisplayName:  createUser.DisplayName,
		Email:        createUser.Email,
		MobileNumber: createUser.MobileNumber,
		Roles:        createUser.Roles,
	}, nil
}

// UpdateUser updates an existing user in the identity service.
func (s *ArkIdentityUsersService) UpdateUser(updateUser *usersmodels.ArkIdentityUpdateUser) error {
	s.Logger.Info("Updating identity user [%s]", updateUser.Username)
	var err error
	if updateUser.Username != "" && updateUser.UserID == "" {
		updateUser.UserID, err = s.UserIDByName(&usersmodels.ArkIdentityUserIDByName{Username: updateUser.Username})
		if err != nil {
			return err
		}
	}
	updateMap := make(map[string]interface{})
	if updateUser.NewUsername != "" {
		if !strings.Contains(updateUser.NewUsername, "@") {
			tenantSuffix := strings.Split(updateUser.Username, "@")[1]
			updateUser.NewUsername = fmt.Sprintf("%s@%s", updateUser.NewUsername, tenantSuffix)
		}
		updateMap["Name"] = updateUser.NewUsername
	}
	if updateUser.DisplayName != "" {
		updateMap["DisplayName"] = updateUser.DisplayName
	}
	if updateUser.Email != "" {
		updateMap["Mail"] = updateUser.Email
	}
	if updateUser.MobileNumber != "" {
		updateMap["MobileNumber"] = updateUser.MobileNumber
	}
	updateMap["ID"] = updateUser.UserID
	response, err := s.client.Post(context.Background(), updateUserURL, updateMap)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update user - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if result["success"].(bool) == false {
		return fmt.Errorf("failed to update user - [%v]", result)
	}
	s.Logger.Info("User updated successfully")
	return nil
}

// DeleteUser deletes an existing user in the identity service.
func (s *ArkIdentityUsersService) DeleteUser(deleteUser *usersmodels.ArkIdentityDeleteUser) error {
	s.Logger.Info("Deleting identity user [%s]", deleteUser.Username)
	if deleteUser.Username == "" && deleteUser.UserID == "" {
		return fmt.Errorf("userID or username is required")
	}
	deleteMap := make(map[string]interface{})
	deleteMap["ID"] = deleteUser.UserID
	if deleteUser.UserID == "" && deleteUser.Username != "" {
		deleteMap["ID"] = deleteUser.Username
	}
	response, err := s.client.Post(context.Background(), deleteUserURL, deleteMap)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete user - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if result["success"].(bool) == false {
		return fmt.Errorf("failed to delete user - [%v]", result)
	}
	s.Logger.Info("User deleted successfully")
	return nil
}

// DeleteUsers deletes multiple users in the identity service.
func (s *ArkIdentityUsersService) DeleteUsers(deleteUsers *usersmodels.ArkIdentityDeleteUsers) error {
	s.Logger.Info("Deleting identity users [%v]", deleteUsers.UserIDs)
	if len(deleteUsers.UserIDs) == 0 {
		return fmt.Errorf("userIDs is required")
	}
	deleteMap := make(map[string]interface{})
	deleteMap["Users"] = deleteUsers.UserIDs
	response, err := s.client.Post(context.Background(), removeUsersURL, deleteMap)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete users - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if result["success"].(bool) == false {
		return fmt.Errorf("failed to delete users - [%v]", result)
	}
	s.Logger.Info("Users deleted successfully")
	return nil
}

// UserIDByName retrieves the user ID by username.
func (s *ArkIdentityUsersService) UserIDByName(user *usersmodels.ArkIdentityUserIDByName) (string, error) {
	s.Logger.Info("Getting identity user ID by name [%s]", user.Username)
	if user.Username == "" {
		return "", fmt.Errorf("username is required")
	}
	redrockQuery := map[string]interface{}{
		"Script": fmt.Sprintf("Select ID, Username from User WHERE Username='%s'", strings.ToLower(user.Username)),
	}
	response, err := s.client.Post(context.Background(), redrockQueryURL, redrockQuery)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user ID - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	if result["success"].(bool) == false {
		return "", fmt.Errorf("failed to get user ID - [%v]", result)
	}
	if len(result["Result"].(map[string]interface{})["Results"].([]interface{})) == 0 {
		return "", fmt.Errorf("failed to retrieve user id by name")
	}
	return result["Result"].(map[string]interface{})["Results"].([]interface{})[0].(map[string]interface{})["Row"].(map[string]interface{})["ID"].(string), nil
}

// UserByName retrieves the user by username.
func (s *ArkIdentityUsersService) UserByName(user *usersmodels.ArkIdentityUserByName) (*usersmodels.ArkIdentityUser, error) {
	s.Logger.Info("Getting identity user by name [%s]", user.Username)
	if user.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	redrockQuery := map[string]interface{}{
		"Script": fmt.Sprintf("Select ID, Username, DisplayName, Email, MobileNumber, LastLogin from User WHERE Username='%s'", strings.ToLower(user.Username)),
	}
	response, err := s.client.Post(context.Background(), redrockQueryURL, redrockQuery)
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
		return nil, fmt.Errorf("failed to get user - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if !result["success"].(bool) || len(result["Result"].(map[string]interface{})["Results"].([]interface{})) == 0 {
		return nil, fmt.Errorf("failed to retrieve user id by name")
	}

	userRow := result["Result"].(map[string]interface{})["Results"].([]interface{})[0].(map[string]interface{})["Row"].(map[string]interface{})
	var lastLogin *time.Time

	if rawLastLogin, ok := userRow["LastLogin"].(string); ok {
		parts := strings.Split(rawLastLogin, "(")
		if len(parts) > 1 {
			timestamp := strings.Split(parts[1], ")")[0]
			timestamp = fmt.Sprintf("%s.%s", timestamp[:10], timestamp[10:]) // for milliseconds
			parsedTime, err := strconv.ParseFloat(timestamp, 64)
			if err == nil {
				t := time.Unix(0, int64(parsedTime*1e6)).UTC()
				lastLogin = &t
			} else {
				s.Logger.Debug(fmt.Sprintf("Failed to parse last login [%s] [%s]", rawLastLogin, err.Error()))
			}
		}
	}
	return &usersmodels.ArkIdentityUser{
		UserID:       userRow["ID"].(string),
		Username:     userRow["Username"].(string),
		DisplayName:  userRow["DisplayName"].(string),
		Email:        userRow["Email"].(string),
		MobileNumber: userRow["MobileNumber"].(string),
		LastLogin:    lastLogin,
	}, nil
}

// UserByID retrieves the user by user ID.
func (s *ArkIdentityUsersService) UserByID(userByID *usersmodels.ArkIdentityUserByID) (*usersmodels.ArkIdentityUser, error) {
	s.Logger.Info("Getting identity user by id [%s]", userByID.UserID)
	if userByID.UserID == "" {
		return nil, fmt.Errorf("userID is required")
	}
	redrockQuery := map[string]interface{}{
		"Script": fmt.Sprintf("Select ID, Username, DisplayName, Email, MobileNumber, LastLogin from User WHERE ID='%s'", userByID.UserID),
	}
	response, err := s.client.Post(context.Background(), redrockQueryURL, redrockQuery)
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
		return nil, fmt.Errorf("failed to get user - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if !result["success"].(bool) || len(result["Result"].(map[string]interface{})["Results"].([]interface{})) == 0 {
		return nil, fmt.Errorf("failed to retrieve user id by name")
	}

	userRow := result["Result"].(map[string]interface{})["Results"].([]interface{})[0].(map[string]interface{})["Row"].(map[string]interface{})
	var lastLogin *time.Time

	if rawLastLogin, ok := userRow["LastLogin"].(string); ok {
		parts := strings.Split(rawLastLogin, "(")
		if len(parts) > 1 {
			timestamp := strings.Split(parts[1], ")")[0]
			timestamp = fmt.Sprintf("%s.%s", timestamp[:10], timestamp[10:]) // for milliseconds
			parsedTime, err := strconv.ParseFloat(timestamp, 64)
			if err == nil {
				t := time.Unix(0, int64(parsedTime*1e6)).UTC()
				lastLogin = &t
			} else {
				s.Logger.Debug(fmt.Sprintf("Failed to parse last login [%s] [%s]", rawLastLogin, err.Error()))
			}
		}
	}
	return &usersmodels.ArkIdentityUser{
		UserID:       userRow["ID"].(string),
		Username:     userRow["Username"].(string),
		DisplayName:  userRow["DisplayName"].(string),
		Email:        userRow["Email"].(string),
		MobileNumber: userRow["MobileNumber"].(string),
		LastLogin:    lastLogin,
	}, nil
}

// ResetUserPassword resets the password for an existing user in the identity service.
func (s *ArkIdentityUsersService) ResetUserPassword(resetUserPassword *usersmodels.ArkIdentityResetUserPassword) error {
	s.Logger.Info("Resetting identity user password [%s]", resetUserPassword.Username)
	userID, err := s.UserIDByName(&usersmodels.ArkIdentityUserIDByName{Username: resetUserPassword.Username})
	if err != nil {
		return err
	}
	resetPasswordMap := make(map[string]interface{})
	resetPasswordMap["ID"] = userID
	resetPasswordMap["newPassword"] = resetUserPassword.NewPassword
	response, err := s.client.Post(context.Background(), resetUserPasswordURL, resetPasswordMap)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to reset user password - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if result["success"].(bool) == false {
		return fmt.Errorf("failed to reset user password - [%v]", result)
	}
	s.Logger.Info("User password reset successfully")
	return nil
}

// UserInfo retrieves the user info from the identity service.
func (s *ArkIdentityUsersService) UserInfo() (*usersmodels.ArkIdentityUserInfo, error) {
	s.Logger.Info("Getting identity user info")
	userInfoMap := map[string]interface{}{
		"Scopes": []string{"userInfo"},
	}
	response, err := s.client.Post(context.Background(), userInfoURL, userInfoMap)
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
		return nil, fmt.Errorf("failed to get user info - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	var userInfo usersmodels.ArkIdentityUserInfo
	err = mapstructure.Decode(result, &userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

// ServiceConfig returns the service configuration for the ArkIdentityUsersService.
func (s *ArkIdentityUsersService) ServiceConfig() services.ArkServiceConfig {
	return IdentityUsersServiceConfig
}
