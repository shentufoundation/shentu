package types

import (
	"encoding/json"

	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
)

type Contract struct {
	Address crypto.Address `json:"address"`
	Code    []byte         `json:"code"`
	Storage []Storage      `json:"storage"`
	Abi     []byte         `json:"abi"`
	Meta    []ContractMeta `json:"meta"`
}

type ContractMeta struct {
	CodeHash     []byte
	MetadataHash []byte
}

type Storage struct {
	Key   binary.Word256 `json:"key"`
	Value []byte         `json:"value"`
}

type Metadata struct {
	Hash     acmstate.MetadataHash `json:"hash"`
	Metadata string                `json:"metadata"`
}

// GenesisState is a cvm genesis state.
type GenesisState struct {
	// GasRate defines the gas exchange rate between Cosmos gas and CVM gas.
	// CVM gas equals to Cosmos Gas * gasRate.
	GasRate   uint64     `json:"gasrate"`
	Contracts []Contract `json:"contracts"`
	Metadata  []Metadata `json:"metadata"`
}

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(rate uint64) GenesisState {
	return GenesisState{
		GasRate: rate,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		GasRate: 1,
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
