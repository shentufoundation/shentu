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

// DequeuePurchase dequeues a purchase from the purchase queue
func (k Keeper) DequeuePurchase(ctx sdk.Context, purchaseList types.PurchaseList, endTime time.Time) {
	timeslice := k.GetPurchaseQueueTimeSlice(ctx, endTime)
	for i, poolPurchaser := range timeslice {
		if (purchaseList.PoolID == poolPurchaser.PoolID) && purchaseList.Purchaser.Equals(poolPurchaser.Purchaser) {
			if len(timeslice) > 1 {
				timeslice = append(timeslice[:i], timeslice[i+1:]...)
				k.SetPurchaseQueueTimeSlice(ctx, endTime, timeslice)
				return
			}
			ctx.KVStore(k.storeKey).Delete(types.GetPurchaseCompletionTimeKey(endTime))
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
	totalServiceFeesLeft := k.GetServiceFeesLeft(ctx)
	totalServiceFeesLeft = totalServiceFeesLeft.Add(types.MixedDecCoins{Native: sdk.NewDecCoinsFromCoins(serviceFees...)})
	k.SetServiceFeesLeft(ctx, totalServiceFeesLeft)

	// Update global pool and project pool's shield.
	totalShield = totalShield.Add(shieldAmt)
	pool.Shield = pool.Shield.Add(shieldAmt)
	k.SetTotalShield(ctx, totalShield)
	k.SetPool(ctx, pool)

	// Set a new purchase.
	protectionEndTime := ctx.BlockTime().Add(poolParams.ProtectionPeriod)
	purchaseID := k.GetNextPurchaseID(ctx)
	purchase := types.NewPurchase(purchaseID, protectionEndTime, description, shieldAmt, types.MixedDecCoins{Native: sdk.NewDecCoinsFromCoins(serviceFees...)})
	purchaseList := k.AddPurchase(ctx, poolID, purchaser, purchase)
	k.InsertPurchaseQueue(ctx, purchaseList, protectionEndTime.Add(k.GetPurchaseDeletionPeriod(ctx)))
	k.SetNextPurchaseID(ctx, purchaseID+1)

	return purchase, nil
}

// PurchaseShield purchases shield of a pool with standard fee rate.
func (k Keeper) PurchaseShield(ctx sdk.Context, poolID uint64, shield sdk.Coins, description string, purchaser sdk.AccAddress) (types.Purchase, error) {
	bondDenom := k.BondDenom(ctx)
	serviceFees := sdk.NewCoins(sdk.NewCoin(bondDenom, shield.AmountOf(bondDenom).ToDec().Mul(k.GetPoolParams(ctx).ShieldFeesRate).TruncateInt()))
	return k.purchaseShield(ctx, poolID, shield, description, purchaser, serviceFees)
}

// IterateAllPurchases iterates over the all the stored purchases and performs a callback function.
func (k Keeper) IterateAllPurchases(ctx sdk.Context, callback func(purchase types.Purchase) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseListKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var purchase types.Purchase
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &purchase)

		if callback(purchase) {
			break
		}
	}
}

// RemoveExpiredPurchasesAndDistributeFees removes expired purchases and distributes fees for current block.
func (k Keeper) RemoveExpiredPurchasesAndDistributeFees(ctx sdk.Context) {
	// Remove expired services and get service fees they should pay.
	serviceFees := k.removeExpiredPurchases(ctx)

	// Add service fees for this block from unexpired purchases.
	serviceFees = serviceFees.Add(k.getServiceFeesForBlock(ctx))

	// Limit service fees by service fees left.
	serviceFeesLeft := k.GetServiceFeesLeft(ctx)
	bondDenom := k.BondDenom(ctx)
	if serviceFeesLeft.Native.AmountOf(bondDenom).LT(serviceFees.Native.AmountOf(bondDenom)) {
		serviceFees.Native = serviceFeesLeft.Native
	}

	// Distribute and update service fees.
	k.distributeFees(ctx, serviceFees)
	k.updateServiceFees(ctx)
}

// removeExpiredPurchases removes expired purchases and return remaining fees
func (k Keeper) removeExpiredPurchases(ctx sdk.Context) types.MixedDecCoins {
	store := ctx.KVStore(k.storeKey)
	totalServiceFees := k.GetServiceFees(ctx)
	totalShield := k.GetTotalShield(ctx)
	deletionPeriod := k.GetPurchaseDeletionPeriod(ctx)
	previousBlockTime := ctx.WithBlockHeight(ctx.BlockHeight() - 1).BlockTime()

	serviceFees := types.InitMixedDecCoins()
	iterator := k.PurchaseQueueIterator(ctx, ctx.BlockTime())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var timeslice []types.PoolPurchaser
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &timeslice)
		for _, poolPurchaser := range timeslice {
			purchaseList, _ := k.GetPurchaseList(ctx, poolPurchaser.PoolID, poolPurchaser.Purchaser)

			for i := 0; i < len(purchaseList.Entries); {
				entry := purchaseList.Entries[i]
				// DeletionTime = ProtectionEndTime - ProtectionPeriod + ClaimPeriod + VotingPeriod
				// If DeletionTime > currentBlockTime, skip.
				if entry.ProtectionEndTime.Add(deletionPeriod).After(ctx.BlockTime()) {
					i++
					continue
				}
				// If previousBlockTime < ProtectionTime <= DeletionTime <= currentBlockTime,
				// calculate remaining service fees to be distributed.
				if entry.ProtectionEndTime.After(previousBlockTime) {
					// Add purchaseServiceFees * (purchaseProtectionEndTime - previousBlockTime) / protectionPeriod.
					serviceFees = serviceFees.Add(entry.ServiceFees.MulDec(
						sdk.NewDec(int64(entry.ProtectionEndTime.Sub(previousBlockTime).Seconds())).Quo(
							sdk.NewDec(int64(k.GetPoolParams(ctx).ProtectionPeriod.Seconds())))))
					// Remove purchaseServiceFees from total service fees.
					totalServiceFees = totalServiceFees.Sub(entry.ServiceFees)
				}
				// If DeletionTime <= currentBlockTime, remove purchase and update shield.
				purchaseList.Entries = append(purchaseList.Entries[:i], purchaseList.Entries[i+1:]...)
				pool, found := k.GetPool(ctx, purchaseList.PoolID)
				if found {
					totalShield = totalShield.Sub(entry.Shield)
					pool.Shield = pool.Shield.Sub(entry.Shield)
					k.SetPool(ctx, pool)
				}
			}
			if len(purchaseList.Entries) == 0 {
				_ = k.DeletePurchaseList(ctx, purchaseList.PoolID, purchaseList.Purchaser)
			} else {
				k.SetPurchaseList(ctx, purchaseList)
			}
		}
		store.Delete(iterator.Key())
	}
	k.SetServiceFees(ctx, totalServiceFees)
	k.SetTotalShield(ctx, totalShield)
	return serviceFees
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

// GetAllPurchaseLists retrieves all purchases.
func (k Keeper) GetAllPurchaseLists(ctx sdk.Context) (purchases []types.PurchaseList) {
	k.IteratePurchaseLists(ctx, func(purchase types.PurchaseList) bool {
		purchases = append(purchases, purchase)
		return false
	})
	return
}

func (k Keeper) InsertPurchaseQueue(ctx sdk.Context, purchaseList types.PurchaseList, endTime time.Time) {
	timeSlice := k.GetPurchaseQueueTimeSlice(ctx, endTime)

	poolPurchaser := types.PoolPurchaser{PoolID: purchaseList.PoolID, Purchaser: purchaseList.Purchaser}
	if len(timeSlice) == 0 {
		k.SetPurchaseQueueTimeSlice(ctx, endTime, []types.PoolPurchaser{poolPurchaser})
		return
	}
	timeSlice = append(timeSlice, poolPurchaser)
	k.SetPurchaseQueueTimeSlice(ctx, endTime, timeSlice)
}

// GetPurchaseQueueTimeSlice gets a specific purchase queue timeslice,
// which is a slice of purchases corresponding to a given time.
func (k Keeper) GetPurchaseQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []types.PoolPurchaser {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPurchaseCompletionTimeKey(timestamp))
	if bz == nil {
		return []types.PoolPurchaser{}
	}
	var ppPairs []types.PoolPurchaser
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &ppPairs)
	return ppPairs
}

func (k Keeper) SetPurchaseQueueTimeSlice(ctx sdk.Context, timestamp time.Time, ppPairs []types.PoolPurchaser) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(ppPairs)
	store.Set(types.GetPurchaseCompletionTimeKey(timestamp), bz)
}

// PurchaseQueueIterator returns all the purchase queue timeslices from time 0 until endTime
func (k Keeper) PurchaseQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.PurchaseQueueKey,
		sdk.InclusiveEndBytes(types.GetPurchaseCompletionTimeKey(endTime)))
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
