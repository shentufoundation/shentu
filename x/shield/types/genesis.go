package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines the shield genesis state.
type GenesisState struct {
	ShieldAdmin          sdk.AccAddress      `json:"shield_admin" yaml:"shield_admin"`
	NextPoolID           uint64              `json:"next_pool_id" yaml:"next_pool_id"`
	NextPurchaseID       uint64              `json:"next_purchase_id" yaml:"next_purchase_id"`
	PoolParams           PoolParams          `json:"pool_params" yaml:"pool_params"`
	ClaimProposalParams  ClaimProposalParams `json:"claim_proposal_params" yaml:"claim_proposal_params"`
	TotalCollateral      sdk.Int             `json:"total_collateral" yaml:"total_collateral"`
	TotalWithdrawing     sdk.Int             `json:"total_withdrawing" yaml:"total_withdrawing"`
	TotalShield          sdk.Int             `json:"total_shield" yaml:"total_shield"`
	TotalLocked          sdk.Int             `json:"total_locked" yaml:"total_locked"`
	ServiceFees          MixedDecCoins       `json:"service_fees" yaml:"service_fees"`
	RemainingServiceFees MixedDecCoins       `json:"remaining_service_fees" yaml:"remaining_service_fees"`
	Pools                []Pool              `json:"pools" yaml:"pools"`
	Providers            []Provider          `json:"providers" yaml:"providers"`
	PurchaseLists        []PurchaseList      `json:"purchases" yaml:"purchases"`
	Withdraws            Withdraws           `json:"withdraws" yaml:"withdraws"`
	LastUpdateTime       time.Time           `json:"last_update_time" yaml:"last_update_time"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(shieldAdmin sdk.AccAddress, nextPoolID, nextPurchaseID uint64, poolParams PoolParams,
	claimProposalParams ClaimProposalParams, totalCollateral, totalWithdrawing, totalShield, totalLocked sdk.Int, serviceFees, remainingServiceFees MixedDecCoins,
	pools []Pool, providers []Provider, purchase []PurchaseList, withdraws Withdraws, lastUpdateTime time.Time) GenesisState {
	return GenesisState{
		ShieldAdmin:          shieldAdmin,
		NextPoolID:           nextPoolID,
		NextPurchaseID:       nextPurchaseID,
		PoolParams:           poolParams,
		ClaimProposalParams:  claimProposalParams,
		TotalCollateral:      totalCollateral,
		TotalWithdrawing:     totalWithdrawing,
		TotalShield:          totalShield,
		TotalLocked:          totalLocked,
		ServiceFees:          serviceFees,
		RemainingServiceFees: remainingServiceFees,
		Pools:                pools,
		Providers:            providers,
		PurchaseLists:        purchase,
		Withdraws:            withdraws,
		LastUpdateTime:       lastUpdateTime,
	}
}

// DefaultGenesisState returns a default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		NextPoolID:           uint64(1),
		NextPurchaseID:       uint64(1),
		PoolParams:           DefaultPoolParams(),
		ClaimProposalParams:  DefaultClaimProposalParams(),
		TotalCollateral:      sdk.ZeroInt(),
		TotalWithdrawing:     sdk.ZeroInt(),
		TotalShield:          sdk.ZeroInt(),
		TotalLocked:          sdk.ZeroInt(),
		ServiceFees:          InitMixedDecCoins(),
		RemainingServiceFees: InitMixedDecCoins(),
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
	if err := validatePoolParams(data.PoolParams); err != nil {
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
