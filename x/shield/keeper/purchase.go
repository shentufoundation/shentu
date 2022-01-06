package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

type pPPTriplet struct {
	poolID     uint64
	purchaseID uint64
	purchaser  sdk.AccAddress
}

// PurchaseShield purchases shield of a pool with standard fee rate.
func (k Keeper) PurchaseShield(ctx sdk.Context, poolID uint64, amount sdk.Coins, description string, purchaser sdk.AccAddress) (types.Purchase, error) {
	poolParams := k.GetPoolParams(ctx)
	if poolParams.MinShieldPurchase.IsAnyGT(amount) {
		return types.Purchase{}, types.ErrPurchaseTooSmall
	}
	bondDenom := k.BondDenom(ctx)
	if amount.AmountOf(bondDenom).Equal(sdk.ZeroInt()) {
		return types.Purchase{}, types.ErrInsufficientStaking
	}
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.Purchase{}, types.ErrNoPoolFound
	}
	if !pool.Active {
		return types.Purchase{}, types.ErrPoolInactive
	}
	if amount.Empty() {
		return types.Purchase{}, types.ErrNoShield
	}

	sp, err := k.AddStaking(ctx, poolID, purchaser, amount)
	return sp, err
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

func (k Keeper) GetPurchaserPurchases(ctx sdk.Context, address sdk.AccAddress) (res []types.Purchase) {
	purchases := k.GetAllPurchase(ctx)

	for _, p := range purchases {
		if p.Purchaser == address.String() {
			res = append(res, p)
		}
	}
	return
}

func (k Keeper) GetPoolPurchases(ctx sdk.Context, poolID uint64) (res []types.Purchase) {
	purchases := k.GetAllPurchase(ctx)

	for _, p := range purchases {
		if p.PoolId == poolID {
			res = append(res, p)
		}
	}
	return
}
