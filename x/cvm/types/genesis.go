package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hyperledger/burrow/crypto"
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
		return fmt.Errorf("failed to validate %s genesis state: GasRate must ???", ModuleName)
	}

	for _, contract := range gs.Contracts {
		if len(contract.Address) != crypto.AddressLength {
			return fmt.Errorf("failed to validate %s genesis state: Incorrect contract address length %s", ModuleName, sdk.AccAddress(contract.Address.Bytes()).String())
		}
	}

	for _, metadata := range gs.Metadatas {
		if len(metadata.Hash) != 32 {
			return fmt.Errorf("failed to validate %s genesis state: A metadata hash is not 256 bits", ModuleName)
		}
	}
	return nil

}
