package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// SetCustomParams sets parameters space for custom.
func (k Keeper) SetCustomParams(ctx sdk.Context, customAddParams typesv1.CustomParams) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&customAddParams)
	if err != nil {
		return err
	}
	store.Set(types.CustomParamsKey, bz)
	return nil
}

// GetCustomParams returns the current CustomParams from the global param store.
func (k Keeper) GetCustomParams(ctx sdk.Context) (customParams typesv1.CustomParams) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CustomParamsKey)
	if bz == nil {
		return customParams
	}

	k.cdc.MustUnmarshal(bz, &customParams)
	return customParams
}
