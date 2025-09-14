package auth

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/auth/identity"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
)

const (
	ispAuthName              = "isp"
	ispAuthHumanReadableName = "Identity Security Platform"
)

// DefaultTokenLifetime is the default token lifetime in seconds.
const (
	DefaultTokenLifetime = 3600
)

var (
	ispAuthMethods               = []auth.ArkAuthMethod{auth.Identity, auth.IdentityServiceUser}
	ispDefaultAuthMethod         = auth.Identity
	ispDefaultAuthMethodSettings = auth.IdentityArkAuthMethodSettings{}
)

// ArkISPAuth is a struct that implements the ArkAuth interface for the Identity Security Platform.
type ArkISPAuth struct {
	ArkAuth
	*ArkAuthBase
}

// NewArkISPAuth creates a new instance of ArkISPAuth.
func NewArkISPAuth(cacheAuthentication bool) ArkAuth {
	authenticator := &ArkISPAuth{}
	var authInterface ArkAuth = authenticator
	baseAuth := NewArkAuthBase(cacheAuthentication, "ArkISPAuth", authInterface)
	authenticator.ArkAuthBase = baseAuth
	return authInterface
}

func (a *ArkISPAuth) performIdentityAuthentication(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, secret *auth.ArkSecret, force bool) (*auth.ArkToken, error) {
	methodSettings := authProfile.AuthMethodSettings.(*auth.IdentityArkAuthMethodSettings)
	identityAuth, err := identity.NewArkIdentity(
		authProfile.Username,
		secret.Secret,
		methodSettings.IdentityURL,
		methodSettings.IdentityTenantSubdomain,
		methodSettings.IdentityMFAMethod,
		a.Logger,
		a.CacheAuthentication,
		a.CacheAuthentication,
		profile,
	)
	if err != nil {
		a.Logger.Error("Failed to create identity security platform object: %v", err)
		return nil, err
	}
	err = identityAuth.AuthIdentity(profile, common.IsInteractive() && methodSettings.IdentityMFAInteractive, force)
	if err != nil {
		a.Logger.Error("Failed to authenticate to identity security platform: %v", err)
		return nil, err
	}
	env := commonmodels.GetDeployEnv()
	tokenLifetime := identityAuth.SessionDetails().TokenLifetime
	if tokenLifetime == 0 {
		tokenLifetime = DefaultTokenLifetime
	}
	marshaledCookies, err := common.MarshalCookies(identityAuth.Session().GetCookieJar())
	if err != nil {
		a.Logger.Error("Failed to marshal cookies: %v", err)
		return nil, err
	}
	return &auth.ArkToken{
		Token:        identityAuth.SessionToken(),
		Username:     authProfile.Username,
		Endpoint:     identityAuth.IdentityURL(),
		TokenType:    auth.JWT,
		AuthMethod:   auth.Identity,
		ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(time.Duration(tokenLifetime) * time.Second)),
		RefreshToken: identityAuth.SessionDetails().RefreshToken,
		Metadata: map[string]interface{}{
			"env":     env,
			"cookies": base64.StdEncoding.EncodeToString(marshaledCookies),
		},
	}, nil
}

func (a *ArkISPAuth) performIdentityRefreshAuthentication(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, token *auth.ArkToken) (*auth.ArkToken, error) {
	methodSettings := authProfile.AuthMethodSettings.(*auth.IdentityArkAuthMethodSettings)
	identityAuth, err := identity.NewArkIdentity(
		authProfile.Username,
		"",
		methodSettings.IdentityURL,
		methodSettings.IdentityTenantSubdomain,
		methodSettings.IdentityMFAMethod,
		a.Logger,
		a.CacheAuthentication,
		a.CacheAuthentication,
		profile,
	)
	if err != nil {
		a.Logger.Error("Failed to create identity security platform object: %v", err)
		return nil, err
	}
	err = identityAuth.RefreshAuthIdentity(profile, methodSettings.IdentityMFAInteractive, false)
	if err != nil {
		a.Logger.Error("Failed to refresh authentication to identity security platform: %v", err)
		return nil, err
	}
	env := commonmodels.GetDeployEnv()
	tokenLifetime := identityAuth.SessionDetails().TokenLifetime
	if tokenLifetime == 0 {
		tokenLifetime = DefaultTokenLifetime
	}
	marshaledCookies, err := common.MarshalCookies(identityAuth.Session().GetCookieJar())
	if err != nil {
		a.Logger.Error("Failed to marshal cookies: %v", err)
		return nil, err
	}
	return &auth.ArkToken{
		Token:        identityAuth.SessionToken(),
		Username:     authProfile.Username,
		Endpoint:     identityAuth.IdentityURL(),
		TokenType:    auth.JWT,
		AuthMethod:   auth.Identity,
		ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(time.Duration(tokenLifetime) * time.Second)),
		RefreshToken: identityAuth.SessionDetails().RefreshToken,
		Metadata: map[string]interface{}{
			"env":     env,
			"cookies": base64.StdEncoding.EncodeToString(marshaledCookies),
		},
	}, nil
}

func (a *ArkISPAuth) performIdentityServiceUserAuthentication(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, secret *auth.ArkSecret, force bool) (*auth.ArkToken, error) {
	if secret == nil {
		return nil, errors.New("token secret is required for identity service user auth")
	}
	methodSettings := authProfile.AuthMethodSettings.(*auth.IdentityServiceUserArkAuthMethodSettings)
	identityAuth, err := identity.NewArkIdentityServiceUser(
		authProfile.Username,
		secret.Secret,
		methodSettings.IdentityAuthorizationApplication,
		methodSettings.IdentityURL,
		methodSettings.IdentityTenantSubdomain,
		a.Logger,
		a.CacheAuthentication,
		a.CacheAuthentication,
		profile,
	)
	if err != nil {
		a.Logger.Error("Failed to create identity security platform object with service user: %v", err)
		return nil, err
	}
	err = identityAuth.AuthIdentity(profile, force)
	if err != nil {
		a.Logger.Error("Failed to authenticate to identity security platform with service user: %v", err)
		return nil, err
	}
	env := commonmodels.GetDeployEnv()
	marshaledCookies, err := common.MarshalCookies(identityAuth.Session().GetCookieJar())
	if err != nil {
		a.Logger.Error("Failed to marshal cookies: %v", err)
		return nil, err
	}
	return &auth.ArkToken{
		Token:      identityAuth.SessionToken(),
		Username:   authProfile.Username,
		Endpoint:   identityAuth.IdentityURL(),
		TokenType:  auth.JWT,
		AuthMethod: auth.Identity,
		ExpiresIn:  commonmodels.ArkRFC3339Time(time.Now().Add(4 * time.Hour)),
		Metadata: map[string]interface{}{
			"env":     env,
			"cookies": base64.StdEncoding.EncodeToString(marshaledCookies),
		},
	}, nil
}

// performAuthentication performs authentication to the ISP using the specified auth method.
func (a *ArkISPAuth) performAuthentication(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, secret *auth.ArkSecret, force bool) (*auth.ArkToken, error) {
	a.Logger.Info("Performing authentication to ISP")
	switch authProfile.AuthMethod {
	case auth.Identity, auth.Default:
		return a.performIdentityAuthentication(profile, authProfile, secret, force)
	case auth.IdentityServiceUser:
		return a.performIdentityServiceUserAuthentication(profile, authProfile, secret, force)
	default:
		return nil, errors.New("given auth method is not supported")
	}
}

// PerformRefreshAuthentication performs refresh authentication to the ISP.
func (a *ArkISPAuth) performRefreshAuthentication(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, token *auth.ArkToken) (*auth.ArkToken, error) {
	a.Logger.Info("Performing refresh authentication to ISP")
	if authProfile.AuthMethod == auth.Identity || authProfile.AuthMethod == auth.Default {
		return a.performIdentityRefreshAuthentication(profile, authProfile, token)
	}
	return token, nil
}

// LoadAuthentication loads the authentication token from the cache or performs authentication if not found.
func (a *ArkISPAuth) LoadAuthentication(profile *models.ArkProfile, refreshAuth bool) (*auth.ArkToken, error) {
	return a.ArkAuthBase.LoadAuthentication(profile, refreshAuth)
}

// Authenticate performs authentication using the specified profile and authentication profile.
func (a *ArkISPAuth) Authenticate(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, secret *auth.ArkSecret, force bool, refreshAuth bool) (*auth.ArkToken, error) {
	return a.ArkAuthBase.Authenticate(profile, authProfile, secret, force, refreshAuth)
}

// IsAuthenticated checks if the user is authenticated using the specified profile.
func (a *ArkISPAuth) IsAuthenticated(profile *models.ArkProfile) bool {
	return a.ArkAuthBase.IsAuthenticated(profile)
}

// AuthenticatorName returns the name of the ISP authenticator.
func (a *ArkISPAuth) AuthenticatorName() string {
	return ispAuthName
}

// AuthenticatorHumanReadableName returns the human-readable name of the ISP authenticator.
func (a *ArkISPAuth) AuthenticatorHumanReadableName() string {
	return ispAuthHumanReadableName
}

// SupportedAuthMethods returns the supported authentication methods for the ISP authenticator.
func (a *ArkISPAuth) SupportedAuthMethods() []auth.ArkAuthMethod {
	return ispAuthMethods
}

// DefaultAuthMethod returns the default authentication method and its settings for the ISP authenticator.
func (a *ArkISPAuth) DefaultAuthMethod() (auth.ArkAuthMethod, auth.ArkAuthMethodSettings) {
	return ispDefaultAuthMethod, ispDefaultAuthMethodSettings
}
