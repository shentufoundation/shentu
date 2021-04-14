package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/certikfoundation/shentu/simapp"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

func TestTask_Basic(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	contract1 := "0x1234567890abcdef"
	function1 := "func1"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 5000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)

	err := ok.CreateTask(ctx, contract1, function1, bounty, description, expiration, addrs[0], waitingBlocks)
	require.Nil(t, err)

	task1, err := ok.GetTask(ctx, contract1, function1)
	require.Nil(t, err)

	contract2 := "0x1234567890fedcba"
	function2 := "func2"
	err = ok.CreateTask(ctx, contract2, function2, bounty, description, expiration, addrs[0], waitingBlocks)
	require.Nil(t, err)

	task2, err := ok.GetTask(ctx, contract2, function2)
	require.Nil(t, err)

	require.Equal(t, []types.Task{task1, task2}, ok.GetAllTasks(ctx))

	err = ok.RemoveTask(ctx, contract1, function1, true, addrs[0])
	require.Error(t, err)

	ctx = ctx.WithBlockHeight(6)
	err = ok.RemoveTask(ctx, contract1, function1, true, addrs[0])
	require.Nil(t, err)

	require.Equal(t, []types.Task{task2}, ok.GetAllTasks(ctx))
}

func TestTask_MinScore(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)

	err := ok.CreateOperator(ctx, addrs[1], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[0], "operator1")
	require.Nil(t, err)

	err = ok.CreateOperator(ctx, addrs[2], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[0], "operator1")
	require.Nil(t, err)

	contract := "0x1234567890abcdef"
	function := "func"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 5000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)

	err = ok.CreateTask(ctx, contract, function, bounty, description, expiration, addrs[0], waitingBlocks)
	require.Nil(t, err)

	err = ok.RespondToTask(ctx, contract, function, 100, addrs[1])
	require.Nil(t, err)

	err = ok.RespondToTask(ctx, contract, function, 0, addrs[2])
	require.Nil(t, err)

	err = ok.Aggregate(ctx, contract, function)
	require.Nil(t, err)

	task, err := ok.GetTask(ctx, contract, function)
	require.Nil(t, err)
	require.Equal(t, types.TaskStatusSucceeded, task.Status)
	require.Equal(t, sdk.NewInt(params.MinimumCollateral), task.Result)

	err = ok.DistributeBounty(ctx, task)
	require.Nil(t, err)

	operator1, err := ok.GetOperator(ctx, addrs[1])
	require.Nil(t, err)
	require.Equal(t, bounty, operator1.AccumulatedRewards)

	operator2, err := ok.GetOperator(ctx, addrs[2])
	require.Nil(t, err)
	require.Nil(t, operator2.AccumulatedRewards)
}

func TestTask_AboveThreshold(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)

	err := ok.CreateOperator(ctx, addrs[1], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[0], "operator1")
	require.Nil(t, err)

	err = ok.CreateOperator(ctx, addrs[2], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[0], "operator1")
	require.Nil(t, err)

	contract := "0x1234567890abcdef"
	function := "func"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 5000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)

	err = ok.CreateTask(ctx, contract, function, bounty, description, expiration, addrs[0], waitingBlocks)
	require.Nil(t, err)

	err = ok.RespondToTask(ctx, contract, function, 80, addrs[1])
	require.Nil(t, err)

	err = ok.RespondToTask(ctx, contract, function, 70, addrs[2])
	require.Nil(t, err)

	err = ok.Aggregate(ctx, contract, function)
	require.Nil(t, err)

	task, err := ok.GetTask(ctx, contract, function)
	require.Nil(t, err)
	require.Equal(t, types.TaskStatusSucceeded, task.Status)
	require.Equal(t, sdk.NewInt(75), task.Result)

	err = ok.DistributeBounty(ctx, task)
	require.Nil(t, err)

	operator1, err := ok.GetOperator(ctx, addrs[1])
	require.Nil(t, err)

	operator2, err := ok.GetOperator(ctx, addrs[2])
	require.Nil(t, err)

	require.Equal(t, bounty, operator1.AccumulatedRewards.Add(operator2.AccumulatedRewards[0]))
}
