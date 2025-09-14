---
title: UAP VM policy CLI workflow
description: Creating a UAP VM Policy using Ark CLI
---

# UAP VM policy CLI workflow
Here is an example workflow for adding a UAP VM policy via the CLI:

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
1. Create UAP VM Policy using a defined json file
    ```json
    {
      "metadata": {
        "name": "Cool Policy",
        "description": "Cool Policy Description",
        "status": {
          "status": "Active",
          "statusDescription": "Example status description"
        },
        "timeFrame": {
          "fromTime": null,
          "toTime": null
        },
        "policyEntitlement": {
          "targetCategory": "VM",
          "locationType": "FQDN/IP",
          "policyType": "Recurring"
        },
        "createdBy": {
          "user": "cool_user",
          "time": "2025-02-08T22:46:06"
        },
        "updatedOn": {
          "user": "cool_user",
          "time": "2025-02-08T22:46:06"
        },
        "policyTags": [
          "cool_tag",
          "cool_tag2"
        ],
        "timeZone": "Asia/Jerusalem"
      },
      "principals": [
        {
          "id": "principal_id",
          "name": "tester@cyberark.cloud",
          "type": "User",
          "sourceDirectoryName": "CyberArk Cloud Directory",
          "sourceDirectoryId": "source_directory_id"
        }
      ],
      "delegationClassification": "Unrestricted",
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
          "fromHour": "05:00",
          "toHour": "23:59"
        },
        "maxSessionDuration": 2,
        "idleTime": 1
      },
      "targets": {
        "fqdnipResource": {
          "fqdnRules": [
            {
              "operator": "EXACTLY",
              "computernamePattern": "myvm.mydomain.com",
              "domain": "domain.com"
            }
          ],
          "ipRules": [
            {
              "operator": "EXACTLY",
              "ipAddresses": [
                "192.168.12.34"
              ],
              "logicalName": "CoolLogicalName"
            }
          ]
        }
      },
      "behavior": {
        "sshProfile": {
          "username": "ssh_user"
        },
        "rdpProfile": {
          "domainEphemeralUser": {
            "assignGroups": [
              "rdp_users"
            ],
            "enableEphemeralUserReconnect": false,
            "assignDomainGroups": [
              "domain_rdp_users"
            ]
          }
        }
      }
    }
    ```

    ```shell
    ark exec --request-file /path/to/policy-request.json uap vm add-policy
    ```
