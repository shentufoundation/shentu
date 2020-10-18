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
	store.Set(types.GetPoolKey(pool.PoolID), bz)
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
func (k Keeper) CreatePool(ctx sdk.Context, creator sdk.AccAddress,
	shield sdk.Coins, deposit types.MixedCoins, sponsor string,
	sponsorAddr sdk.AccAddress, poolLifeTime time.Duration) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !creator.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}

	if _, found := k.GetPoolBySponsor(ctx, sponsor); found {
		return types.Pool{}, types.ErrSponsorAlreadyExists
	}

	if !k.ValidatePoolDuration(ctx, poolLifeTime) {
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
		return types.Pool{}, sdkerrors.Wrapf(types.ErrInsufficientStaking,
			"available %s, shield %s", provider.Available, shield)
	}
	provider.Available = provider.Available.Sub(shieldAmt)

	endTime := ctx.BlockTime().Add(poolLifeTime)
	id := k.GetNextPoolID(ctx)
	depositDec := types.MixedDecCoinsFromMixedCoins(deposit)

	pool := types.NewPool(shield, shieldAmt, depositDec, sponsor, sponsorAddr, endTime, id)

	// Transfer deposit to the Shield module account.
	if err := k.DepositNativePremium(ctx, deposit.Native, creator); err != nil {
		return types.Pool{}, err
	}

	// Make a pseudo-purchase for B.
	purchaseID := k.GetNextPurchaseID(ctx)
	expirationTime := pool.EndTime.Add(-k.gk.GetVotingParams(ctx).VotingPeriod * 2)
	purchase := types.NewPurchase(purchaseID, shield, ctx.BlockHeight(), expirationTime, expirationTime, expirationTime, "shield for sponsor")

	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, id+1)
	k.SetProvider(ctx, admin, provider)
	k.SetCollateral(ctx, pool, admin, types.NewCollateral(pool, admin, shieldAmt))

	k.AddPurchase(ctx, id, sponsorAddr, purchase)
	k.InsertPurchaseQueue(ctx, types.NewPurchaseList(id, sponsorAddr, []types.Purchase{purchase}), expirationTime)
	k.SetNextPurchaseID(ctx, purchaseID+1)

	return pool, nil
}

// UpdatePool updates pool info and sponsor's shield.
func (k Keeper) UpdatePool(ctx sdk.Context, updater sdk.AccAddress, shield sdk.Coins, deposit types.MixedCoins, id uint64, addTime time.Duration, description string) (types.Pool, error) {
	admin := k.GetAdmin(ctx)
	if !updater.Equals(admin) {
		return types.Pool{}, types.ErrNotShieldAdmin
	}

	// Check if shield is backed by admin's delegations.
	provider, found := k.GetProvider(ctx, admin)
	if !found {
		return types.Pool{}, types.ErrNoDelegationAmount
	}
	shieldAmt := shield.AmountOf(k.sk.BondDenom(ctx))
	provider.Collateral = provider.Collateral.Add(shieldAmt)
	if shieldAmt.GT(provider.Available) {
		return types.Pool{}, sdkerrors.Wrapf(types.ErrInsufficientStaking,
			"available %s, shield %s", provider.Available, shield)
	}
	provider.Available = provider.Available.Sub(shieldAmt)

	pool, found := k.GetPool(ctx, id)
	if !found {
		return types.Pool{}, types.ErrNoPoolFound
	}
	pool.EndTime = pool.EndTime.Add(addTime)
	pool.TotalCollateral = pool.TotalCollateral.Add(shieldAmt)
	poolCertiKCollateral, found := k.GetPoolCertiKCollateral(ctx, pool)
	if !found {
		poolCertiKCollateral = types.NewCollateral(pool, admin, sdk.ZeroInt())
	}
	poolCertiKCollateral.Amount = poolCertiKCollateral.Amount.Add(shieldAmt)
	pool.Shield = pool.Shield.Add(shield...)
	pool.Premium = pool.Premium.Add(types.MixedDecCoinsFromMixedCoins(deposit))
	if description != "" {
		pool.Description = description
	}

	// Transfer deposit and store.
	if err := k.DepositNativePremium(ctx, deposit.Native, admin); err != nil {
		return types.Pool{}, err
	}

	// Update sponsor purchase.
	sponsorPurchase, found := k.GetPurchaseList(ctx, id, pool.SponsorAddr)

	purchaseEndTime := pool.EndTime.Add(-k.gk.GetVotingParams(ctx).VotingPeriod * 2)
	if purchaseEndTime.Before(ctx.BlockTime()) {
		return types.Pool{}, types.ErrPoolLifeTooShort
	}
	// Assume there is only one purchase from sponsor address, and add in any if B's purchase expired.
	var purchase types.Purchase
	if !found {
		purchaseID := k.GetNextPurchaseID(ctx)
		purchase = types.NewPurchase(purchaseID, sdk.NewCoins(), ctx.BlockHeight(), purchaseEndTime, purchaseEndTime, purchaseEndTime, "shield for sponsor")
		k.SetNextPurchaseID(ctx, purchaseID+1)
	} else {
		purchase = sponsorPurchase.Entries[0]
		k.DequeuePurchase(ctx, sponsorPurchase, purchase.DeleteTime)
		purchase.DeleteTime = purchaseEndTime
		purchase.ClaimPeriodEndTime = purchaseEndTime
		purchase.ProtectionEndTime = purchaseEndTime
	}

	purchase.Shield = purchase.Shield.Add(shield...)
	newPurchaseList := types.NewPurchaseList(id, pool.SponsorAddr, []types.Purchase{purchase})

	k.SetCollateral(ctx, pool, k.GetAdmin(ctx), poolCertiKCollateral)
	k.SetPool(ctx, pool)
	k.SetProvider(ctx, admin, provider)
	k.SetPurchaseList(ctx, newPurchaseList)
	k.InsertPurchaseQueue(ctx, newPurchaseList, purchaseEndTime)
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

// PoolEnded returns if pool has reached ending time and block height.
func (k Keeper) PoolEnded(ctx sdk.Context, pool types.Pool) bool {
	return ctx.BlockTime().After(pool.EndTime)
}

// ClosePool closes the pool.
func (k Keeper) ClosePool(ctx sdk.Context, pool types.Pool) {
	// TODO: make sure nothing else needs to be done
	k.FreeCollaterals(ctx, pool)
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolKey(pool.PoolID))
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

// WithdrawFromPools withdraws coins from all pools to match total collateral to be less than or equal to total delegation.
func (k Keeper) WithdrawFromPools(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Int) {
	provider, _ := k.GetProvider(ctx, addr)
	withdrawAmtDec := sdk.NewDecFromInt(amount)
	withdrawableAmtDec := sdk.NewDecFromInt(provider.Collateral.Sub(provider.Withdrawing))
	proportion := withdrawAmtDec.Quo(withdrawableAmtDec)
	if withdrawAmtDec.GT(withdrawableAmtDec) {
		// FIXME this could happen. Set an error instead of panic.
		panic(types.ErrNotEnoughCollateral)
	}

	// Initiate proportional withdraws from all of the address's collaterals.
	addrCollaterals := k.GetProviderCollaterals(ctx, addr)
	remainingWithdraw := amount
	for i, collateral := range addrCollaterals {
		var withdrawAmt sdk.Int
		if i == len(addrCollaterals)-1 {
			withdrawAmt = remainingWithdraw
		} else {
			withdrawable := collateral.Amount.Sub(collateral.Withdrawing)
			withdrawAmtDec := sdk.NewDecFromInt(withdrawable).Mul(proportion)
			withdrawAmt = withdrawAmtDec.TruncateInt()
			if remainingWithdraw.LTE(withdrawAmt) {
				withdrawAmt = remainingWithdraw
			} else if remainingWithdraw.GT(withdrawAmt) && withdrawable.GT(withdrawAmt) {
				withdrawAmt = withdrawAmt.Add(sdk.OneInt())
			}
		}
		err := k.WithdrawCollateral(ctx, addr, collateral.PoolID, withdrawAmt)
		if err != nil {
			//TODO: address this error
			continue
		}
		remainingWithdraw = remainingWithdraw.Sub(withdrawAmt)
	}
}
