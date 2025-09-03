// Package actions provides base functionality for Ark SDK command line actions.
//
// This package defines the core interfaces and base implementations for creating
// command line actions in the Ark SDK. It includes configuration management,
// flag handling, and common execution patterns that can be shared across
// different CLI commands.
package actions

import (
	"os"

	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ArkAction is an interface that defines the structure for actions in the Ark SDK.
//
// ArkAction provides a contract for implementing command line actions that can
// be integrated with the Ark SDK CLI framework. Implementations should define
// their specific command behavior through the DefineAction method.
type ArkAction interface {
	// DefineAction configures the provided cobra command with action-specific behavior
	DefineAction(cmd *cobra.Command)
}

// ArkBaseAction is a struct that implements the ArkAction interface as a base action.
//
// ArkBaseAction provides common functionality that can be shared across different
// action implementations. It includes logger management and common flag handling
// patterns. This struct can be embedded in more specific action implementations
// to provide consistent behavior across the CLI.
type ArkBaseAction struct {
	// logger is the internal logger instance for the action
	logger *common.ArkLogger
}

// NewArkBaseAction creates a new instance of ArkBaseAction.
//
// NewArkBaseAction initializes a new ArkBaseAction with a configured logger.
// The logger is set up with a default configuration using the "ArkBaseAction"
// name and Unknown log level.
//
// Returns a new ArkBaseAction instance ready for use.
//
// Example:
//
//	action := NewArkBaseAction()
//	action.CommonActionsConfiguration(cmd)
func NewArkBaseAction() *ArkBaseAction {
	return &ArkBaseAction{
		logger: common.GetLogger("ArkBaseAction", common.Unknown),
	}
}

// CommonActionsConfiguration sets up common flags for the command line interface.
//
// CommonActionsConfiguration adds standard persistent flags to the provided cobra
// command that are commonly used across different Ark SDK actions. These flags
// control logging behavior, output formatting, certificate handling, and other
// common CLI options.
//
// The following flags are added:
//   - raw: Controls whether output should be in raw format
//   - silent: Enables silent execution without interactive prompts
//   - allow-output: Allows stdout/stderr output even in silent mode
//   - verbose: Enables verbose logging
//   - logger-style: Specifies the style for verbose logging
//   - log-level: Sets the log level for verbose mode
//   - disable-cert-verification: Disables HTTPS certificate verification (unsafe)
//   - trusted-cert: Specifies a trusted certificate for HTTPS calls
//
// Parameters:
//   - cmd: The cobra command to configure with persistent flags
//
// Example:
//
//	action := NewArkBaseAction()
//	action.CommonActionsConfiguration(rootCmd)
func (a *ArkBaseAction) CommonActionsConfiguration(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("raw", false, "Whether to raw output")
	cmd.PersistentFlags().Bool("silent", false, "Silent execution, no interactiveness")
	cmd.PersistentFlags().Bool("allow-output", false, "Allow stdout / stderr even when silent and not interactive")
	cmd.PersistentFlags().Bool("verbose", false, "Whether to verbose log")
	cmd.PersistentFlags().String("logger-style", "default", "Which verbose logger style to use")
	cmd.PersistentFlags().String("log-level", "INFO", "Log level to use while verbose")
	cmd.PersistentFlags().Bool("disable-cert-verification", false, "Disables certificate verification on HTTPS calls, unsafe! Avoid using in production environments!")
	cmd.PersistentFlags().String("trusted-cert", "", "Certificate to use for HTTPS calls")
}

// CommonActionsExecution executes common actions based on the command line flags.
//
// CommonActionsExecution processes the standard flags set up by CommonActionsConfiguration
// and applies the corresponding configuration changes to the Ark SDK runtime. This
// function should be called early in command execution to ensure proper setup.
//
// The function performs the following operations:
//  1. Sets default states for color, interactivity, logging, and certificates
//  2. Processes each flag and applies the corresponding configuration
//  3. Handles certificate verification settings (disable or trusted cert)
//  4. Configures profile name if provided
//  5. Sets default DEPLOY_ENV if not already set
//
// Parameters:
//   - cmd: The cobra command containing the parsed flags
//   - args: Command line arguments (not currently used but part of cobra pattern)
//
// The function ignores flag parsing errors and uses default values in such cases,
// following the principle of graceful degradation for CLI flag handling.
//
// Example:
//
//	action := NewArkBaseAction()
//	action.CommonActionsExecution(cmd, args)
func (a *ArkBaseAction) CommonActionsExecution(cmd *cobra.Command, args []string) {
	common.EnableColor()
	common.EnableInteractive()
	common.DisableVerboseLogging()
	common.DisallowOutput()
	common.SetLoggerStyle(viper.GetString("logger-style"))
	common.EnableCertificateVerification()

	if raw, err := cmd.Flags().GetBool("raw"); err == nil && raw {
		common.DisableColor()
	}
	if silent, err := cmd.Flags().GetBool("silent"); err == nil && silent {
		common.DisableInteractive()
	}
	if verbose, err := cmd.Flags().GetBool("verbose"); err == nil && verbose {
		common.EnableVerboseLogging(viper.GetString("log-level"))
	}
	if allowOutput, err := cmd.Flags().GetBool("allow-output"); err == nil && allowOutput {
		common.AllowOutput()
	}
	if disableCertValidation, err := cmd.Flags().GetBool("disable-cert-verification"); err == nil && disableCertValidation {
		common.DisableCertificateVerification()
	} else if trustedCert, err := cmd.Flags().GetString("trusted-cert"); err == nil && trustedCert != "" {
		common.SetTrustedCertificate(viper.GetString("trusted-cert"))
	}
	a.logger = common.GetLogger("ArkBaseAction", common.Unknown)

	if profileName, err := cmd.Flags().GetString("profile-name"); err == nil && profileName != "" {
		viper.Set("profile-name", profiles.DeduceProfileName(profileName))
	}
	if os.Getenv("DEPLOY_ENV") == "" {
		_ = os.Setenv("DEPLOY_ENV", "prod")
	}
}
