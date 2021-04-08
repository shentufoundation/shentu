## certik tx ibc solo create

create new solo machine client

### Synopsis

create a new solo machine client with the specified identifier and public key
	- ConsensusState json example: {"public_key":{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A/3SXL2ONYaOkxpdR5P8tHTlSlPv1AwQwSFxKRee5JQW"},"diversifier":"diversifier","timestamp":"10"}

```
certik tx ibc solo create [sequence] [path/to/consensus_state.json] [flags]
```

### Examples

```
<appd> tx ibc solo machine create [sequence] [path/to/consensus_state] --from node0 --home ../node0/<app>cli --chain-id $CID
```

### Options

```
  -a, --account-number uint           The account number of the signing account (offline mode only)
      --allow_update_after_proposal   allow governance proposal to update client
  -b, --broadcast-mode string         Transaction broadcasting mode (sync|async|block) (default "sync")
      --dry-run                       ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it
      --fees string                   Fees to pay along with transaction; eg: 10uatom
      --from string                   Name or address of private key with which to sign
      --gas string                    gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default 200000)
      --gas-adjustment float          adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string             Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                 Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible)
  -h, --help                          help for create
      --keyring-backend string        Select keyring's backend (os|file|kwallet|pass|test) (default "os")
      --keyring-dir string            The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                        Use a connected Ledger device
      --memo string                   Memo to send along with transaction
      --node string                   <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --offline                       Offline mode (does not allow any online functionality
  -s, --sequence uint                 The sequence number of the signing account (offline mode only)
      --sign-mode string              Choose sign mode (direct|amino-json), this is an advanced feature
      --timeout-height uint           Set a block timeout height to prevent the tx from being committed past a certain height
  -y, --yes                           Skip tx broadcasting prompt confirmation
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

* [certik tx ibc solo](certik_tx_ibc_solo.md)	 - Solo Machine transaction subcommands


