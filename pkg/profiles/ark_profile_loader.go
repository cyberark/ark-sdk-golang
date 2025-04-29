package profiles

import (
	"encoding/json"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"os"
	"path/filepath"
)

// ProfileLoader is an interface that defines methods for loading, saving, and managing Ark profiles.
type ProfileLoader interface {
	LoadProfile(profileName string) (*models.ArkProfile, error)
	SaveProfile(profile *models.ArkProfile) error
	LoadAllProfiles() ([]*models.ArkProfile, error)
	LoadDefaultProfile() (*models.ArkProfile, error)
	DeleteProfile(profileName string) error
	ClearAllProfiles() error
	ProfileExists(profileName string) bool
}

// FileSystemProfilesLoader is a struct that implements the ProfileLoader interface using the file system.
type FileSystemProfilesLoader struct {
	ProfileLoader
}

// DefaultProfilesLoader returns a default implementation of the ProfileLoader interface, which is filesystem.
func DefaultProfilesLoader() *ProfileLoader {
	var profilesLoader ProfileLoader = &FileSystemProfilesLoader{}
	return &profilesLoader
}

// GetProfilesFolder returns the folder path where Ark profiles are stored.
func GetProfilesFolder() string {
	if folder := os.Getenv("ARK_PROFILES_FOLDER"); folder != "" {
		return folder
	}
	return filepath.Join(os.Getenv("HOME"), ".ark_profiles")
}

// DefaultProfileName returns the default profile name.
func DefaultProfileName() string {
	return "ark"
}

// DeduceProfileName deduces the profile name based on the provided name and environment variables.
func DeduceProfileName(profileName string) string {
	if profileName != "" && profileName != DefaultProfileName() {
		return profileName
	}
	if profile := os.Getenv("ARK_PROFILE"); profile != "" {
		return profile
	}
	if profileName != "" {
		return profileName
	}
	return DefaultProfileName()
}

// LoadDefaultProfile loads the default profile from the file system.
func (fspl *FileSystemProfilesLoader) LoadDefaultProfile() (*models.ArkProfile, error) {
	folder := GetProfilesFolder()
	profileName := DeduceProfileName("")
	profilePath := filepath.Join(folder, profileName)
	if _, err := os.Stat(profilePath); err == nil {
		return fspl.LoadProfile(profileName)
	}
	return &models.ArkProfile{}, nil
}

// LoadProfile loads a profile from the file system based on the provided profile name.
func (fspl *FileSystemProfilesLoader) LoadProfile(profileName string) (*models.ArkProfile, error) {
	folder := GetProfilesFolder()
	profilePath := filepath.Join(folder, profileName)
	if _, err := os.Stat(profilePath); err == nil {
		data, err := os.ReadFile(profilePath)
		if err != nil {
			return nil, err
		}
		var profile models.ArkProfile
		if err := json.Unmarshal(data, &profile); err != nil {
			return nil, err
		}
		return &profile, nil
	}
	return nil, nil
}

// SaveProfile saves a profile to the file system, will create needed folders if not already present.
func (fspl *FileSystemProfilesLoader) SaveProfile(profile *models.ArkProfile) error {
	folder := GetProfilesFolder()
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return err
		}
	}
	profilePath := filepath.Join(folder, profile.ProfileName)
	data, err := json.MarshalIndent(profile, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(profilePath, data, 0644)
}

// LoadAllProfiles loads all profiles from the file system.
func (fspl *FileSystemProfilesLoader) LoadAllProfiles() ([]*models.ArkProfile, error) {
	folder := GetProfilesFolder()
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return nil, nil
	}
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	var profiles []*models.ArkProfile
	for _, file := range files {
		if !file.IsDir() {
			profile, err := fspl.LoadProfile(file.Name())
			if err != nil {
				continue
			}
			profiles = append(profiles, profile)
		}
	}
	return profiles, nil
}

// DeleteProfile deletes a profile from the file system.
func (fspl *FileSystemProfilesLoader) DeleteProfile(profileName string) error {
	folder := GetProfilesFolder()
	profilePath := filepath.Join(folder, profileName)
	if _, err := os.Stat(profilePath); err == nil {
		return os.Remove(profilePath)
	}
	return nil
}

// ClearAllProfiles clears all profiles from the file system.
func (fspl *FileSystemProfilesLoader) ClearAllProfiles() error {
	folder := GetProfilesFolder()
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return nil
	}
	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			if err := os.Remove(filepath.Join(folder, file.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

// ProfileExists checks if a profile exists in the file system.
func (fspl *FileSystemProfilesLoader) ProfileExists(profileName string) bool {
	folder := GetProfilesFolder()
	profilePath := filepath.Join(folder, profileName)
	_, err := os.Stat(profilePath)
	return err == nil
}
