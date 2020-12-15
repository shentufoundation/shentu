package types

import (
	"encoding/json"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Platform is a genesis type for certified platform of a validator
type Platform struct {
	Validator   crypto.PubKey
	Description string
}

// GenesisState - crisis genesis state
type GenesisState struct {
	Certifiers   []Certifier   `json:"certifiers"`
	Validators   []Validator   `json:"validators"`
	Platforms    []Platform    `json:"platforms"`
	Certificates []Certificate `json:"certificates"`
	Libraries    []Library     `json:"libraries"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(constantFee sdk.Coin, startingCertificateID CertificateID) GenesisState {
	return GenesisState{}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// ValidateGenesis - validate crisis genesis data
func ValidateGenesis(bz json.RawMessage) error {
	return nil
}

// GetGenesisStateFromAppState returns x/auth GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
