package actions

import (
	"github.com/spf13/cobra"
	commonInternal "github.com/cyberark/ark-sdk-golang/internal/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	commonArgs "github.com/cyberark/ark-sdk-golang/pkg/common/args"
	"os"
	"reflect"
)

// ArkCacheAction is a struct that implements the ArkAction interface for cache management.
type ArkCacheAction struct {
	*ArkBaseAction
}

// NewArkCacheAction Creates a new instance of ArkCacheAction
func NewArkCacheAction() *ArkCacheAction {
	return &ArkCacheAction{
		ArkBaseAction: NewArkBaseAction(),
	}
}

// DefineAction Defines the CLI `cache` action, and adds the clear cache function
func (a *ArkCacheAction) DefineAction(cmd *cobra.Command) {
	cacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "Manage cache",
	}
	cacheCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		a.CommonActionsExecution(cmd, args)
	}
	a.CommonActionsConfiguration(cacheCmd)

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clears all profiles cache",
		Run:   a.runClearCacheAction,
	}

	cacheCmd.AddCommand(clearCmd)
	cmd.AddCommand(cacheCmd)
}

func (a *ArkCacheAction) runClearCacheAction(cmd *cobra.Command, args []string) {
	if keyring, err := common.NewArkKeyring("").GetKeyring(false); err == nil {
		if reflect.TypeOf(keyring) != reflect.TypeOf(&commonInternal.BasicKeyring{}) {
			commonArgs.PrintNormal("Cache clear is only valid for basic keyring implementation at the moment")
			return
		}
		cacheFolderPath := os.ExpandEnv("$HOME") + string(os.PathSeparator) + commonInternal.DefaultBasicKeyringFolder
		if envPath := os.Getenv(commonInternal.ArkBasicKeyringFolderEnvVar); envPath != "" {
			cacheFolderPath = envPath
		}
		_ = os.Remove(cacheFolderPath + string(os.PathSeparator) + "keyring")
		_ = os.Remove(cacheFolderPath + string(os.PathSeparator) + "mac")
	}
}
