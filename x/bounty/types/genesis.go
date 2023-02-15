package types

import "fmt"

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
			return fmt.Errorf("error program id")
		}

		_, ok := programIds[program.ProgramId]
		if ok {
			//repeat program
			return fmt.Errorf("repeat programId:%d", program.ProgramId)
		}

		if !program.Valid() {
			return fmt.Errorf("invalid program [programId:%d]", program.ProgramId)
		}
		programIds[program.ProgramId] = true
	}

	for _, finding := range data.Findings {
		//Check if it is a valid programID
		_, ok := programIds[finding.ProgramId]
		if !ok {
			return fmt.Errorf("programID:%d is invalid programID", finding.ProgramId)
		}

		if finding.FindingId > data.StartingFindingId {
			return fmt.Errorf("error finding id")
		}
		if !finding.Valid() {
			return fmt.Errorf("invalid finding [fingdingId:%d]", finding.FindingId)
		}
	}
	return nil
}
