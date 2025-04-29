package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Iilun/survey/v2"
	"os"
	"regexp"

	"github.com/confluentinc/go-editor"
	"github.com/spf13/cobra"
	commonArgs "github.com/cyberark/ark-sdk-golang/pkg/common/args"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
)

// ArkProfilesAction is a struct that implements the ArkAction interface for managing profiles.
type ArkProfilesAction struct {
	*ArkBaseAction
	profilesLoader *profiles.ProfileLoader
}

// NewArkProfilesAction creates a new instance of ArkProfilesAction.
func NewArkProfilesAction(profilesLoader *profiles.ProfileLoader) *ArkProfilesAction {
	return &ArkProfilesAction{
		ArkBaseAction:  NewArkBaseAction(),
		profilesLoader: profilesLoader,
	}
}

// DefineAction Defines the CLI `profile` action, and adds actions for managing multiple profiles.
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

func (a *ArkProfilesAction) runListAction(cmd *cobra.Command, args []string) {
	// Start by loading all the profiles
	loadedProfiles, err := (*a.profilesLoader).LoadAllProfiles()
	if err != nil || len(loadedProfiles) == 0 {
		commonArgs.PrintWarning("No loadedProfiles were found")
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
		commonArgs.PrintSuccess(string(data))
	} else {
		names := []string{}
		for _, p := range loadedProfiles {
			names = append(names, p.ProfileName)
		}
		data, _ := json.MarshalIndent(names, "", "  ")
		commonArgs.PrintSuccess(string(data))
	}
}

func (a *ArkProfilesAction) runShowAction(cmd *cobra.Command, args []string) {
	profileName, _ := cmd.Flags().GetString("profile-name")
	if profileName == "" {
		profileName = profiles.DeduceProfileName("")
	}

	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil {
		commonArgs.PrintWarning(fmt.Sprintf("No profile was found for the name %s", profileName))
		return
	}

	data, _ := json.MarshalIndent(profile, "", "  ")
	commonArgs.PrintSuccess(string(data))
}

func (a *ArkProfilesAction) runDeleteAction(cmd *cobra.Command, args []string) {
	profileName, _ := cmd.Flags().GetString("profile-name")
	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil || profile == nil {
		commonArgs.PrintWarning(fmt.Sprintf("No profile was found for the name %s", profileName))
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

func (a *ArkProfilesAction) runCloneAction(cmd *cobra.Command, args []string) {
	profileName, _ := cmd.Flags().GetString("profile-name")
	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil {
		commonArgs.PrintWarning(fmt.Sprintf("No profile was found for the name %s", profileName))
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

func (a *ArkProfilesAction) runAddAction(cmd *cobra.Command, args []string) {
	profilePath, _ := cmd.Flags().GetString("profile-path")
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		commonArgs.PrintWarning(fmt.Sprintf("Profile path [%s] does not exist, ignoring", profilePath))
		return
	}
	if _, err := os.Stat(profilePath); err == nil {
		data, err := os.ReadFile(profilePath)
		if err != nil {
			commonArgs.PrintFailure(fmt.Sprintf("Profile path [%s] failed to be read, aborting", profilePath))
			return
		}
		var profile *models.ArkProfile
		if err = json.Unmarshal(data, &profile); err != nil {
			commonArgs.PrintFailure(fmt.Sprintf("Profile path [%s] failed to be parsed, aborting", profilePath))
			return
		}
		err = (*a.profilesLoader).SaveProfile(profile)
		if err != nil {
			return
		}
	}
	commonArgs.PrintFailure(fmt.Sprintf("Profile path [%s] does not exist", profilePath))
}

func (a *ArkProfilesAction) runEditAction(cmd *cobra.Command, args []string) {
	profileName, _ := cmd.Flags().GetString("profile-name")
	if profileName == "" {
		profileName = profiles.DeduceProfileName("")
	}

	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil {
		commonArgs.PrintWarning(fmt.Sprintf("No profile was found for the name %s", profileName))
		return
	}
	edit := editor.NewEditor()
	data, err := json.Marshal(profile)
	if err != nil {
		commonArgs.PrintFailure(fmt.Sprintf("Failed to marshal profile: %s", err))
		return
	}
	edited, path, err := edit.LaunchTempFile(fmt.Sprintf("%s-temp.json", profile.ProfileName), bytes.NewBufferString(string(data)))
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			commonArgs.PrintWarning(fmt.Sprintf("Failed to remove temp file: %s", err))
		}
	}(path)
	if err != nil {
		commonArgs.PrintFailure(fmt.Sprintf("Failed to launch editor: %s", err))
		return
	}
	err = json.Unmarshal(edited, &profile)
	if err != nil {
		commonArgs.PrintFailure(fmt.Sprintf("Failed to unmarshal edited profile: %s", err))
		return
	}
	err = (*a.profilesLoader).SaveProfile(profile)
	if err != nil {
		commonArgs.PrintWarning(fmt.Sprintf("Failed to save edited profile: %s", err))
		return
	}
}
