## shentud query gov tally

Get the tally of a proposal vote

### Synopsis

Query tally of votes on a proposal. You can find
the proposal-id by running "shentud query gov proposals".

Example:
$ shentud query gov tally 1

```
shentud query gov tally [proposal-id] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for tally
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

* [shentud query gov](shentud_query_gov.md)	 - Querying commands for the governance module


