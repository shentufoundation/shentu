package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// GetCustomParams returns the current CustomParams from the global param store.
func (k Keeper) GetCustomParams(ctx sdk.Context) typesv1.CustomParams {
	var customAddParams typesv1.CustomParams
	k.paramSpace.Get(ctx, typesv1.ParamStoreKeyCustomParams, &customAddParams)
	return customAddParams
}

// SetCustomParams sets parameters space for custom.
func (k Keeper) SetCustomParams(ctx sdk.Context, customAddParams typesv1.CustomParams) {
	k.paramSpace.Set(ctx, typesv1.ParamStoreKeyCustomParams, &customAddParams)
}
