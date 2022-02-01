package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// Donate donates the given amount to Shield Donation Pool.
func (k Keeper) Donate(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error {
	donationPool := k.GetDonationPool(ctx)

	if err := k.bk.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amount))); err != nil {
		return err
	}

	donationPool.Amount = donationPool.Amount.Add(amount)
	k.SetDonationPool(ctx, donationPool)

	return nil
}

// SetDonationPool saves Shield Donation Pool.
func (k Keeper) SetDonationPool(ctx sdk.Context, donationPool types.DonationPool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&donationPool)
	store.Set(types.GetDonationPoolKey(), bz)
}

// GetDonationPool retrieves Shield Donation Pool.
func (k Keeper) GetDonationPool(ctx sdk.Context) (donationPool types.DonationPool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDonationPoolKey())
	if bz == nil {
		panic("failed to retrieve Shield Donation Pool")
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &donationPool)
	return
}

// SetPendingPayout stores a pending payout.
func (k Keeper) SetPendingPayout(ctx sdk.Context, pp types.PendingPayout) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&pp)
	store.Set(types.GetPendingPayoutKey(pp.ProposalId), bz)
}

// GetPendingPayout retrieves a pending payout.
func (k Keeper) GetPendingPayout(ctx sdk.Context, proposalId uint64) (pp types.PendingPayout, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPendingPayoutKey(proposalId))
	if bz == nil {
		return pp, false
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &pp)
	return pp, true
}

// GetAllPendingPayouts retrieves all pending payouts.
func (k Keeper) GetAllPendingPayouts(ctx sdk.Context) (payouts []types.PendingPayout) {
	k.IteratePendingPayouts(ctx, func(payout types.PendingPayout) bool {
		payouts = append(payouts, payout)
		return false
	})
	return
}

// IteratePendingPayouts iterates through all pending payouts.
func (k Keeper) IteratePendingPayouts(ctx sdk.Context, callback func(pp types.PendingPayout) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PendingPayoutKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var payout types.PendingPayout
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &payout)

		if callback(payout) {
			break
		}
	}
}

// // MakePayouts ...
// func (k Keeper) MakePayouts(ctx sdk.Context) {

// }
