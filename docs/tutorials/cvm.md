## `x/cvm` module in Shentu

### Introduction
`cvm` module means CertiK Virtual Machine which provides smart contract deployment and execution. It handles the underlying gas calculation and corresponding error handling for contract deployments and calls. 

CertiK believes that implementing smart contracts and blockchain security will give chain operation and VM execution more dynamic, actionable value. A secure smart contract, for example, may choose to interact with secure and non-secure smart contracts differently. This type of security-level difference is common and beneficial in real life.

On the CertiK Virtual Machine, a smart contract can validate the security of any other smart contract, in cryptographic or mathematical forms, allowing for differentiated handling of secure and non-secure peers. Cryptographic certifications are included in on-chain smart contracts as proof of verified security.

The CertiK Virtual Machine's (CVM) goals are:

 -  To provide advanced security for the VM code. The attack vectors from insecure code will be limited by this hyper-secure trusted computing base (TCB). The CVM will eventually impose sandboxing and isolation of any code that is not fully certified via mathematical proofs, while also being fully formally verified.
 - To make security intelligence an expressible value on the blockchain. This is impossible in any of today's blockchain virtual machines since security analysis is done and kept off-chain, making it impossible for programs to refer to the results while considering a transaction. This reduces the real value of security analysis because it is up to an individual to do their due diligence and identify and examine any previous audit reports before using a smart contract. More dynamic and actionable actions are enabled by extending this security intelligence on-chain. For example, while interacting with secure and non-secure smart contracts, a secure smart contract may choose to differentiate its actions; in real world, this differentiation is similar to lenders charging various rates based on a person's credit score (or in this case, security score).

The CVM gives VM code access to smart contract and blockchain security information, allowing for new ways to access, check, rely on, and even dynamically construct blockchain and smart contract security. On-chain security intelligence gives users access to data that helps them make better decisions. While unaudited/unverified smart contracts can still execute in the CVM, people who want to interact with it have more transparency. Additionally, CVM establishes ​a hierarchy of VM code security​, allowing vendors and technologies to co-exist and collaborate with clarity.


 ### `x/cvm` module in Shentu
#### Transactions and Queries 
#### **Transactions**
In order to use the `cvm` module, the user needs to have a smart contract (ERC20) written in .sol file and compile it into .abi and binary files to execute transaction.
- `certik tx cvm deploy <filename> [flags]`: Deploy CVM contract(s). After a smart contract is compiled intoc .abi and binary file to use this command. We need to deploy smart contract in order for it to be available to user. `<filename>`  is the binary that is compiled earlier and goes with flag `--abi <filename>.abi` to deploy the contract.
```{engine = 'sh'}
$ certik tx cvm deploy simple.bin --abi simple.abi --from jack -y --chain-id certikchain

{"height":"1432","txhash":"AED8C146B37B2725CA30EBDD4CC104F28777ED8965C2A8904ED4714F6A77FEF0","codespace":"","code":0,"data":"0A200A066465706C6F7912160A14A8B47B8B87F2B6715477646D8BE0DF176CF4140B","raw_log":"[{\"events\":[{\"type\":\"deploy\",\"attributes\":[{\"key\":\"new-contract-address\",\"value\":\"certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm\"},{\"key\":\"value\",\"value\":\"0\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"deploy\"},{\"key\":\"module\",\"value\":\"cvm\"},{\"key\":\"sender\",\"value\":\"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p\"},{\"key\":\"amount\",\"value\":\"0\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"deploy","attributes":[{"key":"new-contract-address","value":"certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm"},{"key":"value","value":"0"}]},{"type":"message","attributes":[{"key":"action","value":"deploy"},{"key":"module","value":"cvm"},{"key":"sender","value":"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p"},{"key":"amount","value":"0"}]}]}],"info":"","gas_wanted":"200000","gas_used":"80255","tx":null,"timestamp":""}
```

End users can engage with smart contracts by exposing certain functions. When a user wants to call a function, they create a transaction on the blockchain that includes the name of the function and its parameters. When a transaction including a function call is mined and published on the Ethereum network, each machine on the network will execute the function in the smart contract to the application's state in a predictable manner. This entails carrying out the monetary transfers and other state variable modifications that the code specifies.

- `certik tx cvm call <address> <function> [<params>...] [flags]`: Call CVM contract. For example, after the smart contract is deployed, we can call cvm contract with the address of contract and the function inside the contract: 
```{engine = 'sh'}
$ certik tx cvm call certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm set 123 --from jack -y --chain-id certikchain

{"height":"1498","txhash":"08D204505DB059E293FF0979FF0610D5BC9AD0ACA2F659FC72FB5A006D160AC4","codespace":"","code":0,"data":"0A060A0463616C6C","raw_log":"[{\"events\":[{\"type\":\"call\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm\"},{\"key\":\"value\",\"value\":\"0\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"call\"},{\"key\":\"module\",\"value\":\"cvm\"},{\"key\":\"sender\",\"value\":\"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p\"},{\"key\":\"amount\",\"value\":\"0\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"call","attributes":[{"key":"recipient","value":"certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm"},{"key":"value","value":"0"}]},{"type":"message","attributes":[{"key":"action","value":"call"},{"key":"module","value":"cvm"},{"key":"sender","value":"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p"},{"key":"amount","value":"0"}]}]}],"info":"","gas_wanted":"200000","gas_used":"70300","tx":null,"timestamp":""}
```

#### **Queries**
- `certik query cvm abi <address> [flags]`: Get CVM contract code ABI with contract's address. 
```{engine = 'sh'}
$ certik query cvm abi certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm

abi: '[{"constant":true,"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]'
```

- `certik query cvm address-translate <address> [flags]`: Translate a Bech32 address to hex and vice versa. 
```{engine = 'sh'}
$ certik query cvm address-translate certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm

a8b47b8b87f2b6715477646d8be0df176cf4140b
```

- `certik query cvm code <address> [flags]`: Get CVM contract code with contract's address. 
```{engine = 'sh'}
$ certik query cvm code certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm

code: 6080604052348015600f57600080fd5b506004361060325760003560e01c806360fe47b11460375780636d4ce63c146053575b600080fd5b605160048036036020811015604b57600080fd5b5035606b565b005b60596070565b60408051918252519081900360200190f35b600055565b6000549056fea265627a7a723158204241a0de8e4dfaea6fc8bc8f1508a06859552a76a48f4fdafa78d9c8bd05419764736f6c63430005110032
```

- `certik query cvm contract [address] [flags]`: Query contract info by address, revert to query normal account info if the address is not a contract. 
```{engine = 'sh'}
$ certik query cvm contract certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm

Address: A8B47B8B87F2B6715477646D8BE0DF176CF4140B
Balance: "0"
CodeHash: ""
ContractMeta: []
EVMCode: 6080604052348015600F57600080FD5B506004361060325760003560E01C806360FE47B11460375780636D4CE63C146053575B600080FD5B605160048036036020811015604B57600080FD5B5035606B565B005B60596070565B60408051918252519081900360200190F35B600055565B6000549056FEA265627A7A723158204241A0DE8E4DFAEA6FC8BC8F1508A06859552A76A48F4FDAFA78D9C8BD05419764736F6C63430005110032
Forebear: null
NativeName: ""
Permissions:
Base:
Perms: call | createContract
SetBit: root | send | call | createContract | createAccount | bond | name | proposal
| input | batch | identify | hasBase | setBase | unsetBase | setGlobal | hasRole
| addRole | removeRole
Roles: []
PublicKey: null
Sequence: "0"
WASMCode: ""
```

- `certik query cvm meta <address,hash> [flags]`: Get CVM Metadata hash for an address or Metadata for a hash. 
```{engine = 'sh'}
$ certik query cvm meta certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm
#null output
```

- `certik query cvm storage <address> <key> [flags]`: Get CVM storage data. We use the contract address deployed earlier and the txHash to make a query: 
```{engine = 'sh'}
$ certik query cvm storage certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm 08D204505DB059E293FF0979FF0610D5BC9AD0ACA2F659FC72FB5A006D160AC4

value: AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
```

- `certik query cvm view <address> <function> [<params>...] [flags]`: View CVM contract with contract's address and the function inside the .sol file, the flag --caller should be address who deploy the contract. 
```{engine = 'sh'}
$ certik query cvm view certik14z68hzu872m8z4rhv3kchcxlzak0g9qtk8kjkm  get --caller certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p

return_vars:
	- name: "0"
	value: "123"
```


