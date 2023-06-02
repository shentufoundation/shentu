## shentud query ibc-transfer denom-trace

Query the denom trace info from a given trace hash or ibc denom

### Synopsis

Query the denom trace info from a given trace hash or ibc denom

```
shentud query ibc-transfer denom-trace [hash/denom] [flags]
```

### Examples

```
shentud query ibc-transfer denom-trace 27A6394C3F9FF9C9DCF5DFFADF9BB5FE9A37C7E92B006199894CF1824DF9AC7C
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for denom-trace
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

* [shentud query ibc-transfer](shentud_query_ibc-transfer.md)	 - IBC fungible token transfer query subcommands


