---
title: Onboard pCloud Account
description: Onboard pCloud Account
---

# Onboard pCloud Account
Here is an example workflow for onboarding a pCloud safe and creating a Safe:

1. Install Ark SDK:
   ```shell linenums="0"
   go install github.com/cyberark/ark-sdk-golang
   ```
   Make sure that the PATH environment variable points to the go binary. For example:
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
1. Create a new safe:
    ```shell linenums="0"
    ark exec pcloud safes add-safe --safe-name=safe
    ```
1. Create a new account in the Safe:
    ```shell linenums="0"
    ark exec pcloud accounts add-account --name account --safe-name safe --platform-id='UnixSSH' --username root --address 1.2.3.4 --secret-type=password --secret mypass
    ```
