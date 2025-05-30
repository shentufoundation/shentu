package types

import (
	errorsmod "cosmossdk.io/errors"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	startingTheoremID uint64,
	params Params,
) *GenesisState {
	return &GenesisState{
		StartingTheoremId: startingTheoremID,
		Params:            &params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultStartingTheoremID, DefaultParams())
}

// ValidateGenesis - validate bounty genesis data
func ValidateGenesis(data *GenesisState) error {
	if data == nil {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "genesis state cannot be nil")
	}

	if data.StartingTheoremId == 0 {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "starting theorem id cannot be 0")
	}

	if data.Params == nil {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "params cannot be nil")
	}

	if err := data.Params.Validate(); err != nil {
		return errorsmod.Wrap(err, "invalid params")
	}

	programs := make(map[string]int)
	for i, program := range data.Programs {
		programIndex, ok := programs[program.ProgramId]
		if ok {
			//repeat programId
			return errorsmod.Wrapf(ErrProgramID, "already program[%s], this program[%s]",
				data.Programs[programIndex].String(), program.String())
		}

		if err := ValidateProgram(program); err != nil {
			return errorsmod.Wrapf(err, "invalid program %s", program.ProgramId)
		}
		programs[program.ProgramId] = i
	}

	findings := make(map[string]bool)
	for _, finding := range data.Findings {
		//Check if it is a valid programID
		_, ok := programs[finding.ProgramId]
		if !ok {
			return errorsmod.Wrapf(ErrProgramID, "program %s for finding %s does not exist",
				finding.ProgramId, finding.FindingId)
		}

		if findings[finding.FindingId] {
			return errorsmod.Wrapf(ErrFindingID, "duplicate finding id %s", finding.FindingId)
		}

		if err := ValidateFinding(finding); err != nil {
			return errorsmod.Wrapf(err, "invalid finding %s", finding.FindingId)
		}

		findings[finding.FindingId] = true
	}

	theorems := make(map[uint64]bool)
	for _, theorem := range data.Theorems {
		if theorem.Id == 0 {
			return errorsmod.Wrap(ErrNoTheoremMsgs, "theorem id cannot be 0")
		}

		if theorems[theorem.Id] {
			return errorsmod.Wrapf(ErrNoTheoremMsgs, "duplicate theorem id %d", theorem.Id)
		}

		if err := ValidateTheorem(theorem); err != nil {
			return errorsmod.Wrapf(err, "invalid theorem %d", theorem.Id)
		}

		theorems[theorem.Id] = true
	}

	proofs := make(map[string]bool)
	for _, proof := range data.Proofs {
		if len(proof.Id) == 0 {
			return errorsmod.Wrap(ErrProofStatusInvalid, "proof id cannot be empty")
		}

		if proofs[proof.Id] {
			return errorsmod.Wrapf(ErrProofAlreadyExists, "duplicate proof id %s", proof.Id)
		}

		if err := ValidateProof(proof); err != nil {
			return errorsmod.Wrapf(err, "invalid proof %s", proof.Id)
		}

		// Check if theorem exists for this proof
		if proof.TheoremId != 0 && !theorems[proof.TheoremId] {
			return errorsmod.Wrapf(ErrTheoremProposal, "theorem %d for proof %s does not exist",
				proof.TheoremId, proof.Id)
		}

		proofs[proof.Id] = true
	}

	// Validate grants
	for _, grant := range data.Grants {
		if grant.TheoremId == 0 {
			return errorsmod.Wrap(ErrNoTheoremMsgs, "grant theorem id cannot be 0")
		}

		if !theorems[grant.TheoremId] {
			return errorsmod.Wrapf(ErrTheoremProposal, "theorem %d for grant does not exist", grant.TheoremId)
		}

		if err := ValidateGrant(grant); err != nil {
			return errorsmod.Wrapf(err, "invalid grant for theorem %d", grant.TheoremId)
		}
	}

	// Validate deposits
	for _, deposit := range data.Deposits {
		if len(deposit.ProofId) == 0 {
			return errorsmod.Wrap(ErrProofStatusInvalid, "deposit proof id cannot be empty")
		}

		if !proofs[deposit.ProofId] {
			return errorsmod.Wrapf(ErrProofAlreadyExists, "proof %s for deposit does not exist", deposit.ProofId)
		}

		if err := ValidateDeposit(deposit); err != nil {
			return errorsmod.Wrapf(err, "invalid deposit for proof %s", deposit.ProofId)
		}
	}

	// Validate rewards
	for _, reward := range data.Rewards {
		if len(reward.Address) == 0 {
			return errorsmod.Wrap(ErrProofOperatorNotAllowed, "reward address cannot be empty")
		}

		if err := ValidateReward(reward); err != nil {
			return errorsmod.Wrapf(err, "invalid reward for address %s", reward.Address)
		}
	}

	return nil
}
