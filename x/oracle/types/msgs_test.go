package types_test

import (
	"testing"

	"github.com/test-go/testify/require"

	"github.com/certikfoundation/shentu/v2/x/oracle/types"
)

// test ValidateBasic for NewMsgCreateOperator
func Test_NewMsgCreateOperator(t *testing.T) {
	tests := []struct {
		name    string
		msg     *types.MsgCreateOperator
		expPass bool
	}{
		{
			"valid with one denom",
			types.NewMsgCreateOperator(acc1, coins1234, acc2, "operator"),
			true,
		},
		{
			"valid with two denoms",
			types.NewMsgCreateOperator(acc1, multicoins1234, acc2, "operator"),
			true,
		},
		{
			"non-positive coin",
			types.NewMsgCreateOperator(acc1, coins0, acc2, "operator"),
			false,
		},
		{
			"non-positive multicoins",
			types.NewMsgCreateOperator(acc1, multicoins0, acc2, "operator"),
			false,
		},
		{
			"invalid operator address",
			types.NewMsgCreateOperator(emptyAcc, coins1234, acc2, "operator"),
			false,
		},
		{
			"invalid proposer address",
			types.NewMsgCreateOperator(acc1, coins1234, emptyAcc, "operator"),
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expPass {
				require.NoError(t, tc.msg.ValidateBasic())
			} else {
				require.Error(t, tc.msg.ValidateBasic())
			}
		})
	}
}

// test ValidateBasic for NewMsgRemoveOperator
func Test_NewMsgRemoveOperator(t *testing.T) {
	tests := []struct {
		name    string
		msg     *types.MsgRemoveOperator
		expPass bool
	}{
		{
			"valid addresses",
			types.NewMsgRemoveOperator(acc1, acc2),
			true,
		},
		{
			"invalid operator address",
			types.NewMsgRemoveOperator(emptyAcc, acc2),
			false,
		},
		{
			"invalid proposer address",
			types.NewMsgRemoveOperator(acc1, emptyAcc),
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expPass {
				require.Nil(t, tc.msg.ValidateBasic())
			} else {
				require.Error(t, tc.msg.ValidateBasic())
			}
		})
	}
}

// test ValidateBasic for NewMsgAddCollateral
func Test_NewMsgAddCollateral(t *testing.T) {
	tests := []struct {
		name    string
		msg     *types.MsgAddCollateral
		expPass bool
	}{
		{
			"valid with one denom",
			types.NewMsgAddCollateral(acc1, coins1234),
			true,
		},
		{
			"valid with two denoms",
			types.NewMsgAddCollateral(acc1, multicoins1234),
			true,
		},
		{
			"non-positive coin",
			types.NewMsgAddCollateral(acc1, coins0),
			false,
		},
		{
			"non-positive multicoins",
			types.NewMsgAddCollateral(acc1, multicoins0),
			false,
		},
		{
			"empty address",
			types.NewMsgAddCollateral(emptyAcc, coins1234),
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expPass {
				require.Nil(t, tc.msg.ValidateBasic())
			} else {
				require.Error(t, tc.msg.ValidateBasic())
			}
		})
	}
}

// test ValidateBasic for NewMsgReduceCollateral
func Test_NewMsgReduceCollateral(t *testing.T) {
	tests := []struct {
		name    string
		msg     *types.MsgReduceCollateral
		expPass bool
	}{
		{
			"valid with one denom",
			types.NewMsgReduceCollateral(acc1, coins1234),
			true,
		},
		{
			"valid with two denoms",
			types.NewMsgReduceCollateral(acc1, multicoins1234),
			true,
		},
		{
			"non-positive coin",
			types.NewMsgReduceCollateral(acc1, coins0),
			false,
		},
		{
			"non-positive multicoins",
			types.NewMsgReduceCollateral(acc1, multicoins0),
			false,
		},
		{
			"empty address",
			types.NewMsgReduceCollateral(emptyAcc, coins1234),
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expPass {
				require.Nil(t, tc.msg.ValidateBasic())
			} else {
				require.Error(t, tc.msg.ValidateBasic())
			}
		})
	}
}

// test ValidateBasic for NewMsgWithdrawReward
func Test_NewMsgWithdrawReward(t *testing.T) {
	tests := []struct {
		name    string
		msg     *types.MsgWithdrawReward
		expPass bool
	}{
		{
			"valid address",
			types.NewMsgWithdrawReward(acc1),
			true,
		},
		{
			"empty address",
			types.NewMsgWithdrawReward(emptyAcc),
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expPass {
				require.Nil(t, tc.msg.ValidateBasic())
			} else {
				require.Error(t, tc.msg.ValidateBasic())
			}
		})
	}
}
