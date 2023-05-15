## shentud query upgrade module_versions

get the list of module versions

### Synopsis

Gets a list of module names and their respective consensus versions.
Following the command with a specific module name will return only
that module's information.

```
shentud query upgrade module_versions [optional module_name] [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for module_versions
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

* [shentud query upgrade](shentud_query_upgrade.md)	 - Querying commands for the upgrade module


