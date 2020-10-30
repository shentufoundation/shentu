## certikcli query cert certificates

Get certificates information

### Synopsis

Get certificates information

```
certikcli query cert certificates [<flags>] [flags]
```

### Options

```
      --certifier string      certificates issued by certifier
      --content string        certificates by request content
      --content-type string   type of request content
      --height int            Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                  help for certificates
      --indent                Add indent to JSON response
      --ledger                Use a connected Ledger device
      --limit int             pagination limit of certificates to query for (default 100)
      --node string           <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --page int              pagination page of certificates to to query for (default 1)
      --trust-node            Trust connected full node (don't verify proofs for responses)
```

### Options inherited from parent commands

```
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "~/.certikcli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

### SEE ALSO

* [certikcli query cert](certikcli_query_cert.md)	 - Querying commands for the certification module


