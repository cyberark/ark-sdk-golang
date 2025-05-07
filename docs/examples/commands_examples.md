---
title: Commands examples
description: Commands Examples
---

# Commands examples

This page lists some useful CLI examples.

!!! note

    You can disable certificate validation for login to an authenticator using the `--disable-certificate-verification` flag. **This option is not recommended.**

    **Useful environment variables**

    - `ARK_PROFILE`: Sets the profile to be used across the CLI
    - `ARK_DISABLE_CERTIFICATE_VERIFICATION`: Disables certificate verification for REST APIs

## Configure command example

The `configure` command works in interactive or silent mode. When using silent mode, the required parameters need to specified. Here's an example of configuring ISP in silent mode:

```bash linenums="0"
ark configure --profile-name="PROD" --work-with-isp --isp-username="tina@cyberark.cloud.12345" --silent --allow-output
```

## Login commands example

The login command can work in interactive or silent mode. Here's an example of a silent login with the profile configured in the example above:
```bash linenums="0"
ark login -s --isp-secret=CoolPasswordß --profile-name PROD
```

## Exec command examples

Use the `--help` flag to view all `exec` options.

### Generate a short-lived SSO password for a database connection
```shell linenums="0"
ark exec sia sso short-lived-password
```

### Generate a short-lived SSO Oracle wallet for an Oracle database connection
```shell linenums="0"
ark exec sia sso short-lived-oracle-wallet --folder ~/wallet
```

### Generate a kubectl config file
```shell linenums="0"
ark exec sia k8s generate-kubeconfig 
```

### Generate a kubectl config file and save it in the specified path
```shell linenums="0"
ark exec sia k8s generate-kubeconfig --folder=/Users/My.User/.kube
```

### Add SIA VM target set
```shell
ark exec sia workspaces target-sets add-target-set --name mydomain.com --type Domain
```

### Add SIA VM secret
```shell
ark exec sia secrets vm add-secret --secret-type ProvisionerUser --provisioner-username=myuser --provisioner-password=mypassword
```

### List CMGR connector pools
```shell
ark exec cmgr list-pools
```

### Add CMGR network
```shell
ark exec cmgr add-network --name mynetwork
```

### Add CMGR connector pool
```shell
ark exec cmgr add-pool --name mypool --assigned-network-ids mynetwork_id
```

### Get connector installation script
```shell
ark exec sia access connector-setup-script --connector-type ON-PREMISE --connector-os windows --connector-pool-id 588741d5-e059-479d-b4c4-3d821a87f012
```

### Install a connector on windows remotely
```shell
ark exec sia access install-connector --connector-pool-id 89b4f0ff-9b06-445a-9ca8-4ca9a4d72e8c --username myuser --password mypassword --target-machine 1.2.3.4 --connector-os windows --connector-type ON-PREMISE
```

### Install a connector on linux with private key remotely
```shell
ark exec sia access install-connector --connector-pool-id 89b4f0ff-9b06-445a-9ca8-4ca9a4d72e8c --username myuser --private-key-path /path/to/private_key.pem --target-machine 1.2.3.4 --connector-os linux --connector-type ON-PREMISE
```

### Uninstall a connector remotely
```shell
ark exec sia access uninstall-connector --connector-id CMSConnector_588741d5-e059-479d-b4c4-3d821a87f012_1a8b3734-8e1d-43a3-bb99-8a587609e653 --username myuser --password mypassword --target-machine 1.2.3.4 --connector-os windows
```

### Create a pCloud Safe
```shell
ark exec pcloud safes add-safe --safe-name=safe
```

### Create a pCloud account
```shell
ark exec pcloud accounts add-account --name account --safe-name safe --platform-id='UnixSSH' --username root --address 1.2.3.4 --secret-type=password --secret mypass
```

### Retrieve a pCloud account credentials
```shell
ark exec pcloud accounts get-account-credentials --account-id 11_1
```

### Create an Identity user
```shell
ark exec identity users create-user --roles "DpaAdmin" --username "myuser"
```

### Create an Identity role
```shell
ark exec identity roles create-role --role-name myrole
```

### List all directories identities
```shell
ark exec identity directories list-directories-entities
```

### Add SIA database secret

```shell linenums="0"
ark exec sia secrets db add-secret --secret-name mysecret --secret-type username_password --username user --password mypass
```

### Delete SIA database secret

```shell linenums="0"
ark exec sia secrets db delete-secret --secret-name mysecret
```

### Add SIA database
```shell linenums="0"
ark exec sia workspaces db add-database --name mydatabase --provider-engine aurora-mysql --read-write-endpoint myrds.com
```

### Delete SIA database
```shell linenums="0"
ark exec sia workspaces db delete-database --id databaseid
```
