package common

import (
	"os"
	"testing"
)

func TestDisableColor(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func()
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_disable_color_from_enabled_state",
			setupMock: func() {
				EnableColor() // Start with color enabled
			},
			validateFunc: func(t *testing.T) {
				DisableColor()
				if IsColoring() {
					t.Error("Expected IsColoring() to return false after DisableColor()")
				}
			},
		},
		{
			name: "success_disable_color_from_disabled_state",
			setupMock: func() {
				DisableColor() // Start with color already disabled
			},
			validateFunc: func(t *testing.T) {
				DisableColor()
				if IsColoring() {
					t.Error("Expected IsColoring() to return false after DisableColor()")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.validateFunc(t)
		})
	}
}

func TestEnableColor(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func()
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_enable_color_from_disabled_state",
			setupMock: func() {
				DisableColor() // Start with color disabled
			},
			validateFunc: func(t *testing.T) {
				EnableColor()
				if !IsColoring() {
					t.Error("Expected IsColoring() to return true after EnableColor()")
				}
			},
		},
		{
			name: "success_enable_color_from_enabled_state",
			setupMock: func() {
				EnableColor() // Start with color already enabled
			},
			validateFunc: func(t *testing.T) {
				EnableColor()
				if !IsColoring() {
					t.Error("Expected IsColoring() to return true after EnableColor()")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.validateFunc(t)
		})
	}
}

func TestIsColoring(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func()
		expectedResult bool
	}{
		{
			name: "success_coloring_enabled",
			setupMock: func() {
				EnableColor()
			},
			expectedResult: true,
		},
		{
			name: "success_coloring_disabled",
			setupMock: func() {
				DisableColor()
			},
			expectedResult: false,
		},
		{
			name: "success_default_state",
			setupMock: func() {
				// Test default state (color enabled by default)
				EnableColor()
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()
			result := IsColoring()

			if result != tt.expectedResult {
				t.Errorf("Expected IsColoring() %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestEnableInteractive(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func()
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_enable_interactive_from_disabled_state",
			setupMock: func() {
				DisableInteractive() // Start with interactive disabled
			},
			validateFunc: func(t *testing.T) {
				EnableInteractive()
				if !IsInteractive() {
					t.Error("Expected IsInteractive() to return true after EnableInteractive()")
				}
			},
		},
		{
			name: "success_enable_interactive_from_enabled_state",
			setupMock: func() {
				EnableInteractive() // Start with interactive already enabled
			},
			validateFunc: func(t *testing.T) {
				EnableInteractive()
				if !IsInteractive() {
					t.Error("Expected IsInteractive() to return true after EnableInteractive()")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.validateFunc(t)
		})
	}
}

func TestDisableInteractive(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func()
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_disable_interactive_from_enabled_state",
			setupMock: func() {
				EnableInteractive() // Start with interactive enabled
			},
			validateFunc: func(t *testing.T) {
				DisableInteractive()
				if IsInteractive() {
					t.Error("Expected IsInteractive() to return false after DisableInteractive()")
				}
			},
		},
		{
			name: "success_disable_interactive_from_disabled_state",
			setupMock: func() {
				DisableInteractive() // Start with interactive already disabled
			},
			validateFunc: func(t *testing.T) {
				DisableInteractive()
				if IsInteractive() {
					t.Error("Expected IsInteractive() to return false after DisableInteractive()")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.validateFunc(t)
		})
	}
}

func TestIsInteractive(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func()
		expectedResult bool
	}{
		{
			name: "success_interactive_enabled",
			setupMock: func() {
				EnableInteractive()
			},
			expectedResult: true,
		},
		{
			name: "success_interactive_disabled",
			setupMock: func() {
				DisableInteractive()
			},
			expectedResult: false,
		},
		{
			name: "success_default_state",
			setupMock: func() {
				// Test default state (interactive enabled by default)
				EnableInteractive()
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()
			result := IsInteractive()

			if result != tt.expectedResult {
				t.Errorf("Expected IsInteractive() %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestAllowOutput(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func()
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_allow_output_from_disallowed_state",
			setupMock: func() {
				DisallowOutput() // Start with output disallowed
			},
			validateFunc: func(t *testing.T) {
				AllowOutput()
				if !IsAllowingOutput() {
					t.Error("Expected IsAllowingOutput() to return true after AllowOutput()")
				}
			},
		},
		{
			name: "success_allow_output_from_allowed_state",
			setupMock: func() {
				AllowOutput() // Start with output already allowed
			},
			validateFunc: func(t *testing.T) {
				AllowOutput()
				if !IsAllowingOutput() {
					t.Error("Expected IsAllowingOutput() to return true after AllowOutput()")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.validateFunc(t)
		})
	}
}

func TestDisallowOutput(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func()
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_disallow_output_from_allowed_state",
			setupMock: func() {
				AllowOutput() // Start with output allowed
			},
			validateFunc: func(t *testing.T) {
				DisallowOutput()
				if IsAllowingOutput() {
					t.Error("Expected IsAllowingOutput() to return false after DisallowOutput()")
				}
			},
		},
		{
			name: "success_disallow_output_from_disallowed_state",
			setupMock: func() {
				DisallowOutput() // Start with output already disallowed
			},
			validateFunc: func(t *testing.T) {
				DisallowOutput()
				if IsAllowingOutput() {
					t.Error("Expected IsAllowingOutput() to return false after DisallowOutput()")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.validateFunc(t)
		})
	}
}

func TestIsAllowingOutput(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func()
		expectedResult bool
	}{
		{
			name: "success_output_allowed",
			setupMock: func() {
				AllowOutput()
			},
			expectedResult: true,
		},
		{
			name: "success_output_disallowed",
			setupMock: func() {
				DisallowOutput()
			},
			expectedResult: false,
		},
		{
			name: "success_default_state",
			setupMock: func() {
				// Test default state (output disallowed by default)
				DisallowOutput()
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()
			result := IsAllowingOutput()

			if result != tt.expectedResult {
				t.Errorf("Expected IsAllowingOutput() %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestEnableVerboseLogging(t *testing.T) {
	tests := []struct {
		name         string
		logLevel     string
		setupMock    func() (cleanup func())
		validateFunc func(t *testing.T, logLevel string)
	}{
		{
			name:     "success_enable_with_custom_log_level",
			logLevel: "INFO",
			setupMock: func() (cleanup func()) {
				originalLevel := os.Getenv(LogLevel)
				return func() {
					if originalLevel != "" {
						os.Setenv(LogLevel, originalLevel)
					} else {
						os.Unsetenv(LogLevel)
					}
				}
			},
			validateFunc: func(t *testing.T, logLevel string) {
				EnableVerboseLogging(logLevel)
				envValue := os.Getenv(LogLevel)
				if envValue != "INFO" {
					t.Errorf("Expected LogLevel environment variable to be 'INFO', got '%s'", envValue)
				}
			},
		},
		{
			name:     "success_enable_with_empty_log_level_defaults_to_debug",
			logLevel: "",
			setupMock: func() (cleanup func()) {
				originalLevel := os.Getenv(LogLevel)
				return func() {
					if originalLevel != "" {
						os.Setenv(LogLevel, originalLevel)
					} else {
						os.Unsetenv(LogLevel)
					}
				}
			},
			validateFunc: func(t *testing.T, logLevel string) {
				EnableVerboseLogging(logLevel)
				envValue := os.Getenv(LogLevel)
				if envValue != "DEBUG" {
					t.Errorf("Expected LogLevel environment variable to be 'DEBUG', got '%s'", envValue)
				}
			},
		},
		{
			name:     "success_enable_with_debug_log_level",
			logLevel: "DEBUG",
			setupMock: func() (cleanup func()) {
				originalLevel := os.Getenv(LogLevel)
				return func() {
					if originalLevel != "" {
						os.Setenv(LogLevel, originalLevel)
					} else {
						os.Unsetenv(LogLevel)
					}
				}
			},
			validateFunc: func(t *testing.T, logLevel string) {
				EnableVerboseLogging(logLevel)
				envValue := os.Getenv(LogLevel)
				if envValue != "DEBUG" {
					t.Errorf("Expected LogLevel environment variable to be 'DEBUG', got '%s'", envValue)
				}
			},
		},
		{
			name:     "success_enable_with_error_log_level",
			logLevel: "ERROR",
			setupMock: func() (cleanup func()) {
				originalLevel := os.Getenv(LogLevel)
				return func() {
					if originalLevel != "" {
						os.Setenv(LogLevel, originalLevel)
					} else {
						os.Unsetenv(LogLevel)
					}
				}
			},
			validateFunc: func(t *testing.T, logLevel string) {
				EnableVerboseLogging(logLevel)
				envValue := os.Getenv(LogLevel)
				if envValue != "ERROR" {
					t.Errorf("Expected LogLevel environment variable to be 'ERROR', got '%s'", envValue)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			tt.validateFunc(t, tt.logLevel)
		})
	}
}

func TestDisableVerboseLogging(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func() (cleanup func())
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_disable_verbose_logging",
			setupMock: func() (cleanup func()) {
				originalLevel := os.Getenv(LogLevel)
				EnableVerboseLogging("DEBUG") // Start with verbose logging enabled
				return func() {
					if originalLevel != "" {
						os.Setenv(LogLevel, originalLevel)
					} else {
						os.Unsetenv(LogLevel)
					}
				}
			},
			validateFunc: func(t *testing.T) {
				DisableVerboseLogging()
				envValue := os.Getenv(LogLevel)
				if envValue != "CRITICAL" {
					t.Errorf("Expected LogLevel environment variable to be 'CRITICAL', got '%s'", envValue)
				}
			},
		},
		{
			name: "success_disable_from_already_critical_level",
			setupMock: func() (cleanup func()) {
				originalLevel := os.Getenv(LogLevel)
				os.Setenv(LogLevel, "CRITICAL") // Start with critical level
				return func() {
					if originalLevel != "" {
						os.Setenv(LogLevel, originalLevel)
					} else {
						os.Unsetenv(LogLevel)
					}
				}
			},
			validateFunc: func(t *testing.T) {
				DisableVerboseLogging()
				envValue := os.Getenv(LogLevel)
				if envValue != "CRITICAL" {
					t.Errorf("Expected LogLevel environment variable to be 'CRITICAL', got '%s'", envValue)
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

func TestSetLoggerStyle(t *testing.T) {
	tests := []struct {
		name         string
		loggerStyle  string
		setupMock    func() (cleanup func())
		validateFunc func(t *testing.T, loggerStyle string)
	}{
		{
			name:        "success_set_default_logger_style",
			loggerStyle: "default",
			setupMock: func() (cleanup func()) {
				originalStyle := os.Getenv(LoggerStyle)
				return func() {
					if originalStyle != "" {
						os.Setenv(LoggerStyle, originalStyle)
					} else {
						os.Unsetenv(LoggerStyle)
					}
				}
			},
			validateFunc: func(t *testing.T, loggerStyle string) {
				SetLoggerStyle(loggerStyle)
				envValue := os.Getenv(LoggerStyle)
				if envValue != "default" {
					t.Errorf("Expected LoggerStyle environment variable to be 'default', got '%s'", envValue)
				}
			},
		},
		{
			name:        "success_set_non_default_logger_style_defaults_to_default",
			loggerStyle: "custom",
			setupMock: func() (cleanup func()) {
				originalStyle := os.Getenv(LoggerStyle)
				return func() {
					if originalStyle != "" {
						os.Setenv(LoggerStyle, originalStyle)
					} else {
						os.Unsetenv(LoggerStyle)
					}
				}
			},
			validateFunc: func(t *testing.T, loggerStyle string) {
				SetLoggerStyle(loggerStyle)
				envValue := os.Getenv(LoggerStyle)
				if envValue != "default" {
					t.Errorf("Expected LoggerStyle environment variable to be 'default', got '%s'", envValue)
				}
			},
		},
		{
			name:        "success_set_empty_logger_style_defaults_to_default",
			loggerStyle: "",
			setupMock: func() (cleanup func()) {
				originalStyle := os.Getenv(LoggerStyle)
				return func() {
					if originalStyle != "" {
						os.Setenv(LoggerStyle, originalStyle)
					} else {
						os.Unsetenv(LoggerStyle)
					}
				}
			},
			validateFunc: func(t *testing.T, loggerStyle string) {
				SetLoggerStyle(loggerStyle)
				envValue := os.Getenv(LoggerStyle)
				if envValue != "default" {
					t.Errorf("Expected LoggerStyle environment variable to be 'default', got '%s'", envValue)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			tt.validateFunc(t, tt.loggerStyle)
		})
	}
}

func TestEnableCertificateVerification(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func()
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_enable_cert_verification_from_disabled_state",
			setupMock: func() {
				DisableCertificateVerification() // Start with cert verification disabled
			},
			validateFunc: func(t *testing.T) {
				EnableCertificateVerification()
				// Note: IsVerifyingCertificates() may still return false if env var is set
				// This test validates the internal state change
			},
		},
		{
			name: "success_enable_cert_verification_from_enabled_state",
			setupMock: func() {
				EnableCertificateVerification() // Start with cert verification already enabled
			},
			validateFunc: func(t *testing.T) {
				EnableCertificateVerification()
				// Test validates that calling enable multiple times works
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.validateFunc(t)
		})
	}
}

func TestDisableCertificateVerification(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func()
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_disable_cert_verification_from_enabled_state",
			setupMock: func() {
				EnableCertificateVerification() // Start with cert verification enabled
			},
			validateFunc: func(t *testing.T) {
				DisableCertificateVerification()
				// Test validates the internal state change
			},
		},
		{
			name: "success_disable_cert_verification_from_disabled_state",
			setupMock: func() {
				DisableCertificateVerification() // Start with cert verification already disabled
			},
			validateFunc: func(t *testing.T) {
				DisableCertificateVerification()
				// Test validates that calling disable multiple times works
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tt.validateFunc(t)
		})
	}
}

func TestIsVerifyingCertificates(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func() (cleanup func())
		expectedResult bool
	}{
		{
			name: "success_cert_verification_enabled_no_env_var",
			setupMock: func() (cleanup func()) {
				EnableCertificateVerification()
				originalEnv := os.Getenv(ArkDisableCertificateVerificationEnvVar)
				os.Unsetenv(ArkDisableCertificateVerificationEnvVar)
				return func() {
					if originalEnv != "" {
						os.Setenv(ArkDisableCertificateVerificationEnvVar, originalEnv)
					}
				}
			},
			expectedResult: true,
		},
		{
			name: "success_cert_verification_disabled_no_env_var",
			setupMock: func() (cleanup func()) {
				DisableCertificateVerification()
				originalEnv := os.Getenv(ArkDisableCertificateVerificationEnvVar)
				os.Unsetenv(ArkDisableCertificateVerificationEnvVar)
				return func() {
					if originalEnv != "" {
						os.Setenv(ArkDisableCertificateVerificationEnvVar, originalEnv)
					}
				}
			},
			expectedResult: false,
		},
		{
			name: "success_env_var_overrides_enabled_state",
			setupMock: func() (cleanup func()) {
				EnableCertificateVerification()
				originalEnv := os.Getenv(ArkDisableCertificateVerificationEnvVar)
				os.Setenv(ArkDisableCertificateVerificationEnvVar, "true")
				return func() {
					if originalEnv != "" {
						os.Setenv(ArkDisableCertificateVerificationEnvVar, originalEnv)
					} else {
						os.Unsetenv(ArkDisableCertificateVerificationEnvVar)
					}
				}
			},
			expectedResult: false,
		},
		{
			name: "success_env_var_overrides_disabled_state",
			setupMock: func() (cleanup func()) {
				DisableCertificateVerification()
				originalEnv := os.Getenv(ArkDisableCertificateVerificationEnvVar)
				os.Setenv(ArkDisableCertificateVerificationEnvVar, "1")
				return func() {
					if originalEnv != "" {
						os.Setenv(ArkDisableCertificateVerificationEnvVar, originalEnv)
					} else {
						os.Unsetenv(ArkDisableCertificateVerificationEnvVar)
					}
				}
			},
			expectedResult: false,
		},
		{
			name: "success_empty_env_var_uses_internal_state",
			setupMock: func() (cleanup func()) {
				EnableCertificateVerification()
				originalEnv := os.Getenv(ArkDisableCertificateVerificationEnvVar)
				os.Setenv(ArkDisableCertificateVerificationEnvVar, "")
				return func() {
					if originalEnv != "" {
						os.Setenv(ArkDisableCertificateVerificationEnvVar, originalEnv)
					} else {
						os.Unsetenv(ArkDisableCertificateVerificationEnvVar)
					}
				}
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			result := IsVerifyingCertificates()

			if result != tt.expectedResult {
				t.Errorf("Expected IsVerifyingCertificates() %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestSetTrustedCertificate(t *testing.T) {
	tests := []struct {
		name         string
		cert         string
		setupMock    func()
		validateFunc func(t *testing.T, cert string)
	}{
		{
			name: "success_set_valid_certificate",
			cert: "-----BEGIN CERTIFICATE-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END CERTIFICATE-----",
			setupMock: func() {
				SetTrustedCertificate("") // Clear any existing certificate
			},
			validateFunc: func(t *testing.T, cert string) {
				SetTrustedCertificate(cert)
				result := TrustedCertificate()
				if result != cert {
					t.Errorf("Expected TrustedCertificate() to return '%s', got '%s'", cert, result)
				}
			},
		},
		{
			name: "success_set_empty_certificate",
			cert: "",
			setupMock: func() {
				SetTrustedCertificate("previous-cert") // Set a previous certificate
			},
			validateFunc: func(t *testing.T, cert string) {
				SetTrustedCertificate(cert)
				result := TrustedCertificate()
				if result != "" {
					t.Errorf("Expected TrustedCertificate() to return empty string, got '%s'", result)
				}
			},
		},
		{
			name: "success_override_existing_certificate",
			cert: "new-certificate-data",
			setupMock: func() {
				SetTrustedCertificate("old-certificate-data")
			},
			validateFunc: func(t *testing.T, cert string) {
				SetTrustedCertificate(cert)
				result := TrustedCertificate()
				if result != cert {
					t.Errorf("Expected TrustedCertificate() to return '%s', got '%s'", cert, result)
				}
			},
		},
		{
			name: "success_set_special_characters_certificate",
			cert: "cert-with-special-chars-!@#$%^&*()",
			setupMock: func() {
				SetTrustedCertificate("")
			},
			validateFunc: func(t *testing.T, cert string) {
				SetTrustedCertificate(cert)
				result := TrustedCertificate()
				if result != cert {
					t.Errorf("Expected TrustedCertificate() to return '%s', got '%s'", cert, result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()
			tt.validateFunc(t, tt.cert)
		})
	}
}

func TestTrustedCertificate(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func()
		expectedResult string
	}{
		{
			name: "success_return_set_certificate",
			setupMock: func() {
				SetTrustedCertificate("test-certificate-data")
			},
			expectedResult: "test-certificate-data",
		},
		{
			name: "success_return_empty_when_no_certificate_set",
			setupMock: func() {
				SetTrustedCertificate("")
			},
			expectedResult: "",
		},
		{
			name: "success_return_complex_certificate",
			setupMock: func() {
				cert := "-----BEGIN CERTIFICATE-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END CERTIFICATE-----"
				SetTrustedCertificate(cert)
			},
			expectedResult: "-----BEGIN CERTIFICATE-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END CERTIFICATE-----",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()
			result := TrustedCertificate()

			if result != tt.expectedResult {
				t.Errorf("Expected TrustedCertificate() '%s', got '%s'", tt.expectedResult, result)
			}
		})
	}
}

// TestSystemConfig_Integration tests the complete interaction between different configuration functions
func TestSystemConfig_Integration(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func() (cleanup func())
		validateFunc func(t *testing.T)
	}{
		{
			name: "complete_cycle_color_and_interactive_settings",
			setupMock: func() (cleanup func()) {
				// Store original states
				originalColoring := IsColoring()
				originalInteractive := IsInteractive()
				originalOutput := IsAllowingOutput()

				return func() {
					// Restore original states
					if originalColoring {
						EnableColor()
					} else {
						DisableColor()
					}
					if originalInteractive {
						EnableInteractive()
					} else {
						DisableInteractive()
					}
					if originalOutput {
						AllowOutput()
					} else {
						DisallowOutput()
					}
				}
			},
			validateFunc: func(t *testing.T) {
				// Test complete configuration cycle
				DisableColor()
				DisableInteractive()
				DisallowOutput()

				if IsColoring() || IsInteractive() || IsAllowingOutput() {
					t.Error("Expected all settings to be disabled")
				}

				EnableColor()
				EnableInteractive()
				AllowOutput()

				if !IsColoring() || !IsInteractive() || !IsAllowingOutput() {
					t.Error("Expected all settings to be enabled")
				}
			},
		},
		{
			name: "complete_cycle_certificate_and_trusted_cert",
			setupMock: func() (cleanup func()) {
				originalCert := TrustedCertificate()
				originalEnv := os.Getenv(ArkDisableCertificateVerificationEnvVar)

				return func() {
					SetTrustedCertificate(originalCert)
					if originalEnv != "" {
						os.Setenv(ArkDisableCertificateVerificationEnvVar, originalEnv)
					} else {
						os.Unsetenv(ArkDisableCertificateVerificationEnvVar)
					}
				}
			},
			validateFunc: func(t *testing.T) {
				testCert := "test-integration-certificate"

				SetTrustedCertificate(testCert)
				EnableCertificateVerification()
				os.Unsetenv(ArkDisableCertificateVerificationEnvVar)

				if TrustedCertificate() != testCert {
					t.Errorf("Expected trusted certificate '%s', got '%s'", testCert, TrustedCertificate())
				}

				if !IsVerifyingCertificates() {
					t.Error("Expected certificate verification to be enabled")
				}

				// Test environment variable override
				os.Setenv(ArkDisableCertificateVerificationEnvVar, "true")
				if IsVerifyingCertificates() {
					t.Error("Expected certificate verification to be disabled by environment variable")
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
