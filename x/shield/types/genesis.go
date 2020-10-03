package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	ShieldAdmin         sdk.AccAddress      `json:"shield_admin" yaml:"shield_admin"`
	NextPoolID          uint64              `json:"next_pool_id" yaml:"next_pool_id"`
	PoolParams          PoolParams          `json:"pool_params" yaml:"pool_params"`
	ClaimProposalParams ClaimProposalParams `json:"claim_proposal_params" yaml:"claim_proposal_params"`
	Pools               []Pool              `json:"pools" yaml:"pools"`
	Collaterals         []Collateral        `json:"collaterals" yaml:"collaterals"`
	Providers           []Provider          `json:"providers" yaml:"providers"`
	Purchases           []Purchase          `json:"purchases" yaml:"purchases"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	shieldAdmin sdk.AccAddress, nextPoolID uint64, poolParams PoolParams, claimProposalParams ClaimProposalParams,
	pools []Pool, collaterals []Collateral, providers []Provider, purchase []Purchase) GenesisState {
	return GenesisState{
		ShieldAdmin:         shieldAdmin,
		NextPoolID:          nextPoolID,
		PoolParams:          poolParams,
		ClaimProposalParams: claimProposalParams,
		Pools:               pools,
		Collaterals:         collaterals,
		Providers:           providers,
		Purchases:           purchase,
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
