package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(constantFee sdk.Coin, startingCertificateID uint64) GenesisState {
	return GenesisState{NextCertificateId: startingCertificateID}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{NextCertificateId: 1}
}

// ValidateGenesis - validate crisis genesis data
func ValidateGenesis(bz json.RawMessage) error {
	return nil
}

// GetGenesisStateFromAppState returns cert module GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Marshaler, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (p Platform) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(p.ValidatorPubkey, &pubKey)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, certificate := range g.Certificates {
		err := certificate.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	for _, platform := range g.Platforms {
		err := platform.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}
