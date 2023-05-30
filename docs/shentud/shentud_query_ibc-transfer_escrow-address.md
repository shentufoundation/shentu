## shentud query ibc-transfer escrow-address

Get the escrow address for a channel

### Synopsis

Get the escrow address for a channel

```
shentud query ibc-transfer escrow-address [flags]
```

### Examples

```
shentud query ibc-transfer escrow-address [port] [channel-id]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for escrow-address
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


