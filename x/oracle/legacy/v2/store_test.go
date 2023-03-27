package v2_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

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
