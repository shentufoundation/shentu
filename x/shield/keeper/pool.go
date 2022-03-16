package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

func (k Keeper) SetTotalCollateral(ctx sdk.Context, totalCollateral sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&sdk.IntProto{Int: totalCollateral})
	store.Set(types.GetTotalCollateralKey(), bz)
}

func (k Keeper) GetTotalCollateral(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTotalCollateralKey())
	if bz == nil {
		panic("total collateral is not found")
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &ip)
	return ip.Int
}

func (k Keeper) SetTotalWithdrawing(ctx sdk.Context, totalWithdrawing sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&sdk.IntProto{Int: totalWithdrawing})
	store.Set(types.GetTotalWithdrawingKey(), bz)
}

func (k Keeper) GetTotalWithdrawing(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTotalWithdrawingKey())
	if bz == nil {
		panic("total withdrawing is not found")
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &ip)
	return ip.Int
}

func (k Keeper) SetTotalShield(ctx sdk.Context, totalShield sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&sdk.IntProto{Int: totalShield})
	store.Set(types.GetTotalShieldKey(), bz)
}

func (k Keeper) GetTotalShield(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTotalShieldKey())
	if bz == nil {
		panic("total shield is not found")
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &ip)
	return ip.Int
}

func (k Keeper) SetTotalClaimed(ctx sdk.Context, totalClaimed sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&sdk.IntProto{Int: totalClaimed})
	store.Set(types.GetTotalClaimedKey(), bz)
}

func (k Keeper) GetTotalClaimed(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTotalClaimedKey())
	if bz == nil {
		panic("total shield is not found")
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &ip)
	return ip.Int
}

func (k Keeper) SetServiceFees(ctx sdk.Context, serviceFees sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&v1beta1.Fees{Amount: serviceFees})
	store.Set(types.GetServiceFeesKey(), bz)
}

func (k Keeper) GetServiceFees(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetServiceFeesKey())
	if bz == nil {
		return sdk.DecCoins{}
	}
	var serviceFees v1beta1.Fees
	k.cdc.MustUnmarshalLengthPrefixed(bz, &serviceFees)
	return serviceFees.Amount
}

func (k Keeper) DeleteServiceFees(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetServiceFeesKey())
}

// SetPool sets data of a pool in kv-store.
func (k Keeper) SetPool(ctx sdk.Context, pool v1beta1.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(types.GetPoolKey(pool.Id), bz)
}

// GetPool gets data of a pool given pool ID.
func (k Keeper) GetPool(ctx sdk.Context, id uint64) (v1beta1.Pool, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolKey(id))
	if bz == nil {
		return v1beta1.Pool{}, false
	}
	var pool v1beta1.Pool
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}

// CreatePool creates a pool and sponsor's shield.
func (k Keeper) CreatePool(ctx sdk.Context, msg v1beta1.MsgCreatePool) (uint64, error) {
	creator, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return 0, err
	}
	sponsorAddr, err := sdk.AccAddressFromBech32(msg.SponsorAddr)
	if err != nil {
		return 0, err
	}

	admin := k.GetAdmin(ctx)
	if !creator.Equals(admin) {
		return 0, types.ErrNotShieldAdmin
	}
	if _, found := k.GetPoolsBySponsor(ctx, msg.SponsorAddr); found {
		return 0, types.ErrSponsorAlreadyExists
	}

	// Set the new project pool.
	poolID := k.GetNextPoolID(ctx)

	pool := v1beta1.NewPool(poolID, msg.Description, sponsorAddr, sdk.ZeroInt(), msg.ShieldRate, msg.ShieldLimit)

	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, poolID+1)
	return poolID, nil
}

// UpdatePool updates pool info and shield for B.
func (k Keeper) UpdatePool(ctx sdk.Context, msg v1beta1.MsgUpdatePool) (v1beta1.Pool, error) {
	updater, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return v1beta1.Pool{}, err
	}

	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return v1beta1.Pool{}, types.ErrNotShieldAdmin
	}

	// Update pool info.
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return v1beta1.Pool{}, types.ErrNoPoolFound
	}
	if msg.Description != "" {
		pool.Description = msg.Description
	}
	if !msg.ShieldRate.IsZero() {
		pool.ShieldRate = msg.ShieldRate
	}

	if !msg.ShieldLimit.IsZero() {
		pool.ShieldLimit = msg.ShieldLimit
	}
	pool.Active = msg.Active
	k.SetPool(ctx, pool)
	return pool, nil
}

// PausePool sets an active pool to be inactive.
func (k Keeper) PausePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (v1beta1.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return v1beta1.Pool{}, types.ErrNotShieldAdmin
	}
	pool, found := k.GetPool(ctx, id)
	if !found {
		return v1beta1.Pool{}, types.ErrNoPoolFound
	}
	if !pool.Active {
		return v1beta1.Pool{}, types.ErrPoolAlreadyPaused
	}
	pool.Active = false
	k.SetPool(ctx, pool)
	return pool, nil
}

// ResumePool sets an inactive pool to be active.
func (k Keeper) ResumePool(ctx sdk.Context, updater sdk.AccAddress, id uint64) (v1beta1.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return v1beta1.Pool{}, types.ErrNotShieldAdmin
	}
	pool, found := k.GetPool(ctx, id)
	if !found {
		return v1beta1.Pool{}, types.ErrNoPoolFound
	}
	if pool.Active {
		return v1beta1.Pool{}, types.ErrPoolAlreadyActive
	}
	pool.Active = true
	k.SetPool(ctx, pool)
	return pool, nil
}

// GetAllPools retrieves all pools in the store.
func (k Keeper) GetAllPools(ctx sdk.Context) (pools []v1beta1.Pool) {
	k.IterateAllPools(ctx, func(pool v1beta1.Pool) bool {
		pools = append(pools, pool)
		return false
	})
	return pools
}

// IterateAllPools iterates over the all the stored pools and performs a callback function.
func (k Keeper) IterateAllPools(ctx sdk.Context, callback func(pool v1beta1.Pool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PoolKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pool v1beta1.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		if callback(pool) {
			break
		}
	}
}

// UpdateSponsor updates the sponsor information of a given pool.
func (k Keeper) UpdateSponsor(ctx sdk.Context, poolID uint64, newSponsor string, newSponsorAddr, updater sdk.AccAddress) (v1beta1.Pool, error) {
	// Check admin status of the updater.
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return v1beta1.Pool{}, types.ErrNotShieldAdmin
	}

	// Retrieve the pool and update its sponsor information.
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return v1beta1.Pool{}, types.ErrNoPoolFound
	}
	pool.SponsorAddr = newSponsorAddr.String()
	k.SetPool(ctx, pool)

	return pool, nil
}
