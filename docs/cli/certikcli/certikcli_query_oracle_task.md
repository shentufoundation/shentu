## certikcli query oracle task

Get task information

### Synopsis

Get task information

```
certikcli query oracle task <flags> [flags]
```

### Options

```
      --contract string   Provide the contract address
      --function string   Provide the function
      --height int        Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help              help for task
      --indent            Add indent to JSON response
      --ledger            Use a connected Ledger device
      --node string       <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --trust-node        Trust connected full node (don't verify proofs for responses)
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

* [certikcli query oracle](certikcli_query_oracle.md)	 - Oracle staking subcommands


