## certikcli query gov votes

Query votes on a proposal

### Synopsis

Query vote details for a single proposal by its identifier.

Example:
$ certikcli query gov votes 1
$ certikcli query gov votes 1 --page=2 --limit=100

```
certikcli query gov votes [proposal-id] [flags]
```

### Options

```
      --height int    Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help          help for votes
      --indent        Add indent to JSON response
      --ledger        Use a connected Ledger device
      --limit int     pagination limit of votes to query for (default 100)
      --node string   <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --page int      pagination page of votes to to query for (default 1)
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

* [certikcli query gov](certikcli_query_gov.md)	 - Querying commands for the governance module


