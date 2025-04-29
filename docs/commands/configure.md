---
title: Configure
description: Configure Command
---

# Configure command

The `configure` command is used to create a profile. Profiles define user and authentication information, such as which authentication methods to use, the method settings, and other information like MFA.

Profiles are saved in the `~/.ark_profiles` folder.

## Run

```shell linenums="0"
ark configure
```

When you run the command without arguments, you are prompted for the required information (alternatively, add the `--silent` flag with the required arguments).

## Usage

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
