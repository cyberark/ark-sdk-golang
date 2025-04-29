package auth

import (
	"errors"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"net/url"
	"slices"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/models/auth"
)

const (
	defaultExpirationGraceDeltaSeconds = 60
)

// ArkAuth is an interface that defines the methods for authentication in the Ark SDK.
type ArkAuth interface {
	// AuthenticatorName returns the name of the authenticator.
	AuthenticatorName() string
	// AuthenticatorHumanReadableName returns a human-readable name for the authenticator.
	AuthenticatorHumanReadableName() string
	// SupportedAuthMethods returns a list of supported authentication methods.
	SupportedAuthMethods() []auth.ArkAuthMethod
	// IsAuthenticated checks if the authentication is already loaded for the specified profile.
	IsAuthenticated(profile *models.ArkProfile) bool
	// DefaultAuthMethod returns the default authentication method and its settings.
	DefaultAuthMethod() (auth.ArkAuthMethod, auth.ArkAuthMethodSettings)
	// LoadAuthentication loads the authentication token for the specified profile and refreshes it if necessary.
	// It returns the authentication token and an error if any occurred.
	LoadAuthentication(profile *models.ArkProfile, refreshAuth bool) (*auth.ArkToken, error)
	// Authenticate performs authentication using the specified profile and authentication profile.
	// If profile is not passed (nil), will try to use the auth profile alone, but at least one of them needs to be passed
	// Secret may optionally be passed if needed for the authentication type
	// If force is true, it will force re-authentication even if a valid token is already present
	// If refreshAuth is true, it will attempt to refresh the token if it is expired
	// It returns the authentication token and an error if any occurred.
	Authenticate(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, secret *auth.ArkSecret, force bool, refreshAuth bool) (*auth.ArkToken, error)

	performAuthentication(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, secret *auth.ArkSecret, force bool) (*auth.ArkToken, error)
	performRefreshAuthentication(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, token *auth.ArkToken) (*auth.ArkToken, error)
}

// ArkAuthBase is a struct that implements the ArkAuth interface and provides common functionality for authentication.
type ArkAuthBase struct {
	Authenticator       ArkAuth
	Logger              *common.ArkLogger
	CacheAuthentication bool
	CacheKeyring        *common.ArkKeyring
	Token               *auth.ArkToken
	ActiveProfile       *models.ArkProfile
	ActiveAuthProfile   *auth.ArkAuthProfile
}

// NewArkAuthBase creates a new instance of ArkAuthBase.
func NewArkAuthBase(cacheAuthentication bool, name string, authenticator ArkAuth) *ArkAuthBase {
	logger := common.GetLogger(name, common.Unknown)
	var cacheKeyring *common.ArkKeyring
	if cacheAuthentication {
		cacheKeyring = common.NewArkKeyring(name)
	}
	return &ArkAuthBase{
		Authenticator:       authenticator,
		Logger:              logger,
		CacheAuthentication: cacheAuthentication,
		CacheKeyring:        cacheKeyring,
	}
}

// ResolveCachePostfix resolves the cache postfix for the authentication profile.
func (a *ArkAuthBase) ResolveCachePostfix(authProfile *auth.ArkAuthProfile) string {
	postfix := authProfile.Username
	if authProfile.AuthMethod == auth.Direct && authProfile.AuthMethodSettings != nil {
		directMethodSettings := authProfile.AuthMethodSettings.(auth.DirectArkAuthMethodSettings)
		if directMethodSettings.Endpoint != "" {
			parsedURL, _ := url.Parse(directMethodSettings.Endpoint)
			postfix = postfix + "_" + parsedURL.Host
		}
	}
	return postfix
}

// Authenticate performs authentication using the specified profile and authentication profile.
func (a *ArkAuthBase) Authenticate(profile *models.ArkProfile, authProfile *auth.ArkAuthProfile, secret *auth.ArkSecret, force bool, refreshAuth bool) (*auth.ArkToken, error) {
	if authProfile == nil && profile == nil {
		return nil, errors.New("either a profile or a specific auth profile must be supplied")
	}
	if authProfile == nil && profile != nil {
		if ap, ok := profile.AuthProfiles[a.Authenticator.AuthenticatorName()]; ok {
			authProfile = ap
		} else {
			return nil, errors.New(a.Authenticator.AuthenticatorHumanReadableName() + " [" + a.Authenticator.AuthenticatorName() + "] is not defined within the authentication profiles")
		}
	}
	if profile == nil {
		profilesLoader := profiles.DefaultProfilesLoader()
		profile, _ = (*profilesLoader).LoadDefaultProfile()
	}
	if !slices.Contains(a.Authenticator.SupportedAuthMethods(), authProfile.AuthMethod) && authProfile.AuthMethod != auth.Default {
		return nil, errors.New(a.Authenticator.AuthenticatorHumanReadableName() + " does not support authentication method " + string(authProfile.AuthMethod))
	}
	if authProfile.AuthMethod == auth.Default {
		authProfile.AuthMethod, authProfile.AuthMethodSettings = a.Authenticator.DefaultAuthMethod()
	}
	if slices.Contains(auth.ArkAuthMethodsRequireCredentials, authProfile.AuthMethod) && authProfile.Username == "" {
		return nil, errors.New(a.Authenticator.AuthenticatorHumanReadableName() + " requires a username and optionally a secret")
	}
	var token *auth.ArkToken
	var err error
	tokenRefreshed := false
	if a.CacheAuthentication && a.CacheKeyring != nil && !force {
		token, err = a.CacheKeyring.LoadToken(profile, a.ResolveCachePostfix(authProfile), false)
		if err != nil {
			return nil, err
		}
		if token != nil && time.Time(token.ExpiresIn).Before(time.Now()) {
			if refreshAuth && token.RefreshToken != "" {
				token, _ = a.Authenticator.performRefreshAuthentication(profile, authProfile, token)
				if token != nil {
					tokenRefreshed = true
				} else {
					token = nil
				}
			} else {
				token = nil
			}
		}
	}
	if token == nil {
		token, err = a.Authenticator.performAuthentication(profile, authProfile, secret, force)
		if err != nil {
			return nil, err
		}
		if token != nil && a.CacheAuthentication && a.CacheKeyring != nil {
			err := a.CacheKeyring.SaveToken(profile, token, a.ResolveCachePostfix(authProfile), false)
			if err != nil {
				return nil, err
			}
		}
	} else if refreshAuth && !tokenRefreshed {
		token, err = a.Authenticator.performRefreshAuthentication(profile, authProfile, token)
		if err != nil {
			return nil, err
		}
		if token != nil && a.CacheAuthentication && a.CacheKeyring != nil {
			err := a.CacheKeyring.SaveToken(profile, token, a.ResolveCachePostfix(authProfile), false)
			if err != nil {
				return nil, err
			}
		}
	}
	a.Token = token
	a.ActiveProfile = profile
	a.ActiveAuthProfile = authProfile
	return token, nil
}

// IsAuthenticated checks if the authentication is already loaded for the specified profile.
func (a *ArkAuthBase) IsAuthenticated(profile *models.ArkProfile) bool {
	var err error
	a.Logger.Info("Checking if [" + a.Authenticator.AuthenticatorName() + "] is authenticated")
	if a.Token != nil {
		a.Logger.Info("Token is already loaded")
		return true
	}
	if ap, ok := profile.AuthProfiles[a.Authenticator.AuthenticatorName()]; ok && a.CacheKeyring != nil {
		a.Token, err = a.CacheKeyring.LoadToken(profile, ap.Username, false)
		if err != nil {
			return false
		}
		if a.Token != nil && time.Time(a.Token.ExpiresIn).Before(time.Now()) {
			a.Token = nil
		} else {
			a.Logger.Info("Loaded token from cache successfully")
		}
		return a.Token != nil
	}
	return false
}

// LoadAuthentication loads the authentication token for the specified profile and refreshes it if necessary.
func (a *ArkAuthBase) LoadAuthentication(profile *models.ArkProfile, refreshAuth bool) (*auth.ArkToken, error) {
	var err error
	a.Logger.Info("Trying to load [" + a.Authenticator.AuthenticatorName() + "] authentication")
	if profile == nil {
		if a.ActiveProfile != nil {
			profile = a.ActiveProfile
		} else {
			profilesLoader := profiles.DefaultProfilesLoader()
			profile, _ = (*profilesLoader).LoadDefaultProfile()
		}
	}
	authProfile := a.ActiveAuthProfile
	if authProfile == nil {
		if ap, ok := profile.AuthProfiles[a.Authenticator.AuthenticatorName()]; ok {
			authProfile = ap
		}
	}
	if authProfile != nil {
		a.Logger.Info("Loading authentication for profile [" + profile.ProfileName + "] and auth profile [" + a.Authenticator.AuthenticatorName() + "] of type [" + string(authProfile.AuthMethod) + "]")
		if a.CacheKeyring != nil {
			a.Token, err = a.CacheKeyring.LoadToken(profile, a.ResolveCachePostfix(authProfile), false)
			if err != nil {
				return nil, err
			}
		}
		if refreshAuth {
			if a.Token != nil && time.Time(a.Token.ExpiresIn).Add(-time.Duration(defaultExpirationGraceDeltaSeconds)*time.Second).After(time.Now()) {
				a.Logger.Info("Token did not pass grace expiration, no need to refresh")
			} else {
				a.Logger.Info("Trying to refresh token authentication")
				a.Token, _ = a.Authenticator.performRefreshAuthentication(profile, authProfile, a.Token)
				if a.Token != nil && time.Time(a.Token.ExpiresIn).After(time.Now()) {
					a.Logger.Info("Token refreshed")
				}
				if a.Token != nil && a.CacheAuthentication && a.CacheKeyring != nil {
					err = a.CacheKeyring.SaveToken(profile, a.Token, a.ResolveCachePostfix(authProfile), false)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if a.Token != nil && time.Time(a.Token.ExpiresIn).Before(time.Now()) {
			a.Token = nil
		}
		if a.Token != nil {
			a.ActiveProfile = profile
			a.ActiveAuthProfile = authProfile
		}
		return a.Token, nil
	}
	return nil, nil
}
