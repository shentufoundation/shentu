package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func MigrateTaskStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TaskStoreKeyPrefix)

	var taskIDs []types.TaskID

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var task types.Task
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &task)

		newTask := types.NewTask(
			task.Contract,
			task.Function,
			task.BeginBlock,
			task.Bounty,
			task.Description,
			task.Expiration,
			sdk.AccAddress(task.Creator),
			task.ClosingBlock,
			task.WaitingBlocks,
		)

		// task
		store.Set(types.TaskStoreKey(newTask.GetID()), cdc.MustMarshalLengthPrefixed(&newTask))
		// task IDs
		newTaskID := types.TaskID{Tid: task.GetID()}
		taskIDs = append(taskIDs, newTaskID)
		bz := cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})
		store.Set(types.ClosingTaskIDsStoreKey(task.ClosingBlock), bz)

		// delete old task
		oldTaskKey := append(append(types.TaskStoreKeyPrefix, []byte(task.Contract)...), []byte(task.Function)...)
		store.Delete(oldTaskKey)

	}

	return nil
}
