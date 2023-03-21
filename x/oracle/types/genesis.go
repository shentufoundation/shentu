package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState constructs a GenesisState object.
func NewGenesisState(operators []Operator, totalCollateral sdk.Coins, poolParams LockedPoolParams, taskParams TaskParams,
	withdraws []Withdraw, tasks []Task, txTasks []TxTask) GenesisState {
	return GenesisState{
		Operators:       operators,
		TotalCollateral: totalCollateral,
		PoolParams:      &poolParams,
		TaskParams:      &taskParams,
		Withdraws:       withdraws,
		Tasks:           tasks,
		TxTasks:         txTasks,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	state := NewGenesisState(nil, nil, DefaultLockedPoolParams(), DefaultTaskParams(), nil, nil, nil)
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
		return ErrTotalCollateralNotEqual
	}
	if err := validatePoolParams(*gs.PoolParams); err != nil {
		return err
	}
	if err := validateTaskParams(*gs.TaskParams); err != nil {
		return err
	}
	return nil
}
