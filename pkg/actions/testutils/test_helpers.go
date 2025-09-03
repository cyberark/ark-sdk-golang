package testutils

import (
	"fmt"

	"github.com/cyberark/ark-sdk-golang/pkg/models"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/spf13/cobra"
)

// CreateTestProfile creates a basic test profile with the given name.
//
// CreateTestProfile generates a standardized ArkProfile instance for testing
// purposes. The profile includes ISP authentication configuration which is
// commonly used in login action tests.
//
// Parameters:
//   - name: The profile name to set
//
// Returns a pointer to a fully configured ArkProfile for testing.
//
// Example:
//
//	profile := CreateTestProfile("test-profile")
//	// profile.ProfileName == "test-profile"
func CreateTestProfile(name string) *models.ArkProfile {
	return &models.ArkProfile{
		ProfileName:        name,
		ProfileDescription: "Test profile for " + name,
		AuthProfiles: map[string]*authmodels.ArkAuthProfile{
			"isp": {
				Username:   "testuser",
				AuthMethod: authmodels.Identity,
				AuthMethodSettings: &authmodels.IdentityArkAuthMethodSettings{
					IdentityURL:             "https://identity.example.com",
					IdentityTenantSubdomain: "test",
				},
			},
		},
	}
}

// CreateTestProfileWithAuth creates a test profile with specific authentication configuration.
//
// CreateTestProfileWithAuth generates an ArkProfile with customized authentication
// settings, allowing tests to specify the auth profile name and authentication method.
//
// Parameters:
//   - name: The profile name to set
//   - authName: The name of the auth profile within the profile
//   - authMethod: The authentication method to use
//
// Returns a pointer to a configured ArkProfile with custom auth settings.
//
// Example:
//
//	profile := CreateTestProfileWithAuth("prod", "production", authmodels.Identity)
//	// profile has auth profile named "production" with Identity method
func CreateTestProfileWithAuth(name, authName string, authMethod authmodels.ArkAuthMethod) *models.ArkProfile {
	return &models.ArkProfile{
		ProfileName:        name,
		ProfileDescription: "Test profile for " + name,
		AuthProfiles: map[string]*authmodels.ArkAuthProfile{
			authName: {
				Username:   "testuser",
				AuthMethod: authMethod,
			},
		},
	}
}

// CreateTestProfiles creates multiple test profiles with sequential naming.
//
// CreateTestProfiles generates a slice of ArkProfile instances with names
// following the pattern "profile1", "profile2", etc. This is useful for
// tests that need multiple profiles.
//
// Parameters:
//   - count: The number of profiles to create
//
// Returns a slice of ArkProfile pointers.
//
// Example:
//
//	profiles := CreateTestProfiles(3)
//	// Returns profiles named "profile1", "profile2", "profile3"
func CreateTestProfiles(count int) []*models.ArkProfile {
	profiles := make([]*models.ArkProfile, count)
	for i := 0; i < count; i++ {
		profiles[i] = CreateTestProfile(fmt.Sprintf("profile%d", i+1))
	}
	return profiles
}

// SetupTestCommand creates a basic cobra command for testing.
//
// SetupTestCommand creates a minimal cobra command with the given name
// that can be used as a parent command in tests. It includes basic
// configuration suitable for most test scenarios.
//
// Parameters:
//   - name: The command name to use
//
// Returns a pointer to a configured cobra.Command.
//
// Example:
//
//	cmd := SetupTestCommand("test")
//	// cmd.Use == "test"
func SetupTestCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "Test command for " + name,
	}
}

// SetupTestCommandWithFlags creates a cobra command with predefined flags.
//
// SetupTestCommandWithFlags creates a cobra command and adds flags based
// on the provided flag definitions. The flagDefs map should have flag names
// as keys and flag configurations as values.
//
// Parameters:
//   - name: The command name to use
//   - flagDefs: Map of flag name to flag configuration
//
// Returns a pointer to a configured cobra.Command with flags.
//
// Note: This is a basic implementation that can be extended based on
// specific flag type requirements.
func SetupTestCommandWithFlags(name string, flagDefs map[string]string) *cobra.Command {
	cmd := SetupTestCommand(name)
	for flagName, flagDesc := range flagDefs {
		cmd.Flags().String(flagName, "", flagDesc)
	}
	return cmd
}

// ProfileLoaderPtr converts a ProfileLoader to a ProfileLoader pointer.
//
// ProfileLoaderPtr is a helper function for creating ProfileLoader pointers
// from ProfileLoader values, which is commonly needed for action constructors.
//
// Parameters:
//   - loader: The ProfileLoader interface to convert to pointer
//
// Returns a pointer to the ProfileLoader interface.
//
// Example:
//
//	mockLoader := &MockProfileLoader{}
//	loaderPtr := ProfileLoaderPtr(mockLoader)
//	action := NewArkLoginAction(loaderPtr)
func ProfileLoaderPtr(loader profiles.ProfileLoader) *profiles.ProfileLoader {
	return &loader
}

// StringPtr creates a pointer to a string value.
//
// StringPtr is a helper function for creating string pointers in tests,
// which is useful when testing functions that accept optional string parameters.
//
// Parameters:
//   - s: The string value to convert to pointer
//
// Returns a pointer to the string value.
//
// Example:
//
//	namePtr := StringPtr("test-name")
//	// Use namePtr where *string is expected
func StringPtr(s string) *string {
	return &s
}

// BoolPtr creates a pointer to a bool value.
//
// BoolPtr is a helper function for creating bool pointers in tests,
// which is useful when testing functions that accept optional bool parameters.
//
// Parameters:
//   - b: The bool value to convert to pointer
//
// Returns a pointer to the bool value.
//
// Example:
//
//	enabledPtr := BoolPtr(true)
//	// Use enabledPtr where *bool is expected
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr creates a pointer to an int value.
//
// IntPtr is a helper function for creating int pointers in tests,
// which is useful when testing functions that accept optional int parameters.
//
// Parameters:
//   - i: The int value to convert to pointer
//
// Returns a pointer to the int value.
//
// Example:
//
//	countPtr := IntPtr(42)
//	// Use countPtr where *int is expected
func IntPtr(i int) *int {
	return &i
}
