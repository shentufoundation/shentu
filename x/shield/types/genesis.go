package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	ShieldAdmin         sdk.AccAddress      `json:"shield_admin" yaml:"shield_admin"`
	NextPoolID          uint64              `json:"next_pool_id" yaml:"next_pool_id"`
	NextPurchaseID      uint64              `json:"next_purchase_id" yaml:"next_purchase_id"`
	PoolParams          PoolParams          `json:"pool_params" yaml:"pool_params"`
	ClaimProposalParams ClaimProposalParams `json:"claim_proposal_params" yaml:"claim_proposal_params"`
	Pools               []Pool              `json:"pools" yaml:"pools"`
	Collaterals         []Collateral        `json:"collaterals" yaml:"collaterals"`
	Providers           []Provider          `json:"providers" yaml:"providers"`
	PurchaseLists       []PurchaseList      `json:"purchases" yaml:"purchases"`
	Withdraws           Withdraws           `json:"withdraws" yaml:"withdraws"`
}

type WithdrawTimeSlice struct {
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(shieldAdmin sdk.AccAddress, nextPoolID, nextPurchaseID uint64, poolParams PoolParams,
	claimProposalParams ClaimProposalParams, pools []Pool, collaterals []Collateral,
	providers []Provider, purchase []PurchaseList, withdraws Withdraws) GenesisState {
	return GenesisState{
		ShieldAdmin:         shieldAdmin,
		NextPoolID:          nextPoolID,
		NextPurchaseID:      nextPurchaseID,
		PoolParams:          poolParams,
		ClaimProposalParams: claimProposalParams,
		Pools:               pools,
		Collaterals:         collaterals,
		Providers:           providers,
		PurchaseLists:       purchase,
		Withdraws:           withdraws,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		NextPoolID:          uint64(1),
		NextPurchaseID:      uint64(1),
		PoolParams:          DefaultPoolParams(),
		ClaimProposalParams: DefaultClaimProposalParams(),
	}
}

// ValidateGenesis validates shield genesis data.
func ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	if data.NextPoolID < 1 {
		return fmt.Errorf("failed to validate %s genesis state: NextPoolID must be positive ", ModuleName)
	}
	if err := ValidatePoolParams(data.PoolParams); err != nil {
		return fmt.Errorf("failed to validate %s pool params: %w", ModuleName, err)
	}
	if err := validateClaimProposalParams(data.ClaimProposalParams); err != nil {
		return fmt.Errorf("failed to validate %s claim proposal params: %w", ModuleName, err)
	}

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
