---
title: Simple SDK Workflow
description: Simple SDK Workflow
---

# Simple SDK Workflow
Let's say we want to generate a short lived password for SIA connection

The following example shows how to do that using the SDK

## Short lived password example
```go
package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	ssomodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/sso"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso"
	"os"
)

func main() {
	// Perform authentication using ArkISPAuth to the platform
	// First, create an ISP authentication class
	// Afterwards, perform the authentication
	ispAuth := auth.NewArkISPAuth(false)
	_, err := ispAuth.Authenticate(
		nil,
		&authmodels.ArkAuthProfile{
			Username:           "user@cyberark.cloud.12345",
			AuthMethod:         authmodels.Identity,
			AuthMethodSettings: &authmodels.IdentityArkAuthMethodSettings{},
		},
		&authmodels.ArkSecret{
			Secret: os.Getenv("ARK_SECRET"),
		},
		false,
		false,
	)
	if err != nil {
		panic(err)
	}

	// Create an SSO service from the authenticator above
	ssoService, err := sso.NewArkSIASSOService(ispAuth)
	if err != nil {
		panic(err)
	}

	// Generate a short-lived password
	ssoPassword, err := ssoService.ShortLivedPassword(
		&ssomodels.ArkSIASSOGetShortLivedPassword{},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", ssoPassword)
}
```

We first start by importing the required packages. The `auth` package is used to authenticate to the platform, and the `sso` package is used to generate a short-lived password.

We then create an instance of `ArkISPAuth` and call the `Authenticate` method to authenticate to the platform. The `Authenticate` method takes in the username, authentication method and relevant authentication method settings, and password as parameters.

Once we have authenticated, we create an instance of `ArkSIASSOService` using the `ispAuth` instance. This service is used to generate a short-lived password.

Finally, we call the `ShortLivedPassword` method on the `ssoService` instance to generate a short-lived password. The generated password is then printed to the console.
