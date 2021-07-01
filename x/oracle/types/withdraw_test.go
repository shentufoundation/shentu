package types_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/assert"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

func Test_withdraw(t *testing.T) {
	t.Run("test ", func(t *testing.T) {
		acc := sdk.AccAddress([]byte{10})
		amount := sdk.Coins{sdk.NewInt64Coin("uctk", 1234)}
		dueBlock := int64(50)
		wd := types.NewWithdraw(acc, amount, dueBlock)
		s := wd.String()
		wds := []types.Withdraw{wd, wd}
		var wdss types.Withdraws = wds

		assert.Equal(t, wdss.String(), strings.TrimSpace(s+"\n"+s))
	})
}
