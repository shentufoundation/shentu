package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/assert"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// test OperatorStoreKey
func Test_OperatorStoreKey(t *testing.T) {
	t.Run("Operator", func(t *testing.T) {
		acc := sdk.AccAddress([]byte{10})
		tmp := types.OperatorStoreKey(acc)
		assert.Equal(t, tmp, []byte{1, 10})
	})
}

// test WithdrawStoreKey
func Test_WithdrawStoreKey(t *testing.T) {
	t.Run("Withdraw", func(t *testing.T) {
		acc := sdk.AccAddress([]byte{10})
		var n int64 = 34
		tmp := types.WithdrawStoreKey(acc, n)
		assert.Equal(t, tmp, []byte{2, 34, 0, 0, 0, 0, 0, 0, 0, 10})
	})
}

func Test_TotalCollateralKey(t *testing.T) {
	t.Run("TotalCollateral", func(t *testing.T) {
		tmp := types.TotalCollateralKey()
		assert.Equal(t, tmp, []byte{3})
	})
}

// test TaskStoreKey
func Test_TaskStoreKey(t *testing.T) {
	t.Run("Task", func(t *testing.T) {
		s1 := "abc"
		s2 := "ghj"
		tmp := types.TaskStoreKey(s1, s2)
		assert.Equal(t, tmp, []byte{4, 97, 98, 99, 103, 104, 106})
	})
}

// test ClosingTaskIDsStoreKey
func Test_ClosingTaskIDsStoreKey(t *testing.T) {
	t.Run("ClosingTaskIDs", func(t *testing.T) {
		var n int64 = 34
		tmp := types.ClosingTaskIDsStoreKey(n)
		assert.Equal(t, tmp, []byte{5, 34, 0, 0, 0, 0, 0, 0, 0})
	})
}
