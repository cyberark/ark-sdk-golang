---
title: Profiles
description: Profiles Command
---

# Profiles

Use the `profiles` command to manage multiple users and tenants. You can create, copy, modify, and delete profiles for different users and tenant, and to list all existing profiles.

## Running
```shell linenums="0"
ark profiles
```

## Usage
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
