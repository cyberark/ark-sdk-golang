package actions

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/octago/sflags"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/args"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"slices"
	"strings"
)

// ArkConfigureAction is a struct that implements the ArkAction interface for configuring the CLI profiles.
type ArkConfigureAction struct {
	*ArkBaseAction
	profilesLoader *profiles.ProfileLoader
}

// NewArkConfigureAction Creates a new instance of ArkConfigureAction
func NewArkConfigureAction(profilesLoader *profiles.ProfileLoader) *ArkConfigureAction {
	return &ArkConfigureAction{
		ArkBaseAction:  NewArkBaseAction(),
		profilesLoader: profilesLoader,
	}
}

// DefineAction Defines the CLI `configure` action, and adds the clear cache function
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
