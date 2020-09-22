package cvm

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/hyperledger/burrow/crypto"

	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

// NewHandler returns a handler for "cvm" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {

		case MsgCall:
			return handleMsgCall(ctx, keeper, msg)

		case MsgDeploy:
			return handleMsgDeploy(ctx, keeper, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized cert Msg type: %v", msg)
		}
	}
}

func handleMsgCall(ctx sdk.Context, keeper Keeper, msg MsgCall) (*sdk.Result, error) {
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	result, err := keeper.Call(ctx, msg.Caller, msg.Callee, msg.Value, msg.Data, nil, false, false, false)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Caller.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, strconv.FormatUint(msg.Value, 10)),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCall,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Callee.String()),
			sdk.NewAttribute(types.AttributeKeyValue, strconv.FormatUint(msg.Value, 10)),
		),
	)

	return &sdk.Result{
		Data:   result,
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgDeploy(ctx sdk.Context, keeper Keeper, msg MsgDeploy) (*sdk.Result, error) {
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	result, err := keeper.Call(ctx, msg.Caller, nil, msg.Value, msg.Code, msg.Meta, false, msg.IsEWASM, msg.IsRuntime)
	if err != nil {
		return nil, err
	}
	keeper.SetAbi(ctx, crypto.MustAddressFromBytes(result), []byte(msg.Abi))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Caller.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, strconv.FormatUint(msg.Value, 10)),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeploy,
			sdk.NewAttribute(types.AttributeKeyNewContractAddress, sdk.AccAddress(result).String()),
			sdk.NewAttribute(types.AttributeKeyValue, strconv.FormatUint(msg.Value, 10)),
		),
	)
	// Print bech32 address of the new contract in the log.
	return &sdk.Result{
		Data:   result,
		Events: ctx.EventManager().Events(),
	}, nil
}
