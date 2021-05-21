# CertiK Chain

[![PkgGoDev](https://pkg.go.dev/badge/github.com/certikfoundation/shentu)](https://pkg.go.dev/github.com/certikfoundation/shentu)
<a href="https://circleci.com/gh/certikfoundation/shentu/tree/master">
<img src="https://circleci.com/gh/certikfoundation/shentu/tree/master.svg?style=svg&circle-token=b948d67100954a74a11c21fbb8cb6202b83e5f3a">
</a>
<a href="https://codecov.io/gh/certikfoundation/shentu">
<img src="https://codecov.io/gh/certikfoundation/shentu/branch/master/graph/badge.svg">
</a>

# Prerequisites

Install Golang 1.14+ following instructions from https://golang.org/doc/install.

Install Solidity Compiler following instructions from https://solidity.readthedocs.io/en/latest/installing-solidity.html#binary-packages.

# Development

Please install `solc` first.

During the development, for cross-referencing dependent packages' source, run the following commandline will get them into `/vendor`.

```bash
$ go mod vendor
```

# Style

- Mechanically copied code (that we intend to copy again in the future) from external dependencies should be kept in original style as close as possible.

- All exported definitions (types, methods, functions, variables, constants, etc.) should have _whole sentence_ comments starting with the definition name.

- Group imports in the following order, separated by an empty line in between.

  - Standard packages
  - Other 3rd-party packages
  - Tendermint packages
  - Cosmos SDK packages
  - Burrow packages
  - CertiK packages

- Before creating pull requests of your changes, run the following commandlines for style and dependencies.

  ```bash
  $ make tidy
  ```

  Also make sure to fix linter and testing errors.

  ```bash
  $ make lint
  $ make test
  ```

  [golangci-lint](https://github.com/golangci/golangci-lint) needs to be installed locally to run the linters.

# Build

```bash
$ make install
```

# Unit Test

Run unit test:

```bash
$ make test
```

Run unit test and view brief coverage in cli:

```bash
$ make test-cov
```

Run unit test and view detailed coverage report in browser:

```bash
$ make test-cov-html
```

# Build Docker Image

```bash
$ make image
```

This builds the docker image for the shentu node daemon & client. (Note that this command needs to be run prior to the multi node testnet setup below.)

# Multi Node Testnet

Run `make image` first to build the images initially (no need to be rerun for usual code update).

Setup a localnet (per code change):

```bash
$ make localnet
```

This will run (four) testnet nodes in the background as well as an interactive client. If you want to see node output, do `docker-compose logs -f` (in other terminal) which will display output from all four nodes until you exit the process. Stops the localnet by `docker-compose down`.

Enter the client to interact with the localnet:

```bash
$ make localnet.client
```

You can also setup a localnet and start a client in one command:

```bash
$ make localnet.both
```

Run the scripts for testing:

```bash
$ certik query account $NODE0_KEY
$ certik query supply total uctk

# Create a new account
$ certik keys add jack
$ JACK_KEY=$(certik keys show jack -a)
# default password = jkljkljkl
$ certik tx send node0 $JACK_KEY 100000uctk --gas-prices=0.025uctk --from node0
$ certik query account $JACK_KEY
```

Finally, run `make localnet.down` to shutdown the localnet.

# Deploy Testnet on AWS

Please see e2e/README.md for instructions.

# Single Node Chain (without Docker)

Recommend the use of the provided Docker containers for simplicity, but here are instructions for manually setting up a single node chain for testing purposes.

```bash
$ certik unsafe-reset-all
$ rm -rf ~/.certik
$ rm -rf ~/.certik

$ certik init node0 --chain-id certikchain

$ certik config chain-id certikchain
$ certik keys add jack

$ certik add-genesis-account $(certik keys show jack -a) 200000000uctk
```

Notification: Every transaction will need 5000uctk (certik token), so you'd better start with more than 5000uctk here.

```bash
$ certik gentx --name jack --amount 100000000uctk
$ certik collect-gentxs

$ certik start

$ certik query account $(certik keys show jack -a)

$ certik keys add alice
$ certik tx send jack $(certik keys show alice -a) 70000uctk --gas-prices=0.025uctk --from jack
$ certik query account $(certik keys show jack -a)
$ certik query account $(certik keys show alice -a)
```

# CVM module test commands

## simple.sol
We have a very simple storage setter / getter contract `tests/simple.sol`.

```solidity
pragma solidity >=0.4.0 <0.7.0;

contract SimpleStorage {
    uint storedData;

    function set(uint x) public {
        storedData = x;
    }

    function get() public view returns (uint) {
        return storedData;
    }
}
```

To deploy `tests/simple.sol` contract (be sure to have `solc` installed), first compile it into a bytecode and abi.

Assuming you have `simple.sol` compiled into `simple.bytecode` and `simple.abi`
```bash
$ cd tests
$ certik tx cvm deploy simple.bytecode --abi simple.abi --from node0
```

Printed on the main terminal

```bash
Response:
  TxHash: 8067DBC001BE239E5A44843CCEF4C71A87B802352989F97664AF8F265E7B888E
```

Printed on the certik node server terminals

```bash
I[2019-06-27|09:05:33.281] CVM Start                                    module=main txHash=8067dbc001be239e5a44843ccef4c71a87b802352989f97664af8f265e7b888e
I[2019-06-27|09:05:33.282] CVM                                          module=main log_channel=Trace scope=NewVM message="(1) (89203146CD64A945BE5B02B0C44FB0040855C02A) 07BC7F3C21C34643A90AA1138C950FAC5025B693 (code=229) gas: 1000000 (d) 608060405234801561001057600080FD5B5060C68061001F6000396000F3FE6080604052348015600F57600080FD5B506004361060325760003560E01C806360FE47B11460375780636D4CE63C146062575B600080FD5B606060048036036020811015604B57600080FD5B8101908080359060200190929190505050607E565B005B60686088565B6040518082815260200191505060405180910390F35B8060008190555050565B6000805490509056FEA265627A7A723058205FEC64D09C278453AB74A855DCC214EA05BF9541E35E851AF41570397593055564736F6C63430005090032\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.283] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 0   (op) PUSH1          (st) 0    (gas) 1000000" tag=DebugOpcodes
I[2019-06-27|09:05:33.283] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0x0000000000000000000000000000000000000000000000000000000000000080\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.283] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 2   (op) PUSH1          (st) 1    (gas) 999999" tag=DebugOpcodes
I[2019-06-27|09:05:33.284] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0x0000000000000000000000000000000000000000000000000000000000000040\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.284] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 4   (op) MSTORE         (st) 2    (gas) 999998" tag=DebugOpcodes
I[2019-06-27|09:05:33.284] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0x0000000000000000000000000000000000000000000000000000000000000080 @ 0x40\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.284] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 5   (op) CALLVALUE      (st) 0    (gas) 999996" tag=DebugOpcodes
I[2019-06-27|09:05:33.284] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.284] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 6   (op) DUP1           (st) 1    (gas) 999995" tag=DebugOpcodes
I[2019-06-27|09:05:33.284] CVM                                          module=main log_channel=Trace scope=NewVM message=" => [1] 0x0000000000000000000000000000000000000000000000000000000000000000\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.284] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 7   (op) ISZERO         (st) 2    (gas) 999993" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message=" 0000000000000000000000000000000000000000000000000000000000000000 == 0 = 1\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 8   (op) PUSH2          (st) 2    (gas) 999991" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0x0000000000000000000000000000000000000000000000000000000000000010\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 11  (op) JUMPI          (st) 3    (gas) 999990" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message=" ~> 16\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 16  (op) JUMPDEST       (st) 1    (gas) 999988" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message="\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 17  (op) POP            (st) 1    (gas) 999988" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0x0000000000000000000000000000000000000000000000000000000000000000\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 18  (op) PUSH1          (st) 0    (gas) 999987" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0x00000000000000000000000000000000000000000000000000000000000000C6\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.285] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 20  (op) DUP1           (st) 1    (gas) 999986" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message=" => [1] 0x00000000000000000000000000000000000000000000000000000000000000C6\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 21  (op) PUSH2          (st) 2    (gas) 999984" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0x000000000000000000000000000000000000000000000000000000000000001F\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 24  (op) PUSH1          (st) 3    (gas) 999983" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0x0000000000000000000000000000000000000000000000000000000000000000\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 26  (op) CODECOPY       (st) 4    (gas) 999982" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message=" => [0, 31, 198] 6080604052348015600F57600080FD5B506004361060325760003560E01C806360FE47B11460375780636D4CE63C146062575B600080FD5B606060048036036020811015604B57600080FD5B8101908080359060200190929190505050607E565B005B60686088565B6040518082815260200191505060405180910390F35B8060008190555050565B6000805490509056FEA265627A7A723058205FEC64D09C278453AB74A855DCC214EA05BF9541E35E851AF41570397593055564736F6C63430005090032\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 27  (op) PUSH1          (st) 1    (gas) 999979" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message=" => 0x0000000000000000000000000000000000000000000000000000000000000000\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message="(pc) 29  (op) RETURN         (st) 2    (gas) 999978" tag=DebugOpcodes
I[2019-06-27|09:05:33.286] CVM                                          module=main log_channel=Trace scope=NewVM message=" => [0, 198] (198) 0x6080604052348015600F57600080FD5B506004361060325760003560E01C806360FE47B11460375780636D4CE63C146062575B600080FD5B606060048036036020811015604B57600080FD5B8101908080359060200190929190505050607E565B005B60686088565B6040518082815260200191505060405180910390F35B8060008190555050565B6000805490509056FEA265627A7A723058205FEC64D09C278453AB74A855DCC214EA05BF9541E35E851AF41570397593055564736F6C63430005090032\n" tag=DebugOpcodes
I[2019-06-27|09:05:33.287] CVM Stop                                     module=main result=07bc7f3c21c34643a90aa1138c950fac5025b693
```

To inspect deploy transaction details and to obtain Bech32 contract address

```bash
$ certik query tx 8067DBC001BE239E5A44843CCEF4C71A87B802352989F97664AF8F265E7B888E
Response:
  Height: 169
  TxHash: 8067DBC001BE239E5A44843CCEF4C71A87B802352989F97664AF8F265E7B888E
  Data: 07BC7F3C21C34643A90AA1138C950FAC5025B693
  Raw Log: [{"msg_index":"0","success":true,"log":"certik1q77870ppcdry82g25yfce9g043gztd5nd3z8uy"}]
  Logs: [{"msg_index":0,"success":true,"log":"certik1q77870ppcdry82g25yfce9g043gztd5nd3z8uy"}]
  GasWanted: 200000
  GasUsed: 41849
  Tags:
    - action = deploy

  Timestamp: 2019-06-27T16:05:27Z
```

To inspect contract code bytes deployed at `certik1q77870ppcdry82g25yfce9g043gztd5nd3z8uy`

```bash
$ certik query cvm code certik1q77870ppcdry82g25yfce9g043gztd5nd3z8uy
6080604052348015600F57600080FD5B506004361060325760003560E01C806360FE47B114603757
80636D4CE63C146062575B600080FD5B606060048036036020811015604B57600080FD5B81019080
80359060200190929190505050607E565B005B60686088565B604051808281526020019150506040
5180910390F35B8060008190555050565B6000805490509056FEA265627A7A723058205FEC64D09C
278453AB74A855DCC214EA05BF9541E35E851AF41570397593055564736F6C63430005090032
```

To call `SimpleStorage.set(123)` at the contract

```bash
$ certik tx cvm call certik1q77870ppcdry82g25yfce9g043gztd5nd3z8uy set 123 --from node0
```

Then we can verify the storage setting by calling `SimpleStorage.get()` at the
contract

```bash
$ certik tx cvm call certik1q77870ppcdry82g25yfce9g043gztd5nd3z8uy get --from node0
Response:
  TxHash: 6EAABFDF5022F21F88D9DBBE8A0837F3CF6819F06F801BEF301188F22DF16C9B
```

We can inspect and verify the read out value is indeed 123 (0x7b) by either
looking at the certik node server terminal

```bash
I[2019-06-27|09:20:44.674] CVM Stop                                     module=main result=000000000000000000000000000000000000000000000000000000000000007b
```

or query the get call transaction

```bash
$ certik query tx 6EAABFDF5022F21F88D9DBBE8A0837F3CF6819F06F801BEF301188F22DF16C9B
Response:
  Height: 333
  TxHash: 6EAABFDF5022F21F88D9DBBE8A0837F3CF6819F06F801BEF301188F22DF16C9B
  Data: 000000000000000000000000000000000000000000000000000000000000007B
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 17677
  Tags:
    - action = call

  Timestamp: 2019-06-27T16:20:38Z
```

## DeepSEA Integration

Analogous to the steps above, a DeepSEA contract `contract.ds` can be deployed with:

```bash
$ certik tx cvm deploy contract.ds --from node0
```

Make sure that `dsc` (DeepSEA compiler) is in your `PATH`.

# REST test commands

Start the `rest-server` in another terminal window:

```bash
$ certik rest-server --trust-node
```

Back to the previous terminal window where `NODE0_KEY` is defined:

```bash
# Initiate Transaction: Burn 2 ctk from node0.
$ curl -XPOST -s http://localhost:1317/ctk/burn --data-binary '{"base_req":{"from":"'$NODE0_KEY'","chain_id":"certikchain"},"src":"'$NODE0_KEY'","amount":"2"}' > unsignedTx.json

# Sign Transaction
# Note: sequence and account-number can be found in account information
$ certik tx sign unsignedTx.json --from $NODE0_KEY --offline --chain-id certikchain --sequence 5 --account-number 0 > signedTx.json

# Broadcast Transaction
$ certik tx broadcast signedTx.json

# Check the balance of node0
$ curl -s http:/localhost:1317/ctk/balance/$NODE0_KEY
```
