package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func MigrateTaskStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	expireHeightMap := make(map[int64]bool)
	iterator := sdk.KVStorePrefixIterator(store, types.TaskStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldTask Task
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &oldTask)

		newTask := types.Task{
			Contract:      oldTask.Contract,
			Function:      oldTask.Function,
			BeginBlock:    oldTask.BeginBlock,
			Bounty:        oldTask.Bounty,
			Description:   oldTask.Description,
			Expiration:    oldTask.Expiration,
			Creator:       oldTask.Creator,
			Responses:     nil,
			Result:        oldTask.Result,
			ExpireHeight:  oldTask.ClosingBlock,
			WaitingBlocks: oldTask.WaitingBlocks,
			Status:        types.TaskStatus(oldTask.Status),
		}
		for _, response := range oldTask.Responses {
			newResponse := types.Response{
				Operator: response.Operator,
				Score:    response.Score,
				Weight:   response.Weight,
				Reward:   response.Reward,
			}
			newTask.Responses = append(newTask.Responses, newResponse)
		}

		// delete old task
		store.Delete(iterator.Key())
		// set task
		bz, err := cdc.MarshalInterface(&newTask)
		if err != nil {
			return err
		}
		store.Set(types.TaskStoreKey(newTask.GetID()), bz)
		// get all ExpireHeight
		expireHeightMap[newTask.ExpireHeight] = true
	}

	//  Migrate ClosingTaskIDs
	for height := range expireHeightMap {
		closingTaskIDsData := store.Get(types.ClosingTaskIDsStoreKey(height))
		if closingTaskIDsData == nil {
			continue
		}

		var oldTaskIDs TaskIDs
		var newTaskIDs []types.TaskID
		cdc.MustUnmarshalLengthPrefixed(closingTaskIDsData, &oldTaskIDs)
		for _, t := range oldTaskIDs.TaskIds {
			newTaskID := types.TaskID{Tid: append([]byte(t.Contract), []byte(t.Function)...)}
			newTaskIDs = append(newTaskIDs, newTaskID)
		}

		bz := cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: newTaskIDs})
		store.Set(types.ClosingTaskIDsStoreKey(height), bz)
	}

	return nil
}

func UpdateParams(ctx sdk.Context, paramSubspace types.ParamSubspace) {
	var taskParams types.TaskParams
	paramSubspace.Get(ctx, types.ParamsStoreKeyTaskParams, &taskParams)
	newTaskParams := types.NewTaskParams(
		taskParams.ExpirationDuration,
		taskParams.AggregationWindow,
		taskParams.AggregationResult,
		taskParams.ThresholdScore,
		taskParams.Epsilon1,
		taskParams.Epsilon2,
		types.DefaultShortcutQuorum,
	)
	paramSubspace.Set(ctx, types.ParamsStoreKeyTaskParams, &newTaskParams)
}
