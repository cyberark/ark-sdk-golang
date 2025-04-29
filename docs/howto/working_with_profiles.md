---
title: Work with profiles
description: Working With Profiles
---

# Work with profiles
Profiles define authentication methods for users. They are used with the CLI and, to a lesser extent, the SDK. Different profiles can be created and configured via the Ark `configure` command.

You can specify which profile a command uses with the `--profile-name` flag or setting the `ARK_PROFILE` environment variable.

Profiles are stored as JSON files in the `$HOME/.ark_profiles` folder.

!!! note

    When There are multiple profiles configured, if profile is not specified explicitly in the command via<br>
    **--profile-name** or **ARK_PROFILE** environment variable not set, the default profile called **"ark"** will be used.


Here is an example profile file:

``` json
{
    "profile_name": "ark",
    "profile_description": "Default Ark Profile",
    "auth_profiles": {
        "isp": {
            "username": "tina@cyberark.cloud.1234567",
            "auth_method": "identity",
            "auth_method_settings": {
                "identity_mfa_method": "email",
                "identity_mfa_interactive": true,
                "identity_application": null,
                "identity_url": null
            }
        }
    }
}
```

As well as using the CLI to manage profiles, you can create, modify, and delete profiles directly in the `$HOME/.ark_profiles` folder.
