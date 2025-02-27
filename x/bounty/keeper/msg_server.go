package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	createTime := ctx.BlockHeader().Time
	program, err := types.NewProgram(msg.ProgramId, msg.Name, msg.Detail, operatorAddr, types.ProgramStatusInactive, createTime)
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
		return nil, types.ErrProgramNotExists
	}

	// check the status.
	// inactive: program admin, bounty certificate
	// active: bounty certificate
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
	if err = k.Keeper.ActivateProgram(ctx, msg.ProgramId, operatorAddr); err != nil {
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

	operatorAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
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

	createTime := ctx.BlockHeader().Time
	finding, err := types.NewFinding(msg.ProgramId, msg.FindingId, "", "", msg.FindingHash, operatorAddr, createTime, msg.SeverityLevel)
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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgSubmitFindingResponse{}, nil
}

func (k msgServer) EditFinding(goCtx context.Context, msg *types.MsgEditFinding) (*types.MsgEditFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	finding, found := k.GetFinding(ctx, msg.FindingId)
	if !found {
		return nil, types.ErrFindingNotExists
	}
	// check program
	program, isExist := k.GetProgram(ctx, finding.ProgramId)
	if !isExist {
		return nil, types.ErrProgramNotExists
	}
	if program.Status != types.ProgramStatusActive {
		return nil, types.ErrProgramNotActive
	}

	// program admin edit paymentHash
	if len(msg.PaymentHash) > 0 {
		// check status
		if finding.Status != types.FindingStatusConfirmed {
			return nil, types.ErrFindingStatusInvalid
		}
		// check operator is program admin
		if program.AdminAddress != msg.OperatorAddress {
			return nil, types.ErrFindingOperatorNotAllowed
		}
		finding.PaymentHash = msg.PaymentHash

		k.SetFinding(ctx, finding)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeEditFindingPaymentHash,
				sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
				sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
			),
		})
		return &types.MsgEditFindingResponse{}, nil
	}

	// whitehat edit finding
	//  StatusSubmitted and StatusActive can be edited
	if finding.Status != types.FindingStatusSubmitted && finding.Status != types.FindingStatusActive {
		return nil, types.ErrFindingStatusInvalid
	}

	// check operator is whitehat
	if finding.SubmitterAddress != msg.OperatorAddress {
		return nil, types.ErrFindingOperatorNotAllowed
	}
	if len(msg.FindingHash) > 0 {
		finding.FindingHash = msg.FindingHash
	}
	if msg.SeverityLevel != types.Unspecified {
		finding.SeverityLevel = msg.SeverityLevel
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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgEditFindingResponse{}, nil
}

func (k msgServer) ActivateFinding(goCtx context.Context, msg *types.MsgActivateFinding) (*types.MsgActivateFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.AccAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	// get finding
	finding, found := k.GetFinding(ctx, msg.FindingId)
	if !found {
		return nil, types.ErrFindingNotExists
	}
	// only StatusSubmitted can activate
	if finding.Status != types.FindingStatusSubmitted {
		return nil, types.ErrFindingStatusInvalid
	}

	// get program
	program, isExist := k.GetProgram(ctx, finding.ProgramId)
	if !isExist {
		return nil, types.ErrProgramNotExists
	}
	if program.Status != types.ProgramStatusActive {
		return nil, types.ErrProgramNotActive
	}

	// program admin and bounty certificate can activate finding
	if program.AdminAddress != msg.OperatorAddress && !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) {
		return nil, types.ErrFindingOperatorNotAllowed
	}

	finding.Status = types.FindingStatusActive
	k.SetFinding(ctx, finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeActivateFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgActivateFindingResponse{}, nil
}

func (k msgServer) ConfirmFinding(goCtx context.Context, msg *types.MsgConfirmFinding) (*types.MsgConfirmFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	finding, err := k.Keeper.ConfirmFinding(ctx, msg)
	if err != nil {
		return nil, err
	}

	finding.Status = types.FindingStatusConfirmed
	k.SetFinding(ctx, finding)

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

func (k msgServer) ConfirmFindingPaid(goCtx context.Context, msg *types.MsgConfirmFindingPaid) (*types.MsgConfirmFindingPaidResponse, error) {
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
	if finding.Status != types.FindingStatusConfirmed {
		return nil, types.ErrFindingStatusInvalid
	}

	// check operator: finding owner, certificate
	if finding.SubmitterAddress != msg.OperatorAddress && !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) {
		return nil, types.ErrFindingOperatorNotAllowed
	}

	finding.Status = types.FindingStatusPaid
	k.SetFinding(ctx, finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeConfirmFindingPaid,
			sdk.NewAttribute(types.AttributeKeyFindingID, msg.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgConfirmFindingPaidResponse{}, nil
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
	if finding.SubmitterAddress != msg.OperatorAddress && !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) && program.AdminAddress != msg.OperatorAddress {
		return nil, types.ErrFindingOperatorNotAllowed
	}
	finding.Status = types.FindingStatusClosed
	k.SetFinding(ctx, finding)

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

func (k msgServer) PublishFinding(goCtx context.Context, msg *types.MsgPublishFinding) (*types.MsgPublishFindingResponse, error) {
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
			return nil, types.ErrFindingOperatorNotAllowed
		}
	case types.FindingStatusPaid:
		if program.AdminAddress != msg.OperatorAddress {
			return nil, types.ErrProgramOperatorNotAllowed
		}
	default:
		return nil, types.ErrFindingStatusInvalid
	}

	// check hash
	hash := sha256.Sum256([]byte(msg.Description + msg.ProofOfConcept + finding.SubmitterAddress))
	if finding.FindingHash != hex.EncodeToString(hash[:]) {
		return nil, types.ErrFindingHashInvalid
	}
	finding.Title = msg.Title
	finding.Detail = msg.Detail
	finding.Description = msg.Description
	finding.ProofOfConcept = msg.ProofOfConcept
	k.SetFinding(ctx, finding)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePublishFinding,
			sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, program.ProgramId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OperatorAddress),
		),
	})

	return &types.MsgPublishFindingResponse{}, nil
}

func (k msgServer) CreateTheorem(goCtx context.Context, msg *types.MsgCreateTheorem) (*types.MsgCreateTheoremResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Title == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "theorem title cannot be empty")
	}
	if msg.Description == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "theorem description cannot be empty")

	}
	if msg.Code == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "theorem code cannot be empty")
	}
	if len(msg.Title+msg.Description+msg.Code) > 5000 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "theorem description too large")
	}

	proposer, err := k.authKeeper.AddressCodec().StringToBytes(msg.GetProposer())
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid proposer address: %s", err)
	}

	initialGrant := msg.GetInitialGrant()
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get theorem parameters: %w", err)
	}

	if err := k.validateMinGrant(ctx, params, initialGrant); err != nil {
		return nil, err
	}
	if err := k.validateDepositDenom(ctx, params, initialGrant); err != nil {
		return nil, err
	}

	submitTime := ctx.BlockHeader().Time
	theorem, err := k.Keeper.CreateTheorem(ctx, proposer, msg.Title, msg.Description, msg.Code, submitTime, submitTime.Add(*params.TheoremMaxGrantPeriod), *params.TheoremMaxProofPeriod)
	if err != nil {
		return nil, err
	}

	if err = k.Keeper.AddGrant(ctx, theorem.Id, proposer, msg.GetInitialGrant()); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeCreateTheorem,
			sdk.NewAttribute(types.AttributeKeyTheoremProofPeriodStart, fmt.Sprintf("%d", theorem.Id)),
			sdk.NewAttribute(types.AttributeKeyTheoremProposer, msg.GetProposer()),
		),
	)

	return &types.MsgCreateTheoremResponse{
		TheoremId: theorem.Id,
	}, nil
}

func (k msgServer) SubmitProofHash(goCtx context.Context, msg *types.MsgSubmitProofHash) (*types.MsgSubmitProofHashResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// msg check
	if msg.ProofHash == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proof hash cannot be empty")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get theorem parameters: %w", err)
	}
	if err := k.validateMinDeposit(ctx, params, msg.Deposit); err != nil {
		return nil, err
	}
	if err := k.validateDepositDenom(ctx, params, msg.Deposit); err != nil {
		return nil, err
	}

	proposer, err := k.authKeeper.AddressCodec().StringToBytes(msg.GetProver())
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid proposer address: %s", err)
	}

	proof, err := k.Keeper.SubmitProofHash(goCtx, msg.TheoremId, msg.ProofHash, msg.Prover, msg.Deposit)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.AddDeposit(ctx, proof.Id, proposer, msg.Deposit); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeSubmitProofHash,
			sdk.NewAttribute(types.AttributeKeyProofHashLockPeriodStart, proof.Id),
			sdk.NewAttribute(types.AttributeKeyTheoremProposer, msg.GetProver()),
		),
	)

	return &types.MsgSubmitProofHashResponse{}, nil
}

func (k msgServer) SubmitProofDetail(goCtx context.Context, msg *types.MsgSubmitProofDetail) (*types.MsgSubmitProofHashResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if proof exists
	proof, err := k.Proofs.Get(ctx, msg.ProofId)
	if err != nil {
		return nil, err
	}

	// hash check
	if msg.Detail == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proof hash cannot be empty")
	}
	_, err = k.Proofs.Get(ctx, msg.ProofId)
	if err != nil {
		return nil, err
	}

	hash := k.Keeper.GetProofHash(proof.TheoremId, msg.GetProver(), msg.Detail)

	if proof.Id != hash {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proof hash inconsistent")
	}

	err = k.Keeper.SubmitProofDetail(ctx, msg.ProofId, msg.Detail)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeSubmitProofDetail,
			sdk.NewAttribute(types.AttributeKeyProofHashDetailPeriodStart, proof.Id),
			sdk.NewAttribute(types.AttributeKeyTheoremProposer, msg.GetProver()),
		),
	)

	return &types.MsgSubmitProofHashResponse{}, nil
}

func (k msgServer) Grant(goCtx context.Context, msg *types.MsgGrant) (*types.MsgGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	grantor, err := k.authKeeper.AddressCodec().StringToBytes(msg.GetGrantor())
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid proposer address: %s", err)
	}

	err = k.Keeper.AddGrant(ctx, msg.TheoremId, grantor, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgGrantResponse{}, nil
}
