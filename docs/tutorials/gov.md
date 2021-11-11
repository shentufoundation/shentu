## `x/gov` module

### Introduction
The Cosmos-SDK-based blockchain can now support an on-chain governance system with to this module. Holders of the chain's native staking token can vote on proposals on a 1 token, 1 vote basis in this system. The `gov` module enables on-chain governance which allows Certik Chain token holder to participate in the decision-making processes. For example, users can:

 - Form a proposal and seek feedback.
 - Create the proposal and make any necessary changes based on feedback.
 - Submit a proposal along with a deposit.
 - Tokens must be deposited in order to fund an active proposal.
 - Vote for a proposal that is currently active.

The `gov` module currently supports these features:

 - Proposal submission: Users can submit proposals with a deposit. The proposal enters the voting session once the minimal deposit is met.
 - Voting: Participants can vote on proposals that have passed the MinDeposit requirement.
 - Inheritance and penalties: If a delegate does not vote, their validator's vote is inherited.
 - Claiming deposit: Users who made deposits on proposals can get their money back if the proposal was accepted or if the idea was never presented up for a vote.

#### The Governance Procedure in Shentu

 1. **Proposals**
	 Validators and certifiers regulate the chain by means of proposals. A proposal can be submitted by any user, but some types of proposals must be submitted by a certifier. Before a proposal can be voted on, it must first meet a minimum deposit threshold; this helps eliminate spam. The deposit is usually returned at the end of the proposal's life cycle, regardless of whether it passes or not. This deposit time is waived for proposals submitted by a validator or certifier. All three types of participants can submit a proposal: Stake Delegators, Validator Operators, and Security Certifiers.
	 
 2. **Deposit Period** 
	 In this period, the proposal has the status "PROPOSAL_STATUS_DEPOSIT_PERIOD" which indicates a proposal status during the deposit. 

	A deposit is required for Stake Delegates to submit a proposal. Once the minimum deposit is paid, the proposal will be confirmed and the voting period will begin. Depending on the type of proposal, there will be one or two voting rounds, depending on the security changes that the plan will come up. For instance, software upgrade proposals have two different voting stages and they must pass the scrutiny of both certifiers and validators while plain text proposals do not automatically trigger chain action when they pass, they just require a single validator voting period to pass.

	Stake Delegates can gain greater attention through the deposit process, which helps them get their proposals into the Voting Period. Because they have already been entrusted by delegators during the delegation and election process, Security Certifiers and Validator Operators can skip the deposit period when submitting proposals.

 3. **Voting Period**: 
	 	There are four options for voting:
	 - Yes: want the proposal to pass
	 - No: do not want the proposal to pass and want to return the deposit
	 - No with Veto: do not want the proposal to pass and opt to burn the deposit
	 - Abstain: not to participate in the vote

	Moreover, two security voting options in this procedure:
	 - Yes: want the proposal passed and there is no security issue.
	 - No / Abstain:  do not want the proposal passed or elect not to participate in the vote, and there are potential security issues.
	 
	 During the voting period, there are two passes of voting (thus the name of the dual-pass governance model) for functional and security considerations.

	In the voting period, the proposal has 2 status:

	 - `PROPOSAL_STATUS_CERTIFIER_VOTING_PERIOD`: a certifier voting period status. Certifiers can vote with two options: yes/no. Each vote is equally weighted; each certifier receives one vote. The percentage of `yes` votes out of total votes must meet a certain level. This status will be passed after the certifier votes and move to the next status;
	 
	 - `PROPOSAL_STATUS_VALIDATOR_VOTING_PERIOD`: a validator voting period status. There are four options for validators to vote: yes/no/abstain/no_with_veto. The total amount of tokens staked on the chain is used to weight each validator's vote. The proposal fails if a sufficient number of `no_with_veto` votes are cast, and the deposit is destroyed. The deposit is refunded in all other scenarios. The percentage of `yes` votes out of all non-abstinent ballots must be greater than a threshold for a proposal to pass. The proposal will fail if it does not meet the threshold.
	 
	 Only staked tokens are allowed to vote in the voting pass for practical reasons. The voting power, or influence on the decision, is determined by the quantity of tokens staked. Unless they decide to cast their own vote, which would override the Validators' voting choice, Stake Delegators adopt the vote of the Validators they have decided to delegate. On each proposal, each token is entitled to one vote.
	 
	 For the avoidance of doubt, the right to vote is limited to voting on CertiK Platform features; the right to vote does not entitle CTK holders to vote on the Foundation, its affiliates, or their assets, and it does not imply any equity interest in any of these entities.

 5. **Results**
	 There are 3 status for the voting result:

	 - `PROPOSAL_STATUS_PASSED`: a proposal status of a proposal that has passed. The number of voters vote "yes" dominates.
	 - `PROPOSAL_STATUS_REJECTED`: a proposal status of a proposal that has been rejected. The number of voters vote "no" dominates.
	 - `PROPOSAL_STATUS_FAILED`: a proposal status of a proposal that has failed.

	 If the proposal was approved as a software upgrade, nodes must update their software to the new version approved. This procedure is broken down into two parts: a signal and a switch.

	Validator Operators are expected to download and install the updated version of the software while continuing to use the old version during the signal step. When a validator node is deployed as part of the upgrade, it will begin broadcasting to the network that it is ready to switch over.

	At every one time, there is just one signal. If several software upgrade proposals are accepted in a short period of time, a pipeline will form, with each request being implemented in the sequence in which it was accepted.

	When a majority of validator nodes signal a common software upgrade proposal, all nodes (validator nodes, non-validator nodes, and light nodes) are expected to move to the new software version at the same time in the switch phase.
	

### `x/gov` module in Shentu 
#### Transactions and Queries 
#### **Transactions**
- `certik tx gov submit-proposal shield-claim [proposal file] [flags]`: Submit a Shield claim proposal along with an initial deposit. The proposal needs to go to 2 phases voting periods: certifier and validator. The proposal details must be supplied via a JSON file. Here is a sample of a proposal for shield claim:
```
{
	"pool_id": 1,
	"loss": [
		{
			"denom": "uctk",
			"amount": "1000000"
		}
	],
	"evidence": "Attack happened on <time> caused loss of <amount> to <account> by <txhash>",
	"purchase_id": 2,
	"description": "Details of the attack",
	"deposit": [
		{
			"denom": "uctk",
			"amount": "100000000"
		}
	]
}
```

Here is the transaction for the Shield claim proposal: 
```{engine='sh'}
$ certik tx gov submit-proposal shield-claim SampleShieldProposal.json --from bob --fees 5000uctk -y -b block --chain-id certikchain --gas=auto

gas estimate: 287069

{"height":"28292","txhash":"36F2D547C998D039386DE7A4799FDAB67F924C90211D3529AC69C7C9C9E46506","codespace":"","code":0,"data":"0A150A0F7375626D69745F70726F706F73616C12020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"submit_proposal\"},{\"key\":\"sender\",\"value\":\"certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03\"},{\"key\":\"module\",\"value\":\"governance\"},{\"key\":\"sender\",\"value\":\"certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03\"}]},{\"type\":\"proposal_deposit\",\"attributes\":[{\"key\":\"amount\",\"value\":\"100000000uctk\"},{\"key\":\"proposal_id\",\"value\":\"1\"},{\"key\":\"depositor\",\"value\":\"certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03\"}]},{\"type\":\"submit_proposal\",\"attributes\":[{\"key\":\"proposal_type\",\"value\":\"ShieldClaim\"},{\"key\":\"proposal_id\",\"value\":\"1\"},{\"key\":\"voting_period_start\",\"value\":\"1\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"certik10d07y265gmmuvt4z0w9aw880jnsr700ja20jhc\"},{\"key\":\"sender\",\"value\":\"certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03\"},{\"key\":\"amount\",\"value\":\"100000000uctk\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"submit_proposal"},{"key":"sender","value":"certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03"},{"key":"module","value":"governance"},{"key":"sender","value":"certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03"}]},{"type":"proposal_deposit","attributes":[{"key":"amount","value":"100000000uctk"},{"key":"proposal_id","value":"1"},{"key":"depositor","value":"certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03"}]},{"type":"submit_proposal","attributes":[{"key":"proposal_type","value":"ShieldClaim"},{"key":"proposal_id","value":"1"},{"key":"voting_period_start","value":"1"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"certik10d07y265gmmuvt4z0w9aw880jnsr700ja20jhc"},{"key":"sender","value":"certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03"},{"key":"amount","value":"100000000uctk"}]}]}],"info":"","gas_wanted":"287069","gas_used":"285542","tx":null,"timestamp":""}
```

- `certik tx gov submit-proposal certifier-update [proposal file] [flags]`: Submit a certifier update proposal along with an initial deposit. Certifier Update Proposals add a new certifier or remove a certifier if they pass. This proposal type must be submitted by a certifier. Its voting protocol is unique, in that it must pass either the certifier voting round or the validator voting round. The proposal details must be supplied via a JSON file. Here is a sample of a proposal for certifier update:
```
{
	"title": "New Certifier, Joe Shmoe",
	"description": "Why we should make Joe Shmoe a certifier",
	"certifier": "certik1fdyv6hpukqj6kqdtwc42qacq9lpxm0pn85w6l9",
	"add_or_remove": "add",
	"alias": "joe",
	"deposit": [
		{
		"denom": "ctk",
		"amount": "100"
		}
	]
}
```

Here is the transaction for the Certifier update proposal:
``` {engine = 'sh'}
#This command submits proposal for certifier update. Alice is proposed to be a certifier by Jack.
$ certik tx gov submit-proposal certifier-update certifier-update.json --from jack -y --chain-id certikchain

{"height":"29469","txhash":"876B18777488BA25BCA47A1CF957FE183F82E3CFF1359F3D1178E1B9001FFDE8","codespace":"","code":0,"data":"0A150A0F7375626D69745F70726F706F73616C12020802","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"submit_proposal\"},{\"key\":\"module\",\"value\":\"governance\"},{\"key\":\"sender\",\"value\":\"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p\"}]},{\"type\":\"submit_proposal\",\"attributes\":[{\"key\":\"proposal_type\",\"value\":\"CertifierUpdate\"},{\"key\":\"proposal_id\",\"value\":\"2\"},{\"key\":\"voting_period_start\",\"value\":\"2\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"submit_proposal"},{"key":"module","value":"governance"},{"key":"sender","value":"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p"}]},{"type":"submit_proposal","attributes":[{"key":"proposal_type","value":"CertifierUpdate"},{"key":"proposal_id","value":"2"},{"key":"voting_period_start","value":"2"}]}]}],"info":"","gas_wanted":"200000","gas_used":"97397","tx":null,"timestamp":""}
```

- `certik tx gov submit-proposal community-pool-spend [proposal file] [flags]`: Submit a community pool spend proposal along with an initial deposit. Community Pool Spend Proposals transfer tokens from the community pool to an address if they pass. The recipient of the tokens could be a chain user who has done—or plans to do—development work or security work for the chain. The recipient could also be a smart contract that distributes the tokens to multiple addresses after conditions have been met. A community pool spend proposal only needs to pass the validator voting period. The proposal details must be supplied via a JSON file. Here is a sample of a proposal for community pool spend:
```
{
	"title": "Community Pool Spend",
	"description": "Pay me some Atoms!",
	"recipient": "certik1s5afhd6gxevu37mkqcvvsj8qeylhn0rz46zdlq",
	"amount": "1000stake",
	"deposit": "1000stake"
}
```

Here is the transaction for the community pool spend proposal:
``` {engine = 'sh'}
$ certik tx gov submit-proposal community-pool-spend community-pool.json --from jack -y --chain-id certikchain

{"height":"29555","txhash":"CF4D7650B8CE3C87361ACE936EA4896BEE2B03100FE0B862C03BB65D964FB232","codespace":"","code":0,"data":"0A150A0F7375626D69745F70726F706F73616C12020803","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"submit_proposal\"},{\"key\":\"module\",\"value\":\"governance\"},{\"key\":\"sender\",\"value\":\"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p\"}]},{\"type\":\"submit_proposal\",\"attributes\":[{\"key\":\"proposal_type\",\"value\":\"CommunityPoolSpend\"},{\"key\":\"proposal_id\",\"value\":\"3\"},{\"key\":\"voting_period_start\",\"value\":\"3\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"submit_proposal"},{"key":"module","value":"governance"},{"key":"sender","value":"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p"}]},{"type":"submit_proposal","attributes":[{"key":"proposal_type","value":"CommunityPoolSpend"},{"key":"proposal_id","value":"3"},{"key":"voting_period_start","value":"3"}]}]}],"info":"","gas_wanted":"200000","gas_used":"83203","tx":null,"timestamp":""}
```

- `certik tx gov submit-proposal param-change [proposal file] [flags]`: Submit a parameter proposal along with an initial deposit. The proposal specifies the new value for the parameter as well as the parameter to be changed. A parameter change proposal, like a plain text proposal, merely goes through the validator voting phase. The proposal details must be supplied via a JSON file. For values that contains objects, only non-empty fields will be updated. Here is a sample of a proposal for parameter change:
```
# Currently parameter changes are evaluated but not validated, so it is very important that any "value" change is valid.
{
	"title": "Staking Param Change",
	"description": "Update max validators",
	"changes": [
	{
		"subspace": "staking",
		"key": "MaxValidators", #the param wants to be changed
		"value": 105
	}
	],
	"deposit": "1000stake"
}
```

Here is the transaction for the Shield claim proposal. For example, we want to change MaxValidators to 105: 
```{engine = 'sh'}
$ certik tx gov submit-proposal param-change param-change.json --from jack -y --chain-id certikchain

{"height":"29575","txhash":"1A079E06261B9C920F604D1AACDBB72361418C8D8DFAA50B8D63FC63B1911B87","codespace":"","code":0,"data":"0A150A0F7375626D69745F70726F706F73616C12020804","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"submit_proposal\"},{\"key\":\"module\",\"value\":\"governance\"},{\"key\":\"sender\",\"value\":\"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p\"}]},{\"type\":\"submit_proposal\",\"attributes\":[{\"key\":\"proposal_type\",\"value\":\"ParameterChange\"},{\"key\":\"proposal_id\",\"value\":\"4\"},{\"key\":\"voting_period_start\",\"value\":\"4\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"submit_proposal"},{"key":"module","value":"governance"},{"key":"sender","value":"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p"}]},{"type":"submit_proposal","attributes":[{"key":"proposal_type","value":"ParameterChange"},{"key":"proposal_id","value":"4"},{"key":"voting_period_start","value":"4"}]}]}],"info":"","gas_wanted":"200000","gas_used":"81036","tx":null,"timestamp":""}
```

- `certik tx gov submit-proposal software-upgrade [name] (--upgrade-height [height] | --upgrade-time [time]) (--upgrade-info [info]) [flags]`: Submit a software upgrade along with an initial deposit. Please specify a unique name and height OR time for the upgrade to take effect. Software Upgrade Proposals require chain code modifications, such as changing the range/scope of the chain parameters or adding new chain features. If a software upgrade proposal passes, it will introduce a change to the chain's code that all validators must adopt. A user will usually submit a plain text proposal before submitting a software upgrade proposal to assess interest in their possible code improvement. A software update proposal must pass both a certifier and a validator voting round because software upgrades can bring security risks.
```
#Upgrade is the name of the software update, update at height=10000000000000000 (height must not be in the past) along with flags.

$ certik tx gov submit-proposal software-upgrade Upgrade --from jack -y --chain-id certikchain --upgrade-height 10000000000000000 --title "Upgrade Test" --description "Upgrade to latest version" --fees 500uctk

{"height":"29881","txhash":"8479FAA52696F99DCCB9075AC79A9979A0F3DDD2CA347748DF556FC32E56B403","codespace":"","code":0,"data":"0A150A0F7375626D69745F70726F706F73616C12020806","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"submit_proposal\"},{\"key\":\"module\",\"value\":\"governance\"},{\"key\":\"sender\",\"value\":\"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p\"}]},{\"type\":\"submit_proposal\",\"attributes\":[{\"key\":\"proposal_type\",\"value\":\"SoftwareUpgrade\"},{\"key\":\"proposal_id\",\"value\":\"6\"},{\"key\":\"voting_period_start\",\"value\":\"6\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"submit_proposal"},{"key":"module","value":"governance"},{"key":"sender","value":"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p"}]},{"type":"submit_proposal","attributes":[{"key":"proposal_type","value":"SoftwareUpgrade"},{"key":"proposal_id","value":"6"},{"key":"voting_period_start","value":"6"}]}]}],"info":"","gas_wanted":"200000","gas_used":"92107","tx":null,"timestamp":""}
```
After we have several proposals created, we can move to deposit phase if some proposals are not activated due to the minimum deposit earlier.

- `certik tx gov submit-proposal cancel-software-upgrade [flags]`: Cancel a software upgrade along with an initial deposit. The command requires flags --description and --title:

```{engine ='sh'}
certik tx gov submit-proposal cancel-software-upgrade --fees 50uctk --description "Cancel Software" --title "Test Cancel" --from jack -y --chain-id certikchain

{"height":"29660","txhash":"F9F16C140C507119C4AA0DFA0F5337C81DB7BFB1988A77F596178FF0172A9358","codespace":"","code":0,"data":"0A150A0F7375626D69745F70726F706F73616C12020805","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"submit_proposal\"},{\"key\":\"module\",\"value\":\"governance\"},{\"key\":\"sender\",\"value\":\"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p\"}]},{\"type\":\"submit_proposal\",\"attributes\":[{\"key\":\"proposal_type\",\"value\":\"CancelSoftwareUpgrade\"},{\"key\":\"proposal_id\",\"value\":\"5\"},{\"key\":\"voting_period_start\",\"value\":\"5\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"submit_proposal"},{"key":"module","value":"governance"},{"key":"sender","value":"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p"}]},{"type":"submit_proposal","attributes":[{"key":"proposal_type","value":"CancelSoftwareUpgrade"},{"key":"proposal_id","value":"5"},{"key":"voting_period_start","value":"5"}]}]}],"info":"","gas_wanted":"200000","gas_used":"85477","tx":null,"timestamp":""}
```

- `certik tx gov deposit [proposal-id] [deposit] [flags]`: Submit a deposit for an active proposal. You can find the proposal-id by running "`certik query gov proposals`". If the proposal doesn't reach the minimum deposit, use will use this transaction to deposit more into the proposal to make it active: 
```{engine='sh'}
$ certik tx gov deposit 1 1000uctk --from bob --chain-id certikchain

{"body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgDeposit","proposal_id":"1","depositor":"certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03","amount":[{"denom":"uctk","amount":"100"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
```

- `certik tx gov vote [proposal-id] [option] [flags]`: Submit a vote for an active proposal. You can find the proposal-id by running "`certik query gov proposals`". Options here are: yes/no/no_with_veto/abstain/unspecified. The voters must be certified identity in order to vote. (using `certik tx cert issue-certificate identity`).
- 
```{engine='sh'}
$ certik tx gov vote 1 yes --from jack --fees 5000uctk -y --chain-id certikchain

{"height":"28618","txhash":"8236D3BF5BCA568E5636D677DCB912143EE3F34954CC2EC71F2E49D5DAB6E00A","codespace":"","code":0,"data":"0A060A04766F7465","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"vote\"},{\"key\":\"module\",\"value\":\"governance\"},{\"key\":\"sender\",\"value\":\"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p\"}]},{\"type\":\"proposal_vote\",\"attributes\":[{\"key\":\"option\",\"value\":\"VOTE_OPTION_YES\"},{\"key\":\"proposal_id\",\"value\":\"1\"},{\"key\":\"voter\",\"value\":\"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p\"},{\"key\":\"txhash\",\"value\":\"8236d3bf5bca568e5636d677dcb912143ee3f34954cc2ec71f2e49d5dab6e00a\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"vote"},{"key":"module","value":"governance"},{"key":"sender","value":"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p"}]},{"type":"proposal_vote","attributes":[{"key":"option","value":"VOTE_OPTION_YES"},{"key":"proposal_id","value":"1"},{"key":"voter","value":"certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p"},{"key":"txhash","value":"8236d3bf5bca568e5636d677dcb912143ee3f34954cc2ec71f2e49d5dab6e00a"}]}]}],"info":"","gas_wanted":"200000","gas_used":"57755","tx":null,"timestamp":""}
```

#### **Queries**
- `certik query gov deposit [proposal-id] [depositer-addr] [flags]`: Query details for a single proposal deposit on a proposal by its identifier. Currently deprecated. 

- `certik query gov deposits [proposal-id] [flags]`: Query details for all deposits on a proposal. 
```{engine = 'sh'}
$ certik query gov deposits 1

deposits:
	- deposit:
	amount:
	- amount: "100000000"
	denom: uctk
	depositor: certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03
	proposal_id: "1"
	tx_hash: 36f2d547c998d039386de7a4799fdab67f924c90211d3529ac69c7c9c9e46506
	pagination:
	next_key: null
	total: "0"
```
- `certik query gov params [flags]`: Query the all the parameters for the governance process.
```{engine = 'sh'}
$ certik query gov params

deposit_params:
	max_deposit_period: "172800000000000"
	min_deposit:
	- amount: "512000000"
	denom: uctk
	min_initial_deposit:
	- amount: "0"
	denom: uctk
	tally_params:
	certifier_update_security_vote_tally:
	quorum: "0.334000000000000000"
	threshold: "0.667000000000000000"
	veto_threshold: "0.334000000000000000"
	certifier_update_stake_vote_tally:
	quorum: "0.334000000000000000"
	threshold: "0.900000000000000000"
	veto_threshold: "0.334000000000000000"
	default_tally:
	quorum: "0.334000000000000000"
	threshold: "0.500000000000000000"
	veto_threshold: "0.334000000000000000"
	voting_params:
	voting_period: "172800000000000"
```

- `certik query gov param [param-type] [flags]`: Query the parameter that user wants to see for the governance process. For the parameter, user can refer to `certik query gov params` to see the type of parameter. For example, we want to see `voting` param, we use:
```{engine = 'sh'}
$ certik query gov param voting
voting_period: "172800000000000"
```

- `certik query gov proposal [proposal-id] [flags]`: Query details for a proposal. You can find the proposal-id by running "`certik query gov proposal`	".
```{engine = 'sh'}
$ certik query gov proposal 1

proposal:
	content:
	'@type': /shentu.shield.v1alpha1.ShieldClaimProposal
	description: Details of the attack
	evidence: Attack happened on <time> caused loss of <amount> to <account> by <txhash>
	loss:
	- amount: "1000000"
	denom: uctk
	pool_id: "1"
	proposal_id: "1"
	proposer: certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03
	purchase_id: "2"
	deposit_end_time: "2021-08-09T15:25:09.408965Z"
	final_tally_result:
	abstain: "0"
	"no": "0"
	no_with_veto: "0"
	"yes": "0"
	is_proposer_council_member: false
	proposal_id: "1"
	proposer_address: certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03
	status: PROPOSAL_STATUS_VALIDATOR_VOTING_PERIOD
	submit_time: "2021-08-09T15:25:09.408965Z"
	total_deposit:
	- amount: "100000000"
	denom: uctk
	voting_end_time: "2021-08-11T15:25:09.408965Z"
	voting_start_time: "2021-08-09T15:25:09.408965Z"
```

- `certik query gov proposals [proposal-id] [flags]`: Query details for a proposal. You can find the proposal-id by running "`certik query gov proposal`	".
```{engine = 'sh'}
$ certik query gov proposals

pagination:
next_key: null
total: "0"
proposals:

- content:
'@type': /shentu.shield.v1alpha1.ShieldClaimProposal
description: Details of the attack
evidence: Attack happened on <time> caused loss of <amount> to <account> by <txhash>
loss:
- amount: "1000000"
denom: uctk
pool_id: "1"
proposal_id: "1"
proposer: certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03
purchase_id: "2"
deposit_end_time: "2021-08-09T15:25:09.408965Z"
final_tally_result:
abstain: "0"
"no": "0"
no_with_veto: "0"
"yes": "0"
is_proposer_council_member: false
proposal_id: "1"
proposer_address: certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03
status: PROPOSAL_STATUS_VALIDATOR_VOTING_PERIOD
submit_time: "2021-08-09T15:25:09.408965Z"
total_deposit:
- amount: "100000000"
denom: uctk
voting_end_time: "2021-08-11T15:25:09.408965Z"
voting_start_time: "2021-08-09T15:25:09.408965Z"

- content:
'@type': /shentu.cert.v1alpha1.CertifierUpdateProposal
add_or_remove: add
alias: alice
certifier: certik1kgv37gjtl2vj6gg7n2ernkdynr8hj7cfrz0n94
description: Why we should make Alice a certifier
proposer: certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p
title: New Certifier, Alice
deposit_end_time: "2021-08-11T17:04:19.743009Z"
final_tally_result:
abstain: "0"
"no": "0"
no_with_veto: "0"
"yes": "1"
is_proposer_council_member: true
proposal_id: "2"
proposer_address: certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p
status: PROPOSAL_STATUS_PASSED
submit_time: "2021-08-09T17:04:19.743009Z"
total_deposit: []
voting_end_time: "2021-08-11T17:04:19.743009Z"
voting_start_time: "2021-08-09T17:04:19.743009Z"

- content:
'@type': /cosmos.distribution.v1beta1.CommunityPoolSpendProposal
amount: []
description: Pay me some Atoms!
recipient: certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03
title: Community Pool Spend
deposit_end_time: "2021-08-09T17:12:41.689127Z"
final_tally_result:
abstain: "0"
"no": "0"
no_with_veto: "0"
"yes": "0"
is_proposer_council_member: true
proposal_id: "3"
proposer_address: certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p
status: PROPOSAL_STATUS_VALIDATOR_VOTING_PERIOD
submit_time: "2021-08-09T17:12:41.689127Z"
total_deposit: []
voting_end_time: "2021-08-11T17:12:41.689127Z"
voting_start_time: "2021-08-09T17:12:41.689127Z"

- content:
'@type': /cosmos.params.v1beta1.ParameterChangeProposal
changes:
- key: MaxValidators
subspace: staking
value: "105"
description: Update max validators
title: Staking Param Change
deposit_end_time: "2021-08-09T17:14:22.812295Z"
final_tally_result:
abstain: "0"
"no": "0"
no_with_veto: "0"
"yes": "0"
is_proposer_council_member: true
proposal_id: "4"
proposer_address: certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p
status: PROPOSAL_STATUS_VALIDATOR_VOTING_PERIOD
submit_time: "2021-08-09T17:14:22.812295Z"
total_deposit: []
voting_end_time: "2021-08-11T17:14:22.812295Z"
voting_start_time: "2021-08-09T17:14:22.812295Z"

- content:
'@type': /cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal
description: Cancel Software
title: Test Cancel
deposit_end_time: "2021-08-09T17:21:32.447488Z"
final_tally_result:
abstain: "0"
"no": "0"
no_with_veto: "0"
"yes": "0"
is_proposer_council_member: true
proposal_id: "5"
proposer_address: certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p
status: PROPOSAL_STATUS_VALIDATOR_VOTING_PERIOD
submit_time: "2021-08-09T17:21:32.447488Z"
total_deposit: []
voting_end_time: "2021-08-11T17:21:32.447488Z"
voting_start_time: "2021-08-09T17:21:32.447488Z"
```

- `certik query gov proposer [proposal-id] [flags]`: Query which address proposed a proposal with a given ID. For example, user want to see the proposer of proposal 1: 
```{engine = 'sh'}
$ certik query gov proposer 1

proposal_id: "1"
proposer: certik1r9sanqz6vvv54htes625z4epfedl6y36eujg03
```

- `certik query gov tally [proposal-id] [flags]`: Query tally of votes on a proposal. You can find the proposal-id by running "`certik query gov proposals`". Currently deprecated. 


- `certik query gov vote [proposal-id] [voter-addr] [flags]`: Query details for a single vote on a proposal given its identifier. Currently deprecated. 

- `certik query gov votes [proposal-id] [flags]`: Query all vote details for a single proposal by its id. For example, user want to see all the votes from proposal 1: 
```{engine = 'sh'}
$ certik query gov votes 1

pagination:
next_key: null
total: "0"
votes:
- deposit:
option: VOTE_OPTION_YES
proposal_id: "1"
voter: certik1yvzwq93mnmlxeawza88evszxew74qj6tvw5h0p
tx_hash: 8236d3bf5bca568e5636d677dcb912143ee3f34954cc2ec71f2e49d5dab6e00a
- deposit:
option: VOTE_OPTION_YES
proposal_id: "1"
voter: certik1kgv37gjtl2vj6gg7n2ernkdynr8hj7cfrz0n94
tx_hash: 85ade5bb9c4bee246f5c10c0419fe66c455912db7fab4e8b8e92e9e019c5e275
```


