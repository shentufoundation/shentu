## shentud query ibc channel next-sequence-receive

Query a next receive sequence

### Synopsis

Query the next receive sequence for a given channel

```
shentud query ibc channel next-sequence-receive [port-id] [channel-id] [flags]
```

### Examples

```
shentud query ibc channel next-sequence-receive [port-id] [channel-id]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for next-sequence-receive
      --node string     <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
      --prove           show proofs for the query results (default true)
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

* [shentud query ibc channel](shentud_query_ibc_channel.md)	 - IBC channel query subcommands


