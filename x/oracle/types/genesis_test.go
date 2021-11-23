package types_test

import (
	"testing"

	"github.com/test-go/testify/require"

	"github.com/certikfoundation/shentu/v2/x/oracle/types"
)

func Test_GenesisStateValidate(t *testing.T) {
	tests := []struct {
		name         string
		genesisState types.GenesisState
		expPass      bool
	}{
		{
			"Genesis: valid genesis state",
			types.GenesisState{
				Operators:       []types.Operator{operator1, operator2},
				TotalCollateral: coins1e5,
				PoolParams:      &validPoolParams,
				TaskParams:      &validTaskParams,
				Withdraws:       []types.Withdraw{validWithdraw},
				Tasks:           []types.Task{validTask},
			},
			true,
		},
		{
			"Genesis: mismatched total collateral",
			types.GenesisState{
				Operators:       []types.Operator{operator1, operator2},
				TotalCollateral: coins1234,
				PoolParams:      &validPoolParams,
				TaskParams:      &validTaskParams,
				Withdraws:       []types.Withdraw{validWithdraw},
				Tasks:           []types.Task{validTask},
			},
			false,
		},
		{
			"Genesis: invalid pool params",
			types.GenesisState{
				Operators:       []types.Operator{operator1, operator2},
				TotalCollateral: coins1e5,
				PoolParams:      &invalidPoolParams,
				TaskParams:      &validTaskParams,
				Withdraws:       []types.Withdraw{validWithdraw},
				Tasks:           []types.Task{validTask},
			},
			false,
		},
		{
			"Genesis: invalid task params",
			types.GenesisState{
				Operators:       []types.Operator{operator1, operator2},
				TotalCollateral: coins1e5,
				PoolParams:      &validPoolParams,
				TaskParams:      &invalidTaskParams,
				Withdraws:       []types.Withdraw{validWithdraw},
				Tasks:           []types.Task{validTask},
			},
			false,
		},
		{
			"Genesis: malformed withdraw",
			types.GenesisState{
				Operators:       []types.Operator{operator1, operator2},
				TotalCollateral: coins1e5,
				PoolParams:      &validPoolParams,
				TaskParams:      &validTaskParams,
				Withdraws:       []types.Withdraw{invalidWithdraw},
				Tasks:           []types.Task{validTask},
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateGenesis(tc.genesisState)
			if tc.expPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
