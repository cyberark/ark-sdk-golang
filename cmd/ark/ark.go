// Package main provides the entry point for the Ark CLI application.
//
// The Ark CLI is a command-line interface that provides access to various
// Ark services and functionality including profile management, authentication,
// configuration, caching, and service execution.
//
// The application uses the Cobra library for command-line interface management
// and supports multiple subcommands for different operations. Build information
// including version, build number, build date, and git commit are embedded
// at compile time through build variables.
//
// Example usage:
//
//	ark --version
//	ark profiles list
//	ark login
//	ark configure
package main

import (
	"fmt"
	"os"

	"github.com/cyberark/ark-sdk-golang/pkg/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/spf13/cobra"
)

var (
	// GitCommit is the commit hash of the Ark CLI application.
	GitCommit = "N/A"

	// BuildDate is the build date of the Ark CLI application.
	BuildDate = "N/A"
	// Version is the version of the Ark CLI application.
	Version = "N/A"

	// BuildNumber is the build number of the Ark CLI application.
	BuildNumber = "N/A"
)

// main is the entry point for the Ark CLI application.
//
// This function initializes the Cobra root command with version information,
// sets up the application version in the common package, creates a profiles
// loader, and registers all available actions (profiles, cache, configure,
// login, and service execution) with the root command.
//
// The function handles command execution and exits with code 1 if an error
// occurs during command execution. The version template is customized to
// display build information in a specific format.
//
// Build variables (GitCommit, BuildDate, Version, BuildNumber) are expected
// to be set at compile time using ldflags but will default to "N/A" if not
// provided.
//
// Available commands after initialization:
//   - profiles: Manage user profiles
//   - cache: Manage application cache
//   - configure: Configure the CLI
//   - login: Authenticate with services
//   - exec: Execute service actions
//
// The function will call os.Exit(1) if command execution fails.
func main() {
	var rootCmd = &cobra.Command{
		Use:     "ark",
		Version: fmt.Sprintf("Version: %s\nBuild Number: %s\nBuild Date: %s\nGit Commit: %s", Version, BuildNumber, BuildDate, GitCommit),
		Short:   "Ark CLI",
	}
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	common.SetArkVersion(Version)
	profilesLoader := profiles.DefaultProfilesLoader()
	arkActions := []actions.ArkAction{
		actions.NewArkProfilesAction(profilesLoader),
		actions.NewArkCacheAction(),
		actions.NewArkConfigureAction(profilesLoader),
		actions.NewArkLoginAction(profilesLoader),
		actions.NewArkServiceExecAction(profilesLoader),
	}

	for _, action := range arkActions {
		action.DefineAction(rootCmd)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
