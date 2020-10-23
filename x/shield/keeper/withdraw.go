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
