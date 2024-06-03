package keeper

import (
	"bytes"
	"errors"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

var (
	amplifier = sdk.NewInt(1000000)
)

// SetTask sets a task in KVStore.
func (k Keeper) SetTask(ctx sdk.Context, task types.TaskI) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalInterface(task)
	if err != nil {
		panic(err)
	}
	store.Set(types.TaskStoreKey(task.GetID()), bz)
}

// DeleteTask deletes a task from KVStore.
func (k Keeper) DeleteTask(ctx sdk.Context, task types.TaskI) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.TaskStoreKey(task.GetID()))
	return nil
}

// UpdateAndSetTask updates a task and set it in KVStore.
func (k Keeper) UpdateAndSetTask(ctx sdk.Context, task *types.Task) {
	task.ExpireHeight = ctx.BlockHeight() + task.WaitingBlocks
	if task.GetStatus() == types.TaskStatusPending {
		k.AddToClosingTaskIDs(ctx, task)
	}
	k.SetTask(ctx, task)
}

func (k Keeper) SetTxTask(ctx sdk.Context, task *types.TxTask) {
	if !task.IsExpired(ctx) {
		k.SaveExpireTxTask(ctx, task)
		k.SetTask(ctx, task)
		if task.GetStatus() == types.TaskStatusPending {
			k.AddToClosingTaskIDs(ctx, task)
		}
	}
}

func (k Keeper) SaveExpireTxTask(ctx sdk.Context, task *types.TxTask) {
	ids := k.GetTaskIDsByTime(ctx, types.ExpireTaskStoreKeyPrefix, task.Expiration)
	ids = append(ids, types.TaskID{Tid: task.GetID()})
	bz := k.cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: ids})
	ctx.KVStore(k.storeKey).Set(
		types.TimeStoreKey(types.ExpireTaskStoreKeyPrefix, task.Expiration),
		bz,
	)
}

func (k Keeper) DeleteFromExpireIDs(ctx sdk.Context, task types.TxTask) {
	store := ctx.KVStore(k.storeKey)
	taskIDs := k.GetTaskIDsByTime(ctx, types.ExpireTaskStoreKeyPrefix, task.Expiration)
	if taskIDs == nil {
		return
	}
	for i, taskID := range taskIDs {
		if bytes.Equal(taskID.Tid, task.GetID()) {
			taskIDs = append(taskIDs[:i], taskIDs[i+1:]...)
			break
		}
	}
	if len(taskIDs) == 0 {
		store.Delete(types.TimeStoreKey(types.ExpireTaskStoreKeyPrefix, task.Expiration))
	} else {
		bz := k.cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})
		store.Set(
			types.TimeStoreKey(types.ExpireTaskStoreKeyPrefix, task.Expiration),
			bz,
		)
	}
}

// GetTask returns a task given taskID.
func (k Keeper) GetTask(ctx sdk.Context, taskID []byte) (task types.TaskI, err error) {
	TaskData := ctx.KVStore(k.storeKey).Get(types.TaskStoreKey(taskID))
	if TaskData == nil {
		return nil, types.ErrTaskNotExists
	}
	err = k.cdc.UnmarshalInterface(TaskData, &task)
	return
}

// DeleteFromClosingTaskIDs remove ID of the task from closingBlockStore because it has been handled in shortcut
func (k Keeper) DeleteFromClosingTaskIDs(ctx sdk.Context, task types.TaskI) {
	taskIDs := k.GetClosingTaskIDs(ctx, task)
	for i := range taskIDs {
		if bytes.Equal(taskIDs[i].Tid, task.GetID()) {
			taskIDs = append(taskIDs[:i], taskIDs[i+1:]...)
			break
		}
	}
	k.SetClosingTaskIDs(ctx, task, taskIDs)
}

// AddToClosingTaskIDs sets the store of task IDs for aggregation on time.
func (k Keeper) AddToClosingTaskIDs(ctx sdk.Context, task types.TaskI) {
	newTaskID := types.TaskID{Tid: task.GetID()}
	taskIDs := append(k.GetClosingTaskIDs(ctx, task), newTaskID)
	k.SetClosingTaskIDs(ctx, task, taskIDs)
}

func (k Keeper) SetClosingTaskIDs(ctx sdk.Context, task types.TaskI, taskIDs []types.TaskID) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})
	switch task := task.(type) {
	case *types.Task:
		store.Set(types.ClosingTaskIDsStoreKey(task.ExpireHeight), bz)
		return
	case *types.TxTask:
		store.Set(types.TimeStoreKey(types.ClosingTaskStoreKeyTimedPrefix, task.ValidTime), bz)
		return
	default:
		panic(errors.New("oracle: unknown implementation of TaskI"))
	}
}

// GetClosingTaskIDs returns a list of task IDs by the closing block and valid time.
func (k Keeper) GetClosingTaskIDs(ctx sdk.Context, task types.TaskI) (resIDs []types.TaskID) {
	height, theTime := task.GetValidTime()
	if height > 0 {
		resIDs = append(resIDs, k.GetClosingTaskIDsByHeight(ctx, height)...)
	}
	if !theTime.IsZero() {
		resIDs = append(resIDs, k.GetTaskIDsByTime(ctx, types.ClosingTaskStoreKeyTimedPrefix, theTime)...)
	}
	return
}

func (k Keeper) GetClosingTaskIDsByHeight(ctx sdk.Context, blockHeight int64) []types.TaskID {
	bz := ctx.KVStore(k.storeKey).Get(types.ClosingTaskIDsStoreKey(blockHeight))

	var taskIDsProto types.TaskIDs
	if bz != nil {
		k.cdc.MustUnmarshalLengthPrefixed(bz, &taskIDsProto)
	}
	return taskIDsProto.TaskIds
}

func (k Keeper) IteratorTaskIDsByEndTime(ctx sdk.Context, prefix []byte, endTime time.Time, callback func(key, value []byte) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(prefix,
		sdk.InclusiveEndBytes(types.TimeStoreKey(prefix, endTime)))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if callback(iterator.Key(), iterator.Value()) {
			break
		}
	}
}

func (k Keeper) GetTaskIDsByTime(ctx sdk.Context, prefix []byte, theTime time.Time) []types.TaskID {
	bz := ctx.KVStore(k.storeKey).Get(types.TimeStoreKey(prefix, theTime))
	var taskIDsProto types.TaskIDs
	if bz != nil {
		k.cdc.MustUnmarshalLengthPrefixed(bz, &taskIDsProto)
	}
	return taskIDsProto.TaskIds
}

// DeleteClosingTaskIDs deletes stores for task IDs closed at given block.
func (k Keeper) DeleteClosingTaskIDs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ClosingTaskIDsStoreKey(ctx.BlockHeight()))
	k.IteratorTaskIDsByEndTime(
		ctx, types.ClosingTaskStoreKeyTimedPrefix, ctx.BlockTime(),
		func(key, value []byte) bool {
			store.Delete(key)
			return false
		})
}

// delete tasks whose expiration >= BlockTime
// the taget task may already be gone due to explicitally removed by user
func (k Keeper) DeleteExpiredTasks(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	k.IteratorTaskIDsByEndTime(
		ctx, types.ExpireTaskStoreKeyPrefix, ctx.BlockTime(),
		func(key, value []byte) bool {
			var taskIDsProto types.TaskIDs
			k.cdc.MustUnmarshalLengthPrefixed(value, &taskIDsProto)
			protoTids := taskIDsProto.TaskIds
			for i := range protoTids {
				storeKey := types.TaskStoreKey(protoTids[i].Tid)
				if store.Has(storeKey) {
					store.Delete(storeKey)
				}
			}
			store.Delete(key)
			return false
		})
}

func (k Keeper) DeleteShortcutTasks(ctx sdk.Context) {
	ctx.KVStore(k.storeKey).Delete(types.ShortcutTasksKeyPrefix)
}

func (k Keeper) GetInvalidTaskIDs(ctx sdk.Context) (resIDs []types.TaskID) {
	resIDs = append(resIDs, k.GetClosingTaskIDsByHeight(ctx, ctx.BlockHeight())...)
	k.IteratorTaskIDsByEndTime(
		ctx, types.ClosingTaskStoreKeyTimedPrefix, ctx.BlockTime(),
		func(key, value []byte) bool {
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
func (k Keeper) CreateTask(ctx sdk.Context, creator sdk.AccAddress, task types.TaskI) error {
	savedTask, err := k.GetTask(ctx, task.GetID())
	if err == nil {
		if savedTask.GetStatus() == types.TaskStatusPending {
			return types.ErrTaskNotClosed
		}
		if err := k.DeleteTask(ctx, savedTask); err != nil {
			return err
		}
		if txTask, ok := savedTask.(*types.TxTask); ok && savedTask.GetStatus() != types.TaskStatusNil {
			k.DeleteFromExpireIDs(ctx, *txTask)
		}
	}

	k.SetTask(ctx, task)
	// if task's status is TaskStatusNil, it will not go to ClosingBlockStore,
	// therefor will not be handled in EndBlocker
	if task.GetStatus() == types.TaskStatusPending {
		k.AddToClosingTaskIDs(ctx, task)
		if err := k.CollectBounty(ctx, task.GetBounty(), creator); err != nil {
			return err
		}
		if _, ok := task.(*types.TxTask); ok {
			k.TryShortcut(ctx, task)
		}
	}
	return nil
}

func (k Keeper) BuildTxTaskWithExpire(ctx sdk.Context, txHash []byte, creator string, bounty sdk.Coins, validTime time.Time, status types.TaskStatus) *types.TxTask {
	taskParams := k.GetTaskParams(ctx)
	txTask := types.NewTxTask(txHash, creator, bounty, validTime, status)
	txTask.Expiration = ctx.BlockTime().Add(taskParams.ExpirationDuration)

	k.SaveExpireTxTask(ctx, txTask)
	return txTask
}

func (k Keeper) BuildTxTask(ctx sdk.Context, txHash []byte, creator string, bounty sdk.Coins, validTime time.Time) (types.TaskI, error) {
	var txTask *types.TxTask
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
		txTask = k.BuildTxTaskWithExpire(ctx, txHash, creator, bounty, validTime, types.TaskStatusPending)
	}

	if validTime.After(txTask.Expiration) {
		return nil, types.ErrTooLateValidTime
	}
	return txTask, nil
}

func (k Keeper) GetShortcutTasks(ctx sdk.Context) []types.TaskID {
	bz := ctx.KVStore(k.storeKey).Get(types.ShortcutTasksKeyPrefix)
	var taskIDsProto types.TaskIDs
	if bz != nil {
		k.cdc.MustUnmarshalLengthPrefixed(bz, &taskIDsProto)
	}
	return taskIDsProto.TaskIds
}

func (k Keeper) SetShortcutTasks(ctx sdk.Context, tid []byte) {
	// here, the uniqueness of tid is not checked deliberately
	// later on when iterating the ShortcutTasks, the task
	// status will be checked to make sure every one is touched once.
	taskIDs := append(k.GetShortcutTasks(ctx), types.TaskID{Tid: tid})
	bz := k.cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})
	ctx.KVStore(k.storeKey).Set(types.ShortcutTasksKeyPrefix, bz)
}

func (k Keeper) TryShortcut(ctx sdk.Context, task types.TaskI) {
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
	var respondedCollateral = sdk.ZeroInt()
	for _, response := range task.GetResponses() {
		operatorAddr := sdk.MustAccAddressFromBech32(response.Operator)
		amount, err := k.GetCollateralAmount(ctx, operatorAddr)
		if err != nil {
			continue
		}
		respondedCollateral = respondedCollateral.Add(amount)
	}

	taskParams := k.GetTaskParams(ctx)
	if sdk.NewDecFromInt(respondedCollateral).
		Quo(sdk.NewDecFromInt(totalCollateral[0].Amount)).
		GTE(taskParams.ShortcutQuorum) {
		k.SetShortcutTasks(ctx, task.GetID())
		k.DeleteFromClosingTaskIDs(ctx, task)
	}
}

// RemoveTask removes a task from kvstore if it is closed, expired and requested by its creator.
// The id of the removed task may still remain in the ExpireTaskIDsStore.
//
//	in such case, when it's expired, the unfound task will be simply skipped
func (k Keeper) RemoveTask(ctx sdk.Context, taskID []byte, force bool, deleter sdk.AccAddress) error {
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
		k.DeleteFromExpireIDs(ctx, *txTask)
	}
	return nil
}

// IteratorAllTasks iterates over all the stored tasks and performs a callback function.
func (k Keeper) IteratorAllTasks(ctx sdk.Context, callback func(task types.TaskI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TaskStoreKeyPrefix)

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
func (k Keeper) GetAllTasks(ctx sdk.Context) (tasks []types.TaskI) {
	k.IteratorAllTasks(ctx, func(task types.TaskI) bool {
		tasks = append(tasks, task)
		return false
	})
	return
}

// UpdateAndGetAllTasks updates all tasks and returns them.
func (k Keeper) UpdateAndGetAllTasks(ctx sdk.Context) (tasks []types.Task, txTasks []types.TxTask) {
	k.IteratorAllTasks(ctx, func(task types.TaskI) bool {
		if t, ok := task.(*types.Task); ok {
			t.WaitingBlocks = t.ExpireHeight - ctx.BlockHeight()
			tasks = append(tasks, *t)
		} else if t, ok := task.(*types.TxTask); ok {
			txTasks = append(txTasks, *t)
		}
		return false
	})
	return
}

// IsValidResponse returns error if a response is not valid.
func (k Keeper) IsValidResponse(ctx sdk.Context, task types.TaskI, response types.Response) error {
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

func (k Keeper) HandleNoneTxTaskForResponse(ctx sdk.Context, txHash []byte) error {
	if _, err := k.GetTask(ctx, txHash); err != nil {
		//if the corresponding TxTask doesn't exit,
		//create one as a placeholder (statue being set as TaskStatusNil),
		//waiting for the MsgCreateTxTask coming to fill in necessary fields
		txTask := k.BuildTxTaskWithExpire(ctx, txHash, "", nil, time.Time{}, types.TaskStatusNil)
		return k.CreateTask(ctx, nil, txTask)
	}
	return nil
}

// RespondToTask records the response from an operator for a task.
func (k Keeper) RespondToTask(ctx sdk.Context, taskID []byte, score int64, operatorAddress sdk.AccAddress) error {
	if !k.IsOperator(ctx, operatorAddress) {
		return types.ErrUnqualifiedOperator
	}

	task, err := k.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	response := types.NewResponse(sdk.NewInt(score), operatorAddress)
	err = k.IsValidResponse(ctx, task, response)
	if err != nil {
		return err
	}

	task.AddResponse(response)
	k.SetTask(ctx, task)

	if _, ok := task.(*types.TxTask); ok {
		k.TryShortcut(ctx, task)
	}

	return nil
}

// Aggregate does an aggregation of responses for a task and updated task result.
func (k Keeper) Aggregate(ctx sdk.Context, taskID []byte) error {
	taskParams := k.GetTaskParams(ctx)
	task, err := k.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	if task.GetStatus() != types.TaskStatusPending {
		return types.ErrTaskClosed
	}

	result := taskParams.AggregationResult
	totalCollateral := sdk.NewInt(0)
	minScoreCollateral := sdk.NewInt(0)
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
					responses[i].Weight = sdk.NewInt(0)
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
	k.SetTask(ctx, task)
	return nil
}

// TotalValidTaskCollateral calculates the total amount of valid collateral of a task.
func (k Keeper) TotalValidTaskCollateral(ctx sdk.Context, task types.TaskI) math.Int {
	taskParams := k.GetTaskParams(ctx)
	totalValidTaskCollateral := sdk.NewInt(0)
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
func (k Keeper) DistributeBounty(ctx sdk.Context, task types.TaskI) error {
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
	k.SetTask(ctx, task)
	return nil
}

func (k Keeper) RefundBounty(ctx sdk.Context, task types.TaskI) error {
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
