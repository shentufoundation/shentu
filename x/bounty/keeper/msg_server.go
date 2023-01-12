package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
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

func (k msgServer) CreateProgram(goCtx context.Context, msg *types.MsgCreateProgram) (*types.MsgCreateProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creatorAddr, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, msg.Deposit)
	if err != nil {
		return nil, err
	}

	programID, err := k.GetNextProgramID(ctx)
	if err != nil {
		return nil, err
	}

	program := types.Program{
		ProgramId:         programID,
		CreatorAddress:    msg.CreatorAddress,
		SubmissionEndTime: msg.SubmissionEndTime,
		Description:       msg.Description,
		EncryptionKey:     msg.EncryptionKey,
		Deposit:           msg.Deposit,
		CommissionRate:    msg.CommissionRate,
	}

	k.SetProgram(ctx, program)

	k.SetNextProgramID(ctx, programID+1)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateProgram,
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(program.ProgramId, 10)),
			sdk.NewAttribute(types.AttributeKeyDeposit, sdk.NewCoins(msg.Deposit...).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CreatorAddress),
		),
	})

	return &types.MsgCreateProgramResponse{ProgramId: programID}, nil
}
