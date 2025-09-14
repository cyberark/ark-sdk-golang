---
title: UAP SCA policy CLI workflow
description: Creating a UAP SCA Policy using Ark CLI
---

# UAP SCA policy CLI workflow
Here is an example workflow for adding a UAP SCA policy via the CLI:

1. Install Ark SDK:
   ```shell linenums="0"
   go install github.com/cyberark/ark-sdk-golang/cmd/ark@latest
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
1. Create UAP SCA Policy using a defined json file
    ```json
    {
      "metadata": {
        "name": "Cool Cloud Policy",
        "description": "Cool Cloud Policy Description",
        "policyTags": [
          "cool_tag",
          "cool_tag2"
        ],
        "policyEntitlement": {
          "targetCategory": "Cloud console",
          "locationType": "AWS",
          "policyType": "Recurring"
        },
        "timeFrame": {
          "fromTime": null,
          "toTime": null
        },
        "status": {
          "status": "Validating",
          "statusCode": null,
          "statusDescription": "Example status description",
          "link": null
        }
      },
      "principals": [
        {
          "id": "c2c7bcc6-9560-44e0-8dff-5be221cd37ee",
          "name": "user@cyberark.cloud.12345",
          "type": "User",
          "sourceDirectoryName": "CyberArk Cloud Directory",
          "sourceDirectoryId": "09B9A9B0-6CE8-465F-AB03-65766D33B05E"
        }
      ],
      "conditions": {
        "accessWindow": {
          "daysOfTheWeek": [
            0,
            1,
            2,
            3,
            4,
            5,
            6
          ],
          "fromHour": "05:00:00",
          "toHour": "23:59:00"
        },
        "maxSessionDuration": 2
      },
      "delegationClassification": "Unrestricted",
      "targets": {
        "awsAccountTargets": [
          {
            "roleId": "arn:aws:iam::123456789012:role/RoleName",
            "workspaceId": "123456789012",
            "roleName": "RoleName",
            "workspaceName": "WorkspaceName"
          }
        ]
      }
    }
    ```

    ```shell
    ark exec --request-file /path/to/policy-request.json uap sca add-policy
    ```
