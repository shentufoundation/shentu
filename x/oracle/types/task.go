package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

type TaskI interface {
	GetID() []byte
	GetCreator() string
	GetResponses() []Response
	IsExpired(ctx sdk.Context) bool
	GetStatus() TaskStatus
	GetScore() int64
}


