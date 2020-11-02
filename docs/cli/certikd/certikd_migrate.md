## certikd migrate

Migrate genesis

### Synopsis

Migrate the source genesis into the target version and print to STDOUT.
Example:
$ certikd migrate /path/to/genesis.json --chain-id=shentu-incentivized-2 --genesis-time=2020-07-24T17:00:00Z


```
certikd migrate [genesis-file] [flags]
```

### Options

```
      --chain-id string       override chain_id with this flag
      --genesis-time string   override genesis_time with this flag
  -h, --help                  help for migrate
```

### Options inherited from parent commands

```
      --home string        directory for config and data (default "~/.certikd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors
```

### SEE ALSO

* [certikd](certikd.md)	 - CertiK App Daemon (server)


