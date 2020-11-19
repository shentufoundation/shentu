package types

import (
	"encoding/json"
	"fmt"

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
	Certifiers        []Certifier   `json:"certifiers"`
	Validators        []Validator   `json:"validators"`
	Platforms         []Platform    `json:"platforms"`
	Certificates      []Certificate `json:"certificates"`
	Libraries         []Library     `json:"libraries"`
	NextCertificateID uint64        `json:"next_certificate_id" yaml:"next_certificate_id"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(constantFee sdk.Coin, startingCertificateID CertificateID) GenesisState {
	return GenesisState{}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		NextCertificateID:           uint64(1),
	}
}

// ValidateGenesis - validate crisis genesis data
func ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	if data.NextPoolID < 1 {
		return fmt.Errorf("failed to validate %s genesis state: NextPoolID must be positive ", ModuleName)
	}
	
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
