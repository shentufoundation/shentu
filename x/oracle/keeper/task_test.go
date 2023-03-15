package keeper_test

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/oracle"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func TestTaskBasic(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	tth := TTHelper{
		app: app, ctx: ctx, t: t, creator: addrs[0],
		bounty:       sdk.Coins{sdk.NewInt64Coin("uctk", 100000)},
		contract:     "0x1234567890abcdef",
		function:     "func1",
		desc:         "testing",
		waitingBlock: int64(5),
		expiration:   time.Now().Add(time.Hour).UTC(),
		validTime:    time.Time{},
	}
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTask()))
	tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	task1ID := tth.TaskID()

	tth.contract, tth.function = "0x1234567890fedcba", "func2"
	tth.expiration = time.Now().Add(time.Hour * 2).UTC()
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTask()))
	task2Res := tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	task2ID := tth.TaskID()
	task2Expiration := tth.expiration

	tasks := ok.GetAllTasks(ctx)
	require.Len(t, tasks, 2)

	require.Error(t, ok.RemoveTask(ctx, task1ID, false, addrs[0]))
	require.Error(t, ok.RemoveTask(ctx, task2ID, false, addrs[0]))

	ctx = ctx.WithBlockTime(task2Expiration)
	require.Error(t, ok.RemoveTask(ctx, task1ID, false, addrs[0]))
	require.Error(t, ok.RemoveTask(ctx, task2ID, false, addrs[0]))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 6)
	require.NoError(t, ok.RemoveTask(ctx, task1ID, false, addrs[0]))
	require.Error(t, ok.RemoveTask(ctx, task2ID, false, addrs[0]))

	tasks = ok.GetAllTasks(ctx)
	require.Len(t, tasks, 1)
	require.Equal(t, []types.TaskI{task2Res}, ok.GetAllTasks(ctx))
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

	tth := TTHelper{
		app: app, ctx: ctx, t: t, creator: addrs[0],
		bounty:       sdk.Coins{sdk.NewInt64Coin("uctk", 100000)},
		contract:     "0x1234567890abcdef",
		function:     "func",
		desc:         "testing",
		waitingBlock: int64(5),
		expiration:   time.Now().Add(time.Hour).UTC(),
		validTime:    time.Time{},
	}
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTask()))
	taskRes := tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))

	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 100, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 100, addrs[2]))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 6)
	require.Error(t, ok.RespondToTask(ctx, tth.TaskID(), 100, addrs[0]))

	task, result := taskRes.(*types.Task)
	require.True(t, result)
	ok.UpdateAndSetTask(ctx, task)
	taskRes.SetStatus(types.TaskStatusFailed)
	ok.SetTask(ctx, task)
	require.Error(t, ok.Aggregate(ctx, tth.TaskID()))

	task.SetStatus(types.TaskStatusPending)
	ok.SetTask(ctx, task)
	require.NoError(t, ok.Aggregate(ctx, tth.TaskID()))
}

func TestTaskNoResponses(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	tth := TTHelper{
		app: app, ctx: ctx, t: t, creator: addrs[0],
		bounty:       sdk.Coins{sdk.NewInt64Coin("uctk", 100000)},
		contract:     "0x1234567890abcdef",
		function:     "func",
		desc:         "testing",
		waitingBlock: int64(5),
		expiration:   time.Now().Add(time.Hour).UTC(),
		validTime:    time.Time{},
	}
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTask()))
	require.NoError(t, ok.Aggregate(ctx, tth.TaskID()))
	taskRes := tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	require.Equal(t, types.TaskStatusFailed, taskRes.GetStatus())
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

	tth := TTHelper{
		app: app, ctx: ctx, t: t, creator: addrs[0],
		bounty:       sdk.Coins{sdk.NewInt64Coin("uctk", 100000)},
		contract:     "0x1234567890abcdef",
		function:     "func",
		desc:         "testing",
		waitingBlock: int64(5),
		expiration:   time.Now().Add(time.Hour).UTC(),
		validTime:    time.Time{},
	}
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTask()))

	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 100, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 0, addrs[2]))
	require.NoError(t, ok.Aggregate(ctx, tth.TaskID()))

	taskRes := tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	require.Equal(t, types.TaskStatusSucceeded, taskRes.GetStatus())
	require.Equal(t, params.MinimumCollateral, taskRes.GetScore())

	require.NoError(t, ok.DistributeBounty(ctx, taskRes))

	operator1, err := ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator1.Address)
	require.Equal(t, tth.bounty, operator1.AccumulatedRewards)

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

	tth := TTHelper{
		app: app, ctx: ctx, t: t, creator: addrs[0],
		bounty:       sdk.Coins{sdk.NewInt64Coin("uctk", 100000)},
		contract:     "0x1234567890abcdef",
		function:     "func",
		desc:         "testing",
		waitingBlock: int64(5),
		expiration:   time.Now().Add(time.Hour).UTC(),
		validTime:    time.Time{},
	}
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTask()))
	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 40, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 20, addrs[2]))
	require.NoError(t, ok.Aggregate(ctx, tth.TaskID()))

	taskRes := tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	require.Equal(t, types.TaskStatusSucceeded, taskRes.GetStatus())
	require.Equal(t, int64(30), taskRes.GetScore())

	require.NoError(t, ok.DistributeBounty(ctx, taskRes))

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

	tth := TTHelper{
		app: app, ctx: ctx, t: t, creator: addrs[0],
		bounty:       sdk.Coins{sdk.NewInt64Coin("uctk", 100000)},
		contract:     "0x1234567890abcdef",
		function:     "func",
		desc:         "testing",
		waitingBlock: int64(5),
		expiration:   time.Now().Add(time.Hour).UTC(),
		validTime:    time.Time{},
	}
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTask()))

	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 100, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 60, addrs[2]))
	require.NoError(t, ok.Aggregate(ctx, tth.TaskID()))

	taskRes := tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	require.Equal(t, types.TaskStatusSucceeded, taskRes.GetStatus())
	require.Equal(t, int64(80), taskRes.GetScore())

	require.NoError(t, ok.DistributeBounty(ctx, taskRes))

	operator1, err := ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator1.Address)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 58333)}, operator1.AccumulatedRewards)

	operator2, err := ok.GetOperator(ctx, addrs[2])
	require.Nil(t, err)
	require.Equal(t, addrs[2].String(), operator2.Address)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin("uctk", 41666)}, operator2.AccumulatedRewards)
}

func TestTxTaskBasic(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	tth := TTHelper{
		app: app, ctx: ctx, t: t, creator: addrs[0],
		bounty:    sdk.Coins{sdk.NewInt64Coin("uctk", 100000)},
		contract:  "hello",
		validTime: time.Now().Add(time.Hour).UTC(),
	}
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTxTask(types.TaskStatusPending)))
	tth.CheckTxTask(ok.GetTask(ctx, tth.TxTaskID()))

	tth.contract = "hello world"
	tth.validTime = time.Now().Add(time.Hour * 2).UTC()
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTxTask(types.TaskStatusNil)))
	txtask2 := tth.CheckTxTask(ok.GetTask(ctx, tth.TxTaskID()))

	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTxTask(types.TaskStatusPending)))

	_ = ok.DeleteTask(ctx, txtask2)
	_, err := ok.GetTask(ctx, tth.TxTaskID())
	require.Error(t, err)
}

func TestTimer(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 5, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper
	initTime := ctx.BlockTime()

	tth := TTHelper{
		app: app, ctx: ctx, t: t, creator: addrs[0],
		bounty:       sdk.Coins{sdk.NewInt64Coin("uctk", 100000)},
		contract:     "0x1234567890abcdef",
		function:     "func",
		desc:         "testing",
		waitingBlock: int64(5),
		expiration:   time.Now().Add(time.Hour).UTC(),
		validTime:    time.Now().Add(time.Second * 60).UTC(),
	}
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTask()))
	tth.creator = addrs[1]
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTxTask(types.TaskStatusPending)))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 2)
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetAllTasks(ctx), 2)
	taskRes := tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	txTaskRes := tth.CheckTxTask(ok.GetTask(ctx, tth.TxTaskID()))
	require.Equal(t, types.TaskStatusPending, taskRes.GetStatus())
	require.Equal(t, types.TaskStatusPending, txTaskRes.GetStatus())

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)
	require.Len(t, ok.GetInvalidTaskIDs(ctx), 1)
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetInvalidTaskIDs(ctx), 0)
	taskRes = tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	txTaskRes = tth.CheckTxTask(ok.GetTask(ctx, tth.TxTaskID()))
	require.Equal(t, types.TaskStatusFailed, taskRes.GetStatus())
	require.Equal(t, types.TaskStatusPending, txTaskRes.GetStatus())

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 6)
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetAllTasks(ctx), 2)
	require.Len(t, ok.GetInvalidTaskIDs(ctx), 0)

	ctx = ctx.WithBlockTime(tth.validTime.Add(time.Second * 60))
	require.Len(t, ok.GetInvalidTaskIDs(ctx), 1)
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetAllTasks(ctx), 2)
	txTaskRes = tth.CheckTxTask(ok.GetTask(ctx, tth.TxTaskID()))
	require.Equal(t, types.TaskStatusFailed, txTaskRes.GetStatus())

	ctx = ctx.WithBlockTime(tth.validTime.Add(time.Second * 61))
	require.Len(t, ok.GetInvalidTaskIDs(ctx), 0)

	//the taskParams.ExpirationDuration is one day
	ctx = ctx.WithBlockTime(initTime.Add(time.Second * 86400))
	require.Len(t, ok.GetAllTasks(ctx), 2)
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetAllTasks(ctx), 1) //the tx task is supposed to be removed

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Second))
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetAllTasks(ctx), 1)
	taskRes = tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	require.Equal(t, types.TaskStatusFailed, taskRes.GetStatus())
	require.Equal(t, int64(50), taskRes.GetScore())

	tth.contract = "0x12345678909090abc"
	tth.creator = addrs[2]
	tth.ctx = ctx
	require.NoError(t, ok.CreateTask(tth.ctx, tth.creator, tth.GetTxTask(types.TaskStatusNil)))
	tth.ctx = tth.ctx.WithBlockHeight(tth.ctx.BlockHeight() + 2)
	require.NoError(t, ok.CreateTask(tth.ctx, tth.creator, tth.GetTxTask(types.TaskStatusPending)))
	tth.ctx = tth.ctx.WithBlockHeight(tth.ctx.BlockHeight() + 2)
	require.Error(t, ok.RemoveTask(tth.ctx, tth.TxTaskID(), true, addrs[3]))
	require.NoError(t, ok.RemoveTask(tth.ctx, tth.TxTaskID(), true, addrs[2]))
	require.Len(t, ok.GetAllTasks(tth.ctx), 1)
}

type TTHelper struct {
	//for both Task and TxTask
	app     *shentuapp.ShentuApp
	ctx     sdk.Context
	t       *testing.T
	creator sdk.AccAddress
	bounty  sdk.Coins
	//only for Task
	contract     string
	function     string
	desc         string
	waitingBlock int64
	expiration   time.Time
	//only for TxTask
	validTime time.Time
}

func (t *TTHelper) GetTask() *types.Task {
	tk := types.NewTask(
		t.contract, t.function, t.ctx.BlockHeight(),
		t.bounty, t.desc, t.expiration,
		t.creator, t.ctx.BlockHeight()+t.waitingBlock,
		t.waitingBlock,
	)
	return &tk
}

func (t *TTHelper) TaskID() []byte {
	return types.NewTaskID(t.contract, t.function)
}

func (t *TTHelper) CheckTask(i types.TaskI, err error) types.TaskI {
	require.Nil(t.t, err)
	res := i.(*types.Task)
	require.Equal(t.t, t.contract, res.Contract)
	require.Equal(t.t, t.function, res.Function)
	require.Equal(t.t, t.expiration, res.Expiration)
	require.Equal(t.t, t.desc, res.Description)
	return i
}

func (t *TTHelper) GetTxTask(status types.TaskStatus) *types.TxTask {
	return t.app.OracleKeeper.BuildTxTaskWithExpire(
		t.ctx, t.TxTaskID(), t.creator.String(),
		t.bounty, t.validTime, status)
}

func (t *TTHelper) TxTaskID() []byte {
	txHash := sha256.Sum256([]byte(t.contract))
	return txHash[:]
}

func (t *TTHelper) CheckTxTask(i types.TaskI, err error) types.TaskI {
	require.Nil(t.t, err)
	res := i.(*types.TxTask)
	require.Equal(t.t, t.creator, sdk.MustAccAddressFromBech32(res.Creator))
	require.Equal(t.t, t.bounty, res.Bounty)
	require.Equal(t.t, t.TxTaskID(), res.GetID())
	require.Equal(t.t, t.validTime, res.ValidTime)
	return i
}
