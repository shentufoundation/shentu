package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type GenesisState struct {
	ShieldOperator sdk.AccAddress `json:"shield_operator" yaml:"shield_operator"`
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
