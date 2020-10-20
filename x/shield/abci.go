package shield

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlock executes logics to begin a block.
func BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
}

// EndBlocker processes premium payment at every block.
func EndBlocker(ctx sdk.Context, k Keeper) {
	// Distribute service fees to providers.
	k.DistributeFees(ctx)

	// Remove expired purchases.
	k.RemoveExpiredPurchases(ctx)

	// Process completed withdraws.
	k.DequeueCompletedWithdrawQueue(ctx)
}
