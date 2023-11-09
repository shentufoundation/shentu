package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

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

	program, err := types.NewProgram(msg.ProgramId, msg.Name, msg.Detail, operatorAddr, types.ProgramStatusInactive, msg.BountyLevels)
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

	// check the status.
	// inactive: program admin, cert
	// active: certificate
	switch program.Status {
	case types.ProgramStatusInactive:
		if program.AdminAddress != msg.OperatorAddress && !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) {
			return nil, types.ErrProgramOperatorNotAllowed
		}
	case types.ProgramStatusActive:
		if !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) {
			return nil, types.ErrProgramOperatorNotAllowed
		}
	default:
		return nil, types.ErrProgramOperatorNotAllowed
	}

	if len(msg.Name) > 0 {
		program.Name = msg.Name
	}
	if len(msg.Detail) > 0 {
		program.Detail = msg.Detail
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

func (k msgServer) ActivateProgram(goCtx context.Context, msg *types.MsgActivateProgram) (*types.MsgActivateProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	if err = k.Keeper.OpenProgram(ctx, msg.ProgramId, operatorAddr); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeActivateProgram,
			sdk.NewAttribute(types.AttributeKeyProgramID, msg.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})
	return &types.MsgActivateProgramResponse{}, nil
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
			types.EventTypeCloseProgram,
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

	_, found := k.GetFinding(ctx, msg.FindingId)
	if found {
		return nil, types.ErrFindingAlreadyExists
	}

	submitTime := ctx.BlockHeader().Time
	finding, err := types.NewFinding(msg.ProgramId, msg.FindingId, msg.Title, msg.Detail, msg.FindingHash, operatorAddr, submitTime, msg.SeverityLevel)
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

func (k msgServer) EditFinding(goCtx context.Context, msg *types.MsgEditFinding) (*types.MsgEditFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	program, isExist := k.GetProgram(ctx, msg.ProgramId)
	if !isExist {
		return nil, types.ErrProgramNotExists
	}
	if program.Status != types.ProgramStatusActive {
		return nil, types.ErrProgramNotActive
	}

	finding, found := k.GetFinding(ctx, msg.FindingId)
	if !found {
		return nil, types.ErrFindingNotExists
	}
	// check submitter
	if finding.SubmitterAddress != msg.SubmitterAddress {
		return nil, types.ErrFindingSubmitterInvalid
	}
	if len(msg.Title) > 0 {
		finding.Title = msg.Title
	}
	if len(msg.FindingHash) > 0 {
		finding.FindingHash = msg.FindingHash
	}
	if msg.SeverityLevel != types.Unspecified {
		finding.SeverityLevel = msg.SeverityLevel
	}
	if len(msg.Detail) > 0 {
		finding.Detail = msg.Detail
	}

	k.SetFinding(ctx, finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SubmitterAddress),
		),
	})

	return &types.MsgEditFindingResponse{}, nil
}

func (k msgServer) ConfirmFinding(goCtx context.Context, msg *types.MsgConfirmFinding) (*types.MsgConfirmFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	finding, err := k.hostProcess(ctx, msg.FindingId, msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	finding.Status = types.FindingStatusConfirmed
	k.SetFinding(ctx, *finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeConfirmFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgConfirmFindingResponse{}, nil
}

func (k msgServer) hostProcess(ctx sdk.Context, fid, hostAddr string) (*types.Finding, error) {

	// get finding
	finding, found := k.GetFinding(ctx, fid)
	if !found {
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

	// todo fingerprint comparison
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
	// check finding status: StatusSubmitted and StatusActive can be closed
	if finding.Status != types.FindingStatusSubmitted && finding.Status != types.FindingStatusActive {
		return nil, types.ErrFindingStatusInvalid
	}
	// get program
	program, isExist := k.GetProgram(ctx, finding.ProgramId)
	if !isExist {
		return nil, types.ErrProgramNotExists
	}

	// check operator
	// program, certificate, finding owner
	if finding.SubmitterAddress != msg.OperatorAddress && !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) && finding.SubmitterAddress != program.AdminAddress {
		return nil, types.ErrFindingSubmitterInvalid
	}

	k.DeleteFidFromFidList(ctx, finding.ProgramId, finding.FindingId)
	k.DeleteFinding(ctx, finding.FindingId)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCloseFinding,
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
	if program.Status == types.ProgramStatusInactive {
		return nil, types.ErrProgramNotActive
	}

	// close，finding owner can release
	// paid, program admin can release
	switch finding.Status {
	case types.FindingStatusClosed:
		if finding.SubmitterAddress != msg.OperatorAddress {
			return nil, types.ErrFindingSubmitterInvalid
		}
	case types.FindingStatusPaid:
		if program.AdminAddress != msg.OperatorAddress {
			return nil, types.ErrProgramCreatorInvalid
		}
	}

	// check hash
	hash := sha256.Sum256([]byte(msg.Description + msg.ProofOfConcept + finding.SubmitterAddress))
	if finding.FindingHash != hex.EncodeToString(hash[:]) {
		return nil, types.ErrProgramCreatorInvalid
	}

	finding.Description = msg.Description
	finding.ProofOfConcept = msg.ProofOfConcept
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
