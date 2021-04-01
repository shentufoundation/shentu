## certik query cert certificates

Get certificates information

```
certik query cert certificates [<flags>] [flags]
```

### Options

```
      --certifier string      certificates issued by certifier
      --content string        certificates by request content
      --content-type string   type of request content
      --count-total           count total number of records in certificates to query for
      --height int            Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                  help for certificates
      --limit uint            pagination limit of certificates to query for (default 100)
      --node string           <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint           pagination offset of certificates to query for
  -o, --output string         Output format (text|json) (default "text")
      --page uint             pagination page of certificates to query for. This sets offset to a multiple of limit (default 1)
      --page-key string       pagination page-key of certificates to query for
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

* [certik query cert](certik_query_cert.md)	 - Querying commands for the certification module


