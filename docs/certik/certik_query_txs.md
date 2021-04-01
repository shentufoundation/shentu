## certik query txs

Query for paginated transactions that match a set of events

### Synopsis

Search for transactions that match the exact given events where results are paginated.
Each event takes the form of '{eventType}.{eventAttribute}={value}'. Please refer
to each module's documentation for the full set of events to query for. Each module
documents its respective events under 'xx_events.md'.

Example:
$ <appd> query txs --events 'message.sender=cosmos1...&message.action=withdraw_delegator_reward' --page 1 --limit 30

```
certik query txs [flags]
```

### Options

```
      --events string            list of transaction events in the form of {eventType}.{eventAttribute}={value}
  -h, --help                     help for txs
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test) (default "os")
      --limit int                Query number of transactions results per page returned (default 30)
  -n, --node string              Node to connect to (default "tcp://localhost:26657")
      --page int                 Query a specific page of paginated results (default 1)
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

* [certik query](certik_query.md)	 - Querying subcommands


