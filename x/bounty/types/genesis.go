package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(constantFee sdk.Coin) *GenesisState {
	return &GenesisState{}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		StartingProgramId: 0,
		StartingFindingId: 0,
		Programs:          []Program{},
		Findings:          []Finding{},
	}
}

// ValidateGenesis - validate bounty genesis data
func ValidateGenesis(data *GenesisState) error {
	// TODO: implement ValidateGenesis
	return nil
}
