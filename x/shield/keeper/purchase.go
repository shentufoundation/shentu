package keeper

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// SetPurchaseList sets a purchase list.
func (k Keeper) SetPurchaseList(ctx sdk.Context, purchaseList types.PurchaseList) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(purchaseList)
	store.Set(types.GetPurchaseListKey(purchaseList.PoolID, purchaseList.Purchaser), bz)
}

// AddPurchase sets a purchase of shield.
func (k Keeper) AddPurchase(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchase types.Purchase) types.PurchaseList {
	purchaseList, found := k.GetPurchaseList(ctx, poolID, purchaser)
	if !found {
		purchaseList = types.NewPurchaseList(poolID, purchaser, []types.Purchase{purchase})
	} else {
		purchaseList.Entries = append(purchaseList.Entries, purchase)
	}
	k.SetPurchaseList(ctx, purchaseList)
	return purchaseList
}

// GetPurchaseList gets a purchase from store by txhash.
func (k Keeper) GetPurchaseList(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (types.PurchaseList, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPurchaseListKey(poolID, purchaser))
	if bz != nil {
		var purchase types.PurchaseList
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &purchase)
		return purchase, true
	}
	return types.PurchaseList{}, false
}

// GetPurchase gets a purchase out of a purchase list
func GetPurchase(purchaseList types.PurchaseList, purchaseID uint64) (types.Purchase, bool) {
	for _, entry := range purchaseList.Entries {
		if entry.PurchaseID == purchaseID {
			return entry, true
		}
	}
	return types.Purchase{}, false
}

// DeletePurchaseList deletes a purchase of shield.
func (k Keeper) DeletePurchaseList(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	_, found := k.GetPurchaseList(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	store.Delete(types.GetPurchaseListKey(poolID, purchaser))
	return nil
}

// DequeuePurchase dequeues a purchase from the purchase queue.
func (k Keeper) DequeuePurchase(ctx sdk.Context, purchaseList types.PurchaseList, endTime time.Time) {
	timeslice := k.GetExpiringPurchaseQueueTimeSlice(ctx, endTime)
	for i, poolPurchaser := range timeslice {
		if (purchaseList.PoolID == poolPurchaser.PoolID) && purchaseList.Purchaser.Equals(poolPurchaser.Purchaser) {
			if len(timeslice) > 1 {
				timeslice = append(timeslice[:i], timeslice[i+1:]...)
				k.SetExpiringPurchaseQueueTimeSlice(ctx, endTime, timeslice)
				return
			}
			ctx.KVStore(k.storeKey).Delete(types.GetPurchaseExpirationTimeKey(endTime))
			return
		}
	}
}

// PurchaseShield purchases shield of a pool.
func (k Keeper) purchaseShield(ctx sdk.Context, poolID uint64, shield sdk.Coins, description string, purchaser sdk.AccAddress, serviceFees sdk.Coins) (types.Purchase, error) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.Purchase{}, types.ErrNoPoolFound
	}
	if !pool.Active {
		return types.Purchase{}, types.ErrPoolInactive
	}

	// Check available collaterals.
	shieldAmt := shield.AmountOf(k.sk.BondDenom(ctx))
	totalCollateral := k.GetTotalCollateral(ctx)
	totalWithdrawing := k.GetTotalWithdrawing(ctx)
	totalShield := k.GetTotalShield(ctx)
	if totalShield.Add(shieldAmt).GT(totalCollateral.Sub(totalWithdrawing)) {
		return types.Purchase{}, types.ErrNotEnoughCollateral
	}

	// Check pool shield limit.
	poolParams := k.GetPoolParams(ctx)
	maxShield := sdk.MinInt(pool.ShieldLimit, totalCollateral.Sub(totalWithdrawing).ToDec().Mul(poolParams.PoolShieldLimit).TruncateInt())
	if shieldAmt.Add(pool.Shield).GT(maxShield) {
		return types.Purchase{}, types.ErrPoolShieldExceedsLimit
	}

	// Send service fees to the shield module account and update service fees.
	if err := k.DepositNativeServiceFees(ctx, serviceFees, purchaser); err != nil {
		return types.Purchase{}, err
	}
	totalServiceFees := k.GetServiceFees(ctx)
	totalServiceFees = totalServiceFees.Add(types.MixedDecCoins{Native: sdk.NewDecCoinsFromCoins(serviceFees...)})
	k.SetServiceFees(ctx, totalServiceFees)
	totalRemainingServiceFees := k.GetRemainingServiceFees(ctx)
	totalRemainingServiceFees = totalRemainingServiceFees.Add(types.MixedDecCoins{Native: sdk.NewDecCoinsFromCoins(serviceFees...)})
	k.SetRemainingServiceFees(ctx, totalRemainingServiceFees)

	// Update global pool and project pool's shield.
	totalShield = totalShield.Add(shieldAmt)
	pool.Shield = pool.Shield.Add(shieldAmt)
	k.SetTotalShield(ctx, totalShield)
	k.SetPool(ctx, pool)

	// Set a new purchase.
	protectionEndTime := ctx.BlockTime().Add(poolParams.ProtectionPeriod)
	purchaseID := k.GetNextPurchaseID(ctx)
	purchase := types.NewPurchase(purchaseID, protectionEndTime, protectionEndTime, description, shieldAmt, types.MixedDecCoins{Native: sdk.NewDecCoinsFromCoins(serviceFees...)})
	purchaseList := k.AddPurchase(ctx, poolID, purchaser, purchase)
	k.InsertExpiringPurchaseQueue(ctx, purchaseList, protectionEndTime)
	k.SetNextPurchaseID(ctx, purchaseID+1)

	lastUpdateTime, found := k.GetLastUpdateTime(ctx)
	if !found || lastUpdateTime.IsZero() {
		k.SetLastUpdateTime(ctx, ctx.BlockTime())
	}

	return purchase, nil
}

// PurchaseShield purchases shield of a pool with standard fee rate.
func (k Keeper) PurchaseShield(ctx sdk.Context, poolID uint64, shield sdk.Coins, description string, purchaser sdk.AccAddress) (types.Purchase, error) {
	poolParams := k.GetPoolParams(ctx)
	if poolParams.MinShieldPurchase.IsAnyGT(shield) {
		return types.Purchase{}, types.ErrPurchaseTooSmall
	}
	bondDenom := k.BondDenom(ctx)
	serviceFees := sdk.NewCoins(sdk.NewCoin(bondDenom, shield.AmountOf(bondDenom).ToDec().Mul(k.GetPoolParams(ctx).ShieldFeesRate).TruncateInt()))
	return k.purchaseShield(ctx, poolID, shield, description, purchaser, serviceFees)
}

// RemoveExpiredPurchasesAndDistributeFees removes expired purchases and distributes fees for current block.
func (k Keeper) RemoveExpiredPurchasesAndDistributeFees(ctx sdk.Context) {
	lastUpdateTime, found := k.GetLastUpdateTime(ctx)
	if !found || lastUpdateTime.IsZero() {
		// Last update time will be set when a purchase is made.
		return
	}

	store := ctx.KVStore(k.storeKey)
	totalServiceFees := k.GetServiceFees(ctx)
	totalShield := k.GetTotalShield(ctx)
	serviceFees := types.InitMixedDecCoins()

	// Check all purchases whose protection end time is before current block time.
	// 1) Update service fees for purchases whose protection end time is before current block time.
	// 2) Remove purchases whose deletion time is before current block time.
	iterator := k.ExpiringPurchaseQueueIterator(ctx, ctx.BlockTime())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var timeslice []types.PoolPurchaser
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &timeslice)
		for _, poolPurchaser := range timeslice {
			purchaseList, _ := k.GetPurchaseList(ctx, poolPurchaser.PoolID, poolPurchaser.Purchaser)
			for i := 0; i < len(purchaseList.Entries); i++ {
				entry := purchaseList.Entries[i]

				// If purchaseProtectionEndTime > previousBlockTime, update service fees.
				// Otherwise services fees were updated in the last block.
				if entry.ProtectionEndTime.After(lastUpdateTime) && entry.ServiceFees.Native.IsAllPositive() {
					// Add purchaseServiceFees * (purchaseProtectionEndTime - previousBlockTime) / protectionPeriod.
					serviceFees = serviceFees.Add(entry.ServiceFees.MulDec(
						sdk.NewDec(entry.ProtectionEndTime.Sub(lastUpdateTime).Nanoseconds()).Quo(
							sdk.NewDec(k.GetPoolParams(ctx).ProtectionPeriod.Nanoseconds()))))
					// Remove purchaseServiceFees from total service fees.
					totalServiceFees = totalServiceFees.Sub(entry.ServiceFees)
					// Set purchaseServiceFees to zero because it can be reached again.
					purchaseList.Entries[i].ServiceFees = types.InitMixedDecCoins()
				}

				// If purchaseDeletionTime < currentBlockTime, remove the purchase.
				if entry.DeletionTime.Before(ctx.BlockTime()) {
					// If purchaseProtectionEndTime > previousBlockTime, calculate and set service fees before removing the purchase.
					purchaseList.Entries = append(purchaseList.Entries[:i], purchaseList.Entries[i+1:]...)
					// Update pool shield and total shield.
					pool, found := k.GetPool(ctx, purchaseList.PoolID)
					if !found {
						panic("cannot find the pool for an expired purchase")
					}
					totalShield = totalShield.Sub(entry.Shield)
					pool.Shield = pool.Shield.Sub(entry.Shield)
					k.SetPool(ctx, pool)
					// Minus one because the current entry is deleted.
					i--
				}
			}
			if len(purchaseList.Entries) == 0 {
				_ = k.DeletePurchaseList(ctx, purchaseList.PoolID, purchaseList.Purchaser)
			} else {
				k.SetPurchaseList(ctx, purchaseList)
			}
		}
		// TODO: For phase I only. Need to modify the logic here after claims are enabled.
		store.Delete(iterator.Key())
	}
	k.SetServiceFees(ctx, totalServiceFees)
	k.SetTotalShield(ctx, totalShield)

	// Add service fees for this block from unexpired purchases.
	// totalServiceFees * (currentBlockTime - previousBlockTime) / protectionPeriodTime
	serviceFees = serviceFees.Add(totalServiceFees.MulDec(
		sdk.NewDec(ctx.BlockTime().Sub(lastUpdateTime).Nanoseconds())).QuoDec(
		sdk.NewDec(k.GetPoolParams(ctx).ProtectionPeriod.Nanoseconds())))

	// Limit service fees by remaining service fees.
	remainingServiceFees := k.GetRemainingServiceFees(ctx)
	bondDenom := k.BondDenom(ctx)
	if remainingServiceFees.Native.AmountOf(bondDenom).LT(serviceFees.Native.AmountOf(bondDenom)) {
		serviceFees.Native = remainingServiceFees.Native
	}

	// Distribute service fees.
	totalCollateral := k.GetTotalCollateral(ctx)
	providers := k.GetAllProviders(ctx)
	for _, provider := range providers {
		// fees * providerCollateral / totalCollateral
		nativeFees := serviceFees.Native.MulDec(sdk.NewDecFromInt(provider.Collateral).QuoInt(totalCollateral))
		if nativeFees.AmountOf(bondDenom).GT(remainingServiceFees.Native.AmountOf(bondDenom)) {
			nativeFees = remainingServiceFees.Native
		}
		provider.Rewards = provider.Rewards.Add(types.MixedDecCoins{Native: nativeFees})
		k.SetProvider(ctx, provider.Address, provider)

		remainingServiceFees.Native = remainingServiceFees.Native.Sub(nativeFees)
	}
	k.SetRemainingServiceFees(ctx, remainingServiceFees)
	k.SetLastUpdateTime(ctx, ctx.BlockTime())
}

// GetPurchaserPurchases returns all purchases by a given purchaser.
func (k Keeper) GetPurchaserPurchases(ctx sdk.Context, address sdk.AccAddress) (res []types.PurchaseList) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		pList, found := k.GetPurchaseList(ctx, pool.ID, address)
		if !found {
			continue
		}
		res = append(res, pList)
	}
	return
}

// GetPoolPurchaseLists returns all purchases in a given pool.
func (k Keeper) GetPoolPurchaseLists(ctx sdk.Context, poolID uint64) (purchases []types.PurchaseList) {
	k.IteratePoolPurchaseLists(ctx, poolID, func(purchaseList types.PurchaseList) bool {
		if purchaseList.PoolID == poolID {
			purchases = append(purchases, purchaseList)
		}
		return false
	})
	return purchases
}

// IteratePurchaseLists iterates through purchase lists in a pool
func (k Keeper) IteratePurchaseLists(ctx sdk.Context, callback func(purchase types.PurchaseList) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseListKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchaseList types.PurchaseList
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchaseList)

		if callback(purchaseList) {
			break
		}
	}
}

// IteratePoolPurchaseLists iterates through purchases in a pool
func (k Keeper) IteratePoolPurchaseLists(ctx sdk.Context, poolID uint64, callback func(purchaseList types.PurchaseList) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, poolID)
	iterator := sdk.KVStorePrefixIterator(store, append(types.PurchaseListKey, bz...))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchaseList types.PurchaseList
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchaseList)

		if callback(purchaseList) {
			break
		}
	}
}

// GetAllPurchaseLists retrieves all purchase lists.
func (k Keeper) GetAllPurchaseLists(ctx sdk.Context) (purchases []types.PurchaseList) {
	k.IteratePurchaseLists(ctx, func(purchase types.PurchaseList) bool {
		purchases = append(purchases, purchase)
		return false
	})
	return
}

// InsertExpiringPurchaseQueue inserts a purchase into the expired purchase queue.
func (k Keeper) InsertExpiringPurchaseQueue(ctx sdk.Context, purchaseList types.PurchaseList, endTime time.Time) {
	timeSlice := k.GetExpiringPurchaseQueueTimeSlice(ctx, endTime)

	poolPurchaser := types.PoolPurchaser{PoolID: purchaseList.PoolID, Purchaser: purchaseList.Purchaser}
	if len(timeSlice) == 0 {
		k.SetExpiringPurchaseQueueTimeSlice(ctx, endTime, []types.PoolPurchaser{poolPurchaser})
		return
	}
	timeSlice = append(timeSlice, poolPurchaser)
	k.SetExpiringPurchaseQueueTimeSlice(ctx, endTime, timeSlice)
}

// GetExpiringPurchaseQueueTimeSlice gets a specific purchase queue timeslice,
// which is a slice of purchases corresponding to a given time.
func (k Keeper) GetExpiringPurchaseQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []types.PoolPurchaser {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPurchaseExpirationTimeKey(timestamp))
	if bz == nil {
		return []types.PoolPurchaser{}
	}
	var ppPairs []types.PoolPurchaser
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &ppPairs)
	return ppPairs
}

// SetExpiringPurchaseQueueTimeSlice sets a time slice for a purchase expiring at give time.
func (k Keeper) SetExpiringPurchaseQueueTimeSlice(ctx sdk.Context, timestamp time.Time, ppPairs []types.PoolPurchaser) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(ppPairs)
	store.Set(types.GetPurchaseExpirationTimeKey(timestamp), bz)
}

// ExpiringPurchaseQueueIterator returns a iterator of purchases expiring before endTime
func (k Keeper) ExpiringPurchaseQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.PurchaseQueueKey,
		sdk.InclusiveEndBytes(types.GetPurchaseExpirationTimeKey(endTime)))
}

// SetNextPurchaseID sets the latest pool ID to store.
func (k Keeper) SetNextPurchaseID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.GetNextPurchaseIDKey(), bz)
}

// GetNextPurchaseID gets the latest pool ID from store.
func (k Keeper) GetNextPurchaseID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.GetNextPurchaseIDKey())
	return binary.LittleEndian.Uint64(opBz)
}

// GetAllPurchases retrieves all purchases.
func (k Keeper) GetAllPurchases(ctx sdk.Context) (purchases []types.Purchase) {
	k.IteratePurchaseListEntries(ctx, func(purchase types.Purchase) bool {
		purchases = append(purchases, purchase)
		return false
	})
	return
}

// IteratePurchaseListEntries iterates through entries of
// all purchase lists.
func (k Keeper) IteratePurchaseListEntries(ctx sdk.Context, callback func(purchase types.Purchase) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseListKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchaseList types.PurchaseList
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchaseList)

		for _, entry := range purchaseList.Entries {
			if callback(entry) {
				break
			}
		}
	}
}
