## certikcli query tx

Query for a transaction by hash in a committed block

### Synopsis

Query for a transaction by hash in a committed block

```
certikcli query tx [hash] [flags]
```

### Options

```
  -h, --help          help for tx
  -n, --node string   Node to connect to (default "tcp://localhost:26657")
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


