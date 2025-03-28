package types

import (
	errorsmod "cosmossdk.io/errors"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	startingTheoremId uint64,
	params Params,
) *GenesisState {
	return &GenesisState{
		StartingTheoremId: startingTheoremId,
		Params:            &params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(1, DefaultParams())
}

// ValidateGenesis - validate bounty genesis data
func ValidateGenesis(data *GenesisState) error {
	if data.StartingTheoremId == 0 {
		return errorsmod.Wrap(ErrNoTheoremMsgs, "starting theorem id cannot be 0")
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

		//if err := program.ValidateBasic(); err != nil {
		//	return err
		//}
		programs[program.ProgramId] = i
	}

	for _, finding := range data.Findings {
		//Check if it is a valid programID
		_, ok := programs[finding.ProgramId]
		if !ok {
			return ErrProgramID
		}

		if err := finding.ValidateBasic(); err != nil {
			return err
		}
	}

	//theorems := make(map[uint64]bool)
	//for _, theorem := range data.Theorems {
	//	if theorem.Id == 0 {
	//		return errorsmod.Wrap(ErrNoTheoremMsgs, "theorem id cannot be 0")
	//	}
	//
	//	if theorems[theorem.Id] {
	//		return errorsmod.Wrapf(ErrNoTheoremMsgs, "duplicate theorem id %d", theorem.Id)
	//	}
	//
	//	if err := theorem.ValidateBasic(); err != nil {
	//		return errorsmod.Wrapf(err, "invalid theorem %d", theorem.Id)
	//	}
	//
	//	theorems[theorem.Id] = true
	//}
	//
	//proofs := make(map[string]bool)
	//for _, proof := range data.Proofs {
	//	if proof.Id == "" {
	//		return errorsmod.Wrap(ErrProofStatusInvalid, "proof id cannot be empty")
	//	}
	//
	//	if proofs[proof.Id] {
	//		return errorsmod.Wrapf(ErrProofAlreadyExists, "duplicate proof id %s", proof.Id)
	//	}
	//
	//	if err := proof.ValidateBasic(); err != nil {
	//		return errorsmod.Wrapf(err, "invalid proof %s", proof.Id)
	//	}
	//
	//	// Check if theorem exists for this proof
	//	if proof.TheoremId != 0 && !theorems[proof.TheoremId] {
	//		return errorsmod.Wrapf(ErrTheoremProposal, "theorem %d for proof %s does not exist",
	//			proof.TheoremId, proof.Id)
	//	}
	//
	//	proofs[proof.Id] = true
	//}
	//
	//// Validate rewards
	//for _, reward := range data.Rewards {
	//	if reward.Address == "" {
	//		return errorsmod.Wrap(ErrProofOperatorNotAllowed, "reward address cannot be empty")
	//	}
	//
	//	if err := reward.ValidateBasic(); err != nil {
	//		return errorsmod.Wrapf(err, "invalid reward for address %s", reward.Address)
	//	}
	//}

	return nil
}
