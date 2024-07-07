package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

func (k Keeper) GetBlockServiceFees(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBlockServiceFeesKey())
	if bz == nil {
		return sdk.DecCoins{}
	}
	var blockServiceFees types.Fees
	k.cdc.MustUnmarshalLengthPrefixed(bz, &blockServiceFees)
	return blockServiceFees.Fees
}

func (k Keeper) SetRemainingServiceFees(ctx sdk.Context, fees sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	serviceFee := types.Fees{
		Fees: fees,
	}
	bz := k.cdc.MustMarshalLengthPrefixed(&serviceFee)
	store.Set(types.GetRemainingServiceFeesKey(), bz)
}

func (k Keeper) GetRemainingServiceFees(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRemainingServiceFeesKey())
	if bz == nil {
		panic("remaining service fees are not found")
	}
	var remainingServiceFees types.Fees
	k.cdc.MustUnmarshalLengthPrefixed(bz, &remainingServiceFees)
	return remainingServiceFees.Fees
}
