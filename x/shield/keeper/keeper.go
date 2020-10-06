package keeper

import (
	"encoding/binary"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	sk           types.StakingKeeper
	supplyKeeper types.SupplyKeeper
	paramSpace   params.Subspace
}

// NewKeeper creates a shield keeper.
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, sk types.StakingKeeper, supplyKeeper types.SupplyKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		sk:           sk,
		supplyKeeper: supplyKeeper,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (staking.ValidatorI, bool) {
	return k.sk.GetValidator(ctx, addr)
}

// SetLatestPoolID sets the latest pool ID to store.
func (k Keeper) SetNextPoolID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.GetNextPoolIDKey(), bz)
}

// GetNextPoolID gets the latest pool ID from store.
func (k Keeper) GetNextPoolID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.GetNextPoolIDKey())
	return binary.LittleEndian.Uint64(opBz)
}

// GetPoolBySponsor search store for a pool object with given pool ID.
func (k Keeper) GetPoolBySponsor(ctx sdk.Context, sponsor string) (types.Pool, error) {
	ret := types.Pool{
		PoolID: 0,
	}
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		if pool.Sponsor == sponsor {
			ret = pool
			return true
		} else {
			return false
		}
	})
	if ret.PoolID == 0 {
		return ret, types.ErrNoPoolFound
	}
	return ret, nil
}

// DepositNativePremium deposits premium in native tokens from the shield admin or purchasers.
func (k Keeper) DepositNativePremium(ctx sdk.Context, premium sdk.Coins, from sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, premium)
}

func (k Keeper) InsertWithdrawalQueue(ctx sdk.Context, withdrawal types.Withdrawal, completionTime time.Time) {
	timeSlice := k.GetWithdrawalQueueTimeSlice(ctx, completionTime)
	timeSlice = append(timeSlice, withdrawal)
	k.SetWithdrawalQueueTimeSlice(ctx, completionTime, timeSlice)
}

// GetWithdrawalQueueTimeSlice gets a specific withdrawal queue timeslice,
// which is a slice of withdrawals corresponding to a given time.
func (k Keeper) GetWithdrawalQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []types.Withdrawal {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetWithdrawalCompletionTimeKey(timestamp))
	if bz == nil {
		return []types.Withdrawal{}
	}
	var withdrawals []types.Withdrawal
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &withdrawals)
	return withdrawals
}

func (k Keeper) SetWithdrawalQueueTimeSlice(ctx sdk.Context, timestamp time.Time, withdrawals []types.Withdrawal) {
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

	var withdrawals []types.Withdrawal
	for ; withdrawalTimesliceIterator.Valid(); withdrawalTimesliceIterator.Next() {
		timeslice := []types.Withdrawal{}
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
			panic("withdrawal collateral not found!")
		}
		collateral.Amount = collateral.Amount.Sub(withdrawal.Amount)
		collateral.Withdrawal = collateral.Withdrawal.Sub(withdrawal.Amount)
		if collateral.Amount.IsZero() {
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
		if withdrawal.Amount.IsAnyGT(provider.Collateral) {
			panic("withdrawal amount is greater than the provider's total collateral amount")
		}

		provider.Collateral = provider.Collateral.Sub(withdrawal.Amount)
		provider.Available = provider.Available.Add(withdrawal.Amount.AmountOf(k.sk.BondDenom(ctx)))
		k.SetProvider(ctx, withdrawal.Address, provider)
	}
}
