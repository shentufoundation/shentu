## certikcli query tendermint-validator-set

Get the full tendermint validator set at given height

### Synopsis

Get the full tendermint validator set at given height

```
certikcli query tendermint-validator-set [height] [flags]
```

### Options

```
  -h, --help          help for tendermint-validator-set
      --indent        indent JSON response
      --limit int     Query number of results returned per page (default 100)
  -n, --node string   Node to connect to (default "tcp://localhost:26657")
      --page int      Query a specific page of paginated results
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

* [certikcli query](certikcli_query.md)	 - Querying subcommands


