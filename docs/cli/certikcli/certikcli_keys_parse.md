## certikcli keys parse

Parse address from hex to bech32 and vice versa

### Synopsis

Convert and print to stdout key addresses and fingerprints from
hexadecimal into bech32 cosmos prefixed format and vice versa.


```
certikcli keys parse <hex-or-bech32-address> [flags]
```

### Options

```
  -h, --help     help for parse
      --indent   Indent JSON output
```

### Options inherited from parent commands

```
      --chain-id string          Chain ID of tendermint node
  -e, --encoding string          Binary encoding (hex|b64|btc) (default "hex")
      --home string              directory for config and data (default "~/.certikcli")
      --keyring-backend string   Select keyring's backend (os|file|test) (default "os")
  -o, --output string            Output format (text|json) (default "text")
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [certikcli keys](certikcli_keys.md)	 - Add or view local private keys


