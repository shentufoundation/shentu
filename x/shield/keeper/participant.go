package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) SetParticipant(ctx sdk.Context, delAddr sdk.AccAddress, participant types.Participant) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(participant)
	store.Set(types.GetParticipantKey(delAddr), bz)
}

func (k Keeper) GetParticipant(ctx sdk.Context, delegator sdk.AccAddress) (dt types.Participant, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetParticipantKey(delegator))
	if bz == nil {
		return types.Participant{}, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &dt)
	return dt, true
}

func (k Keeper) updateDelegationAmount(ctx sdk.Context, delAddr sdk.AccAddress) {
	// go through delAddr's delegations to recompute total amount of bonded delegation
	delegations := k.sk.GetAllDelegatorDelegations(ctx, delAddr)
	bondedAmount := sdk.Coins{}
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("expected validator, not found")
		}
		bondedAmount = bondedAmount.Add(sdk.NewCoin(k.sk.BondDenom(ctx), val.TokensFromShares(del.GetShares()).TruncateInt()))
	}

	// update or create a new entry
	participant, found := k.GetParticipant(ctx, delAddr)
	if !found {
		return // ignore non-participating addr
	}
	originalBonded := participant.DelegationBonded
	participant.DelegationBonded = bondedAmount

	// if bonded decreased, an unbonding began
	if bondedAmount.IsAllLT(originalBonded) {
		participant.DelegationUnbonding = participant.DelegationUnbonding.Add(originalBonded.Sub(bondedAmount)...)
	}

	totalDelegation := participant.DelegationBonded.Add(participant.DelegationUnbonding...)
	if totalDelegation.IsAllLT(participant.Collateral) {
		participant.Collateral = totalDelegation
		withdrawAmount := participant.Collateral.Sub(totalDelegation)
		k.WithdrawFromPools(ctx, delAddr, withdrawAmount)
	}
	k.SetParticipant(ctx, delAddr, participant)
}
