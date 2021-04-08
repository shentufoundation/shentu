## certik query ibc channel connections

Query all channels associated with a connection

### Synopsis

Query all channels associated with a connection

```
certik query ibc channel connections [connection-id] [flags]
```

### Examples

```
<appd> query ibc channel connections [connection-id]
```

### Options

```
      --count-total       count total number of records in channels associated with a connection to query for
      --height int        Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help              help for connections
      --limit uint        pagination limit of channels associated with a connection to query for (default 100)
      --node string       <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint       pagination offset of channels associated with a connection to query for
  -o, --output string     Output format (text|json) (default "text")
      --page uint         pagination page of channels associated with a connection to query for. This sets offset to a multiple of limit (default 1)
      --page-key string   pagination page-key of channels associated with a connection to query for
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


