package actions

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/cli"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/args"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"time"
)

// ArkExecAction is an interface that defines the structure for executing actions in the Ark SDK.
type ArkExecAction interface {
	DefineExecAction(cmd *cobra.Command) error
	RunExecAction(api *cli.ArkCLIAPI, cmd *cobra.Command, execCmd *cobra.Command, args []string) error
}

// ArkBaseExecAction is a struct that implements the ArkExecAction interface as a base action.
type ArkBaseExecAction struct {
	*ArkBaseAction
	profilesLoader *profiles.ProfileLoader
	execAction     *ArkExecAction
	logger         *common.ArkLogger
}

// NewArkBaseExecAction creates a new instance of ArkBaseExecAction.
func NewArkBaseExecAction(execAction *ArkExecAction, name string, profilesLoader *profiles.ProfileLoader) *ArkBaseExecAction {
	return &ArkBaseExecAction{
		ArkBaseAction:  NewArkBaseAction(),
		profilesLoader: profilesLoader,
		execAction:     execAction,
		logger:         common.GetLogger(name, common.Unknown),
	}
}

// DefineAction defines the CLI `exec` action, and adds the exec action function.
func (a *ArkBaseExecAction) DefineAction(cmd *cobra.Command) {
	execCmd := &cobra.Command{
		Use:   "exec",
		Short: "Exec an action",
	}
	execCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		a.CommonActionsExecution(cmd, args)
	}
	a.CommonActionsConfiguration(execCmd)

	execCmd.PersistentFlags().String("profile-name", profiles.DefaultProfileName(), "Profile name to load")
	execCmd.PersistentFlags().String("output-path", "", "Output file to write data to")
	execCmd.PersistentFlags().String("request-file", "", "Request file containing the parameters for the exec action")
	execCmd.PersistentFlags().Int("retry-count", 1, "Retry count for execution")
	execCmd.PersistentFlags().Bool("refresh-auth", false, "If a cache exists, will also try to refresh it")
	err := (*a.execAction).DefineExecAction(execCmd)
	if err != nil {
		args.PrintFailure(fmt.Sprintf("Error defining exec action %v", err))
		panic(err)
	}
	cmd.AddCommand(execCmd)
}

func (a *ArkBaseExecAction) runExecAction(cmd *cobra.Command, execArgs []string) {
	a.CommonActionsExecution(cmd, execArgs)
	var execCmd *cobra.Command
	currentCmd := cmd
	for currentCmd != nil {
		if currentCmd.Use == "exec" {
			execCmd = currentCmd
			break
		}
		currentCmd = currentCmd.Parent()
	}
	if execCmd == nil {
		args.PrintFailure("Failed to find exec command")
		return
	}
	profileName, _ := execCmd.Flags().GetString("profile-name")
	profile, err := (*a.profilesLoader).LoadProfile(profiles.DeduceProfileName(profileName))
	if err != nil || profile == nil {
		args.PrintFailure("Please configure a profile before trying to login")
		return
	}

	var authenticators []auth.ArkAuth
	for authenticatorName := range profile.AuthProfiles {
		authenticator := auth.SupportedAuthenticators[authenticatorName]
		refreshAuth, _ := cmd.Flags().GetBool("refresh-auth")
		token, err := authenticator.LoadAuthentication(profile, refreshAuth)
		if err != nil || token == nil {
			continue
		}
		if time.Now().After(time.Time(token.ExpiresIn)) {
			continue
		}
		authenticators = append(authenticators, authenticator)
	}

	if len(authenticators) == 0 {
		args.PrintFailure("Failed to load authenticators, tokens are either expired or authenticators are not logged in, please login first")
		return
	}
	if len(authenticators) != len(profile.AuthProfiles) && common.IsInteractive() {
		args.PrintColored("Not all authenticators are logged in, some of the functionality will be disabled", color.New())
	}

	// Create the CLI API with the authenticators
	api, err := cli.NewArkCLIAPI(authenticators, profile)
	if err != nil {
		args.PrintFailure(fmt.Sprintf("Failed to create CLI API: %s", err))
		return
	}

	// Run the actual exec fitting action with the api
	// Run it with retries as per defined by user
	retryCount, _ := execCmd.Flags().GetInt("retry-count")
	err = common.RetryCall(func() error {
		return (*a.execAction).RunExecAction(api, cmd, execCmd, execArgs)
	}, retryCount, 1, nil, 1, 0, func(err error, delay int) {
		args.PrintFailure(fmt.Sprintf("Retrying in %d seconds", delay))
	})

	if err != nil {
		args.PrintFailure(fmt.Sprintf("Failed to execute action: %s", err))
	}
}
