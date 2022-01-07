package keeper

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

type pPPTriplet struct {
	poolID     uint64
	purchaseID uint64
	purchaser  sdk.AccAddress
}

// PurchaseShield purchases shield of a pool.
func (k Keeper) purchaseShield(ctx sdk.Context, poolID uint64, shield sdk.Coins, description string, purchaser sdk.AccAddress, serviceFees sdk.Coins, stakingCoins sdk.Coins) (types.Purchase, error) {

}

// PurchaseShield purchases shield of a pool with standard fee rate.
func (k Keeper) PurchaseShield(ctx sdk.Context, poolID uint64, shield sdk.Coins, description string, purchaser sdk.AccAddress, staking bool) (types.StakingPurchase, error) {
	poolParams := k.GetPoolParams(ctx)
	if poolParams.MinShieldPurchase.IsAnyGT(shield) {
		return types.StakingPurchase{}, types.ErrPurchaseTooSmall
	}
	bondDenom := k.BondDenom(ctx)
	serviceFees := sdk.NewCoins()
	stakingCoins := sdk.NewCoins()

	// TODO: consider direct purchase
	// stake to the staking purchase pool
	stakingAmt := k.GetShieldStakingRate(ctx).MulInt(shield.AmountOf(bondDenom)).TruncateInt()
	stakingCoins = sdk.NewCoins(sdk.NewCoin(bondDenom, stakingAmt))
	return k.purchaseShield(ctx, poolID, shield, description, purchaser, serviceFees, stakingCoins)
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
