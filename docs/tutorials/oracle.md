## `x/oracle` module in Shentu

### Introduction
`oracle` module wants to break down complex audit reports into smaller security primitives that can be invoked on-chain to verify a smart contract's security in real time. These Security Oracle scores are dynamic, aggregating scores and providing insights into the reliability of the underlying code by querying most latest security primitives and tests.

CertiK Security Oracles (`oracle` module) can also be used to submit requests for unaudited smart contracts. These requests are forwarded to a decentralized group of security operators that compete for the CTK transaction fee. The CertiK Oracle Combinator integrates the different results from each operator into an on-chain score. The transaction's fees are shared among the operators who contributed security primitives to the request.

Users can make better decisions about their potential transactions and external invocations by using CertiK Security Oracles to retrieve security intelligence. This decentralized information system gives communities the ability to conduct real-time security checks, such as those active in the booming DeFi ecosystem. This innovation, in the spirit of complete decentralization, decentralizes security intelligence from a few security auditors to the entire blockchain community, making it available on-chain on demand.

#### Security Oracle Architecture in Certik Chain
CertiK Chain is envisioned as the Guardian of the Blockchain Galaxy, and it offers a variety of Combinators designed to address various aspects of security issues. 

 - Oracle Combinator: CertiK Chain has built-in frameworks that make it easier to fulfill general oracle workflows with decentralization and transparency. Oracle tasks and result aggregation calculations will be broadcast to CertiK Chain and recorded in states as proofs. The system is designed to reward good actors and punish bad ones by implementing a set of critical rules and reinforcements.
 - Security Primitive: Security Providers can register their on-chain services or off-chain API endpoints as Security Primitives, which Oracle Operators can then use to invoke them. Security Primitives are a collection of service functionalities that address security issues from a variety of perspectives. It is best practice to include a carefully selected list of Security Primitives in order to make the most informed decision possible about the security score of a given smart contract address and function signature.


#### Security Oracle Workflow in Certik Chain

 - End users submit oracle tasks for security insights they want on the Business Chain, which are funded with CTKs.
 - By subscribing to CertiK Chain events, Oracle Operators will receive the task.
 - It will send the task specifics to each Operator's customized Primitive Combination for real-time security checks;
 - The operator will respond to the oracle task by broadcasting a transaction to CertiK Chain after generating a security score.
 - When the task response window closes, CertiK Chain's Oracle Combinator collects all responses for that task and aggregates them into a final security score. Operators will be given task bounties based on their performance.
 - The final security score will then be pushed to the Business Chain's Security Oracle contract via a cross-chain bridge component.

### `x/oracle` module in Shentu 
#### Transactions and Queries 
#### **Transactions**
- `certik tx oracle create-operator <address> <collateral> [flags]`: Create an operator and deposit collateral. For instance, you want to create new  operator with alice's address from jack's key by using:
```{engine='sh'}
$ certik tx oracle create-operator certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we 500000uctk --fees 10000uctk --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.oracle.v1alpha1.MsgCreateOperator","address":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we","collateral":[{"denom":"uctk","amount":"500000"}],"proposer":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we","name":""}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[{"denom":"uctk","amount":"10000"}],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx oracle create-task <contract_address> <function> <bounty> [flags]`: Create a task. For example, we want to create a task with a contract has Txhash (=0xD...) and a function 0x00000000, we can use command below:
```{engine='sh'}
$ certik tx oracle create-task 0xD850942eF8811f2A866692A623011bDE52a462C1 0x00000000 1000uctk --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.oracle.v1alpha1.MsgCreateTask","contract":"0xD850942eF8811f2A866692A623011bDE52a462C1","function":"0x00000000","bounty":[{"denom":"uctk","amount":"1000"}],"description":"","creator":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we","wait":"0","valid_duration":"0s"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx oracle deposit-collateral <address> <amount> [flags]`: Increase an operator's collateral. After we assign Alice as operator, Jack can increase 10000uctk in Alice's collateral with this command:
```{engine='sh'}
$ certik tx oracle deposit-collateral certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we 10000uctk --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.oracle.v1alpha1.MsgAddCollateral","address":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we","collateral_increment":[{"denom":"uctk","amount":"10000"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx oracle withdraw-collateral <address> <amount> [flags]`: Reduce an operator's collateral. As mentioned above, Jack can reduce Alice's collateral to 10000uctk if he doesn't want to give Alice:
```{engine = 'sh'}
$ certik tx oracle withdraw-collateral certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we 10000uctk --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.oracle.v1alpha1.MsgReduceCollateral","address":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we","collateral_decrement":[{"denom":"uctk","amount":"10000"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx oracle respond-to-task <contract_address> <function> <score> [flags]`: Respond to a task. After creating a task, this command is used to response the task created and judge the score for the task: 
```{engine = 'sh'}
$ certik tx oracle respond-to-task 0xD850942eF8811f2A866692A623011bDE52a462C1 0x00000000 9 --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.oracle.v1alpha1.MsgTaskResponse","contract":"0xD850942eF8811f2A866692A623011bDE52a462C1","function":"0x00000000","score":"9","operator":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx oracle claim-reward <address> [flags]`: Withdraw all of an operator's accumulated rewards. As Alice is assigned to be operator, Jack wants to withdraw rewards of Alice by using this command: 
```{engine='sh'}
$ certik tx oracle claim-reward certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.oracle.v1alpha1.MsgWithdrawReward","address":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx oracle remove-operator <address> [flags]`: Remove an operator and withdraw collateral & rewards. For example, we remove the operator (Alice) that we created earlier above:
```{engine = 'sh'}
$ certik tx oracle remove-operator certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.oracle.v1alpha1.MsgRemoveOperator","address":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we","proposer":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```
- `certik tx oracle delete-task <contract_address> <function> [flags]`: Delete a finished task. After creating and responding the task, we can delete the task:
```{engine = 'sh'}
$ certik tx oracle delete-task 0xD850942eF8811f2A866692A623011bDE52a462C1 0x00000000 --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.oracle.v1alpha1.MsgDeleteTask","contract":"0xD850942eF8811f2A866692A623011bDE52a462C1","function":"0x00000000","force":false,"deleter":"certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

#### **Queries**
- `certik query oracle operator <address> [flags]`: Get operator information. (In this case, Alice is operator as transaction above).

```{engine = 'sh'}
$ certik query oracle operator certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we --chain-id certikchain

operator:
accumulated_rewards: []
address: certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we
collateral:
- amount: "500000"
denom: uctk
name: ""
proposer: certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we
```
- `certik query oracle operators [flags]`: Get operators information. Since only Alice is assigned as operator, only one operator listed below. We can assign many operators as we want.
```{engine = 'sh'}
$ certik query oracle operator certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we --chain-id certikchain

operator: #We just add 1 operator above.
accumulated_rewards: []
address: certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we
collateral:
- amount: "500000"
denom: uctk
name: ""
proposer: certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we
```
- `certik query oracle withdraws [flags]`: Get all withdrawals. As Jack withdraw 10000uctk from Alice, this query show the withdraw information here.
```{engine = 'sh'}
$ certik query oracle withdraws --chain-id certikchain

withdraws:
- address: certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we
amount:
- amount: "10000"
denom: uctk
due_block: "12271"
```

- `certik query oracle response <operator_address> <contract_address> <function> [flags]`: Get response information. After creating `tx oracle respond-to-task` above, we can get the information task's response here.
```{engine = 'sh'}
$ certik query oracle response certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we 0xD850942eF8811f2A866692A623011bDE52a462C1 0x00000000 --chain-id certikchain

response:
	operator: certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we
	reward:
	- amount: "1000"
	denom: uctk
	score: "9"
	weight: "500000"
```

- `certik query oracle task <contract_address> <function> [flags]`: Get task information. After creating a task, we can query to see the task information by this command:
```{engine = 'sh'}
$ certik query oracle task 0xD850942eF8811f2A866692A623011bDE52a462C1 0x00000000 --chain-id certikchain

task:
	begin_block: "23850"
	bounty:
	- amount: "1000"
	denom: uctk
	closing_block: "23870"
	contract: 0xD850942eF8811f2A866692A623011bDE52a462C1
	creator: certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we
	description: ""
	expiration: "2021-07-28T15:38:08.510323Z"
	function: "0x00000000"
	responses:
	- operator: certik1nc6v8tme0env488494ys09ld39dn9xzw6gc7we
	reward:
	- amount: "1000"
	denom: uctk
	score: "9"
	weight: "500000"
	result: "9"
	status: TASK_STATUS_SUCCEEDED
	waiting_blocks: "20"
```
