---
title: End-user Kubernetes workflow
description: End-user Kubernetes Workflow
---

# End-user Kubernetes workflow

To securely access a Kubernetes cluster, do the following:

1. Install Ark SDK with your artifactory credentials:
   ```shell linenums="0"
   go install github.com/cyberark/ark-sdk-golang
   ```
   Make sure that the PATH environment variable points to the go binary path, for example
   ```shell linenums="0"
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
1. Configure a profile:
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
1. To generate a kubectl kubeconfig file, which defines the user's permissions and accessible clusters, do **one** of the following:
    * Run this command:
        ```shell linenums="0"
        ark exec sia k8s generate-kubeconfig
        ```
    * Use the `-f` flag to generate the config file in the specified path (this option **overrides** existing files with the same name):
        ```shell
        ark exec sia k8s generate-kubeconfig -f ~/.kube
        ```

## Refresh SSO certificate workflow

When you refresh the certificate, you can keep using its associated kubeconfig file and only need to refresh the MFA authentication data. To refresh the certificate, run **one** of the following:

* To generate two files (certificate and private key files), where the required `-f` flag defines the generated files' location:
    ```shell
    ark exec sia sso short-lived-client-certificate -of file -f ~/home
    ```
* To print the certificate and private key to the console as plaintext:
    ```shell
    ark exec sia sso short-lived-client-certificate -of raw
    ```
* To print the certificate and private key to the console as base64-encoded strings:
    ```shell
    ark ark exec sia sso short-lived-client-certificate -of base64
    ```
