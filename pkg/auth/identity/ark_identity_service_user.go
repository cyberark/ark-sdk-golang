package identity

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
)

// ArkIdentityServiceUser is a struct that represents identity authentication with service user.
type ArkIdentityServiceUser struct {
	username            string
	token               string
	appName             string
	identityURL         string
	logger              *common.ArkLogger
	keyring             *common.ArkKeyring
	cacheAuthentication bool
	session             *common.ArkClient
	sessionToken        string
	sessionExp          commonmodels.ArkRFC3339Time
}

// NewArkIdentityServiceUser creates a new instance of ArkIdentityServiceUser.
func NewArkIdentityServiceUser(username string, token string, appName string, identityURL string, identityTenantSubdomain string, logger *common.ArkLogger, cacheAuthentication bool, loadCache bool, cacheProfile *models.ArkProfile) (*ArkIdentityServiceUser, error) {
	identityServiceAuth := &ArkIdentityServiceUser{
		username:            username,
		token:               token,
		appName:             appName,
		identityURL:         identityURL,
		logger:              logger,
		cacheAuthentication: cacheAuthentication,
	}
	var err error
	if identityURL == "" {
		if identityTenantSubdomain != "" {
			identityURL, err = ResolveTenantFqdnFromTenantSubdomain(identityTenantSubdomain, commonmodels.GetDeployEnv())
		} else {
			tenantSuffix := username[strings.Index(username, "@"):]
			identityURL, err = ResolveTenantFqdnFromTenantSuffix(tenantSuffix, commonmodels.IdentityEnvUrls[commonmodels.GetDeployEnv()])
		}
	}
	if err != nil {
		return nil, err
	}
	identityServiceAuth.session = common.NewSimpleArkClient(identityURL)

	if cacheAuthentication {
		identityServiceAuth.keyring = common.NewArkKeyring(strings.ToLower("ArkIdentity"))
	}
	if loadCache && cacheAuthentication && cacheProfile != nil {
		identityServiceAuth.loadCache(cacheProfile)
	}
	return identityServiceAuth, nil
}

func (ai *ArkIdentityServiceUser) loadCache(profile *models.ArkProfile) bool {
	if ai.keyring != nil && profile != nil {
		token, err := ai.keyring.LoadToken(profile, ai.username+"_identity_service_user", false)
		if err != nil {
			ai.logger.Error("Error loading token from cache: %v", err.Error())
			return false
		}
		if token != nil && token.Username == ai.username {
			ai.sessionToken = token.Token
			ai.sessionExp = token.ExpiresIn
			ai.session.UpdateToken(ai.sessionToken, "Bearer")
			return true
		}
	}
	return false
}

func (ai *ArkIdentityServiceUser) saveCache(profile *models.ArkProfile) error {
	if ai.keyring != nil && profile != nil && ai.sessionToken != "" {
		ai.sessionExp = commonmodels.ArkRFC3339Time(time.Now().Add(4 * time.Hour))
		err := ai.keyring.SaveToken(profile, &auth.ArkToken{
			Token:      ai.sessionToken,
			Username:   ai.username,
			Endpoint:   ai.session.BaseURL,
			TokenType:  auth.Internal,
			AuthMethod: auth.Other,
			ExpiresIn:  ai.sessionExp,
		}, ai.username+"_identity_service_user", false)
		if err != nil {
			return err
		}
	}
	return nil
}

// AuthIdentity Authenticates to Identity with a service user.
// This method creates an auth token and authorizes to the service.
func (ai *ArkIdentityServiceUser) AuthIdentity(profile *models.ArkProfile, force bool) error {
	ai.logger.Info("Authenticating to service user via endpoint [%s]", ai.identityURL)
	if ai.cacheAuthentication && !force && ai.loadCache(profile) {
		if time.Time(ai.sessionExp).After(time.Now()) {
			ai.logger.Info("Loaded identity service user details from cache")
			return nil
		}
	}
	ai.session.UpdateToken(
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", ai.username, ai.token))),
		"Basic",
	)
	response, err := ai.session.Post(
		context.Background(),
		fmt.Sprintf("Oauth2/Token/%s", ai.appName),
		map[string]interface{}{
			"grant_type": "client_credentials",
			"scope":      "api",
		},
	)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ai.logger.Warning("Error closing response body")
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed logging in to identity service user")
	}

	var authResult map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&authResult); err != nil {
		return err
	}

	accessToken, ok := authResult["access_token"].(string)
	if !ok {
		return fmt.Errorf("failed logging in to identity service user, access token not found")
	}
	ai.session.UpdateToken(accessToken, "Bearer")
	response, err = ai.session.Get(
		context.Background(),
		fmt.Sprintf("OAuth2/Authorize/%s", ai.appName),
		map[string]string{
			"client_id":     ai.appName,
			"response_type": "id_token",
			"scope":         "openid profile api",
			"redirect_uri":  "https://cyberark.cloud/redirect",
		},
	)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ai.logger.Warning("Error closing response body")
		}
	}(response.Body)

	if response.StatusCode != http.StatusFound || response.Header.Get("Location") == "" {
		return fmt.Errorf("failed to authorize to application")
	}

	locationHeader := response.Header.Get("Location")
	locationHeaderSplitted := strings.Split(locationHeader, "#")
	if len(locationHeaderSplitted) != 2 {
		return fmt.Errorf("failed to parse location header to retrieve token from")
	}

	parsedQuery, err := url.ParseQuery(locationHeaderSplitted[1])
	if err != nil {
		return err
	}

	idTokens, ok := parsedQuery["id_token"]
	if !ok || len(idTokens) != 1 {
		return fmt.Errorf("failed to parse id token from location header")
	}

	ai.sessionToken = idTokens[0]
	ai.session.UpdateToken(ai.sessionToken, "Bearer")
	ai.sessionExp = commonmodels.ArkRFC3339Time(time.Now().Add(4 * time.Hour))
	ai.logger.Info("Created a service user session via endpoint [%s] with user [%s] to platform", ai.identityURL, ai.username)

	if ai.cacheAuthentication {
		if err := ai.saveCache(profile); err != nil {
			return err
		}
	}

	return nil
}

// Session returns the current identity session
func (ai *ArkIdentityServiceUser) Session() *common.ArkClient {
	return ai.session
}

// SessionToken returns the current identity session token if logged in
func (ai *ArkIdentityServiceUser) SessionToken() string {
	return ai.sessionToken
}

// IdentityURL returns the current identity URL
func (ai *ArkIdentityServiceUser) IdentityURL() string {
	return ai.session.BaseURL
}
