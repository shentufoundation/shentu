package keeper

import (
	"context"
	"fmt"
	"strconv"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/client/cli"
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

	nextID, err := k.GetNextProgramID(ctx)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, msg.Deposit)
	if err != nil {
		return nil, err
	}

	if msg.SubmissionEndTime.Before(ctx.BlockTime()) {
		return nil, fmt.Errorf("submission end time is invalid")
	}

	program := types.Program{
		ProgramId:         nextID,
		CreatorAddress:    msg.CreatorAddress,
		SubmissionEndTime: msg.SubmissionEndTime,
		Description:       msg.Description,
		EncryptionKey:     msg.EncryptionKey,
		Deposit:           msg.Deposit,
		CommissionRate:    msg.CommissionRate,
		Active:            true,
	}

	k.SetProgram(ctx, program)

	k.SetNextProgramID(ctx, nextID+1)

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

	return &types.MsgCreateProgramResponse{ProgramId: nextID}, nil
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

	if !program.Active {
		return nil, types.ErrProgramInactive
	}

	findingID, err := k.GetNextFindingID(ctx)
	if err != nil {
		return nil, err
	}

	finding := types.Finding{
		FindingId:        findingID,
		Title:            msg.Title,
		FindingDesc:      msg.EncryptedDesc,
		ProgramId:        msg.ProgramId,
		SeverityLevel:    msg.SeverityLevel,
		FindingPoc:       msg.EncryptedPoc,
		SubmitterAddress: msg.SubmitterAddress,
		FindingStatus:    types.FindingStatusUnConfirmed,
	}

	err = k.AppendFidToFidList(ctx, msg.ProgramId, findingID)
	if err != nil {
		return nil, err
	}

	k.SetFinding(ctx, finding)
	k.SetNextFindingID(ctx, findingID+1)

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

	finding, err := k.hostProcess(ctx, msg.FindingId, msg.HostAddress, msg.EncryptedComment)
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

	finding, err := k.hostProcess(ctx, msg.FindingId, msg.HostAddress, msg.EncryptedComment)
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

func (k msgServer) hostProcess(ctx sdk.Context, fid uint64, hostAddr string, encryptedCommentAny *codectypes.Any) (*types.Finding, error) {

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
	if !program.Active {
		return nil, types.ErrProgramInactive
	}

	// only creator can update finding comment
	if program.CreatorAddress != hostAddr {
		return nil, types.ErrProgramCreatorInvalid
	}

	finding.FindingComment = encryptedCommentAny
	return &finding, nil
}

func (k msgServer) CancelFinding(goCtx context.Context, msg *types.MsgCancelFinding) (*types.MsgCancelFindingResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.SubmitterAddress)
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
	if finding.SubmitterAddress != msg.SubmitterAddress {
		return nil, types.ErrFindingSubmitterInvalid
	}

	// check status
	if finding.FindingStatus != types.FindingStatusUnConfirmed {
		return nil, types.ErrFindingStatusInvalid
	}

	k.DeleteFidFromFidList(ctx, finding.ProgramId, finding.FindingId)
	k.DeleteFinding(ctx, finding.FindingId)

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
	if !program.Active {
		return nil, types.ErrProgramInactive
	}

	// only creator can update finding comment
	if program.CreatorAddress != msg.HostAddress {
		return nil, types.ErrProgramCreatorInvalid
	}

	pubKey, err := cli.KeyAnyToPubKey(program.EncryptionKey)
	if err != nil {
		return nil, types.ErrProgramPubKey
	}

	if err = CheckPlainText(pubKey, msg, finding); err != nil {
		return nil, err
	}

	plainTextDesc := types.PlainTextDesc{
		FindingDesc: []byte(msg.Desc),
	}
	descAny, err := codectypes.NewAnyWithValue(&plainTextDesc)
	if err != nil {
		return nil, err
	}
	finding.FindingDesc = descAny

	plainTextPoc := types.PlainTextPoc{
		FindingPoc: []byte(msg.Poc),
	}
	pocAny, err := codectypes.NewAnyWithValue(&plainTextPoc)
	if err != nil {
		return nil, err
	}
	finding.FindingPoc = pocAny

	plainTextComment := types.PlainTextComment{
		FindingComment: []byte(msg.Comment),
	}
	commentAny, err := codectypes.NewAnyWithValue(&plainTextComment)
	if err != nil {
		return nil, err
	}
	finding.FindingComment = commentAny

	k.SetFinding(ctx, finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReleaseFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, strconv.FormatUint(finding.FindingId, 10)),
			sdk.NewAttribute(types.AttributeKeyProgramID, strconv.FormatUint(program.ProgramId, 10)),
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
