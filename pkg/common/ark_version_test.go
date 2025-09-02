package common

import (
	"testing"
)

func TestSetArkVersion(t *testing.T) {
	tests := []struct {
		name         string
		version      string
		setupMock    func()
		validateFunc func(t *testing.T, version string)
	}{
		{
			name:    "success_set_valid_version",
			version: "1.2.3",
			setupMock: func() {
				SetArkVersion("0.0.0") // Reset to default
			},
			validateFunc: func(t *testing.T, version string) {
				SetArkVersion(version)
				result := ArkVersion()
				if result != "1.2.3" {
					t.Errorf("Expected version '1.2.3', got '%s'", result)
				}
			},
		},
		{
			name:    "success_set_semantic_version",
			version: "2.1.0-beta.1",
			setupMock: func() {
				SetArkVersion("1.0.0")
			},
			validateFunc: func(t *testing.T, version string) {
				SetArkVersion(version)
				result := ArkVersion()
				if result != "2.1.0-beta.1" {
					t.Errorf("Expected version '2.1.0-beta.1', got '%s'", result)
				}
			},
		},
		{
			name:    "success_ignore_empty_version",
			version: "",
			setupMock: func() {
				SetArkVersion("1.5.0") // Set initial version
			},
			validateFunc: func(t *testing.T, version string) {
				originalVersion := ArkVersion()
				SetArkVersion(version) // Should be ignored
				result := ArkVersion()
				if result != originalVersion {
					t.Errorf("Expected version to remain '%s', got '%s'", originalVersion, result)
				}
			},
		},
		{
			name:    "success_set_version_with_build_metadata",
			version: "1.0.0+20230815.abcd123",
			setupMock: func() {
				SetArkVersion("0.0.0")
			},
			validateFunc: func(t *testing.T, version string) {
				SetArkVersion(version)
				result := ArkVersion()
				if result != "1.0.0+20230815.abcd123" {
					t.Errorf("Expected version '1.0.0+20230815.abcd123', got '%s'", result)
				}
			},
		},
		{
			name:    "success_override_existing_version",
			version: "3.0.0",
			setupMock: func() {
				SetArkVersion("2.5.1") // Set initial version
			},
			validateFunc: func(t *testing.T, version string) {
				SetArkVersion(version)
				result := ArkVersion()
				if result != "3.0.0" {
					t.Errorf("Expected version '3.0.0', got '%s'", result)
				}
			},
		},
		{
			name:    "success_set_development_version",
			version: "dev-snapshot",
			setupMock: func() {
				SetArkVersion("1.0.0")
			},
			validateFunc: func(t *testing.T, version string) {
				SetArkVersion(version)
				result := ArkVersion()
				if result != "dev-snapshot" {
					t.Errorf("Expected version 'dev-snapshot', got '%s'", result)
				}
			},
		},
		{
			name:    "success_set_version_with_special_characters",
			version: "v1.2.3-rc.1+build.456",
			setupMock: func() {
				SetArkVersion("0.0.0")
			},
			validateFunc: func(t *testing.T, version string) {
				SetArkVersion(version)
				result := ArkVersion()
				if result != "v1.2.3-rc.1+build.456" {
					t.Errorf("Expected version 'v1.2.3-rc.1+build.456', got '%s'", result)
				}
			},
		},
		{
			name:    "edge_case_whitespace_only_version",
			version: "   ",
			setupMock: func() {
				SetArkVersion("1.0.0")
			},
			validateFunc: func(t *testing.T, version string) {
				SetArkVersion(version)
				result := ArkVersion()
				if result != "   " {
					t.Errorf("Expected version '   ', got '%s'", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.validateFunc(t, tt.version)
		})
	}
}

func TestArkVersion(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func()
		expectedResult string
	}{
		{
			name: "success_return_default_version",
			setupMock: func() {
				SetArkVersion("0.0.0") // Ensure default state
			},
			expectedResult: "0.0.0",
		},
		{
			name: "success_return_set_version",
			setupMock: func() {
				SetArkVersion("1.2.3")
			},
			expectedResult: "1.2.3",
		},
		{
			name: "success_return_semantic_version",
			setupMock: func() {
				SetArkVersion("2.1.0-alpha.1")
			},
			expectedResult: "2.1.0-alpha.1",
		},
		{
			name: "success_return_version_with_build_metadata",
			setupMock: func() {
				SetArkVersion("1.0.0+build.123")
			},
			expectedResult: "1.0.0+build.123",
		},
		{
			name: "success_return_development_version",
			setupMock: func() {
				SetArkVersion("dev-latest")
			},
			expectedResult: "dev-latest",
		},
		{
			name: "success_return_complex_version",
			setupMock: func() {
				SetArkVersion("v2.0.0-rc.1+exp.sha.5114f85")
			},
			expectedResult: "v2.0.0-rc.1+exp.sha.5114f85",
		},
		{
			name: "success_return_version_after_multiple_sets",
			setupMock: func() {
				SetArkVersion("1.0.0")
				SetArkVersion("2.0.0")
				SetArkVersion("") // Should be ignored
				SetArkVersion("3.0.0")
			},
			expectedResult: "3.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()
			result := ArkVersion()

			if result != tt.expectedResult {
				t.Errorf("Expected ArkVersion() '%s', got '%s'", tt.expectedResult, result)
			}
		})
	}
}

// TestArkVersion_EmptyStringBehavior tests the specific behavior with empty strings
func TestArkVersion_EmptyStringBehavior(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func()
		operation      func()
		expectedResult string
	}{
		{
			name: "success_empty_string_preserves_previous_version",
			setupMock: func() {
				SetArkVersion("1.5.0")
			},
			operation: func() {
				SetArkVersion("") // Should not change version
			},
			expectedResult: "1.5.0",
		},
		{
			name: "success_multiple_empty_strings_preserve_version",
			setupMock: func() {
				SetArkVersion("2.0.0")
			},
			operation: func() {
				SetArkVersion("")
				SetArkVersion("")
				SetArkVersion("")
			},
			expectedResult: "2.0.0",
		},
		{
			name: "success_empty_string_then_valid_version",
			setupMock: func() {
				SetArkVersion("1.0.0")
			},
			operation: func() {
				SetArkVersion("")      // Should be ignored
				SetArkVersion("2.0.0") // Should update
			},
			expectedResult: "2.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.operation()

			result := ArkVersion()
			if result != tt.expectedResult {
				t.Errorf("Expected version '%s', got '%s'", tt.expectedResult, result)
			}
		})
	}
}

// TestArkVersion_Integration tests the complete interaction between SetArkVersion and ArkVersion
func TestArkVersion_Integration(t *testing.T) {
	tests := []struct {
		name          string
		operations    []string
		expectedFinal string
	}{
		{
			name:          "complete_cycle_multiple_version_updates",
			operations:    []string{"1.0.0", "1.1.0", "", "1.2.0", "2.0.0-beta", ""},
			expectedFinal: "2.0.0-beta",
		},
		{
			name:          "complete_cycle_with_special_versions",
			operations:    []string{"dev", "1.0.0-alpha", "1.0.0-beta", "1.0.0-rc.1", "1.0.0"},
			expectedFinal: "1.0.0",
		},
		{
			name:          "complete_cycle_build_metadata_versions",
			operations:    []string{"1.0.0+build.1", "1.0.0+build.2", "", "1.0.0+build.3"},
			expectedFinal: "1.0.0+build.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to known state
			SetArkVersion("0.0.0")

			// Apply all operations
			for _, version := range tt.operations {
				SetArkVersion(version)
			}

			// Verify final result
			result := ArkVersion()
			if result != tt.expectedFinal {
				t.Errorf("Expected final version '%s', got '%s'", tt.expectedFinal, result)
			}
		})
	}
}

// TestArkVersion_ThreadSafety tests basic concurrent access patterns
func TestArkVersion_ConcurrentAccess(t *testing.T) {
	tests := []struct {
		name         string
		setupVersion string
		validateFunc func(t *testing.T)
	}{
		{
			name:         "success_concurrent_reads",
			setupVersion: "1.0.0",
			validateFunc: func(t *testing.T) {
				// Set initial version
				SetArkVersion("1.0.0")

				// Multiple concurrent reads should all return the same value
				done := make(chan string, 5)

				for i := 0; i < 5; i++ {
					go func() {
						done <- ArkVersion()
					}()
				}

				// Collect all results
				for i := 0; i < 5; i++ {
					result := <-done
					if result != "1.0.0" {
						t.Errorf("Expected version '1.0.0', got '%s'", result)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validateFunc(t)
		})
	}
}
