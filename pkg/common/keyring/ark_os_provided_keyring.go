// Package keyring provides keyring utilities for the ARK SDK.
//
// This package contains internal implementations including a OS based keyring
// system for secure password storage. The ArkOSProvidedKeyring
// provides OS based storage.
package keyring

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/99designs/keyring"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
)

const (
	keyringPrefix = "ark_sdk_golang"
)

// OSKeyringOpen is a variable that points to the keyring.Open function from the go-keyring package.
var OSKeyringOpen = keyring.Open

// ArkOSProvidedKeyring provides an OS-based keyring implementation for secure password storage.
//
// ArkOSProvidedKeyring wraps the underlying OS keyring and falls back to a basic keyring
// implementation if the OS keyring is unavailable or encounters errors. It is intended for
// use in environments where secure credential storage is required.
//
// Fields:
//   - ArkKeyringImpl: Embedded base implementation for keyring operations.
//   - fallbackKeyring: The fallback keyring implementation used if OS keyring fails.
//   - logger: Logger instance for warning and error messages.
type ArkOSProvidedKeyring struct {
	ArkKeyringImpl
	fallbackKeyring ArkKeyringImpl
	logger          *common.ArkLogger
}

func (b *ArkOSProvidedKeyring) keyringForService(serviceName string) (keyring.Keyring, error) {
	cfg := keyring.Config{
		ServiceName: serviceName,
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,      // macOS Keychain
			keyring.WinCredBackend,       // Windows Credential Manager
			keyring.SecretServiceBackend, // Linux (libsecret/gnome-keyring)
			keyring.KWalletBackend,       // Linux KDE
			keyring.PassBackend,          // Linux 'pass'
		},
	}
	switch runtime.GOOS {
	case "darwin":
		cfg.KeychainName = "login"
		cfg.KeychainTrustApplication = true
		cfg.KeychainSynchronizable = false
	case "linux":
		cfg.LibSecretCollectionName = "default"
	}
	keyringStore, err := OSKeyringOpen(cfg)
	if err != nil {
		return nil, err
	}
	return keyringStore, nil
}

// NewArkOSProvidedKeyring creates a new ArkOSProvidedKeyring instance.
//
// NewArkOSProvidedKeyring initializes the OS-based keyring with a fallback keyring implementation.
// The fallbackKeyring is used if the OS keyring is unavailable or encounters errors.
//
// Parameters:
//   - fallbackKeyring: ArkKeyringImpl. The fallback keyring implementation to use.
//
// Returns a pointer to the initialized ArkOSProvidedKeyring.
//
// Example:
//
//	fallback := NewBasicKeyring()
//	keyring := NewArkOSProvidedKeyring(fallback)
func NewArkOSProvidedKeyring(fallbackKeyring ArkKeyringImpl) *ArkOSProvidedKeyring {
	return &ArkOSProvidedKeyring{
		fallbackKeyring: fallbackKeyring,
		logger:          common.GetLogger("ArkOSProvidedKeyring", common.Unknown),
	}
}

// SetPassword stores a password in the OS keyring for the specified service and username.
//
// SetPassword attempts to securely store the password using the OS keyring. If the OS keyring
// operation fails, it logs a warning and falls back to the provided basic keyring implementation.
//
// Parameters:
//   - serviceName: string. The name of the service for which the password is stored.
//   - username: string. The username associated with the password.
//   - password: string. The password to store.
//
// Returns an error if the password could not be stored in either the OS keyring or the fallback keyring.
//
// Example:
//
//	err := keyring.SetPassword("my_service", "user1", "secret")
//	if err != nil {
//	    // handle error
//	}
func (b *ArkOSProvidedKeyring) SetPassword(serviceName string, username string, password string) error {
	keyringStore, err := b.keyringForService(serviceName)
	if err != nil {
		b.logger.Warning("Failed to open OS keyring: %v. Falling back to basic keyring.", err)
		return b.fallbackKeyring.SetPassword(serviceName, username, password)
	}
	err = keyringStore.Set(keyring.Item{
		Key:  fmt.Sprintf("%s_%s", keyringPrefix, username),
		Data: []byte(password),
	})
	if err != nil {
		b.logger.Warning("Failed to set password in OS keyring: %v. Falling back to basic keyring.", err)
		return b.fallbackKeyring.SetPassword(serviceName, username, password)
	}
	return nil
}

// GetPassword retrieves a password from the OS keyring for the specified service and username.
//
// GetPassword attempts to retrieve the password using the OS keyring. If the OS keyring
// operation fails, it logs a warning and falls back to the provided basic keyring implementation.
//
// Parameters:
//   - serviceName: string. The name of the service for which the password is retrieved.
//   - username: string. The username associated with the password.
//
// Returns the password as a string and an error if retrieval fails from both the OS keyring and the fallback keyring.
//
// Example:
//
//	password, err := keyring.GetPassword("my_service", "user1")
//	if err != nil {
//	    // handle error
//	}
func (b *ArkOSProvidedKeyring) GetPassword(serviceName string, username string) (string, error) {
	keyringStore, err := b.keyringForService(serviceName)
	if err != nil {
		b.logger.Warning("Failed to open OS keyring: %v. Falling back to basic keyring.", err)
		return b.fallbackKeyring.GetPassword(serviceName, username)
	}
	item, err := keyringStore.Get(fmt.Sprintf("%s_%s", keyringPrefix, username))
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return "", nil
		}
		b.logger.Warning("Failed to get password from OS keyring: %v. Falling back to basic keyring.", err)
		return b.fallbackKeyring.GetPassword(serviceName, username)
	}
	return string(item.Data), nil
}

// DeletePassword removes a password from the OS keyring for the specified service and username.
//
// DeletePassword attempts to delete the password using the OS keyring. If the OS keyring
// operation fails, it logs a warning and falls back to the provided basic keyring implementation.
//
// Parameters:
//   - serviceName: string. The name of the service for which the password is deleted.
//   - username: string. The username associated with the password.
//
// Returns an error if the password could not be deleted from either the OS keyring or the fallback keyring.
//
// Example:
//
//	err := keyring.DeletePassword("my_service", "user1")
//	if err != nil {
//	    // handle error
//	}
func (b *ArkOSProvidedKeyring) DeletePassword(serviceName string, username string) error {
	keyringStore, err := b.keyringForService(serviceName)
	if err != nil {
		b.logger.Warning("Failed to open OS keyring: %v. Falling back to basic keyring.", err)
		return b.fallbackKeyring.DeletePassword(serviceName, username)
	}
	err = keyringStore.Remove(fmt.Sprintf("%s_%s", keyringPrefix, username))
	if err != nil {
		b.logger.Warning("Failed to delete password from OS keyring: %v. Falling back to basic keyring.", err)
		return b.fallbackKeyring.DeletePassword(serviceName, username)
	}
	return nil
}

// ClearAllPasswords removes all stored passwords from the OS keyring and the fallback keyring.
//
// ClearAllPasswords attempts to remove all password entries from the OS keyring that match the ARK SDK prefix.
// If the OS keyring cannot be opened or keys cannot be listed, it logs a warning and falls back to the basic keyring's
// ClearAllPasswords method. The function is idempotent and will not return an error if no entries exist.
//
// Returns an error if removal fails for either the OS keyring or the fallback keyring.
//
// Example:
//
//	err := keyring.ClearAllPasswords()
//	if err != nil {
//	    // Handle error
//	}
func (b *ArkOSProvidedKeyring) ClearAllPasswords() error {
	keyringStore, err := b.keyringForService("")
	if err != nil {
		b.logger.Warning("Failed to open OS keyring: %v. Falling back to basic keyring.", err)
		return b.fallbackKeyring.ClearAllPasswords()
	}
	items, err := keyringStore.Keys()
	if err != nil {
		b.logger.Warning("Failed to list keys from OS keyring: %v. Falling back to basic keyring.", err)
		return b.fallbackKeyring.ClearAllPasswords()
	}
	for _, key := range items {
		if strings.HasPrefix(key, keyringPrefix) {
			_ = keyringStore.Remove(key)
		}
	}
	return nil
}
