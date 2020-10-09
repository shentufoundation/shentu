package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) InsertWithdrawalQueue(ctx sdk.Context, withdrawal types.Withdraw) {
	timeSlice := k.GetWithdrawalQueueTimeSlice(ctx, withdrawal.CompletionTime)
	timeSlice = append(timeSlice, withdrawal)
	k.SetWithdrawalQueueTimeSlice(ctx, withdrawal.CompletionTime, timeSlice)
}

// GetWithdrawalQueueTimeSlice gets a specific withdrawal queue timeslice,
// which is a slice of withdrawals corresponding to a given time.
func (k Keeper) GetWithdrawalQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []types.Withdraw {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetWithdrawalCompletionTimeKey(timestamp))
	if bz == nil {
		return []types.Withdraw{}
	}
	var withdrawals []types.Withdraw
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &withdrawals)
	return withdrawals
}

func (k Keeper) SetWithdrawalQueueTimeSlice(ctx sdk.Context, timestamp time.Time, withdrawals []types.Withdraw) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(withdrawals)
	store.Set(types.GetWithdrawalCompletionTimeKey(timestamp), bz)
}

// WithdrawalQueueIterator returns all the withdrawal queue timeslices from time 0 until endTime
func (k Keeper) WithdrawalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.WithdrawalQueueKey,
		sdk.InclusiveEndBytes(types.GetWithdrawalCompletionTimeKey(endTime)))
}

func (k Keeper) DequeueCompletedWithdrawalQueue(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	withdrawalTimesliceIterator := k.WithdrawalQueueIterator(ctx, ctx.BlockHeader().Time)
	defer withdrawalTimesliceIterator.Close()

	var withdrawals []types.Withdraw
	for ; withdrawalTimesliceIterator.Valid(); withdrawalTimesliceIterator.Next() {
		timeslice := []types.Withdraw{}
		value := withdrawalTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)
		withdrawals = append(withdrawals, timeslice...)
		store.Delete(withdrawalTimesliceIterator.Key())
	}

	for _, withdrawal := range withdrawals {
		// update pool community or CertiK first in case the pool is closed
		pool, err := k.GetPool(ctx, withdrawal.PoolID)
		if err != nil {
			// do not update provider if the pool has been closed
			continue
		}
		pool.TotalCollateral = pool.TotalCollateral.Sub(withdrawal.Amount)
		collateral, found := k.GetCollateral(ctx, pool, withdrawal.Address)
		if !found {
			panic("withdrawal collateral not found")
		}
		// allow collateral amount to be negative, which could happen when it is locked
		collateral.Amount = collateral.Amount.Sub(withdrawal.Amount)
		collateral.Withdrawing = collateral.Withdrawing.Sub(withdrawal.Amount)
		if collateral.Amount.IsNegative() && len(collateral.LockedCollaterals) == 0 {
			store.Delete(types.GetCollateralKey(pool.PoolID, collateral.Provider))
		} else {
			k.SetCollateral(ctx, pool, collateral.Provider, collateral)
		}
		k.SetPool(ctx, pool)

		// update provider's collateral amount
		provider, found := k.GetProvider(ctx, withdrawal.Address)
		if !found {
			panic("provider not found but its collaterals are being withdrawn")
		}

		provider.Collateral = provider.Collateral.Sub(withdrawal.Amount)
		provider.Available = provider.Available.Add(withdrawal.Amount)
		provider.Withdrawing = provider.Withdrawing.Sub(withdrawal.Amount)
		k.SetProvider(ctx, withdrawal.Address, provider)
	}
}

// IterateWithdraws iterates through all ongoing withdraws.
func (k Keeper) IterateWithdraws(ctx sdk.Context, callback func(withdraw types.Withdraws) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WithdrawalQueueKey)

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

// GetAllWithdraws gets all collaterals.
func (k Keeper) GetAllWithdraws(ctx sdk.Context) (withdraws types.Withdraws) {
	k.IterateWithdraws(ctx, func(withdraw types.Withdraws) bool {
		withdraws = append(withdraws, withdraw...)
		return false
	})
	return withdraws
}
