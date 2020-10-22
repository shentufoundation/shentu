package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

var (
	amplifier = sdk.NewInt(1000000)
)

// SetTask sets a task in KVStore.
func (k Keeper) SetTask(ctx sdk.Context, task types.Task) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TaskStoreKey(task.Contract, task.Function), k.cdc.MustMarshalBinaryLengthPrefixed(task))
}

// DeleteTask deletes a task from KVStore.
func (k Keeper) DeleteTask(ctx sdk.Context, task types.Task) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.TaskStoreKey(task.Contract, task.Function))
	return nil
}

// UpdateAndSetTask updates a task and set it in KVStore.
func (k Keeper) UpdateAndSetTask(ctx sdk.Context, task types.Task) {
	task.ClosingBlock = ctx.BlockHeight() + task.WaitingBlocks
	k.SetTask(ctx, task)
	k.SetClosingBlockStore(ctx, task)
}

// SetClosingBlockStore sets the store of the aggregation block for a task.
func (k Keeper) SetClosingBlockStore(ctx sdk.Context, task types.Task) {
	store := ctx.KVStore(k.storeKey)
	taskIDs := k.GetClosingTaskIDs(ctx, task.ClosingBlock)
	taskIDs = append(taskIDs, types.TaskID{Contract: task.Contract, Function: task.Function})
	store.Set(types.ClosingTaskIDsStoreKey(task.ClosingBlock), k.cdc.MustMarshalBinaryLengthPrefixed(taskIDs))
}

// GetTask returns a task given contract and function.
func (k Keeper) GetTask(ctx sdk.Context, contract, function string) (types.Task, error) {
	TaskData := ctx.KVStore(k.storeKey).Get(types.TaskStoreKey(contract, function))
	if TaskData == nil {
		return types.Task{}, types.ErrTaskNotExists
	}
	var task types.Task
	k.cdc.MustUnmarshalBinaryLengthPrefixed(TaskData, &task)
	return task, nil
}

// GetClosingTaskIDs returns a list of task IDs by the closing block.
func (k Keeper) GetClosingTaskIDs(ctx sdk.Context, closingBlock int64) []types.TaskID {
	closingTaskIDsData := ctx.KVStore(k.storeKey).Get(types.ClosingTaskIDsStoreKey(closingBlock))
	var closingTaskIDs []types.TaskID
	if closingTaskIDsData != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(closingTaskIDsData, &closingTaskIDs)
	}
	return closingTaskIDs
}

// DeleteClosingTaskIDs deletes stores for task IDs closed at given block.
func (k Keeper) DeleteClosingTaskIDs(ctx sdk.Context, closingBlock int64) {
	ctx.KVStore(k.storeKey).Delete(types.ClosingTaskIDsStoreKey(closingBlock))
}

// CreateTask creates a new task.
func (k Keeper) CreateTask(ctx sdk.Context, contract string, function string, bounty sdk.Coins,
	description string, expiration time.Time, creator sdk.AccAddress, waitingBlocks int64) error {
	task, err := k.GetTask(ctx, contract, function)
	if err == nil {
		if task.ClosingBlock > ctx.BlockHeight() {
			return types.ErrTaskNotClosed
		}
		if err := k.DeleteTask(ctx, task); err != nil {
			return err
		}
	}
	closingBlock := ctx.BlockHeight() + waitingBlocks
	task = types.NewTask(contract, function, ctx.BlockHeight(), bounty, description, expiration, creator, closingBlock, waitingBlocks)
	k.SetTask(ctx, task)
	k.SetClosingBlockStore(ctx, task)
	if err := k.CollectBounty(ctx, bounty, creator); err != nil {
		return err
	}
	return nil
}

// RemoveTask removes a task from kvstore if it is closed, expired and requested by its creator.
func (k Keeper) RemoveTask(ctx sdk.Context, contract, function string, force bool, creator sdk.AccAddress) error {
	task, err := k.GetTask(ctx, contract, function)
	if err != nil {
		return err
	}
	if !force && !task.Expiration.Before(ctx.BlockTime()) {
		return types.ErrNotExpired
	}

	if ctx.BlockHeight() <= task.ClosingBlock {
		return types.ErrNotFinished
	}

	// TODO: only creator can delete the task for now
	if !task.Creator.Equals(creator) {
		return types.ErrNotCreator
	}
	err = k.DeleteTask(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

// UpdateWaitingBlocks updates the number of blocks before aggregating results.
func (k Keeper) UpdateWaitingBlocks(ctx sdk.Context, task types.Task) error {
	if task.ClosingBlock < ctx.BlockHeight() {
		return types.ErrTaskClosed
	}
	task.WaitingBlocks = task.ClosingBlock - ctx.BlockHeight()
	k.SetTask(ctx, task)
	return nil
}

// IteratorAllTasks iterates over all the stored tasks and performs a callback function.
func (k Keeper) IteratorAllTasks(ctx sdk.Context, callback func(task types.Task) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TaskStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var task types.Task
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &task)

		if callback(task) {
			break
		}
	}
}

// GetAllTasks gets all tasks.
func (k Keeper) GetAllTasks(ctx sdk.Context) (tasks []types.Task) {
	k.IteratorAllTasks(ctx, func(task types.Task) bool {
		tasks = append(tasks, task)
		return false
	})
	return
}

// UpdateAndGetAllTasks updates all tasks and returns them.
func (k Keeper) UpdateAndGetAllTasks(ctx sdk.Context) (tasks []types.Task) {
	k.IteratorAllTasks(ctx, func(task types.Task) bool {
		_ = k.UpdateWaitingBlocks(ctx, task)
		tasks = append(tasks, task)
		return false
	})
	return
}

// IsValidResponse returns error if a response is not valid.
func (k Keeper) IsValidResponse(ctx sdk.Context, task types.Task, response types.Response) error {
	if ctx.BlockHeight() > task.ClosingBlock {
		return types.ErrTaskClosed
	}
	for _, r := range task.Responses {
		if r.Operator.Equals(response.Operator) {
			return types.ErrDuplicateResponse
		}
	}
	if response.Score.LT(types.MinScore) || response.Score.GT(types.MaxScore) {
		return types.ErrInvalidScore
	}
	return nil
}

// RespondToTask records the response from an operator for a task.
func (k Keeper) RespondToTask(ctx sdk.Context, contract string, function string, score int64, operatorAddress sdk.AccAddress) error {
	if !k.IsOperator(ctx, operatorAddress) {
		return types.ErrUnqualifiedOperator
	}

	task, err := k.GetTask(ctx, contract, function)
	if err != nil {
		return err
	}

	response := types.NewResponse(sdk.NewInt(score), operatorAddress)
	err = k.IsValidResponse(ctx, task, response)
	if err != nil {
		return err
	}

	task.Responses = append(task.Responses, response)
	k.SetTask(ctx, task)

	return nil
}

// Aggregate does an aggregation of responses for a task and updated task result.
func (k Keeper) Aggregate(ctx sdk.Context, contract, function string) error {
	taskParams := k.GetTaskParams(ctx)
	task, err := k.GetTask(ctx, contract, function)
	if err != nil {
		return err
	}

	if task.Status != types.TaskStatusPending {
		return types.ErrTaskClosed
	}

	result := taskParams.AggregationResult
	totalCollateral := sdk.NewInt(0)
	minScoreCollateral := sdk.NewInt(0)
	for i, response := range task.Responses {
		amount, err := k.GetCollateralAmount(ctx, response.Operator)
		if err != nil {
			continue
		}
		result = result.Add(response.Score.Mul(amount))
		task.Responses[i].Weight = amount
		totalCollateral = totalCollateral.Add(amount)
		if response.Score.Equal(types.MinScore) {
			minScoreCollateral = minScoreCollateral.Add(amount)
		}
	}

	if totalCollateral.IsPositive() {
		if minScoreCollateral.MulRaw(3).GTE(totalCollateral) {
			result = minScoreCollateral
			for i, response := range task.Responses {
				if !response.Score.Equal(types.MinScore) {
					task.Responses[i].Weight = sdk.NewInt(0)
				}
			}
		} else {
			result = result.Quo(totalCollateral)
		}
		task.Status = types.TaskStatusSucceeded
	} else {
		task.Status = types.TaskStatusFailed
	}
	task.Result = result
	k.SetTask(ctx, task)
	return nil
}

// TotalValidTaskCollateral calculates the total amount of valid collateral of a task.
func (k Keeper) TotalValidTaskCollateral(ctx sdk.Context, task types.Task) sdk.Int {
	taskParams := k.GetTaskParams(ctx)
	totalValidTaskCollateral := sdk.NewInt(0)
	if task.Result.Equal(types.MinScore) {
		for _, response := range task.Responses {
			if response.Score.Equal(types.MinScore) {
				collateral, err := k.GetCollateralAmount(ctx, response.Operator)
				if err != nil {
					continue
				}
				totalValidTaskCollateral = totalValidTaskCollateral.Add(collateral)
			}
		}
	} else if task.Result.LT(taskParams.ThresholdScore) {
		for _, response := range task.Responses {
			if response.Score.LT(taskParams.ThresholdScore) {
				collateral, err := k.GetCollateralAmount(ctx, response.Operator)
				if err != nil {
					continue
				}
				totalValidTaskCollateral = totalValidTaskCollateral.Add(
					amplifier.Mul(collateral).Quo(response.Score.Add(taskParams.Epsilon1)),
				)
			}
		}
	} else {
		for _, response := range task.Responses {
			if response.Score.GTE(taskParams.ThresholdScore) {
				collateral, err := k.GetCollateralAmount(ctx, response.Operator)
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
func (k Keeper) DistributeBounty(ctx sdk.Context, task types.Task) error {
	taskParams := k.GetTaskParams(ctx)
	totalValidTaskCollateral := k.TotalValidTaskCollateral(ctx, task)
	if totalValidTaskCollateral.IsZero() {
		return types.ErrTaskFailed
	}

	for _, bounty := range task.Bounty {
		if task.Result.Equal(types.MinScore) {
			for i, response := range task.Responses {
				if response.Score.Equal(types.MinScore) {
					collateral, err := k.GetCollateralAmount(ctx, response.Operator)
					if err != nil {
						continue
					}
					amount := bounty.Amount.Mul(collateral).Quo(totalValidTaskCollateral)
					reward := sdk.NewCoins(sdk.NewCoin(bounty.Denom, amount))
					if err := k.AddReward(ctx, response.Operator, reward); err != nil {
						continue
					}
					task.Responses[i].Reward = reward
				}
			}
		} else if task.Result.LT(taskParams.ThresholdScore) {
			for i, response := range task.Responses {
				if response.Score.LT(taskParams.ThresholdScore) {
					collateral, err := k.GetCollateralAmount(ctx, response.Operator)
					if err != nil {
						continue
					}
					amount := bounty.Amount.Mul(
						amplifier.Mul(collateral).Quo(response.Score.Add(taskParams.Epsilon1)),
					).Quo(totalValidTaskCollateral)
					reward := sdk.NewCoins(sdk.NewCoin(task.Bounty[0].Denom, amount))
					if err := k.AddReward(ctx, response.Operator, reward); err != nil {
						continue
					}
					task.Responses[i].Reward = reward
				}
			}
		} else {
			for i, response := range task.Responses {
				if response.Score.GTE(taskParams.ThresholdScore) {
					collateral, err := k.GetCollateralAmount(ctx, response.Operator)
					if err != nil {
						continue
					}
					amount := bounty.Amount.Mul(
						amplifier.Mul(collateral).Quo(types.MaxScore.Sub(response.Score).Add(taskParams.Epsilon2)),
					).Quo(totalValidTaskCollateral)
					reward := sdk.NewCoins(sdk.NewCoin(bounty.Denom, amount))
					if err := k.AddReward(ctx, response.Operator, reward); err != nil {
						continue
					}
					task.Responses[i].Reward = reward
				}
			}
		}
	}
	k.SetTask(ctx, task)
	return nil
}
