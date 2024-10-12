package keeper

import (
	"context"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// SetCustomParams sets parameters space for custom.
func (k Keeper) SetCustomParams(ctx context.Context, customAddParams typesv1.CustomParams) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&customAddParams)
	if err != nil {
		return err
	}
	return store.Set(types.CustomParamsKey, bz)
}

// GetCustomParams returns the current CustomParams from the global param store.
func (k Keeper) GetCustomParams(ctx context.Context) (customParams typesv1.CustomParams, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.CustomParamsKey)
	if err != nil {
		return customParams, err
	}
	if bz == nil {
		return customParams, nil
	}

	k.cdc.MustUnmarshal(bz, &customParams)
	return customParams, nil
}
