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
	programs := make(map[uint64]int)
	for i, program := range data.Programs {
		if program.ProgramId > data.StartingProgramId {
			return ErrProgramID
		}

		programIndex, ok := programs[program.ProgramId]
		if ok {
			//repeat programId
			return sdkerrors.Wrapf(ErrProgramID, "already program[%s], this program[%s]",
				data.Programs[programIndex].String(), program.String())
		}

		if err := program.ValidateBasic(); err != nil {
			return err
		}
		programs[program.ProgramId] = i
	}

	for _, finding := range data.Findings {
		//Check if it is a valid programID
		_, ok := programs[finding.ProgramId]
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
