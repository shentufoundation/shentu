## shentud replace-validators

Replace top N validators in a given genesis with a set json

### Synopsis

Migrate the source genesis into the target version and print to STDOUT.
Example:
$ shentud migrate /path/to/genesis.json --chain-id=cosmoshub-4 --genesis-time=2019-04-22T17:00:00Z --initial-height=5000


```
shentud replace-validators [genesis-file] [replacement-cons-keys] [flags]
```

### Options

```
  -h, --help                           help for replace-validators
      --replacement-cons-keys string   Proviide a JSON file to replace the consensus keys of validators
```

### Options inherited from parent commands

```
      --home string         directory for config and data (default "~/.shentud")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [shentud](shentud.md)	 - Stargate Shentu Chain App


