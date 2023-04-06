package keeper_test

import (
	"bytes"
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

	ctx = PassBlocks(ctx, ok, t, 5, 2)
	require.NoError(t, ok.RemoveTask(ctx, task1ID, false, addrs[0]))

	tasks = ok.GetAllTasks(ctx)
	require.Len(t, tasks, 1)
	tt := ok.GetAllTasks(ctx)[0]
	require.Equal(t, task2ID, tt.GetID())
	require.Equal(t, task2Res.GetBounty(), tt.GetBounty())
	require.Equal(t, types.TaskStatusFailed, tt.GetStatus())
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

	require.Error(t, ok.RespondToTask(ctx, tth.TaskID(), 110, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 100, addrs[2]))
	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 100, addrs[0]))

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

	require.NoError(t, ok.RefundBounty(ctx, taskRes))
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
		validTime:    time.Now().Add(time.Second).UTC(),
	}
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetTask()))

	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 100, addrs[2]))
	require.Error(t, ok.RespondToTask(ctx, tth.TaskID(), -1, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, tth.TaskID(), 0, addrs[0]))

	require.NoError(t, ok.Aggregate(ctx, tth.TaskID()))

	taskRes := tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	require.Equal(t, types.TaskStatusSucceeded, taskRes.GetStatus())
	require.Equal(t, types.MinScore.Int64(), taskRes.GetScore())

	require.NoError(t, ok.DistributeBounty(ctx, taskRes))

	taskRes = tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	require.NoError(t, ok.RefundBounty(ctx, taskRes))

	operator1, err := ok.GetOperator(ctx, addrs[0])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), operator1.Address)
	require.Equal(t, tth.bounty, operator1.AccumulatedRewards)

	operator2, err := ok.GetOperator(ctx, addrs[2])
	require.Nil(t, err)
	require.Equal(t, addrs[2].String(), operator2.Address)
	require.Nil(t, operator2.AccumulatedRewards)

	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetAtxTask(types.TaskStatusPending)))

	require.NoError(t, ok.RespondToTask(ctx, tth.AtxTaskID(), 0, addrs[0]))
	require.NoError(t, ok.RespondToTask(ctx, tth.AtxTaskID(), 0, addrs[2]))
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Second * 6))
	oracle.EndBlocker(ctx, ok)
	atxTaskRes := tth.CheckAtxTask(ok.GetTask(ctx, tth.AtxTaskID()))
	require.Equal(t, types.TaskStatusSucceeded, atxTaskRes.GetStatus())
	require.Equal(t, types.MinScore.Int64(), atxTaskRes.GetScore())
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

	taskRes = tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	require.NoError(t, ok.RefundBounty(ctx, taskRes))

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

func TestAtxTaskBasic(t *testing.T) {
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
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetAtxTask(types.TaskStatusPending)))
	tth.CheckAtxTask(ok.GetTask(ctx, tth.AtxTaskID()))

	tth.contract = "hello world"
	tth.validTime = time.Now().Add(time.Hour * 2).UTC()
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetAtxTask(types.TaskStatusNil)))
	atxtask2 := tth.CheckAtxTask(ok.GetTask(ctx, tth.AtxTaskID()))

	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetAtxTask(types.TaskStatusPending)))

	_ = ok.DeleteTask(ctx, atxtask2)
	_, err := ok.GetTask(ctx, tth.AtxTaskID())
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
	require.NoError(t, ok.CreateTask(ctx, tth.creator, tth.GetAtxTask(types.TaskStatusPending)))

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 2)
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetAllTasks(ctx), 2)
	taskRes := tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	atxTaskRes := tth.CheckAtxTask(ok.GetTask(ctx, tth.AtxTaskID()))
	require.Equal(t, types.TaskStatusPending, taskRes.GetStatus())
	require.Equal(t, types.TaskStatusPending, atxTaskRes.GetStatus())

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 3)
	require.Len(t, ok.GetInvalidTaskIDs(ctx), 1)
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetInvalidTaskIDs(ctx), 0)
	taskRes = tth.CheckTask(ok.GetTask(ctx, tth.TaskID()))
	atxTaskRes = tth.CheckAtxTask(ok.GetTask(ctx, tth.AtxTaskID()))
	require.Equal(t, types.TaskStatusFailed, taskRes.GetStatus())
	require.Equal(t, types.TaskStatusPending, atxTaskRes.GetStatus())

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 6)
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetAllTasks(ctx), 2)
	require.Len(t, ok.GetInvalidTaskIDs(ctx), 0)

	ctx = ctx.WithBlockTime(tth.validTime.Add(time.Second * 60))
	require.Len(t, ok.GetInvalidTaskIDs(ctx), 1)
	oracle.EndBlocker(ctx, ok)
	require.Len(t, ok.GetAllTasks(ctx), 2)
	atxTaskRes = tth.CheckAtxTask(ok.GetTask(ctx, tth.AtxTaskID()))
	require.Equal(t, types.TaskStatusFailed, atxTaskRes.GetStatus())

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
	tth.validTime = ctx.BlockTime().Add(time.Second * 50).UTC()
	require.NoError(t, ok.CreateTask(tth.ctx, tth.creator, tth.GetAtxTask(types.TaskStatusNil)))
	tth.ctx = PassBlocks(tth.ctx, ok, t, 2, 0)
	require.NoError(t, ok.CreateTask(tth.ctx, tth.creator, tth.GetAtxTask(types.TaskStatusPending)))
	tth.ctx = PassBlocks(tth.ctx, ok, t, 2, 0)
	require.Error(t, ok.RemoveTask(tth.ctx, tth.AtxTaskID(), true, addrs[3])) // status is still pending
	tth.ctx = PassBlocks(tth.ctx, ok, t, 6, 1)
	require.Error(t, ok.RemoveTask(tth.ctx, tth.AtxTaskID(), true, addrs[3])) // not the creator
	require.NoError(t, ok.RemoveTask(tth.ctx, tth.AtxTaskID(), true, addrs[2]))
	require.Len(t, ok.GetAllTasks(tth.ctx), 1)
}

type TTHelper struct {
	//for both Task and AtxTask
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
	//only for AtxTask
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

func (t *TTHelper) GetAtxTask(status types.TaskStatus) *types.AtxTask {
	return t.app.OracleKeeper.BuildAtxTaskWithExpire(
		t.ctx, t.AtxTaskID(), t.creator.String(),
		t.bounty, t.validTime, status)
}

func (t *TTHelper) AtxTaskID() []byte {
	atxHash := sha256.Sum256([]byte(t.contract))
	return atxHash[:]
}

func (t *TTHelper) CheckAtxTask(i types.TaskI, err error) types.TaskI {
	require.Nil(t.t, err)
	res := i.(*types.AtxTask)
	require.Equal(t.t, t.creator, sdk.MustAccAddressFromBech32(res.Creator))
	require.Equal(t.t, t.bounty, res.Bounty)
	require.Equal(t.t, t.AtxTaskID(), res.GetID())
	require.Equal(t.t, t.validTime, res.ValidTime)
	t.CheckClosingTaskIDsShortcutTasks(i)
	return i
}

func (t *TTHelper) CheckClosingTaskIDsShortcutTasks(task types.TaskI) {
	ClosingTaskIDs := t.app.OracleKeeper.GetClosingTaskIDs(t.ctx, task)
	ShortcutTasksIDs := t.app.OracleKeeper.GetShortcutTasks(t.ctx)

	var duplicates = false
	for _, shortcutTasksID := range ShortcutTasksIDs {
		for _, closingTaskID := range ClosingTaskIDs {
			if bytes.Equal(shortcutTasksID.Tid, closingTaskID.Tid) {
				duplicates = true
				break
			}
		}
		if duplicates {
			break
		}
	}
	require.False(t.t, duplicates)
}
