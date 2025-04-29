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
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	nonceSize = 16
	tagSize   = 16
	blockSize = 32
)

// Variables for basic keyring
const (
	DefaultBasicKeyringFolder   = ".ark_cache/keyring"
	ArkBasicKeyringFolderEnvVar = "ARK_KEYRING_FOLDER"
)

// BasicKeyring is a simple keyring implementation that uses AES encryption to store passwords.
type BasicKeyring struct {
	basicFolderPath string
	keyringFilePath string
	macFilePath     string
}

// NewBasicKeyring creates a new BasicKeyring instance. It initializes the keyring folder and file paths.
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
	data, err := ioutil.ReadFile(b.macFilePath)
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
	data, err := ioutil.ReadFile(b.keyringFilePath)
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
	if err := ioutil.WriteFile(b.keyringFilePath, data, 0644); err != nil {
		return err
	}
	return b.updateMac()
}

// GetPassword retrieves a password for a given service and username from the keyring.
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
