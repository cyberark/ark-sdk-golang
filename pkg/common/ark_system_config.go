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

// ArkDisableCertificateVerificationEnvVar Environment variable for disabling certificate validation
const (
	ArkDisableCertificateVerificationEnvVar = "ARK_DISABLE_CERTIFICATE_VERIFICATION"
)

// DisableColor disables colored output in the console.
func DisableColor() {
	noColor = true
}

// EnableColor enables colored output in the console.
func EnableColor() {
	noColor = false
}

// IsColoring checks if colored output is enabled.
func IsColoring() bool {
	return !noColor
}

// EnableInteractive enables interactive mode.
func EnableInteractive() {
	isInteractive = true
}

// DisableInteractive disables interactive mode.
func DisableInteractive() {
	isInteractive = false
}

// IsInteractive checks if interactive mode is enabled.
func IsInteractive() bool {
	return isInteractive
}

// AllowOutput allows output to be displayed.
func AllowOutput() {
	isAllowingOutput = true
}

// DisallowOutput disallows output to be displayed.
func DisallowOutput() {
	isAllowingOutput = false
}

// IsAllowingOutput checks if output is allowed to be displayed.
func IsAllowingOutput() bool {
	return isAllowingOutput
}

// EnableVerboseLogging enables verbose logging with the specified log level.
func EnableVerboseLogging(logLevel string) {
	if logLevel == "" {
		logLevel = "DEBUG"
	}
	_ = os.Setenv(LogLevel, logLevel)
}

// DisableVerboseLogging disables verbose logging.
func DisableVerboseLogging() {
	_ = os.Setenv(LogLevel, "CRITICAL")
}

// SetLoggerStyle sets the logger style based on the provided string.
func SetLoggerStyle(loggerStyle string) {
	if loggerStyle == "default" {
		_ = os.Setenv(LoggerStyle, loggerStyle)
	} else {
		_ = os.Setenv(LoggerStyle, "default")
	}
}

// EnableCertificateVerification enables certificate verification.
func EnableCertificateVerification() {
	isCertificateVerification = true
}

// DisableCertificateVerification disables certificate verification.
func DisableCertificateVerification() {
	isCertificateVerification = false
}

// IsVerifyingCertificates checks if certificate verification is enabled.
func IsVerifyingCertificates() bool {
	if os.Getenv(ArkDisableCertificateVerificationEnvVar) != "" {
		return false
	}
	return isCertificateVerification
}

// SetTrustedCertificate sets the trusted certificate for verification.
func SetTrustedCertificate(cert string) {
	trustedCert = cert
}

// TrustedCertificate returns the trusted certificate for verification.
func TrustedCertificate() string {
	return trustedCert
}
