---
title: Work with Ark cache
description: Working With Ark Cache
---

# Work with Ark cache

Both the CLI and SDK cache login information in the local machine's keystore or, when a keystore does not exist, in an encrypted folder (located in `$HOME/.ark_cache`). The cached information is used to run commands until the authentication tokens expire or are otherwise invalided.

You can set the cache folder with the `ARK_KEYRING_FOLDER` env variable. To force Ark SDK to work only with the filesystem cache, use the `ARK_BASIC_KEYRING` environment variable

If you want to ignore the cache when logging in, use the `-f` flag:
``` bash  linenums="0"
ark login --force
```

To clear the cache, run `ark cache clear` or, when using an encrypted folder, remove the files from the `$HOME/.ark_cache` folder.
