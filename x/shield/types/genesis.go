package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	ShieldOperator      sdk.AccAddress      `json:"shield_operator" yaml:"shield_operator"`
	NextPoolID          uint64              `json:"next_pool_id" yaml:"next_pool_id"`
	PoolParams          PoolParams          `json:"pool_params" yaml:"pool_params"`
	ClaimProposalParams ClaimProposalParams `json:"claim_proposal_params" yaml:"claim_proposal_params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	shieldOperator sdk.AccAddress, nextPoolID uint64, poolParams PoolParams, claimProposalParams ClaimProposalParams) GenesisState {
	return GenesisState{
		ShieldOperator:      shieldOperator,
		NextPoolID:          nextPoolID,
		PoolParams:          poolParams,
		ClaimProposalParams: claimProposalParams,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		NextPoolID:          uint64(1),
		PoolParams:          DefaultPoolParams(),
		ClaimProposalParams: DefaultClaimProposalParams(),
	}
}

// ValidateGenesis returns a default genesis state
func ValidateGenesis(bz json.RawMessage) error {
	return nil
}

// GetGenesisStateFromAppState returns GenesisState given raw application genesis state.
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
