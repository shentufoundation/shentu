package v1beta1

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(shieldAdmin sdk.AccAddress, nextPoolID uint64, poolParams PoolParams,
	claimProposalParams ClaimProposalParams, totalCollateral, totalWithdrawing, totalShield, totalClaimed sdk.Int,
	serviceFees sdk.DecCoins, pools []Pool, providers []Provider, withdraws []Withdraw,
	globalStakingPool sdk.Int, stakingPurchases []Purchase,
	reserve Reserve, pendingPayouts []PendingPayout, blockRewardParams BlockRewardParams) GenesisState {
	return GenesisState{
		ShieldAdmin:    shieldAdmin.String(),
		NextPoolId:     nextPoolID,
		Fees:           serviceFees,
		Pools:          pools,
		Providers:      providers,
		Withdraws:      withdraws,
		Purchases:      stakingPurchases,
		Reserve:        reserve,
		PendingPayouts: pendingPayouts,
		GlobalPools:    NewGlobalPools(totalCollateral, totalWithdrawing, totalShield, totalClaimed, globalStakingPool),
		ShieldParams:   NewShieldParams(poolParams, claimProposalParams, blockRewardParams),
	}
}

// DefaultGenesisState returns a default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		NextPoolId:   uint64(1),
		Fees:         sdk.NewDecCoins(),
		Reserve:      InitialReserve(),
		GlobalPools:  DefaultGlobalPools(),
		ShieldParams: DefaultShieldParams(),
	}
}

// ValidateGenesis validates shield genesis data.
func ValidateGenesis(data GenesisState) error {
	if data.NextPoolId < 1 {
		return fmt.Errorf("failed to validate %s genesis state: NextPoolID must be positive ", types.ModuleName)
	}
	if data.Reserve.Amount.IsNegative() {
		return fmt.Errorf("reserve amount is negative %v", data.Reserve.Amount)
	}
	if err := ValidatePoolParams(data.ShieldParams.PoolParams); err != nil {
		return fmt.Errorf("failed to validate %s pool params: %w", types.ModuleName, err)
	}
	if err := ValidateClaimProposalParams(data.ShieldParams.ClaimProposalParams); err != nil {
		return fmt.Errorf("failed to validate %s claim proposal params: %w", types.ModuleName, err)
	}
	if err := ValidateBlockRewardParams(data.ShieldParams.BlockRewardParams); err != nil {
		return fmt.Errorf("failed to validate %s block reward params: %w", types.ModuleName, err)
	}

	return nil
}

// GetGenesisStateFromAppState returns GenesisState given raw application genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[types.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[types.ModuleName], &genesisState)
	}
	return genesisState
}

// NewGlobalPools creates a new GlobalPools object.
func NewGlobalPools(totalCollateral, totalWithdrawing, totalShield, totalClaimed, globalStakingPool sdk.Int) GlobalPools {
	return GlobalPools{
		TotalCollateral:   totalCollateral,
		TotalWithdrawing:  totalWithdrawing,
		TotalShield:       totalShield,
		TotalClaimed:      totalClaimed,
		GlobalStakingPool: globalStakingPool,
	}
}

// DefaultGlobalPools returns a default GlobalPools object.
func DefaultGlobalPools() GlobalPools {
	return GlobalPools{
		TotalCollateral:  sdk.ZeroInt(),
		TotalWithdrawing: sdk.ZeroInt(),
		TotalShield:      sdk.ZeroInt(),
		TotalClaimed:     sdk.ZeroInt(),
	}
}

// NewShieldParams creates a new ShieldParams object.
func NewShieldParams(poolParams PoolParams, claimProposalParams ClaimProposalParams, blockRewardParams BlockRewardParams) ShieldParams {
	return ShieldParams{
		PoolParams:          poolParams,
		ClaimProposalParams: claimProposalParams,
		BlockRewardParams:   blockRewardParams,
	}
}

// DefaultShieldParams returns a default ShieldParams object.
func DefaultShieldParams() ShieldParams {
	return ShieldParams{
		PoolParams:          DefaultPoolParams(),
		ClaimProposalParams: DefaultClaimProposalParams(),
		BlockRewardParams:   DefaultBlockRewardParams(),
	}
}
