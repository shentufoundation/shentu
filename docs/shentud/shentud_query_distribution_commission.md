## shentud query distribution commission

Query distribution validator commission

### Synopsis

Query validator commission rewards from delegators to that validator.

Example:
$ shentud query distribution commission shentuvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldtshsu6

```
shentud query distribution commission [validator] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for commission
      --node string     <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
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

* [shentud query distribution](shentud_query_distribution.md)	 - Querying commands for the distribution module


