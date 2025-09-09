---
title: UAP database policy CLI workflow
description: Creating a UAP DB Policy using Ark CLI
---

# UAP database policy CLI workflow
Here is an example workflow for adding a UAP DB policy alongside all needed assets via the CLI:

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
1. Add SIA DB User Secret
    ```shell
    ark exec sia secrets db add-secret --secret-name mysecret --secret-type username_password --username user --password mypass
    ```
1. Add SIA Database
    ```shell
    ark exec sia workspaces db add-database \
      --name mydomain.com \
      --provider-engine postgres-sh \
      --read-write-endpoint myendpoint.mydomain.com \
      --secret-id <SECRET_ID_FROM_PREVIOUS_STEP>
    ```
1. Create UAP DB Policy using a defined json file
    ```json
    {
      "metadata": {
        "name": "Cool Policy",
        "description": "Cool Policy Description",
        "status": { "status": "ACTIVE" },
        "timeFrame": { "fromTime": null, "toTime": null },
        "policyEntitlement": {
          "targetCategory": "DB",
          "locationType": "FQDN_IP",
          "policyType": "RECURRING"
        },
        "policyTags": ["cool_tag", "cool_tag2"],
        "timeZone": "Asia/Jerusalem"
      },
      "principals": [
        {
          "id": "principal_id",
          "name": "tester@cyberark.cloud",
          "sourceDirectoryName": "CyberArk Cloud Directory",
          "sourceDirectoryId": "source_directory_id",
          "type": "USER"
        }
      ],
      "conditions": {
        "accessWindow": {
          "daysOfTheWeek": [0, 1, 2, 3, 4, 5, 6],
          "fromHour": "05:00",
          "toHour": "23:59"
        },
        "maxSessionDuration": 2,
        "idleTime": 1
      },
      "targets": {
        "FQDN_IP": {
          "instances": [
            {
              "instanceName": "Mongo-atlas_ephemeral_user",
              "instanceType": "Mongo",
              "instanceId": "1234",
              "authenticationMethod": "MONGO_AUTH",
              "mongoAuthProfile": {
                "globalBuiltinRoles": ["readWriteAnyDatabase"],
                "databaseBuiltinRoles": {
                  "mydb1": ["userAdmin"],
                  "mydb2": ["dbAdmin"]
                },
                "databaseCustomRoles": {
                  "mydb1": ["myCoolRole"]
                }
              }
            }
          ]
        }
      }
    }
    ```

    ```shell
    ark exec --request-file /path/to/policy-request.json uap db add-policy
    ```
