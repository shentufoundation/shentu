package keeper

import (
	"errors"
	"time"

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
func (k Keeper) UpdateAndSetTask(ctx sdk.Context, task types.TaskI) {
	if task.IsValid(ctx) {
		k.SetClosingBlockStore(ctx, task)
	}
	if scTask, ok := task.(*types.Task); ok {
		scTask.ExpireHeight = ctx.BlockHeight() + scTask.WaitingBlocks
		k.SetTask(ctx, scTask)
	} else {
		k.SetTask(ctx, task)
	}
}

// GetTask returns a task given contract and function.
func (k Keeper) GetTask(ctx sdk.Context, taskID []byte) (task types.TaskI, err error) {
	TaskData := ctx.KVStore(k.storeKey).Get(types.TaskStoreKey(taskID))
	if TaskData == nil {
		return nil, types.ErrTaskNotExists
	}
	err = k.cdc.UnmarshalInterface(TaskData, &task)
	return
}

// SetClosingBlockStore sets the store of the aggregation block for a task.
func (k Keeper) SetClosingBlockStore(ctx sdk.Context, task types.TaskI) {
	store := ctx.KVStore(k.storeKey)

	newTaskID := types.TaskID{Tid: task.GetID()}
	taskIDs := append(k.GetClosingTaskIDs(ctx, task), newTaskID)

	bz := k.cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})
	switch task := task.(type) {
	case *types.Task:
		store.Set(types.ClosingTaskIDsStoreKey(task.ExpireHeight), bz)
		return
	case *types.TxTask:
		store.Set(types.ClosingTaskIDsTimedStoreKey(task.ValidTime), bz)
		return
	default:
		panic(errors.New("oracle: unknown implementation of TaskI"))
	}
}

// GetClosingTaskIDs returns a list of task IDs by the closing block.
func (k Keeper) GetClosingTaskIDs(ctx sdk.Context, task types.TaskI) (resIDs []types.TaskID) {
	height, endTime := ctx.BlockHeight(), ctx.BlockTime()
	if task != nil {
		height, endTime = task.GetValidTime()
	}
	if height > 0 {
		resIDs = append(resIDs, k.GetClosingTaskIDsByHeight(ctx, height)...)
	}
	if !endTime.IsZero() {
		resIDs = append(resIDs, k.GetClosingTaskIDsByTime(ctx, endTime)...)
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

func (k Keeper) IteratorClosingTaskIDsByTime(ctx sdk.Context, endTime time.Time, callback func(key, value []byte) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.ClosingTaskStoreKeyTimedPrefix,
		sdk.InclusiveEndBytes(types.ClosingTaskIDsTimedStoreKey(endTime)))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if callback(iterator.Key(), iterator.Value()) {
			break
		}
	}
}

func (k Keeper) GetClosingTaskIDsByTime(ctx sdk.Context, endTime time.Time) (resIDs []types.TaskID) {
	k.IteratorClosingTaskIDsByTime(ctx, endTime, func(key, value []byte) bool {
		var taskIDsProto types.TaskIDs
		k.cdc.MustUnmarshalLengthPrefixed(value, &taskIDsProto)
		resIDs = append(resIDs, taskIDsProto.TaskIds...)
		return false
	})
	return
}

// DeleteClosingTaskIDs deletes stores for task IDs closed at given block.
func (k Keeper) DeleteClosingTaskIDs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ClosingTaskIDsStoreKey(ctx.BlockHeight()))
	k.IteratorClosingTaskIDsByTime(ctx, ctx.BlockTime(), func(key, value []byte) bool {
		store.Delete(key)
		return false
	})
}

// CreateTask creates a new task.
func (k Keeper) CreateTask(ctx sdk.Context, creator sdk.AccAddress, task types.TaskI) error {
	savedTask, err := k.GetTask(ctx, task.GetID())
	if err == nil {
		if _, ok := task.(*types.Task); ok {
			if savedTask.IsValid(ctx) {
				return types.ErrTaskNotClosed
			}
			if err := k.DeleteTask(ctx, savedTask); err != nil {
				return err
			}
		} else if reqTask, ok := task.(*types.TxTask); ok {
			oldTask, ok := savedTask.(*types.TxTask)
			if !ok {
				return types.ErrInvalidTask
			}
			if !reqTask.IsValid(ctx) {
				return types.ErrInvalidTask
			}

			if oldTask.GetStatus() != types.TaskStatusNil {
				if savedTask.IsValid(ctx) {
					return types.ErrTaskNotClosed
				}
			} else {
				//created by fast path
				reqTask.Responses = oldTask.Responses
				reqTask.Expiration = oldTask.Expiration
				reqTask.Score = oldTask.Score
			}
			if err := k.DeleteTask(ctx, savedTask); err != nil {
				return err
			}
		}
	}

	k.SetTask(ctx, task)
	k.SetClosingBlockStore(ctx, task)
	if err := k.CollectBounty(ctx, task.GetBounty(), creator); err != nil {
		return err
	}
	return nil
}

// RemoveTask removes a task from kvstore if it is closed, expired and requested by its creator.
func (k Keeper) RemoveTask(ctx sdk.Context, taskID []byte, force bool, creator sdk.AccAddress) error {
	task, err := k.GetTask(ctx, taskID)
	if err != nil {
		return err
	}
	if !force && !task.IsExpired(ctx) {
		return types.ErrNotExpired
	}

	if task.IsValid(ctx) {
		return types.ErrNotFinished
	}

	// TODO: only creator can delete the task for now
	creatorAddr, err := sdk.AccAddressFromBech32(task.GetCreator())
	if err != nil {
		panic(err)
	}
	if !creatorAddr.Equals(creator) {
		return types.ErrNotCreator
	}
	err = k.DeleteTask(ctx, task)
	if err != nil {
		return err
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
func (k Keeper) UpdateAndGetAllTasks(ctx sdk.Context) (tasks []types.TaskI) {
	k.IteratorAllTasks(ctx, func(task types.TaskI) bool {
		if t, ok := task.(*types.Task); ok {
			t.WaitingBlocks = t.ExpireHeight - ctx.BlockHeight()
			tasks = append(tasks, t)
		} else {
			tasks = append(tasks, task)
		}
		return false
	})
	return
}

// IsValidResponse returns error if a response is not valid.
func (k Keeper) IsValidResponse(ctx sdk.Context, task types.TaskI, response types.Response) error {
	if !task.IsValid(ctx) {
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
		operatorAddr, err := sdk.AccAddressFromBech32(response.Operator)
		if err != nil {
			panic(err)
		}
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
			result = minScoreCollateral
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
func (k Keeper) TotalValidTaskCollateral(ctx sdk.Context, task types.TaskI) sdk.Int {
	taskParams := k.GetTaskParams(ctx)
	totalValidTaskCollateral := sdk.NewInt(0)
	responses := task.GetResponses()
	if task.GetScore() == types.MinScore.Int64() {
		for _, response := range responses {
			if response.Score.Equal(types.MinScore) {
				operatorAddr, err := sdk.AccAddressFromBech32(response.Operator)
				if err != nil {
					panic(err)
				}
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
				operatorAddr, err := sdk.AccAddressFromBech32(response.Operator)
				if err != nil {
					panic(err)
				}
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
				operatorAddr, err := sdk.AccAddressFromBech32(response.Operator)
				if err != nil {
					panic(err)
				}
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
					operatorAddr, err := sdk.AccAddressFromBech32(responses[i].Operator)
					if err != nil {
						panic(err)
					}
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
					operatorAddr, err := sdk.AccAddressFromBech32(responses[i].Operator)
					if err != nil {
						panic(err)
					}
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
