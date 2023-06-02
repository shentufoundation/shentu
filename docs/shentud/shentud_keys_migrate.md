## shentud keys migrate

Migrate keys from the legacy (db-based) Keybase

### Synopsis

Migrate key information from the legacy (db-based) Keybase to the new keyring-based Keyring.
The legacy Keybase used to persist keys in a LevelDB database stored in a 'keys' sub-directory of
the old client application's home directory, e.g. $HOME/.gaiacli/keys/.
For each key material entry, the command will prompt if the key should be skipped or not. If the key
is not to be skipped, the passphrase must be entered. The key will only be migrated if the passphrase
is correct. Otherwise, the command will exit and migration must be repeated.

It is recommended to run in 'dry-run' mode first to verify all key migration material.


```
shentud keys migrate <old_home_dir> [flags]
```

### Options

```
      --dry-run   Run migration without actually persisting any changes to the new Keybase
  -h, --help      help for migrate
```

### Options inherited from parent commands

```
      --home string              The application home directory (default "~/.shentud")
      --keyring-backend string   Select keyring's backend (os|file|test) (default "os")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) (default "plain")
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --output string            Output format (text|json) (default "text")
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [shentud keys](shentud_keys.md)	 - Manage your application's keys


