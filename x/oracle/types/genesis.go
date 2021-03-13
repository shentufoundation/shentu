package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState constructs a GenesisState object.
func NewGenesisState(operators []Operator, totalCollateral sdk.Coins, poolParams LockedPoolParams, taskParams TaskParams,
	withdraws []Withdraw, tasks []Task) GenesisState {
	return GenesisState{
		Operators:       operators,
		TotalCollateral: totalCollateral,
		PoolParams:      &poolParams,
		TaskParams:      &taskParams,
		Withdraws:       withdraws,
		Tasks:           tasks,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	state := NewGenesisState(nil, nil, DefaultLockedPoolParams(), DefaultTaskParams(), nil, nil)
	return &state
}

// ValidateGenesis validates oracle genesis data.
func ValidateGenesis(gs GenesisState) error {
	operators := gs.Operators
	withdraws := gs.Withdraws
	sum := sdk.NewCoins()
	for _, operator := range operators {
		sum = sum.Add(operator.Collateral...)
	}
	for _, withdraw := range withdraws {
		if withdraw.DueBlock < 0 {
			return ErrInvalidDueBlock
		}
	}
	if !sum.IsEqual(gs.TotalCollateral) {
		panic(ErrTotalCollateralNotEqual)
	}
	if gs.PoolParams.LockedInBlocks < 0 || gs.PoolParams.MinimumCollateral < 0 {
		panic(ErrInvalidPoolParams)
	}
	return nil
}
