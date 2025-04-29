---
title: Install SIA connectors
description: Install SIA connectors
---

# Install SIA connectors
Here is an example workflow for installing a connector on a linux / windows box:

1. Install Ark SDK with your artifactory credentials:
   ```shell linenums="0"
   go install github.com/cyberark/ark-sdk-golang
   ```
   Make sure that the PATH environment variable points to the go binary path, for example
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
1. Create a network and connector pool:
    ```shell linenums="0"
    ark exec cmgr add-network --name mynetwork
    ark exec cmgr add-pool --name mypool --assigned-network-ids mynetwork_id
    ```
   1. Install a connector:
       * Windows:
           ```shell linenums="0"
           ark exec sia access install-connector --connector-pool-id 89b4f0ff-9b06-445a-9ca8-4ca9a4d72e8c --username myuser --password mypassword --target-machine 1.1.1.1 --connector-os windows --connector-type ON-PREMISE
           ```
       * Linux:
           ```shell linenums="0"
           ark exec sia access install-connector --connector-pool-id 89b4f0ff-9b06-445a-9ca8-4ca9a4d72e8c --username myuser --private-key-path /path/to/private_key.pem --target-machine 1.1.1.1 --connector-os linux --connector-type ON-PREMISE
           ```
