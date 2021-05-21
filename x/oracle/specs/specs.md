# Oracle

The `oracle` module facilitates real-time security checks by providing on-chain security scores from audited smart contracts that can be queried by a Security Oracle on any chain.
CertiK Chain users can serve as `Operator`s, which provide queriers with security scores for audited smart contracts and earn rewards for doing so. The reward an operator receives is proportional to the amount of collateral the operator locks up. Also, operators' security scores for a given smart contract are weighted by their collateral. In this way, the `oracle` module serves as the Oracle Combinator, which combines the various results from each operator into a composite score.

See the [whitepaper](https://www.certik.foundation/whitepaper#2-CertiK-Security-Oracle) for more information on the CertiK Security Oracle.


## State

### Operators

`Operator` is a chain operator, which can be added by `CreateOperator`.

- Operator: `0x1 | Address -> amino(operator)`

```go
type Operator struct {
	Address            string		`json:"address" yaml:"address"`
	Proposer           string		`json:"proposer" yaml:"proposer"`
	Collateral         sdk.Coins	`json:"collateral" yaml:"collateral"`
	AccumulatedRewards sdk.Coins	`json:"accumulated_rewards" yaml:"accumulated_rewards"`
	Name               string		`json:"name" yaml:"name"`
}
```

### Withdraws

`Withdraw` stores a withdraw of `Amount` scheduled for a given `DueBlock`. A withdraw is scheduled when `ReduceCollateral` or `RemoveOperator` is called.

- Withdraw: `0x2 | LittleEndian(DueBlock) | Address -> amino(withdraw)`

```go
type Withdraw struct {
	Address  string		`json:"address" yaml:"address"`
	Amount   sdk.Coins	`json:"amount" yaml:"amount"`
	DueBlock int64		`json:"due_block" yaml:"due_block"`
}
```

### Collateral

`TotalCollateral` stores the total amount of collateral accumulated in the pool. It is stored as a `CoinsProto` object that includes all collateralized coins in the pool.

- TotalCollateral: `0x3 -> amino(collateral)`

```go
type CoinsProto struct {
	Coins sdk.Coins	`json:"coins" yaml:"amount"`
}
```

### Tasks

`Task` stores a request to generate a score for a given smart contract.

- Task: `0x4 | Contract | Function -> amino(task)`

```go
type Task struct {
	Contract      string		`json:"contract" yaml:"contract"`
	Function      string		`json:"function" yaml:"function"`
	BeginBlock    int64			`json:"begin_block" yaml:"begin_block"`
	Bounty        sdk.Coins 	`json:"bounty" yaml:"bounty"`
	Description   string		`json:"description" yaml:"description"`
	Expiration    time.Time		`json:"expiration" yaml:"expiration"`
	Creator       string		`json:"creator" yaml:"creator"`
	Responses     Responses		`json:"responses" yaml:"responses"`
	Result        sdk.Int   	`json:"result" yaml:"result"`
	ClosingBlock  int64			`json:"closing_block" yaml:"closing_block"`
	WaitingBlocks int64			`json:"waiting_blocks" yaml:"waiting_blocks"`
	Status        TaskStatus	`json:"status" yaml:"status"`
}

type TaskID struct {
	Contract string `json:"contract" yaml:"contract"`
	Function string `json:"function" yaml:"function"`
}
```

An operator can respond to the task by providing a valid `Response`, which contains the score from the operator. Scores from multiple responses will be combined to yield the aggregate score for a smart contract.

```go
type Response struct {
	Operator string		`json:"operator" yaml:"operator"`
	Score    sdk.Int	`json:"score" yaml:"score"`
	Weight   sdk.Int	`json:"weight" yaml:"weight"`
	Reward   sdk.Coins	`json:"reward" yaml:"reward"`
}
```

A list of tasks in the aggregation block is stored as a `TaskIDs` object. The list is updated upon every update of a task until closing block is reached, from which aggregation occurs.

- ClosingTaskID: `0x5 | LittleEndian(BlockHeight) -> amino(TaskIds)`

```go
type TaskIDs struct {
	TaskIds []TaskID `json:"task_ids"`
}
```

## Messages

### Operators

`MsgCreateOperator` adds `Address` as a new `Operator` and adds their `Collateral` to the collateral pool. Likewise, `MsgRemoveOperator` removes an operator.

```go
type MsgCreateOperator struct {
	Address    string		`json:"address" yaml:"address"`
	Collateral sdk.Coins	`json:"collateral" yaml:"collateral"`
	Proposer   string		`json:"proposer" yaml:"proposer"`
	Name       string		`json:"name" yaml:"name"`
}

type MsgRemoveOperator struct {
	Address  string `json:"address" yaml:"address"`
	Proposer string `json:"proposer" yaml:"proposer"`
}
```

`MsgAddCollateral` and `MsgReduceCollateral` increase or decrease collateral for a given operator at `Address`.

```go
type MsgAddCollateral struct {
	Address             string		`json:"address" yaml:"address"`
	CollateralIncrement sdk.Coins	`json:"collateral_increment" yaml:"collateral_increment"`
}

type MsgReduceCollateral struct {
	Address             string		`json:"address" yaml:"address"`
	CollateralDecrement sdk.Coins	`json:"collateral_decrement" yaml:"collateral_decrement"`
}
```

`MsgWithdrawReward` transfers the accumulated rewards to the operator at `Address`.

```go
type MsgWithdrawReward struct {
	Address string `json:"address" yaml:"address"`
}
```

### Tasks

`MsgCreateTask` creates a new `Task`. After the `ValidDuration` has passed, it can be removed with `MsgDeleteTask` by its `Creator`. It is not removed automatically.

```go
type MsgCreateTask struct {
	Contract      string		`json:"contract" yaml:"contract"`
	Function      string		`json:"function" yaml:"function"`
	Bounty        sdk.Coins		`json:"bounty" yaml:"bounty"`
	Description   string		`json:"description" yaml:"description"`
	Creator       string		`json:"creator," yaml:"creator"`
	Wait          int64			`json:"wait" yaml:"wait"`
	ValidDuration time.Duration	`json:"valid_duration" yaml:"valid_duration"`
}

type MsgDeleteTask struct {
	Contract string `json:"contract" yaml:"contract"`
	Function string `json:"function" yaml:"function"`
	Force    bool   `json:"force" yaml:"force"`
	Deleter  string `json:"deleter" yaml:"deleter"`
}
```

While a `Task` is active, operators can submit scores for the task's contract. The task's `Result` can be queried with `MsgInquiryTask`.

```go
type MsgTaskResponse struct {
	Contract string `json:"contract" yaml:"contract"`
	Function string `json:"function" yaml:"function"`
	Score    int64  `json:"score" yaml:"score"`
	Operator string `json:"operator" yaml:"operator"`
}

type MsgInquiryTask struct {
	Contract string `json:"contract" yaml:"contract"`
	Function string `json:"function" yaml:"function"`
	TxHash   string `json:"tx_hash" yaml:"tx_hash"`
	Inquirer string `json:"inquirer" yaml:"inquirer"`
}
```

## Parameters

### TaskParams

| Parameter            | Info                                                                         | Default  |
|----------------------|------------------------------------------------------------------------------|----------|
| `ExpirationDuration` | default task duration, for tasks with unspecified durations                  | 24 hours |
| `AggregationWindow`  | number of blocks between task creation and calculation of final score        | 20       |
| `AggregationResult`  | aggregation result for a task with no responses                              | 50       |
| `ThresholdScore`     | threshold above/below which a contract is considered secure/insecure         | 50       |
| `Epsilon1`           | distribution curve parameter                                                 | 1        |
| `Epsilon2`           | distribution curve parameter                                                 | 100      |


### LockedPoolParams

| Parameter            | Info                                                                         | Default  |
|----------------------|------------------------------------------------------------------------------|----------|
| `LockedInBlocks`     | number of blocks operators need to wait before getting their collateral back | 30       |
| `MinimumCollateral`  | minimum amount of collateral in a pool										  | 50000    |
