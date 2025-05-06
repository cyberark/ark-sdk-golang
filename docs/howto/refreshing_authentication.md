---
title: Refresh authentication
description: Refreshing Authentication
---

# Refresh authentication

When you want to continue working with an existing authenticator, you can refresh the authentications. You can refresh authentications for the following:

- The login command
- The exec command

## Login command

To try to authenticate with an existing authenticator, use the `--refresh-auth` CLI flag:

```bash  linenums="0"
ark login --refresh-auth
```

The `--refresh-auth` flag indicates that the user's profile authenticator should be refreshed and used for authentication. The user is only prompted for additional authentication values when the refresh fails.

## Exec command

To try to run any command with an existing authenticator, use the `--refresh-auth` CLI flag:
```bash  linenums="0"
ark exec --refresh-auth sia sso short-lived-password
```

The `--refresh-auth` flag indicates that the user's profile authenticator should be refreshed and used before executing the command. When the refresh fails, an error is returned and you must log in again.
