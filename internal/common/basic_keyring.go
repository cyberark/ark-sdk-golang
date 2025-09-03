// Package common provides internal shared utilities for the ARK SDK.
//
// This package contains internal implementations including a basic keyring
// system for secure password storage using AES encryption. The BasicKeyring
// provides file-based storage with MAC (Message Authentication Code) validation
// to ensure data integrity.
package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	// nonceSize defines the size in bytes for AES-GCM nonce generation.
	nonceSize = 16
	// tagSize defines the size in bytes for AES-GCM authentication tag.
	tagSize = 16
	// blockSize defines the block size in bytes for PKCS7 padding operations.
	blockSize = 32
)

// Keyring configuration constants
const (
	// DefaultBasicKeyringFolder is the default folder path relative to HOME directory
	// where the basic keyring files are stored.
	DefaultBasicKeyringFolder = ".ark_cache/keyring"

	// ArkBasicKeyringFolderEnvVar is the environment variable name that can be used
	// to override the default keyring folder location.
	ArkBasicKeyringFolderEnvVar = "ARK_KEYRING_FOLDER"
)

// BasicKeyring is a simple keyring implementation that uses AES encryption to store passwords.
//
// BasicKeyring provides secure password storage using AES-GCM encryption with the hostname
// as the encryption key. Passwords are stored in a JSON file with MAC (Message Authentication
// Code) validation to ensure data integrity. The keyring supports multiple services and
// usernames within each service.
//
// The encryption key is derived from the system hostname and padded using PKCS7 padding
// to ensure consistent key length. Each password entry is encrypted separately with its
// own nonce for security.
//
// File Structure:
//   - keyring: JSON file containing encrypted password data
//   - mac: File containing SHA256 hash of keyring file for integrity validation
type BasicKeyring struct {
	// basicFolderPath is the absolute path to the keyring folder
	basicFolderPath string
	// keyringFilePath is the absolute path to the keyring data file
	keyringFilePath string
	// macFilePath is the absolute path to the MAC validation file
	macFilePath string
}

// NewBasicKeyring creates a new BasicKeyring instance with initialized folder and file paths.
//
// NewBasicKeyring initializes the keyring folder structure and returns a new BasicKeyring
// instance. The folder location is determined by the ArkBasicKeyringFolderEnvVar environment
// variable, or defaults to DefaultBasicKeyringFolder within the user's HOME directory.
//
// The function automatically creates the keyring folder if it doesn't exist. If folder
// creation fails, the function returns nil.
//
// Returns a new BasicKeyring instance or nil if folder creation fails.
//
// Environment Variables:
//   - ARK_KEYRING_FOLDER: Override default keyring folder location
//   - HOME: Used for default keyring folder path construction
//
// Example:
//
//	keyring := NewBasicKeyring()
//	if keyring == nil {
//	    // Handle keyring initialization failure
//	}
func NewBasicKeyring() *BasicKeyring {
	basicFolderPath := filepath.Join(os.Getenv("HOME"), DefaultBasicKeyringFolder)
	if folder := os.Getenv(ArkBasicKeyringFolderEnvVar); folder != "" {
		basicFolderPath = folder
	}
	if _, err := os.Stat(basicFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(basicFolderPath, os.ModePerm)
		if err != nil {
			return nil
		}
	}
	return &BasicKeyring{
		basicFolderPath: basicFolderPath,
		keyringFilePath: filepath.Join(basicFolderPath, "keyring"),
		macFilePath:     filepath.Join(basicFolderPath, "mac"),
	}
}

func (b *BasicKeyring) encrypt(secret []byte, data string) (map[string]string, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCMWithNonceSize(block, nonceSize)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	ciphertextWithTag := aesGCM.Seal(nil, nonce, []byte(data), nil)
	tag := ciphertextWithTag[len(ciphertextWithTag)-tagSize:]
	ciphertext := ciphertextWithTag[:len(ciphertextWithTag)-tagSize]
	return map[string]string{
		"nonce":      base64.StdEncoding.EncodeToString(nonce),
		"ciphertext": base64.StdEncoding.EncodeToString(ciphertext),
		"tag":        base64.StdEncoding.EncodeToString(tag),
	}, nil
}

func (b *BasicKeyring) decrypt(secret []byte, data map[string]string) (string, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCMWithNonceSize(block, nonceSize)
	if err != nil {
		return "", err
	}
	nonce, err := base64.StdEncoding.DecodeString(data["nonce"])
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(data["ciphertext"])
	if err != nil {
		return "", err
	}
	tag, err := base64.StdEncoding.DecodeString(data["tag"])
	if err != nil {
		return "", err
	}
	fullCiphertext := append(ciphertext, tag...)
	plaintext, err := aesGCM.Open(nil, nonce, fullCiphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func (b *BasicKeyring) getCurrentMac() (string, error) {
	if _, err := os.Stat(b.macFilePath); os.IsNotExist(err) {
		return "", errors.New("invalid keyring")
	}
	data, err := os.ReadFile(b.macFilePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (b *BasicKeyring) validateMacAndGetData() (string, error) {
	mac, err := b.getCurrentMac()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(b.keyringFilePath)
	if err != nil {
		return "", err
	}
	dataMac := sha256.Sum256(data)
	if mac == hex.EncodeToString(dataMac[:]) {
		return string(data), nil
	}
	return "", errors.New("invalid keyring")
}

func (b *BasicKeyring) updateMac() error {
	data, err := os.ReadFile(b.keyringFilePath)
	if err != nil {
		return err
	}
	dataMac := sha256.Sum256(data)
	return os.WriteFile(b.macFilePath, []byte(hex.EncodeToString(dataMac[:])), 0644)
}

func (b *BasicKeyring) pKCS7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// SetPassword sets a password for a given service and username in the keyring.
//
// SetPassword encrypts and stores a password for the specified service and username
// combination. The password is encrypted using AES-GCM with a key derived from the
// system hostname. If a keyring file already exists, it loads the existing data and
// adds the new entry. The function updates the MAC file after successful storage
// to maintain data integrity validation.
//
// Parameters:
//   - serviceName: The name of the service (e.g., "github", "aws")
//   - username: The username for the service
//   - password: The password to encrypt and store
//
// Returns an error if encryption, file operations, or MAC update fails.
//
// Example:
//
//	err := keyring.SetPassword("github", "myuser", "mypassword")
//	if err != nil {
//	    // Handle password storage error
//	}
func (b *BasicKeyring) SetPassword(serviceName, username, password string) error {
	key := make([]byte, blockSize)
	hostname, _ := os.Hostname()
	copy(key, b.pKCS7Pad([]byte(hostname), blockSize))
	existingKeyring := make(map[string]map[string]map[string]string)
	if _, err := os.Stat(b.keyringFilePath); err == nil {
		data, err := b.validateMacAndGetData()
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(data), &existingKeyring); err != nil {
			return err
		}
	}
	if _, ok := existingKeyring[serviceName]; !ok {
		existingKeyring[serviceName] = make(map[string]map[string]string)
	}
	encryptedPassword, err := b.encrypt(key, password)
	if err != nil {
		return err
	}
	existingKeyring[serviceName][username] = encryptedPassword
	data, err := json.Marshal(existingKeyring)
	if err != nil {
		return err
	}
	if err := os.WriteFile(b.keyringFilePath, data, 0644); err != nil {
		return err
	}
	return b.updateMac()
}

// GetPassword retrieves a password for a given service and username from the keyring.
//
// GetPassword decrypts and returns the stored password for the specified service
// and username combination. The function validates the keyring MAC before accessing
// the data to ensure integrity. If the keyring file doesn't exist, the service
// doesn't exist, or the username doesn't exist, an empty string is returned without error.
//
// Parameters:
//   - serviceName: The name of the service to retrieve password for
//   - username: The username to retrieve password for
//
// Returns the decrypted password string and any error encountered during retrieval.
// Returns empty string with nil error if the entry doesn't exist.
//
// Example:
//
//	password, err := keyring.GetPassword("github", "myuser")
//	if err != nil {
//	    // Handle retrieval error
//	}
//	if password == "" {
//	    // Password not found
//	}
func (b *BasicKeyring) GetPassword(serviceName, username string) (string, error) {
	key := make([]byte, blockSize)
	hostname, _ := os.Hostname()
	copy(key, b.pKCS7Pad([]byte(hostname), blockSize))
	if _, err := os.Stat(b.keyringFilePath); os.IsNotExist(err) {
		return "", nil
	}
	data, err := b.validateMacAndGetData()
	if err != nil {
		return "", err
	}
	existingKeyring := make(map[string]map[string]map[string]string)
	if err := json.Unmarshal([]byte(data), &existingKeyring); err != nil {
		return "", err
	}
	if _, ok := existingKeyring[serviceName]; !ok {
		return "", nil
	}
	if _, ok := existingKeyring[serviceName][username]; !ok {
		return "", nil
	}
	return b.decrypt(key, existingKeyring[serviceName][username])
}

// DeletePassword deletes a password for a given service and username from the keyring.
//
// DeletePassword removes the specified password entry from the keyring and updates
// the MAC file to maintain data integrity. If the keyring file doesn't exist, the
// service doesn't exist, or the username doesn't exist, the function returns nil
// without error (idempotent behavior).
//
// Parameters:
//   - serviceName: The name of the service to delete password from
//   - username: The username to delete password for
//
// Returns an error if file operations, JSON marshaling, or MAC update fails.
//
// Example:
//
//	err := keyring.DeletePassword("github", "myuser")
//	if err != nil {
//	    // Handle deletion error
//	}
func (b *BasicKeyring) DeletePassword(serviceName, username string) error {
	if _, err := os.Stat(b.keyringFilePath); os.IsNotExist(err) {
		return nil
	}
	data, err := b.validateMacAndGetData()
	if err != nil {
		return err
	}
	existingKeyring := make(map[string]map[string]map[string]string)
	if err := json.Unmarshal([]byte(data), &existingKeyring); err != nil {
		return err
	}
	if _, ok := existingKeyring[serviceName]; !ok {
		return nil
	}
	if _, ok := existingKeyring[serviceName][username]; !ok {
		return nil
	}
	delete(existingKeyring[serviceName], username)
	dataBytes, err := json.Marshal(existingKeyring)
	if err != nil {
		return err
	}
	if err := os.WriteFile(b.keyringFilePath, dataBytes, 0644); err != nil {
		return err
	}
	return b.updateMac()
}
