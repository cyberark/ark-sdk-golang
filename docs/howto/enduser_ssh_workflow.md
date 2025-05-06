---
title: End-user ssh workflow
description: End-user ssh Workflow
---

# End-user SSH workflow
Here is an example workflow for connecting to a linux box using SSH:

1. Install Ark SDK:
   ```shell linenums="0"
   go install github.com/cyberark/ark-sdk-golang
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
1. Get a short-lived SSH private key for a linux box from the SIA service:
    ```shell linenums="0"
    ark exec sia sso short-lived-ssh-key
    ```
1. Log in directly to the linux box:
    ```shell linenums="0"
    ssh -i ~/.ssh/sia_ssh_key.pem myuser@suffix@targetuser@targetaddress@sia_proxy
    ```
