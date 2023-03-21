package types_test

import (
	"time"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

var (
	acc1     = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	acc2     = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	emptyAcc = sdk.AccAddress([]byte(""))

	coins0         = sdk.Coins{sdk.NewInt64Coin("uctk", 0)}
	coins1234      = sdk.NewCoins(sdk.NewInt64Coin("uctk", 1234))
	coins5000      = sdk.NewCoins(sdk.NewInt64Coin("uctk", 50000))
	coins1e5       = sdk.NewCoins(sdk.NewInt64Coin("uctk", 1e5))
	multicoins1234 = sdk.NewCoins(sdk.NewInt64Coin("uctk", 1234), sdk.NewInt64Coin("eth", 1234))
	multicoins0    = sdk.Coins{sdk.NewInt64Coin("uctk", 1234), sdk.NewInt64Coin("eth", 0)}

	operator1 = types.NewOperator(acc1, acc1, coins5000, nil, "operator1")
	operator2 = types.NewOperator(acc2, acc2, coins5000, nil, "operator2")

	validPoolParams   = types.DefaultLockedPoolParams()
	invalidPoolParams = types.NewLockedPoolParams(int64(-1), types.DefaultMinimumCollateral)

	validTaskParams   = types.DefaultTaskParams()
	invalidTaskParams = types.NewTaskParams(
		time.Duration(-1),
		types.DefaultAggregationWindow,
		types.DefaultAggregationResult,
		types.DefaultThresholdScore,
		types.DefaultEpsilon1,
		types.DefaultEpsilon2,
	)

	validWithdraw   = types.NewWithdraw(acc1, coins1234, int64(100))
	invalidWithdraw = types.NewWithdraw(acc1, coins1234, int64(-1))

	validTask = types.NewTask(
		"0x1234567890abcdef",
		"func",
		int64(0),
		coins1234,
		"",
		time.Now().Add(time.Hour).UTC(),
		acc1,
		int64(100),
		int64(5),
	)

	validTxTask = types.NewTxTask(
		[]byte("testtxtask"),
		acc1.String(),
		coins5000,
		time.Now().Add(time.Hour).UTC(),
		types.TaskStatusPending,
	)
)
