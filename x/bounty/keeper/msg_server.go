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

	operatorAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}
	_, found := k.GetProgram(ctx, msg.ProgramId)
	if found {
		return nil, types.ErrProgramAlreadyExists
	}

	program, err := types.NewProgram(msg.ProgramId, msg.Name, msg.Description, operatorAddr, msg.MemberAccounts, types.ProgramStatusInactive, msg.BountyLevels)
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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgCreateProgramResponse{}, nil
}

func (k msgServer) EditProgram(goCtx context.Context, msg *types.MsgEditProgram) (*types.MsgEditProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	program, found := k.GetProgram(ctx, msg.ProgramId)
	if !found {
		return nil, types.ErrNoProgramFound
	}

	// Check the permissions. Only the admin of the program or cert address can operate.
	if program.AdminAddress != msg.OperatorAddress {
		if !k.certKeeper.IsCertifier(ctx, operatorAddr) {
			return nil, types.ErrProgramOperatorNotAllowed
		}
	}

	if len(msg.Name) > 0 {
		program.Name = msg.Name
	}
	if len(msg.Description) > 0 {
		program.Description = msg.Description
	}
	if len(msg.MemberAccounts) > 0 {
		program.MemberAccounts = msg.MemberAccounts
	}
	if len(msg.BountyLevels) > 0 {
		program.BountyLevels = msg.BountyLevels
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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgEditProgramResponse{}, nil
}

func (k msgServer) OpenProgram(goCtx context.Context, msg *types.MsgOpenProgram) (*types.MsgOpenProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := k.GetProgram(ctx, msg.ProgramId)
	if !found {
		return nil, types.ErrNoProgramFound
	}

	operatorAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	if err = k.Keeper.OpenProgram(ctx, msg.ProgramId, operatorAddr); err != nil {
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
	return &types.MsgOpenProgramResponse{}, nil
}

func (k msgServer) CloseProgram(goCtx context.Context, msg *types.MsgCloseProgram) (*types.MsgCloseProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.CloseProgram(ctx, msg.ProgramId, operatorAddr)
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
	return &types.MsgCloseProgramResponse{}, nil
}

func (k msgServer) SubmitFinding(goCtx context.Context, msg *types.MsgSubmitFinding) (*types.MsgSubmitFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
	if err != nil {
		return nil, err
	}

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

	finding, err := types.NewFinding(msg.ProgramId, msg.FindingId, msg.Title, msg.Description, operatorAddr, submitTime, msg.SeverityLevel)
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

func (k msgServer) AcceptFinding(goCtx context.Context, msg *types.MsgAcceptFinding) (*types.MsgAcceptFindingResponse, error) {
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

	return &types.MsgAcceptFindingResponse{}, nil
}

func (k msgServer) RejectFinding(goCtx context.Context, msg *types.MsgRejectFinding) (*types.MsgRejectFindingResponse, error) {
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

	return &types.MsgRejectFindingResponse{}, nil
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

func (k msgServer) CloseFinding(goCtx context.Context, msg *types.MsgCloseFinding) (*types.MsgCloseFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	// get finding
	finding, ok := k.GetFinding(ctx, msg.FindingId)
	if !ok {
		return nil, types.ErrFindingNotExists
	}

	// check submitter
	if finding.SubmitterAddress != msg.OperatorAddress || !k.certKeeper.IsCertifier(ctx, operatorAddr) {
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

	return &types.MsgCloseFindingResponse{}, nil
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

	if program.AdminAddress != msg.OperatorAddress {
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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgReleaseFindingResponse{}, nil
}
