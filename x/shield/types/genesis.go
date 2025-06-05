package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(remainingServiceFees sdk.DecCoins, providers []Provider) GenesisState {
	return GenesisState{
		RemainingServiceFees: remainingServiceFees,
		Providers:            providers,
	}
}

// DefaultGenesisState returns a default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		RemainingServiceFees: sdk.DecCoins{},
	}
}

// ValidateGenesis validates shield genesis data.
func ValidateGenesis(_ GenesisState) error {
	return nil
}

// GetGenesisStateFromAppState returns GenesisState given raw application genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
