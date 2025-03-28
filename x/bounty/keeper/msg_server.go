package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bounty MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateProgram creates a new bounty program
func (k msgServer) CreateProgram(goCtx context.Context, msg *types.MsgCreateProgram) (*types.MsgCreateProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"programId": msg.ProgramId,
		"name":      msg.Name,
		"detail":    msg.Detail,
	}); err != nil {
		return nil, err
	}

	operatorAddr, err := k.validateAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	exist, err := k.Programs.Has(ctx, msg.ProgramId)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, types.ErrProgramAlreadyExists
	}

	createTime := ctx.BlockHeader().Time
	program, err := types.NewProgram(msg.ProgramId, msg.Name, msg.Detail, operatorAddr, types.ProgramStatusInactive, createTime)
	if err != nil {
		return nil, err
	}

	if err = k.Programs.Set(ctx, program.ProgramId, program); err != nil {
		return nil, err
	}

	// emit event
	k.emitProgramEvent(ctx, types.EventTypeCreateProgram, msg.ProgramId, msg.OperatorAddress)

	return &types.MsgCreateProgramResponse{}, nil
}

// EditProgram modifies an existing bounty program
func (k msgServer) EditProgram(goCtx context.Context, msg *types.MsgEditProgram) (*types.MsgEditProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"programId": msg.ProgramId,
	}); err != nil {
		return nil, err
	}

	operatorAddr, err := k.validateAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	program, err := k.Programs.Get(goCtx, msg.ProgramId)
	if err != nil {
		return nil, err
	}

	// check permissions based on program status
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

	// update program fields
	if len(msg.Name) > 0 {
		program.Name = msg.Name
	}
	if len(msg.Detail) > 0 {
		program.Detail = msg.Detail
	}

	if err = k.Programs.Set(ctx, program.ProgramId, program); err != nil {
		return nil, err
	}

	// emit event
	k.emitProgramEvent(ctx, types.EventTypeEditProgram, msg.ProgramId, msg.OperatorAddress)

	return &types.MsgEditProgramResponse{}, nil
}

// ActivateProgram changes a program's status to active
func (k msgServer) ActivateProgram(goCtx context.Context, msg *types.MsgActivateProgram) (*types.MsgActivateProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"programId": msg.ProgramId,
	}); err != nil {
		return nil, err
	}

	operatorAddr, err := k.validateAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	// get program
	program, err := k.Programs.Get(ctx, msg.ProgramId)
	if err != nil {
		return nil, err
	}

	// check if the program is already active
	if program.Status == types.ProgramStatusActive {
		return nil, types.ErrProgramAlreadyActive
	}

	// check the permissions. Only the bounty cert address can operate.
	if !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) {
		return nil, types.ErrProgramOperatorNotAllowed
	}

	// update program status
	program.Status = types.ProgramStatusActive
	if err = k.Programs.Set(ctx, program.ProgramId, program); err != nil {
		return nil, err
	}

	// emit event
	k.emitProgramEvent(ctx, types.EventTypeActivateProgram, msg.ProgramId, msg.OperatorAddress)

	return &types.MsgActivateProgramResponse{}, nil
}

// CloseProgram closes a bounty program
func (k msgServer) CloseProgram(goCtx context.Context, msg *types.MsgCloseProgram) (*types.MsgCloseProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"programId": msg.ProgramId,
	}); err != nil {
		return nil, err
	}

	operatorAddr, err := k.validateAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	// get program
	program, err := k.Programs.Get(ctx, msg.ProgramId)
	if err != nil {
		return nil, err
	}

	// check if the program is already closed
	if program.Status == types.ProgramStatusClosed {
		return nil, types.ErrProgramAlreadyClosed
	}

	// the program cannot be closed if there are findings in certain states
	// there are 3 finding states: FindingStatusSubmitted FindingStatusActive FindingStatusConfirmed
	fidsList, err := k.getProgramFindings(ctx, msg.ProgramId)
	if err != nil {
		return nil, err
	}
	for _, fid := range fidsList {
		finding, err := k.Findings.Get(ctx, fid)
		if err != nil {
			return nil, err
		}
		if finding.Status == types.FindingStatusSubmitted ||
			finding.Status == types.FindingStatusActive ||
			finding.Status == types.FindingStatusConfirmed {
			return nil, types.ErrProgramCloseNotAllowed
		}
	}

	// check the permissions. Only the admin of the program or bounty cert address can operate.
	if program.AdminAddress != msg.OperatorAddress && !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) {
		return nil, types.ErrProgramOperatorNotAllowed
	}

	// close the program and update its status
	program.Status = types.ProgramStatusClosed
	if err = k.Programs.Set(ctx, program.ProgramId, program); err != nil {
		return nil, err
	}

	// emit event
	k.emitProgramEvent(ctx, types.EventTypeCloseProgram, msg.ProgramId, msg.OperatorAddress)

	return &types.MsgCloseProgramResponse{}, nil
}

// SubmitFinding creates a new security finding for a program
func (k msgServer) SubmitFinding(goCtx context.Context, msg *types.MsgSubmitFinding) (*types.MsgSubmitFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"programId":   msg.ProgramId,
		"findingId":   msg.FindingId,
		"findingHash": msg.FindingHash,
	}); err != nil {
		return nil, err
	}

	if !types.ValidFindingSeverityLevel(msg.SeverityLevel) {
		return nil, errors.Wrap(types.ErrFindingSeverityLevelInvalid, msg.SeverityLevel.String())
	}

	operatorAddr, err := k.validateAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	_, err = k.validateProgramStatus(ctx, msg.ProgramId, types.ProgramStatusActive)
	if err != nil {
		return nil, err
	}

	// check if finding already exists - corrected logic
	exist, err := k.Findings.Has(ctx, msg.FindingId)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, types.ErrFindingAlreadyExists
	}

	createTime := ctx.BlockHeader().Time
	finding, err := types.NewFinding(msg.ProgramId, msg.FindingId, "", "", msg.FindingHash, operatorAddr, createTime, msg.SeverityLevel)
	if err != nil {
		return nil, err
	}

	if err = k.ProgramFindings.Set(ctx, collections.Join(msg.ProgramId, msg.FindingId)); err != nil {
		return nil, err
	}

	if err = k.Findings.Set(ctx, finding.FindingId, finding); err != nil {
		return nil, err
	}

	// emit event
	k.emitFindingEvent(ctx, types.EventTypeSubmitFinding, finding, msg.OperatorAddress)

	return &types.MsgSubmitFindingResponse{}, nil
}

// EditFinding modifies an existing security finding
func (k msgServer) EditFinding(goCtx context.Context, msg *types.MsgEditFinding) (*types.MsgEditFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"findingId": msg.FindingId,
	}); err != nil {
		return nil, err
	}

	if !types.ValidFindingSeverityLevel(msg.SeverityLevel) {
		return nil, errors.Wrap(types.ErrFindingSeverityLevelInvalid, msg.SeverityLevel.String())
	}

	// validate operator address
	if _, err := k.validateAddress(msg.OperatorAddress); err != nil {
		return nil, err
	}

	findingPtr, err := k.validateFindingStatus(ctx, msg.FindingId, types.FindingStatusSubmitted, types.FindingStatusActive)
	if err != nil {
		return nil, err
	}
	finding := *findingPtr

	program, err := k.validateProgramStatus(ctx, finding.ProgramId, types.ProgramStatusActive)
	if err != nil {
		return nil, err
	}

	// program admin edit paymentHash
	if len(msg.PaymentHash) > 0 {
		if finding.Status != types.FindingStatusConfirmed {
			return nil, types.ErrFindingStatusInvalid
		}
		if program.AdminAddress != msg.OperatorAddress {
			return nil, types.ErrFindingOperatorNotAllowed
		}
		finding.PaymentHash = msg.PaymentHash

		if err = k.Findings.Set(ctx, finding.FindingId, finding); err != nil {
			return nil, err
		}

		// emit event for payment hash update
		k.emitFindingEvent(ctx, types.EventTypeEditFindingPaymentHash, finding, msg.OperatorAddress)

		return &types.MsgEditFindingResponse{}, nil
	}

	// whitehat edit finding
	if finding.SubmitterAddress != msg.OperatorAddress {
		return nil, types.ErrFindingOperatorNotAllowed
	}
	if len(msg.FindingHash) > 0 {
		finding.FindingHash = msg.FindingHash
	}
	if msg.SeverityLevel != types.Unspecified {
		finding.SeverityLevel = msg.SeverityLevel
	}

	if err = k.Findings.Set(ctx, finding.FindingId, finding); err != nil {
		return nil, err
	}

	// emit event for general finding edit
	k.emitFindingEvent(ctx, types.EventTypeEditFinding, finding, msg.OperatorAddress)

	return &types.MsgEditFindingResponse{}, nil
}

// ActivateFinding changes a finding's status to active
// Only program admins and bounty certificate holders can activate findings
func (k msgServer) ActivateFinding(goCtx context.Context, msg *types.MsgActivateFinding) (*types.MsgActivateFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"findingId": msg.FindingId,
	}); err != nil {
		return nil, err
	}

	operatorAddr, err := k.validateAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	// get finding
	finding, err := k.Findings.Get(ctx, msg.FindingId)
	if err != nil {
		return nil, err
	}

	// only StatusSubmitted can activate
	if finding.Status != types.FindingStatusSubmitted {
		return nil, types.ErrFindingStatusInvalid
	}

	program, err := k.Programs.Get(ctx, finding.ProgramId)
	if err != nil {
		return nil, err
	}
	if program.Status != types.ProgramStatusActive {
		return nil, types.ErrProgramNotActive
	}

	// check permissions
	if program.AdminAddress != msg.OperatorAddress && !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) {
		return nil, types.ErrFindingOperatorNotAllowed
	}

	// update finding status
	finding.Status = types.FindingStatusActive
	if err = k.Findings.Set(ctx, finding.FindingId, finding); err != nil {
		return nil, err
	}

	// emit event
	k.emitFindingEvent(ctx, types.EventTypeActivateFinding, finding, msg.OperatorAddress)

	return &types.MsgActivateFindingResponse{}, nil
}

// ConfirmFinding confirms a security finding with the given fingerprint
func (k msgServer) ConfirmFinding(goCtx context.Context, msg *types.MsgConfirmFinding) (*types.MsgConfirmFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"findingId":   msg.FindingId,
		"fingerprint": msg.Fingerprint,
	}); err != nil {
		return nil, err
	}

	// validate operator address
	if _, err := k.validateAddress(msg.OperatorAddress); err != nil {
		return nil, err
	}

	// get finding
	finding, err := k.Findings.Get(ctx, msg.FindingId)
	if err != nil {
		return nil, err
	}

	// only StatusActive can be confirmed
	if finding.Status != types.FindingStatusActive {
		return nil, types.ErrFindingStatusInvalid
	}

	// get program
	program, err := k.Programs.Get(ctx, finding.ProgramId)
	if err != nil {
		return nil, err
	}
	if program.Status != types.ProgramStatusActive {
		return nil, types.ErrProgramNotActive
	}

	// only program admin can confirm finding
	if program.AdminAddress != msg.OperatorAddress {
		return nil, types.ErrProgramOperatorNotAllowed
	}

	// fingerprint comparison
	fingerprintHash := k.GetFindingFingerprintHash(&finding)
	if msg.Fingerprint != fingerprintHash {
		return nil, types.ErrFindingHashInvalid
	}

	// update finding status
	finding.Status = types.FindingStatusConfirmed
	if err = k.Findings.Set(ctx, finding.FindingId, finding); err != nil {
		return nil, err
	}

	// emit event
	k.emitFindingEvent(ctx, types.EventTypeConfirmFinding, finding, msg.OperatorAddress)

	return &types.MsgConfirmFindingResponse{}, nil
}

// ConfirmFindingPaid marks a finding as paid
// Can be called by the finding submitter or a bounty admin
func (k msgServer) ConfirmFindingPaid(goCtx context.Context, msg *types.MsgConfirmFindingPaid) (*types.MsgConfirmFindingPaidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"findingId": msg.FindingId,
	}); err != nil {
		return nil, err
	}

	operatorAddr, err := k.validateAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	findingPtr, err := k.validateFindingStatus(ctx, msg.FindingId, types.FindingStatusConfirmed)
	if err != nil {
		return nil, err
	}
	finding := *findingPtr

	// check operator permissions: finding owner or bounty admin
	if finding.SubmitterAddress != msg.OperatorAddress && !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) {
		return nil, types.ErrFindingOperatorNotAllowed
	}

	finding.Status = types.FindingStatusPaid
	if err = k.Findings.Set(ctx, finding.FindingId, finding); err != nil {
		return nil, err
	}

	// emit event
	k.emitFindingEvent(ctx, types.EventTypeConfirmFindingPaid, finding, msg.OperatorAddress)

	return &types.MsgConfirmFindingPaidResponse{}, nil
}

// CloseFinding closes a security finding
// Only available for findings in submitted or active state
func (k msgServer) CloseFinding(goCtx context.Context, msg *types.MsgCloseFinding) (*types.MsgCloseFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"findingId": msg.FindingId,
	}); err != nil {
		return nil, err
	}

	operatorAddr, err := k.validateAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	// get finding
	finding, err := k.Findings.Get(goCtx, msg.FindingId)
	if err != nil {
		return nil, err
	}

	// check finding status: StatusSubmitted and StatusActive can be closed
	if finding.Status != types.FindingStatusSubmitted && finding.Status != types.FindingStatusActive {
		return nil, types.ErrFindingStatusInvalid
	}
	// get program
	program, err := k.Programs.Get(ctx, finding.ProgramId)
	if err != nil {
		return nil, err
	}

	// check operator
	// program, certificate, finding owner
	if finding.SubmitterAddress != msg.OperatorAddress && !k.certKeeper.IsBountyAdmin(ctx, operatorAddr) && program.AdminAddress != msg.OperatorAddress {
		return nil, types.ErrFindingOperatorNotAllowed
	}
	finding.Status = types.FindingStatusClosed
	if err = k.Findings.Set(ctx, finding.FindingId, finding); err != nil {
		return nil, err
	}

	// emit event
	k.emitFindingEvent(ctx, types.EventTypeCloseFinding, finding, msg.OperatorAddress)

	return &types.MsgCloseFindingResponse{}, nil
}

// PublishFinding publishes a security finding with full details
func (k msgServer) PublishFinding(goCtx context.Context, msg *types.MsgPublishFinding) (*types.MsgPublishFindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"findingId":      msg.FindingId,
		"description":    msg.Description,
		"proofOfConcept": msg.ProofOfConcept,
	}); err != nil {
		return nil, err
	}

	// validate operator address
	if _, err := k.validateAddress(msg.OperatorAddress); err != nil {
		return nil, err
	}

	// get finding
	finding, err := k.Findings.Get(ctx, msg.FindingId)
	if err != nil {
		return nil, err
	}

	// get program
	program, err := k.Programs.Get(goCtx, finding.ProgramId)
	if err != nil {
		return nil, err
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

	// verify hash
	hash := sha256.Sum256([]byte(msg.Description + msg.ProofOfConcept + finding.SubmitterAddress))
	if finding.FindingHash != hex.EncodeToString(hash[:]) {
		return nil, types.ErrFindingHashInvalid
	}

	// update finding details
	finding.Title = msg.Title
	finding.Detail = msg.Detail
	finding.Description = msg.Description
	finding.ProofOfConcept = msg.ProofOfConcept
	if err = k.Findings.Set(ctx, finding.FindingId, finding); err != nil {
		return nil, err
	}

	// emit event
	k.emitFindingEvent(ctx, types.EventTypePublishFinding, finding, msg.OperatorAddress)

	return &types.MsgPublishFindingResponse{}, nil
}

func (k msgServer) CreateTheorem(goCtx context.Context, msg *types.MsgCreateTheorem) (*types.MsgCreateTheoremResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate basic message fields
	if err := validateMsgFields(map[string]string{
		"title":       msg.Title,
		"description": msg.Description,
		"code":        msg.Code,
	}); err != nil {
		return nil, err
	}

	proposer, err := k.validateAddress(msg.GetProposer())
	if err != nil {
		return nil, err
	}

	initialGrant := msg.GetInitialGrant()

	// validate grant funds
	params, err := k.ValidateFunds(ctx, initialGrant, "grant")
	if err != nil {
		return nil, err
	}

	submitTime := ctx.BlockHeader().Time
	endTime := submitTime.Add(*params.TheoremMaxProofPeriod)
	theoremID, err := k.TheoremID.Next(ctx)
	if err != nil {
		return nil, err
	}

	theorem, err := types.NewTheorem(theoremID, proposer, msg.Title, msg.Description, msg.Code, submitTime, endTime)
	if err != nil {
		return nil, err
	}
	if err = k.Theorems.Set(ctx, theorem.Id, theorem); err != nil {
		return nil, err
	}
	if err = k.ActiveTheoremsQueue.Set(ctx, collections.Join(endTime, theoremID), theoremID); err != nil {
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

func (k msgServer) Grant(goCtx context.Context, msg *types.MsgGrant) (*types.MsgGrantResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	grantor, err := k.validateAddress(msg.GetGrantor())
	if err != nil {
		return nil, err
	}

	// validate grant funds
	_, err = k.ValidateFunds(ctx, msg.Amount, "grant")
	if err != nil {
		return nil, err
	}

	err = k.Keeper.AddGrant(ctx, msg.TheoremId, grantor, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgGrantResponse{}, nil
}

func (k msgServer) SubmitProofHash(goCtx context.Context, msg *types.MsgSubmitProofHash) (*types.MsgSubmitProofHashResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate proofHash length (SHA-256 hash as hex string should be 64 characters)
	if len(msg.ProofHash) != 64 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proofHash must be a 64-character SHA-256 hash")
	}

	// validate proofHash is a valid hex string
	if _, err := hex.DecodeString(msg.ProofHash); err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proofHash must be a valid hex string")
	}

	proposer, err := k.validateAddress(msg.GetProver())
	if err != nil {
		return nil, err
	}

	// check if proof already exists
	exists, err := k.Proofs.Has(ctx, msg.ProofHash)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, types.ErrProofAlreadyExists
	}
	// check if theorem exists
	theorem, err := k.Theorems.Get(ctx, msg.TheoremId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "theorem %d doesn't exist", msg.TheoremId)
		}
		return nil, err
	}

	// check theorem is still proof able
	if theorem.Status != types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD {
		return nil, types.ErrTheoremStatusInvalid
	}

	// validate deposit funds
	params, err := k.ValidateFunds(ctx, msg.Deposit, "deposit")
	if err != nil {
		return nil, err
	}

	submitTime := ctx.BlockHeader().Time
	endTime := submitTime.Add(*params.ProofMaxLockPeriod)
	proof, err := types.NewProof(msg.TheoremId, msg.ProofHash, msg.Prover, submitTime, endTime, msg.Deposit)
	if err != nil {
		return nil, err
	}

	if err = k.Proofs.Set(ctx, proof.Id, proof); err != nil {
		return nil, err
	}
	if err = k.TheoremProof.Set(ctx, msg.TheoremId, msg.ProofHash); err != nil {
		return nil, err
	}
	if err = k.ActiveProofsQueue.Set(ctx, collections.Join(endTime, msg.ProofHash), proof); err != nil {
		return nil, err
	}
	if err = k.ProofsByTheorem.Set(ctx, collections.Join(msg.TheoremId, msg.ProofHash), []byte{}); err != nil {
		return nil, err
	}
	if err := k.Keeper.AddDeposit(ctx, proof.Id, proposer, msg.Deposit); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeSubmitProofHash,
			sdk.NewAttribute(types.AttributeKeyProofID, proof.Id),
			sdk.NewAttribute(types.AttributeKeyTheoremProposer, msg.GetProver()),
		),
	)

	return &types.MsgSubmitProofHashResponse{}, nil
}

// SubmitProofDetail submits a proof detail for a theorem
func (k msgServer) SubmitProofDetail(goCtx context.Context, msg *types.MsgSubmitProofDetail) (*types.MsgSubmitProofDetailResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic message fields
	if err := validateMsgFields(map[string]string{
		"detail": msg.Detail,
	}); err != nil {
		return nil, err
	}

	// Get and validate the proof
	proof, err := k.Proofs.Get(ctx, msg.ProofId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "proof %d doesn't exist", msg.ProofId)
		}
		return nil, err
	}

	// Check proof status - must be in hash lock period
	if proof.Status != types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD {
		return nil, types.ErrProofStatusInvalid
	}

	// Verify the hash matches
	hash := k.GetProofHash(proof.TheoremId, msg.GetProver(), msg.Detail)
	if proof.Id != hash {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proof hash inconsistent")
	}

	// Remove from active proofs queue since we're transitioning out of hash lock period
	// The proof will no longer need to be tracked for potential expiration
	if err = k.ActiveProofsQueue.Remove(ctx, collections.Join(*proof.EndTime, proof.Id)); err != nil {
		return nil, err
	}

	// Update proof
	proof.Detail = msg.Detail
	proof.Status = types.ProofStatus_PROOF_STATUS_HASH_DETAIL_PERIOD

	// Save updated proof
	if err = k.Proofs.Set(ctx, proof.Id, proof); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProofDetail,
			sdk.NewAttribute(types.AttributeKeyProofID, msg.ProofId),
			sdk.NewAttribute(types.AttributeKeyTheoremProposer, msg.GetProver()),
		),
	)

	return &types.MsgSubmitProofDetailResponse{}, nil
}

// SubmitProofVerification submits a proof verification for a theorem
func (k msgServer) SubmitProofVerification(goCtx context.Context, msg *types.MsgSubmitProofVerification) (*types.MsgSubmitProofVerificationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate checker address and authority
	checkerAddr, err := k.validateAddress(msg.Checker)
	if err != nil {
		return nil, err
	}
	if !k.certKeeper.IsBountyAdmin(ctx, checkerAddr) {
		return nil, types.ErrProofOperatorNotAllowed
	}

	// Validate proof status
	if !isValidProofStatus(msg.Status) {
		return nil, types.ErrProofStatusInvalid
	}

	// Get and validate proof
	proof, err := k.validateProofStatus(ctx, msg.ProofId, types.ProofStatus_PROOF_STATUS_HASH_DETAIL_PERIOD)
	if err != nil {
		return nil, err
	}

	// Get prover address
	proverAddr, err := k.validateAddress(proof.Prover)
	if err != nil {
		return nil, err
	}

	// Handle proof verification based on status
	if err = k.handleProofVerification(ctx, msg.Status, *proof, checkerAddr, proverAddr); err != nil {
		return nil, err
	}

	// Remove theorem proof mapping
	if err = k.TheoremProof.Remove(ctx, proof.TheoremId); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProofVerification,
			sdk.NewAttribute(types.AttributeKeyProofID, msg.ProofId),
			sdk.NewAttribute(types.AttributeKeyChecker, msg.Checker),
			sdk.NewAttribute(types.AttributeKeyProofStatus, msg.Status.String()),
		),
	)

	return &types.MsgSubmitProofVerificationResponse{}, nil
}

// WithdrawReward withdraws a reward from a bounty program
func (k msgServer) WithdrawReward(goCtx context.Context, msg *types.MsgWithdrawReward) (*types.MsgWithdrawRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := k.validateAddress(msg.Address)
	if err != nil {
		return nil, err
	}

	reward, err := k.Rewards.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	finalRewards, _ := reward.Reward.TruncateDecimal()
	if !finalRewards.IsZero() {
		if err = k.Rewards.Remove(ctx, addr); err != nil {
			return nil, err
		}

		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, finalRewards); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawReward,
			sdk.NewAttribute(types.AttributeKeyReward, finalRewards.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
		),
	)

	return &types.MsgWithdrawRewardResponse{}, nil
}

// validateAddress validates the address string and returns the decoded address
func (k msgServer) validateAddress(address string) (sdk.AccAddress, error) {
	return k.authKeeper.AddressCodec().StringToBytes(address)
}

// validateProgramStatus checks if the program exists and is in the expected status
func (k msgServer) validateProgramStatus(ctx sdk.Context, programID string, expectedStatus types.ProgramStatus) (*types.Program, error) {
	program, err := k.Programs.Get(ctx, programID)
	if err != nil {
		return nil, err
	}
	if program.Status != expectedStatus {
		return nil, types.ErrProgramNotActive
	}
	return &program, nil
}

// validateFindingStatus checks if the finding exists and is in any of the expected statuses
func (k msgServer) validateFindingStatus(ctx sdk.Context, findingID string, expectedStatuses ...types.FindingStatus) (*types.Finding, error) {
	finding, err := k.Findings.Get(ctx, findingID)
	if err != nil {
		return nil, err
	}

	for _, status := range expectedStatuses {
		if finding.Status == status {
			return &finding, nil
		}
	}

	return nil, types.ErrFindingStatusInvalid
}

// validateProofStatus checks if the proof exists and is in the expected status
func (k msgServer) validateProofStatus(ctx sdk.Context, proofID string, expectedStatus types.ProofStatus) (*types.Proof, error) {
	proof, err := k.Proofs.Get(ctx, proofID)
	if err != nil {
		return nil, err
	}
	if proof.Status != expectedStatus {
		return nil, types.ErrProofStatusInvalid
	}
	return &proof, nil
}

// isValidProofStatus checks if the given proof status is valid
func isValidProofStatus(status types.ProofStatus) bool {
	return status == types.ProofStatus_PROOF_STATUS_PASSED ||
		status == types.ProofStatus_PROOF_STATUS_FAILED
}

// handleProofVerification processes proof verification based on status
// For passed proofs, updates theorem status and distributes rewards
// For failed proofs, deletes the proof and removes theorem-proof mapping
func (k msgServer) handleProofVerification(
	ctx sdk.Context,
	status types.ProofStatus,
	proof types.Proof,
	checkerAddr, proverAddr sdk.AccAddress,
) error {
	proof.Status = status

	switch status {
	case types.ProofStatus_PROOF_STATUS_PASSED:
		return k.handlePassedProof(ctx, proof, checkerAddr, proverAddr)
	case types.ProofStatus_PROOF_STATUS_FAILED:
		return k.handleFailedProof(ctx, proof, proverAddr)
	default:
		return types.ErrProofStatusInvalid
	}
}

// handlePassedProof processes a proof that passed verification
// Updates proof and theorem status, and distributes rewards
func (k msgServer) handlePassedProof(
	ctx sdk.Context,
	proof types.Proof,
	checkerAddr, proverAddr sdk.AccAddress,
) error {
	// Update proof status
	if err := k.Proofs.Set(ctx, proof.Id, proof); err != nil {
		return err
	}

	// Update theorem status
	theorem, err := k.Theorems.Get(ctx, proof.TheoremId)
	if err != nil {
		return err
	}

	theorem.Status = types.TheoremStatus_THEOREM_STATUS_PASSED
	if err := k.Theorems.Set(ctx, theorem.Id, theorem); err != nil {
		return err
	}

	// Remove from active theorems queue
	if err := k.ActiveTheoremsQueue.Remove(ctx, collections.Join(*theorem.EndTime, theorem.Id)); err != nil {
		return err
	}

	// Distribute grants
	return k.DistributionGrants(ctx, proof.TheoremId, checkerAddr, proverAddr)
}

// handleFailedProof processes a proof that failed verification
// Deletes the proof and cleans up related mappings
func (k msgServer) handleFailedProof(
	ctx sdk.Context,
	proof types.Proof,
	proverAddr sdk.AccAddress,
) error {
	if err := k.DeleteProof(ctx, proof.Id); err != nil {
		return err
	}

	if err := k.TheoremProof.Remove(ctx, proof.TheoremId); err != nil {
		return err
	}

	return k.Deposits.Remove(ctx, collections.Join(proof.Id, proverAddr))
}

// validateMsgFields validates that required fields are not empty
func validateMsgFields(fields map[string]string) error {
	for fieldName, fieldValue := range fields {
		if len(fieldValue) == 0 {
			return errors.Wrap(sdkerrors.ErrInvalidRequest, "empty "+fieldName)
		}
	}
	return nil
}

// emitProgramEvent emits a standardized event for program-related operations
func (k msgServer) emitProgramEvent(ctx sdk.Context, eventType string, programID, operatorAddress string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeKeyProgramID, programID),
			sdk.NewAttribute(sdk.AttributeKeySender, operatorAddress),
		),
	)
}

// emitFindingEvent emits a standardized event for finding-related operations
func (k msgServer) emitFindingEvent(ctx sdk.Context, eventType string, finding types.Finding, operatorAddress string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeKeyFindingID, finding.FindingId),
			sdk.NewAttribute(types.AttributeKeyProgramID, finding.ProgramId),
			sdk.NewAttribute(sdk.AttributeKeySender, operatorAddress),
		),
	)
}

// getProgramFindings returns all finding IDs associated with a program
func (k msgServer) getProgramFindings(ctx context.Context, programID string) ([]string, error) {
	var findingIDs []string

	rng := collections.NewPrefixedPairRange[string, string](programID)
	err := k.ProgramFindings.Walk(ctx, rng, func(key collections.Pair[string, string]) (stop bool, err error) {
		if key.K1() == programID {
			findingIDs = append(findingIDs, key.K2())
		}
		return false, nil
	})
	if err != nil {
		return findingIDs, err
	}

	return findingIDs, nil
}
