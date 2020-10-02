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

// DepositCollateral deposits a community member's collateral for a pool.
func (k Keeper) DepositCollateral(ctx sdk.Context, from sdk.AccAddress, id uint64, amount sdk.Coins) error {
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return err
	}

	// check eligibility and update participant
	participant, found := k.GetParticipant(ctx, from)
	if !found {
		k.addParticipant(ctx, from)
	}
	participant.Collateral = participant.Collateral.Add(amount...)
	if participant.Collateral.IsAnyGT(participant.DelegationBonded) {
		return types.ErrInsufficientStaking
	}
	k.SetParticipant(ctx, from, participant)

	// update the pool and pool community
	found = false
	for i, collateral := range pool.Community {
		if collateral.Provider.Equals(from) {
			pool.Community[i].Amount = pool.Community[i].Amount.Add(amount...)
			found = true
			break
		}
	}
	if !found {
		pool.Community = append(pool.Community, types.NewCollateral(from, amount))
	}
	pool.TotalCollateral = pool.TotalCollateral.Add(amount...)
	pool.Available = pool.Available.Add(amount.AmountOf(k.sk.BondDenom(ctx)))
	k.SetPool(ctx, pool)

	return nil
}

// GetOnesCollaterals returns a community member's all collaterals.
func (k Keeper) GetOnesCollaterals(ctx sdk.Context, address sdk.AccAddress) (collaterals []types.Collateral) {
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		for _, collateral := range pool.Community {
			if collateral.Provider.Equals(address) {
				collaterals = append(collaterals, collateral)
				break
			}
		}
		return false
	})
	return collaterals
}

// WithdrawCollateral withdraws a community member's collateral for a pool.
func (k Keeper) WithdrawCollateral(ctx sdk.Context, from sdk.AccAddress, id uint64, amount sdk.Coins) error {
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return err
	}

	// check eligibility
	participant, found := k.GetParticipant(ctx, from)
	if !found {
		return types.ErrNoDelegationAmount
	}
	if amount.IsAnyGT(participant.Collateral) {
		return types.ErrInvalidCollateralAmount
	}

	// update the pool available coins, but not pool total collateral or community which should be updated 21 days later
	pool.Available = pool.Available.Sub(amount.AmountOf(k.sk.BondDenom(ctx)))
	k.SetPool(ctx, pool)

	// insert to withdrawal queue
	poolParams := k.GetPoolParams(ctx)
	completionTime := ctx.BlockHeader().Time.Add(poolParams.WithdrawalPeriod)
	withdrawal := types.NewWithdrawal(id, from, amount)
	k.InsertWithdrawalQueue(ctx, withdrawal, completionTime)

	return nil
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

// SetPoolParams sets parameters subspace for shield pool parameters.
func (k Keeper) SetPoolParams(ctx sdk.Context, poolParams types.PoolParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyPoolParams, &poolParams)
}

// GetPoolParams returns shield pool parameters.
func (k Keeper) GetPoolParams(ctx sdk.Context) types.PoolParams {
	var poolParams types.PoolParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyPoolParams, &poolParams)
	return poolParams
}

// SetClaimProposalParams sets parameters subspace for shield claim proposal parameters.
func (k Keeper) SetClaimProposalParams(ctx sdk.Context, claimProposalParams types.ClaimProposalParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyClaimProposalParams, &claimProposalParams)
}

// GetClaimProposalParams returns shield claim proposal parameters.
func (k Keeper) GetClaimProposalParams(ctx sdk.Context) types.ClaimProposalParams {
	var claimProposalParams types.ClaimProposalParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyClaimProposalParams, &claimProposalParams)
	return claimProposalParams
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
			// do not update participant if the pool has been closed
			continue
		}
		pool.TotalCollateral = pool.TotalCollateral.Sub(withdrawal.Amount)
		for i := range pool.Community {
			if pool.Community[i].Provider.Equals(withdrawal.Address) {
				pool.Community[i].Amount = pool.Community[i].Amount.Sub(withdrawal.Amount)
				break
			}
		}
		if pool.CertiK.Provider.Equals(withdrawal.Address) {
			pool.CertiK.Amount = pool.CertiK.Amount.Sub(withdrawal.Amount)
		}
		k.SetPool(ctx, pool)

		// update participant
		participant, found := k.GetParticipant(ctx, withdrawal.Address)
		if !found {
			// TODO will this happen?
			continue
		}
		if withdrawal.Amount.IsAnyGT(participant.Collateral) {
			// TODO will this happen?
			participant.Collateral = sdk.Coins{}
		} else {
			participant.Collateral = participant.Collateral.Sub(withdrawal.Amount)
		}
		k.SetParticipant(ctx, withdrawal.Address, participant)
	}
}
