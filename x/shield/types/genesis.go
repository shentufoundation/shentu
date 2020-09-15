package types

type GenesisState struct {
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState()
}

// ValidateGenesis returns a default genesis state
func ValidateGenesis(data GenesisState) error {
	return nil
}
