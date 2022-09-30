## certik keys list

List all keys

### Synopsis

Return a list of all public keys stored by this key manager
along with their associated name and address.

```
certik keys list [flags]
```

### Options

```
  -h, --help         help for list
  -n, --list-names   List names only
```

### Options inherited from parent commands

```
      --home string              The application home directory (default "~/.certik")
      --keyring-backend string   Select keyring's backend (os|file|test) (default "os")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) (default "plain")
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --output string            Output format (text|json) (default "text")
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [certik keys](certik_keys.md)	 - Manage your application's keys


