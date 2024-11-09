package types

import (
	errorsmod "cosmossdk.io/errors"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(programs []Program, findings []Finding) *GenesisState {
	return &GenesisState{
		Programs: programs,
		Findings: findings,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Programs: []Program{},
		Findings: []Finding{},
	}
}

// ValidateGenesis - validate bounty genesis data
func ValidateGenesis(data *GenesisState) error {
	programs := make(map[string]int)
	for i, program := range data.Programs {
		programIndex, ok := programs[program.ProgramId]
		if ok {
			//repeat programId
			return errorsmod.Wrapf(ErrProgramID, "already program[%s], this program[%s]",
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

		if err := finding.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}
