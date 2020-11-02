## certikcli query shield providers

query all providers

### Synopsis

Query providers with pagination parameters

Example:
$ certikcli query shield providers
$ certikcli query shield providers --page=2 --limit=100

```
certikcli query shield providers [flags]
```

### Options

```
      --height int    Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help          help for providers
      --indent        Add indent to JSON response
      --ledger        Use a connected Ledger device
      --limit int     pagination limit of providers to query for (default 100)
      --node string   <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --page int      pagination page of providers to to query for (default 1)
      --trust-node    Trust connected full node (don't verify proofs for responses)
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

* [certikcli query shield](certikcli_query_shield.md)	 - Querying commands for the shield module


