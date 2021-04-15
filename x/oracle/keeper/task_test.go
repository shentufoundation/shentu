package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/oracle/types"
)

func TestTaskBasic(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	description := "testing"
	waitingBlocks := int64(5)

	contract1 := "0x1234567890abcdef"
	function1 := "func1"
	expiration1 := time.Now().Add(time.Hour).UTC()
	require.NoError(t, ok.CreateTask(ctx, contract1, function1, bounty, description, expiration1, addrs[0], waitingBlocks))

	task1, err := ok.GetTask(ctx, contract1, function1)
	require.Nil(t, err)
	require.Equal(t, contract1, task1.Contract)
	require.Equal(t, function1, task1.Function)
	require.Equal(t, expiration1, task1.Expiration)

	contract2 := "0x1234567890fedcba"
	function2 := "func2"
	expiration2 := time.Now().Add(time.Hour * 2).UTC()
	require.NoError(t, ok.CreateTask(ctx, contract2, function2, bounty, description, expiration2, addrs[0], waitingBlocks))

	task2, err := ok.GetTask(ctx, contract2, function2)
	require.Nil(t, err)
	require.Equal(t, contract2, task2.Contract)
	require.Equal(t, function2, task2.Function)
	require.Equal(t, expiration2, task2.Expiration)

	tasks := ok.GetAllTasks(ctx)
	require.Len(t, tasks, 2)

	require.Error(t, ok.RemoveTask(ctx, contract1, function1, false, addrs[0]))
	require.Error(t, ok.RemoveTask(ctx, contract2, function2, false, addrs[0]))

	ctx = ctx.WithBlockTime(expiration2)
	require.Error(t, ok.RemoveTask(ctx, contract1, function1, false, addrs[0]))
	require.Error(t, ok.RemoveTask(ctx, contract2, function2, false, addrs[0]))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 6)
	require.NoError(t, ok.RemoveTask(ctx, contract1, function1, false, addrs[0]))
	require.Error(t, ok.RemoveTask(ctx, contract2, function2, false, addrs[0]))

	tasks = ok.GetAllTasks(ctx)
	require.Len(t, tasks, 1)
	require.Equal(t, []types.Task{task2}, ok.GetAllTasks(ctx))
}

func TestTaskAggregateFail(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)
	collateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}

	require.NoError(t, ok.CreateOperator(ctx, addrs[0], collateral, addrs[1], "operator1"))
	require.NoError(t, ok.CreateOperator(ctx, addrs[2], collateral, addrs[3], "operator2"))

	contract := "0x1234567890abcdef"
	function := "func"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)

	require.NoError(t, ok.CreateTask(ctx, contract, function, bounty, description, expiration, addrs[0], waitingBlocks))

	task, err := ok.GetTask(ctx, contract, function)
	require.Nil(t, err)
	require.Equal(t, contract, task.Contract)
	require.Equal(t, function, task.Function)

	require.NoError(t, ok.RespondToTask(ctx, contract, function, 100, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, contract, function, 100, addrs[2]))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 6)
	require.Error(t, ok.RespondToTask(ctx, contract, function, 100, addrs[0]))

	ok.UpdateAndSetTask(ctx, task)
	task.Status = types.TaskStatusFailed
	ok.SetTask(ctx, task)
	require.Error(t, ok.Aggregate(ctx, contract, function))

	task.Status = types.TaskStatusPending
	ok.SetTask(ctx, task)
	require.NoError(t, ok.Aggregate(ctx, contract, function))
}

func TestTaskNoResponses(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	contract := "0x1234567890abcdef"
	function := "func"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)

	require.NoError(t, ok.CreateTask(ctx, contract, function, bounty, description, expiration, addrs[0], waitingBlocks))
	require.NoError(t, ok.Aggregate(ctx, contract, function))

	task, err := ok.GetTask(ctx, contract, function)
	require.Nil(t, err)
	require.Equal(t, contract, task.Contract)
	require.Equal(t, function, task.Function)
	require.Equal(t, types.TaskStatusFailed, task.Status)
}

func TestTaskMinScore(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)
	collateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}

	require.NoError(t, ok.CreateOperator(ctx, addrs[0], collateral, addrs[1], "operator1"))
	require.NoError(t, ok.CreateOperator(ctx, addrs[2], collateral, addrs[3], "operator2"))

	contract := "0x1234567890abcdef"
	function := "func"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)

	require.NoError(t, ok.CreateTask(ctx, contract, function, bounty, description, expiration, addrs[0], waitingBlocks))

	require.NoError(t, ok.RespondToTask(ctx, contract, function, 100, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, contract, function, 0, addrs[2]))

	require.NoError(t, ok.Aggregate(ctx, contract, function))

	task, err := ok.GetTask(ctx, contract, function)
	require.Nil(t, err)
	require.Equal(t, contract, task.Contract)
	require.Equal(t, function, task.Function)
	require.Equal(t, types.TaskStatusSucceeded, task.Status)
	require.Equal(t, sdk.NewInt(params.MinimumCollateral), task.Result)

	require.NoError(t, ok.DistributeBounty(ctx, task))

	operator1, err := ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator1.Address)
	require.Equal(t, bounty, operator1.AccumulatedRewards)

	operator2, err := ok.GetOperator(ctx, addrs[2])
	require.Nil(t, err)
	require.Equal(t, addrs[2].String(), operator2.Address)
	require.Nil(t, operator2.AccumulatedRewards)
}

func TestTaskBelowThreshold(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)
	collateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}

	require.NoError(t, ok.CreateOperator(ctx, addrs[0], collateral, addrs[1], "operator1"))
	require.NoError(t, ok.CreateOperator(ctx, addrs[2], collateral, addrs[3], "operator2"))

	contract := "0x1234567890abcdef"
	function := "func"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)

	require.NoError(t, ok.CreateTask(ctx, contract, function, bounty, description, expiration, addrs[0], waitingBlocks))

	require.NoError(t, ok.RespondToTask(ctx, contract, function, 40, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, contract, function, 20, addrs[2]))

	require.NoError(t, ok.Aggregate(ctx, contract, function))

	task, err := ok.GetTask(ctx, contract, function)
	require.Nil(t, err)
	require.Equal(t, contract, task.Contract)
	require.Equal(t, function, task.Function)
	require.Equal(t, types.TaskStatusSucceeded, task.Status)
	require.Equal(t, sdk.NewInt(30), task.Result)

	require.NoError(t, ok.DistributeBounty(ctx, task))

	operator1, err := ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator1.Address)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 33870)}, operator1.AccumulatedRewards)

	operator2, err := ok.GetOperator(ctx, addrs[2])
	require.Nil(t, err)
	require.Equal(t, addrs[2].String(), operator2.Address)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 66129)}, operator2.AccumulatedRewards)
}

func TestTaskAboveThreshold(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)
	collateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}

	require.NoError(t, ok.CreateOperator(ctx, addrs[0], collateral, addrs[1], "operator1"))
	require.NoError(t, ok.CreateOperator(ctx, addrs[2], collateral, addrs[3], "operator2"))

	contract := "0x1234567890abcdef"
	function := "func"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)

	require.NoError(t, ok.CreateTask(ctx, contract, function, bounty, description, expiration, addrs[0], waitingBlocks))

	require.NoError(t, ok.RespondToTask(ctx, contract, function, 100, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, contract, function, 60, addrs[2]))

	require.NoError(t, ok.Aggregate(ctx, contract, function))

	task, err := ok.GetTask(ctx, contract, function)
	require.Nil(t, err)
	require.Equal(t, contract, task.Contract)
	require.Equal(t, function, task.Function)
	require.Equal(t, types.TaskStatusSucceeded, task.Status)
	require.Equal(t, sdk.NewInt(80), task.Result)

	require.NoError(t, ok.DistributeBounty(ctx, task))

	operator1, err := ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator1.Address)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 58333)}, operator1.AccumulatedRewards)

	operator2, err := ok.GetOperator(ctx, addrs[2])
	require.Nil(t, err)
	require.Equal(t, addrs[2].String(), operator2.Address)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 41666)}, operator2.AccumulatedRewards)
}
