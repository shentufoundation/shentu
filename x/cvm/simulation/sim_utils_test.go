package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

var denom = "uctk"
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func MakeCoins(amount int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(amount)))
}

func CallRandomRFees(t *testing.T, amount int64) sdk.Coins {
	res, err := RandomReasonableFees(r, sdk.Context{}, MakeCoins(amount))
	require.NoError(t, err)
	return res
}

func FeeAfterRandomLessThan(t *testing.T, orig int64, target int64) {
	res := CallRandomRFees(t, int64(orig))
	require.True(t, res[0].Amount.LT(sdk.NewInt(target)))
}

func TestMigrateStore(t *testing.T) {
	require.Equal(t, sdk.Coins(nil), CallRandomRFees(t, 0))
	require.Equal(t, MakeCoins(1), CallRandomRFees(t, 1))
	require.Equal(t, MakeCoins(1), CallRandomRFees(t, 8))
	FeeAfterRandomLessThan(t, 9, 2)
	FeeAfterRandomLessThan(t, 81, 9)
	FeeAfterRandomLessThan(t, 1789034678, 198781631)
	FeeAfterRandomLessThan(t, 2000, 223)
}
