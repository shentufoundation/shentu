package bounty

import (
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k *keeper.Keeper) error {

	logger := ctx.Logger().With("module", "x/"+types.ModuleName)

	// delete dead theorems from store and returns theirs grant.
	// A theorem is dead when it's active and didn't get correct proof on time to get into pass phase.
	rng := collections.NewPrefixUntilPairRange[time.Time, uint64](ctx.BlockTime())
	err := k.ActiveTheoremsQueue.Walk(ctx, rng, func(key collections.Pair[time.Time, uint64], value uint64) (stop bool, err error) {
		theorem, err := k.Theorems.Get(ctx, key.K2())
		if err != nil {
			return false, err
		}

		// TODO add to a func
		err = k.ActiveTheoremsQueue.Remove(ctx, collections.Join(*theorem.EndTime, theorem.Id))
		if err != nil {
			return false, err
		}
		err = k.Theorems.Remove(ctx, theorem.Id)
		if err != nil {
			return false, err
		}

		err = k.RefundAndDeleteGrants(ctx, theorem.Id)
		if err != nil {
			return false, err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeActiveTheorem,
				sdk.NewAttribute(types.AttributeKeyTheoremID, fmt.Sprintf("%d", theorem.Id)),
				sdk.NewAttribute(types.AttributeKeyTheoremResult, types.AttributeValueTheoremDropped),
			),
		)

		logger.Info(
			"theorem did not meet correct proof on time; deleted",
			"theorem", theorem.Id,
			"title", theorem.Title,
			"total_deposit", sdk.NewCoins(theorem.TotalGrant...).String(),
		)

		return false, nil
	})

	if err != nil {
		return err
	}

	return nil
}
