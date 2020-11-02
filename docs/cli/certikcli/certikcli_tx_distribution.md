## certikcli tx distribution

Distribution transactions subcommands

### Synopsis

Distribution transactions subcommands

```
certikcli tx distribution [flags]
```

### Options

```
  -h, --help   help for distribution
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

* [certikcli tx](certikcli_tx.md)	 - Transactions subcommands
* [certikcli tx distribution fund-community-pool](certikcli_tx_distribution_fund-community-pool.md)	 - Funds the community pool with the specified amount
* [certikcli tx distribution set-withdraw-addr](certikcli_tx_distribution_set-withdraw-addr.md)	 - change the default withdraw address for rewards associated with an address
* [certikcli tx distribution withdraw-all-rewards](certikcli_tx_distribution_withdraw-all-rewards.md)	 - withdraw all delegations rewards for a delegator
* [certikcli tx distribution withdraw-rewards](certikcli_tx_distribution_withdraw-rewards.md)	 - Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator


