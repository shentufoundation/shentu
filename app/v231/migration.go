package v231

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	bankkeeper "github.com/certikfoundation/shentu/v2/x/bank/keeper"
	shieldkeeper "github.com/certikfoundation/shentu/v2/x/shield/keeper"
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1alpha1"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
	stakingkeeper "github.com/certikfoundation/shentu/v2/x/staking/keeper"
)

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

func ExpireStakingPurchase(ctx sdk.Context, cdc codec.BinaryCodec, bk bankkeeper.Keeper, k shieldkeeper.Keeper, sk *stakingkeeper.Keeper, storeKey sdk.StoreKey) {
	bondDenom := sk.BondDenom(ctx)
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, shieldtypes.PurchaseKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ss v1alpha1.ShieldStaking
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &ss)
		addr, err := sdk.AccAddressFromBech32(ss.Purchaser)
		if err != nil {
			panic(err)
		}
		err = bk.SendCoinsFromModuleToAccount(ctx, shieldtypes.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(bondDenom, ss.Amount)))
		if err != nil {
			panic(err)
		}
		store.Delete(iterator.Key())
	}
}

func RefundPurchasers(ctx sdk.Context, cdc codec.BinaryCodec, ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, sk *stakingkeeper.Keeper, k shieldkeeper.Keeper, storeKey sdk.StoreKey) {
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

	// directly pay out provider rewards
	remainingFees := totalFees
	defer iterator2.Close()
	for ; iterator2.Valid(); iterator2.Next() {
		var pv v1alpha1.Provider
		cdc.MustUnmarshalLengthPrefixed(iterator2.Value(), &pv)
		rewardsInt := sdk.NewCoins()
		remainders := sdk.NewDecCoins()
		for _, r := range pv.Rewards.Native {
			rInt, remainder := r.TruncateDecimal()
			remainders = remainders.Add(remainder)
			remainingFees = remainingFees.Sub(r.Amount)
			rewardsInt = rewardsInt.Add(rInt)
		}
		addr, err := sdk.AccAddressFromBech32(pv.Address)
		if err != nil {
			panic(err)
		}
		err = bk.SendCoinsFromModuleToAccount(ctx, shieldtypes.ModuleName, addr, rewardsInt)
		if err != nil {
			panic(err)
		}
		pv.Rewards.Native = sdk.NewDecCoins()
		pvBz := cdc.MustMarshalLengthPrefixed(&pv)
		store.Set(iterator2.Key(), pvBz)
	}

	var refundRatio sdk.Dec
	if !totalFees.IsZero() && remainingFees.IsPositive() {
		refundRatio = remainingFees.Quo(totalFees)
	} else {
		refundRatio = sdk.ZeroDec()
	}

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
		pRInt := purchaserReimbursement.TruncateInt()
		if pRInt.IsPositive() {
			if err := bk.SendCoinsFromModuleToAccount(ctx, shieldtypes.ModuleName, addr,
				sdk.NewCoins(sdk.NewCoin(bondDenom, pRInt))); err != nil {
				panic(err)
			}
			remainingFees = remainingFees.Sub(pRInt.ToDec())
		}
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

	// initialize reserve with the remaining module account balance, since nothing else is left
	acc := ak.GetModuleAccount(ctx, shieldtypes.ModuleName)
	mAccBalances := bk.GetAllBalances(ctx, acc.GetAddress())
	reserve := v1beta1.NewReserve()
	reserve.Amount = mAccBalances.AmountOf(bondDenom).ToDec()
	k.SetReserve(ctx, reserve)

	// initialize and zero out any other trackers
	k.SetServiceFees(ctx, sdk.NewDecCoins())
	k.SetTotalShield(ctx, sdk.ZeroInt())
	k.SetGlobalStakingPool(ctx, sdk.ZeroInt())
}
