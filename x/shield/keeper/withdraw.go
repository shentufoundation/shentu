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

// WithdrawQueueIterator returns all the withdraw queue timeslices from time 0 until endTime
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
	// Retrieve completed withdraws from the queue.
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

	for _, withdraw := range withdraws {
		provider, found := k.GetProvider(ctx, withdraw.Address)
		if !found {
			panic("provider not found but its collaterals are being withdrawn")
		}
		pool, err := k.GetPool(ctx, withdraw.PoolID)
		if err != nil {
			// Pools, collaterals and providers have been updated for closed pools.
			continue
		}
		collateral, found := k.GetCollateral(ctx, pool, withdraw.Address)
		if !found {
			panic("withdraw collateral not found")
		}

		// Update withdrawing.
		collateral.Withdrawing = collateral.Withdrawing.Sub(withdraw.Amount)
		provider.Withdrawing = provider.Withdrawing.Sub(withdraw.Amount)

		// Update collateral.
		// It is possible that the withdraw amount exceeds the collateral amount because of locking collaterals by claim proposals.
		// Do not allow collateral amount to be negative. Set overdraft for this situation.
		var validWithdrawAmount sdk.Int
		if collateral.Amount.GTE(withdraw.Amount) {
			validWithdrawAmount = withdraw.Amount
		} else {
			validWithdrawAmount = collateral.Amount
			collateral.Overdraft = collateral.Overdraft.Add(withdraw.Amount.Sub(collateral.Amount))
		}
		if collateral.Overdraft.GT(collateral.TotalLocked) {
			panic("overdraft amount is greater than locked amount")
		}
		pool.TotalCollateral = pool.TotalCollateral.Sub(validWithdrawAmount)
		collateral.Amount = collateral.Amount.Sub(validWithdrawAmount)
		provider.Collateral = provider.Collateral.Sub(validWithdrawAmount)

		// Update provider's available delegations.
		provider.Available = provider.Available.Add(validWithdrawAmount)

		if collateral.Amount.IsZero() && len(collateral.LockedCollaterals) == 0 {
			store.Delete(types.GetCollateralKey(pool.PoolID, collateral.Provider))
		} else {
			k.SetCollateral(ctx, pool, collateral.Provider, collateral)
		}
		k.SetPool(ctx, pool)
		k.SetProvider(ctx, withdraw.Address, provider)
	}
}
