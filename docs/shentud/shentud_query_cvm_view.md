## shentud query cvm view

View CVM contract

```
shentud query cvm view <address> <function> [<params>...] [flags]
```

### Options

```
      --caller string   optional caller parameter to run the view function with
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for view
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

* [shentud query cvm](shentud_query_cvm.md)	 - Querying commands for the CVM module


