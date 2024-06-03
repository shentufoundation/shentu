package shield

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/shield/keeper"
)

// EndBlocker processes premium payment at every block.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Remove expired purchases and distribute service fees.
	k.RemoveExpiredPurchasesAndDistributeFees(ctx)

	// Process completed withdraws.
	k.DequeueCompletedWithdrawQueue(ctx)

	// Close pools who do not have any shield and shield limits are set to zero.
	k.ClosePools(ctx)
}
