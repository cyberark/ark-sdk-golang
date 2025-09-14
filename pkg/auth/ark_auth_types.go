package auth

import (
	"slices"

	"github.com/cyberark/ark-sdk-golang/pkg/models/auth"
)

var (
	// SupportedAuthenticatorsList is a list of supported authenticators.
	SupportedAuthenticatorsList = []ArkAuth{
		NewArkISPAuth(true),
	}

	// SupportedAuthenticators is a map of supported authenticators.
	SupportedAuthenticators = func() map[string]ArkAuth {
		authenticators := make(map[string]ArkAuth)
		for _, auth := range SupportedAuthenticatorsList {
			authenticators[auth.AuthenticatorName()] = auth
		}
		return authenticators
	}()

	// SupportedAuthMethods is a list of supported authentication methods.
	SupportedAuthMethods = func() []auth.ArkAuthMethod {
		authMethods := make([]auth.ArkAuthMethod, 0)
		for _, auth := range SupportedAuthenticatorsList {
			for _, method := range auth.SupportedAuthMethods() {
				if !slices.Contains(authMethods, method) {
					authMethods = append(authMethods, method)
				}
			}
		}
		return authMethods
	}()
)
