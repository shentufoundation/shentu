# Oracle Operator

Oracle Operator listens to the `creat_task` event from CertiK Chain,
queries the primitives and pushes the result back to CertiK Chain.  

## How to Config and Run

1. Register the operator on CertiK Chain (through CLI or RESTful API) and lock a certain amount `ctk`.
```bash
$ certikcli tx oracle create-operator <account address> <collateral> --name <operator name> --from <key> --fees 5000uctk -y -b block
```

2. Create the oracle-operator configuration file in `certikcli` home, e.g. `.certikcli/oracle-operator.toml` and write the following configurations in it.
See template at [oracle-operator.toml](oracle-operator.toml).
```
# configurations related to oracle operator
# strategy type
type = "linear"
# primitive configuration
[[runner.strategy.primitive]]
primitive_contract_address = "certik1r4834vyyu8vrarxgyatn34j8lsguyhn7csl0ju"
weight = 0.1
[[runner.strategy.primitive]]
primitive_contract_address = "certik1r4834vyyu8vrarxgyatn34j8lsguyhn7csl0ju"
weight = 0.1
```
The field of `primitive_contract_address` should be filled in with provided security primitive contract address.
The `weight` decides the weight of the result from the corresponding primitive to the final result.
The number of primitives is not limited.

3. Run the oracle operator by the following command.
```bash
$ certikcli oracle-operator --log_level "debug" --from <key> 
```

A sample shell script of running Oracle Operator:
```bash
# start `alice` operator in `test` mode
node0Addr=$(certikcli keys show node0 --keyring-backend test -a)
certikcli keys add alice --keyring-backend test
aliceAddr=$(certikcli keys show alice --keyring-backend test -a)
certikcli tx send node0 "$aliceAddr" 100000000000uctk --from node0 --fees 5000uctk -b block -y
certikcli tx oracle create-operator "$aliceAddr" 100000uctk --from alice --fees 5000uctk -y -b block

# start oracle operator
certikcli oracle-operator --log_level "debug" --keyring-backend test --from alice
```

### How to deploy SecurityPrimitive contract on CertiK Chain

1. Revise function `getEndpointUrl` in [SecurityPrimtive.sol](contracts/SecurityPrimitive.sol) 
to set your own url pattern.
2. Deploy your SecurityPrimitive Contract on CertiK Chain 

```bash
certikcli tx cvm deploy SecurityPrimitive.sol --args <your/person/base/endpoint/url> --from node0 --gas-prices 0.025uctk --gas-adjustment 2.0 --gas auto -y -b block
```
3. Record your Primitive Contract address `new-contract-address` from screen output.
5. Check your Primitive Contract by query `getInsight` function.
```bash
certikcli tx cvm call <your/primitive/contract/address> "getInsight" "0x00000000000000000000" "0x0100" --from node0 --gas-prices 0.025uctk --gas-adjustment 2.0 --gas auto -y -b block
```
4. Set your contract address the oracle-operator configuration file (<home>/.certikcli/config/oracle-operator.toml).
See `primitive` in template at [oracle-operator.toml](oracle-operator.toml).

## Modules

### Querier

Helper module queries the primitive endpoint for result, which receives task request from `runner` 
and deliver the result back to runner.

### Runner

Runner is in charge of receiving / dlivering messages for certik chain.

It listens to the `create_task` event published by certik chain, queries the primitive endpoint
for the task result and delivers the result back to certik chain.

### Chain `x/oracle`

The `x/oracle` module in chain handles operator registry and task management (create, response, delete, etc...)
It is located at the chain repo's `x/oracle`.