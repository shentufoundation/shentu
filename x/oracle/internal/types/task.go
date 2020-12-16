package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// func (t TaskStatus) String() string {
// 	switch t {
// 	case TaskStatusPending:
// 		return "pending"
// 	case TaskStatusSucceeded:
// 		return "succeeded"
// 	case TaskStatusFailed:
// 		return "failed"
// 	default:
// 		return "unknown task status"
// 	}
// }

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
		Creator:       creator.String(),
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

// NewResponse returns a new response.
func NewResponse(score sdk.Int, operator sdk.AccAddress) Response {
	return Response{
		Operator: operator.String(),
		Score:    score,
	}
}

type Responses []Response

// String implements the Stringer interface.
func (r Responses) String() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}
