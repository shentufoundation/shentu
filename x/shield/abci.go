package shield

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/keeper"
)

// BeginBlocker executes logics to begin a block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
}

// EndBlocker processes premium payment at every block.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Remove expired purchases and distribute service fees.
	k.RemoveExpiredPurchasesAndDistributeFees(ctx)

	// Process completed withdraws.
	k.DequeueCompletedWithdrawQueue(ctx)

	// Close pools who do not have any shield and shield limits are set to zero.
	k.ClosePools(ctx)
}
