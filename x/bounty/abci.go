package bounty

import (
	"time"

	"cosmossdk.io/collections"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k *keeper.Keeper) error {

	logger := ctx.Logger().With("module", "x/"+types.ModuleName)

	rngProof := collections.NewPrefixUntilPairRange[time.Time, string](ctx.BlockTime())
	err := k.ActiveProofsQueue.Walk(ctx, rngProof, func(key collections.Pair[time.Time, string], value types.Proof) (stop bool, err error) {
		// Only delete proofs in hash_lock phase
		if value.Status != types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD {
			return false, nil
		}

		err = k.DeleteProof(ctx, value.Id)
		if err != nil {
			return false, err
		}

		logger.Info(
			"proof did not submit detail on time; expired",
			"proof_id", value.Id,
		)

		return false, nil
	})
	if err != nil {
		return err
	}

	// delete dead theorems from store and returns theirs grant.
	// A theorem is dead when it's active and didn't get correct proof on time to get into pass phase.
	rngTheorem := collections.NewPrefixUntilPairRange[time.Time, uint64](ctx.BlockTime())
	err = k.ActiveTheoremsQueue.Walk(ctx, rngTheorem, func(key collections.Pair[time.Time, uint64], value uint64) (stop bool, err error) {
		theorem, err := k.Theorems.Get(ctx, key.K2())
		if err != nil {
			return false, err
		}

		// Check if theorem has any related proofs
		hasActiveProof, _, err := k.HasActiveProofs(ctx, theorem.Id)
		if err != nil {
			return false, err
		}

		// If there's a valid proof, keep the theorem
		if hasActiveProof {
			return false, nil
		}

		if err = k.DeleteTheorem(ctx, theorem.Id); err != nil {
			return false, err
		}

		err = k.RefundAndDeleteGrants(ctx, theorem.Id)
		if err != nil {
			return false, err
		}

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
