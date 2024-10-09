package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

func (k Keeper) GetBlockServiceFees(ctx sdk.Context) (sdk.DecCoins, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GetBlockServiceFeesKey())
	if err != nil {
		return nil, err
	}
	var blockServiceFees types.Fees
	k.cdc.MustUnmarshalLengthPrefixed(bz, &blockServiceFees)
	return blockServiceFees.Fees, nil
}

func (k Keeper) SetRemainingServiceFees(ctx sdk.Context, fees sdk.DecCoins) error {
	store := k.storeService.OpenKVStore(ctx)
	serviceFee := types.Fees{
		Fees: fees,
	}
	bz := k.cdc.MustMarshalLengthPrefixed(&serviceFee)
	return store.Set(types.GetRemainingServiceFeesKey(), bz)
}

func (k Keeper) GetRemainingServiceFees(ctx sdk.Context) (sdk.DecCoins, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GetRemainingServiceFeesKey())
	if err != nil {
		return nil, err
	}

	var remainingServiceFees types.Fees
	k.cdc.MustUnmarshalLengthPrefixed(bz, &remainingServiceFees)
	return remainingServiceFees.Fees, nil
}
