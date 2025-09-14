// Package keyring provides keyring utilities for the ARK SDK.
//
// This package implements cross-platform keyring support with automatic fallback mechanisms
// for different environments including Docker containers, WSL, and various operating systems.
// The keyring handles token expiration, automatic cleanup, and secure storage of authentication
// credentials for ARK SDK applications.
package keyring

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/models/auth"
)

// Env vars and definitions
const (
	ArkBasicKeyringOverrideEnvVar      = "ARK_BASIC_KEYRING"
	DBusSessionEnvVar                  = "DBUS_SESSION_BUS_ADDRESS"
	DefaultExpirationGraceDeltaSeconds = 60
	MaxKeyringRecordTimeHours          = 12
)

// ArkKeyringImpl defines the interface for keyring operations.
type ArkKeyringImpl interface {
	SetPassword(serviceName string, username string, password string) error
	GetPassword(serviceName string, username string) (string, error)
	DeletePassword(serviceName string, username string) error
	ClearAllPasswords() error
}

// ArkKeyring represents a keyring for storing and retrieving authentication tokens.
//
// ArkKeyring provides a secure storage mechanism for authentication tokens with
// automatic fallback to basic keyring when system keyrings are unavailable. It
// supports different keyring backends based on the operating system and environment,
// including Docker containers and WSL environments.
//
// The keyring handles token expiration, refresh token management, and automatic
// cleanup of expired tokens based on configurable time limits.
type ArkKeyring struct {
	serviceName string
	logger      *common.ArkLogger
}

// NewArkKeyring creates a new instance of ArkKeyring with the specified service name.
//
// The service name is used as a namespace for storing tokens in the keyring,
// allowing multiple applications or services to use the same keyring without
// conflicts.
//
// Parameters:
//   - serviceName: The name used to identify this service's tokens in the keyring
//
// Returns a new ArkKeyring instance configured with the provided service name
// and a logger for keyring operations.
//
// Example:
//
//	keyring := NewArkKeyring("myapp")
func NewArkKeyring(serviceName string) *ArkKeyring {
	return &ArkKeyring{
		serviceName: serviceName,
		logger:      common.GetLogger("ArkKeyring", common.Unknown),
	}
}

func (a *ArkKeyring) isDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	if data, err := os.ReadFile("/proc/self/cgroup"); err == nil {
		return strings.Contains(string(data), "docker")
	}
	return false
}

func (a *ArkKeyring) isWSL() bool {
	if data, err := os.ReadFile("/proc/version"); err == nil {
		return strings.Contains(string(data), "Microsoft")
	}
	return false
}

// GetKeyring returns a keyring instance based on the operating system and environment.
//
// GetKeyring automatically selects the most appropriate keyring backend based on
// the current environment. It falls back to basic keyring in Docker containers,
// WSL environments, when the ARK_BASIC_KEYRING environment variable is set, or
// when enforceBasicKeyring is true. For Windows, macOS, and Linux systems with
// proper D-Bus session, it attempts to use system-specific secure storage.
//
// Parameters:
//   - enforceBasicKeyring: When true, forces the use of basic keyring regardless of environment
//
// Returns a ArkBasicKeyring instance configured for the current environment, or an
// error if keyring initialization fails.
//
// Example:
//
//	keyring, err := arkKeyring.GetKeyring(false)
//	if err != nil {
//	    // handle error
//	}
func (a *ArkKeyring) GetKeyring(enforceBasicKeyring bool) (ArkKeyringImpl, error) {
	if a.isDocker() || a.isWSL() || os.Getenv(ArkBasicKeyringOverrideEnvVar) != "" || enforceBasicKeyring {
		return NewArkBasicKeyring(), nil
	}
	if runtime.GOOS == "windows" {
		return NewArkOSProvidedKeyring(NewArkBasicKeyring()), nil
	}
	if runtime.GOOS == "darwin" {
		return NewArkOSProvidedKeyring(NewArkBasicKeyring()), nil
	}
	if runtime.GOOS == "linux" && os.Getenv(DBusSessionEnvVar) != "" {
		return NewArkOSProvidedKeyring(NewArkBasicKeyring()), nil
	}
	return NewArkBasicKeyring(), nil
}

// SaveToken saves an authentication token to the keyring for the specified profile and postfix.
//
// SaveToken stores the provided token in the keyring using a composite key format
// of "serviceName-postfix" and the profile name. The token is serialized to JSON
// before storage. If the initial save fails and enforceBasicKeyring is false,
// it automatically falls back to basic keyring storage.
//
// Parameters:
//   - profile: The ARK profile containing the profile name used as the keyring username
//   - token: The authentication token to store in the keyring
//   - postfix: A suffix added to the service name to create unique keys for different token types
//   - enforceBasicKeyring: When true, uses basic keyring without attempting system keyring first
//
// Returns an error if the token cannot be saved to any available keyring backend.
//
// Example:
//
//	err := keyring.SaveToken(profile, token, "access", false)
//	if err != nil {
//	    // handle save error
//	}
func (a *ArkKeyring) SaveToken(profile *models.ArkProfile, token *auth.ArkToken, postfix string, enforceBasicKeyring bool) error {
	a.logger.Info("Trying to save token [%s-%s] of profile [%s]", a.serviceName, postfix, profile.ProfileName)
	kr, err := a.GetKeyring(enforceBasicKeyring)
	if err != nil {
		return err
	}
	tokenData, err := json.Marshal(token)
	if err != nil {
		return err
	}
	if err := kr.SetPassword(profile.ProfileName, a.serviceName+"-"+postfix, string(tokenData)); err != nil {
		if !enforceBasicKeyring {
			a.logger.Warning("Falling back to basic keyring as we failed to save token with keyring [%v]", kr)
			return a.SaveToken(profile, token, postfix, true)
		}
		a.logger.Warning("Failed to save token [%v]", err)
		return err
	}
	a.logger.Info("Saved token successfully")
	return nil
}

// LoadToken loads an authentication token from the keyring for the specified profile and postfix.
//
// LoadToken retrieves and validates a stored authentication token from the keyring.
// It performs automatic token expiration checking and cleanup. For tokens without
// refresh capability that are expired beyond the grace period, the token is removed
// and nil is returned. For tokens with refresh capability that have been cached
// too long, they are also removed and nil is returned. If the initial load fails
// and enforceBasicKeyring is false, it automatically falls back to basic keyring.
//
// Parameters:
//   - profile: The ARK profile containing the profile name used as the keyring username
//   - postfix: A suffix added to the service name to match the key used during SaveToken
//   - enforceBasicKeyring: When true, uses basic keyring without attempting system keyring first
//
// Returns the loaded token if found and valid, nil if no token exists or token
// is expired, or an error if the keyring operation fails.
//
// Example:
//
//	token, err := keyring.LoadToken(profile, "access", false)
//	if err != nil {
//	    // handle load error
//	}
//	if token == nil {
//	    // no valid token found, need to authenticate
//	}
func (a *ArkKeyring) LoadToken(profile *models.ArkProfile, postfix string, enforceBasicKeyring bool) (*auth.ArkToken, error) {
	a.logger.Info("Trying to load token [%s-%s] of profile [%s]", a.serviceName, postfix, profile.ProfileName)
	kr, err := a.GetKeyring(enforceBasicKeyring)
	if err != nil {
		return nil, err
	}
	tokenData, err := kr.GetPassword(profile.ProfileName, a.serviceName+"-"+postfix)
	if err != nil {
		if !enforceBasicKeyring {
			a.logger.Warning("Falling back to basic keyring as we failed to load token with keyring [%v]", kr)
			return a.LoadToken(profile, postfix, true)
		}
		a.logger.Warning("Failed to load cached token [%v]", err)
		return nil, err
	}
	if tokenData == "" {
		a.logger.Info("No token found")
		return nil, nil
	}
	var token auth.ArkToken
	if err := json.Unmarshal([]byte(tokenData), &token); err != nil {
		a.logger.Info("Token failed to be parsed [%v]", err)
		return nil, err
	}
	if !time.Time(token.ExpiresIn).IsZero() {
		if token.RefreshToken == "" && token.TokenType != auth.Internal && time.Time(token.ExpiresIn).Before(time.Now().Add(-DefaultExpirationGraceDeltaSeconds*time.Second)) {
			a.logger.Info("Token is expired and no refresh token exists")
			err := kr.DeletePassword(profile.ProfileName, a.serviceName+"-"+postfix)
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
		if token.RefreshToken != "" && time.Time(token.ExpiresIn).Add(MaxKeyringRecordTimeHours*time.Hour).Before(time.Now()) {
			a.logger.Info("Token is expired and has been in the cache for too long before another usage")
			err := kr.DeletePassword(profile.ProfileName, a.serviceName+"-"+postfix)
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
	}
	a.logger.Info("Loaded token successfully")
	return &token, nil
}
