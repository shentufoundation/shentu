## certik tx ibc tendermint-client create

create new tendermint client

### Synopsis

Create a new tendermint IBC client.
  - 'trust-level' flag can be a fraction (eg: '1/3') or 'default'
  - 'proof-specs' flag can be JSON input, a path to a .json file or 'default'
  - 'upgrade-path' flag is a string specifying the upgrade path for this chain where a future upgraded client will be stored. The path is a comma-separated list representing the keys in order of the keyPath to the committed upgraded client.
  e.g. 'upgrade/upgradedClient'

```
certik tx ibc tendermint-client create [path/to/consensus_state.json] [trusting_period] [unbonding_period] [max_clock_drift] [flags]
```

### Examples

```
<appd> tx ibc tendermint-client create [path/to/consensus_state.json] [trusting_period] [unbonding_period] [max_clock_drift] --trust-level default --consensus-params [path/to/consensus-params.json] --proof-specs [path/to/proof-specs.json] --upgrade-path upgrade/upgradedClient --from node0 --home ../node0/<app>cli --chain-id $CID
```

### Options

```
  -a, --account-number uint               The account number of the signing account (offline mode only)
      --allow_update_after_expiry         allow governance proposal to update client after expiry
      --allow_update_after_misbehaviour   allow governance proposal to update client after misbehaviour
  -b, --broadcast-mode string             Transaction broadcasting mode (sync|async|block) (default "sync")
      --dry-run                           ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it
      --fees string                       Fees to pay along with transaction; eg: 10uatom
      --from string                       Name or address of private key with which to sign
      --gas string                        gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default 200000)
      --gas-adjustment float              adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string                 Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                     Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible)
  -h, --help                              help for create
      --keyring-backend string            Select keyring's backend (os|file|kwallet|pass|test) (default "os")
      --keyring-dir string                The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                            Use a connected Ledger device
      --memo string                       Memo to send along with transaction
      --node string                       <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --offline                           Offline mode (does not allow any online functionality
      --proof-specs string                proof specs format to be used for verification (default "default")
  -s, --sequence uint                     The sequence number of the signing account (offline mode only)
      --sign-mode string                  Choose sign mode (direct|amino-json), this is an advanced feature
      --timeout-height uint               Set a block timeout height to prevent the tx from being committed past a certain height
      --trust-level string                light client trust level fraction for header updates (default "default")
  -y, --yes                               Skip tx broadcasting prompt confirmation
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

* [certik tx ibc tendermint-client](certik_tx_ibc_tendermint-client.md)	 - Tendermint client transaction subcommands


