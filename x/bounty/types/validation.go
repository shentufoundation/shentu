package types

import (
	"encoding/hex"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateParams validates the module parameters
func ValidateParams(params *Params) error {
	if params == nil {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "params cannot be nil")
	}

	// Validate ProofMaxLockPeriod
	if params.ProofMaxLockPeriod == nil {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "proof max lock period cannot be nil")
	}

	// Validate TheoremMaxProofPeriod
	if params.TheoremMaxProofPeriod == nil {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "theorem max proof period cannot be nil")
	}

	return params.Validate()
}

// ValidateProgram validates a program
func ValidateProgram(program *Program) error {
	if program == nil {
		return errorsmod.Wrap(ErrProgramID, "program cannot be nil")
	}

	if len(program.ProgramId) == 0 {
		return errorsmod.Wrap(ErrProgramID, "program id cannot be empty")
	}

	if len(program.Name) == 0 {
		return errorsmod.Wrap(ErrProgramID, "program name cannot be empty")
	}

	// Check if AdminAddress is a valid address
	if _, err := sdk.AccAddressFromBech32(program.AdminAddress); err != nil {
		return errorsmod.Wrapf(err, "invalid admin address %s", program.AdminAddress)
	}

	// Other program validations can be added here

	return nil
}

// ValidateFinding validates a finding
func ValidateFinding(finding *Finding) error {
	if finding == nil {
		return errorsmod.Wrap(ErrFindingID, "finding cannot be nil")
	}

	if len(finding.ProgramId) == 0 {
		return errorsmod.Wrap(ErrProgramID, "program id cannot be empty")
	}

	if len(finding.FindingId) == 0 {
		return errorsmod.Wrap(ErrFindingID, "finding id cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(finding.SubmitterAddress); err != nil {
		return errorsmod.Wrapf(err, "invalid submitter address %s", finding.SubmitterAddress)
	}

	if !ValidFindingStatus(finding.Status) {
		return errorsmod.Wrap(ErrFindingStatusInvalid, "invalid finding status")
	}

	if !ValidFindingSeverityLevel(finding.SeverityLevel) {
		return errorsmod.Wrap(ErrFindingSeverityLevelInvalid, "invalid finding severity level")
	}

	return nil
}

// ValidFindingStatus returns true if the finding status is valid and false
// otherwise.
func ValidFindingStatus(status FindingStatus) bool {
	if status == FindingStatusSubmitted ||
		status == FindingStatusActive ||
		status == FindingStatusConfirmed ||
		status == FindingStatusPaid ||
		status == FindingStatusClosed {
		return true
	}
	return false
}

// ValidateTheorem validates a theorem
func ValidateTheorem(theorem *Theorem) error {
	if theorem == nil {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "theorem cannot be nil")
	}

	if theorem.Id == 0 {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "theorem id cannot be 0")
	}

	if len(theorem.Title) == 0 {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "theorem title cannot be empty")
	}

	if len(theorem.Description) == 0 {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "theorem description cannot be empty")
	}

	// Check if Proposer is a valid address
	if _, err := sdk.AccAddressFromBech32(theorem.Proposer); err != nil {
		return errorsmod.Wrapf(err, "invalid proposer address %s", theorem.Proposer)
	}

	// If the theorem is active, make sure the end time is set
	if theorem.Status == TheoremStatus_THEOREM_STATUS_PROOF_PERIOD && theorem.EndTime == nil {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "active theorem must have an end time")
	}

	return nil
}

// ValidateGrant validates a grant
func ValidateGrant(grant *Grant) error {
	if grant == nil {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "grant cannot be nil")
	}

	if grant.TheoremId == 0 {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "grant theorem id cannot be 0")
	}

	// Check if Grantor is a valid address
	if _, err := sdk.AccAddressFromBech32(grant.Grantor); err != nil {
		return errorsmod.Wrapf(err, "invalid grantor address %s", grant.Grantor)
	}

	// Validate amount is not empty
	if len(grant.Amount) == 0 {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "grant amount cannot be empty")
	}

	return nil
}

// ValidateDeposit validates a deposit
func ValidateDeposit(deposit *Deposit) error {
	if deposit == nil {
		return errorsmod.Wrap(ErrProofStatusInvalid, "deposit cannot be nil")
	}

	if len(deposit.ProofId) == 0 {
		return errorsmod.Wrap(ErrProofStatusInvalid, "deposit proof id cannot be empty")
	}

	// Check if Depositor is a valid address
	if _, err := sdk.AccAddressFromBech32(deposit.Depositor); err != nil {
		return errorsmod.Wrapf(err, "invalid depositor address %s", deposit.Depositor)
	}

	// Validate amount is not empty
	if len(deposit.Amount) == 0 {
		return errorsmod.Wrap(ErrProofStatusInvalid, "deposit amount cannot be empty")
	}

	return nil
}

// ValidateReward validates a reward
func ValidateReward(reward *Reward) error {
	if reward == nil {
		return errorsmod.Wrap(ErrProofOperatorNotAllowed, "reward cannot be nil")
	}

	if len(reward.Address) == 0 {
		return errorsmod.Wrap(ErrProofOperatorNotAllowed, "reward address cannot be empty")
	}

	// Check if Address is a valid address
	if _, err := sdk.AccAddressFromBech32(reward.Address); err != nil {
		return errorsmod.Wrapf(err, "invalid reward address %s", reward.Address)
	}

	// Validate reward is not empty
	if len(reward.Reward) == 0 {
		return errorsmod.Wrap(ErrProofOperatorNotAllowed, "reward amount cannot be empty")
	}

	return nil
}

// ValidateProof validates a proof
func ValidateProof(proof *Proof) error {
	if proof == nil {
		return errorsmod.Wrap(ErrProofStatusInvalid, "proof cannot be nil")
	}

	// validate proof.Id is a valid hex string
	if len(proof.Id) != 64 {
		return errorsmod.Wrap(ErrProofHashInvalid, "proof id must be a 64-character SHA-256 hash")
	}
	if _, err := hex.DecodeString(proof.Id); err != nil {
		return errorsmod.Wrap(ErrProofHashInvalid, "proof id must be a valid hex string")
	}

	if proof.TheoremId == 0 {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "proof theorem id cannot be 0")
	}

	// Check if Prover is a valid address
	if _, err := sdk.AccAddressFromBech32(proof.Prover); err != nil {
		return errorsmod.Wrapf(err, "invalid prover address %s", proof.Prover)
	}

	// If the proof is in hash lock period, make sure the SubmitTime is set
	if proof.Status == ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD {
		if proof.SubmitTime == nil {
			return errorsmod.Wrap(ErrProofStatusInvalid, "proof in hash lock period must have a submit time")
		}
		if proof.EndTime == nil {
			return errorsmod.Wrap(ErrProofStatusInvalid, "proof in hash lock period must have an end time")
		}
		if !proof.EndTime.After(*proof.SubmitTime) {
			return errorsmod.Wrap(ErrProofStatusInvalid, "proof end time must be after submit time")
		}
	}

	return nil
}
