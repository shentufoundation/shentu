package migrate

import (
	"time"

	oracletypes "github.com/certikfoundation/shentu/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Operator struct {
	Address            sdk.AccAddress `json:"address"`
	Proposer           sdk.AccAddress `json:"proposer"`
	Collateral         sdk.Coins      `json:"collateral"`
	AccumulatedRewards sdk.Coins      `json:"accumulated_rewards"`
	Name               string         `json:"name"`
}

type ShieldWithdraw struct {
	Address  sdk.AccAddress `json:"address"`
	Amount   sdk.Coins      `json:"amount"`
	DueBlock int64          `json:"due_block"`
}

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

// Response defines the data structure of a response.
type Response struct {
	Operator sdk.AccAddress `json:"operator"`
	Score    sdk.Int        `json:"score"`
	Weight   sdk.Int        `json:"weight"`
	Reward   sdk.Coins      `json:"reward"`
}

type Responses []Response

// Task defines the data structure of a task.
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

// TaskStatus defines the data type of the status of a task.
type TaskStatus byte

const (
	TaskStatusNil = iota
	TaskStatusPending
	TaskStatusSucceeded
	TaskStatusFailed
)

type OracleGenesisState struct {
	Operators       []Operator       `json:"operators"`
	TotalCollateral sdk.Coins        `json:"total_collateral"`
	PoolParams      LockedPoolParams `json:"pool_params"`
	TaskParams      TaskParams       `json:"task_params"`
	Withdraws       []ShieldWithdraw `json:"withdraws"`
	Tasks           []Task           `json:"tasks"`
}

func migrateOracle(oldState OracleGenesisState) *oracletypes.GenesisState {
	var newOperators oracletypes.Operators
	for _, o := range oldState.Operators {
		newOperators = append(newOperators, oracletypes.Operator{
			Address:            o.Address.String(),
			Proposer:           o.Proposer.String(),
			Collateral:         o.Collateral,
			AccumulatedRewards: o.AccumulatedRewards,
			Name:               o.Name,
		})
	}

	var newWithdraws oracletypes.Withdraws
	for _, w := range oldState.Withdraws {
		newWithdraws = append(newWithdraws, oracletypes.Withdraw{
			Address:  w.Address.String(),
			Amount:   w.Amount,
			DueBlock: w.DueBlock,
		})
	}

	var newTasks []oracletypes.Task
	for _, t := range oldState.Tasks {
		var newResponses oracletypes.Responses
		for _, r := range t.Responses {
			newResponses = append(newResponses, oracletypes.Response{
				Operator: r.Operator.String(),
				Score:    r.Score,
				Weight:   r.Weight,
				Reward:   r.Reward,
			})
		}

		newTasks = append(newTasks, oracletypes.Task{
			Contract:      t.Contract,
			Function:      t.Function,
			BeginBlock:    t.BeginBlock,
			Bounty:        t.Bounty,
			Description:   t.Description,
			Expiration:    t.Expiration,
			Creator:       t.Creator.String(),
			Responses:     newResponses,
			Result:        t.Result,
			ClosingBlock:  t.ClosingBlock,
			WaitingBlocks: t.WaitingBlocks,
			Status:        oracletypes.TaskStatus(t.Status),
		})
	}

	return &oracletypes.GenesisState{
		Operators:       newOperators,
		TotalCollateral: oldState.TotalCollateral,
		PoolParams: &oracletypes.LockedPoolParams{
			LockedInBlocks:    oldState.PoolParams.LockedInBlocks,
			MinimumCollateral: oldState.PoolParams.MinimumCollateral,
		},
		TaskParams: &oracletypes.TaskParams{
			ExpirationDuration: oldState.TaskParams.ExpirationDuration,
			AggregationWindow:  oldState.TaskParams.AggregationWindow,
			AggregationResult:  oldState.TaskParams.AggregationResult,
			ThresholdScore:     oldState.TaskParams.ThresholdScore,
			Epsilon1:           oldState.TaskParams.Epsilon1,
			Epsilon2:           oldState.TaskParams.Epsilon2,
		},
		Withdraws: newWithdraws,
		Tasks:     newTasks,
	}
}
