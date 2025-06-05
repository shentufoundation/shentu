package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState constructs a GenesisState object.
func NewGenesisState(
	operators []Operator,
	totalCollateral sdk.Coins,
	poolParams LockedPoolParams,
	taskParams TaskParams,
	withdraws []Withdraw,
	tasks []Task,
	txTasks []TxTask,
) GenesisState {
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
func ValidateGenesis(genesisState GenesisState) error {
	operators := genesisState.Operators
	withdraws := genesisState.Withdraws
	totalCollateral := sdk.NewCoins()
	for _, operator := range operators {
		totalCollateral = totalCollateral.Add(operator.Collateral...)
	}
	for _, withdraw := range withdraws {
		if withdraw.DueBlock < 0 {
			return ErrInvalidDueBlock
		}
	}
	if !totalCollateral.Equal(genesisState.TotalCollateral) {
		return ErrTotalCollateralNotEqual
	}
	if err := validatePoolParams(*genesisState.PoolParams); err != nil {
		return err
	}
	if err := validateTaskParams(*genesisState.TaskParams); err != nil {
		return err
	}
	return nil
}
