## shentud init

Initialize private validator, p2p, genesis, and application configuration files

### Synopsis

Initialize validators's and node's configuration files.

```
shentud init [moniker] [flags]
```

### Options

```
      --chain-id string   genesis file chain-id, if left blank will be randomly created
  -h, --help              help for init
  -o, --overwrite         overwrite the genesis.json file
      --recover           provide seed phrase to recover existing key instead of creating
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


