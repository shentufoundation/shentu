### How to deploy SecurityPrimitive contract on CertiK Chain

1. Revise function `getEndpointUrl` in [SecurityPrimitive.sol](SecurityPrimitive.sol) 
to set your own url pattern.
    
    e.g. `return string(abi.encodePacked(_endpoint, "?address=", contractAddress, "&functionSignature=", functionSignature));`

    Or custom your personal score evaluation method and change the first return value of
    `getInsight` to true, which represents `isUrl`.
    
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
See `primitive` in template at [oracle-operator.toml](../oracle-operator.toml).
