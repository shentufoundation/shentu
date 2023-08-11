## shentud query distribution rewards

Query all distribution delegator rewards or rewards from a particular validator

### Synopsis

Query all rewards earned by a delegator, optionally restrict to rewards from a single validator.

Example:
$ shentud query distribution rewards shentu1gghjut3ccd8ay0zduzj64hwre2fxs9ldkq89hu
$ shentud query distribution rewards shentu1gghjut3ccd8ay0zduzj64hwre2fxs9ldkq89hu shentuvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldtshsu6

```
shentud query distribution rewards [delegator-addr] [validator-addr] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for rewards
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


