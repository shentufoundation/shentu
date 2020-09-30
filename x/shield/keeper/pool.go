package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.GetPoolKey(pool.PoolID), bz)
}

func (k Keeper) GetPool(ctx sdk.Context, id uint64) (types.Pool, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolKey(id))
	if bz != nil {
		var pool types.Pool
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
		return pool, nil
	}
	return types.Pool{}, types.ErrNoPoolFound
}

func (k Keeper) CreatePool(
	ctx sdk.Context, creator sdk.AccAddress, shield sdk.Coins, deposit types.MixedCoins, sponsor string,
	timeOfCoverage, blocksOfCoverage int64) (types.Pool, error) {
	operator := k.GetAdmin(ctx)
	if !creator.Equals(operator) {
		return types.Pool{}, types.ErrNotShieldOperator
	}
	if err := k.DepositNativePremium(ctx, deposit.Native, creator); err != nil {
		return types.Pool{}, err
	}

	if !k.ValidatePoolDuration(ctx, timeOfCoverage, blocksOfCoverage) {
		return types.Pool{}, types.ErrPoolLifeTooShort
	}
	// check if shield is backed by operator's delegations
	participant, found := k.GetParticipant(ctx, operator)
	if !found {
		return types.Pool{}, types.ErrNoDelegationAmount
	}
	participant.TotalCollateral = participant.TotalCollateral.Add(shield...)
	if participant.TotalCollateral.IsAnyGT(participant.TotalDelegation) {
		return types.Pool{}, types.ErrInsufficientStaking
	}

	// Store endTime. If not available, store endBlockHeight.
	var endTime, endBlockHeight int64
	startBlockHeight := ctx.BlockHeight()
	if timeOfCoverage != 0 {
		endTime = ctx.BlockTime().Unix() + timeOfCoverage
	} else if blocksOfCoverage != 0 {
		endBlockHeight = startBlockHeight + blocksOfCoverage
	}

	id := k.GetNextPoolID(ctx)
	depositDec := types.MixedDecCoinsFromMixedCoins(deposit)

	pool := types.NewPool(operator, shield, depositDec, sponsor, endTime, startBlockHeight, endBlockHeight, id)

	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, id+1)
	k.SetParticipant(ctx, operator, participant)

	return pool, nil
}

func (k Keeper) UpdatePool(
	ctx sdk.Context, updater sdk.AccAddress, shield sdk.Coins, deposit types.MixedCoins, id uint64,
	additionalTime, additionalBlocks int64) (types.Pool, error) {
	operator := k.GetAdmin(ctx)
	if !updater.Equals(operator) {
		return types.Pool{}, types.ErrNotShieldOperator
	}

	// check if shield is backed by operator's delegations
	participant, found := k.GetParticipant(ctx, operator)
	if !found {
		return types.Pool{}, types.ErrNoDelegationAmount
	}
	participant.TotalCollateral = participant.TotalCollateral.Add(shield...)
	if participant.TotalCollateral.IsAnyGT(participant.TotalDelegation) {
		return types.Pool{}, types.ErrInsufficientStaking
	}

	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return types.Pool{}, err
	}

	newCoverageTime := additionalTime + pool.EndTime - ctx.BlockTime().Unix()
	newCoverageBlocks := additionalBlocks + pool.EndBlockHeight - ctx.BlockHeight()
	if !k.ValidatePoolDuration(ctx, newCoverageTime, newCoverageBlocks) {
		return types.Pool{}, types.ErrPoolLifeTooShort
	}
	// Extend EndTime. If not available, extend EndBlockHeight.
	if additionalTime != 0 {
		if pool.EndTime == 0 {
			return types.Pool{}, types.ErrCannotExtend
		}
		pool.EndTime += additionalTime
	} else if additionalBlocks != 0 {
		if pool.EndBlockHeight == 0 {
			return types.Pool{}, types.ErrCannotExtend
		}
		pool.EndBlockHeight += additionalBlocks
	}

	pool.Shield = pool.Shield.Add(shield...)
	pool.Premium = pool.Premium.Add(types.MixedDecCoinsFromMixedCoins(deposit))
	k.SetPool(ctx, pool)
	k.SetParticipant(ctx, operator, participant)

	return pool, nil
}

func (k Keeper) PausePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (types.Pool, error) {
	operator := k.GetAdmin(ctx)
	if !updater.Equals(operator) {
		return types.Pool{}, types.ErrNotShieldOperator
	}
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return types.Pool{}, err
	}
	if !pool.Active {
		return types.Pool{}, types.ErrPoolAlreadyPaused
	}
	pool.Active = false
	k.SetPool(ctx, pool)
	return pool, nil
}

func (k Keeper) ResumePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (types.Pool, error) {
	operator := k.GetAdmin(ctx)
	if !updater.Equals(operator) {
		return types.Pool{}, types.ErrNotShieldOperator
	}
	pool, err := k.GetPool(ctx, id)
	if err != nil {
		return types.Pool{}, err
	}
	if pool.Active {
		return types.Pool{}, types.ErrPoolAlreadyActive
	}
	pool.Active = true
	k.SetPool(ctx, pool)
	return pool, nil
}

// GetAllPools retrieves all pools in the store.
func (k Keeper) GetAllPools(ctx sdk.Context) (pools []types.Pool) {
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		pools = append(pools, pool)
		return false
	})
	return pools
}

// PoolEnded returns if pool has reached ending time and block height
func (k Keeper) PoolEnded(ctx sdk.Context, pool types.Pool) bool {
	if ctx.BlockTime().Unix() > pool.EndTime && ctx.BlockHeight() > pool.EndBlockHeight {
		return true
	}
	return false
}

// ClosePool closes the pool
func (k Keeper) ClosePool(ctx sdk.Context, pool types.Pool) {
	// TODO: make sure nothing else needs to be done
	k.FreeCollateral(ctx, pool)
}

// FreeCollateral frees collaterals depositted in a pool.
func (k Keeper) FreeCollateral(ctx sdk.Context, pool types.Pool) {
	participants := append(pool.Community, pool.CertiK)
	for _, member := range participants {
		participant, _ := k.GetParticipant(ctx, member.Provider)
		participant.TotalCollateral = participant.TotalCollateral.Sub(member.Amount)
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolKey(pool.PoolID))
}

// IterateAllPools iterates over the all the stored pools and performs a callback function.
func (k Keeper) IterateAllPools(ctx sdk.Context, callback func(certificate types.Pool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PoolKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &pool)

		if callback(pool) {
			break
		}
	}
}

// ValidatePoolDuration validates new pool duration to be valid
func (k Keeper) ValidatePoolDuration(ctx sdk.Context, timeDuration, numBlocks int64) bool {
	poolparams := k.GetPoolParams(ctx)
	minPoolDuration := int64(poolparams.MinPoolLife)
	return timeDuration > minPoolDuration || numBlocks*5 > minPoolDuration
}
