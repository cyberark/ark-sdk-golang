package common

import (
	"encoding/json"
	"github.com/cyberark/ark-sdk-golang/internal/common"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"os"
	"runtime"
	"strings"
	"time"
)

// Env vars and definitions
const (
	ArkBasicKeyringOverrideEnvVar      = "ARK_BASIC_KEYRING"
	DBusSessionEnvVar                  = "DBUS_SESSION_BUS_ADDRESS"
	DefaultExpirationGraceDeltaSeconds = 60
	MaxKeyringRecordTimeHours          = 12
)

// ArkKeyring is a struct that represents a keyring for storing and retrieving tokens.
type ArkKeyring struct {
	serviceName string
	logger      *ArkLogger
}

// NewArkKeyring creates a new instance of ArkKeyring.
func NewArkKeyring(serviceName string) *ArkKeyring {
	return &ArkKeyring{
		serviceName: serviceName,
		logger:      GetLogger("ArkKeyring", Unknown),
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
func (a *ArkKeyring) GetKeyring(enforceBasicKeyring bool) (*common.BasicKeyring, error) {
	if a.isDocker() || a.isWSL() || os.Getenv(ArkBasicKeyringOverrideEnvVar) != "" || enforceBasicKeyring {
		return common.NewBasicKeyring(), nil
	}
	if runtime.GOOS == "windows" {
		// TODO - Implement Windows-specific keyring logic here
		return common.NewBasicKeyring(), nil
	}
	if os.Getenv(DBusSessionEnvVar) == "" {
		return common.NewBasicKeyring(), nil
	}
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		// TODO - Implement SecretService keyring logic here
		return common.NewBasicKeyring(), nil
	}
	return common.NewBasicKeyring(), nil
}

// SaveToken saves the token to the keyring for the specified profile and postfix.
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
	if err := kr.SetPassword(a.serviceName+"-"+postfix, profile.ProfileName, string(tokenData)); err != nil {
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

// LoadToken loads the token from the keyring for the specified profile and postfix.
func (a *ArkKeyring) LoadToken(profile *models.ArkProfile, postfix string, enforceBasicKeyring bool) (*auth.ArkToken, error) {
	a.logger.Info("Trying to load token [%s-%s] of profile [%s]", a.serviceName, postfix, profile.ProfileName)
	kr, err := a.GetKeyring(enforceBasicKeyring)
	if err != nil {
		return nil, err
	}
	tokenData, err := kr.GetPassword(a.serviceName+"-"+postfix, profile.ProfileName)
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
			err := kr.DeletePassword(a.serviceName+"-"+postfix, profile.ProfileName)
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
		if token.RefreshToken != "" && time.Time(token.ExpiresIn).Add(MaxKeyringRecordTimeHours*time.Hour).Before(time.Now()) {
			a.logger.Info("Token is expired and has been in the cache for too long before another usage")
			err := kr.DeletePassword(a.serviceName+"-"+postfix, profile.ProfileName)
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
	}
	a.logger.Info("Loaded token successfully")
	return &token, nil
}
