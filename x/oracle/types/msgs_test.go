package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/require"

	"github.com/certikfoundation/shentu/v2/x/oracle/types"
)

// test ValidateBasic for NewMsgCreateOperator
func Test_NewMsgCreateOperator(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	addrEmpty := sdk.AccAddress([]byte(""))
	// addrLong := sdk.AccAddress([]byte("Purposefully long address"))

	uctk123 := sdk.NewCoins(sdk.NewInt64Coin("uctk", 123))
	uctk0 := sdk.NewCoins(sdk.NewInt64Coin("uctk", 0))
	uctk123eth123 := sdk.NewCoins(sdk.NewInt64Coin("uctk", 123), sdk.NewInt64Coin("eth", 123))
	uctk123eth0 := sdk.Coins{sdk.NewInt64Coin("uctk", 123), sdk.NewInt64Coin("eth", 0)}

	name := "abc"

	cases := []struct {
		name       string
		expectPass bool
		msg        *types.MsgCreateOperator
	}{
		{"valid with one denom", true, types.NewMsgCreateOperator(addr1, uctk123, addr2, name)},
		{"valid with two denoms", true, types.NewMsgCreateOperator(addr1, uctk123eth123, addr2, name)},
		// {true, types.NewMsgCreateOperator(addrLong, uctk123, addr2, name)},
		// {true, types.NewMsgCreateOperator(addr1, uctk123, addrLong, name)},
		{"non-positive coin", false, types.NewMsgCreateOperator(addr1, uctk0, addr2, name)},
		{"non-positive multicoins", false, types.NewMsgCreateOperator(addr1, uctk123eth0, addr2, name)},
		{"invalid operator address", false, types.NewMsgCreateOperator(addrEmpty, uctk123, addr2, name)},
		{"invalid proposer address", false, types.NewMsgCreateOperator(addr1, uctk123, addrEmpty, name)},
	}

	for _, tc := range cases {
		if tc.expectPass {
			require.Nil(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for NewMsgRemoveOperator
func Test_NewMsgRemoveOperator(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))
	addrEmpty := sdk.AccAddress([]byte(""))

	cases := []struct {
		name       string
		expectPass bool
		msg        *types.MsgRemoveOperator
	}{
		{"valid addresses", true, types.NewMsgRemoveOperator(addr1, addr2)},
		{"invalid operator address", false, types.NewMsgRemoveOperator(addrEmpty, addr2)},
		{"invalid proposer address", false, types.NewMsgRemoveOperator(addr1, addrEmpty)},
	}

	for _, tc := range cases {
		if tc.expectPass {
			require.Nil(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for NewMsgAddCollateral
func Test_NewMsgAddCollateral(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr1_______________"))
	addrEmpty := sdk.AccAddress([]byte(""))

	uctk123 := sdk.NewCoins(sdk.NewInt64Coin("uctk", 123))
	uctk0 := sdk.NewCoins(sdk.NewInt64Coin("uctk", 0))
	uctk123eth123 := sdk.NewCoins(sdk.NewInt64Coin("uctk", 123), sdk.NewInt64Coin("eth", 123))
	uctk123eth0 := sdk.Coins{sdk.NewInt64Coin("uctk", 123), sdk.NewInt64Coin("eth", 0)}

	cases := []struct {
		name       string
		expectPass bool
		msg        *types.MsgAddCollateral
	}{
		{"valid with one denom", true, types.NewMsgAddCollateral(addr, uctk123)},
		{"valid with two denoms", true, types.NewMsgAddCollateral(addr, uctk123eth123)},
		{"non-positive coin", false, types.NewMsgAddCollateral(addr, uctk0)},
		{"non-positive multicoins", false, types.NewMsgAddCollateral(addr, uctk123eth0)},
		{"empty address", false, types.NewMsgAddCollateral(addrEmpty, uctk123)},
	}

	for _, tc := range cases {
		if tc.expectPass {
			require.Nil(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for NewMsgReduceCollateral
func Test_NewMsgReduceCollateral(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr1_______________"))
	addrEmpty := sdk.AccAddress([]byte(""))

	uctk123 := sdk.NewCoins(sdk.NewInt64Coin("uctk", 123))
	uctk0 := sdk.NewCoins(sdk.NewInt64Coin("uctk", 0))
	uctk123eth123 := sdk.NewCoins(sdk.NewInt64Coin("uctk", 123), sdk.NewInt64Coin("eth", 123))
	uctk123eth0 := sdk.Coins{sdk.NewInt64Coin("uctk", 123), sdk.NewInt64Coin("eth", 0)}

	cases := []struct {
		name       string
		expectPass bool
		msg        *types.MsgReduceCollateral
	}{
		{"valid with one denom", true, types.NewMsgReduceCollateral(addr, uctk123)},
		{"valid with two denoms", true, types.NewMsgReduceCollateral(addr, uctk123eth123)},
		{"non-positive coin", false, types.NewMsgReduceCollateral(addr, uctk0)},
		{"non-positive multicoins", false, types.NewMsgReduceCollateral(addr, uctk123eth0)},
		{"empty address", false, types.NewMsgReduceCollateral(addrEmpty, uctk123)},
	}

	for _, tc := range cases {
		if tc.expectPass {
			require.Nil(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

// test ValidateBasic for NewMsgWithdrawReward
func Test_NewMsgWithdrawReward(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr1_______________"))
	addrEmpty := sdk.AccAddress([]byte(""))

	cases := []struct {
		name       string
		expectPass bool
		msg        *types.MsgWithdrawReward
	}{
		{"valid with one denom", true, types.NewMsgWithdrawReward(addr)},
		{"empty address", false, types.NewMsgWithdrawReward(addrEmpty)},
	}

	for _, tc := range cases {
		if tc.expectPass {
			require.Nil(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.Error(t, tc.msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}
