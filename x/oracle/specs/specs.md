# Oracle

The `oracle` module facilitates real-time security checks by providing on-chain security scores from audited smart contracts that can be queried by a Security Oracle on any chain.
CertiK Chain users can serve as `Operator`s, which provide queriers with security scores for audited smart contracts and earn rewards for doing so. The reward an operator receives is proportional to the amount of collateral the operator locks up. Also, operators' security scores for a given smart contract are weighted by their collateral. In this way, the `oracle` module serves as the Oracle Combinator, which combines the various results from each operator into a composite score.

See the [whitepaper](https://www.certik.foundation/whitepaper#2-CertiK-Security-Oracle) for more information on the CertiK Security Oracle.


## State

`Withdraw` stores a withdraw of `Amount` scheduled for a given `DueBlock`. A withdraw is scheduled when `ReduceCollateral` or `RemoveOperator` is called.

```go
type Withdraw struct {
	Address  sdk.AccAddress `json:"address"`
	Amount   sdk.Coins      `json:"amount"`
	DueBlock int64          `json:"due_block"`
}
```

`Operator` is a chain operator, which can be added by `CreateOperator`.

```go
type Operator struct {
	Address            sdk.AccAddress `json:"address"`
	Proposer           sdk.AccAddress `json:"proposer"`
	Collateral         sdk.Coins      `json:"collateral"`
	AccumulatedRewards sdk.Coins      `json:"accumulated_rewards"`
	Name               string         `json:"name"`
}
```

`Task` stores a request to generate a score for a given smart contract.

```go
type Task struct {
	Contract      string         `json:"contract"`
	Function      string         `json:"function"`
	BeginBlock    int64          `json:"begin_block"`
	Bounty        sdk.Coins      `json:"bounty"`
	Description   string         `json:"string"`
	Expiration    time.Time      `json:"expiration"`
	Creator       sdk.AccAddress `json:"creator"`
	Responses     Responses      `json:"responses"`
	Result        sdk.Int        `json:"result"`
	ClosingBlock  int64          `json:"closing_block"`
	WaitingBlocks int64          `json:"waiting_blocks"`
	Status        TaskStatus     `json:"status"`
}

type TaskID struct {
	Contract string `json:"contract"`
	Function string `json:"function"`
}
```

`Response` contains the score from an operator, which will be combined with other responses to yield the aggregate score for a smart contract.

```go
type Response struct {
	Operator sdk.AccAddress `json:"operator"`
	Score    sdk.Int        `json:"score"`
	Weight   sdk.Int        `json:"weight"`
	Reward   sdk.Coins      `json:"reward"`
}
```

## Messages

### Operators

`MsgCreateOperator` adds `Address` as a new `Operator` and adds their `Collateral` to the collateral pool. Likewise, `MsgRemoveOperator` removes an operator.

```go
type MsgCreateOperator struct {
	Address    sdk.AccAddress
	Collateral sdk.Coins
	Proposer   sdk.AccAddress
	Name       string
}

type MsgRemoveOperator struct {
	Address  sdk.AccAddress
	Proposer sdk.AccAddress
}
```

`MsgAddCollateral` and `MsgReduceCollateral` increase or decrease collateral for a given operator at `Address`.

```go
type MsgAddCollateral struct {
	Address             sdk.AccAddress
	CollateralIncrement sdk.Coins
}

type MsgReduceCollateral struct {
	Address             sdk.AccAddress
	CollateralDecrement sdk.Coins
}
```

`MsgWithdrawReward` transfers the accumulated rewards to the operator at `Address`.

```go
type MsgWithdrawReward struct {
	Address sdk.AccAddress
}
```

### Tasks

`MsgCreateTask` creates a new `Task`. After the `ValidDuration` has passed, it can be removed with `MsgDeleteTask` by its `Creator`. It is not removed automatically.

```go
type MsgCreateTask struct {
	Contract      string
	Function      string
	Bounty        sdk.Coins
	Description   string
	Creator       sdk.AccAddress
	Wait          int64
	ValidDuration time.Duration
}

type MsgDeleteTask struct {
	Contract string
	Function string
	Force    bool
	Deleter  sdk.AccAddress
}
```

While a `Task` is active, operators can submit scores for the task's contract. The task's `Result` can be queried with `MsgInquiryTask`.

```go
type MsgTaskResponse struct {
	Contract string
	Function string
	Score    int64
	Operator sdk.AccAddress
}

type MsgInquiryTask struct {
	Contract string
	Function string
	TxHash   string
	Inquirer sdk.AccAddress
}
```

## Parameters
| Parameter            | Info                                                                         | Default  |
|----------------------|------------------------------------------------------------------------------|----------|
| `ExpirationDuration` | default task duration, for tasks with unspecified durations                  | 24 hours |
| `AggregationWindow`  | number of blocks between task creation and calculation of final score        | 20       |
| `AggregationResult`  | aggregation result for a task with no responses                              | 50       |
| `ThresholdScore`     | threshold above/below which a contract is considered secure/insecure         | 50       |
| `Epsilon1`           | distribution curve parameter                                                 | 1        |
| `Epsilon2`           | distribution curve parameter                                                 | 100      |
| `LockedInBlocks`     | number of blocks operators need to wait before getting their collateral back | 30       |
