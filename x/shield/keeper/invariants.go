package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// RegisterInvariants registers all shield invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "account-collaterals",
		AccountCollateralsInvariants(k))
	ir.RegisterRoute(types.ModuleName, "purchased-collaterals",
		PurchasedCollateralsInvariants(k))
	ir.RegisterRoute(types.ModuleName, "module-coins",
		ModuleCoinsInvariants(k))
}

// AllInvariants runs all invariants of the shield module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return AccountCollateralsInvariants(k)(ctx)
	}
}

// AccountCollateralInvariants checks that the total collaterals for an accounts has to equal the sum of
// the account's registered collaterals
func AccountCollateralsInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		broken := false
		providerCollateral := sdk.Coins{}
		providerCollateralSum := sdk.Coins{}
		k.IterateProviders(ctx, func(provider types.Provider) bool {
			providerCollaterals := k.GetOnesCollaterals(ctx, provider.Address)
			sum := sdk.Coins{}
			for _, collateral := range providerCollaterals {
				sum = sum.Add(collateral.Amount...)
			}
			providerCollateral = provider.Collateral
			providerCollateralSum = sum
			broken = !(sum.IsEqual(provider.Collateral))
			return broken
		})
		return sdk.FormatInvariant(types.ModuleName, "account collateral and total sum of deposited collateral",
			fmt.Sprintf("\tSum of Provider's deposited tokens: %v\n"+
				"\tAccount's tracked collaterals: %v\n", providerCollateralSum, providerCollateral)), broken
	}
}

// PurchasedCollateralsInvariants checks the total purchased amount is less than or equal to the pool's total collateral amount.
func PurchasedCollateralsInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		broken := false
		currentPool := types.Pool{}
		purchased := sdk.Coins{}
		k.IterateAllPools(ctx, func(pool types.Pool) bool {
			purchases := k.GetPoolPurchases(ctx, pool.PoolID)
			purchased = sdk.Coins{}
			for _, purchase := range purchases {
				purchased = purchased.Add(purchase.Shield...)
			}
			currentPool = pool
			broken = pool.TotalCollateral.IsAllLT(purchased)
			return broken
		})
		return sdk.FormatInvariant(types.ModuleName, "pool total collateral and total sum of purchased collateral",
			fmt.Sprintf("\tPool ID: %v\n"+
				"\tSum of purchased Shield: %v\n"+
				"\tPool's total collaterals: %v\n", currentPool.PoolID, purchased, currentPool.TotalCollateral)), broken
	}
}

// ModuleCoinsInvariants checks the total sum of coins equals to module account's balance.
func ModuleCoinsInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		unbondings := k.sk.GetAllUnbondingDelegations(ctx, k.supplyKeeper.GetModuleAddress(types.ModuleName))
		actualModuleCoinsAmt := sdk.NewInt(0)
		for _, ubd := range unbondings {
			for _, entry := range ubd.Entries {
				actualModuleCoinsAmt = actualModuleCoinsAmt.Add(entry.Balance)
			}
		}
		providers := k.GetAllProviders(ctx)
		rewardsDec := sdk.NewDec(0)
		for _, provider := range providers {
			rewardsDec = rewardsDec.Add(provider.Rewards.Native.AmountOf(k.sk.BondDenom(ctx)))
		}
		pools := k.GetAllPools(ctx)
		for _, pool := range pools {
			rewardsDec = rewardsDec.Add(pool.Premium.Native.AmountOf(k.sk.BondDenom(ctx)))
		}

		actualModuleCoinsAmt = actualModuleCoinsAmt.Add(rewardsDec.TruncateInt())

		expectedModuleCoinsAmt := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetCoins().AmountOf(k.sk.BondDenom(ctx))

		broken := !expectedModuleCoinsAmt.Equal(actualModuleCoinsAmt)
		return sdk.FormatInvariant(types.ModuleName, "module total sum of coins and module account coins",
			fmt.Sprintf("\tSum of premiums and unbondings: %v\n"+
				"\tmodule coins amount: %v\n", actualModuleCoinsAmt, expectedModuleCoinsAmt)), broken
	}
}
