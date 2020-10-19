package shield

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// BeginBlock executes logics to begin a block.
func BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
}

// EndBlocker processes premium payment at every block.
func EndBlocker(ctx sdk.Context, k Keeper) {
	// Distribute service fees to all providers.
	serviceFees := k.GetServiceFees(ctx)
	totalCollateral := k.GetTotalCollateral(ctx)
	totalLocked := k.GetTotalLocked(ctx)

	totalCollateralAmount := totalCollateral.Add(totalLocked)
	providers := k.GetAllProviders(ctx)
	for _, provider := range providers {
		proportion := sdk.NewDecFromInt(sdk.MaxInt(provider.Collateral.Add(provider.TotalLocked), sdk.ZeroInt())).QuoInt(totalCollateralAmount)
		nativeFees := serviceFees.Native.MulDecTruncate(proportion)
		foreignFees := serviceFees.Foreign.MulDecTruncate(proportion)

		serviceFees.Native = serviceFees.Native.Sub(nativeFees)
		serviceFees.Foreign = serviceFees.Foreign.Sub(foreignFees)

		rewards := types.NewMixedDecCoins(nativeFees, foreignFees)
		provider.Rewards.Add(rewards)
		k.SetProvider(ctx, provider.Address, provider)
	}
	k.SetServiceFees(ctx, serviceFees)

	// Remove expired purchases.
	k.RemoveExpiredPurchases(ctx)

	// Process completed withdraws.
	k.DequeueCompletedWithdrawQueue(ctx)
}
