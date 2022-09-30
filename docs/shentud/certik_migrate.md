## certik migrate

Migrate 1.X.X genesis to 2.0.X

### Synopsis

Migrate the source genesis into the target version and print to STDOUT.

Example:
$ <appd> migrate /path/to/genesis.json --chain-id=cosmoshub-4 --genesis-time=2019-04-22T17:00:00Z --initial-height=5000


```
certik migrate [genesis-file] [flags]
```

### Options

```
      --chain-id string                override chain_id with this flag
      --genesis-time string            override genesis_time with this flag
  -h, --help                           help for migrate
      --initial-height int             Set the starting height for the chain
      --no-prop-29                     Do not implement fund recovery from prop29
      --replacement-cons-keys string   Proviide a JSON file to replace the consensus keys of validators
```

### Options inherited from parent commands

```
      --home string         directory for config and data (default "~/.certik")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [certik](certik.md)	 - Stargate CosmosHub App


