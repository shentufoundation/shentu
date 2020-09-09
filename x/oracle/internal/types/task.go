package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TaskStatus defines the data type of the status of a task.
type TaskStatus byte

const (
	TaskStatusNil = iota
	TaskStatusPending
	TaskStatusSucceeded
	TaskStatusFailed
)

func (t TaskStatus) String() string {
	switch t {
	case TaskStatusPending:
		return "pending"
	case TaskStatusSucceeded:
		return "succeeded"
	case TaskStatusFailed:
		return "failed"
	default:
		return "unknown task status"
	}
}

// Task defines the data structure of a task.
type Task struct {
	Contract      string         `json:"contract"`
	Function      string         `json:"function"`
	Bounty        sdk.Coins      `json:"bounty"`
	Description   string         `json:"string"`
	Expiration    time.Time      `json:"expiration"`
	Creator       sdk.AccAddress `json:"creator"`
	Responses     []Response     `json:"responses"`
	Result        sdk.Int        `json:"result"`
	ClosingBlock  int64          `json:"closing_block"`
	WaitingBlocks int64          `json:"waiting_blocks"`
	Status        TaskStatus     `json:"status"`
}

// NewTask returns a new task.
func NewTask(
	contract string,
	function string,
	bounty sdk.Coins,
	description string,
	expiration time.Time,
	creator sdk.AccAddress,
	closingBlock int64,
	waitingBlocks int64,
) Task {
	return Task{
		Contract:      contract,
		Function:      function,
		Bounty:        bounty,
		Description:   description,
		Expiration:    expiration,
		Creator:       creator,
		ClosingBlock:  closingBlock,
		WaitingBlocks: waitingBlocks,
		Status:        TaskStatusPending,
	}
}

// TaskID defines the data structure of the ID of a task.
type TaskID struct {
	Contract string `json:"contract"`
	Function string `json:"function"`
}

// Response defines the data structure of a response.
type Response struct {
	Contract string         `json:"contract"`
	Function string         `json:"function"`
	Score    sdk.Int        `json:"score"`
	Operator sdk.AccAddress `json:"operator"`
}

// NewResponse returns a new response.
func NewResponse(contract, function string, score sdk.Int, operator sdk.AccAddress) Response {
	return Response{
		Contract: contract,
		Function: function,
		Score:    score,
		Operator: operator,
	}
}
