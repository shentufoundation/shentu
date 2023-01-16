package types

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
	// TODO: implement ValidateGenesis
	return nil
}
