## `x/shield` module in Shentu
### Introduction
`shield` module is developed to utilize the unique staking, governance, and security features of CertiK Chain. CertiKShield is a decentralized pool of CTK that utilizes the CertiK Chain on-chain governance system to reimburse assets that have been lost, stolen, or are otherwise inaccessible on any blockchain network. Every member of a CertiKShield Pool can actively participate in selecting appropriate coverage circumstances thanks to this decentralized, on-chain voting process, leading in a dynamic and fully flexible coverage model. The cost of reserving funds from the CertiKShield Pool for personal reimbursement of lost assets will be directly related to the CertiK Security Oracle score, with lower scores (which indicate greater risk) demanding higher costs for protection.

There are two members of the CertiKShield system: Collateral Providers and Shield Purchasers:
 - Collateral Providers: Members who contribute their cryptocurrency (CTK or another accepted cryptocurrency) to the CertiKShield Pool as collateral. Because these collateralized funds are used to pay out any authorized reimbursement requests, these Collateral Providers may end up with less crypto than they started with. These members receive staking rewards for their staked CTK, as well as a share of the costs paid by Shield Purchasers seeking protection.
 - Shield Purchasers: Members who are seeking for a protection for their crypto assets. These Members must select the amount of protection they require for their crypto assets (referred to as a "Shield") and pay a fee that goes straight to the Collateral Providers who provided funds for reservation. The funds that were used as collateral for active Shield Purchasers' Shields can no longer be reserved until the Shields expire, allowing the Pool to preserve full collateralization.

#### Benefits for Blockchain Projects and Their Holders
CertiKShield Pools offer a flexible type protection for blockchain network supporters, safeguarding both the project and its community. CertiKShield allows members to create a discretionary fund that can be utilized for reimbursements if their supporters experience any unexpected issues. 

Any qualified member, whether from the blockchain project or an individual, must submit a Claim Proposal with a Submission Fee in order to request a reimbursement. This Submission Fee is designed to prevent illegitimate requests from being spammed. After the Submission Fee has been paid and the Claim Proposal has been carefully crafted, a decentralized voting process start, in which all CertiKShield system members can vote to accept or reject the Claim Proposal. 

The blockchain project is responsible for paying a recurring fee to the Members who deposit collateral into the CertiKShield Pool in order to keep it active (the Collateral Providers). Members are incentivized to contribute to the Pool in order to earn a part of these fees, which keeps incentives aligned. Blockchain projects directly reward Members who provide collateral to support their ecosystem's protection, and Members are incentivized to contribute to the Pool in order to earn a chunk of these fees.

#### Benefits and Risks for Staking Members / Liquidity Providers
All CertiKShield Pool Collateral Providers will get typical staking benefits for staked CTK, as well as a part of the costs paid by Pool Shield Purchasers.

CertiKShield Pools are a higher-risk, higher-reward staking option for CertiKShield. Staking on CertiK Nodes is another alternative that does not require any CTK as collateral. The main difference between staking and the CertiKShield Pool is that each Collateral Provider's stake can be utilized to reimburse Shield Purchasers for approved Claim Proposals.

All Collateral Providers should be aware that their own collateral stake could be used for reimbursements. As a result, Collateral Providers are responsible for completing thorough due diligence for all CertiKShield Pool Members; the Security Oracle score can be one factor of security, but all Members are encouraged to thoroughly examine all aspects of the blockchain project.

 ### `x/shield` module in Shentu
#### Transactions and Queries 
#### **Transactions**
- `certik tx shield clear-payouts [denom] [flags]`: Clear pending payouts after they have been distributed. Currently deprecated.

- `certik tx shield create-pool [shield amount] [sponsor] [sponsor-address] [flags]`: Create a Shield pool. Can only be executed from the Shield admin address. Every project wants to buy a Shield need to create a pool then the project can purchase Shield. For example, jack is the Shield admin and he can create a pool as described below with alice as sponsor:
```{engine = 'sh'}
$ certik tx shield create-pool 10000uctk abc $(certik keys show alice -a) --native-deposit 1000uctk --shield-limit 10000000 --from jack --fees 500uctk -y -b block --chain-id certikchain

{"height":"69","txhash":"8978C7CB0712943FB5081A14AAA6CB870810240BD3703956FC60564FE9274D21","codespace":"","code":0,"data":"0A0D0A0B6372656174655F706F6F6C","raw_log":"[{\"events\":[{\"type\":\"create_pool\",\"attributes\":[{\"key\":\"shield\",\"value\":\"10000uctk\"},{\"key\":\"deposit\",\"value\":\"native:\\u003cdenom:\\\"uctk\\\" amount:\\\"1000\\\" \\u003e \"},{\"key\":\"sponsor\",\"value\":\"abc\"},{\"key\":\"pool_id\",\"value\":\"1\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"create_pool\"},{\"key\":\"sender\",\"value\":\"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y\"},{\"key\":\"module\",\"value\":\"shield\"},{\"key\":\"sender\",\"value\":\"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"certik1qu4xymkj6mx86hdqlnpc3ucx0dyujw7dp9w7dr\"},{\"key\":\"sender\",\"value\":\"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y\"},{\"key\":\"amount\",\"value\":\"1000uctk\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"create_pool","attributes":[{"key":"shield","value":"10000uctk"},{"key":"deposit","value":"native:\u003cdenom:\"uctk\" amount:\"1000\" \u003e "},{"key":"sponsor","value":"abc"},{"key":"pool_id","value":"1"}]},{"type":"message","attributes":[{"key":"action","value":"create_pool"},{"key":"sender","value":"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y"},{"key":"module","value":"shield"},{"key":"sender","value":"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"certik1qu4xymkj6mx86hdqlnpc3ucx0dyujw7dp9w7dr"},{"key":"sender","value":"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y"},{"key":"amount","value":"1000uctk"}]}]}],"info":"","gas_wanted":"200000","gas_used":"115797","tx":null,"timestamp":""}
```

- `certik tx shield deposit-collateral [collateral] [flags]`: Join a Shield pool as a community member by depositing collateral. For example, Alice wants to join a Shield pool, she needs to deposit 100000uctk as below:
```{engine = 'sh'}
$ certik tx shield deposit-collateral 100000uctk --from alice --fees 500uctk -y -b block --chain-id certikchain

{"height":"57","txhash":"7E40EA6C5B5D04C87D4E51F08A038EC44542177BD6CAB8DB7ABABECCE0EE4036","codespace":"","code":0,"data":"0A140A126465706F7369745F636F6C6C61746572616C","raw_log":"[{\"events\":[{\"type\":\"deposit_collateral\",\"attributes\":[{\"key\":\"collateral\",\"value\":\"100000\"},{\"key\":\"sender\",\"value\":\"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y\"}]},{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"deposit_collateral\"},{\"key\":\"module\",\"value\":\"shield\"},{\"key\":\"sender\",\"value\":\"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"deposit_collateral","attributes":[{"key":"collateral","value":"100000"},{"key":"sender","value":"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y"}]},{"type":"message","attributes":[{"key":"action","value":"deposit_collateral"},{"key":"module","value":"shield"},{"key":"sender","value":"certik1ltsnt5y5jgzelffdw0cmdyfltyls493naaj
```

- `certik tx shield pause-pool [pool id] [flags]`: Pause a Shield pool to prevent new Shield purchases. Can only be executed from the Shield. admin address. After creating a pool above, we can pause the pool by this command: (flag --from jack indicate jack is the shield admin and he can create/pause/resume the pool).
```{engine = 'sh'}
$ certik tx shield pause-pool 1 --from jack -y --chain-id certikchain

{"height":"273","txhash":"BF329A55C34628B506D874507497775AEA3052E9CA95ED32C6E2D7E42A774455","codespace":"","code":0,"data":"0A0C0A0A70617573655F706F6F6C","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"pause_pool\"},{\"key\":\"module\",\"value\":\"shield\"},{\"key\":\"sender\",\"value\":\"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed\"}]},{\"type\":\"pause_pool\",\"attributes\":[{\"key\":\"pool_id\",\"value\":\"1\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"pause_pool"},{"key":"module","value":"shield"},{"key":"sender","value":"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed"}]},{"type":"pause_pool","attributes":[{"key":"pool_id","value":"1"}]}]}],"info":"","gas_wanted":"200000","gas_used":"44372","tx":null,"timestamp":""}
```

- `certik tx shield purchase [pool id] [shield amount] [description] [flags]`: Purchase Shield. Requires purchaser to provide descriptions of accounts to be protected. For example, Bob purchases a pool that has id 2 and he needs to provide description "test purchase 1" as below:
```{engine = 'sh'}
$ certik tx shield purchase 2 50000000uctk "test purchase 1" --from bob -y --chain-id certikchain

{"height":"808","txhash":"D9E07C3A566AEC8ED01B918C961D1BEB3A718FC50835AB1FA2475F6E3FF7CB34","codespace":"","code":0,"data":"0A110A0F70757263686173655F736869656C64","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"purchase_shield\"},{\"key\":\"sender\",\"value\":\"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed\"},{\"key\":\"module\",\"value\":\"shield\"},{\"key\":\"sender\",\"value\":\"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed\"}]},{\"type\":\"purchase_shield\",\"attributes\":[{\"key\":\"purchase_id\",\"value\":\"3\"},{\"key\":\"pool_id\",\"value\":\"2\"},{\"key\":\"protection_end_time\",\"value\":\"2021-08-18 17:23:28.41884 +0000 UTC\"},{\"key\":\"purchase_description\",\"value\":\"test purchase 1\"},{\"key\":\"shield\",\"value\":\"50000000\"},{\"key\":\"service_fees\",\"value\":\"native:\\u003cdenom:\\\"uctk\\\" amount:\\\"384500000000000000000000\\\" \\u003e \"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"certik1qu4xymkj6mx86hdqlnpc3ucx0dyujw7dp9w7dr\"},{\"key\":\"sender\",\"value\":\"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed\"},{\"key\":\"amount\",\"value\":\"384500uctk\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"purchase_shield"},{"key":"sender","value":"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed"},{"key":"module","value":"shield"},{"key":"sender","value":"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed"}]},{"type":"purchase_shield","attributes":[{"key":"purchase_id","value":"3"},{"key":"pool_id","value":"2"},{"key":"protection_end_time","value":"2021-08-18 17:23:28.41884 +0000 UTC"},{"key":"purchase_description","value":"test purchase 1"},{"key":"shield","value":"50000000"},{"key":"service_fees","value":"native:\u003cdenom:\"uctk\" amount:\"384500000000000000000000\" \u003e "}]},{"type":"transfer","attributes":[{"key":"recipient","value":"certik1qu4xymkj6mx86hdqlnpc3ucx0dyujw7dp9w7dr"},{"key":"sender","value":"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed"},{"key":"amount","value":"384500uctk"}]}]}],"info":"","gas_wanted":"200000","gas_used":"96859","tx":null,"timestamp":""}
```

- `certik tx shield resume-pool [pool id] [flags]`: Resume a Shield pool to reactivate Shield purchase. Can only be executed from the Shield admin address. After pausing pool 1 above, jack can resume the pool so that members in Shield pool can purchase it.
```{engine = 'sh'}
$ certik tx shield resume-pool 1 --from jack -y --chain-id certikchain

{"height":"282","txhash":"24C7755605CBCE48489E312B619C21418F27B6E994F9B9F67D0110D04F758FDB","codespace":"","code":0,"data":"0A0D0A0B726573756D655F706F6F6C","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"resume_pool\"},{\"key\":\"module\",\"value\":\"shield\"},{\"key\":\"sender\",\"value\":\"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed\"}]},{\"type\":\"resume_pool\",\"attributes\":[{\"key\":\"pool_id\",\"value\":\"1\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"resume_pool"},{"key":"module","value":"shield"},{"key":"sender","value":"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed"}]},{"type":"resume_pool","attributes":[{"key":"pool_id","value":"1"}]}]}],"info":"","gas_wanted":"200000","gas_used":"44436","tx":null,"timestamp":""}
```

- `certik tx shield stake-for-shield [pool id] [shield amount] [description] [flags]`: Obtain shield through staking. Requires purchaser to provide descriptions of accounts to be protected. For example, jack wants to obtain the shield, he can stake 5000000uctk for pool with id 1 and provide the description "test 1" as described below:
```{engine = 'sh'}
$ certik tx shield stake-for-shield 1 50000000uctk "test 1" --chain-id certikchain --from jack

{"body":{"messages":[{"@type":"/shentu.shield.v1alpha1.MsgStakeForShield","pool_id":"1","shield":[{"denom":"uctk","amount":"50000000"}],"description":"test 1","from":"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}
  
confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx shield unstake-from-shield [pool id] [amount]  [flags]`: Withdraw staking from shield. Requires existing shield purchase through staking. After staking for shield above, jack can unstake  the amount with this command:
```{engine = 'sh'}
$ certik tx shield unstake-from-shield 1 5000000uctk  --chain-id certikchain --from jack

{"body":{"messages":[{"@type":"/shentu.shield.v1alpha1.MsgUnstakeFromShield","pool_id":"1","shield":[{"denom":"uctk","amount":"5000000"}],"from":"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}  

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx shield update-pool [pool id] [flags]`: Update a Shield pool. Can only be executed from the Shield admin address. We can update an existing Shield pool by adding more deposit or updating Shield amount:
```{engine = 'sh'}
$ certik tx shield update-pool 1 --shield-limit 100000000 --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.shield.v1alpha1.MsgUpdatePool","from":"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed","shield":[],"service_fees":{"native":[],"foreign":[]},"pool_id":"1","description":"","shield_limit":"100000000"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx shield update-sponsor [pool id] [new_sponsor] [new_sponsor_address] [flags]`: Update a pool's sponsor. Can only be executed from the Shield admin address. As alice is the sponsor for pool 1, we can update bob as sponsor instead of alice with this command: 
```{engine = 'sh'}
$ certik tx shield update-sponsor 1 testb $(certik keys show bob -a) --from bob --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.shield.v1alpha1.MsgUpdateSponsor","pool_id":"1","sponsor":"testb","sponsor_addr":"certik1wymn4cm35qmqamgtpaljsdjah620dc5dj9560l","from":"certik1wymn4cm35qmqamgtpaljsdjah620dc5dj9560l"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx shield withdraw-collateral [collateral] [flags]`: Withdraw deposited collateral from Shield pool. Alice can withdraw her collateral after she deposited it by using: 
```{engine = 'sh'}
$ certik tx shield withdraw-collateral 100000uctk --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.shield.v1alpha1.MsgWithdrawCollateral","from":"certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed","collateral":[{"denom":"uctk","amount":"100000"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx shield withdraw-foreign-rewards [denom] [address] [flags]`: Withdraw foreign rewards coins to their original chain. Currently deprecated.

- `certik tx shield withdraw-reimbursement [proposal id] [flags]`: Withdraw reimbursement by proposal id. A proposal is created when someone purchase a pool and they want to submit a proposal through `gov` module. After the voting period and proposal status is passed, the reimbursement is directly   approved. After the reimbursement is made, we can withdraw it by using:
```{engine = 'sh'}
$ certik tx shield withdraw-reimbursement 1 --from bob --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.shield.v1alpha1.MsgWithdrawReimbursement","proposal_id":"3","from":"certik17k6gljthxvuexhfmqja43rxete2h6ua9ahm3s4"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx shield withdraw-rewards [flags]`: Withdraw CTK rewards.
```{engine = 'sh'}
$ certik tx shield withdraw-rewards --from jack --chain-id certikchain

{"body":{"messages":[{"@type":"/shentu.shield.v1alpha1.MsgWithdrawRewards","from":"certik1a58devkfvteuqrre4aurj8h3hz9dlsmn4ynwkf"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```
#### **Queries**
- `certik query shield claim-params [flags]`: Get claim parameters.
```{engine = 'sh'}
$ certik query shield claim-params --chain-id certikchain

params:
	claim_period: 1814400s
	deposit_rate: "0.100000000000000000"
	fees_rate: "0.010000000000000000"
	min_deposit:
	- amount: "100000000"
	denom: uctk
	payout_period: 4838400s	
```

- `certik query shield pool [pool_ID] [flags]`: Query a pool. For example, after the pool 1 is created, we can query to get the information:
```{engine = 'sh'}
$ certik query shield pool 1 --chain-id certikchain
	pool:
	active: true
	description: ""
	id: "1"
	shield: "10000"
	shield_limit: "10000000"
	sponsor: abc
	sponsor_addr: certik1ltsnt5y5jgzelffdw0cmdyfltyls493naajx2y
```

- `certik query shield pool-params [flags]`: Get pool parameters. 
```{engine = 'sh'}
$ certik query shield pool-params --chain-id certikchain

params:
	min_shield_purchase:
	- amount: "50000000"
	denom: uctk
	pool_shield_limit: "0.500000000000000000"
	protection_period: 1814400s
	shield_fees_rate: "0.007690000000000000"
	withdraw_period: 1814400s
```

- `certik query shield pool-purchaser [pool_ID] [purchaser_address] [flags]`: Get purchases corresponding to a given pool-purchaser pair. For example, jack is the purchaser for pool 2, we can query to get information from jack and transaction of pool 2:
```{engine = 'sh'}
$ certik query shield pool-purchaser 2 $(certik keys show jack -a) --chain-id certikchain

purchase_list:
	entries:
	- deletion_time: "2021-08-18T16:35:53.986569Z"
	description: shield for sponsor
	protection_end_time: "2021-08-18T16:35:53.986569Z"
	purchase_id: "2"
	service_fees:
	foreign: []
	native:
	- amount: "10000000.000000000000000000"
	denom: uctk
	shield: "10000"
	- deletion_time: "2021-08-18T17:23:28.418840Z"
	description: test purchase 1
	protection_end_time: "2021-08-18T17:23:28.418840Z"
	purchase_id: "3"
	service_fees:
	foreign: []
	native:
	- amount: "384500.000000000000000000"
	denom: uctk
	shield: "50000000"
	pool_id: "2"
	purchaser: certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed
```

- `certik query shield pool-purchases [pool_ID] [flags]`: Query purchases in a given pool. Get all the purchasers from a pool 1 by using: 
```{engine = 'sh'}
$ certik query shield pool-purchases 1 --chain-id certikchain

purchase_lists:
	- entries:
	- deletion_time: "2021-08-18T16:20:18.048781Z"
	description: shield for sponsor
	protection_end_time: "2021-08-18T16:20:18.048781Z"
	purchase_id: "1"
	service_fees:
	foreign: []
	native:
	- amount: "1000000.000000000000000000"
	denom: uctk
	shield: "1000000"
	pool_id: "1"
	purchaser: certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed
```

- `certik query shield pools [flags]`: Query a complete list of pools. After pools are created, we can get all the information of the pools:
```{engine = 'sh'}
$ certik query shield pools --chain-id certikchain

pools:

	- active: true
	description: ""
	id: "1"
	shield: "1000000"
	shield_limit: "100000000"
	sponsor: abc
	sponsor_addr: certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed

	- active: true
	description: ""
	id: "2"
	shield: "10000"
	shield_limit: "100000000"
	sponsor: test
	sponsor_addr: certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed
```

- `certik query shield provider [provider_address] [flags]`: Get provider information.
```{engine = 'sh'}
$ certik query shield provider certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed --chain-id certikchain
providers:
	- address: certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed
	collateral: "100000000"
	delegation_bonded: "10020000000"
	rewards:
	foreign: []
	native:
	- amount: "648091.781665013227513228"
	denom: uctk
	total_locked: "0"
	withdrawing: "100000"
```

- `certik query shield providers [flags]`: Query providers with pagination parameters.
```{engine = 'sh'}
$ certik query shield providers --chain-id certikchain

providers:
	- address: certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed
	collateral: "100000000"
	delegation_bonded: "10020000000"
	rewards:
	foreign: []
	native:
	- amount: "648091.781665013227513228"
	denom: uctk
	total_locked: "0"
	withdrawing: "100000"
```

- `certik query shield purchases [flags]`: Query all purchases. For example, each pool is purchased by a purchaser, we can query all the information of the purchase lists:
```{engine = 'sh'}
$ purchases:
	- deletion_time: "2021-08-18T16:20:18.048781Z"
	description: shield for sponsor
	protection_end_time: "2021-08-18T16:20:18.048781Z"
	purchase_id: "1"
	service_fees:
	foreign: []
	native:
	- amount: "1000000.000000000000000000"
	denom: uctk
	shield: "1000000"
	- deletion_time: "2021-08-18T16:35:53.986569Z"
	description: shield for sponsor
	protection_end_time: "2021-08-18T16:35:53.986569Z"
	purchase_id: "2"
	service_fees:
	foreign: []
	native:
	- amount: "10000000.000000000000000000"
	denom: uctk
	shield: "10000"
	- deletion_time: "2021-08-18T17:23:28.418840Z"
	description: test purchase 1
	protection_end_time: "2021-08-18T17:23:28.418840Z"
	purchase_id: "3"
	service_fees:
	foreign: []
	native:
	- amount: "384500.000000000000000000"
	denom: uctk
	shield: "50000000"
```

- `certik query shield purchases-by [purchaser_address] [flags]`: Query purchase information of a given account. For example, we want to get information from purchaser jack, we can use this command:
```{engine = 'sh'}
certik query shield purchases-by $(certik keys show jack -a) --chain-id certikchain

purchase_lists:
	- entries:
	- deletion_time: "2021-08-18T16:20:18.048781Z"
	description: shield for sponsor
	protection_end_time: "2021-08-18T16:20:18.048781Z"
	purchase_id: "1"
	service_fees:
	foreign: []
	native:
	- amount: "1000000.000000000000000000"
	denom: uctk
	shield: "1000000"
	pool_id: "1"
	purchaser: certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed
	- entries:
	- deletion_time: "2021-08-18T16:35:53.986569Z"
	description: shield for sponsor
	protection_end_time: "2021-08-18T16:35:53.986569Z"
	purchase_id: "2"
	service_fees:
	foreign: []
	native:
	- amount: "10000000.000000000000000000"
	denom: uctk
	shield: "10000"
	- deletion_time: "2021-08-18T17:23:28.418840Z"
	description: test purchase 1
	protection_end_time: "2021-08-18T17:23:28.418840Z"
	purchase_id: "3"
	service_fees:
	foreign: []
	native:
	- amount: "384500.000000000000000000"
	denom: uctk
	shield: "50000000"
	pool_id: "2"
	purchaser: certik1p49jqcz6zs7n603qr74p5tmzf3068uee9dc5ed
```

- `certik query shield reimbursement [proposal ID] [flags]`: Query a reimbursement. After the proposal is passed, a reimbursement is made and we can get the information of that reimbursement:
```{engine = 'sh'}
$ certik query shield reimbursement 1

reimbursement:
	amount:
	- amount: "1000000"
	denom: uctk
	beneficiary: certik1vz43j346g7mkcnrd46hz02vvm7glkxrdx5fj8u
	payout_time: "2021-09-27T14:11:25.195151Z"
```

- `certik query shield reimbursements [flags]`: Query all reimbursements. 
```{engine = 'sh'}
certik query shield reimbursements

pairs:
	- proposal_id: "1"
	reimbursement:
	amount:
	- amount: "1000000"
	denom: uctk
	beneficiary: certik1vz43j346g7mkcnrd46hz02vvm7glkxrdx5fj8u
	payout_time: "2021-09-27T14:11:25.195151Z"
```

- `certik query shield shield-staking-rate [flags]`: Get shield staking rate for stake-for-shield.
```{engine = 'sh'}
$ certik query shield shield-staking-rate --chain-id certikchain

rate: "2.000000000000000000"
```

- `certik query shield sponsor [pool_ID] [flags]`: Query pools for a sponsor. #the [pool_ID] should be change to [sponsor_address]. Each pool is created with a sponsor and we can get information of that sponsor: 
```{engine = 'sh'}
$ certik query shield sponsor $(certik keys show alice -a)

pools:
	- active: false
	description: ""
	id: "1"
	shield: "50000000"
	shield_limit: "100000000"
	sponsor: abc
	sponsor_addr: certik199tr500yufs4r9xyzgwjr5saql6942gnewqn2h
	- active: true
	description: ""
	id: "2"
	shield: "50990000"
	shield_limit: "100000000"
	sponsor: alice
	sponsor_addr: certik199tr500yufs4r9xyzgwjr5saql6942gnewqn2h
```

- `certik query shield staked-for-shield [pool_ID] [purchaser_address] [flags]`: Get staked CTK for shield corresponding to a given pool-purchaser pair.
```{engine = 'sh'}
$ certik query shield staked-for-shield 1 $(certik keys show bob -a)

shield_staking:
	amount: "100000000"
	pool_id: "1"
	purchaser: certik1a7zukuunk727yflrl09cdxlpzdr5q5zk9skanw
	withdraw_requested: "0"
```

- `certik query shield status [flags]`: Get shield status.
```{engine = 'sh'}
$ certik query shield status --chain-id certikchain

current_service_fees:
foreign: []
native: []
global_shield_staking_pool: "0"
remaining_service_fees:
foreign: []
native: []
total_collateral: "1000"
total_shield: "0"
total_withdrawing: "0"
```
