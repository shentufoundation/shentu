package keeper

import (
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/exec"

	"github.com/shentufoundation/shentu/v2/x/cvm/types"
)

type EventSink struct {
	ctx sdk.Context
}

func NewEventSink(ctx sdk.Context) *EventSink {
	return &EventSink{ctx}
}

func (es *EventSink) Call(call *exec.CallEvent, exception *errors.Exception) error {
	// do not log anything on the first call
	if call.StackDepth == 0 {
		return nil
	}

	caller, err := sdk.AccAddressFromHex(call.CallData.Caller.String())
	if err != nil {
		return err
	}

	callee, err := sdk.AccAddressFromHex(call.CallData.Callee.String())
	if err != nil {
		return err
	}

	es.ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInternalCall,
			sdk.Attribute{
				Key:   "caller",
				Value: caller.String(),
			},
			sdk.Attribute{
				Key:   "callee",
				Value: callee.String(),
			},
			sdk.Attribute{
				Key:   "data",
				Value: call.CallData.Data.String(),
			},
			sdk.Attribute{
				Key:   "value",
				Value: hex.EncodeToString(call.CallData.Value),
			},
			sdk.Attribute{
				Key:   "stack-depth",
				Value: strconv.FormatUint(call.StackDepth, 10),
			},
		),
	)
	return nil
}

func (es *EventSink) Log(log *exec.LogEvent) error {
	topicsString := ""
	for _, topic := range log.Topics {
		topicsString += topic.String()
	}

	b32addr, err := sdk.AccAddressFromHex(log.Address.String())
	if err != nil {
		panic("address data in CVM is corrupted")
	}

	es.ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCVMEvent,
			sdk.Attribute{
				Key:   "address",
				Value: b32addr.String(),
			},
			sdk.Attribute{
				Key:   "topics",
				Value: topicsString,
			},
			sdk.Attribute{
				Key:   "data",
				Value: log.Data.String(),
			},
		),
	)
	return nil
}

func (es *EventSink) Print(print *exec.PrintEvent) error {
	return nil
}
