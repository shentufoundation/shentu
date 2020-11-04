# Oracle

## State

type Withdraw struct {
	Address  sdk.AccAddress `json:"address"`
	Amount   sdk.Coins      `json:"amount"`
	DueBlock int64          `json:"due_block"`
}

type Operator struct {
	Address            sdk.AccAddress `json:"address"`
	Proposer           sdk.AccAddress `json:"proposer"`
	Collateral         sdk.Coins      `json:"collateral"`
	AccumulatedRewards sdk.Coins      `json:"accumulated_rewards"`
	Name               string         `json:"name"`
}

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

type Response struct {
	Operator sdk.AccAddress `json:"operator"`
	Score    sdk.Int        `json:"score"`
	Weight   sdk.Int        `json:"weight"`
	Reward   sdk.Coins      `json:"reward"`
}


## Messages

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
type MsgAddCollateral struct {
	Address             sdk.AccAddress
	CollateralIncrement sdk.Coins
}
type MsgReduceCollateral struct {
	Address             sdk.AccAddress
	CollateralDecrement sdk.Coins
}
type MsgWithdrawReward struct {
	Address sdk.AccAddress
}
type MsgCreateTask struct {
	Contract      string
	Function      string
	Bounty        sdk.Coins
	Description   string
	Creator       sdk.AccAddress
	Wait          int64
	ValidDuration time.Duration
}
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

type MsgDeleteTask struct {
	Contract string
	Function string
	Force    bool
	Deleter  sdk.AccAddress
}

## Events

## Parameters

type TaskParams struct {
	ExpirationDuration time.Duration `json:"task_expiration_duration"`
	AggregationWindow  int64         `json:"task_aggregation_window"`
	AggregationResult  sdk.Int       `json:"task_aggregation_result"`
	ThresholdScore     sdk.Int       `json:"task_threshold_score"`
	Epsilon1           sdk.Int       `json:"task_epsilon1"`
	Epsilon2           sdk.Int       `json:"task_epsilon2"`
}
type LockedPoolParams struct {
	LockedInBlocks    int64 `json:"locked_in_blocks"`
	MinimumCollateral int64 `json:"minimum_collateral"`
}
