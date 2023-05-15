## shentud tx authz grant

Grant authorization to an address

### Synopsis

grant authorization to an address to execute a transaction on your behalf:

Examples:
 $ shentud tx authz grant cosmos1skjw.. send /cosmos.bank.v1beta1.MsgSend --spend-limit=1000stake --from=cosmos1skl..
 $ shentud tx authz grant cosmos1skjw.. generic --msg-type=/cosmos.gov.v1beta1.MsgVote --from=cosmos1sk..

```
shentud tx authz grant <grantee> <authorization_type="send"|"generic"|"delegate"|"unbond"|"redelegate"> --from <granter> [flags]
```

### Options

```
  -a, --account-number uint          The account number of the signing account (offline mode only)
      --allowed-validators strings   Allowed validators addresses separated by ,
  -b, --broadcast-mode string        Transaction broadcasting mode (sync|async|block) (default "sync")
      --deny-validators strings      Deny validators addresses separated by ,
      --dry-run                      ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --expiration int               The Unix timestamp. Default is one year. (default 1714208473)
      --fee-account string           Fee account pays fees for the transaction instead of deducting from the signer
      --fees string                  Fees to pay along with transaction; eg: 10uatom
      --from string                  Name or address of private key with which to sign
      --gas string                   gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default 200000)
      --gas-adjustment float         adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string            Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible)
  -h, --help                         help for grant
      --keyring-backend string       Select keyring's backend (os|file|kwallet|pass|test|memory) (default "os")
      --keyring-dir string           The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                       Use a connected Ledger device
      --msg-type string              The Msg method name for which we are creating a GenericAuthorization
      --node string                  <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --note string                  Note to add a description to the transaction (previously --memo)
      --offline                      Offline mode (does not allow any online functionality
  -o, --output string                Output format (text|json) (default "json")
  -s, --sequence uint                The sequence number of the signing account (offline mode only)
      --sign-mode string             Choose sign mode (direct|amino-json), this is an advanced feature
      --spend-limit string           SpendLimit for Send Authorization, an array of Coins allowed spend
      --timeout-height uint          Set a block timeout height to prevent the tx from being committed past a certain height
  -y, --yes                          Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data (default "~/.shentud")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [shentud tx authz](shentud_tx_authz.md)	 - Authorization transactions subcommands


