// Package actions provides base implementation for CLI exec actions in the ARK SDK.
//
// This package contains the core functionality for defining and executing CLI commands
// that integrate with various authentication providers and profiles.
package actions

import (
	"fmt"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/cli"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/args"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ArkExecAction defines the interface for executing CLI actions in the ARK SDK.
//
// Implementations of this interface provide the actual command definition and execution logic
// for specific ARK SDK operations. The interface enables pluggable action implementations
// while maintaining consistent command structure and behavior.
type ArkExecAction interface {
	// DefineExecAction configures the specific exec command structure and flags.
	//
	// This method is called during command initialization to set up the CLI command
	// structure, including any subcommands, flags, and help text specific to the action.
	//
	// Parameters:
	//   - cmd: The parent cobra command to which the exec action will be added
	//
	// Returns an error if the command definition fails.
	DefineExecAction(cmd *cobra.Command) error

	// RunExecAction executes the specific action logic with an authenticated API client.
	//
	// This method is called after authentication and profile loading are complete.
	// It receives a fully configured CLI API client and should implement the actual
	// business logic for the specific action.
	//
	// Parameters:
	//   - api: Authenticated ARK CLI API client
	//   - cmd: The root command being executed
	//   - execCmd: The specific exec command being run
	//   - args: Command line arguments passed to the action
	//
	// Returns an error if the action execution fails.
	RunExecAction(api *cli.ArkCLIAPI, cmd *cobra.Command, execCmd *cobra.Command, args []string) error
}

// ArkBaseExecAction provides a base implementation for ARK CLI exec actions.
//
// This struct serves as a foundation for implementing specific exec actions by combining
// the common ARK action functionality with exec-specific behavior. It handles profile
// loading, authentication, and action execution orchestration.
//
// The struct embeds ArkBaseAction to inherit common CLI functionality and adds
// exec-specific capabilities like profile management and action delegation.
type ArkBaseExecAction struct {
	*ArkBaseAction
	profilesLoader *profiles.ProfileLoader
	execAction     *ArkExecAction
	logger         *common.ArkLogger
}

// NewArkBaseExecAction creates a new instance of ArkBaseExecAction with the specified configuration.
//
// NewArkBaseExecAction initializes a base exec action with the provided action implementation,
// name for logging, and profile loader. The returned instance can be used to define
// and execute CLI commands that require profile-based authentication. The constructor
// sets up logging with the specified name and configures all necessary dependencies
// for command execution.
//
// Parameters:
//   - execAction: Implementation of the specific exec action behavior
//   - name: Identifier used for logging and error reporting
//   - profilesLoader: Service for loading and managing authentication profiles
//
// Returns a configured ArkBaseExecAction ready for command definition and execution.
//
// Example:
//
//	action := NewArkBaseExecAction(
//	    &myExecAction,
//	    "my-action",
//	    profileLoader,
//	)
//	action.DefineAction(rootCmd)
func NewArkBaseExecAction(execAction *ArkExecAction, name string, profilesLoader *profiles.ProfileLoader) *ArkBaseExecAction {
	return &ArkBaseExecAction{
		ArkBaseAction:  NewArkBaseAction(),
		profilesLoader: profilesLoader,
		execAction:     execAction,
		logger:         common.GetLogger(name, common.Unknown),
	}
}

// DefineAction defines the CLI `exec` command structure and configuration.
//
// DefineAction creates and configures the `exec` subcommand with all necessary flags
// and execution behavior. It sets up the command hierarchy, adds persistent flags
// for profile management and execution control, and delegates specific action
// definition to the embedded exec action implementation.
//
// The method configures the following persistent flags:
//   - profile-name: Specifies which authentication profile to use
//   - output-path: Optional file path for writing command output
//   - request-file: Optional file containing action parameters
//   - retry-count: Number of retry attempts for failed executions
//   - refresh-auth: Forces authentication token refresh
//
// Parameters:
//   - cmd: The parent cobra command to which the exec command will be added
//
// The method panics if the embedded exec action fails to define its specific behavior.
//
// Example:
//
//	baseAction := NewArkBaseExecAction(&myAction, "test", loader)
//	baseAction.DefineAction(rootCmd)
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

// runExecAction executes the configured action with profile-based authentication.
//
// runExecAction orchestrates the complete execution flow including profile loading,
// authentication validation, API client creation, and action execution with retry logic.
// The method performs the following steps:
//  1. Locates the exec command in the command hierarchy
//  2. Loads the specified authentication profile
//  3. Validates and loads available authenticators
//  4. Creates an authenticated CLI API client
//  5. Executes the action with configured retry logic
//
// The method handles authentication token expiration and provides user feedback
// for authentication status. If not all authenticators are available, it warns
// the user that some functionality may be disabled.
//
// Parameters:
//   - cmd: The root command being executed
//   - execArgs: Command line arguments passed to the exec action
//
// The method prints failure messages and returns early if any critical step fails,
// including profile loading, authentication, or API client creation.
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
