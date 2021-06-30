package types

import (
	fmt "fmt"
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
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		GasRate: DefaultGasRate,
	}
}

// Validate performs validation of cvm genesis data.
func (gs GenesisState) Validate() error {
	if gs.GasRate > 100 {
		return fmt.Errorf("failed to validate %s genesis state: GasRate is too high", ModuleName)
	}

	for _, metadata := range gs.Metadatas {
		if len(metadata.Hash) != 32 {
			return fmt.Errorf("failed to validate %s genesis state: A metadata hash is not 256 bits", ModuleName)
		}
	}
	return nil

}
