package types

import (
	"encoding/json"
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

// NewTask returns a new task.
func NewTask(
	contract string,
	function string,
	beginBlock int64,
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
		BeginBlock:    beginBlock,
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
	Operator sdk.AccAddress `json:"operator"`
	Score    sdk.Int        `json:"score"`
	Weight   sdk.Int        `json:"weight"`
	Reward   sdk.Coins      `json:"reward"`
}

// NewResponse returns a new response.
func NewResponse(score sdk.Int, operator sdk.AccAddress) Response {
	return Response{
		Operator: operator,
		Score:    score,
	}
}

// String implements the Stringer interface.
func (r Response) String() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// Responses defines a list of responses.
type Responses []Response

// String implements the Stringer interface.
func (r Responses) String() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}
