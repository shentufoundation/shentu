package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// RegisterInvariants registers all shield invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "IDs", IDsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "Tasks", TasksInvariant(k))
	ir.RegisterRoute(types.ModuleName, "fund-check", FundInvariant(k))
}

type IDSet = map[string]struct{}

func IDsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var expireTaskIDs, valid1TaskIDs, valid2TaskIDs, shortcutTaskIDs IDSet
		var err error
		expireTaskIDs, err = checkDuplicatedIDs(ctx, k, types.ExpireTaskStoreKeyPrefix, true, nil, "expireTaskIDs")
		if err != nil {
			return err.Error(), true
		}
		valid1TaskIDs, err = checkDuplicatedIDs(ctx, k, types.ClosingTaskStoreKeyPrefix, true, nil, "ClosingTaskIDs")
		if err != nil {
			return err.Error(), true
		}
		valid2TaskIDs, err = checkDuplicatedIDs(ctx, k, types.ClosingTaskStoreKeyTimedPrefix, true, valid1TaskIDs, "TimedClosingTaskIDs")
		if err != nil {
			return err.Error(), true
		}
		shortcutTaskIDs, err = checkDuplicatedIDs(ctx, k, types.ShortcutTasksKeyPrefix, false, valid2TaskIDs, "shortcutTaskIDs")
		if err != nil {
			return err.Error(), true
		}

		pendingTaskAmount, txTaskAmount := 0, 0
		k.IteratorAllTasks(ctx, func(task types.TaskI) bool {
			if task.GetStatus() == types.TaskStatusPending {
				pendingTaskAmount++
			}
			if _, ok := task.(*types.TxTask); ok {
				txTaskAmount++
			}
			return false
		})

		if len(valid1TaskIDs)+len(valid2TaskIDs)+len(shortcutTaskIDs) != pendingTaskAmount {
			return "pending task amount doesn't match", true
		}
		if len(expireTaskIDs) != txTaskAmount {
			return "txTask amount doesn't match", true
		}
		return "", false
	}
}

func TasksInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		invariantMsg, broken := "", false
		k.IteratorAllTasks(ctx, func(task types.TaskI) bool {
			if err := checkTask(task); err != nil {
				invariantMsg = err.Error()
				broken = true
				return true
			}
			return false
		})
		return invariantMsg, broken
	}
}

// three invariants are checked
//  1. sum of operator's collater should be equal to stored totalCollateral
//  2. module account balance minus totalCollateral should be larger or equal to
//     sum of operator's reward plus sum of pending task's bounty
//  3. module account balance should be larger or equal to
//     sum of operator's collateral plus sum of withdraw amount
func FundInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		broken := false
		var (
			sumCollateral, sumReward, sumPendingBounty, sumWithdraw sdk.Coins
		)
		k.IterateAllOperators(ctx, func(op types.Operator) bool {
			sumCollateral = sumCollateral.Add(op.Collateral...)
			sumReward = sumReward.Add(op.AccumulatedRewards...)
			return false
		})
		k.IteratorAllTasks(ctx, func(task types.TaskI) bool {
			if task.GetStatus() == types.TaskStatusPending {
				sumPendingBounty = sumPendingBounty.Add(task.GetBounty()...)
			}
			return false
		})
		k.IterateAllWithdraws(ctx, func(wd types.Withdraw) bool {
			sumWithdraw = sumWithdraw.Add(wd.Amount...)
			return false
		})
		totalCollateral, err := k.GetTotalCollateral(ctx)
		oracleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
		moduleBalance := k.bankKeeper.SpendableCoins(ctx, oracleAddress)
		if err != nil {
			return sdk.FormatInvariant(types.ModuleName, "fund-check",
				"cannot get module account balance"), false
		}
		if !totalCollateral.IsEqual(sumCollateral) {
			broken = true
		}
		if !moduleBalance.Sub(totalCollateral...).
			IsAllGTE(sumReward.Add(sumPendingBounty...)) {
			broken = true
		}
		if !moduleBalance.IsAllGTE(sumCollateral.Add(sumWithdraw...)) {
			broken = true
		}
		return sdk.FormatInvariant(types.ModuleName, "fund-check",
			fmt.Sprintf("\n\ttotal collateral amount: %s"+
				"\n\tmodule account balance:  %s"+
				"\n\tsum of operator's collateral:  %s"+
				"\n\tsum of operator's reward:  %s"+
				"\n\tsum of operator's withdraw:  %s"+
				"\n\tsum of pending task's bounty:  %s\n",
				totalCollateral, moduleBalance, sumCollateral,
				sumReward, sumWithdraw, sumPendingBounty)), broken
	}
}

func checkDuplicatedIDs(ctx sdk.Context, k Keeper, prefixKey []byte, checkSelf bool, extraSet IDSet, typeStr string) (IDSet, error) {
	idSet := make(IDSet)
	for _, tid := range getAllTaskIDs(ctx, k, prefixKey) {
		if _, ok := idSet[tid]; ok && checkSelf {
			return nil, fmt.Errorf("found duplicated ids(%s) in %s", tid, typeStr)
		}
		if extraSet != nil {
			if _, ok := extraSet[tid]; ok {
				return nil, fmt.Errorf("found duplicated ids(%s) in %s and extra", tid, typeStr)
			}
		}
		idSet[tid] = struct{}{}
	}
	return idSet, nil
}

func getAllTaskIDs(ctx sdk.Context, k Keeper, prefixKey []byte) (taskIDs []string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, prefixKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var taskIDsProto types.TaskIDs
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &taskIDsProto)
		for _, tidProto := range taskIDsProto.TaskIds {
			taskIDs = append(taskIDs, string(tidProto.Tid))
		}
	}
	return
}

func checkTask(task types.TaskI) error {
	status, respCount := task.GetStatus(), len(task.GetResponses())
	if status == types.TaskStatusNil {
		txTask, ok := task.(*types.TxTask)
		if !ok { // only txTask could be in nil status
			return fmt.Errorf("the nil task(%s) is not a txTask", task.GetID())
		}
		if respCount == 0 {
			return fmt.Errorf("the nil task(%s) has empty response", task.GetID())
		}
		if !txTask.ValidTime.IsZero() {
			return fmt.Errorf("the nil task(%s)'s valid time is not zero", task.GetID())
		}
	}
	if status == types.TaskStatusSucceeded && respCount == 0 {
		return fmt.Errorf("succeeded task(%s) has empty response", task.GetID())
	}
	return nil
}
