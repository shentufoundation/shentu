package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func MigrateTaskStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	expireHeightTaskIDs := make(map[int64][]types.TaskID)
	iterator := sdk.KVStorePrefixIterator(store, types.TaskStoreKeyPrefix)

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
		store.Delete(iterator.Key())
		// set task
		bz, err := cdc.MarshalInterface(&newTask)
		if err != nil {
			panic(err)
		}
		store.Set(types.TaskStoreKey(newTask.GetID()), bz)
		// get all ExpireHeight and taskIDs
		if newTask.Status == types.TaskStatusPending {
			expireHeight := newTask.ExpireHeight
			newTaskID := types.TaskID{Tid: newTask.GetID()}
			expireHeightTaskIDs[expireHeight] = append(expireHeightTaskIDs[expireHeight], newTaskID)
		}
	}

	//  Migrate ExpireHeight and taskIDs
	for height, taskIDs := range expireHeightTaskIDs {
		bz := cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})
		store.Set(types.ClosingTaskIDsStoreKey(height), bz)
	}

	return nil
}
