## certik query ibc channel channels

Query all channels

### Synopsis

Query all channels from a chain

```
certik query ibc channel channels [flags]
```

### Examples

```
<appd> query ibc channel channels
```

### Options

```
      --count-total       count total number of records in channels to query for
      --height int        Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help              help for channels
      --limit uint        pagination limit of channels to query for (default 100)
      --node string       <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint       pagination offset of channels to query for
  -o, --output string     Output format (text|json) (default "text")
      --page uint         pagination page of channels to query for. This sets offset to a multiple of limit (default 1)
      --page-key string   pagination page-key of channels to query for
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


