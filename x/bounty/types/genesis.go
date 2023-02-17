package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(startingProgramID, startingFindingID uint64, programs []Program, findings []Finding) *GenesisState {
	return &GenesisState{
		StartingProgramId: startingProgramID,
		StartingFindingId: startingFindingID,
		Programs:          programs,
		Findings:          findings,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		StartingProgramId: 1,
		StartingFindingId: 1,
		Programs:          []Program{},
		Findings:          []Finding{},
	}
}

// ValidateGenesis - validate bounty genesis data
func ValidateGenesis(data *GenesisState) error {
	programIds := make(map[uint64]bool)
	for _, program := range data.Programs {
		if program.ProgramId > data.StartingProgramId {
			return ErrProgramID
		}

		_, ok := programIds[program.ProgramId]
		if ok {
			return sdkerrors.Wrapf(ErrProgramID, "repeat programId:%d", program.ProgramId)
		}

		if err := program.ValidateBasic(); err != nil {
			return err
		}
		programIds[program.ProgramId] = true
	}

	for _, finding := range data.Findings {
		//Check if it is a valid programID
		_, ok := programIds[finding.ProgramId]
		if !ok {
			return ErrProgramID
		}

		if finding.FindingId > data.StartingFindingId {
			return ErrFindingID
		}
		if err := finding.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}
