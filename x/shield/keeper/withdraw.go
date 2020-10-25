package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// InsertWithdrawQueue prepares a withdraw queue timeslice
// for insertion into the queue.
func (k Keeper) InsertWithdrawQueue(ctx sdk.Context, withdraw types.Withdraw) {
	timeSlice := k.GetWithdrawQueueTimeSlice(ctx, withdraw.CompletionTime)
	timeSlice = append(timeSlice, withdraw)
	k.SetWithdrawQueueTimeSlice(ctx, withdraw.CompletionTime, timeSlice)
}

// SetWithdrawQueueTimeSlice stores a withdraw queue timeslice
// using the timestamp as the key.
func (k Keeper) SetWithdrawQueueTimeSlice(ctx sdk.Context, timestamp time.Time, withdraws []types.Withdraw) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(withdraws)
	store.Set(types.GetWithdrawCompletionTimeKey(timestamp), bz)
}

// GetWithdrawQueueTimeSlice gets a specific withdraw queue timeslice,
// which is a slice of withdraws corresponding to a given time.
func (k Keeper) GetWithdrawQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []types.Withdraw {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetWithdrawCompletionTimeKey(timestamp))
	if bz == nil {
		return []types.Withdraw{}
	}
	var withdraws []types.Withdraw
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &withdraws)
	return withdraws
}

// WithdrawQueueIterator returns all the withdraw queue timeslices from time 0 until endTime.
func (k Keeper) WithdrawQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.WithdrawQueueKey,
		sdk.InclusiveEndBytes(types.GetWithdrawCompletionTimeKey(endTime)))
}

// IterateWithdraws iterates through all ongoing withdraws.
func (k Keeper) IterateWithdraws(ctx sdk.Context, callback func(withdraw types.Withdraws) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WithdrawQueueKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		timeslice := types.Withdraws{}
		value := iterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)

		if callback(timeslice) {
			break
		}
	}
}

// GetAllWithdraws gets all collaterals that are being withdrawn.
func (k Keeper) GetAllWithdraws(ctx sdk.Context) (withdraws types.Withdraws) {
	k.IterateWithdraws(ctx, func(withdraw types.Withdraws) bool {
		withdraws = append(withdraws, withdraw...)
		return false
	})
	return withdraws
}

func (k Keeper) RemoveTimeSliceFromWithdrawQueue(ctx sdk.Context, timestamp time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetWithdrawCompletionTimeKey(timestamp))
}

// DequeueCompletedWithdrawQueue dequeues completed withdraws
// and processes their completions.
func (k Keeper) DequeueCompletedWithdrawQueue(ctx sdk.Context) {
	// retrieve completed withdraws from the queue
	store := ctx.KVStore(k.storeKey)
	withdrawTimesliceIterator := k.WithdrawQueueIterator(ctx, ctx.BlockHeader().Time)
	defer withdrawTimesliceIterator.Close()

	var withdraws []types.Withdraw
	for ; withdrawTimesliceIterator.Valid(); withdrawTimesliceIterator.Next() {
		timeslice := []types.Withdraw{}
		value := withdrawTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)
		withdraws = append(withdraws, timeslice...)
		store.Delete(withdrawTimesliceIterator.Key())
	}

	// For each completed withdraw, process adjustments.
	totalCollateral := k.GetTotalCollateral(ctx)
	totalWithdrawing := k.GetTotalWithdrawing(ctx)
	for _, withdraw := range withdraws {
		provider, found := k.GetProvider(ctx, withdraw.Address)
		if !found {
			panic("provider not found but its collaterals are being withdrawn")
		}
		provider.Collateral = provider.Collateral.Sub(withdraw.Amount)
		provider.Withdrawing = provider.Withdrawing.Sub(withdraw.Amount)
		k.SetProvider(ctx, withdraw.Address, provider)

		totalCollateral = totalCollateral.Sub(withdraw.Amount)
		totalWithdrawing = totalWithdrawing.Sub(withdraw.Amount)
	}
	k.SetTotalCollateral(ctx, totalCollateral)
	k.SetTotalWithdrawing(ctx, totalWithdrawing)
}

// ComputeWithdrawAmountByTime computes the amount of collaterals
// that will be dequeued from the withdraw queue by a given time.
func (k Keeper) ComputeWithdrawAmountByTime(ctx sdk.Context, time time.Time) sdk.Int {
	withdrawTimesliceIterator := k.WithdrawQueueIterator(ctx, time)
	defer withdrawTimesliceIterator.Close()

	amount := sdk.ZeroInt()
	for ; withdrawTimesliceIterator.Valid(); withdrawTimesliceIterator.Next() {
		timeslice := []types.Withdraw{}
		value := withdrawTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)

		for _, withdraw := range timeslice {
			amount = amount.Add(withdraw.Amount)
		}
	}
	return amount
}

// DelayWithdraws looks at the provider's withdraws ending before the delay
// duration from now and delays the given amount of withdraws by the specified
// delay duration.
func (k Keeper) DelayWithdraws(ctx sdk.Context, delay time.Duration, amount sdk.Int, provider sdk.AccAddress) error {
	// Retrieve all withdrawals ending before the delay duration from now.
	delayedTime := ctx.BlockTime().Add(delay)
	withdrawTimesliceIterator := k.WithdrawQueueIterator(ctx, delayedTime)
	defer withdrawTimesliceIterator.Close()

	withdraws := []types.Withdraw{}
	for ; withdrawTimesliceIterator.Valid(); withdrawTimesliceIterator.Next() {
		timeslice := []types.Withdraw{}
		value := withdrawTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)

		for _, withdraw := range timeslice {
			if withdraw.Address.Equals(provider) {
				withdraws = append(withdraws, withdraw)
			}
		}
	}

	// Delay withdrawals.
	// TODO: Withdraw the exact amount?
	remaining := amount
	for _, withdraw := range withdraws {
		if !remaining.IsPositive() {
			break
		}

		// Remove from withdraw queue.
		timeSlice := k.GetWithdrawQueueTimeSlice(ctx, withdraw.CompletionTime)
		if len(timeSlice) > 1 {
			for i := 0; i < len(timeSlice); i++ {
				if timeSlice[i].Address.Equals(provider) {
					timeSlice = append(timeSlice[:i], timeSlice[i+1:]...)
					k.SetWithdrawQueueTimeSlice(ctx, withdraw.CompletionTime, timeSlice)
					break
				}
			}
		} else {
			k.RemoveTimeSliceFromWithdrawQueue(ctx, withdraw.CompletionTime)
		}

		// Adjust the withdraw end time and re-insert.
		withdraw.CompletionTime = delayedTime
		k.InsertWithdrawQueue(ctx, withdraw)

		// Delay linked unbonding, if exists.
		if !withdraw.LinkedUnbonding.CompletionTime.IsZero() {
			// TODO: only works if withdraw period == unbonding period?
			k.DelayUnbonding(ctx, provider, withdraw.LinkedUnbonding, delayedTime)	
		}

		remaining = remaining.Sub(withdraw.Amount)
	} // for each withdraw

	if remaining.IsPositive() {
		panic("failed to delay enough withdraws") // TODO
	}

	return nil
}

// DelayUnbonding delays the completion time of an unbonding identified
// by provider (delegator) and timestamp (unbonding completion time).
func (k Keeper) DelayUnbonding(ctx sdk.Context, delAddr sdk.AccAddress, ubdInfo types.UnbondingInfo, delayedTime time.Time) {
	valAddr := ubdInfo.ValidatorAddress
	timestamp := ubdInfo.CompletionTime

	unbonding, found := k.sk.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		panic("unbonding list was not found for the given delegator-validator pair")
	}

	// Identify the particular unbonding entry from the unbonding list.
	// TODO: Can we identify the particular UnbondingDelegationEntry with completionTime?
	// That is, there can be no unbonding entry with the same completionTime?
	found = false
	for i := 0; i < len(unbonding.Entries); i++ {
		if !found && unbonding.Entries[i].CompletionTime.Equal(timestamp) {
			unbonding.Entries[i].CompletionTime = delayedTime
			found = true
		} else if found && unbonding.Entries[i].CompletionTime.Before(unbonding.Entries[i-1].CompletionTime) {
			unbonding.Entries[i-1], unbonding.Entries[i] = unbonding.Entries[i], unbonding.Entries[i-1]
		} else if found {
			break
		}
	}
	if !found {
		panic("particular unbonding entry not found for the given timestamp")
	}

	// Update the stores.
	timeSlice := k.sk.GetUBDQueueTimeSlice(ctx, timestamp)
	timeSliceLength := len(timeSlice)
	if timeSliceLength == 0 {
		panic("unbonding was not found from the unbonding queue")
	}
	for i := 0; i < timeSliceLength; i++ {
		if timeSlice[i].DelegatorAddress.Equals(delAddr) {
			if timeSliceLength > 1 {
				timeSlice = append(timeSlice[:i], timeSlice[i+1:]...)
				k.sk.SetUBDQueueTimeSlice(ctx, timestamp, timeSlice)
			} else {
				k.sk.RemoveUBDQueue(ctx, timestamp)
			}
			k.sk.InsertUBDQueue(ctx, unbonding, delayedTime)
			k.sk.SetUnbondingDelegation(ctx, unbonding)
			return
		}
	}

	panic("particular unbonding entry not found from the unbonding queue")
}
