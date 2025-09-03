---
title: Exec
description: Exec Command
---

# Exec

Use the `exec` command to run commands on available services (the available services depend on the authorized user's account).

## SIA services

The following SIA commands are supported:

- `ark exec sia`: Root command for the SIA service (aliases: dpa)
    - `sso` - SSO end-user operations
    - `k8s` - Kubernetes service
    - `db` - DB service
    - `workspaces` - Workspaces service
      - `target-sets` - Target sets operations
      - `db` - Database operations
    - `secrets` - Secrets service
      - `vm` - VM operations
      - `db` - Database operations
    - `access` - Access service
    - `ssh-ca` - SSH CA key service
- `ark exec cmgr`: Root command for the CMGR service (aliases: connectormanager,cm)
- `ark exec pcloud`: Root command for PCloud service (aliases: privilegecloud,pc)
    - `accounts` - Accounts management
    - `safes` - Safes management
- `ark exec identity`: Root command for the Identity service (aliases: idaptive,id)
    - `directories` - Directories management
    - `users` - Users management
    - `roles` - Roles management
-  `ark exec sechub`: Root command for the Secrets Hub Service (aliases: secretshub,sh)
    - `configuration` - Configuration management
    - `service-info` - Service Info management
    - `secrets` - Secrets management
    - `scans` - Scans management
    - `secret-stores` - Secret Stores management
    - `sync-policies` - Sync Policies management
- `ark exec sm`: Root command for the SM service (aliases: sessionmonitoring)
- `ark exec uap`: Root command for the UAP service (aliases: useraccesspolicies)
    - `sca` - SCA management
    - `db` - SIA DB management
    - `vm` - SIA VM management

All commands have their own subcommands and respective arguments and aliases.

## Running
```shell linenums="0"
ark exec
```

## Usage
```shell
Exec an action

Usage:
  ark exec [command]

Available Commands:
  cmgr        (aliases: connectormanager, cm)
  identity    (aliases: idaptive, id)
  pcloud      (aliases: privilegecloud, pc)
  sechub      (aliases: secretshub, sh)
  sia         (aliases: dpa)
  sm          (aliases: sessionmonitoring)
  uap         (aliases: useraccesspolicies)

Flags:
      --allow-output                Allow stdout / stderr even when silent and not interactive
      --disable-cert-verification   Disables certificate verification on HTTPS calls, unsafe!
  -h, --help                        help for exec
      --log-level string            Log level to use while verbose (default "INFO")
      --logger-style string         Which verbose logger style to use (default "default")
      --output-path string          Output file to write data to
      --profile-name string         Profile name to load (default "ark")
      --raw                         Whether to raw output
      --refresh-auth                If a cache exists, will also try to refresh it
      --request-file string         Request file containing the parameters for the exec action
      --retry-count int             Retry count for execution (default 1)
      --silent                      Silent execution, no interactiveness
      --trusted-cert string         Certificate to use for HTTPS calls
      --verbose                     Whether to verbose log

Use "ark exec [command] --help" for more information about a command.
```
