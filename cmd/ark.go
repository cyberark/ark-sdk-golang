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
