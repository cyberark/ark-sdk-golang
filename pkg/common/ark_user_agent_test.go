package common

import (
	"strings"
	"testing"
)

func TestUserAgent(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func() (cleanup func())
		validateFunc func(t *testing.T, result string)
	}{
		{
			name: "success_user_agent_contains_chrome_and_sdk_version",
			setupMock: func() (cleanup func()) {
				// Store original version
				originalVersion := ArkVersion()
				SetArkVersion("1.0.0")
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, result string) {
				if result == "" {
					t.Error("Expected non-empty user agent string")
				}

				if !strings.Contains(result, "Ark-SDK-Golang/1.0.0") {
					t.Errorf("Expected user agent to contain 'Ark-SDK-Golang/1.0.0', got '%s'", result)
				}

				// Should contain some browser-like elements
				if !strings.Contains(result, "Mozilla") && !strings.Contains(result, "Chrome") && !strings.Contains(result, "Safari") {
					t.Errorf("Expected user agent to contain browser-like elements, got '%s'", result)
				}
			},
		},
		{
			name: "success_user_agent_with_different_version",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				SetArkVersion("2.5.10")
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, result string) {
				if !strings.Contains(result, "Ark-SDK-Golang/2.5.10") {
					t.Errorf("Expected user agent to contain 'Ark-SDK-Golang/2.5.10', got '%s'", result)
				}
			},
		},
		{
			name: "success_user_agent_with_default_version",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				SetArkVersion("0.0.0")
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, result string) {
				if !strings.Contains(result, "Ark-SDK-Golang/0.0.0") {
					t.Errorf("Expected user agent to contain 'Ark-SDK-Golang/0.0.0', got '%s'", result)
				}
			},
		},
		{
			name: "success_user_agent_with_beta_version",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				SetArkVersion("1.0.0-beta.1")
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, result string) {
				if !strings.Contains(result, "Ark-SDK-Golang/1.0.0-beta.1") {
					t.Errorf("Expected user agent to contain 'Ark-SDK-Golang/1.0.0-beta.1', got '%s'", result)
				}
			},
		},
		{
			name: "success_user_agent_format_consistency",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				SetArkVersion("3.1.4")
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, result string) {
				// Verify the format is consistent: should end with " Ark-SDK-Golang/{version}"
				expectedSuffix := " Ark-SDK-Golang/3.1.4"
				if !strings.HasSuffix(result, expectedSuffix) {
					t.Errorf("Expected user agent to end with '%s', got '%s'", expectedSuffix, result)
				}

				// Should have content before the SDK part
				if len(result) <= len(expectedSuffix) {
					t.Errorf("Expected user agent to have browser content before SDK suffix, got '%s'", result)
				}
			},
		},
		{
			name: "success_user_agent_contains_space_separator",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				SetArkVersion("1.2.3")
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, result string) {
				// Verify there's a space before "Ark-SDK-Golang"
				if !strings.Contains(result, " Ark-SDK-Golang/") {
					t.Errorf("Expected user agent to contain ' Ark-SDK-Golang/', got '%s'", result)
				}
			},
		},
		{
			name: "success_user_agent_with_empty_version",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				SetArkVersion("")
				// Since SetArkVersion ignores empty strings, version should remain unchanged
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, result string) {
				// Should still contain the Ark-SDK-Golang part with whatever version is set
				if !strings.Contains(result, "Ark-SDK-Golang/") {
					t.Errorf("Expected user agent to contain 'Ark-SDK-Golang/', got '%s'", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			result := UserAgent()
			tt.validateFunc(t, result)
		})
	}
}

// TestUserAgent_MultipleCallsConsistency tests that multiple calls return consistent results
func TestUserAgent_MultipleCallsConsistency(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func() (cleanup func())
		validateFunc func(t *testing.T, results []string)
	}{
		{
			name: "success_multiple_calls_same_version_consistent_format",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				SetArkVersion("1.0.0")
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, results []string) {
				if len(results) < 2 {
					t.Fatal("Expected at least 2 results for consistency test")
				}

				// All results should end with the same SDK version part
				expectedSuffix := " Ark-SDK-Golang/1.0.0"
				for i, result := range results {
					if !strings.HasSuffix(result, expectedSuffix) {
						t.Errorf("Result %d: expected suffix '%s', got '%s'", i, expectedSuffix, result)
					}
				}

				// All results should be non-empty
				for i, result := range results {
					if result == "" {
						t.Errorf("Result %d: expected non-empty string, got empty", i)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			var results []string

			if tt.name == "success_version_change_affects_user_agent" {
				// Special case: test with version changes
				SetArkVersion("1.0.0")
				results = append(results, UserAgent())

				SetArkVersion("2.0.0")
				results = append(results, UserAgent())
			} else {
				// Normal case: multiple calls with same version
				for i := 0; i < 3; i++ {
					results = append(results, UserAgent())
				}
			}

			tt.validateFunc(t, results)
		})
	}
}

// TestUserAgent_BrowserComponent tests the browser component behavior
func TestUserAgent_BrowserComponent(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func() (cleanup func())
		validateFunc func(t *testing.T, result string)
	}{
		{
			name: "success_browser_component_not_empty",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				SetArkVersion("1.0.0")
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, result string) {
				// Split by the SDK part to get browser component
				parts := strings.Split(result, " Ark-SDK-Golang/")
				if len(parts) != 2 {
					t.Fatalf("Expected user agent to have exactly one ' Ark-SDK-Golang/' separator, got %d parts", len(parts))
				}

				browserPart := parts[0]
				if browserPart == "" {
					t.Error("Expected non-empty browser component")
				}

				// Browser part should contain typical user agent elements
				lowerBrowser := strings.ToLower(browserPart)
				hasValidElements := strings.Contains(lowerBrowser, "mozilla") ||
					strings.Contains(lowerBrowser, "chrome") ||
					strings.Contains(lowerBrowser, "safari") ||
					strings.Contains(lowerBrowser, "webkit")

				if !hasValidElements {
					t.Errorf("Expected browser component to contain typical user agent elements, got '%s'", browserPart)
				}
			},
		},
		{
			name: "success_sdk_component_format",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				SetArkVersion("1.2.3")
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T, result string) {
				// Split by the SDK part
				parts := strings.Split(result, " Ark-SDK-Golang/")
				if len(parts) != 2 {
					t.Fatalf("Expected user agent to have exactly one ' Ark-SDK-Golang/' separator, got %d parts", len(parts))
				}

				sdkVersion := parts[1]
				if sdkVersion != "1.2.3" {
					t.Errorf("Expected SDK version '1.2.3', got '%s'", sdkVersion)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cleanup := tt.setupMock()
			defer cleanup()

			result := UserAgent()
			tt.validateFunc(t, result)
		})
	}
}

// TestUserAgent_Integration tests the complete integration with version management
func TestUserAgent_Integration(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func() (cleanup func())
		validateFunc func(t *testing.T)
	}{
		{
			name: "complete_cycle_version_and_user_agent",
			setupMock: func() (cleanup func()) {
				originalVersion := ArkVersion()
				return func() {
					SetArkVersion(originalVersion)
				}
			},
			validateFunc: func(t *testing.T) {
				// Test complete workflow
				versions := []string{"1.0.0", "2.1.5", "3.0.0-alpha.1"}

				for _, version := range versions {
					SetArkVersion(version)
					userAgent := UserAgent()

					expectedSuffix := " Ark-SDK-Golang/" + version
					if !strings.HasSuffix(userAgent, expectedSuffix) {
						t.Errorf("For version '%s', expected user agent to end with '%s', got '%s'",
							version, expectedSuffix, userAgent)
					}

					if !strings.Contains(userAgent, " Ark-SDK-Golang/") {
						t.Errorf("For version '%s', expected user agent to contain ' Ark-SDK-Golang/', got '%s'",
							version, userAgent)
					}

					// Verify the current version matches
					currentVersion := ArkVersion()
					if currentVersion != version {
						t.Errorf("Expected current version '%s', got '%s'", version, currentVersion)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			tt.validateFunc(t)
		})
	}
}
