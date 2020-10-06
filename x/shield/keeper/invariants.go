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
		var id uint64
		broken := false
		pool := types.Pool{}
		purchased := sdk.Coins{}
		k.IterateAllPools(ctx, func(pool types.Pool) bool {
			purchases := k.GetAllPurchases(ctx)
			purchased = sdk.Coins{}
			for _, purchase := range purchases {
				if purchase.PoolID == pool.PoolID {
					purchased.Add(purchase.Shield...)
				}
			}
			id = pool.PoolID
			broken = pool.TotalCollateral.IsAllLT(purchased)
			return broken
		})
		return sdk.FormatInvariant(types.ModuleName, "account collateral and total sum of deposited collateral",
			fmt.Sprintf("\tPool ID: %v\n"+
				"\tSum of purchased Shield: %v\n"+
				"\tPool's total collaterals: %v\n", id, purchased, pool.TotalCollateral)), broken
	}
}
