## shentud query gov proposals

Query proposals with optional filters

### Synopsis

Query for a all paginated proposals that match optional filters:

Example:
$ shentud query gov proposals --depositor shentu1skjwj5whet0lpe65qaq4rpq03hjxlwd9ma4udt
$ shentud query gov proposals --voter shentu1skjwj5whet0lpe65qaq4rpq03hjxlwd9ma4udt
$ shentud query gov proposals --status (DepositPeriod|CertifierVotingPeriod|ValidatorVotingPeriod|Passed|Rejected)
$ shentud query gov proposals --page=2 --limit=100

```
shentud query gov proposals [flags]
```

### Options

```
      --count-total        count total number of records in proposals to query for
      --depositor string   (optional) filter by proposals deposited on by depositor
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for proposals
      --limit uint         pagination limit of proposals to query for (default 100)
      --node string        <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
      --offset uint        pagination offset of proposals to query for
  -o, --output string      Output format (text|json) (default "text")
      --page uint          pagination page of proposals to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of proposals to query for
      --reverse            results are sorted in descending order
      --status string      (optional) filter proposals by proposal status, status: deposit_period/voting_period/passed/rejected
      --voter string       (optional) filter by proposals voted on by voted
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

* [shentud query gov](shentud_query_gov.md)	 - Querying commands for the governance module


