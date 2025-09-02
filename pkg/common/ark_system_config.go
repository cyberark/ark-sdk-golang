// Package common provides shared utilities and types for the ARK SDK.
//
// This package handles configuration for colored output, interactive mode, certificate
// verification, output control, logging levels, and trusted certificates. It provides
// a centralized way to control various system behaviors through global state and
// environment variables.
package common

import (
	"os"
)

var (
	noColor                   = false
	isInteractive             = true
	isCertificateVerification = true
	isAllowingOutput          = false
	trustedCert               = ""
)

// ArkDisableCertificateVerificationEnvVar is the environment variable name for disabling certificate validation.
//
// When this environment variable is set to any non-empty value, certificate verification
// will be disabled regardless of the internal isCertificateVerification setting.
const (
	ArkDisableCertificateVerificationEnvVar = "ARK_DISABLE_CERTIFICATE_VERIFICATION"
)

// DisableColor disables colored output in the console.
//
// DisableColor sets the global noColor flag to true, which will cause IsColoring()
// to return false. This affects all subsequent console output that checks for
// color support throughout the application.
//
// Example:
//
//	DisableColor()
//	if IsColoring() {
//	    // This block will not execute
//	}
func DisableColor() {
	noColor = true
}

// EnableColor enables colored output in the console.
//
// EnableColor sets the global noColor flag to false, which will cause IsColoring()
// to return true. This enables colored console output throughout the application.
//
// Example:
//
//	EnableColor()
//	if IsColoring() {
//	    // This block will execute, allowing colored output
//	}
func EnableColor() {
	noColor = false
}

// IsColoring checks if colored output is enabled.
//
// IsColoring returns true when colored output is enabled (noColor is false) and
// false when colored output is disabled. This function is used throughout the
// application to determine whether to apply color formatting to console output.
//
// Returns true if colored output is enabled, false otherwise.
//
// Example:
//
//	if IsColoring() {
//	    fmt.Print("\033[31mRed text\033[0m")
//	} else {
//	    fmt.Print("Plain text")
//	}
func IsColoring() bool {
	return !noColor
}

// EnableInteractive enables interactive mode.
//
// EnableInteractive sets the global isInteractive flag to true, allowing the
// application to prompt for user input and display interactive elements.
//
// Example:
//
//	EnableInteractive()
//	if IsInteractive() {
//	    // Show interactive prompts
//	}
func EnableInteractive() {
	isInteractive = true
}

// DisableInteractive disables interactive mode.
//
// DisableInteractive sets the global isInteractive flag to false, preventing
// the application from displaying interactive prompts or requiring user input.
// This is useful for automated scripts or CI/CD environments.
//
// Example:
//
//	DisableInteractive()
//	if IsInteractive() {
//	    // This block will not execute
//	}
func DisableInteractive() {
	isInteractive = false
}

// IsInteractive checks if interactive mode is enabled.
//
// IsInteractive returns true when the application is allowed to display
// interactive prompts and request user input, and false when running in
// non-interactive mode (suitable for automation).
//
// Returns true if interactive mode is enabled, false otherwise.
//
// Example:
//
//	if IsInteractive() {
//	    response := promptUser("Continue? (y/n): ")
//	}
func IsInteractive() bool {
	return isInteractive
}

// AllowOutput allows output to be displayed.
//
// AllowOutput sets the global isAllowingOutput flag to true, enabling the
// application to display output messages, logs, and other information to
// the console or other output destinations.
//
// Example:
//
//	AllowOutput()
//	if IsAllowingOutput() {
//	    fmt.Println("This message will be displayed")
//	}
func AllowOutput() {
	isAllowingOutput = true
}

// DisallowOutput disallows output to be displayed.
//
// DisallowOutput sets the global isAllowingOutput flag to false, preventing
// the application from displaying output. This is useful for silent operation
// modes or when output needs to be suppressed.
//
// Example:
//
//	DisallowOutput()
//	if IsAllowingOutput() {
//	    // This block will not execute
//	}
func DisallowOutput() {
	isAllowingOutput = false
}

// IsAllowingOutput checks if output is allowed to be displayed.
//
// IsAllowingOutput returns true when the application is permitted to display
// output messages and false when output should be suppressed.
//
// Returns true if output is allowed, false otherwise.
//
// Example:
//
//	if IsAllowingOutput() {
//	    logger.Info("Operation completed successfully")
//	}
func IsAllowingOutput() bool {
	return isAllowingOutput
}

// EnableVerboseLogging enables verbose logging with the specified log level.
//
// EnableVerboseLogging sets the LogLevel environment variable to the provided
// log level string. If an empty string is provided, it defaults to "DEBUG".
// This affects the logging verbosity throughout the application.
//
// Parameters:
//   - logLevel: The desired log level string (defaults to "DEBUG" if empty)
//
// Example:
//
//	EnableVerboseLogging("INFO")
//	EnableVerboseLogging("") // Uses "DEBUG" as default
func EnableVerboseLogging(logLevel string) {
	if logLevel == "" {
		logLevel = "DEBUG"
	}
	_ = os.Setenv(LogLevel, logLevel)
}

// DisableVerboseLogging disables verbose logging.
//
// DisableVerboseLogging sets the LogLevel environment variable to "CRITICAL",
// effectively reducing the logging output to only critical messages.
//
// Example:
//
//	DisableVerboseLogging()
//	// Only critical log messages will be displayed
func DisableVerboseLogging() {
	_ = os.Setenv(LogLevel, "CRITICAL")
}

// SetLoggerStyle sets the logger style based on the provided string.
//
// SetLoggerStyle configures the LoggerStyle environment variable. If the
// provided style is "default", it sets the style to "default"; otherwise,
// it defaults to "default" regardless of the input value.
//
// Parameters:
//   - loggerStyle: The desired logger style ("default" or any other value defaults to "default")
//
// Example:
//
//	SetLoggerStyle("default")
//	SetLoggerStyle("custom") // Also sets to "default"
func SetLoggerStyle(loggerStyle string) {
	if loggerStyle == "default" {
		_ = os.Setenv(LoggerStyle, loggerStyle)
	} else {
		_ = os.Setenv(LoggerStyle, "default")
	}
}

// EnableCertificateVerification enables certificate verification.
//
// EnableCertificateVerification sets the global isCertificateVerification flag
// to true, enabling SSL/TLS certificate validation for network connections.
// Note that if the ArkDisableCertificateVerificationEnvVar environment variable
// is set, certificate verification will still be disabled.
//
// Example:
//
//	EnableCertificateVerification()
//	if IsVerifyingCertificates() {
//	    // Certificates will be verified
//	}
func EnableCertificateVerification() {
	isCertificateVerification = true
}

// DisableCertificateVerification disables certificate verification.
//
// DisableCertificateVerification sets the global isCertificateVerification flag
// to false, disabling SSL/TLS certificate validation for network connections.
// This should be used with caution as it reduces security.
//
// Example:
//
//	DisableCertificateVerification()
//	if IsVerifyingCertificates() {
//	    // This block will not execute
//	}
func DisableCertificateVerification() {
	isCertificateVerification = false
}

// IsVerifyingCertificates checks if certificate verification is enabled.
//
// IsVerifyingCertificates returns false if the ArkDisableCertificateVerificationEnvVar
// environment variable is set to any non-empty value, regardless of the internal
// isCertificateVerification setting. Otherwise, it returns the value of the
// isCertificateVerification flag.
//
// Returns true if certificate verification is enabled, false otherwise.
//
// Example:
//
//	if IsVerifyingCertificates() {
//	    // Use secure connection with certificate validation
//	} else {
//	    // Use connection without certificate validation
//	}
func IsVerifyingCertificates() bool {
	if os.Getenv(ArkDisableCertificateVerificationEnvVar) != "" {
		return false
	}
	return isCertificateVerification
}

// SetTrustedCertificate sets the trusted certificate for verification.
//
// SetTrustedCertificate stores the provided certificate string in the global
// trustedCert variable. This certificate can be used for custom certificate
// validation scenarios.
//
// Parameters:
//   - cert: The certificate string to be stored as trusted
//
// Example:
//
//	cert := "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
//	SetTrustedCertificate(cert)
func SetTrustedCertificate(cert string) {
	trustedCert = cert
}

// TrustedCertificate returns the trusted certificate for verification.
//
// TrustedCertificate retrieves the currently stored trusted certificate string
// that was previously set using SetTrustedCertificate. Returns an empty string
// if no certificate has been set.
//
// Returns the trusted certificate string, or empty string if none is set.
//
// Example:
//
//	cert := TrustedCertificate()
//	if cert != "" {
//	    // Use the trusted certificate for validation
//	}
func TrustedCertificate() string {
	return trustedCert
}
