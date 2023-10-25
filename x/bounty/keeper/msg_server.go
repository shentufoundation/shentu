package keeper

import (
	"context"

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

	_, found := k.GetProgram(ctx, msg.ProgramId)
	if found {
		return nil, types.ErrProgramAlreadyExists
	}

	program, err := types.NewProgram(msg.ProgramId, msg.Name, msg.CreatorAddress, msg.Detail, msg.MemberAccounts, types.ProgramStatusInactive)
	if err != nil {
		return nil, err
	}

	k.SetProgram(ctx, program)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateProgram,
			sdk.NewAttribute(types.AttributeKeyProgramID, msg.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CreatorAddress),
		),
	})

	return &types.MsgCreateProgramResponse{}, nil
}

func (k msgServer) EditProgram(goCtx context.Context, msg *types.MsgEditProgram) (*types.MsgEditProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, found := k.GetProgram(ctx, msg.ProgramId)
	if !found {
		return nil, types.ErrNoProgramFound
	}

	// todo add UpdateDescription
	program.Detail = msg.Detail

	if len(msg.Name) != 0 {
		program.Name = msg.Name
	}
	if len(msg.MemberAccounts) != 0 {
		program.MemberAccounts = msg.MemberAccounts
	}

	k.SetProgram(ctx, program)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditProgram,
			sdk.NewAttribute(types.AttributeKeyProgramID, msg.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CreatorAddress),
		),
	})

	return &types.MsgEditProgramResponse{}, nil
}

func (k msgServer) OpenProgram(goCtx context.Context, msg *types.MsgModifyProgramStatus) (*types.MsgModifyProgramStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}
	if err = k.Keeper.OpenProgram(ctx, fromAddr, msg.ProgramId); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOpenProgram,
			sdk.NewAttribute(types.AttributeKeyProgramID, msg.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})
	return &types.MsgModifyProgramStatusResponse{}, nil
}

func (k msgServer) CloseProgram(goCtx context.Context, msg *types.MsgModifyProgramStatus) (*types.MsgModifyProgramStatusResponse, error) {
	fromAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	err = k.Keeper.CloseProgram(ctx, fromAddr, msg.ProgramId)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOpenProgram,
			sdk.NewAttribute(types.AttributeKeyProgramID, msg.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})
	return &types.MsgModifyProgramStatusResponse{}, nil
}

func (k msgServer) SubmitFinding(goCtx context.Context, msg *types.MsgSubmitFinding) (*types.MsgSubmitFindingResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	program, isExist := k.GetProgram(ctx, msg.ProgramId)
	if !isExist {
		return nil, types.ErrProgramNotExists
	}
	if program.Status != types.ProgramStatusActive {
		return nil, types.ErrProgramNotActive
	}

	submitTime := ctx.BlockHeader().Time

	_, found := k.GetFinding(ctx, msg.FindingId)
	if found {
		return nil, types.ErrFindingAlreadyExists
	}

	finding, err := types.NewFinding(msg.ProgramId, msg.FindingId, msg.Title, msg.SubmitterAddress, msg.Detail, submitTime)
	if err != nil {
		return nil, err
	}

	if err = k.AppendFidToFidList(ctx, msg.ProgramId, msg.FindingId); err != nil {
		return nil, err
	}

	k.SetFinding(ctx, finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSubmitFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SubmitterAddress),
		),
	})

	return &types.MsgSubmitFindingResponse{}, nil
}

func (k msgServer) AcceptFinding(goCtx context.Context, msg *types.MsgModifyFindingStatus) (*types.MsgModifyFindingStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	finding, err := k.hostProcess(ctx, msg.FindingId, msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	finding.Status = types.FindingStatusConfirmed
	k.SetFinding(ctx, *finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAcceptFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgModifyFindingStatusResponse{}, nil
}

func (k msgServer) RejectFinding(goCtx context.Context, msg *types.MsgModifyFindingStatus) (*types.MsgModifyFindingStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	finding, err := k.hostProcess(ctx, msg.FindingId, msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	finding.Status = types.FindingStatusClosed
	k.SetFinding(ctx, *finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRejectFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgModifyFindingStatusResponse{}, nil
}

func (k msgServer) hostProcess(ctx sdk.Context, fid, hostAddr string) (*types.Finding, error) {

	// get finding
	finding, isExist := k.GetFinding(ctx, fid)
	if !isExist {
		return nil, types.ErrFindingNotExists
	}
	// get program
	program, isExist := k.GetProgram(ctx, finding.ProgramId)
	if !isExist {
		return nil, types.ErrProgramNotExists
	}
	if program.Status != types.ProgramStatusActive {
		return nil, types.ErrProgramNotActive
	}

	// only host can update finding comment
	if program.AdminAddress != hostAddr {
		return nil, types.ErrProgramCreatorInvalid
	}

	return &finding, nil
}

func (k msgServer) CancelFinding(goCtx context.Context, msg *types.MsgModifyFindingStatus) (*types.MsgModifyFindingStatusResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// get finding
	finding, ok := k.GetFinding(ctx, msg.FindingId)
	if !ok {
		return nil, types.ErrFindingNotExists
	}

	// check submitter
	if finding.SubmitterAddress != msg.OperatorAddress {
		return nil, types.ErrFindingSubmitterInvalid
	}

	// check status
	if finding.Status != types.FindingStatusReported {
		return nil, types.ErrFindingStatusInvalid
	}

	k.DeleteFidFromFidList(ctx, finding.ProgramId, finding.FindingId)
	k.DeleteFinding(ctx, finding.FindingId)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, msg.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgModifyFindingStatusResponse{}, nil
}

func (k msgServer) ReleaseFinding(goCtx context.Context, msg *types.MsgReleaseFinding) (*types.MsgReleaseFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get finding
	finding, isExist := k.GetFinding(ctx, msg.FindingId)
	if !isExist {
		return nil, types.ErrFindingNotExists
	}
	// get program
	program, isExist := k.GetProgram(ctx, finding.ProgramId)
	if !isExist {
		return nil, types.ErrProgramNotExists
	}
	if program.Status != types.ProgramStatusActive {
		return nil, types.ErrProgramNotActive
	}

	if program.AdminAddress != msg.HostAddress {
		return nil, types.ErrProgramCreatorInvalid
	}

	k.SetFinding(ctx, finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReleaseFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, program.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.HostAddress),
		),
	})

	return &types.MsgReleaseFindingResponse{}, nil
}
