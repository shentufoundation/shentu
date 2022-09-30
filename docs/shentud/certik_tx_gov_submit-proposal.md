## certik tx gov submit-proposal

Submit a proposal along with an initial deposit

### Synopsis

Submit a proposal along with an initial deposit.
Proposal title, description, type and deposit can be given directly or through a proposal JSON file.

Example:
$ <appd> tx gov submit-proposal --proposal="path/to/proposal.json" --from mykey

Where proposal.json contains:

{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "10ctk"
}

Which is equivalent to:

$ <appd> tx gov submit-proposal --title="Test Proposal" --description="My awesome proposal" --type="Text" --deposit="10ctk" --from mykey

```
certik tx gov submit-proposal [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async|block) (default "sync")
      --deposit string           The proposal deposit
      --description string       The proposal description
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible)
  -h, --help                     help for submit-proposal
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test) (default "os")
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                   Use a connected Ledger device
      --memo string              Memo to send along with transaction
      --node string              <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --offline                  Offline mode (does not allow any online functionality
      --proposal string          Proposal file path (if this path is given, other proposal flags are ignored)
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json), this is an advanced feature
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --title string             The proposal title
      --type string              The proposal Type
  -y, --yes                      Skip tx broadcasting prompt confirmation
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

* [certik tx gov](certik_tx_gov.md)	 - Governance transactions subcommands
* [certik tx gov submit-proposal cancel-software-upgrade](certik_tx_gov_submit-proposal_cancel-software-upgrade.md)	 - Cancel the current software upgrade proposal
* [certik tx gov submit-proposal certifier-update](certik_tx_gov_submit-proposal_certifier-update.md)	 - Submit a certifier update proposal
* [certik tx gov submit-proposal community-pool-spend](certik_tx_gov_submit-proposal_community-pool-spend.md)	 - Submit a community pool spend proposal
* [certik tx gov submit-proposal param-change](certik_tx_gov_submit-proposal_param-change.md)	 - Submit a parameter change proposal
* [certik tx gov submit-proposal shield-claim](certik_tx_gov_submit-proposal_shield-claim.md)	 - Submit a Shield claim proposal
* [certik tx gov submit-proposal software-upgrade](certik_tx_gov_submit-proposal_software-upgrade.md)	 - Submit a software upgrade proposal


