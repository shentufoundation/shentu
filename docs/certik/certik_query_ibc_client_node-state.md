## certik query ibc client node-state

Query a node consensus state

### Synopsis

Query a node consensus state. This result is feed to the client creation transaction.

```
certik query ibc client node-state [flags]
```

### Examples

```
<appd> query ibc client node-state
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for node-state
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

* [certik query ibc client](certik_query_ibc_client.md)	 - IBC client query subcommands


