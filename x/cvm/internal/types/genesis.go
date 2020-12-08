package types

import (
	"encoding/json"
)

type Contracts = []Contract

type Metadatas = []Metadata

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(rate uint64) GenesisState {
	return GenesisState{
		GasRate: rate,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		GasRate: DefaultGasRate,
	}
}

// ValidateGenesis validates cvm genesis data.
func ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return nil
}
