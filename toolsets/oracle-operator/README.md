# Oracle Operator

Oracle Operator listens to the `creat_task` event from CertiK Chain, queries the primitives and pushes the result back to CertiK Chain.

## How to Config and Run

1. Register the operator on CertiK Chain (through CLI or RESTful API) and lock a certain amount of `CTK`.
  ```bash
  $ certikcli tx oracle create-operator <account address> <collateral> --name <operator name> --from <account> --fees 5000uctk -y -b block
  ```
2. Create the oracle operator configuration file in `certikcli` home (default `.certikcli/config/oracle-operator.toml`). See template at [oracle-operator.toml](oracle-operator.toml):
  - `type`: Aggregation type, e.g. `linear`. Check [Strategy](STRATEGY.md).
  - `primitive_contract_address`: security primitive contract address.
  - `weight`: the weight of the result from the corresponding primitive to the final result.
3. Run the oracle operator by the following command.
  ```bash
  $ certikcli oracle-operator --from <account>
  ```

A sample shell script of running Oracle Operator:

```bash
certikcli tx oracle create-operator $(certikcli keys show alice --keyring-backend test -a) 100000uctk --from alice --fees 5000uctk -y -b block
certikcli oracle-operator --log_level "debug" --keyring-backend test --from alice
```

## Support of Multiple Client Chain

Contract addresses in security oracle tasks are prefixed with the chain identifier, e.g. `eth:0xabc`(Ethereum), `bsc:0xdef`(Binance Smart Chain). To enable oracle operator handle tasks for multiple chains, the configuration file can be specified as:

```toml
[strategy.eth]
type = "linear"
[[stragety.eth.primitive]]
primitive_contract_address = "certik111..."
weight = 0.1
[[stragety.eth.primitive]]
primitive_contract_address = "certik222..."
weight = 0.1

[strategy.bsc]
type = "linear"
[[stragety.bsc.primitive]]
primitive_contract_address = "certik333..."
weight = 0.1
[[stragety.bsc.primitive]]
primitive_contract_address = "certik444..."
weight = 0.1
```

### How to Deploy SecurityPrimitive Contract on CertiK Chain

Security Primitive Contract wraps a security endpoint to chain.

1. Set function `getEndpointUrl` based on template at [SecurityPrimtive.sol](contracts/SecurityPrimitive.sol).
2. Deploy your SecurityPrimitive Contract on CertiK Chain
  ```bash
  certikcli tx cvm deploy SecurityPrimitive.sol --args <security-primitive-endpoint> --from <account> --gas-prices 0.025uctk --gas-adjustment 2.0 --gas auto -y -b block
  ```
3. Record your SecurityPrimitive Contract's address `new-contract-address` from screen output.
4. Check your PrimitivePrimitive Contract by querying `getInsight` function.
  ```bash
  certikcli tx cvm call <primitive-contract-address> "getInsight" "0x00000000000000000000" "0x0100" --from <account> --gas-prices 0.025uctk --gas-adjustment 2.0 --gas auto -y -b block
  ```
5. Set your contract address in the oracle-operator configuration file (`~/.certikcli/config/oracle-operator.toml`).

## Modules

### Chain `x/oracle`

The `x/oracle` module in chain handles operator registry and task management (create, response, delete, etc...)

## Oracle Operator

The oracle operator listens to the `create_task` event published by certik chain, queries the primitive endpoint for the task result and delivers the result back to certik chain.
