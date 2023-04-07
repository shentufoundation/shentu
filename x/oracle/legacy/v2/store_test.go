package v2_test

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	params "github.com/cosmos/cosmos-sdk/x/params/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	v2 "github.com/shentufoundation/shentu/v2/x/oracle/legacy/v2"
	oracletypes "github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func Test_MigrateTaskStore(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cdc := shentuapp.MakeEncodingConfig().Marshaler

	// mock old data
	var tasks []v2.Task
	for i := 0; i < 10; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		beginBlock := r.Int63n(100)
		waitingBlocks := r.Int63n(10) + 1
		ClosingBlock := beginBlock + waitingBlocks
		status := r.Intn(4)

		task := v2.Task{
			Contract:      simtypes.RandStringOfLength(r, 10),
			Function:      simtypes.RandStringOfLength(r, 5),
			BeginBlock:    beginBlock,
			Bounty:        nil,
			Description:   simtypes.RandStringOfLength(r, 5),
			Expiration:    time.Time{},
			Creator:       "",
			Responses:     nil,
			Result:        sdk.Int{},
			ClosingBlock:  ClosingBlock,
			WaitingBlocks: waitingBlocks,
			Status:        v2.TaskStatus(status),
		}
		tasks = append(tasks, task)
	}

	store := ctx.KVStore(app.GetKey(oracletypes.StoreKey))
	for _, task := range tasks {
		// SetTask
		store.Set(TaskStoreKey(task.Contract, task.Function), cdc.MustMarshalLengthPrefixed(&task))
		// SetClosingBlockStore
		newTaskID := v2.TaskID{Contract: task.Contract, Function: task.Function}
		closingTaskIDsData := store.Get(oracletypes.ClosingTaskIDsStoreKey(task.ClosingBlock))
		var taskIDsProto v2.TaskIDs
		if closingTaskIDsData != nil {
			cdc.MustUnmarshalLengthPrefixed(closingTaskIDsData, &taskIDsProto)
		}
		taskIds := append(taskIDsProto.TaskIds, newTaskID)
		bz := cdc.MustMarshalLengthPrefixed(&v2.TaskIDs{TaskIds: taskIds})
		store.Set(oracletypes.ClosingTaskIDsStoreKey(task.ClosingBlock), bz)
	}

	err := v2.MigrateTaskStore(ctx, app.GetKey(oracletypes.StoreKey), cdc)
	require.Nil(t, err)
}

func TaskStoreKey(contract, function string) []byte {
	return append(append(oracletypes.TaskStoreKeyPrefix, []byte(contract)...), []byte(function)...)
}

func Test_MigrateTaskParams(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	ok := app.OracleKeeper
	oracleSubspace := app.GetSubspace(oracletypes.ModuleName)
	oldTaskParams := v2.TaskParams{
		ExpirationDuration: time.Duration(569000),
		AggregationWindow:  123,
		AggregationResult:  sdk.NewInt(58),
		ThresholdScore:     sdk.NewInt(75),
		Epsilon1:           sdk.NewInt(9),
		Epsilon2:           sdk.NewInt(7),
	}
	oldTable := params.NewKeyTable(
		params.NewParamSetPair(oracletypes.ParamsStoreKeyTaskParams, v2.TaskParams{}, func(i interface{}) error { return nil }),
		params.NewParamSetPair(oracletypes.ParamsStoreKeyPoolParams, v2.LockedPoolParams{}, func(i interface{}) error { return nil }),
	)

	tableField := reflect.ValueOf(&oracleSubspace).Elem().FieldByName("table")
	// save the KeyTable for restoring later
	cachedTable := GetUnexportedField(tableField)
	// set the KeyTable as old version of Oracle module
	SetUnexportedField(tableField, oldTable)
	// set the params in the form of old version of Oracle module
	oracleSubspace.Set(ctx, oracletypes.ParamsStoreKeyTaskParams, &oldTaskParams)
	tp := ok.GetTaskParams(ctx)
	require.True(t, tp.ShortcutQuorum.IsNil())
	// restore the KeyTable as this version of Oracle module
	SetUnexportedField(tableField, cachedTable)
	v2.UpdateParams(ctx, oracleSubspace)
	tp = ok.GetTaskParams(ctx)
	require.Equal(t, oracletypes.DefaultShortcutQuorum, tp.ShortcutQuorum)
	require.Equal(t, time.Duration(1800)*time.Second, tp.ExpirationDuration)
	require.Equal(t, oldTaskParams.AggregationWindow, tp.AggregationWindow)
	require.Equal(t, oldTaskParams.AggregationResult, tp.AggregationResult)
	require.Equal(t, oldTaskParams.ThresholdScore, tp.ThresholdScore)
	require.Equal(t, oldTaskParams.Epsilon1, tp.Epsilon1)
	require.Equal(t, oldTaskParams.Epsilon2, tp.Epsilon2)
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().Set(reflect.ValueOf(value))
}
