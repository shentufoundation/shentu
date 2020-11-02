## certikcli keys show

Show key info for the given name

### Synopsis

Return public details of a single local key. If multiple names are
provided, then an ephemeral multisig key will be created under the name "multi"
consisting of all the keys provided by name and multisig threshold.

```
certikcli keys show [name [name...]] [flags]
```

### Options

```
  -a, --address                   Output the address only (overrides --output)
      --bech string               The Bech32 prefix encoding for a key (acc|val|cons) (default "acc")
  -d, --device                    Output the address in a ledger device
  -h, --help                      help for show
      --indent                    Add indent to JSON response
      --multisig-threshold uint   K out of N required signatures (default 1)
  -p, --pubkey                    Output the public key only (overrides --output)
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


