---
title: Login
description: Login Command
---

# Login

The `login` command is used to authenticate to Ark using the configured profile. When you run the command, you are prompted for the required login information (such as a password and MFA verifications).

After you have logged in, the returned access tokens are stored in a secure location on your machine. After the tokens expire, a token refresh maybe attempted (see [Refresh token](../howto/refreshing_authentication.md)) or a new login is required.

## Run
```shell linenums="0"
ark login
```

## Usage
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
