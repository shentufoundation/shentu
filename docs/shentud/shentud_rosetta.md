## shentud rosetta

spin up a rosetta server

```
shentud rosetta [flags]
```

### Options

```
      --addr string         the address rosetta will bind to (default ":8080")
      --blockchain string   the blockchain type (default "app")
      --grpc string         the app gRPC endpoint (default "localhost:9090")
  -h, --help                help for rosetta
      --network string      the network name (default "network")
      --offline             run rosetta only with construction API
      --retries int         the number of retries that will be done before quitting (default 5)
      --tendermint string   the tendermint rpc endpoint, without tcp:// (default "localhost:26657")
```

### Options inherited from parent commands

```
      --home string         directory for config and data (default "~/.shentud")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [shentud](shentud.md)	 - Stargate Shentu Chain App


