package testutils

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
)

// MockProfileLoader is a shared mock implementation of the profiles.ProfileLoader interface.
//
// MockProfileLoader provides configurable mock behavior for all ProfileLoader methods
// through function fields. This allows tests to customize behavior for specific
// test scenarios while maintaining a consistent mock implementation across all tests.
//
// Each method checks if a corresponding function field is set and calls it,
// otherwise returns sensible defaults (nil, nil for most cases).
type MockProfileLoader struct {
	LoadProfileFunc        func(string) (*models.ArkProfile, error)
	SaveProfileFunc        func(*models.ArkProfile) error
	LoadAllProfilesFunc    func() ([]*models.ArkProfile, error)
	LoadDefaultProfileFunc func() (*models.ArkProfile, error)
	DeleteProfileFunc      func(string) error
	ClearAllProfilesFunc   func() error
	ProfileExistsFunc      func(string) bool
}

// LoadProfile loads a profile by name using the configured mock function or returns defaults.
func (m *MockProfileLoader) LoadProfile(name string) (*models.ArkProfile, error) {
	if m.LoadProfileFunc != nil {
		return m.LoadProfileFunc(name)
	}
	return nil, nil
}

// SaveProfile saves a profile using the configured mock function or returns nil.
func (m *MockProfileLoader) SaveProfile(profile *models.ArkProfile) error {
	if m.SaveProfileFunc != nil {
		return m.SaveProfileFunc(profile)
	}
	return nil
}

// LoadAllProfiles loads all profiles using the configured mock function or returns empty slice.
func (m *MockProfileLoader) LoadAllProfiles() ([]*models.ArkProfile, error) {
	if m.LoadAllProfilesFunc != nil {
		return m.LoadAllProfilesFunc()
	}
	return nil, nil
}

// LoadDefaultProfile loads the default profile using the configured mock function or returns nil.
func (m *MockProfileLoader) LoadDefaultProfile() (*models.ArkProfile, error) {
	if m.LoadDefaultProfileFunc != nil {
		return m.LoadDefaultProfileFunc()
	}
	return nil, nil
}

// DeleteProfile deletes a profile using the configured mock function or returns nil.
func (m *MockProfileLoader) DeleteProfile(profileName string) error {
	if m.DeleteProfileFunc != nil {
		return m.DeleteProfileFunc(profileName)
	}
	return nil
}

// ClearAllProfiles clears all profiles using the configured mock function or returns nil.
func (m *MockProfileLoader) ClearAllProfiles() error {
	if m.ClearAllProfilesFunc != nil {
		return m.ClearAllProfilesFunc()
	}
	return nil
}

// ProfileExists checks if a profile exists using the configured mock function or returns false.
func (m *MockProfileLoader) ProfileExists(profileName string) bool {
	if m.ProfileExistsFunc != nil {
		return m.ProfileExistsFunc(profileName)
	}
	return false
}

// NewMockProfileLoader creates a new MockProfileLoader instance.
//
// NewMockProfileLoader returns a MockProfileLoader with all function fields
// set to nil, providing default behavior. Tests can customize behavior by
// setting the appropriate function fields.
//
// Returns a pointer to a MockProfileLoader that implements profiles.ProfileLoader.
//
// Example:
//
//	mock := NewMockProfileLoader()
//	mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
//	    return &models.ArkProfile{ProfileName: name}, nil
//	}
//	var loader profiles.ProfileLoader = mock
func NewMockProfileLoader() *MockProfileLoader {
	return &MockProfileLoader{}
}

// AsProfileLoader returns the mock as a profiles.ProfileLoader interface.
//
// AsProfileLoader provides a convenient way to get the MockProfileLoader
// as a ProfileLoader interface pointer, which is commonly needed in test
// setup functions.
//
// Returns a pointer to the ProfileLoader interface.
//
// Example:
//
//	mock := NewMockProfileLoader()
//	loader := mock.AsProfileLoader()
//	action := NewArkProfilesAction(loader)
func (m *MockProfileLoader) AsProfileLoader() *profiles.ProfileLoader {
	var loader profiles.ProfileLoader = m
	return &loader
}
