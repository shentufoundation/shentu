## certik query ibc channel client-state

Query the client state associated with a channel

### Synopsis

Query the client state associated with a channel, by providing its port and channel identifiers.

```
certik query ibc channel client-state [port-id] [channel-id] [flags]
```

### Examples

```
<appd> query ibc channel client-state [port-id] [channel-id]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for client-state
      --node string     <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
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

* [certik query ibc channel](certik_query_ibc_channel.md)	 - IBC channel query subcommands


