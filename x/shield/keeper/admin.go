package keeper

import (
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetAdmin sets the Shield admin account address.
func (k Keeper) SetAdmin(ctx sdk.Context, admin sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := admin
	store.Set(types.GetShieldAdminKey(), bz)
}

// GetAdmin gets the Shield admin account address.
func (k Keeper) GetAdmin(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.GetShieldAdminKey())
}
