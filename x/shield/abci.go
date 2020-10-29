package shield

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/certikfoundation/shentu/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlock executes logics to begin a block.
func BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	if ctx.BlockHeight() == common.UpdateHeight {
		k.SetStakingPurchaseRate(ctx, sdk.NewDec(2))
	}
}

// EndBlocker processes premium payment at every block.
func EndBlocker(ctx sdk.Context, k Keeper) {
	// Remove expired purchases and distribute service fees.
	k.RemoveExpiredPurchasesAndDistributeFees(ctx)

	// Process completed withdraws.
	k.DequeueCompletedWithdrawQueue(ctx)

	// Close pools who do not have any shield and shield limits are set to zero.
	k.ClosePools(ctx)
}
