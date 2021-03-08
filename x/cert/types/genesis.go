package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(constantFee sdk.Coin, startingCertificateID CertificateID) GenesisState {
	return GenesisState{}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
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
func (g GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for i := 0; i < len(g.Certificates); i++ {
		var cert Certificate
		err := unpacker.UnpackAny(g.Certificates[i], &cert)
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

	for _, validator := range g.Validators {
		err := validator.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}
