package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

type unbondingInfo struct {
	delegator      sdk.AccAddress
	validator      sdk.ValAddress
	completionTime time.Time
}

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
	return store.Iterator(types.WithdrawQueueKey, sdk.InclusiveEndBytes(types.GetWithdrawCompletionTimeKey(endTime)))
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

// GetWithdrawsByProvider gets all withdraws of a provider.
func (k Keeper) GetWithdrawsByProvider(ctx sdk.Context, providerAddr sdk.AccAddress) types.Withdraws {
	var withdraws types.Withdraws
	k.IterateWithdraws(ctx, func(timeSlice types.Withdraws) bool {
		for _, withdraw := range timeSlice {
			if withdraw.Address.Equals(providerAddr) {
				withdraws = append(withdraws, withdraw)
			}
		}
		return false
	})
	return withdraws
}

// RemoveTimeSliceFromWithdrawQueue removes a time slice from the withdraw queue.
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
// from a given provider that will be dequeued from the withdraw
// queue by a given time.
func (k Keeper) ComputeWithdrawAmountByTime(ctx sdk.Context, provider sdk.AccAddress, time time.Time) sdk.Int {
	withdrawTimesliceIterator := k.WithdrawQueueIterator(ctx, time)
	defer withdrawTimesliceIterator.Close()

	amount := sdk.ZeroInt()
	for ; withdrawTimesliceIterator.Valid(); withdrawTimesliceIterator.Next() {
		timeslice := []types.Withdraw{}
		value := withdrawTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)

		for _, withdraw := range timeslice {
			if withdraw.Address.Equals(provider) {
				amount = amount.Add(withdraw.Amount)
			}
		}
	}
	return amount
}

func (k Keeper) ComputeTotalUnbondingAmount(ctx sdk.Context, provider sdk.AccAddress) sdk.Int {
	unbondings := k.sk.GetAllUnbondingDelegations(ctx, provider)

	sum := sdk.ZeroInt()
	for _, ubd := range unbondings {
		for _, entry := range ubd.Entries {
			sum = sum.Add(entry.Balance)
		}
	}
	return sum
}

func (k Keeper) ComputeUnbondingAmountByTime(ctx sdk.Context, provider sdk.AccAddress, time time.Time) sdk.Int {
	dvPairs := k.getUnbondingsByProviderMaturingByTime(ctx, provider, time)

	sum := sdk.ZeroInt()
	seen := make(map[string]bool)
	for _, dvPair := range dvPairs {
		valAddr := dvPair.validator
		if seen[valAddr.String()] {
			continue
		}
		seen[valAddr.String()] = true

		// obtain unbonding entries and iterate through them
		ubd, found := k.sk.GetUnbondingDelegation(ctx, provider, valAddr)
		if !found {
			continue //TODO
		}
		for i := 0; i < len(ubd.Entries); i++ {
			entry := ubd.Entries[i]
			if !entry.IsMature(time) {
				break
			}
			sum = sum.Add(entry.Balance)
		}
	}
	return sum
}

func (k Keeper) getUnbondingsByProviderMaturingByTime(ctx sdk.Context, provider sdk.AccAddress, time time.Time) (results []unbondingInfo) {
	unbondingTimesliceIterator := k.sk.UBDQueueIterator(ctx, time)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := []staking.DVPair{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)

		for _, ubd := range timeslice {
			if ubd.DelegatorAddress.Equals(provider) {
				completionTime, _ := sdk.ParseTimeBytes(unbondingTimesliceIterator.Key()[1:])
				ubdInfo := unbondingInfo{
					delegator:      ubd.DelegatorAddress,
					validator:      ubd.ValidatorAddress,
					completionTime: completionTime,
				}
				results = append(results, ubdInfo)
			}
		}
	}
	return results
}

// DelayWithdraws delays the given amount of withdraws maturing
// before the delay duration until the end of the delay duration.
func (k Keeper) DelayWithdraws(ctx sdk.Context, provider sdk.AccAddress, amount sdk.Int, delayedTime time.Time) error {
	// Retrieve delay candidates, which are withdraws
	// ending before the delay duration from now.
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

	// Delay withdraws, starting with the candidates
	// with the oldest withdraw completion time.
	remaining := amount
	for i := len(withdraws) - 1; i >= 0; i-- {
		if !remaining.IsPositive() {
			break
		}
		// Remove from withdraw queue.
		if timeSlice := k.GetWithdrawQueueTimeSlice(ctx, withdraws[i].CompletionTime); len(timeSlice) > 1 {
			for j := len(timeSlice) - 1; j >= 0; j-- {
				if timeSlice[j].Address.Equals(provider) && timeSlice[j].Amount.Equal(withdraws[i].Amount) {
					timeSlice = append(timeSlice[:j], timeSlice[j+1:]...)
					k.SetWithdrawQueueTimeSlice(ctx, withdraws[i].CompletionTime, timeSlice)
					break
				}
			}
		} else {
			k.RemoveTimeSliceFromWithdrawQueue(ctx, withdraws[i].CompletionTime)
		}

		// Adjust the withdraw completion time and re-insert.
		withdraws[i].CompletionTime = delayedTime
		k.InsertWithdrawQueue(ctx, withdraws[i])

		remaining = remaining.Sub(withdraws[i].Amount)
	}

	if remaining.IsPositive() {
		panic("failed to delay enough withdraws")
	}

	return nil
}

func (k Keeper) DelayUnbonding(ctx sdk.Context, provider sdk.AccAddress, amount sdk.Int, delayedTime time.Time) error {
	// Retrieve delay candidates, which are unbondings
	// ending before the delay duration from now.
	ubds := k.getUnbondingsByProviderMaturingByTime(ctx, provider, delayedTime)

	// Delay unbondings, starting with the candidates
	// with the oldest unbonding completion time.
	remaining := amount
	for i := len(ubds) - 1; i >= 0; i-- {
		if !remaining.IsPositive() {
			break
		}
		// Remove from unbonding queue.
		if timeSlice := k.sk.GetUBDQueueTimeSlice(ctx, ubds[i].completionTime); len(timeSlice) > 1 {
			for j := len(timeSlice) - 1; j >= 0; j-- {
				if timeSlice[j].DelegatorAddress.Equals(provider) && timeSlice[j].ValidatorAddress.Equals(ubds[i].validator) {
					timeSlice = append(timeSlice[:j], timeSlice[j+1:]...)
					k.sk.SetUBDQueueTimeSlice(ctx, ubds[i].completionTime, timeSlice)
					break
				}
			}
		} else {
			k.sk.RemoveUBDQueue(ctx, ubds[i].completionTime)
		}

		unbondingDels, found := k.sk.GetUnbondingDelegation(ctx, provider, ubds[i].validator)
		if !found {
			panic("unbonding list was not found for the given provider-validator pair")
		}

		found = false
		amount := sdk.ZeroInt()
		for j := 0; j < len(unbondingDels.Entries); j++ {
			if !found {
				if unbondingDels.Entries[j].CompletionTime.Equal(ubds[i].completionTime) {
					unbondingDels.Entries[j].CompletionTime = delayedTime
					found = true
					amount = unbondingDels.Entries[j].Balance
				}
				continue
			}

			if unbondingDels.Entries[j].CompletionTime.Before(unbondingDels.Entries[j-1].CompletionTime) {
				unbondingDels.Entries[j-1], unbondingDels.Entries[j] = unbondingDels.Entries[j], unbondingDels.Entries[j-1]
				continue
			}

			break
		}
		if !found {
			panic("particular unbonding entry not found for the given timestamp")
		}

		k.sk.InsertUBDQueue(ctx, unbondingDels, delayedTime)
		k.sk.SetUnbondingDelegation(ctx, unbondingDels)

		remaining = remaining.Sub(amount)
	}

	if remaining.IsPositive() {
		panic("failed to delay enough unbondings")
	}

	return nil
}
