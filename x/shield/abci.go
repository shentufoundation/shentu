package shield

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/shield/types"
)

// BeginBlock executes logics to begin a block
func BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
}

// EndBlocker processes premium payment at every block.
func EndBlocker(ctx sdk.Context, k Keeper, stk types.StakingKeeper) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		if k.PoolEnded(ctx, pool) || (pool.Premium.Native.Empty() && pool.Premium.Foreign.Empty()) {
			k.ClosePool(ctx, pool)
			continue
		}
		// compute premiums for current block
		var currentBlockPremium types.MixedDecCoins
		if pool.EndTime != 0 {
			// use endTime to compute premiums
			timeUntilEnd := pool.EndTime - ctx.BlockTime().Unix()
			blocksUntilEnd := sdk.MaxDec(common.BlocksPerSecondDec.Mul(sdk.NewDec(timeUntilEnd)), sdk.OneDec())
			if ctx.BlockTime().Unix() > pool.EndTime {
				// must spend all premium
				currentBlockPremium = pool.Premium
			} else {
				currentBlockPremium = pool.Premium.QuoDec(blocksUntilEnd)
			}
		} else {
			// use block height to compute premiums
			blocksUntilEnd := sdk.NewDec(pool.EndBlockHeight - ctx.BlockHeight())
			if ctx.BlockHeight() >= pool.EndBlockHeight {
				// must spend all premium
				currentBlockPremium = pool.Premium
			} else {
				currentBlockPremium = pool.Premium.QuoDec(blocksUntilEnd)
			}
		}

		// distribute to A and C in proportion
		bondDenom := stk.BondDenom(ctx) // common.MicroCTKDenom
		totalCollatInt := pool.TotalCollateral.AmountOf(bondDenom)
		recipients := append(pool.Community, pool.CertiK)
		for _, recipient := range recipients {
			stakeProportion := sdk.NewDecFromInt(recipient.Amount.AmountOf(bondDenom)).QuoInt(totalCollatInt)
			nativePremium := currentBlockPremium.Native.MulDecTruncate(stakeProportion)
			foreignPremium := currentBlockPremium.Foreign.MulDecTruncate(stakeProportion)

			pool.Premium.Native = pool.Premium.Native.Sub(nativePremium)

			pool.Premium.Foreign = pool.Premium.Foreign.Sub(foreignPremium)

			rewards := types.NewMixedDecCoins(nativePremium, foreignPremium)
			k.AddRewards(ctx, recipient.Provider, rewards)
		}

		k.SetPool(ctx, pool)
	} // for each pool

	// remove expired purchases
	k.RemoveExpiredPurchases(ctx)

	// track unbonding amounts of participants
	TrackUnbondingAmount(ctx, k, stk)

	// process completed withdrawals
	// Remove all mature unbonding delegations from the ubd queue.
	k.DequeueCompletedWithdrawalQueue(ctx)
}

// TrackUnbondingAmount tracks the amount to be unbonded by staking end blocker.
func TrackUnbondingAmount(ctx sdk.Context, k Keeper, stk types.StakingKeeper) {
	matureUnbonds := k.GetAllMatureUBDQueue(ctx, ctx.BlockHeader().Time, stk)
	for _, dvPair := range matureUnbonds {
		delAddr := dvPair.DelegatorAddress
		valAddr := dvPair.ValidatorAddress
		ubd, _ := stk.GetUnbondingDelegation(ctx, delAddr, valAddr)
		bondDenom := stk.BondDenom(ctx)
		balances := sdk.NewCoins()
		ctxTime := ctx.BlockHeader().Time

		// loop through all the entries and complete unbonding mature entries
		for _, entry := range ubd.Entries {
			if entry.IsMature(ctxTime) && !entry.Balance.IsZero() {
				balances = balances.Add(sdk.NewCoin(bondDenom, entry.Balance))
			}
		}

		participant, _ := k.GetParticipant(ctx, delAddr)

		if participant.DelegationUnbonding.Empty() {
			continue
		}
		participant.DelegationUnbonding = participant.DelegationUnbonding.Sub(balances)
		k.SetParticipant(ctx, delAddr, participant)
	} // for each mature unbond dv pair
}
