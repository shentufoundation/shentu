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

	program, err := k.Keeper.CreateProgram(ctx, msg)
	if err != nil {
		return nil, err
	}

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

	return &types.MsgCreateProgramResponse{ProgramId: program.ProgramId}, nil
}

func (k msgServer) SubmitFinding(goCtx context.Context, msg *types.MsgSubmitFinding) (*types.MsgSubmitFindingResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	finding, err := k.Keeper.SubmitFinding(ctx, msg)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSubmitFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, strconv.FormatUint(finding.FindingId, 10)),
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(finding.ProgramId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SubmitterAddress),
		),
	})

	return &types.MsgSubmitFindingResponse{
		FindingId: finding.FindingId,
	}, nil
}

func (k msgServer) HostAcceptFinding(goCtx context.Context, msg *types.MsgHostAcceptFinding) (*types.MsgHostAcceptFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	finding, err := k.Keeper.hostProcessFinding(ctx, msg.FindingId, msg.HostAddress, msg.EncryptedComment)
	if err != nil {
		return nil, err
	}

	finding.FindingStatus = types.FindingStatusValid
	k.SetFinding(ctx, *finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAcceptFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, strconv.FormatUint(finding.FindingId, 10)),
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(finding.ProgramId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.HostAddress),
		),
	})

	return &types.MsgHostAcceptFindingResponse{}, nil
}

func (k msgServer) HostRejectFinding(goCtx context.Context, msg *types.MsgHostRejectFinding) (*types.MsgHostRejectFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	finding, err := k.Keeper.hostProcessFinding(ctx, msg.FindingId, msg.HostAddress, msg.EncryptedComment)
	if err != nil {
		return nil, err
	}

	finding.FindingStatus = types.FindingStatusInvalid
	k.SetFinding(ctx, *finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRejectFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, strconv.FormatUint(finding.FindingId, 10)),
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(finding.ProgramId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.HostAddress),
		),
	})

	return &types.MsgHostRejectFindingResponse{}, nil
}

func (k msgServer) CancelFinding(goCtx context.Context, msg *types.MsgCancelFinding) (*types.MsgCancelFindingResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	finding, err := k.Keeper.CancelFinding(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, strconv.FormatUint(msg.FindingId, 10)),
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(finding.ProgramId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SubmitterAddress),
		),
	})

	return &types.MsgCancelFindingResponse{}, nil
}

func (k msgServer) ReleaseFinding(goCtx context.Context, msg *types.MsgReleaseFinding) (*types.MsgReleaseFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	finding, err := k.Keeper.ReleaseFinding(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReleaseFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, strconv.FormatUint(finding.FindingId, 10)),
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(finding.ProgramId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.HostAddress),
		),
	})

	return &types.MsgReleaseFindingResponse{}, nil
}

func (k msgServer) EndProgram(goCtx context.Context, msg *types.MsgEndProgram) (*types.MsgEndProgramResponse, error) {
	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	err = k.Keeper.EndProgram(ctx, fromAddr, msg.ProgramId)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEndProgram,
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(msg.ProgramId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})
	return &types.MsgEndProgramResponse{}, nil
}
