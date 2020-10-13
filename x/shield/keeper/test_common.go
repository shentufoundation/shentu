package keeper

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// PrintPoolInfo prints info of a pool.
func PrintPoolInfo(k Keeper, ctx sdk.Context, poolID uint64, description string) {
	pool, err := k.GetPool(ctx, poolID)
	if err == nil {
		fmt.Printf("%s: pool ID %d, total collateral %s, available %s, shield %s\n",
			description, pool.PoolID, pool.TotalCollateral, pool.Available, pool.Shield)
	} else {
		fmt.Printf("%s error: pool is not found\n", description)
	}
}

// PrintCollateralInfo prints info of a collateral.
func PrintCollateralInfo(k Keeper, ctx sdk.Context, poolID uint64, addr sdk.AccAddress, description string) {
	pool, err := k.GetPool(ctx, poolID)
	if err != nil {
		fmt.Printf("%s error: pool is not found\n", description)
		return
	}
	collateral, found := k.GetCollateral(ctx, pool, addr)
	if found {
		fmt.Printf("%s: collateral of %s, pool ID %d, collateral %s\n", description, addr, pool.PoolID, collateral.Amount)
	} else {
		fmt.Printf("%s error: collateral is not found\n", description)
	}
}

// PrintProviderInfo prints info of a provider.
func PrintProviderInfo(k Keeper, ctx sdk.Context, addr sdk.AccAddress, description string) {
	provider, found := k.GetProvider(ctx, addr)
	if found {
		fmt.Printf("%s: provider %s, delegation %s, collateral %s, available %s\n",
			description, addr, provider.DelegationBonded, provider.Collateral, provider.Available)
	} else {
		fmt.Printf("%s error: provider is not found\n", description)
	}
}

// RandomValidator returns a random validator given access to the keeper and ctx.
func RandomValidator(r *rand.Rand, k Keeper, ctx sdk.Context) (staking.Validator, bool) {
	vals := k.sk.GetAllValidators(ctx)
	if len(vals) == 0 {
		return staking.Validator{}, false
	}

	i := r.Intn(len(vals))
	return vals[i], true
}

// RandomDelegation returns a random delegation info given access to the keeper and ctx.
func RandomDelegation(r *rand.Rand, k Keeper, ctx sdk.Context) (sdk.AccAddress, sdk.Int, bool) {
	val, ok := RandomValidator(r, k, ctx)
	if !ok {
		return nil, sdk.Int{}, false
	}

	dels := k.sk.GetValidatorDelegations(ctx, val.OperatorAddress)

	i := r.Intn(len(dels))
	return dels[i].DelegatorAddress, val.TokensFromShares(dels[i].Shares).TruncateInt(), true
}

// RandomPoolInfo returns info of a random pool given access to the keeper and ctx.
func RandomPoolInfo(r *rand.Rand, k Keeper, ctx sdk.Context) (uint64, string, bool) {
	pools := k.GetAllPools(ctx)
	if len(pools) == 0 {
		return 0, "", false
	}
	i := r.Intn(len(pools))
	return pools[i].PoolID, pools[i].Sponsor, true
}

// RandomCollateral returns a random collateral given access to the keeper and ctx.
func RandomCollateral(r *rand.Rand, k Keeper, ctx sdk.Context) (types.Collateral, bool) {
	poolID, _, found := RandomPoolInfo(r, k, ctx)
	if !found {
		return types.Collateral{}, false
	}
	pool, err := k.GetPool(ctx, poolID)
	if err != nil {
		return types.Collateral{}, false
	}
	collaterals := k.GetAllPoolCollaterals(ctx, pool)
	if len(collaterals) == 0 {
		return types.Collateral{}, false
	}
	i := r.Intn(len(collaterals))
	return collaterals[i], true
}

// RandomPurchaseList returns a random purchase given access to the keeper and ctx.
func RandomPurchaseList(r *rand.Rand, k Keeper, ctx sdk.Context) (types.PurchaseList, bool) {
	purchases := k.GetAllPurchaseLists(ctx)
	if len(purchases) == 0 {
		return types.PurchaseList{}, false
	}
	i := r.Intn(len(purchases))
	return purchases[i], true
}
