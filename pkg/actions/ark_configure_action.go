package actions

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/args"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/mitchellh/mapstructure"
	"github.com/octago/sflags"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ArkConfigureAction is a struct that implements the ArkAction interface for configuring the CLI profiles.
//
// ArkConfigureAction provides functionality for managing CLI configuration profiles,
// including both interactive and silent configuration modes. It handles profile
// creation, modification, and persistence, as well as configuring authentication
// profiles for different authenticators supported by the Ark SDK.
//
// The action supports dynamic flag generation based on available authenticators
// and their supported authentication methods, allowing for flexible configuration
// of different authentication backends.
type ArkConfigureAction struct {
	// ArkBaseAction provides common action functionality
	*ArkBaseAction
	// profilesLoader handles loading and saving of profile configurations
	profilesLoader *profiles.ProfileLoader
}

// NewArkConfigureAction creates a new instance of ArkConfigureAction.
//
// NewArkConfigureAction initializes a new ArkConfigureAction with the provided
// profile loader and an embedded ArkBaseAction for common CLI functionality.
// The profile loader is used for loading existing profiles and saving new
// or modified profile configurations.
//
// Parameters:
//   - profilesLoader: A pointer to a ProfileLoader for handling profile operations
//
// Returns a new ArkConfigureAction instance ready for defining configure commands.
//
// Example:
//
//	loader := profiles.NewProfileLoader()
//	configAction := NewArkConfigureAction(loader)
//	configAction.DefineAction(rootCmd)
func NewArkConfigureAction(profilesLoader *profiles.ProfileLoader) *ArkConfigureAction {
	return &ArkConfigureAction{
		ArkBaseAction:  NewArkBaseAction(),
		profilesLoader: profilesLoader,
	}
}

// DefineAction defines the CLI configure action and adds configuration management commands.
//
// DefineAction creates a "configure" command that allows users to set up and modify
// CLI profiles for the Ark SDK. The command supports both interactive and silent
// modes of operation, dynamically generates flags based on available authenticators,
// and handles complex profile configuration scenarios.
//
// The function performs the following setup:
//  1. Creates a "configure" command with appropriate usage and help text
//  2. Sets up common action configuration through embedded ArkBaseAction
//  3. Dynamically generates flags for ArkProfile struct fields
//  4. Iterates through supported authenticators and adds their specific flags
//  5. Generates flags for authentication methods and their settings
//  6. Handles flag conflicts by filtering duplicate flags
//
// Parameters:
//   - cmd: The parent cobra command to which the configure command will be added
//
// The configure command supports extensive flag generation including:
//   - Profile-level settings (from models.ArkProfile)
//   - Authenticator enablement flags (work-with-<authenticator>)
//   - Authentication method selection flags
//   - Authenticator-specific configuration flags
//   - Authentication method-specific settings flags
//
// Example:
//
//	configAction := NewArkConfigureAction(loader)
//	configAction.DefineAction(rootCmd)
//	// This adds: myapp configure [flags]
func (a *ArkConfigureAction) DefineAction(cmd *cobra.Command) {
	confCmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure the CLI",
		Run:   a.runConfigureAction,
	}
	confCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		a.CommonActionsExecution(cmd, args)
	}
	a.CommonActionsConfiguration(confCmd)

	// Add the profile settings to the arguments
	err := gpflag.ParseTo(&models.ArkProfile{}, confCmd.Flags())
	if err != nil {
		a.logger.Error("Error parsing flags to ArkProfile %v", err)
		panic(err)
	}

	// Add the supported authenticator settings and whether to work with them or not
	for _, authenticator := range auth.SupportedAuthenticatorsList {
		confCmd.Flags().Bool(
			"work-with-"+strings.Replace(authenticator.AuthenticatorName(), "_", "-", -1),
			false,
			"Whether to work with "+authenticator.AuthenticatorHumanReadableName()+" services",
		)
		if len(authenticator.SupportedAuthMethods()) > 1 {
			confCmd.Flags().String(
				strings.Replace(authenticator.AuthenticatorName(), "_", "-", -1)+"-auth-method",
				string(authmodels.Default),
				"Authentication method for "+authenticator.AuthenticatorHumanReadableName(),
			)
		}
		// Add the rest of the ark auth profile params
		err = gpflag.ParseTo(&authmodels.ArkAuthProfile{}, confCmd.Flags(), sflags.Prefix(strings.Replace(authenticator.AuthenticatorName(), "_", "-", -1)+"-"))
		if err != nil {
			a.logger.Error("Error parsing flags to ArkAuthProfile %v", err)
			panic(err)
		}

		// Add the supported authentication methods settings of the authenticators
		for _, authMethod := range authenticator.SupportedAuthMethods() {
			authSettings := authmodels.ArkAuthMethodSettingsMap[authMethod]
			flags, err := sflags.ParseStruct(authSettings, sflags.Prefix(strings.Replace(authenticator.AuthenticatorName(), "_", "-", -1)+"-"))
			if err != nil {
				a.logger.Error("Error parsing flags to ArkAuthMethod settings %v", err)
				panic(err)
			}
			existingFlags := make([]string, 0)
			for _, flag := range flags {
				confCmd.Flags().VisitAll(func(f *pflag.Flag) {
					if f.Name == flag.Name {
						existingFlags = append(existingFlags, flag.Name)
					}
				})
			}
			filteredFlags := make([]*sflags.Flag, 0)
			for _, flag := range flags {
				if !slices.Contains(existingFlags, flag.Name) {
					filteredFlags = append(filteredFlags, flag)
				}
			}
			gpflag.GenerateTo(filteredFlags, confCmd.Flags())
		}
	}
	cmd.AddCommand(confCmd)
}

// runInteractiveConfigureAction handles interactive profile configuration through prompts.
//
// runInteractiveConfigureAction provides an interactive mode for configuring CLI profiles
// where users are prompted for each configuration value. It loads existing profiles,
// presents current values as defaults, and guides users through the configuration
// process for both profile settings and authenticator configurations.
//
// The function performs the following operations:
//  1. Prompts for profile name with intelligent default deduction
//  2. Loads existing profile or creates new one if not found
//  3. Prompts for each profile field with existing values as defaults
//  4. Presents checkbox selection for which authenticators to configure
//  5. For each selected authenticator, prompts for authentication method and settings
//  6. Handles authentication method-specific configuration dynamically
//
// Parameters:
//   - cmd: The cobra command containing flag definitions and values
//   - configureArgs: Command line arguments (not currently used)
//
// Returns the configured ArkProfile and any error encountered during configuration.
//
// The function uses the args package for interactive prompting and supports
// complex scenarios like multiple authenticators and different authentication methods.
func (a *ArkConfigureAction) runInteractiveConfigureAction(cmd *cobra.Command, configureArgs []string) (*models.ArkProfile, error) {
	profileName, err := args.GetArg(
		cmd,
		"profile-name",
		"Profile Name",
		profiles.DeduceProfileName(""),
		false,
		true,
		false,
	)
	if err != nil {
		return nil, err
	}
	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil || profile == nil {
		profile = &models.ArkProfile{
			ProfileName:  profileName,
			AuthProfiles: make(map[string]*authmodels.ArkAuthProfile),
		}
	}
	flags, err := sflags.ParseStruct(profile)
	if err != nil {
		return nil, err
	}
	answers := map[string]interface{}{}
	err = mapstructure.Decode(profile, &answers)
	if err != nil {
		return nil, err
	}
	for _, flag := range flags {
		if flag.Name == "profile-name" {
			continue
		}
		val, err := cmd.Flags().GetString(flag.Name)
		if err != nil {
			return nil, err
		}
		if val == "" {
			existingVal, ok := answers[strings.Replace(flag.Name, "-", "_", -1)]
			if ok {
				val = existingVal.(string)
			}
		}
		val, err = args.GetArg(
			cmd,
			flag.Name,
			flag.Usage,
			val,
			false,
			true,
			true,
		)
		answers[strings.Replace(flag.Name, "-", "_", -1)] = val
	}
	err = mapstructure.Decode(answers, &profile)

	var workWithAuthenticators []string
	if len(auth.SupportedAuthenticatorsList) == 1 {
		workWithAuthenticators = []string{auth.SupportedAuthenticatorsList[0].AuthenticatorHumanReadableName()}
	} else {
		workWithAuthenticators, err = args.GetCheckboxArgs(
			cmd,
			func() []string {
				keys := make([]string, len(auth.SupportedAuthenticatorsList))
				for i, a := range auth.SupportedAuthenticatorsList {
					keys[i] = fmt.Sprintf("work_with_%s", strings.Replace(a.AuthenticatorName(), "-", "_", -1))
				}
				return keys
			}(),
			"Which authenticators would you like to connect to",
			func() []string {
				names := make([]string, len(auth.SupportedAuthenticatorsList))
				for i, a := range auth.SupportedAuthenticatorsList {
					names[i] = a.AuthenticatorHumanReadableName()
				}
				return names
			}(),
			func() map[string]string {
				existingVals := make(map[string]string)
				for _, a := range auth.SupportedAuthenticatorsList {
					if _, exists := profile.AuthProfiles[a.AuthenticatorName()]; exists {
						existingVals[fmt.Sprintf("work_with_%s", strings.Replace(a.AuthenticatorName(), "-", "_", -1))] = a.AuthenticatorHumanReadableName()
					}
				}
				return existingVals
			}(),
			true,
		)
		if err != nil {
			return nil, err
		}
	}
	for _, authenticator := range auth.SupportedAuthenticatorsList {
		authProfile, ok := profile.AuthProfiles[authenticator.AuthenticatorName()]
		if !ok || authProfile == nil {
			authProfile = &authmodels.ArkAuthProfile{}
		}

		if slices.Contains(workWithAuthenticators, authenticator.AuthenticatorHumanReadableName()) {
			args.PrintSuccessBright(fmt.Sprintf("\n◉ Configuring %s", authenticator.AuthenticatorHumanReadableName()))

			var authMethod authmodels.ArkAuthMethod
			if len(authenticator.SupportedAuthMethods()) > 1 {
				authMethodStr, err := args.GetSwitchArg(
					cmd,
					fmt.Sprintf("%s_auth_method", strings.Replace(authenticator.AuthenticatorName(), "-", "_", -1)),
					"Authentication Method",
					func() []string {
						methods := make([]string, len(authenticator.SupportedAuthMethods()))
						for i, m := range authenticator.SupportedAuthMethods() {
							methods[i] = authmodels.ArkAuthMethodsDescriptionMap[m]
						}
						return methods
					}(),
					func() string {
						if _, exists := profile.AuthProfiles[authenticator.AuthenticatorName()]; exists {
							return authmodels.ArkAuthMethodsDescriptionMap[authProfile.AuthMethod]
						}
						authMethod, _ := authenticator.DefaultAuthMethod()
						return authmodels.ArkAuthMethodsDescriptionMap[authMethod]
					}(),
					true,
				)
				if err != nil {
					return nil, err
				}
				for key, val := range authmodels.ArkAuthMethodsDescriptionMap {
					if val == authMethodStr {
						authMethod = key
						break
					}
				}
				if authMethod == authmodels.Default {
					authMethod, _ = authenticator.DefaultAuthMethod()
				}
			} else {
				authMethod, _ = authenticator.DefaultAuthMethod()
			}
			authPrefix := strings.Replace(authenticator.AuthenticatorName(), "_", "-", -1) + "-"
			authProfileFlags, err := sflags.ParseStruct(authProfile, sflags.Prefix(authPrefix))
			if err != nil {
				return nil, err
			}
			authProfileAnswers := map[string]interface{}{}
			err = mapstructure.Decode(authProfile, &authProfileAnswers)
			if err != nil {
				return nil, err
			}
			for _, flag := range authProfileFlags {
				val, err := cmd.Flags().GetString(flag.Name)
				if err != nil {
					return nil, err
				}
				if val == "" {
					existingVal, ok := authProfileAnswers[strings.Replace(strings.TrimPrefix(flag.Name, authPrefix), "-", "_", -1)]
					if ok {
						val = existingVal.(string)
					}
				}
				val, err = args.GetArg(
					cmd,
					flag.Name,
					flag.Usage,
					val,
					false,
					true,
					true,
				)
				authProfileAnswers[strings.Replace(strings.TrimPrefix(flag.Name, authPrefix), "-", "_", -1)] = val
			}
			err = mapstructure.Decode(authProfileAnswers, &authProfile)
			if err != nil {
				return nil, err
			}
			var methodSettings interface{}
			if authMethod != authProfile.AuthMethod {
				methodSettings = authmodels.ArkAuthMethodSettingsMap[authMethod]
			} else {
				methodSettings = authmodels.ArkAuthMethodSettingsMap[authMethod]
				err = mapstructure.Decode(authProfile.AuthMethodSettings, methodSettings)
				if err != nil {
					return nil, err
				}
			}
			methodSettingsFlags, err := sflags.ParseStruct(methodSettings, sflags.Prefix(strings.Replace(authenticator.AuthenticatorName(), "_", "-", -1)+"-"))
			if err != nil {
				return nil, err
			}
			methodSettingsAnswers := map[string]interface{}{}
			err = mapstructure.Decode(methodSettings, &methodSettingsAnswers)
			if err != nil {
				return nil, err
			}
			for _, flag := range methodSettingsFlags {
				if flag.Value.Type() == "bool" {
					val, err := cmd.Flags().GetBool(flag.Name)
					if err != nil {
						return nil, err
					}
					val, err = args.GetBoolArg(
						cmd,
						flag.Name,
						flag.Usage,
						nil,
						true,
					)
					methodSettingsAnswers[strings.Replace(strings.TrimPrefix(flag.Name, authPrefix), "-", "_", -1)] = val
				} else {
					val, err := cmd.Flags().GetString(flag.Name)
					if err != nil {
						return nil, err
					}
					if val == "" {
						existingVal, ok := methodSettingsAnswers[strings.Replace(strings.TrimPrefix(flag.Name, authPrefix), "-", "_", -1)]
						if ok {
							val = existingVal.(string)
						}
					}
					val, err = args.GetArg(
						cmd,
						flag.Name,
						flag.Usage,
						val,
						false,
						true,
						true,
					)
					methodSettingsAnswers[strings.Replace(strings.TrimPrefix(flag.Name, authPrefix), "-", "_", -1)] = val
				}
			}
			err = mapstructure.Decode(methodSettingsAnswers, methodSettings)
			if err != nil {
				return nil, err
			}

			authProfile.AuthMethod = authMethod
			authProfile.AuthMethodSettings = methodSettings
			profile.AuthProfiles[authenticator.AuthenticatorName()] = authProfile
		} else if _, exists := profile.AuthProfiles[authenticator.AuthenticatorName()]; exists {
			delete(profile.AuthProfiles, authenticator.AuthenticatorName())
		}
	}
	return profile, nil
}

// runSilentConfigureAction handles non-interactive profile configuration using command flags.
//
// runSilentConfigureAction processes configuration in silent mode where all configuration
// values are provided through command line flags rather than interactive prompts. It
// loads existing profiles, merges flag values, and validates the configuration for
// completeness and correctness.
//
// The function performs the following operations:
//  1. Extracts profile name from flags or deduces default
//  2. Loads existing profile or creates new one if not found
//  3. Processes all changed flags and merges them into the profile
//  4. Iterates through authenticators based on work-with-* flags
//  5. Validates required fields (like username for credential-based auth methods)
//  6. Processes authenticator-specific flags and settings
//  7. Removes authenticator profiles that are no longer configured
//
// Parameters:
//   - cmd: The cobra command containing parsed flags and values
//   - args: Command line arguments (not currently used)
//
// Returns the configured ArkProfile and any error encountered during configuration.
//
// The function requires all necessary configuration to be provided via flags and
// will return errors for missing required fields or invalid configurations.
func (a *ArkConfigureAction) runSilentConfigureAction(cmd *cobra.Command, args []string) (*models.ArkProfile, error) {
	// Load the profile based on the CLI params and merge the rest of the params
	profileName, err := cmd.Flags().GetString("profile-name")
	if err != nil {
		return nil, err
	}
	if profileName == "" {
		profileName = profiles.DeduceProfileName("")
	}
	profile, err := (*a.profilesLoader).LoadProfile(profileName)
	if err != nil || profile == nil {
		profile = &models.ArkProfile{
			ProfileName:  profileName,
			AuthProfiles: make(map[string]*authmodels.ArkAuthProfile),
		}
	}
	flags := map[string]interface{}{}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			flags[strings.Replace(f.Name, "-", "_", -1)] = f.Value.String()
		}
	})
	err = mapstructure.Decode(flags, &profile)
	if err != nil {
		return nil, err
	}

	// Load the authenticators
	for _, authenticator := range auth.SupportedAuthenticatorsList {
		authProfile, ok := profile.AuthProfiles[authenticator.AuthenticatorName()]
		if !ok || authProfile == nil {
			authProfile = &authmodels.ArkAuthProfile{}
		}
		workWithAuth, err := cmd.Flags().GetBool(fmt.Sprintf("work-with-%s", authenticator.AuthenticatorName()))
		if err != nil {
			continue
		}
		if workWithAuth {
			var authMethod authmodels.ArkAuthMethod
			if len(authenticator.SupportedAuthMethods()) > 1 {
				authMethodStr, err := cmd.Flags().GetString(fmt.Sprintf("%s-auth-method", authenticator.AuthenticatorName()))
				if err != nil {
					return nil, err
				}
				authMethod = authmodels.ArkAuthMethod(authMethodStr)
				if authMethod == authmodels.Default {
					authMethod, _ = authenticator.DefaultAuthMethod()
				} else if !slices.Contains(authenticator.SupportedAuthMethods(), authMethod) {
					return nil, fmt.Errorf("auth method %s is not supported by %s", authMethod, authenticator.AuthenticatorHumanReadableName())
				}
			} else {
				authMethod, _ = authenticator.DefaultAuthMethod()
			}
			authSpecificFlags := map[string]interface{}{}
			authPrefix := strings.Replace(authenticator.AuthenticatorName(), "_", "-", -1) + "-"
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				if f.Changed && strings.HasPrefix(f.Name, authPrefix) {
					authSpecificFlags[strings.Replace(strings.TrimPrefix(f.Name, authPrefix), "-", "_", -1)] = f.Value.String()
				}
			})
			err = mapstructure.Decode(authSpecificFlags, &authProfile)
			if err != nil {
				return nil, err
			}
			if slices.Contains(authmodels.ArkAuthMethodsRequireCredentials, authMethod) && authProfile.Username == "" {
				return nil, fmt.Errorf("missing username for authenticator [%s]", authenticator.AuthenticatorHumanReadableName())
			}

			var methodSettings interface{}
			if authMethod != authProfile.AuthMethod {
				methodSettings = authmodels.ArkAuthMethodSettingsMap[authMethod]
			} else {
				methodSettings = authmodels.ArkAuthMethodSettingsMap[authMethod]
				err = mapstructure.Decode(authProfile.AuthMethodSettings, methodSettings)
				if err != nil {
					return nil, err
				}
			}
			err = mapstructure.Decode(authSpecificFlags, methodSettings)
			if err != nil {
				return nil, err
			}

			authProfile.AuthMethod = authMethod
			authProfile.AuthMethodSettings = methodSettings
			profile.AuthProfiles[authenticator.AuthenticatorName()] = authProfile
		} else if _, exists := profile.AuthProfiles[authenticator.AuthenticatorName()]; exists {
			delete(profile.AuthProfiles, authenticator.AuthenticatorName())
		}
	}
	return profile, nil
}

// runConfigureAction is the main entry point for the configure command execution.
//
// runConfigureAction determines whether to run in interactive or silent mode based
// on the current CLI environment settings, executes the appropriate configuration
// method, saves the resulting profile, and displays the configuration results to
// the user.
//
// The function performs the following operations:
//  1. Checks if the CLI is running in interactive mode
//  2. Calls either runInteractiveConfigureAction or runSilentConfigureAction
//  3. Saves the configured profile using the profile loader
//  4. Marshals the profile to JSON for display
//  5. Prints success messages and profile location information
//
// Parameters:
//   - cmd: The cobra command containing configuration and flags
//   - configureArgs: Command line arguments passed to the configure command
//
// The function handles errors by logging them and panicking for critical failures
// like configuration errors, following the CLI pattern of failing fast for
// configuration issues that prevent proper operation.
//
// Example output:
//
//	{
//	  "profile_name": "default",
//	  "auth_profiles": {...}
//	}
//	Profile has been saved to /home/user/.ark-profiles
func (a *ArkConfigureAction) runConfigureAction(cmd *cobra.Command, configureArgs []string) {
	var profile *models.ArkProfile
	var err error
	if common.IsInteractive() {
		profile, err = a.runInteractiveConfigureAction(cmd, configureArgs)
	} else {
		profile, err = a.runSilentConfigureAction(cmd, configureArgs)
	}
	if err != nil {
		a.logger.Error("Error configuring ark profile %v", err)
		panic(err)
	}
	err = (*a.profilesLoader).SaveProfile(profile)
	if err != nil {
		a.logger.Error("Error saving ark profile %v", err)
		return
	}
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		a.logger.Error("Error serializing ark profile %v", err)
		return
	}
	args.PrintSuccess(string(data))
	args.PrintSuccessBright(fmt.Sprintf("Profile has been saved to %s", profiles.GetProfilesFolder()))
}
