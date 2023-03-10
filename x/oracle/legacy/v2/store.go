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
		var oldTask Task
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &oldTask)

		newTask := types.NewTask(
			oldTask.Contract,
			oldTask.Function,
			oldTask.BeginBlock,
			oldTask.Bounty,
			oldTask.Description,
			oldTask.Expiration,
			sdk.AccAddress(oldTask.Creator),
			oldTask.ClosingBlock,
			oldTask.WaitingBlocks,
		)
		newTask.Status = types.TaskStatus(oldTask.Status)

		// delete old task
		oldTaskKey := append(append(types.TaskStoreKeyPrefix, []byte(oldTask.Contract)...), []byte(oldTask.Function)...)
		store.Delete(oldTaskKey)
		// set task
		bz, err := cdc.MarshalInterface(&newTask)
		if err != nil {
			panic(err)
		}
		store.Set(types.TaskStoreKey(newTask.GetID()), bz)
		// task IDs
		newTaskID := types.TaskID{Tid: newTask.GetID()}
		taskIDs = append(taskIDs, newTaskID)
		bz = cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})
		store.Set(types.ClosingTaskIDsStoreKey(newTask.ExpireHeight), bz)
	}

	return nil
}
