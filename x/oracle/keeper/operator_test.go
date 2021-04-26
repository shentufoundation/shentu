package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/simapp"
)

func TestOperatorBasic(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)
	collateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}

	require.NoError(t, ok.CreateOperator(ctx, addrs[0], collateral, addrs[1], "operator1"))
	require.NoError(t, ok.CreateOperator(ctx, addrs[2], collateral, addrs[3], "operator2"))

	operator1, err := ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator1.Address)
	require.Equal(t, collateral, operator1.Collateral)
	require.Equal(t, addrs[1].String(), operator1.Proposer)

	operator2, err := ok.GetOperator(ctx, addrs[2])
	require.Nil(t, err)
	require.Equal(t, addrs[2].String(), operator2.Address)
	require.Equal(t, collateral, operator2.Collateral)
	require.Equal(t, addrs[3].String(), operator2.Proposer)

	operators := ok.GetAllOperators(ctx)
	require.Len(t, operators, 2)

	require.NoError(t, ok.RemoveOperator(ctx, addrs[0]))

	operators = ok.GetAllOperators(ctx)
	require.Len(t, operators, 1)
	require.Equal(t, operator2, operators[0])
}

func TestOperatorCollateral(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)
	minCollateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}
	hundred := sdk.Coins{sdk.NewInt64Coin("uctk", 100)}

	require.NoError(t, ok.CreateOperator(ctx, addrs[0], minCollateral, addrs[1], "operator1"))

	collateral, err := ok.GetCollateralAmount(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, sdk.NewInt(params.MinimumCollateral), collateral)

	require.NoError(t, ok.AddCollateral(ctx, addrs[0], hundred))

	collateral, err = ok.GetCollateralAmount(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, sdk.NewInt(params.MinimumCollateral+100), collateral)

	require.NoError(t, ok.ReduceCollateral(ctx, addrs[0], hundred))

	collateral, err = ok.GetCollateralAmount(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, sdk.NewInt(params.MinimumCollateral), collateral)

	require.Error(t, ok.ReduceCollateral(ctx, addrs[0], hundred))
	require.Error(t, ok.AddCollateral(ctx, addrs[1], hundred))

	withdraws := ok.GetAllWithdraws(ctx)
	require.Len(t, withdraws, 1)
}

func TestOperatorReward(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)
	minCollateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}

	require.NoError(t, ok.CreateOperator(ctx, addrs[0], minCollateral, addrs[1], "operator1"))

	operator, err := ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator.Address)
	require.Equal(t, minCollateral, operator.Collateral)
	require.Equal(t, addrs[1].String(), operator.Proposer)
	require.Nil(t, operator.AccumulatedRewards)

	require.NoError(t, ok.AddReward(ctx, addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", 100)}))
	require.NoError(t, ok.AddReward(ctx, addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", 200)}))
	require.Error(t, ok.AddReward(ctx, addrs[1], sdk.Coins{sdk.NewInt64Coin("uctk", 300)}))

	operator, err = ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator.Address)
	require.Equal(t, minCollateral, operator.Collateral)
	require.Equal(t, addrs[1].String(), operator.Proposer)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 300)}, operator.AccumulatedRewards)

	withdrawal, err := ok.WithdrawAllReward(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 300)}, withdrawal)

	operator, err = ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator.Address)
	require.Equal(t, minCollateral, operator.Collateral)
	require.Equal(t, addrs[1].String(), operator.Proposer)
	require.Nil(t, operator.AccumulatedRewards)
}
