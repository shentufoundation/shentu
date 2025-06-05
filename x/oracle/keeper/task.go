package keeper

import (
	"bytes"
	"context"
	"errors"
	"time"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

var amplifier = math.NewInt(1000000)

// SetTask sets a task in KVStore.
func (k Keeper) SetTask(ctx context.Context, task types.TaskI) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.MarshalInterface(task)
	if err != nil {
		return err
	}
	return store.Set(types.TaskStoreKey(task.GetID()), bz)
}

// DeleteTask deletes a task from KVStore.
func (k Keeper) DeleteTask(ctx context.Context, task types.TaskI) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.TaskStoreKey(task.GetID()))
}

// UpdateAndSetTask updates a task and set it in KVStore.
func (k Keeper) UpdateAndSetTask(ctx context.Context, task *types.Task) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	task.ExpireHeight = sdkCtx.BlockHeight() + task.WaitingBlocks
	if task.GetStatus() == types.TaskStatusPending {
		if err := k.AddToClosingTaskIDs(ctx, task); err != nil {
			return err
		}
	}
	return k.SetTask(ctx, task)
}

func (k Keeper) SetTxTask(ctx context.Context, task *types.TxTask) error {
	if !task.IsExpired(ctx) {
		err := k.SaveExpireTxTask(ctx, task)
		if err != nil {
			return err
		}
		err = k.SetTask(ctx, task)
		if err != nil {
			return err
		}
		if task.GetStatus() == types.TaskStatusPending {
			if err = k.AddToClosingTaskIDs(ctx, task); err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) SaveExpireTxTask(ctx context.Context, task *types.TxTask) error {
	ids, err := k.GetTaskIDsByTime(ctx, types.ExpireTaskStoreKeyPrefix, task.Expiration)
	if err != nil {
		return err
	}
	ids = append(ids, types.TaskID{Tid: task.GetID()})
	bz := k.cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: ids})
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.TimeStoreKey(types.ExpireTaskStoreKeyPrefix, task.Expiration), bz)
}

func (k Keeper) DeleteFromExpireIDs(ctx context.Context, task types.TxTask) error {
	store := k.storeService.OpenKVStore(ctx)
	taskIDs, err := k.GetTaskIDsByTime(ctx, types.ExpireTaskStoreKeyPrefix, task.Expiration)
	if err != nil {
		return err
	}
	for i, taskID := range taskIDs {
		if bytes.Equal(taskID.Tid, task.GetID()) {
			taskIDs = append(taskIDs[:i], taskIDs[i+1:]...)
			break
		}
	}
	if len(taskIDs) == 0 {
		err := store.Delete(types.TimeStoreKey(types.ExpireTaskStoreKeyPrefix, task.Expiration))
		if err != nil {
			return err
		}
	} else {
		bz := k.cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})
		err := store.Set(types.TimeStoreKey(types.ExpireTaskStoreKeyPrefix, task.Expiration), bz)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetTask returns a task given taskID.
func (k Keeper) GetTask(ctx context.Context, taskID []byte) (task types.TaskI, err error) {
	store := k.storeService.OpenKVStore(ctx)
	TaskData, err := store.Get(types.TaskStoreKey(taskID))
	if err != nil {
		return nil, err
	}
	if len(TaskData) == 0 {
		return nil, errors.New("oracle: task not found")
	}
	err = k.cdc.UnmarshalInterface(TaskData, &task)
	return
}

// DeleteFromClosingTaskIDs remove ID of the task from closingBlockStore because it has been handled in shortcut
func (k Keeper) DeleteFromClosingTaskIDs(ctx context.Context, task types.TaskI) error {
	taskIDs, err := k.GetClosingTaskIDs(ctx, task)
	if err != nil {
		return err
	}
	for i := range taskIDs {
		if bytes.Equal(taskIDs[i].Tid, task.GetID()) {
			taskIDs = append(taskIDs[:i], taskIDs[i+1:]...)
			break
		}
	}
	if err := k.SetClosingTaskIDs(ctx, task, taskIDs); err != nil {
		return err
	}

	return nil
}

// AddToClosingTaskIDs sets the store of task IDs for aggregation on time.
func (k Keeper) AddToClosingTaskIDs(ctx context.Context, task types.TaskI) error {
	newTaskID := types.TaskID{Tid: task.GetID()}
	taskIDs, err := k.GetClosingTaskIDs(ctx, task)
	if err != nil {
		return err
	}
	taskIDs = append(taskIDs, newTaskID)
	if err = k.SetClosingTaskIDs(ctx, task, taskIDs); err != nil {
		return err
	}

	return nil
}

func (k Keeper) SetClosingTaskIDs(ctx context.Context, task types.TaskI, taskIDs []types.TaskID) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})
	switch task := task.(type) {
	case *types.Task:
		err := store.Set(types.ClosingTaskIDsStoreKey(task.ExpireHeight), bz)
		if err != nil {
			return err
		}
	case *types.TxTask:
		err := store.Set(types.TimeStoreKey(types.ClosingTaskStoreKeyTimedPrefix, task.ValidTime), bz)
		if err != nil {
			return err
		}
	default:
		return errors.New("oracle: unknown implementation of TaskI")
	}
	return nil
}

// GetClosingTaskIDs returns a list of task IDs by the closing block and valid time.
func (k Keeper) GetClosingTaskIDs(ctx context.Context, task types.TaskI) (resIDs []types.TaskID, err error) {
	height, theTime := task.GetValidTime()
	if height > 0 {
		taskIDs, err := k.GetClosingTaskIDsByHeight(ctx, height)
		if err != nil {
			return nil, err
		}
		resIDs = append(resIDs, taskIDs...)
	}
	if !theTime.IsZero() {
		taskIDs, err := k.GetTaskIDsByTime(ctx, types.ClosingTaskStoreKeyTimedPrefix, theTime)
		if err != nil {
			return nil, err
		}
		resIDs = append(resIDs, taskIDs...)
	}
	return
}

func (k Keeper) GetClosingTaskIDsByHeight(ctx context.Context, blockHeight int64) ([]types.TaskID, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ClosingTaskIDsStoreKey(blockHeight))
	if err != nil {
		return nil, err
	}
	var taskIDsProto types.TaskIDs
	if bz != nil {
		k.cdc.MustUnmarshalLengthPrefixed(bz, &taskIDsProto)
	}
	return taskIDsProto.TaskIds, nil
}

func (k Keeper) IteratorTaskIDsByEndTime(ctx context.Context, prefix []byte, endTime time.Time, callback func(key, value []byte) (stop bool)) {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(prefix, storetypes.InclusiveEndBytes(types.TimeStoreKey(prefix, endTime)))
	if err != nil {
		return
	}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if callback(iterator.Key(), iterator.Value()) {
			break
		}
	}
}

func (k Keeper) GetTaskIDsByTime(ctx context.Context, prefix []byte, theTime time.Time) ([]types.TaskID, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.TimeStoreKey(prefix, theTime))
	if err != nil {
		return nil, err
	}
	var taskIDsProto types.TaskIDs
	if bz != nil {
		k.cdc.MustUnmarshalLengthPrefixed(bz, &taskIDsProto)
	}
	return taskIDsProto.TaskIds, nil
}

// DeleteClosingTaskIDs deletes stores for task IDs closed at given block.
func (k Keeper) DeleteClosingTaskIDs(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(ctx)
	err := store.Delete(types.ClosingTaskIDsStoreKey(sdkCtx.BlockHeight()))
	if err != nil {
		return err
	}
	k.IteratorTaskIDsByEndTime(
		ctx, types.ClosingTaskStoreKeyTimedPrefix, sdkCtx.BlockTime(),
		func(key, _ []byte) bool {
			err := store.Delete(key)
			if err != nil {
				return false
			}
			return false
		})
	return nil
}

// delete tasks whose expiration >= BlockTime
// the taget task may already be gone due to explicitally removed by user
func (k Keeper) DeleteExpiredTasks(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(ctx)
	k.IteratorTaskIDsByEndTime(
		ctx, types.ExpireTaskStoreKeyPrefix, sdkCtx.BlockTime(),
		func(key, value []byte) bool {
			var taskIDsProto types.TaskIDs
			k.cdc.MustUnmarshalLengthPrefixed(value, &taskIDsProto)
			protoTids := taskIDsProto.TaskIds
			for i := range protoTids {
				storeKey := types.TaskStoreKey(protoTids[i].Tid)
				if _, err := store.Has(storeKey); err == nil {
					err := store.Delete(storeKey)
					if err != nil {
						return false
					}
				}
			}
			err := store.Delete(key)
			if err != nil {
				return false
			}
			return false
		})
	return nil
}

func (k Keeper) DeleteShortcutTasks(ctx context.Context) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.ShortcutTasksKeyPrefix)
}

func (k Keeper) GetInvalidTaskIDs(ctx context.Context) (resIDs []types.TaskID) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	ids, err := k.GetClosingTaskIDsByHeight(ctx, sdkCtx.BlockHeight())
	if err != nil {
		return nil
	}
	resIDs = append(resIDs, ids...)
	k.IteratorTaskIDsByEndTime(
		ctx, types.ClosingTaskStoreKeyTimedPrefix, sdkCtx.BlockTime(),
		func(_, value []byte) bool {
			var taskIDsProto types.TaskIDs
			k.cdc.MustUnmarshalLengthPrefixed(value, &taskIDsProto)
			resIDs = append(resIDs, taskIDsProto.TaskIds...)
			return false
		})
	return resIDs
}

// calling of CreateTask creates one of following
// 1. Task (smart contract task)
// 2. TxTask (transaction task)
// 3. placeholder TxTask (status:TaskStatusNil, creator: nil, bounty:nil)
func (k Keeper) CreateTask(ctx context.Context, creator sdk.AccAddress, task types.TaskI) error {
	savedTask, err := k.GetTask(ctx, task.GetID())
	if err == nil {
		if savedTask.GetStatus() == types.TaskStatusPending {
			return types.ErrTaskNotClosed
		}
		if err = k.DeleteTask(ctx, savedTask); err != nil {
			return err
		}
		if txTask, ok := savedTask.(*types.TxTask); ok && savedTask.GetStatus() != types.TaskStatusNil {
			if err = k.DeleteFromExpireIDs(ctx, *txTask); err != nil {
				return err
			}
		}
	}

	if err = k.SetTask(ctx, task); err != nil {
		return err
	}
	// if task's status is TaskStatusNil, it will not go to ClosingBlockStore,
	// therefor will not be handled in EndBlocker
	if task.GetStatus() == types.TaskStatusPending {
		if err = k.AddToClosingTaskIDs(ctx, task); err != nil {
			return err
		}
		if err := k.CollectBounty(ctx, task.GetBounty(), creator); err != nil {
			return err
		}
		if _, ok := task.(*types.TxTask); ok {
			k.TryShortcut(ctx, task)
		}
	}
	return nil
}

func (k Keeper) BuildTxTaskWithExpire(ctx context.Context, txHash []byte, creator string, bounty sdk.Coins, validTime time.Time, status types.TaskStatus) (*types.TxTask, error) {
	taskParams := k.GetTaskParams(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	txTask := types.NewTxTask(txHash, creator, bounty, validTime, status)
	txTask.Expiration = sdkCtx.BlockTime().Add(taskParams.ExpirationDuration)

	if err := k.SaveExpireTxTask(ctx, txTask); err != nil {
		return nil, err
	}
	return txTask, nil
}

func (k Keeper) BuildTxTask(ctx context.Context, txHash []byte, creator string, bounty sdk.Coins, validTime time.Time) (types.TaskI, error) {
	var txTask *types.TxTask
	var err error
	// if a TaskStatusNil task already exists, overwrite it after copying several fields.
	// please be noted that the expiration hook remains
	if savedTask, err := k.GetTask(ctx, txHash); err == nil {
		if savedTask.GetStatus() == types.TaskStatusNil {
			// in fast-path case, a TxTask could be created before the creatTxTask msg
			savedTask, ok := savedTask.(*types.TxTask)
			if !ok {
				return nil, types.ErrUnexpectedTask
			}
			txTask = types.NewTxTask(txHash, creator, bounty, validTime, types.TaskStatusPending)
			txTask.Expiration = savedTask.Expiration
			txTask.Responses = savedTask.Responses
			txTask.Score = savedTask.Score
		}
	}
	if txTask == nil {
		// BuildTxTaskWithExpire should be called with new TxTask created and expiration hooking up
		txTask, err = k.BuildTxTaskWithExpire(ctx, txHash, creator, bounty, validTime, types.TaskStatusPending)
		if err != nil {
			return txTask, err
		}
	}

	if validTime.After(txTask.Expiration) {
		return nil, types.ErrTooLateValidTime
	}
	return txTask, nil
}

func (k Keeper) GetShortcutTasks(ctx context.Context) ([]types.TaskID, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ShortcutTasksKeyPrefix)
	if err != nil {
		return nil, err
	}
	var taskIDsProto types.TaskIDs
	k.cdc.MustUnmarshalLengthPrefixed(bz, &taskIDsProto)
	return taskIDsProto.TaskIds, nil
}

func (k Keeper) SetShortcutTasks(ctx context.Context, tid []byte) error {
	store := k.storeService.OpenKVStore(ctx)
	// here, the uniqueness of tid is not checked deliberately
	// later on when iterating the ShortcutTasks, the task
	// status will be checked to make sure every one is touched once.
	tasks, err := k.GetShortcutTasks(ctx)
	if err != nil {
		return err
	}
	tasks = append(tasks, types.TaskID{Tid: tid})
	bz := k.cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: tasks})
	return store.Set(types.ShortcutTasksKeyPrefix, bz)
}

func (k Keeper) TryShortcut(ctx context.Context, task types.TaskI) {
	if task.ShouldAgg(ctx) || task.GetStatus() != types.TaskStatusPending {
		// skip checking quorum, if
		// 1. this task will be handled in EndBlocker later on
		// 2. this task is not pending
		return
	}
	totalCollateral, err := k.GetTotalCollateral(ctx)
	if err != nil || totalCollateral.Empty() || totalCollateral[0].Amount.IsZero() {
		return
	}
	respondedCollateral := math.ZeroInt()
	for _, response := range task.GetResponses() {
		operatorAddr := sdk.MustAccAddressFromBech32(response.Operator)
		amount, err := k.GetCollateralAmount(ctx, operatorAddr)
		if err != nil {
			continue
		}
		respondedCollateral = respondedCollateral.Add(amount)
	}

	taskParams := k.GetTaskParams(ctx)
	if math.LegacyNewDecFromInt(respondedCollateral).
		Quo(math.LegacyNewDecFromInt(totalCollateral[0].Amount)).
		GTE(taskParams.ShortcutQuorum) {
		err := k.SetShortcutTasks(ctx, task.GetID())
		if err != nil {
			return
		}
		err = k.DeleteFromClosingTaskIDs(ctx, task)
		if err != nil {
			return
		}
	}
}

// RemoveTask removes a task from kvstore if it is closed, expired and requested by its creator.
// The id of the removed task may still remain in the ExpireTaskIDsStore.
//
//	in such case, when it's expired, the unfound task will be simply skipped
func (k Keeper) RemoveTask(ctx context.Context, taskID []byte, force bool, deleter sdk.AccAddress) error {
	task, err := k.GetTask(ctx, taskID)
	if err != nil {
		return err
	}
	if !force && !task.IsExpired(ctx) {
		return types.ErrNotExpired
	}

	if task.GetStatus() == types.TaskStatusPending ||
		task.GetStatus() == types.TaskStatusNil {
		return types.ErrNotFinished
	}

	// TODO: only creator can delete the task for now
	creatorAddr := sdk.MustAccAddressFromBech32(task.GetCreator())
	if !creatorAddr.Equals(deleter) {
		return types.ErrNotCreator
	}
	err = k.DeleteTask(ctx, task)
	if err != nil {
		return err
	}
	if txTask, ok := task.(*types.TxTask); ok {
		if err := k.DeleteFromExpireIDs(ctx, *txTask); err != nil {
			return err
		}
	}
	return nil
}

// IteratorAllTasks iterates over all the stored tasks and performs a callback function.
func (k Keeper) IteratorAllTasks(ctx context.Context, callback func(task types.TaskI) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.TaskStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var task types.TaskI
		if err := k.cdc.UnmarshalInterface(iterator.Value(), &task); err != nil {
			panic(err)
		}

		if callback(task) {
			break
		}
	}
}

// GetAllTasks gets all tasks.
func (k Keeper) GetAllTasks(ctx context.Context) (tasks []types.TaskI) {
	k.IteratorAllTasks(ctx, func(task types.TaskI) bool {
		tasks = append(tasks, task)
		return false
	})
	return
}

// UpdateAndGetAllTasks updates all tasks and returns them.
func (k Keeper) UpdateAndGetAllTasks(ctx context.Context) (tasks []types.Task, txTasks []types.TxTask) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.IteratorAllTasks(ctx, func(task types.TaskI) bool {
		if t, ok := task.(*types.Task); ok {
			t.WaitingBlocks = t.ExpireHeight - sdkCtx.BlockHeight()
			tasks = append(tasks, *t)
		} else if t, ok := task.(*types.TxTask); ok {
			txTasks = append(txTasks, *t)
		}
		return false
	})
	return
}

// IsValidResponse returns error if a response is not valid.
func (k Keeper) IsValidResponse(_ context.Context, task types.TaskI, response types.Response) error {
	// due to fast-path, response should be allowed to add if it's a TaskStatusNil task
	if task.GetStatus() != types.TaskStatusPending &&
		task.GetStatus() != types.TaskStatusNil {
		return types.ErrTaskClosed
	}
	for _, r := range task.GetResponses() {
		if r.Operator == response.Operator {
			return types.ErrDuplicateResponse
		}
	}
	if response.Score.LT(types.MinScore) || response.Score.GT(types.MaxScore) {
		return types.ErrInvalidScore
	}
	return nil
}

func (k Keeper) HandleNoneTxTaskForResponse(ctx context.Context, txHash []byte) error {
	if _, err := k.GetTask(ctx, txHash); err != nil {
		// if the corresponding TxTask doesn't exit,
		// create one as a placeholder (statue being set as TaskStatusNil),
		// waiting for the MsgCreateTxTask coming to fill in necessary fields
		txTask, err := k.BuildTxTaskWithExpire(ctx, txHash, "", nil, time.Time{}, types.TaskStatusNil)
		if err != nil {
			return err
		}
		return k.CreateTask(ctx, nil, txTask)
	}
	return nil
}

// RespondToTask records the response from an operator for a task.
func (k Keeper) RespondToTask(ctx context.Context, taskID []byte, score int64, operatorAddress sdk.AccAddress) error {
	if _, err := k.IsOperator(ctx, operatorAddress); err != nil {
		return err
	}

	task, err := k.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	response := types.NewResponse(math.NewInt(score), operatorAddress)
	err = k.IsValidResponse(ctx, task, response)
	if err != nil {
		return err
	}

	task.AddResponse(response)
	if err = k.SetTask(ctx, task); err != nil {
		return err
	}

	if _, ok := task.(*types.TxTask); ok {
		k.TryShortcut(ctx, task)
	}

	return nil
}

// Aggregate does an aggregation of responses for a task and updated task result.
func (k Keeper) Aggregate(ctx context.Context, taskID []byte) error {
	taskParams := k.GetTaskParams(ctx)
	task, err := k.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	if task.GetStatus() != types.TaskStatusPending {
		return types.ErrTaskClosed
	}

	result := taskParams.AggregationResult
	totalCollateral := math.NewInt(0)
	minScoreCollateral := math.NewInt(0)
	responses := task.GetResponses()
	for i, response := range responses {
		operatorAddr := sdk.MustAccAddressFromBech32(response.Operator)
		amount, err := k.GetCollateralAmount(ctx, operatorAddr)
		if err != nil {
			continue
		}
		result = result.Add(response.Score.Mul(amount))
		responses[i].Weight = amount
		totalCollateral = totalCollateral.Add(amount)
		if response.Score.Equal(types.MinScore) {
			minScoreCollateral = minScoreCollateral.Add(amount)
		}
	}

	if totalCollateral.IsPositive() {
		if minScoreCollateral.MulRaw(3).GTE(totalCollateral) {
			result = types.MinScore
			for i, response := range responses {
				if !response.Score.Equal(types.MinScore) {
					responses[i].Weight = math.NewInt(0)
				}
			}
		} else {
			result = result.Quo(totalCollateral)
		}
		task.SetStatus(types.TaskStatusSucceeded)
	} else {
		task.SetStatus(types.TaskStatusFailed)
	}
	task.SetScore(result.Int64())
	return k.SetTask(ctx, task)
}

// TotalValidTaskCollateral calculates the total amount of valid collateral of a task.
func (k Keeper) TotalValidTaskCollateral(ctx context.Context, task types.TaskI) math.Int {
	taskParams := k.GetTaskParams(ctx)
	totalValidTaskCollateral := math.NewInt(0)
	responses := task.GetResponses()
	if task.GetScore() == types.MinScore.Int64() {
		for _, response := range responses {
			if response.Score.Equal(types.MinScore) {
				operatorAddr := sdk.MustAccAddressFromBech32(response.Operator)
				collateral, err := k.GetCollateralAmount(ctx, operatorAddr)
				if err != nil {
					continue
				}
				totalValidTaskCollateral = totalValidTaskCollateral.Add(collateral)
			}
		}
	} else if task.GetScore() < taskParams.ThresholdScore.Int64() {
		for _, response := range responses {
			if response.Score.LT(taskParams.ThresholdScore) {
				operatorAddr := sdk.MustAccAddressFromBech32(response.Operator)
				collateral, err := k.GetCollateralAmount(ctx, operatorAddr)
				if err != nil {
					continue
				}
				totalValidTaskCollateral = totalValidTaskCollateral.Add(
					amplifier.Mul(collateral).Quo(response.Score.Add(taskParams.Epsilon1)),
				)
			}
		}
	} else {
		for _, response := range responses {
			if response.Score.GTE(taskParams.ThresholdScore) {
				operatorAddr := sdk.MustAccAddressFromBech32(response.Operator)
				collateral, err := k.GetCollateralAmount(ctx, operatorAddr)
				if err != nil {
					continue
				}
				totalValidTaskCollateral = totalValidTaskCollateral.Add(
					amplifier.Mul(collateral).Quo(types.MaxScore.Sub(response.Score).Add(taskParams.Epsilon2)),
				)
			}
		}
	}
	return totalValidTaskCollateral
}

// TODO: this is a simplified version (without confidence calculation)

// DistributeBounty distributes bounty to operators based on responses and the aggregation result.
func (k Keeper) DistributeBounty(ctx context.Context, task types.TaskI) error {
	taskParams := k.GetTaskParams(ctx)
	totalValidTaskCollateral := k.TotalValidTaskCollateral(ctx, task)
	if totalValidTaskCollateral.IsZero() {
		return types.ErrTaskFailed
	}

	responses := task.GetResponses()
	for _, bounty := range task.GetBounty() {
		if task.GetScore() == types.MinScore.Int64() {
			for i := range responses {
				if responses[i].Score.Equal(types.MinScore) {
					operatorAddr := sdk.MustAccAddressFromBech32(responses[i].Operator)
					collateral, err := k.GetCollateralAmount(ctx, operatorAddr)
					if err != nil {
						continue
					}
					amount := bounty.Amount.Mul(collateral).Quo(totalValidTaskCollateral)
					reward := sdk.NewCoins(sdk.NewCoin(bounty.Denom, amount))
					if err := k.AddReward(ctx, operatorAddr, reward); err == nil {
						responses[i].Reward = reward
					}
				}
			}
		} else if task.GetScore() < taskParams.ThresholdScore.Int64() {
			for i := range responses {
				if responses[i].Score.LT(taskParams.ThresholdScore) {
					operatorAddr, err := sdk.AccAddressFromBech32(responses[i].Operator)
					if err != nil {
						panic(err)
					}
					collateral, err := k.GetCollateralAmount(ctx, operatorAddr)
					if err != nil {
						continue
					}
					amount := bounty.Amount.Mul(
						amplifier.Mul(collateral).Quo(responses[i].Score.Add(taskParams.Epsilon1)),
					).Quo(totalValidTaskCollateral)
					reward := sdk.NewCoins(sdk.NewCoin(bounty.Denom, amount))
					if err := k.AddReward(ctx, operatorAddr, reward); err == nil {
						responses[i].Reward = reward
					}
				}
			}
		} else {
			for i := range responses {
				if responses[i].Score.GTE(taskParams.ThresholdScore) {
					operatorAddr := sdk.MustAccAddressFromBech32(responses[i].Operator)
					collateral, err := k.GetCollateralAmount(ctx, operatorAddr)
					if err != nil {
						continue
					}
					amount := bounty.Amount.Mul(
						amplifier.Mul(collateral).Quo(types.MaxScore.Sub(responses[i].Score).Add(taskParams.Epsilon2)),
					).Quo(totalValidTaskCollateral)
					reward := sdk.NewCoins(sdk.NewCoin(bounty.Denom, amount))
					if err := k.AddReward(ctx, operatorAddr, reward); err == nil {
						responses[i].Reward = reward
					}
				}
			}
		}
	}
	return k.SetTask(ctx, task)
}

func (k Keeper) RefundBounty(ctx context.Context, task types.TaskI) error {
	taskCreator, err := sdk.AccAddressFromBech32(task.GetCreator())
	if err != nil {
		panic(err)
	}

	totalReward := make(sdk.Coins, 0, 1)
	for _, response := range task.GetResponses() {
		if response.Reward != nil {
			totalReward = totalReward.Add(response.Reward...)
		}
	}

	bounties := task.GetBounty()
	leftBounty := bounties.Sub(totalReward...)
	if leftBounty != nil && leftBounty.IsAllPositive() {
		oracleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
		spendableCoins := k.bankKeeper.SpendableCoins(ctx, oracleAddress)
		if ok := spendableCoins.IsAllGTE(leftBounty); !ok {
			panic("Insufficient oracle model balance")
		}

		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, taskCreator, leftBounty); err != nil {
			if err = k.distrKeeper.FundCommunityPool(ctx, leftBounty, oracleAddress); err != nil {
				panic(err)
			}
			return err
		}
	}
	return nil
}
