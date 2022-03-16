package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

// Donate donates the given amount to Shield Donation Pool.
func (k Keeper) Donate(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error {
	reserve := k.GetReserve(ctx)

	if err := k.bk.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amount))); err != nil {
		return err
	}

	reserve.Amount = reserve.Amount.Add(amount.ToDec())
	k.SetReserve(ctx, reserve)

	return nil
}

// SetReserve saves Shield Donation Pool.
func (k Keeper) SetReserve(ctx sdk.Context, reserve v1beta1.Reserve) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&reserve)
	store.Set(types.GetReserveKey(), bz)
}

// GetReserve retrieves Shield Donation Pool.
func (k Keeper) GetReserve(ctx sdk.Context) (reserve v1beta1.Reserve) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetReserveKey())
	if bz == nil {
		panic("failed to retrieve Shield Donation Pool")
	}
	k.cdc.MustUnmarshal(bz, &reserve)
	return
}

// SetPendingPayout stores a pending payout.
func (k Keeper) SetPendingPayout(ctx sdk.Context, pp v1beta1.PendingPayout) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&pp)
	store.Set(types.GetPendingPayoutKey(pp.ProposalId), bz)
}

// DeletePendingPayout deletes a pending payout given its proposal ID.
func (k Keeper) DeletePendingPayout(ctx sdk.Context, proposalID uint64) error {
	store := ctx.KVStore(k.storeKey)
	if _, found := k.GetPendingPayout(ctx, proposalID); !found {
		return types.ErrPendingPayoutNotFound
	}
	store.Delete(types.GetPendingPayoutKey(proposalID))
	return nil
}

// GetPendingPayout retrieves a pending payout.
func (k Keeper) GetPendingPayout(ctx sdk.Context, proposalId uint64) (pp v1beta1.PendingPayout, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPendingPayoutKey(proposalId))
	if bz == nil {
		return pp, false
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &pp)
	return pp, true
}

// GetAllPendingPayouts retrieves all pending payouts.
func (k Keeper) GetAllPendingPayouts(ctx sdk.Context) (payouts []v1beta1.PendingPayout) {
	k.IteratePendingPayouts(ctx, func(payout v1beta1.PendingPayout) bool {
		payouts = append(payouts, payout)
		return false
	})
	return
}

// IteratePendingPayouts iterates through all pending payouts.
func (k Keeper) IteratePendingPayouts(ctx sdk.Context, callback func(pp v1beta1.PendingPayout) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PendingPayoutKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var payout v1beta1.PendingPayout
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &payout)

		if callback(payout) {
			break
		}
	}
}

// ProcessPendingPayout processes the given amount from a pending
// payout.
func (k Keeper) ProcessPendingPayout(ctx sdk.Context, pp v1beta1.PendingPayout, amount sdk.Int) error {
	//reimb, err := k.GetReimbursement(ctx, pp.ProposalId)
	//if err != nil {
	//	return types.ErrReimbursementNotFound
	//}
	//reimb.Amount = reimb.Amount.Add(sdk.NewCoin(k.BondDenom(ctx), amount))
	//k.SetReimbursement(ctx, pp.ProposalId, reimb)

	beneficiary, err := k.gk.GetProposalProposer(ctx, pp.ProposalId)
	if err != nil {
		return err
	}
	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, beneficiary, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amount))); err != nil {
		return err
	}

	pp.Amount = pp.Amount.Sub(amount)
	if pp.Amount.IsZero() {
		if pp.Amount.IsNegative() { //testing purpose
			panic("negative pending payout amount")
		}
		k.DeletePendingPayout(ctx, pp.ProposalId)
	} else {
		k.SetPendingPayout(ctx, pp)
	}
	return nil
}

// MakePayouts makes payouts from reserve to pending payouts.
// It processes as many pending payouts as possible.
// TODO: Order matters??
func (k Keeper) MakePayouts(ctx sdk.Context) {
	reserve := k.GetReserve(ctx)

	k.IteratePendingPayouts(ctx, func(payout v1beta1.PendingPayout) bool {
		if reserve.Amount.IsZero() {
			if reserve.Amount.IsNegative() { //testing purpose
				panic("negative reserve balance")
			}
			return true
		}

		var amount sdk.Int
		if reserve.Amount.GTE(payout.Amount.ToDec()) {
			amount = payout.Amount
		} else {
			amount = reserve.Amount.TruncateInt()
		}

		k.ProcessPendingPayout(ctx, payout, amount)
		reserve.Amount = reserve.Amount.Sub(amount.ToDec())
		return false
	})

	k.SetReserve(ctx, reserve)
}
