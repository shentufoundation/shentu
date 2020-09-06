package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Operators       []Operator       `json:"operators"`
	TotalCollateral sdk.Coins        `json:"total_collateral"`
	PoolParams      LockedPoolParams `json:"pool_params"`
	TaskParams      TaskParams       `json:"task_params"`
	Withdraws       []Withdraw       `json:"withdraws"`
	Tasks           []Task           `json:"tasks"`
}

// NewGenesisState constructs a GenesisState object.
func NewGenesisState(operators []Operator, totalCollateral sdk.Coins, poolParams LockedPoolParams, taskParams TaskParams,
	withdraws []Withdraw, tasks []Task) GenesisState {
	return GenesisState{
		Operators:       operators,
		TotalCollateral: totalCollateral,
		PoolParams:      poolParams,
		TaskParams:      taskParams,
		Withdraws:       withdraws,
		Tasks:           tasks,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() GenesisState {
	return NewGenesisState(
		nil,
		nil,
		DefaultLockedPoolParams(),
		DefaultTaskParams(),
		nil,
		nil,
	)
}

// ValidateGenesis validates oracle genesis data.
func ValidateGenesis(bz json.RawMessage) error {
	var gs GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}
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
