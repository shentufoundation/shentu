package v231

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bankkeeper "github.com/certikfoundation/shentu/v2/x/bank/keeper"
	shieldkeeper "github.com/certikfoundation/shentu/v2/x/shield/keeper"
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1alpha1"
	stakingkeeper "github.com/certikfoundation/shentu/v2/x/staking/keeper"
)

func RefundPurchasers(ctx sdk.Context, cdc codec.BinaryCodec, bk bankkeeper.Keeper, sk *stakingkeeper.Keeper, k shieldkeeper.Keeper, storeKey sdk.StoreKey) {
	bondDenom := sk.BondDenom(ctx)

	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, shieldtypes.PurchaseListKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pl v1alpha1.PurchaseList
		cdc.MustUnmarshal(iterator.Value(), &pl)
		total := sdk.ZeroDec()
		for _, e := range pl.Entries {
			total.Add(e.ServiceFees.Native.AmountOf(bondDenom))
		}
		addr, err := sdk.AccAddressFromBech32(pl.Purchaser)
		if err != nil {
			panic(err)
		}
		if err := bk.SendCoinsFromModuleToAccount(ctx, shieldtypes.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(bondDenom, total.TruncateInt()))); err != nil {
			panic(err)
		}
		store.Delete(shieldtypes.GetPurchaseListKey(pl.PoolId, addr))
	}

	k.SetTotalShield(ctx, sdk.ZeroInt())
	k.SetGlobalStakingPool(ctx, sdk.ZeroInt())
}
