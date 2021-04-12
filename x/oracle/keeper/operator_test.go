package keeper_test

import (
	"testing"
	"time"

	"github.com/certikfoundation/shentu/simapp"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

func TestOperator_Basic(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)

	err := ok.CreateOperator(ctx, addrs[1], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[0], "operator1")
	require.Nil(t, err)

	operator1, err := ok.GetOperator(ctx, addrs[1])
	require.Nil(t, err)

	err = ok.CreateOperator(ctx, addrs[2], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[0], "operator2")
	require.Nil(t, err)

	operator2, err := ok.GetOperator(ctx, addrs[2])
	require.Nil(t, err)

	require.Equal(t, types.Operators{operator1, operator2}, ok.GetAllOperators(ctx))

	err = ok.RemoveOperator(ctx, addrs[1])
	require.Nil(t, err)

	require.Equal(t, types.Operators{operator2}, ok.GetAllOperators(ctx))
}

func TestOperator_Collateral(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)

	err := ok.CreateOperator(ctx, addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[1], "operator1")
	require.Nil(t, err)

	collateral, err := ok.GetCollateralAmount(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, sdk.NewInt(params.MinimumCollateral), collateral)

	err = ok.AddCollateral(ctx, addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", 100)})
	require.Nil(t, err)

	collateral, err = ok.GetCollateralAmount(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, sdk.NewInt(params.MinimumCollateral+100), collateral)

	err = ok.ReduceCollateral(ctx, addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", 100)})
	require.Nil(t, err)

	collateral, err = ok.GetCollateralAmount(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, sdk.NewInt(params.MinimumCollateral), collateral)

	withdraws := ok.GetAllWithdraws(ctx)
	require.Equal(t, 1, len(withdraws))
}

func TestOperator_Reward(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)

	err := ok.CreateOperator(ctx, addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[1], "operator1")
	require.Nil(t, err)

	err = ok.AddReward(ctx, addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", 200)})
	require.Nil(t, err)

	operator, err := ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)

	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 200)}, operator.AccumulatedRewards)

	withdrawal, err := ok.WithdrawAllReward(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 200)}, withdrawal)
}
