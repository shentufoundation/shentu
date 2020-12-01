# Compatibility Integration Test
This feature provides automated binary compatibility check. It sets up a two-node chain on your machine: a validator node running old binary and a non-validator node running latest binary. It then executes a sequence of transactions and constantly checks if a consensus failure happened. 

```
 ----------------                               -----------------
|                |                             |                 |
|     node 0     | port: 26656     port: 27756 |      node 1     | port: 26657
|    validator   | <-----------p2p-----------> |  non-validator  | <--------rpc(abci)
|                |                             |                 |
 ----------------                               -----------------
```
The diagram of the two-node chain is shown above. Because there's only one validator on the chain, consensus failure can only happen to the non-validator node. Therefore, we let the non-validator node listen to abci to better capture a consensus failure. 


```
.
|-- certifier_update.json
|-- certikcli (you have to place it here manually) 
|-- certikd (you have to place it here manually) 
|-- node0.sh
|-- node1.sh
|-- README.md
|-- shield_claim.json
|-- start.sh
`-- txs.sh
```
Before starting the integration test, make sure `.../shentu/tests/sync/` directory looks like the diagram above. `certikcli` and `certikd` are the old binaries run by the validator node, and you have to manually place them in this directory. `txs.sh` covers most of our custom transactions, and you're welcome to replace them with your own desired tx sequence.

To start the integration test, first you need to stop any `certikd` processes on your machine, then run `start.sh`. If something went wrong and you want to terminate the mess, run `killall certikd`. 