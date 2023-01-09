package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

// GetCustomParams returns the current CustomParams from the global param store.
func (k Keeper) GetCustomParams(ctx sdk.Context) types.CustomParams {
	var customAddParams types.CustomParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyCustomParams, &customAddParams)
	return customAddParams
}

// SetCustomParams sets parameters space for custom.
func (k Keeper) SetCustomParams(ctx sdk.Context, customAddParams types.CustomParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyCustomParams, &customAddParams)
}
