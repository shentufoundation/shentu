package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

type UnbondingInfo struct {
	Delegator      sdk.AccAddress
	Validator      sdk.ValAddress
	CompletionTime time.Time
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
		for i := 0; i < len(ubd.Entries); i++ {
			sum = sum.Add(ubd.Entries[i].Balance)
		}
	}
	return sum
}

func (k Keeper) ComputeUnbondingAmountByTime(ctx sdk.Context, provider sdk.AccAddress, time time.Time) sdk.Int {
	dvPairs := k.GetUnbondingsByProviderMaturingByTime(ctx, provider, time)

	sum := sdk.ZeroInt()
	seen := make([]sdk.ValAddress, 0, len(dvPairs))
	for _, dvPair := range dvPairs {
		valAddr := dvPair.Validator
		if find(seen, valAddr) {
			continue
		}
		seen = append(seen, valAddr)

		// obtain unbonding entries and iterate through them
		ubd, found := k.sk.GetUnbondingDelegation(ctx, provider, valAddr)
		if !found {
			continue //TODO
		}
		for i := 0; i < len(ubd.Entries); i++ {
			entry := ubd.Entries[i]
			if entry.IsMature(time) {
				sum = sum.Add(entry.Balance)
			} else {
				break
			}
		}
	}
	return sum
}

func find(list []sdk.ValAddress, item sdk.ValAddress) bool {
	for _, val := range list {
		if val.Equals(item) {
			return true
		}
	}
	return false
}

func (k Keeper) GetUnbondingsByProviderMaturingByTime(ctx sdk.Context, provider sdk.AccAddress, time time.Time) (results []UnbondingInfo) {
	unbondingTimesliceIterator := k.sk.UBDQueueIterator(ctx, time)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := []staking.DVPair{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)

		for _, ubd := range timeslice {
			if ubd.DelegatorAddress.Equals(provider) {
				completionTime, _ := sdk.ParseTimeBytes(unbondingTimesliceIterator.Key()[1:])
				ubdInfo := UnbondingInfo{
					Delegator:      ubd.DelegatorAddress,
					Validator:      ubd.ValidatorAddress,
					CompletionTime: completionTime,
				}
				results = append(results, ubdInfo)
			}
		}
	}
	return results
}

// DelayWithdraws delays the given amount of withdraws maturing
// before the delay duration until the end of the delay duration.
func (k Keeper) DelayWithdraws(ctx sdk.Context, provider sdk.AccAddress, amount sdk.Int, delay time.Duration) error {
	// Retrieve delay candidates, which are withdraws
	// ending before the delay duration from now.
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

	// Delay withdraws, starting with the candidates
	// with the oldest withdraw completion time.
	remaining := amount
	for i := len(withdraws) - 1; i >= 0; i-- {
		if !remaining.IsPositive() {
			break
		}
		// Remove from withdraw queue.
		timeSlice := k.GetWithdrawQueueTimeSlice(ctx, withdraws[i].CompletionTime)
		if len(timeSlice) > 1 {
			for j := 0; j < len(timeSlice); j++ {
				if timeSlice[j].Address.Equals(provider) {
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

func (k Keeper) DelayUnbonding(ctx sdk.Context, provider sdk.AccAddress, amount sdk.Int, delay time.Duration) error {
	// Retrieve delay candidates, which are unbondings
	// ending before the delay duration from now.
	delayedTime := ctx.BlockTime().Add(delay)
	ubds := k.GetUnbondingsByProviderMaturingByTime(ctx, provider, delayedTime)

	// Delay unbondings, starting with the candidates
	// with the oldest unbonding completion time.
	remaining := amount
	for i := len(ubds) - 1; i >= 0; i-- {
		if !remaining.IsPositive() {
			break
		}
		// Remove from unbonding queue.
		timeSlice := k.sk.GetUBDQueueTimeSlice(ctx, ubds[i].CompletionTime)
		if len(timeSlice) > 1 {
			for j := 0; j < len(timeSlice); j++ {
				if timeSlice[j].DelegatorAddress.Equals(provider) && timeSlice[j].ValidatorAddress.Equals(ubds[i].Validator) {
					timeSlice = append(timeSlice[:j], timeSlice[j+1:]...)
					k.sk.SetUBDQueueTimeSlice(ctx, ubds[i].CompletionTime, timeSlice)
					break
				}
			}
		} else {
			k.sk.RemoveUBDQueue(ctx, ubds[i].CompletionTime)
		}

		unbondingDels, found := k.sk.GetUnbondingDelegation(ctx, provider, ubds[i].Validator)
		if !found {
			panic("unbonding list was not found for the given provider-validator pair")
		}

		found = false
		amount := sdk.ZeroInt()
		for j := 0; j < len(unbondingDels.Entries); j++ {
			if !found && unbondingDels.Entries[j].CompletionTime.Equal(ubds[i].CompletionTime) {
				unbondingDels.Entries[j].CompletionTime = delayedTime
				found = true
				amount = unbondingDels.Entries[j].Balance
			} else if found && unbondingDels.Entries[j].CompletionTime.Before(unbondingDels.Entries[j-1].CompletionTime) {
				unbondingDels.Entries[j-1], unbondingDels.Entries[j] = unbondingDels.Entries[j], unbondingDels.Entries[j-1]
			} else if found {
				break
			}
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
