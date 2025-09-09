---
title: End-user database psql
description: End-user Database psql
---

# End-user database Workflow
Here is an example workflow for connecting to a psql DB via ark CLI, which assumes a DB was already onboarded:

1. Install Ark SDK:
   ```shell linenums="0"
   go install github.com/cyberark/ark-sdk-golang/cmd/ark
   ```
   Make sure that the PATH environment variable points to the Go binary. For example:
   ```shell linenums="0"
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
1. Create a profile:
    * Interactively:
        ```shell linenums="0"
        ark configure
        ```
    * Silently:
        ```shell linenums="0"
        ark configure --silent --work-with-isp --isp-username myuser
        ```
1. Log in to Ark:
    ```shell linenums="0"
    ark login --silent --isp-secret <my-ark-secret>
    ```
1. Connect to postgres using the CLI with an MFA caching token behind the scenes:
    ```shell linenums="0"
    ark exec sia db psql --target-address myaddress.com --target-user myuser
    ```
