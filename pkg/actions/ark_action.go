package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ArkAction is an interface that defines the structure for actions in the Ark SDK.
type ArkAction interface {
	DefineAction(cmd *cobra.Command)
}

// ArkBaseAction is a struct that implements the ArkAction interface as a base action.
type ArkBaseAction struct {
	logger *common.ArkLogger
}

// NewArkBaseAction creates a new instance of ArkBaseAction.
func NewArkBaseAction() *ArkBaseAction {
	return &ArkBaseAction{
		logger: common.GetLogger("ArkBaseAction", common.Unknown),
	}
}

// CommonActionsConfiguration sets up common flags for the command line interface.
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
