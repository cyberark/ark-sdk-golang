// Package profiles provides functionality for managing ARK configuration profiles
// in the ARK SDK Golang. This package implements profile loading, saving, and
// management capabilities using a filesystem-based storage approach.
//
// Profiles are used to store and manage ARK configuration settings including
// authentication details, environment settings, and other configuration parameters.
// The package supports multiple profiles with automatic resolution based on
// environment variables and defaults.
//
// Key features:
//   - Filesystem-based profile storage
//   - Multiple profile management
//   - Environment variable-based profile resolution
//   - JSON serialization for profile data
//   - Automatic directory creation for profile storage
//   - Profile existence checking and validation
//
// Profile Resolution Order:
// 1. Explicitly provided profile name (if not default)
// 2. ARK_PROFILE environment variable
// 3. Provided profile name (if any)
// 4. Default profile name ("ark")
//
// Example:
//
//	loader := DefaultProfilesLoader()
//	profile, err := (*loader).LoadDefaultProfile()
//	if err != nil {
//		// handle error
//	}
//
//	// Save a new profile
//	newProfile := &models.ArkProfile{
//		ProfileName: "production",
//		// ... other settings
//	}
//	err = (*loader).SaveProfile(newProfile)
package profiles

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cyberark/ark-sdk-golang/pkg/models"
)

// ProfileLoader is an interface that defines methods for loading, saving, and managing Ark profiles.
//
// This interface provides a contract for profile management operations including
// loading, saving, deleting, and querying profiles. Implementations should handle
// profile persistence and provide consistent behavior across different storage
// mechanisms.
//
// All profile operations work with profile names as string identifiers and
// return appropriate errors for failure conditions such as missing profiles,
// permission issues, or storage failures.
type ProfileLoader interface {
	// LoadProfile loads a profile by name and returns it or an error if loading fails.
	LoadProfile(profileName string) (*models.ArkProfile, error)
	// SaveProfile saves a profile to persistent storage and returns an error if saving fails.
	SaveProfile(profile *models.ArkProfile) error
	// LoadAllProfiles loads all available profiles and returns them as a slice or an error.
	LoadAllProfiles() ([]*models.ArkProfile, error)
	// LoadDefaultProfile loads the default profile based on environment and defaults.
	LoadDefaultProfile() (*models.ArkProfile, error)
	// DeleteProfile removes a profile by name and returns an error if deletion fails.
	DeleteProfile(profileName string) error
	// ClearAllProfiles removes all profiles and returns an error if the operation fails.
	ClearAllProfiles() error
	// ProfileExists checks if a profile with the given name exists and returns a boolean.
	ProfileExists(profileName string) bool
}

// FileSystemProfilesLoader is a struct that implements the ProfileLoader interface using the file system.
//
// This implementation stores profiles as JSON files in a configurable directory.
// Each profile is stored as a separate file named after the profile name.
// The loader handles directory creation, file permissions, and JSON serialization
// automatically.
//
// Profile files are stored in the directory returned by GetProfilesFolder(),
// which can be customized using the ARK_PROFILES_FOLDER environment variable.
type FileSystemProfilesLoader struct {
	ProfileLoader
}

// DefaultProfilesLoader returns a default implementation of the ProfileLoader interface, which is filesystem.
//
// Creates and returns a filesystem-based ProfileLoader implementation that stores
// profiles as JSON files in the local filesystem. This is the standard
// implementation used throughout the ARK SDK.
//
// Returns a pointer to a ProfileLoader interface that can be used for all
// profile management operations.
//
// Example:
//
//	loader := DefaultProfilesLoader()
//	profile, err := (*loader).LoadProfile("production")
//	if err != nil {
//		// handle error
//	}
func DefaultProfilesLoader() *ProfileLoader {
	var profilesLoader ProfileLoader = &FileSystemProfilesLoader{}
	return &profilesLoader
}

// GetProfilesFolder returns the folder path where Ark profiles are stored.
//
// Determines the directory where profile files should be stored by checking
// the ARK_PROFILES_FOLDER environment variable first, and falling back to
// a default location in the user's home directory.
//
// The resolution order is:
// 1. ARK_PROFILES_FOLDER environment variable (if set and non-empty)
// 2. $HOME/.ark_profiles (default fallback)
//
// Returns the absolute path to the profiles directory as a string.
//
// Example:
//
//	folder := GetProfilesFolder()
//	// Returns "/home/user/.ark_profiles" or custom path from environment
func GetProfilesFolder() string {
	if folder := os.Getenv("ARK_PROFILES_FOLDER"); folder != "" {
		return folder
	}
	return filepath.Join(os.Getenv("HOME"), ".ark_profiles")
}

// DefaultProfileName returns the default profile name.
//
// Returns the standard default profile name used when no specific profile
// is requested or when profile resolution falls back to the default.
//
// Returns "ark" as the default profile name.
//
// Example:
//
//	defaultName := DefaultProfileName()
//	// Returns "ark"
func DefaultProfileName() string {
	return "ark"
}

// DeduceProfileName deduces the profile name based on the provided name and environment variables.
//
// Implements the profile name resolution logic that determines which profile
// to use based on multiple sources in order of precedence. This function
// provides consistent profile name resolution across the ARK SDK.
//
// The resolution order is:
// 1. If profileName is provided and differs from default, use it
// 2. If ARK_PROFILE environment variable is set, use its value
// 3. If profileName is provided (even if it's the default), use it
// 4. Fall back to the default profile name
//
// Parameters:
//   - profileName: The explicitly requested profile name, can be empty
//
// Returns the resolved profile name as a string.
//
// Example:
//
//	// With ARK_PROFILE="production" environment variable
//	name := DeduceProfileName("")
//	// Returns "production"
//
//	// With explicit name
//	name := DeduceProfileName("staging")
//	// Returns "staging"
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
//
// Loads the default profile by first determining the profile name using
// DeduceProfileName with an empty string (which triggers environment variable
// and default resolution), then attempting to load that profile from the
// filesystem.
//
// If the profile file exists, it loads and returns the profile. If the profile
// file doesn't exist, it returns an empty ArkProfile struct rather than an error.
//
// Returns a pointer to the loaded ArkProfile and an error if file operations fail.
// If no default profile exists, returns an empty profile without error.
//
// Example:
//
//	profile, err := loader.LoadDefaultProfile()
//	if err != nil {
//		// handle file system error
//	}
//	// profile contains default profile data or empty struct
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
//
// Attempts to load a profile by constructing the file path from the profiles
// folder and the profile name, then reading and deserializing the JSON content.
// If the profile file exists, it reads and unmarshals the JSON data into an
// ArkProfile struct.
//
// If the profile file doesn't exist, the method returns nil for both the profile
// and error, indicating that no profile was found. This allows callers to
// distinguish between "profile not found" and "error loading profile".
//
// Parameters:
//   - profileName: The name of the profile to load
//
// Returns a pointer to the loaded ArkProfile, or nil if the profile doesn't exist.
// Returns an error if file reading or JSON unmarshaling fails.
//
// Example:
//
//	profile, err := loader.LoadProfile("production")
//	if err != nil {
//		// handle file system or JSON error
//	}
//	if profile == nil {
//		// profile doesn't exist
//	}
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
//
// Saves the provided profile to the filesystem by serializing it to JSON and
// writing it to a file named after the profile's ProfileName. The method
// automatically creates the profiles directory if it doesn't exist.
//
// The profile is saved with indented JSON formatting for better readability.
// File permissions are set to 0644 (readable by owner and group, writable by owner).
//
// Parameters:
//   - profile: The ArkProfile to save, must have a valid ProfileName
//
// Returns an error if directory creation, JSON marshaling, or file writing fails.
//
// Example:
//
//	profile := &models.ArkProfile{
//		ProfileName: "production",
//		// ... other fields
//	}
//	err := loader.SaveProfile(profile)
//	if err != nil {
//		// handle save error
//	}
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
//
// Scans the profiles directory and attempts to load all profile files found.
// Only regular files are considered (directories are skipped). If a file
// cannot be loaded as a valid profile, it is silently skipped and processing
// continues with other files.
//
// If the profiles directory doesn't exist, returns nil without error.
//
// Returns a slice of pointers to loaded ArkProfile structs, or an error if
// the directory cannot be read. Individual profile loading errors are ignored
// to allow partial success when some profiles are corrupted.
//
// Example:
//
//	profiles, err := loader.LoadAllProfiles()
//	if err != nil {
//		// handle directory read error
//	}
//	for _, profile := range profiles {
//		fmt.Printf("Loaded profile: %s\n", profile.ProfileName)
//	}
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
//
// Removes the profile file with the specified name from the profiles directory.
// If the profile file exists, it is deleted. If the file doesn't exist,
// the method returns without error (idempotent behavior).
//
// Parameters:
//   - profileName: The name of the profile to delete
//
// Returns an error if the file exists but cannot be deleted due to permissions
// or other filesystem issues.
//
// Example:
//
//	err := loader.DeleteProfile("old-profile")
//	if err != nil {
//		// handle deletion error
//	}
func (fspl *FileSystemProfilesLoader) DeleteProfile(profileName string) error {
	folder := GetProfilesFolder()
	profilePath := filepath.Join(folder, profileName)
	if _, err := os.Stat(profilePath); err == nil {
		return os.Remove(profilePath)
	}
	return nil
}

// ClearAllProfiles clears all profiles from the file system.
//
// Removes all profile files from the profiles directory. Only regular files
// are deleted; subdirectories are left intact. If the profiles directory
// doesn't exist, the method returns without error.
//
// The operation stops and returns an error if any file cannot be deleted,
// which may leave the directory in a partially cleared state.
//
// Returns an error if the directory cannot be read or if any file deletion fails.
//
// Example:
//
//	err := loader.ClearAllProfiles()
//	if err != nil {
//		// handle clear operation error
//	}
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
//
// Checks for the existence of a profile file with the specified name in the
// profiles directory. Uses filesystem stat operations to determine if the
// file exists and is accessible.
//
// Parameters:
//   - profileName: The name of the profile to check for existence
//
// Returns true if the profile file exists and is accessible, false otherwise.
// Returns false for any stat errors including permission denied, file not found,
// or other filesystem issues.
//
// Example:
//
//	if loader.ProfileExists("production") {
//		// profile file exists
//	} else {
//		// profile file doesn't exist or is not accessible
//	}
func (fspl *FileSystemProfilesLoader) ProfileExists(profileName string) bool {
	folder := GetProfilesFolder()
	profilePath := filepath.Join(folder, profileName)
	_, err := os.Stat(profilePath)
	return err == nil
}
