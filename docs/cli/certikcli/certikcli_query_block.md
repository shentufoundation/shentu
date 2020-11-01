## certikcli query block

Get verified data for a the block at given height

### Synopsis

Get verified data for a the block at given height

```
certikcli query block [height] [flags]
```

### Options

```
  -h, --help          help for block
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


