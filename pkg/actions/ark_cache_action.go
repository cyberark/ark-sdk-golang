package actions

import (
	"os"
	"reflect"

	commoninternal "github.com/cyberark/ark-sdk-golang/internal/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	commonargs "github.com/cyberark/ark-sdk-golang/pkg/common/args"
	"github.com/spf13/cobra"
)

// ArkCacheAction is a struct that implements the ArkAction interface for cache management.
//
// ArkCacheAction provides functionality for managing cache operations in the Ark SDK CLI.
// It embeds ArkBaseAction to inherit common CLI functionality and adds specific cache
// management commands such as clearing cached profiles and credentials.
//
// The action specifically supports clearing cache for basic keyring implementations,
// which store credentials and profile information in local files.
type ArkCacheAction struct {
	// ArkBaseAction provides common action functionality
	*ArkBaseAction
}

// NewArkCacheAction creates a new instance of ArkCacheAction.
//
// NewArkCacheAction initializes a new ArkCacheAction with an embedded ArkBaseAction,
// providing all the common CLI functionality along with cache-specific operations.
// The returned instance is ready to be used for defining cache management commands.
//
// Returns a new ArkCacheAction instance with initialized base action functionality.
//
// Example:
//
//	cacheAction := NewArkCacheAction()
//	cacheAction.DefineAction(rootCmd)
func NewArkCacheAction() *ArkCacheAction {
	return &ArkCacheAction{
		ArkBaseAction: NewArkBaseAction(),
	}
}

// DefineAction defines the CLI cache action and adds cache management subcommands.
//
// DefineAction creates a "cache" command with subcommands for cache management operations.
// Currently, it provides a "clear" subcommand that removes cached credentials and profile
// information from the basic keyring implementation.
//
// The function sets up:
//  1. A parent "cache" command for cache management
//  2. Common action configuration through embedded ArkBaseAction
//  3. A "clear" subcommand for clearing cached data
//  4. Persistent pre-run hook for common action execution
//
// Parameters:
//   - cmd: The parent cobra command to which the cache command will be added
//
// The cache command supports the standard CLI flags inherited from ArkBaseAction
// and executes common action setup before running cache-specific operations.
//
// Example:
//
//	cacheAction := NewArkCacheAction()
//	cacheAction.DefineAction(rootCmd)
//	// This adds: myapp cache clear
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

// runClearCacheAction clears cached credentials and profile data for basic keyring implementations.
//
// runClearCacheAction attempts to clear cached data by removing the keyring and MAC files
// from the basic keyring storage location. The function only operates on basic keyring
// implementations and will print a message for other keyring types.
//
// The function performs the following operations:
//  1. Creates a new ArkKeyring instance and retrieves the keyring
//  2. Checks if the keyring is a BasicKeyring implementation
//  3. If not BasicKeyring, prints an informational message and returns
//  4. Determines the cache folder path from HOME or environment variable
//  5. Removes the "keyring" and "mac" files from the cache folder
//
// Parameters:
//   - cmd: The cobra command (not currently used)
//   - args: Command line arguments (not currently used)
//
// The function gracefully handles errors by ignoring file removal failures,
// following the pattern of best-effort cleanup operations.
//
// Cache folder resolution:
//   - Default: $HOME/.ark-keyring/
//   - Override: ARK_BASIC_KEYRING_FOLDER environment variable
func (a *ArkCacheAction) runClearCacheAction(cmd *cobra.Command, args []string) {
	if keyring, err := common.NewArkKeyring("").GetKeyring(false); err == nil {
		if reflect.TypeOf(keyring) != reflect.TypeOf(&commoninternal.BasicKeyring{}) {
			commonargs.PrintNormal("Cache clear is only valid for basic keyring implementation at the moment")
			return
		}
		cacheFolderPath := os.ExpandEnv("$HOME") + string(os.PathSeparator) + commoninternal.DefaultBasicKeyringFolder
		if envPath := os.Getenv(commoninternal.ArkBasicKeyringFolderEnvVar); envPath != "" {
			cacheFolderPath = envPath
		}
		_ = os.Remove(cacheFolderPath + string(os.PathSeparator) + "keyring")
		_ = os.Remove(cacheFolderPath + string(os.PathSeparator) + "mac")
	}
}
