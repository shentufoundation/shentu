package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func TestTaskBasic(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	description := "testing"
	waitingBlocks := int64(5)

	contract1 := "0x1234567890abcdef"
	function1 := "func1"
	expiration1 := time.Now().Add(time.Hour).UTC()
	scTask := types.NewTask(
		contract1, function1, ctx.BlockHeight(),
		bounty, description, expiration1,
		addrs[0], ctx.BlockHeight()+waitingBlocks, waitingBlocks)
	require.NoError(t, ok.CreateTask(ctx, addrs[0], &scTask))

	task1, err := ok.GetTask(ctx, types.NewTaskID(contract1, function1))
	scTaskRes, castOK := task1.(*types.Task)
	require.True(t, castOK)
	require.Nil(t, err)
	require.Equal(t, contract1, scTaskRes.Contract)
	require.Equal(t, function1, scTaskRes.Function)
	require.Equal(t, expiration1, scTaskRes.Expiration)

	contract2 := "0x1234567890fedcba"
	function2 := "func2"
	expiration2 := time.Now().Add(time.Hour * 2).UTC()
	scTask = types.NewTask(
		contract2, function2, ctx.BlockHeight(),
		bounty, description, expiration2,
		addrs[0], ctx.BlockHeight()+waitingBlocks, waitingBlocks)
	require.NoError(t, ok.CreateTask(ctx, addrs[0], &scTask))

	task2, err := ok.GetTask(ctx, types.NewTaskID(contract2, function2))
	scTaskRes, castOK = task2.(*types.Task)
	require.True(t, castOK)
	require.Nil(t, err)
	require.Equal(t, contract2, scTaskRes.Contract)
	require.Equal(t, function2, scTaskRes.Function)
	require.Equal(t, expiration2, scTaskRes.Expiration)

	tasks := ok.GetAllTasks(ctx)
	require.Len(t, tasks, 2)

	require.Error(t, ok.RemoveTask(ctx, types.NewTaskID(contract1, function1), false, addrs[0]))
	require.Error(t, ok.RemoveTask(ctx, types.NewTaskID(contract2, function2), false, addrs[0]))

	ctx = ctx.WithBlockTime(expiration2)
	require.Error(t, ok.RemoveTask(ctx, types.NewTaskID(contract1, function1), false, addrs[0]))
	require.Error(t, ok.RemoveTask(ctx, types.NewTaskID(contract2, function2), false, addrs[0]))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 6)
	require.NoError(t, ok.RemoveTask(ctx, types.NewTaskID(contract1, function1), false, addrs[0]))
	require.Error(t, ok.RemoveTask(ctx, types.NewTaskID(contract2, function2), false, addrs[0]))

	tasks = ok.GetAllTasks(ctx)
	require.Len(t, tasks, 1)
	var returnedScTasks []types.Task
	tasks = ok.GetAllTasks(ctx)
	for _, t := range tasks {
		returnedScTasks = append(returnedScTasks, *t.(*types.Task))
	}
	require.Equal(t, []types.Task{*scTaskRes}, returnedScTasks)
}

func TestTaskAggregateFail(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
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

	scTask := types.NewTask(
		contract, function, ctx.BlockHeight(),
		bounty, description, expiration,
		addrs[0], ctx.BlockHeight()+waitingBlocks, waitingBlocks)
	require.NoError(t, ok.CreateTask(ctx, addrs[0], &scTask))

	task, err := ok.GetTask(ctx, types.NewTaskID(contract, function))
	scTaskRes := task.(*types.Task)
	require.Nil(t, err)
	require.Equal(t, contract, scTaskRes.Contract)
	require.Equal(t, function, scTaskRes.Function)

	require.NoError(t, ok.RespondToTask(ctx, types.NewTaskID(contract, function), 100, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, types.NewTaskID(contract, function), 100, addrs[2]))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 6)
	require.Error(t, ok.RespondToTask(ctx, types.NewTaskID(contract, function), 100, addrs[0]))

	ok.UpdateAndSetTask(ctx, task)
	task.SetStatus(types.TaskStatusFailed)
	ok.SetTask(ctx, task)
	require.Error(t, ok.Aggregate(ctx, types.NewTaskID(contract, function)))

	task.SetStatus(types.TaskStatusPending)
	ok.SetTask(ctx, task)
	require.NoError(t, ok.Aggregate(ctx, types.NewTaskID(contract, function)))
}

func TestTaskNoResponses(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	contract := "0x1234567890abcdef"
	function := "func"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)

	scTask := types.NewTask(
		contract, function, ctx.BlockHeight(),
		bounty, description, expiration,
		addrs[0], ctx.BlockHeight()+waitingBlocks, waitingBlocks)
	require.NoError(t, ok.CreateTask(ctx, addrs[0], &scTask))
	require.NoError(t, ok.Aggregate(ctx, types.NewTaskID(contract, function)))

	task, err := ok.GetTask(ctx, types.NewTaskID(contract, function))
	require.Nil(t, err)
	scTaskRes := task.(*types.Task)
	require.Equal(t, contract, scTaskRes.Contract)
	require.Equal(t, function, scTaskRes.Function)
	require.Equal(t, types.TaskStatusFailed, scTaskRes.Status)
}

func TestTaskMinScore(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
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

	scTask := types.NewTask(
		contract, function, ctx.BlockHeight(),
		bounty, description, expiration,
		addrs[0], ctx.BlockHeight()+waitingBlocks, waitingBlocks)
	require.NoError(t, ok.CreateTask(ctx, addrs[0], &scTask))

	require.NoError(t, ok.RespondToTask(ctx, types.NewTaskID(contract, function), 100, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, types.NewTaskID(contract, function), 0, addrs[2]))

	require.NoError(t, ok.Aggregate(ctx, types.NewTaskID(contract, function)))

	task, err := ok.GetTask(ctx, types.NewTaskID(contract, function))
	require.Nil(t, err)
	scTaskRes := task.(*types.Task)
	require.Equal(t, contract, scTaskRes.Contract)
	require.Equal(t, function, scTaskRes.Function)
	require.Equal(t, types.TaskStatusSucceeded, scTaskRes.Status)
	require.Equal(t, sdk.NewInt(params.MinimumCollateral), scTaskRes.Result)

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
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
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

	scTask := types.NewTask(
		contract, function, ctx.BlockHeight(),
		bounty, description, expiration,
		addrs[0], ctx.BlockHeight()+waitingBlocks, waitingBlocks)
	require.NoError(t, ok.CreateTask(ctx, addrs[0], &scTask))

	require.NoError(t, ok.RespondToTask(ctx, types.NewTaskID(contract, function), 40, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, types.NewTaskID(contract, function), 20, addrs[2]))

	require.NoError(t, ok.Aggregate(ctx, types.NewTaskID(contract, function)))

	task, err := ok.GetTask(ctx, types.NewTaskID(contract, function))
	require.Nil(t, err)
	scTaskRes := task.(*types.Task)
	require.Equal(t, contract, scTaskRes.Contract)
	require.Equal(t, function, scTaskRes.Function)
	require.Equal(t, types.TaskStatusSucceeded, scTaskRes.Status)
	require.Equal(t, sdk.NewInt(30), scTaskRes.Result)

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
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
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

	scTask := types.NewTask(
		contract, function, ctx.BlockHeight(),
		bounty, description, expiration,
		addrs[0], ctx.BlockHeight()+waitingBlocks, waitingBlocks)
	require.NoError(t, ok.CreateTask(ctx, addrs[0], &scTask))

	require.NoError(t, ok.RespondToTask(ctx, types.NewTaskID(contract, function), 100, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, types.NewTaskID(contract, function), 60, addrs[2]))

	require.NoError(t, ok.Aggregate(ctx, types.NewTaskID(contract, function)))

	task, err := ok.GetTask(ctx, types.NewTaskID(contract, function))
	require.Nil(t, err)
	scTaskRes := task.(*types.Task)
	require.Equal(t, contract, scTaskRes.Contract)
	require.Equal(t, function, scTaskRes.Function)
	require.Equal(t, types.TaskStatusSucceeded, scTaskRes.Status)
	require.Equal(t, sdk.NewInt(80), scTaskRes.Result)

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
