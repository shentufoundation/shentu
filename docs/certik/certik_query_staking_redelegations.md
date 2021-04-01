## certik query staking redelegations

Query all redelegations records for one delegator

### Synopsis

Query all redelegation records for an individual delegator.

Example:
$ <appd> query staking redelegation certik1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p

```
certik query staking redelegations [delegator-addr] [flags]
```

### Options

```
      --count-total       count total number of records in delegator redelegations to query for
      --height int        Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help              help for redelegations
      --limit uint        pagination limit of delegator redelegations to query for (default 100)
      --node string       <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint       pagination offset of delegator redelegations to query for
  -o, --output string     Output format (text|json) (default "text")
      --page uint         pagination page of delegator redelegations to query for. This sets offset to a multiple of limit (default 1)
      --page-key string   pagination page-key of delegator redelegations to query for
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data (default "~/.certik")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [certik query staking](certik_query_staking.md)	 - Querying commands for the staking module


