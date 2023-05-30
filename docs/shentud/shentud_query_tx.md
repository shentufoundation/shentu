## shentud query tx

Query for a transaction by hash, "<addr>/<seq>" combination or comma-separated signatures in a committed block

### Synopsis

Example:
$ shentud query tx <hash>
$ shentud query tx --type=acc_seq <addr>/<sequence>
$ shentud query tx --type=signature <sig1_base64>,<sig2_base64...>

```
shentud query tx --type=[hash|acc_seq|signature] [hash|acc_seq|signature] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for tx
      --node string     <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
      --type string     The type to be used when querying tx, can be one of "hash", "acc_seq", "signature" (default "hash")
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data (default "~/.shentud")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [shentud query](shentud_query.md)	 - Querying subcommands


