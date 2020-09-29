package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// SetAdmin sets the Shield Admin account address.
func (k Keeper) SetAdmin(ctx sdk.Context, admin sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := admin
	store.Set(types.GetShieldOperatorKey(), bz)
}

// GetAdmin gets the Shield Admin account address.
func (k Keeper) GetAdmin(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.GetShieldOperatorKey())
}
