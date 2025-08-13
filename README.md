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
  - [x] SIA SSH CA Key Service
  - [x] Connector Manager Service
  - [x] PCloud Accounts Service
  - [x] PCloud Safes Service
  - [x] Identity Directories Service
  - [x] Identity Roles Service
  - [x] Identity Users Service
  - [x] Secrets Hub Secret Stores Service
  - [x] Secrets Hub Secrets Service
  - [x] Secrets Hub Sync Policies Service
  - [x] Secrets Hub Scans Service
  - [x] Secrets Hub Service Info Service
  - [x] Secrets Hub Configuration Service
  - [x] Secrets Hub Filters Service
  - [x] Session Monitoring Service
  - [x] Unified Access Policies Service
    - [x] SCA - Secure Cloud Access
    - [x] DB - Databases
    - [x] VM - Virtual Machines
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
- <b>uap</b> - Unified Access Policies Services
  - <b>sca</b> - secure cloud access policies management
  - <b>db</b> - databases access policies management
  - <b>vm</b> - virtual machines access policies management

Any command has its own subcommands, with respective arguments

For example, generating a short lived password for DB
```shell
ark exec sia sso short-lived-password
```

Or a short lived password for RDP
```shell
ark exec sia sso short-lived-password --service DPA-RDP
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

Get Secrets Hub Configuration
```shell
ark exec sechub configuration get-configuration
```
Set Secrets Hub Configuration
```shell 
ark exec sechub configuration set-configuration --sync-settings 360
```

Get Secrets Hub Filters
```shell
ark exec sechub filters get-filters --store-id store-e488dd22-a59c-418c-bbe3-3f061dd9b667
```
Add Secrets Hub Filter
```shell
ark exec sechub filters add-filter --type "PAM_SAFE" --store-id store-e488dd22-a59c-418c-bbe3-3f061dd9b667 --data-safe-name "example-safe"
```
Delete Secrets Hub Filter
```shell
ark exec sechub filters delete-filter --filter-id filter-7f3d187d-7439-407f-b968-ec27650be692 --store-id store-e488dd22-a59c-418c-bbe3-3f061dd9b667
```

Get Secrets Hub Scans
```shell
ark exec sechub scans get-scans 
```
Trigger Secrets Hub Scan
```shell
ark exec sechub scans trigger-scan --id default --secret-stores-ids store-e488dd22-a59c-418c-bbe3-3f061dd9b667 type secret-store
```

Create Secrets Hub Secret Store
```shell
ark exec sechub secret-stores create-secret-store --type AWS_ASM --description sdk-testing --name "SDK Testing" --state ENABLED --data-aws-account-alias ALIAS-NAME-EXAMPLE --data-aws-region-id us-east-1 --data-aws-account-id 123456789123 --data-aws-rolename Secrets-Hub-IAM-Role-Name-Created-For-Secrets-Hub
```
Retrieve Secrets Hub Secret Store
```shell
ark exec sechub secret-stores get-secret-store --secret-store-id store-e488dd22-a59c-418c-bbe3-3f061dd9b667
```
Update Secrets Hub Secret Store
```shell
ark exec sechub secret-stores update-secret-store --secret-store-id store-7f3d187d-7439-407f-b968-ec27650be692 --name "New Name" --description "Updated Description" --data-aws-account-alias "Test2"
```
Delete Secrets Hub Secret Store
```shell
ark exec sechub secret-stores delete-secret-store --secret-store-id store-fd11bc7c-22d0-4d9b-ac1b-f8458161935f
```

Get Secrets Hub Secrets
```shell
ark exec sechub secrets get-secrets
```
Get Secrets Hub Secrets using a filter
```shell
ark exec sechub secrets get-secrets-by --limit 5 --projection EXTEND --filter "name CONTAINS EXAMPLE"
```

Get Secrets Hub Service Information
```shell
ark exec sechub service-info get-service-info
```

Get Secrets Hub Sync Policies
```shell
ark exec sechub sync-policies get-sync-policies
```
Get Secrets Hub Sync Policy
```shell
ark exec sechub sync-policies get-sync-policy --policy-id policy-7f3d187d-7439-407f-b968-ec27650be692 --projection EXTEND
```
Create Secrets Hub Sync Policy
```shell
ark exec sechub sync-policies create-sync-policy --name "New Sync Policy" --description "New Sync Policy Description" --filter-type PAM_SAFE --filter-data-safe-name EXAMPLE-SAFE-NAME --source-id store-e488dd22-a59c-418c-bbe3-3f061dd12367 --target-id store-e488dd22-a59c-418c-bbe3-3f061dd9b667
```
Delete Secrets Hub Sync Policy
```shell
ark exec sechub sync-policies delete-sync-policy --policy-id policy-7f3d187d-7439-407f-b968-ec27650be692
```

List Sessions
```shell
ark exec sm list-sessions
```
Count Sessions
```shell
ark exec sm count-sessions
```
List Sessions By Filter
```shell
ark exec sm list-sessions-by --search "duration LE 01:00:00"
```
Count Sessions By Filter
```shell
ark exec sm count-sessions-by --search "command STARTSWITH ls"
```
Get Session
```shell
ark exec sm get-session --session-id my-id
```
List Session Activities
```shell
ark exec sm list-session-activities --session-id my-id
```
Count Session Activities
```shell
ark exec sm count-session-activities --session-id my-id
```
List Session Activities By Filter
```shell
ark exec sm list-session-activities-by --session-id my-id --command-contains "ls"
```
Count Session Activities By Filter
```shell
ark exec sm count-sessions-by --session-id my-id --command-contains "chmod"
```
Get Sessions Statistics
```shell
ark exec sm get-sessions-stats
```

List all UAP policies
```shell
ark exec uap list-policies
```

Delete UAP DB Policy
```shell
ark exec uap db delete-policy --policy-id my-policy-id
```

List DB Policies from UAP
```shell
ark exec uap db list-policies
```

Get DB Policy from UAP
```shell
ark exec uap db policy --policy-id my-policy-id
```

Add UAP DB Policy
```shell
ark exec uap db add-policy --request-file /path/to/policy-request.json
```

List UAP SCA Policies
```shell
ark exec uap sca list-policies
```

Get UAP SCA Policy
```shell
ark exec uap sca policy --policy-id my-policy-id
```

Add UAP SCA Policy
```shell
ark exec uap sca add-policy --request-file /path/to/policy-request.json
```

Delete UAP SCA Policy
```shell
ark exec uap sca delete-policy --policy-id my-policy-id
```

List VM Policies from UAP
```shell
ark exec uap vm list-policies
```

Get VM Policy from UAP
```shell
ark exec uap vm policy --policy-id my-policy-id
```

Delete VM Policy from UAP
```shell
ark exec uap vm delete-policy --policy-id my-policy-id
```

Connect to MySQL ZSP with the mysql cli via Ark CLI
```shell
ark exec sia db mysql --target-address myaddress.com
```

Connect to PostgreSQL Vaulted with the psql cli via Ark CLI
```shell
ark exec sia db psql --target-address myaddress.com --target-user myuser
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
