## certik tx ibc connection open-init

Initialize connection on chain A

### Synopsis

Initialize a connection on chain A with a given counterparty chain B.
	- 'version-identifier' flag can be a single pre-selected version identifier to be used in the handshake.
	- 'version-features' flag can be a list of features separated by commas to accompany the version identifier.

```
certik tx ibc connection open-init [client-id] [counterparty-client-id] [path/to/counterparty_prefix.json] [flags]
```

### Examples

```
<appd> tx ibc connection open-init [client-id] [counterparty-client-id] [path/to/counterparty_prefix.json] --version-identifier="1.0" --version-features="ORDER_UNORDERED" --delay-period=500
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async|block) (default "sync")
      --delay-period uint           delay period that must pass before packet verification can pass against a consensus state
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible)
  -h, --help                        help for open-init
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test) (default "os")
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --memo string                 Memo to send along with transaction
      --node string                 <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --offline                     Offline mode (does not allow any online functionality
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json), this is an advanced feature
      --timeout-height uint         Set a block timeout height to prevent the tx from being committed past a certain height
      --version-features string     version features list separated by commas without spaces. The features must function with the version identifier.
      --version-identifier string   version identifier to be used in the connection handshake version negotiation
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [certik tx ibc connection](certik_tx_ibc_connection.md)	 - IBC connection transaction subcommands


