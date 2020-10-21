package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) ClaimLock(ctx sdk.Context, proposalID, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, lockPeriod time.Duration) error {
	return nil
}

func (k Keeper) ClaimUnlock(ctx sdk.Context, proposalID, poolID uint64, loss sdk.Coins) error {
	return nil
}

func (k Keeper) RestoreShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, id uint64, loss sdk.Coins) error {
	return nil
}

// SetReimbursement sets a reimbursement in store.
func (k Keeper) SetReimbursement(ctx sdk.Context, proposalID uint64, payout types.Reimbursement) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(payout)
	store.Set(types.GetReimbursementKey(proposalID), bz)
}

// CreateReimbursement creates a reimbursement.
func (k Keeper) CreateReimbursement(ctx sdk.Context, proposalID uint64, poolID uint64, rmb sdk.Coins, beneficiary sdk.AccAddress) error {
	reimbursement := types.NewReimbursement(rmb, beneficiary, ctx.BlockTime().Add(k.GetClaimProposalParams(ctx).PayoutPeriod))
	k.SetReimbursement(ctx, proposalID, reimbursement)
	return nil
}
