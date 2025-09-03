package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/Iilun/survey/v2"

	"github.com/confluentinc/go-editor"
	commonargs "github.com/cyberark/ark-sdk-golang/pkg/common/args"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/spf13/cobra"
)

// ArkProfilesAction is a struct that implements the ArkAction interface for managing CLI profiles.
//
// ArkProfilesAction provides functionality for managing multiple CLI configuration profiles
// including listing, showing, deleting, clearing, cloning, adding, and editing profiles.
// It handles profile operations through a ProfileLoader and provides both interactive
// and non-interactive modes for various operations.
//
// The action supports comprehensive profile management operations:
//   - List profiles with optional filtering by name pattern or auth type
//   - Show detailed profile information
//   - Delete specific profiles with confirmation prompts
//   - Clear all profiles with confirmation
//   - Clone profiles with automatic or custom naming
//   - Add profiles from file paths
//   - Edit profiles interactively using an external editor
type ArkProfilesAction struct {
	// ArkBaseAction provides common action functionality
	*ArkBaseAction
	// profilesLoader handles loading and saving of profile configurations
	profilesLoader *profiles.ProfileLoader
}

// NewArkProfilesAction creates a new instance of ArkProfilesAction.
//
// NewArkProfilesAction initializes a new ArkProfilesAction with the provided
// profile loader and an embedded ArkBaseAction for common CLI functionality.
// The profile loader is used for all profile operations including loading,
// saving, deleting, and managing profile configurations.
//
// Parameters:
//   - profilesLoader: A pointer to a ProfileLoader for handling profile operations
//
// Returns a new ArkProfilesAction instance ready for defining profile management commands.
//
// Example:
//
//	loader := profiles.NewProfileLoader()
//	profilesAction := NewArkProfilesAction(loader)
//	profilesAction.DefineAction(rootCmd)
func NewArkProfilesAction(profilesLoader *profiles.ProfileLoader) *ArkProfilesAction {
	return &ArkProfilesAction{
		ArkBaseAction:  NewArkBaseAction(),
		profilesLoader: profilesLoader,
	}
}

// DefineAction defines the CLI profiles action and adds profile management commands.
//
// DefineAction creates a "profiles" command that provides comprehensive profile management
// functionality. The command includes multiple subcommands for different profile operations,
// each with their own flags and functionality.
//
// The function creates the following subcommands:
//   - list: Lists all profiles with optional filtering by name pattern or auth type
//   - show: Shows detailed information for a specific profile
//   - delete: Deletes a specific profile with confirmation
//   - clear: Clears all profiles with confirmation
//   - clone: Clones an existing profile with optional renaming
//   - add: Adds a profile from a file path
//   - edit: Opens a profile for interactive editing
//
// Parameters:
//   - cmd: The parent cobra command to which the profiles command will be added
//
// Each subcommand includes appropriate flags for configuration and supports both
// interactive and non-interactive modes where applicable.
//
// Example:
//
//	profilesAction := NewArkProfilesAction(loader)
//	profilesAction.DefineAction(rootCmd)
//	// This adds: myapp profiles [list|show|delete|clear|clone|add|edit] [flags]
func (a *ArkProfilesAction) DefineAction(cmd *cobra.Command) {
	profileCmd := &cobra.Command{
		Use:   "profiles",
		Short: "Manage profiles",
	}
	profileCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		a.CommonActionsExecution(cmd, args)
	}
	a.CommonActionsConfiguration(profileCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		Run:   a.runListAction,
	}
	listCmd.Flags().StringP("name", "", "", "Profile name to filter with by wildcard")
	listCmd.Flags().StringP("auth-profile", "", "", "Filter profiles by auth types")
	listCmd.Flags().BoolP("all", "", false, "Whether to show all profiles data as well and not only their names")

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show a profile",
		Run:   a.runShowAction,
	}
	showCmd.Flags().StringP("profile-name", "", "", "Profile name to show, if not given, shows the current one")

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a specific profile",
		Run:   a.runDeleteAction,
	}
	deleteCmd.Flags().StringP("profile-name", "", "", "Profile name to delete")
	deleteCmd.Flags().BoolP("yes", "", false, "Whether to approve deletion non interactively")

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all profiles",
		Run:   a.runClearAction,
	}
	clearCmd.Flags().BoolP("yes", "", false, "Whether to approve clear non interactively")

	cloneCmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone a profile",
		Run:   a.runCloneAction,
	}
	cloneCmd.Flags().StringP("profile-name", "", "", "Profile name to clone")
	cloneCmd.Flags().StringP("new-profile-name", "", "", "New cloned profile name, if not given, will add _clone as part of the name")
	cloneCmd.Flags().BoolP("yes", "", false, "Whether to override existing profile if exists")

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a profile from a given path",
		Run:   a.runAddAction,
	}
	addCmd.Flags().StringP("profile-path", "", "", "Profile file path to be added")

	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a profile interactively",
		Run:   a.runEditAction,
	}
	editCmd.Flags().StringP("profile-name", "", "", "Profile name to edit, if not given, edits the current one")

	profileCmd.AddCommand(listCmd, showCmd, deleteCmd, clearCmd, cloneCmd, addCmd, editCmd)
	cmd.AddCommand(profileCmd)
}

// runListAction handles the profiles list command execution.
//
// runListAction loads all available profiles and displays them based on the provided
// filtering criteria and output format options. It supports filtering by name pattern
// using regex matching and by authentication profile type.
//
// The function performs the following operations:
//  1. Loads all profiles using the profile loader
//  2. Applies name-based filtering using regex pattern matching
//  3. Applies auth-profile filtering by checking profile auth configurations
//  4. Outputs either profile names only or full profile data based on --all flag
//
// Parameters:
//   - cmd: The cobra command containing flag values for filtering and output options
//   - args: Command line arguments (not currently used)
//
// Supported flags:
//   - name: Regex pattern to filter profile names
//   - auth-profile: Filter profiles by specific auth type
//   - all: Show full profile data instead of just names
//
// The function prints warnings if no profiles are found and outputs JSON-formatted
// results for successful operations.

func (a *ArkProfilesAction) runListAction(cmd *cobra.Command, args []string) {
	// Start by loading all the profiles
	loadedProfiles, err := (*a.profilesLoader).LoadAllProfiles()
	if err != nil || len(loadedProfiles) == 0 {
		commonargs.PrintWarning("No loadedProfiles were found")
		return
	}

	// Filter profiles
	name, _ := cmd.Flags().GetString("name")
	if name != "" {
		var filtered []*models.ArkProfile
		for _, p := range loadedProfiles {
			if matched, err := regexp.MatchString(name, p.ProfileName); err == nil && matched {
				filtered = append(filtered, p)
			}
		}
		loadedProfiles = filtered
	}

	authProfile, _ := cmd.Flags().GetString("auth-profile")
	if authProfile != "" {
		var filtered []*models.ArkProfile
		for _, p := range loadedProfiles {
			if _, ok := p.AuthProfiles[authProfile]; ok {
				filtered = append(filtered, p)
			}
		}
		loadedProfiles = filtered
	}

	// Print them based on request
	showAll, _ := cmd.Flags().GetBool("all")
	if showAll {
		data, _ := json.MarshalIndent(loadedProfiles, "", "  ")
		commonargs.PrintSuccess(string(data))
	} else {
		names := []string{}
		for _, p := range loadedProfiles {
			names = append(names, p.ProfileName)
		}
		data, _ := json.MarshalIndent(names, "", "  ")
		commonargs.PrintSuccess(string(data))
	}
}

// runShowAction handles the profiles show command execution.
//
// runShowAction displays detailed information for a specific profile. If no profile
// name is provided, it uses the default profile name deduction logic to determine
// which profile to show.
//
// Parameters:
//   - cmd: The cobra command containing the profile-name flag
//   - args: Command line arguments (not currently used)
//
// Supported flags:
//   - profile-name: Name of the profile to show (optional, defaults to current profile)
//
// The function prints a warning if the specified profile is not found, otherwise
// it outputs the profile data in JSON format.
func (a *ArkProfilesAction) runShowAction(cmd *cobra.Command, args []string) {
	profileName, _ := cmd.Flags().GetString("profile-name")
	if profileName == "" {
		profileName = profiles.DeduceProfileName("")
	}

	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil {
		commonargs.PrintWarning(fmt.Sprintf("No profile was found for the name %s", profileName))
		return
	}

	data, _ := json.MarshalIndent(profile, "", "  ")
	commonargs.PrintSuccess(string(data))
}

// runDeleteAction handles the profiles delete command execution.
//
// runDeleteAction deletes a specific profile after loading it and optionally
// prompting for user confirmation. The function includes safety checks to
// ensure the profile exists before attempting deletion.
//
// Parameters:
//   - cmd: The cobra command containing the profile-name and yes flags
//   - args: Command line arguments (not currently used)
//
// Supported flags:
//   - profile-name: Name of the profile to delete (required)
//   - yes: Skip confirmation prompt for non-interactive deletion
//
// The function prints warnings if the profile is not found and uses interactive
// confirmation unless the --yes flag is provided.
func (a *ArkProfilesAction) runDeleteAction(cmd *cobra.Command, args []string) {
	profileName, _ := cmd.Flags().GetString("profile-name")
	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil || profile == nil {
		commonargs.PrintWarning(fmt.Sprintf("No profile was found for the name %s", profileName))
		return
	}

	yes, _ := cmd.Flags().GetBool("yes")
	if !yes {
		confirm := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete profile %s?", profileName),
		}
		err := survey.AskOne(prompt, &confirm)
		if err != nil || !confirm {
			return
		}
	}

	err = (*a.profilesLoader).DeleteProfile(profileName)
	if err != nil {
		return
	}
}

// runClearAction handles the profiles clear command execution.
//
// runClearAction clears all profiles after optionally prompting for user
// confirmation. This is a destructive operation that removes all stored
// profile configurations.
//
// Parameters:
//   - cmd: The cobra command containing the yes flag
//   - args: Command line arguments (not currently used)
//
// Supported flags:
//   - yes: Skip confirmation prompt for non-interactive clearing
//
// The function uses interactive confirmation unless the --yes flag is provided
// to prevent accidental deletion of all profiles.
func (a *ArkProfilesAction) runClearAction(cmd *cobra.Command, args []string) {
	yes, _ := cmd.Flags().GetBool("yes")
	if !yes {
		confirm := false
		prompt := &survey.Confirm{
			Message: "Are you sure you want to clear all profiles?",
		}
		err := survey.AskOne(prompt, &confirm)
		if err != nil || !confirm {
			return
		}
	}
	err := (*a.profilesLoader).ClearAllProfiles()
	if err != nil {
		return
	}
}

// runCloneAction handles the profiles clone command execution.
//
// runCloneAction creates a copy of an existing profile with a new name. If no
// new name is provided, it automatically appends "_clone" to the original name.
// The function includes logic to handle name conflicts with existing profiles.
//
// Parameters:
//   - cmd: The cobra command containing profile-name, new-profile-name, and yes flags
//   - args: Command line arguments (not currently used)
//
// Supported flags:
//   - profile-name: Name of the profile to clone (required)
//   - new-profile-name: Name for the cloned profile (optional, defaults to original_clone)
//   - yes: Skip confirmation prompt when overwriting existing profiles
//
// The function prompts for confirmation if the target profile name already exists
// unless the --yes flag is provided.
func (a *ArkProfilesAction) runCloneAction(cmd *cobra.Command, args []string) {
	profileName, _ := cmd.Flags().GetString("profile-name")
	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil {
		commonargs.PrintWarning(fmt.Sprintf("No profile was found for the name %s", profileName))
		return
	}

	newProfileName, _ := cmd.Flags().GetString("new-profile-name")
	if newProfileName == "" {
		newProfileName = profileName + "_clone"
	}

	clonedProfile := profile
	clonedProfile.ProfileName = newProfileName

	if (*a.profilesLoader).ProfileExists(newProfileName) {
		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			confirm := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Profile %s already exists, do you want to override it?", newProfileName),
			}
			err := survey.AskOne(prompt, &confirm)
			if err != nil || !confirm {
				return
			}
		}
	}
	err = (*a.profilesLoader).SaveProfile(clonedProfile)
	if err != nil {
		return
	}
}

// runAddAction handles the profiles add command execution.
//
// runAddAction adds a profile from a specified file path by reading and parsing
// the JSON profile data. The function includes validation to ensure the file
// exists and contains valid profile data.
//
// Parameters:
//   - cmd: The cobra command containing the profile-path flag
//   - args: Command line arguments (not currently used)
//
// Supported flags:
//   - profile-path: File system path to the profile JSON file (required)
//
// The function validates file existence, reads the file content, unmarshals
// the JSON data into a profile structure, and saves it using the profile loader.
// It prints appropriate warnings and failures for various error conditions.
func (a *ArkProfilesAction) runAddAction(cmd *cobra.Command, args []string) {
	profilePath, _ := cmd.Flags().GetString("profile-path")
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		commonargs.PrintWarning(fmt.Sprintf("Profile path [%s] does not exist, ignoring", profilePath))
		return
	}
	if _, err := os.Stat(profilePath); err == nil {
		data, err := os.ReadFile(profilePath)
		if err != nil {
			commonargs.PrintFailure(fmt.Sprintf("Profile path [%s] failed to be read, aborting", profilePath))
			return
		}
		var profile *models.ArkProfile
		if err = json.Unmarshal(data, &profile); err != nil {
			commonargs.PrintFailure(fmt.Sprintf("Profile path [%s] failed to be parsed, aborting", profilePath))
			return
		}
		err = (*a.profilesLoader).SaveProfile(profile)
		if err != nil {
			return
		}
	}
	commonargs.PrintFailure(fmt.Sprintf("Profile path [%s] does not exist", profilePath))
}

// runEditAction handles the profiles edit command execution.
//
// runEditAction opens a profile for interactive editing using an external editor.
// The profile data is marshaled to JSON, opened in a temporary file for editing,
// and then parsed back and saved after the user completes editing.
//
// Parameters:
//   - cmd: The cobra command containing the profile-name flag
//   - args: Command line arguments (not currently used)
//
// Supported flags:
//   - profile-name: Name of the profile to edit (optional, defaults to current profile)
//
// The function creates a temporary JSON file, launches the configured editor,
// waits for the user to complete editing, parses the modified content, and
// saves the updated profile. It includes cleanup logic to remove temporary files
// and comprehensive error handling for various failure scenarios.
func (a *ArkProfilesAction) runEditAction(cmd *cobra.Command, args []string) {
	profileName, _ := cmd.Flags().GetString("profile-name")
	if profileName == "" {
		profileName = profiles.DeduceProfileName("")
	}

	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil {
		commonargs.PrintWarning(fmt.Sprintf("No profile was found for the name %s", profileName))
		return
	}
	edit := editor.NewEditor()
	data, err := json.Marshal(profile)
	if err != nil {
		commonargs.PrintFailure(fmt.Sprintf("Failed to marshal profile: %s", err))
		return
	}
	edited, path, err := edit.LaunchTempFile(fmt.Sprintf("%s-temp.json", profile.ProfileName), bytes.NewBufferString(string(data)))
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			commonargs.PrintWarning(fmt.Sprintf("Failed to remove temp file: %s", err))
		}
	}(path)
	if err != nil {
		commonargs.PrintFailure(fmt.Sprintf("Failed to launch editor: %s", err))
		return
	}
	err = json.Unmarshal(edited, &profile)
	if err != nil {
		commonargs.PrintFailure(fmt.Sprintf("Failed to unmarshal edited profile: %s", err))
		return
	}
	err = (*a.profilesLoader).SaveProfile(profile)
	if err != nil {
		commonargs.PrintWarning(fmt.Sprintf("Failed to save edited profile: %s", err))
		return
	}
}
