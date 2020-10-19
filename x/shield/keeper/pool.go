package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// SetGlobalPool sets data of the global pool in kv-store.
func (k Keeper) SetGlobalPool(ctx sdk.Context, globalPool types.GlobalPool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(globalPool)
	store.Set(types.GetGlobalPoolKey(), bz)
}

// GetGlobalPool gets data of the shield global pool.
func (k Keeper) GetGlobalPool(ctx sdk.Context) types.GlobalPool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGlobalPoolKey())
	if bz == nil {
		panic("global pool is not found")
	}
	var globalPool types.GlobalPool
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &globalPool)
	return globalPool
}

// SetPool sets data of a pool in kv-store.
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(types.GetPoolKey(pool.ID), bz)
}

// GetPool gets data of a pool given pool ID.
func (k Keeper) GetPool(ctx sdk.Context, id uint64) (types.Pool, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolKey(id))
	if bz == nil {
		return types.Pool{}, false
	}
	var pool types.Pool
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	return pool, true
}

// CreatePool creates a pool and sponsor's shield.
func (k Keeper) CreatePool(ctx sdk.Context, creator sdk.AccAddress, shield sdk.Coins, serviceFees types.MixedCoins, sponsor string, sponsorAddr sdk.AccAddress, ProtectionPeriod time.Duration, description string) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !creator.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}

	if _, found := k.GetPoolBySponsor(ctx, sponsor); found {
		return types.Pool{}, types.ErrSponsorAlreadyExists
	}

	if !k.ValidatePoolDuration(ctx, ProtectionPeriod) {
		return types.Pool{}, types.ErrPoolLifeTooShort
	}

	// Check if shield is backed by admin's delegations.
	provider, found := k.GetProvider(ctx, admin)
	if !found {
		provider = k.addProvider(ctx, admin)
	}
	shieldAmt := shield.AmountOf(k.sk.BondDenom(ctx))
	provider.Collateral = provider.Collateral.Add(shieldAmt)
	if shieldAmt.GT(provider.Available) {
		return types.Pool{}, sdkerrors.Wrapf(types.ErrInsufficientStaking, "available %s, shield %s", provider.Available, shield)
	}
	provider.Available = provider.Available.Sub(shieldAmt)

	id := k.GetNextPoolID(ctx)
	pool := types.NewPool(id, description, sponsor, sponsorAddr, shieldAmt)

	// Transfer deposit to the Shield module account.
	if err := k.DepositNativeServiceFees(ctx, serviceFees.Native, creator); err != nil {
		return types.Pool{}, err
	}

	// Update service fees in the global pool.
	globalPool := k.GetGlobalPool(ctx)
	globalPool.ServiceFees = globalPool.ServiceFees.Add(types.MixedDecCoinsFromMixedCoins(serviceFees))
	k.SetGlobalPool(ctx, globalPool)

	// Make a pseudo-purchase for B.
	purchaseID := k.GetNextPurchaseID(ctx)
	protectionEndTime := ctx.BlockTime().Add(ProtectionPeriod)
	votingPeriod := k.gk.GetVotingParams(ctx).VotingPeriod * 2
	// FIXME +4 or +7+4?
	deletionTime := protectionEndTime.Add(votingPeriod)
	purchase := types.NewPurchase(purchaseID, protectionEndTime, "shield for sponsor", shieldAmt)

	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, id+1)
	k.SetProvider(ctx, admin, provider)
	k.SetCollateral(ctx, pool, admin, types.NewCollateral(pool, admin, shieldAmt))

	k.AddPurchase(ctx, id, sponsorAddr, purchase)
	k.InsertPurchaseQueue(ctx, types.NewPurchaseList(id, sponsorAddr, []types.Purchase{purchase}), deletionTime)
	k.SetNextPurchaseID(ctx, purchaseID+1)

	return pool, nil
}

// UpdatePool updates pool info.
func (k Keeper) UpdatePool(ctx sdk.Context, poolID uint64, description string, updater sdk.AccAddress) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}

	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.Pool{}, types.ErrNoPoolFound
	}
	if description != "" {
		pool.Description = description
	}
	k.SetPool(ctx, pool)
	return pool, nil
}

// PausePool sets an active pool to be inactive.
func (k Keeper) PausePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}
	pool, found := k.GetPool(ctx, id)
	if !found {
		return types.Pool{}, types.ErrNoPoolFound
	}
	if !pool.Active {
		return types.Pool{}, types.ErrPoolAlreadyPaused
	}
	pool.Active = false
	k.SetPool(ctx, pool)
	return pool, nil
}

// ResumePool sets an inactive pool to be active.
func (k Keeper) ResumePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}
	pool, found := k.GetPool(ctx, id)
	if !found {
		return types.Pool{}, types.ErrNoPoolFound
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

// ClosePool closes the pool.
func (k Keeper) ClosePool(ctx sdk.Context, pool types.Pool) {
	// TODO: make sure nothing else needs to be done
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolKey(pool.ID))
}

// IterateAllPools iterates over the all the stored pools and performs a callback function.
func (k Keeper) IterateAllPools(ctx sdk.Context, callback func(pool types.Pool) (stop bool)) {
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

// ValidatePoolDuration validates new pool duration to be valid.
func (k Keeper) ValidatePoolDuration(ctx sdk.Context, timeDuration time.Duration) bool {
	poolParams := k.GetPoolParams(ctx)
	minPoolDuration := poolParams.MinPoolLife
	return timeDuration >= minPoolDuration
}
