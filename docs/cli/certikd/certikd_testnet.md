## certikd testnet

Initialize files for a certikd testnet

### Synopsis

testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	certikd testnet --v 4 --output-dir ./output --server-ip-address 192.168.10.2
	

```
certikd testnet [flags]
```

### Options

```
      --chain-id string             genesis file chain-id, if left blank will be randomly created
      --config string               Initialization config.
      --default-password string      (default "12345678")
  -h, --help                        help for testnet
      --keyring-backend string      Select keyring's backend (os|file|test) (default "test")
      --minimum-gas-prices string   Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum.
      --node-cli-home string        Home directory of the node's cli configuration (default "certikcli")
      --node-daemon-home string     Home directory of the node's daemon configuration (default "certikd")
      --node-dir-prefix string      Prefix the directory name for each node with (node results in node0, node1, ...) (default "node")
  -o, --output-dir string           Directory to store initialization data for the testnet (default "./mytestnet")
      --port-increment int           (default 100)
      --server-ip-address string    Server IP Address (default "192.168.253.177")
      --v int                       Number of validators to initialize the testnet with (default 4)
```

### Options inherited from parent commands

```
      --home string        directory for config and data (default "~/.certikd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors
```

### SEE ALSO

* [certikd](certikd.md)	 - CertiK App Daemon (server)


