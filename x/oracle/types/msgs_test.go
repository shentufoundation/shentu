package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/require"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

// test ValidateBasic for NewMsgCreateOperator
func Test_NewMsgCreateOperator(t *testing.T) {
	valAddr := sdk.AccAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	valCollateral := sdk.Coins{sdk.NewInt64Coin("uctk", 1234)}
	proposer := sdk.AccAddress([]byte{20})

	// emptyAddr := sdk.AccAddress{}
	//cannot generate a neg coin
	// negCollateral := sdk.Coins{sdk.NewInt64Coin("uctk", 1234)}
	tests := []struct {
		ts_name, name string
		address       sdk.AccAddress
		collateral    sdk.Coins
		proposer      sdk.AccAddress
		expectPass    bool
	}{
		{"basic good", "abc", valAddr, valCollateral, proposer, true},
		// {"empty address", "abc", emptyAddr, valCollateral, proposer, false},
		// {"negative Collateral", "abc", valAddr, negCollateral, proposer, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgCreateOperator(tc.address, tc.collateral, tc.proposer, tc.name)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		}
	}
}

// test ValidateBasic for NewMsgRemoveOperator
func Test_NewMsgRemoveOperator(t *testing.T) {
	valAddr := sdk.AccAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	proposer := sdk.AccAddress([]byte{20})
	// emptyAddr := sdk.AccAddress{}

	tests := []struct {
		ts_name    string
		address    sdk.AccAddress
		proposer   sdk.AccAddress
		expectPass bool
	}{
		{"basic good", valAddr, proposer, true},
		// {"empty address", emptyAddr, proposer, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgRemoveOperator(tc.address, tc.proposer)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		}
	}
}

// test ValidateBasic for NewMsgAddCollateral
func Test_NewMsgAddCollateral(t *testing.T) {
	valAddr := sdk.AccAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	valCollateral := sdk.Coins{sdk.NewInt64Coin("uctk", 1234)}
	// emptyAddr := sdk.AccAddress{}
	// negCollateral := sdk.Coins{sdk.NewInt64Coin("uctk", -1234)}

	tests := []struct {
		ts_name    string
		address    sdk.AccAddress
		collateral sdk.Coins
		expectPass bool
	}{
		{"basic good", valAddr, valCollateral, true},
		// {"empty address", emptyAddr, valCollateral, false},
		// {"negative Collateral", valAddr, negCollateral, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgAddCollateral(tc.address, tc.collateral)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		}
	}
}

// test ValidateBasic for NewMsgReduceCollateral
func Test_NNewMsgReduceCollateral(t *testing.T) {
	valAddr := sdk.AccAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	valCollateral := sdk.Coins{sdk.NewInt64Coin("uctk", 1234)}
	// emptyAddr := sdk.AccAddress{}
	// negCollateral := sdk.Coins{sdk.NewInt64Coin("uctk", -1234)}

	tests := []struct {
		ts_name    string
		address    sdk.AccAddress
		collateral sdk.Coins
		expectPass bool
	}{
		{"basic good", valAddr, valCollateral, true},
		// {"empty address", emptyAddr, valCollateral, false},
		// {"negative Collateral", valAddr, negCollateral, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgReduceCollateral(tc.address, tc.collateral)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		}
	}
}

// test ValidateBasic for NewMsgWithdrawReward
func Test_NewMsgWithdrawReward(t *testing.T) {
	valAddr := sdk.AccAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	// emptyAddr := sdk.AccAddress{}

	tests := []struct {
		ts_name    string
		address    sdk.AccAddress
		expectPass bool
	}{
		{"basic good", valAddr, true},
		// {"empty address", emptyAddr, false},
	}

	for _, tc := range tests {
		msg := types.NewMsgWithdrawReward(tc.address)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.ts_name)
		}
	}
}
