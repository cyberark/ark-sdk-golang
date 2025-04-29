---
title: Cache
description: Cache Command
---

# Cache

Use the `cache` command to manage the Ark data cached on your machine. Currently, you can only clear the filesystem cache (not data cached in the OS's keystore).

## Running
```shell linenums="0"
ark cache
```


## Usage
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
