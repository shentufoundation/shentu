package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/crypto"

	"github.com/shentufoundation/shentu/v2/x/cvm/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Deploy(goCtx context.Context, msg *types.MsgDeploy) (*types.MsgDeployResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	result, err := k.Keeper.Deploy(ctx, msg)
	if err != nil {
		return nil, err
	}

	addr, err := crypto.AddressFromBytes(result)
	if err != nil {
		return nil, err
	}

	k.Keeper.SetAbi(ctx, addr, []byte(msg.Abi))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Caller),
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

	return &types.MsgDeployResponse{
		Result: result,
	}, nil
}

func (k msgServer) Call(goCtx context.Context, msg *types.MsgCall) (*types.MsgCallResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	result, err := k.Keeper.Call(ctx, msg, false)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Caller),
			sdk.NewAttribute(sdk.AttributeKeyAmount, strconv.FormatUint(msg.Value, 10)),
		),
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCall,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Callee),
			sdk.NewAttribute(types.AttributeKeyValue, strconv.FormatUint(msg.Value, 10)),
		),
	)

	return &types.MsgCallResponse{
		Result: result,
	}, nil
}
