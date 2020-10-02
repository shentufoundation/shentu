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

// addParticipant adds a new participant into shield module. Should only be called from DepositCollateral.
func (k Keeper) addParticipant(ctx sdk.Context, addr sdk.AccAddress) types.Participant {
	delegations := k.sk.GetAllDelegatorDelegations(ctx, addr)

	totalStaked := sdk.Coins{}
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("expected validator, not found")
		}
		totalStaked = totalStaked.Add(sdk.NewCoin(k.sk.BondDenom(ctx), val.TokensFromShares(del.GetShares()).TruncateInt()))
	}
	participant := types.NewParticipant()
	participant.DelegationBonded = totalStaked

	k.SetParticipant(ctx, addr, participant)
	return participant
}

func (k Keeper) updateDelegationAmount(ctx sdk.Context, delAddr sdk.AccAddress) {
	// go through delAddr's delegations to recompute total amount of bonded delegation
	// update or create a new entry
	participant, found := k.GetParticipant(ctx, delAddr)
	if !found {
		return // ignore non-participating addr
	}

	// update delegations
	totalStaked := sdk.Coins{}
	delegations := k.sk.GetAllDelegatorDelegations(ctx, delAddr)
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("expected validator, not found")
		}
		totalStaked = totalStaked.Add(sdk.NewCoin(k.sk.BondDenom(ctx), val.TokensFromShares(del.GetShares()).TruncateInt()))
	}
	participant.DelegationBonded = totalStaked

	if participant.DelegationBonded.IsAllLT(participant.Collateral) {
		withdrawAmount := participant.Collateral.Sub(participant.DelegationBonded)
		k.WithdrawFromPools(ctx, delAddr, withdrawAmount)
	}
	k.SetParticipant(ctx, delAddr, participant)
}
