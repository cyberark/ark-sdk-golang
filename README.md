![Ark SDK Golang](https://github.com/cyberark/ark-sdk-golang/blob/main/assets/sdk.png)

<p align="center">
    <a href="https://actions-badge.atrox.dev/cyberark/ark-sdk-golang/goto?ref=main" alt="Build">
        <img src="https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fcyberark%2Fark-sdk-golang%2Fbadge%3Fref%3Dmain&style=flat" />
    </a>
    <a alt="Go Version">
        <img src="https://img.shields.io/github/go-mod/go-version/cyberark/ark-sdk-golang" />
    </a>
    <a href="https://github.com/cyberark/ark-sdk-golang/blob/main/LICENSE.txt" alt="License">
        <img src="https://img.shields.io/github/license/cyberark/ark-sdk-golang?style=flat" alt="License" />
    </a>
</p>

Ark SDK Golang
==============

📜[**Documentation**](https://cyberark.github.io/ark-sdk-golang/)

CyberArk's Official SDK and CLI for different services operations

## Features and Services
- [x] Extensive and Interactive CLI
- [x] Different Authenticators
    - [x] Identity Authentication Methods
    - [x] MFA Support for Identity
    - [x] Identity Security Platform
- [x] Ready to use SDK in Golang
- [x] Fully Interactive CLI comprising of 3 main actions
    - [x] Configure
    - [x] Login
    - [x] Exec
- [x] Services API
  - [x] SIA SSO Service
  - [x] SIA K8S Service
  - [x] SIA VM Secrets Service
  - [x] SIA DB Secrets Service
  - [x] SIA Target Sets Workspace Service
  - [x] SIA Access Service
  - [x] Connector Manager Service
  - [x] PCloud Accounts Service
  - [x] PCloud Safes Service
  - [x] Identity Directories Service
  - [x] Identity Roles Service
  - [x] Identity Users Service
- [x] Filesystem Inputs and Outputs for the CLI
- [x] Silent and Verbose logging
- [x] Profile Management and Authentication Caching


TL;DR
=====

## Enduser
![Ark SDK Enduser Usage](https://github.com/cyberark/ark-sdk-golang/blob/main/assets/ark_sdk_enduser_tldr.gif)


Installation
============

One can install the SDK via the following command:
```shell
go install github.com/cyberark/ark-sdk-golang
```

CLI Usage
============
Both the SDK and the CLI works with profiles

The profiles can be configured upon need and be used for the consecutive actions

The CLI has the following basic commands:
- <b>configure</b> - Configures profiles and their respective authentication methods
- <b>login</b> - Logs into the profile authentication methods
- <b>exec</b> - Executes different commands based on the supported services
- <b>profiles</b> - Manage multiple profiles on the machine
- <b>cache</b> - Manage the cache of the authentication methods


configure
---------
The configure command is used to create a profile to work on<br>
The profile consists of infomration regarding which authentication methods to use and what are their method settings, along with other related information such as MFA

How to run:
```shell
ark configure
```


The profiles are saved to ~/.ark_profiles

No arguments are required, and interactive questions will be asked

If you wish to only supply arguments in a silent fashion, --silent can be added along with the arugments

Usage:
```shell
Configure the CLI

Usage:
  ark configure [flags]

Flags:
      --allow-output                                    Allow stdout / stderr even when silent and not interactive
      --disable-cert-verification                       Disables certificate verification on HTTPS calls, unsafe!
  -h, --help                                            help for configure
      --isp-auth-method string                          Authentication method for Identity Security Platform (default "default")
      --isp-identity-application string                 Identity Application
      --isp-identity-authorization-application string   Service User Authorization Application
      --isp-identity-mfa-interactive                    Allow Interactive MFA
      --isp-identity-mfa-method string                  MFA Method to use by default [pf, sms, email, otp]
      --isp-identity-tenant-subdomain string            Identity Tenant Subdomain
      --isp-identity-url string                         Identity Url
      --isp-username string                             Username
      --log-level string                                Log level to use while verbose (default "INFO")
      --logger-style string                             Which verbose logger style to use (default "default")
      --profile-description string                      Profile Description
      --profile-name string                             The name of the profile to use
      --raw                                             Whether to raw output
      --silent                                          Silent execution, no interactiveness
      --trusted-cert string                             Certificate to use for HTTPS calls
      --verbose                                         Whether to verbose log
      --work-with-isp                                   Whether to work with Identity Security Platform services
```


login
-----
The login command is used to login to the authentication methods configured for the profile

You will be asked to write a password for each respective authentication method that supports password, and alongside that, any needed MFA prompt

Once the login is done, the access tokens are stored on the computer keystore for their lifetime

Once they are expired, a consecutive login will be required

How to run:
```shell
ark login
```

Usage:
```shell
Login to the system

Usage:
  ark login [flags]

Flags:
      --allow-output                Allow stdout / stderr even when silent and not interactive
      --disable-cert-verification   Disables certificate verification on HTTPS calls, unsafe!
      --force                       Whether to force login even though token has not expired yet
  -h, --help                        help for login
      --isp-secret string           Secret to authenticate with to Identity Security Platform
      --isp-username string         Username to authenticate with to Identity Security Platform
      --log-level string            Log level to use while verbose (default "INFO")
      --logger-style string         Which verbose logger style to use (default "default")
      --no-shared-secrets           Do not share secrets between different authenticators with the same username
      --profile-name string         Profile name to load (default "ark")
      --raw                         Whether to raw output
      --refresh-auth                If a cache exists, will also try to refresh it
      --show-tokens                 Print out tokens as well if not silent
      --silent                      Silent execution, no interactiveness
      --trusted-cert string         Certificate to use for HTTPS calls
      --verbose                     Whether to verbose log
```

Notes:

- You may disable certificate validation for login to different authenticators using the --disable-certificate-verification or supply a certificate to be used, not recommended to disable


exec
----
The exec command is used to execute various commands based on supported services for the fitting logged in authenticators

The following services and commands are supported:
- <b>sia</b> - Secure Infrastructure Access Services
  - <b>sso</b> - SIA SSO Management
  - <b>k8s</b> - SIA K8S Management
  - <b>workspaces</b> - SIA Workspaces Management
    - <b>target-sets</b> - SIA VM Target Sets Management
  - <b>secrets</b> - SIA Secrets Management
    - <b>vm</b> - SIA VM Secrets Management
  - <b>access</b> - SIA Access Management
- <b>cmgr</b> - Connector Manager
- <b>pcloud</b> - PCloud Service
  - <b>accounts</b> - PCloud Accounts Management
  - <b>safes</b> - PCloud Safes Management
- <b>identity</b> - Identity Service
  - <b>directories</b> - Identity Directories Management
  - <b>roles</b> - Identity Roles Management
  - <b>users</b> - Identity Users Management

Any command has its own subcommands, with respective arguments

For example, generating a short lived password
```shell
ark exec sia sso short-lived-password
```

Add SIA VM Target Set
```shell
ark exec sia workspaces target-sets add-target-set --name mydomain.com --type Domain
```

Add SIA VM Secret
```shell
ark exec sia secrets vm add-secret --secret-type ProvisionerUser --provisioner-username=myuser --provisioner-password=mypassword
```

List connector pools
```shell
ark exec exec cmgr list-pools
```

Get connector installation script
```shell
ark exec sia access connector-setup-script --connector-type ON-PREMISE --connector-os windows --connector-pool-id 588741d5-e059-479d-b4c4-3d821a87f012
```

Create a PCloud Safe
```shell
ark exec pcloud safes add-safe --safe-name=safe
```

Create a PCloud Account
```shell
ark exec pcloud accounts add-account --name account --safe-name safe --platform-id='UnixSSH' --username root --address 1.2.3.4 --secret-type=password --secret mypass
```

Retrieve a PCloud Account Credentials
```shell
ark exec pcloud accounts get-account-credentials --account-id 11_1
```

Create an Identity User
```shell
ark exec identity users create-user --roles "DpaAdmin" --username "myuser"
```

Create an Identity Role
```shell
ark exec identity roles create-role --role-name myrole
```

List all directories identities
```shell
ark exec identity directories list-directories-entities
```

Add SIA Database Secret
```shell
ark exec sia secrets db add-secret --secret-name mysecret --secret-type username_password --username user --password mypass
```
Delete SIA Database Secret
```shell
ark exec sia secrets db delete-secret --secret-name mysecret
```

Add SIA database
```shell
ark exec sia workspaces db add-database --name mydatabase --provider-engine aurora-mysql --read-write-endpoint myrds.com
```
Delete SIA database
```shell
ark exec sia workspaces db delete-database --id databaseid
```

You can view all of the commands via the --help for each respective exec action

Notes:

- You may disable certificate validation for login to different authenticators using the --disable-certificate-verification or supply a certificate to be used, not recommended to disable


Usafe Env Vars:
- ARK_PROFILE - Sets the profile to be used across the CLI
- ARK_DISABLE_CERTIFICATE_VERIFICATION - Disables certificate verification on REST API's


profiles
-------
As one may have multiple environments to manage, this would also imply that multiple profiles are required, either for multiple users in the same environment or multiple tenants

Therefore, the profiles command manages those profiles as a convenice set of methods

Using the profiles as simply running commands under:
```shell
ark profiles
```

Usage:
```shell
Manage profiles

Usage:
  ark profiles [command]

Available Commands:
  add         Add a profile from a given path
  clear       Clear all profiles
  clone       Clone a profile
  delete      Delete a specific profile
  edit        Edit a profile interactively
  list        List all profiles
  show        Show a profile

Flags:
      --allow-output                Allow stdout / stderr even when silent and not interactive
      --disable-cert-verification   Disables certificate verification on HTTPS calls, unsafe!
  -h, --help                        help for profiles
      --log-level string            Log level to use while verbose (default "INFO")
      --logger-style string         Which verbose logger style to use (default "default")
      --raw                         Whether to raw output
      --silent                      Silent execution, no interactiveness
      --trusted-cert string         Certificate to use for HTTPS calls
      --verbose                     Whether to verbose log

Use "ark profiles [command] --help" for more information about a command.
```


cache
-------
Use the cache command to manage the Ark data cached on your machine. Currently, you can only clear the filesystem cache (not data cached in the OS's keystore).


Using the cache as simply running commands under:
```shell
ark cache
```

Usage:
```shell
Manage cache

Usage:
  ark cache [command]

Available Commands:
  clear       Clears all profiles cache

Flags:
      --allow-output                Allow stdout / stderr even when silent and not interactive
      --disable-cert-verification   Disables certificate verification on HTTPS calls, unsafe!
  -h, --help                        help for cache
      --log-level string            Log level to use while verbose (default "INFO")
      --logger-style string         Which verbose logger style to use (default "default")
      --raw                         Whether to raw output
      --silent                      Silent execution, no interactiveness
      --trusted-cert string         Certificate to use for HTTPS calls
      --verbose                     Whether to verbose log

Use "ark cache [command] --help" for more information about a command.
```


SDK Usage
=========
As well as using the CLI, one can also develop under the ark sdk using its API / class driven design

The same idea as the CLI applies here as well

Let's say we want to generate a short lived password from the code

To do so, we can use the following script:
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

More examples can be found in the examples folder

## License

This project is licensed under Apache License 2.0 - see [`LICENSE`](LICENSE.txt) for more details

Copyright (c) 2025 CyberArk Software Ltd. All rights reserved.
