package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateTask is the message for creating a task.
type MsgCreateTask struct {
	Contract string
	Function string
}

// MsgRespondToTask is the message for responding to a task
type MsgRespondToTask struct {
	Contract string
	Function string
	Score    uint8
	Operator sdk.AccAddress
}
