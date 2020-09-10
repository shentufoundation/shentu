# Oracle Operator

Oracle Operator is used to listen to the `creat_task` event
from certik-chain, query the primitive and push the result back to
certik-chain.  
Ideally it is started as client process of certik chain, e.g.:
```bash
certikcli oracle-operator --home <~/.certikcli/config/oracle-operator.toml> --log_level "debug"  --from <account>
```

## How to Deploy

The oracle-operator should be mounted onto `certikcli` and started with, e.g.:
```bash
certikcli oracle-operator --log_level "debug" --query_endpoint <endpoint_url> --from <account> 
```

## How to Run

1. Register the operator on certik chain (through certikcli CLI or RESTful API) 
    and lock certain amount `ctk`.
2. Fill in the `rpc_addr` of certain (normally, `tcp://127.0.0.1:26657`), 
    primitive endpoint url and corresponding http method (`GET` or `POST`) the oracle-operator configuration file 
    (<home>/.certikcli/config/oracle-operator.toml). See template at [oracle-operator.toml](oracle-operator.toml).
3. Run the oracle operator as instructed above.
4. Wait for rewards after completing tasks. Get rewards by using withdraw command.

e.g.
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