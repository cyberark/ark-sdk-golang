package actions

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/auth/identity"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/args"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"slices"
)

// ArkLoginAction is a struct that implements the ArkAction interface for login action.
type ArkLoginAction struct {
	*ArkBaseAction
	profilesLoader *profiles.ProfileLoader
}

// NewArkLoginAction creates a new instance of ArkLoginAction.
func NewArkLoginAction(profilesLoader *profiles.ProfileLoader) *ArkLoginAction {
	return &ArkLoginAction{
		ArkBaseAction:  NewArkBaseAction(),
		profilesLoader: profilesLoader,
	}
}

// DefineAction defines the CLI `login` action, and adds the login function.
func (a *ArkLoginAction) DefineAction(cmd *cobra.Command) {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the system",
		Run:   a.runLoginAction,
	}
	loginCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		a.CommonActionsExecution(cmd, args)
	}
	a.CommonActionsConfiguration(loginCmd)

	loginCmd.Flags().String("profile-name", profiles.DefaultProfileName(), "Profile name to load")
	loginCmd.Flags().Bool("force", false, "Whether to force login even though token has not expired yet")
	loginCmd.Flags().Bool("no-shared-secrets", false, "Do not share secrets between different authenticators with the same username")
	loginCmd.Flags().Bool("show-tokens", false, "Print out tokens as well if not silent")
	loginCmd.Flags().Bool("refresh-auth", false, "If a cache exists, will also try to refresh it")

	for _, authenticator := range auth.SupportedAuthenticatorsList {
		loginCmd.Flags().String(fmt.Sprintf("%s-username", authenticator.AuthenticatorName()), "", fmt.Sprintf("Username to authenticate with to %s", authenticator.AuthenticatorHumanReadableName()))
		loginCmd.Flags().String(fmt.Sprintf("%s-secret", authenticator.AuthenticatorName()), "", fmt.Sprintf("Secret to authenticate with to %s", authenticator.AuthenticatorHumanReadableName()))
	}

	cmd.AddCommand(loginCmd)
}

func (a *ArkLoginAction) runLoginAction(cmd *cobra.Command, loginArgs []string) {
	a.CommonActionsExecution(cmd, loginArgs)

	profileName, _ := cmd.Flags().GetString("profile-name")
	profile, err := (*a.profilesLoader).LoadProfile(profiles.DeduceProfileName(profileName))
	if err != nil || profile == nil {
		args.PrintFailure("Please configure a profile before trying to login")
		return
	}

	sharedSecretsMap := make(map[authmodels.ArkAuthMethod][][2]string)
	tokensMap := make(map[string]*authmodels.ArkToken)

	for authenticatorName, authProfile := range profile.AuthProfiles {
		authenticator := auth.SupportedAuthenticators[authenticatorName]
		force, _ := cmd.Flags().GetBool("force")
		refreshAuth, _ := cmd.Flags().GetBool("refresh-auth")
		if authenticator.IsAuthenticated(profile) && !force {
			if refreshAuth {
				_, err := authenticator.LoadAuthentication(profile, true)
				if err == nil {
					args.PrintSuccess(fmt.Sprintf("%s Authentication Refreshed", authenticator.AuthenticatorHumanReadableName()))
					continue
				}
				a.logger.Info(fmt.Sprintf("%s Failed to refresh token, performing normal login [%s]", authenticator.AuthenticatorHumanReadableName(), err))
			} else {
				args.PrintSuccess(fmt.Sprintf("%s Already Authenticated", authenticator.AuthenticatorHumanReadableName()))
				continue
			}
		}
		secretStr, _ := cmd.Flags().GetString(fmt.Sprintf("%s-secret", authenticatorName))
		secret := &authmodels.ArkSecret{Secret: secretStr}
		userName, _ := cmd.Flags().GetString(fmt.Sprintf("%s-username", authenticatorName))
		if userName == "" {
			userName = authProfile.Username
		}
		if common.IsInteractive() && slices.Contains(authmodels.ArkAuthMethodsRequireCredentials, authProfile.AuthMethod) {
			authProfile.Username, err = args.GetArg(
				cmd,
				fmt.Sprintf("%s-username", authenticatorName),
				fmt.Sprintf("%s Username", authenticator.AuthenticatorHumanReadableName()),
				userName,
				false,
				true,
				false,
			)
			if slices.Contains(authmodels.ArkAuthMethodSharableCredentials, authProfile.AuthMethod) && len(sharedSecretsMap[authProfile.AuthMethod]) > 0 && !viper.GetBool("no-shared-secrets") {
				for _, s := range sharedSecretsMap[authProfile.AuthMethod] {
					if s[0] == authProfile.Username {
						secret = &authmodels.ArkSecret{Secret: s[1]}
						break
					}
				}
			} else {
				if authenticatorName == "isp" &&
					authProfile.AuthMethod == authmodels.Identity &&
					!identity.IsPasswordRequired(authProfile.Username,
						authProfile.AuthMethodSettings.(*authmodels.IdentityArkAuthMethodSettings).IdentityURL,
						authProfile.AuthMethodSettings.(*authmodels.IdentityArkAuthMethodSettings).IdentityTenantSubdomain) {
					secret = &authmodels.ArkSecret{Secret: ""}
				} else {
					secretStr, err = args.GetArg(
						cmd,
						fmt.Sprintf("%s-secret", authenticatorName),
						fmt.Sprintf("%s Secret", authenticator.AuthenticatorHumanReadableName()),
						secretStr,
						true,
						false,
						false,
					)
					if err != nil {
						args.PrintFailure(fmt.Sprintf("Failed to get %s secret: %s", authenticatorName, err))
						return
					}
					secret = &authmodels.ArkSecret{Secret: secretStr}
				}
			}
		} else if !common.IsInteractive() && slices.Contains(authmodels.ArkAuthMethodsRequireCredentials, authProfile.AuthMethod) && secret.Secret == "" {
			args.PrintFailure(fmt.Sprintf("%s-secret argument is required if authenticating to %s", authenticatorName, authenticator.AuthenticatorHumanReadableName()))
			return
		}

		token, err := authenticator.Authenticate(profile, nil, secret, force, refreshAuth)
		if err != nil {
			args.PrintFailure(fmt.Sprintf("Failed to authenticate with %s: %s", authenticator.AuthenticatorHumanReadableName(), err))
			return
		}

		noSharedSecrets, _ := cmd.Flags().GetBool("no-shared-secrets")
		if !noSharedSecrets && slices.Contains(authmodels.ArkAuthMethodSharableCredentials, authProfile.AuthMethod) {
			sharedSecretsMap[authProfile.AuthMethod] = append(sharedSecretsMap[authProfile.AuthMethod], [2]string{authProfile.Username, secret.Secret})
		}
		tokensMap[authenticator.AuthenticatorHumanReadableName()] = token
	}
	showTokens, _ := cmd.Flags().GetBool("show-tokens")
	if !showTokens && len(tokensMap) > 0 {
		args.PrintSuccess("Login tokens are hidden")
	}

	for k, v := range tokensMap {
		if v.Metadata != nil {
			if _, ok := v.Metadata["cookies"]; ok {
				delete(v.Metadata, "cookies")
			}
		}
		tokenMap := make(map[string]interface{})
		data, _ := json.Marshal(v)
		_ = json.Unmarshal(data, &tokenMap)
		if !showTokens {
			delete(tokenMap, "token")
			delete(tokenMap, "refresh_token")
		}
		jsonData, _ := json.MarshalIndent(tokenMap, "", "  ")
		args.PrintSuccess(fmt.Sprintf("%s Token\n%s", k, jsonData))
	}
}
