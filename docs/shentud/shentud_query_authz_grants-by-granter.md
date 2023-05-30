## shentud query authz grants-by-granter

query authorization grants granted by granter

### Synopsis

Query authorization grants granted by granter.
Examples:
$ shentud q authz grants-by-granter cosmos1skj..

```
shentud query authz grants-by-granter [granter-addr] [flags]
```

### Options

```
      --count-total       count total number of records in granter-grants to query for
      --height int        Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help              help for grants-by-granter
      --limit uint        pagination limit of granter-grants to query for (default 100)
      --node string       <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint       pagination offset of granter-grants to query for
  -o, --output string     Output format (text|json) (default "text")
      --page uint         pagination page of granter-grants to query for. This sets offset to a multiple of limit (default 1)
      --page-key string   pagination page-key of granter-grants to query for
      --reverse           results are sorted in descending order
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

* [shentud query authz](shentud_query_authz.md)	 - Querying commands for the authz module


