package v231

import (
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
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

	// aggregate total service fees (including the rewards to be paid to the providers)
	totalFees := sdk.ZeroDec()
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pl v1alpha1.PurchaseList
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &pl)
		for _, e := range pl.Entries {
			totalFees = totalFees.Add(e.ServiceFees.Native.AmountOf(bondDenom))
		}
	}

	iterator2 := sdk.KVStorePrefixIterator(store, shieldtypes.ProviderKey)

	// directly pay out rewards
	var prov []v1alpha1.Provider
	remainingFees := totalFees
	defer iterator2.Close()
	for ; iterator2.Valid(); iterator2.Next() {
		var pv v1alpha1.Provider
		cdc.MustUnmarshalLengthPrefixed(iterator2.Value(), &pv)
		prov = append(prov, pv)

		addr, err := sdk.AccAddressFromBech32(pv.Address)
		if err != nil {
			panic(err)
		}
		rewardsInt := sdk.NewCoins()
		for _, r := range pv.Rewards.Native {
			rInt, _ := r.TruncateDecimal()
			remainingFees = remainingFees.Sub(rInt.Amount.ToDec())
			rewardsInt = rewardsInt.Add(rInt)
		}
		err = bk.SendCoinsFromModuleToAccount(ctx, shieldtypes.ModuleName, addr, rewardsInt)
		if err != nil {
			panic(err)
		}
		pv.Rewards = v1alpha1.NewMixedDecCoins(sdk.NewDecCoins(), sdk.NewDecCoins())
		pvBz := cdc.MustMarshalLengthPrefixed(&pv)
		store.Set(iterator2.Key(), pvBz)
	}

	refundRatio := remainingFees.Quo(totalFees)

	iterator3 := sdk.KVStorePrefixIterator(store, shieldtypes.PurchaseListKey)

	// send remaining service fees to purchasers proportionally
	defer iterator3.Close()
	for ; iterator3.Valid(); iterator3.Next() {
		var pl v1alpha1.PurchaseList
		cdc.MustUnmarshalLengthPrefixed(iterator3.Value(), &pl)
		purchaserTotal := sdk.ZeroDec()
		for _, e := range pl.Entries {
			purchaserTotal = purchaserTotal.Add(e.ServiceFees.Native.AmountOf(bondDenom))
		}
		addr, err := sdk.AccAddressFromBech32(pl.Purchaser)
		if err != nil {
			panic(err)
		}
		purchaserReimbursement := purchaserTotal.Mul(refundRatio)
		if err := bk.SendCoinsFromModuleToAccount(ctx, shieldtypes.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(bondDenom, purchaserReimbursement.TruncateInt()))); err != nil {
			panic(err)
		}
		remainingFees = remainingFees.Sub(purchaserReimbursement)
		store.Delete(iterator3.Key())
	}

	// reset pool shield to 0
	iterator4 := sdk.KVStorePrefixIterator(store, shieldtypes.PoolKey)
	defer iterator4.Close()
	for ; iterator4.Valid(); iterator4.Next() {
		var pool v1alpha1.Pool
		cdc.MustUnmarshalLengthPrefixed(iterator4.Value(), &pool)

		pool.Shield = sdk.ZeroInt()
		poolBz := cdc.MustMarshalLengthPrefixed(&pool)
		store.Set(iterator4.Key(), poolBz)
	}

	reserve := v1beta1.NewReserve()
	reserve.Amount = reserve.Amount.Add(remainingFees.TruncateInt())
	k.SetReserve(ctx, reserve)
	k.SetServiceFees(ctx, sdk.NewDecCoins())
	k.SetTotalShield(ctx, sdk.ZeroInt())
	k.SetGlobalStakingPool(ctx, sdk.ZeroInt())
}

func PayoutReimbursements(ctx sdk.Context, cdc codec.BinaryCodec, bk bankkeeper.Keeper, k shieldkeeper.Keeper, storeKey sdk.StoreKey) {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, shieldtypes.ReimbursementKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var reimbursement v1alpha1.Reimbursement
		cdc.MustUnmarshal(iterator.Value(), &reimbursement)
		addr, err := sdk.AccAddressFromBech32(reimbursement.Beneficiary)
		if err != nil {
			panic(err)
		}
		if err := bk.SendCoinsFromModuleToAccount(ctx, shieldtypes.ModuleName, addr, reimbursement.Amount); err != nil {
			panic(err)
		}
		store.Delete(iterator.Key())
	}

	k.SetTotalClaimed(ctx, sdk.ZeroInt())
}
