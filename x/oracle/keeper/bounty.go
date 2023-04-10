package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// SetRemainingBounty This function saves the RemainingBounty struct associated to its address in the store.
func (k Keeper) SetRemainingBounty(ctx sdk.Context, remainingBounty types.RemainingBounty) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&remainingBounty)
	addr := sdk.MustAccAddressFromBech32(remainingBounty.Address)
	store.Set(types.RemainingBountyStoreKey(addr), bz)
}

// GetRemainingBounty This function retrieves the RemainingBounty of a given address from the store.
func (k Keeper) GetRemainingBounty(ctx sdk.Context, address sdk.AccAddress) (types.RemainingBounty, error) {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.RemainingBountyStoreKey(address))
	if opBz != nil {
		var remainingBounty types.RemainingBounty
		k.cdc.MustUnmarshalLengthPrefixed(opBz, &remainingBounty)
		return remainingBounty, nil
	}
	return types.RemainingBounty{}, types.ErrNoRemainingBountyFound
}

// DeleteRemainingBounty This function deletes the store key associated with a Creator's remaining Bounty within the context of the given Keeper.
func (k Keeper) DeleteRemainingBounty(ctx sdk.Context, address sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.RemainingBountyStoreKey(address))
	return nil
}

func (k Keeper) GetAllRemainingBounties(ctx sdk.Context) []types.RemainingBounty {
	var remainingBounties []types.RemainingBounty
	k.IterateAllRemainingBounties(ctx, func(remainingBounty types.RemainingBounty) bool {
		remainingBounties = append(remainingBounties, remainingBounty)
		return false
	})
	return remainingBounties
}

// IterateAllRemainingBounties This function enables iteration over all RemainingBounties stored in the store while giving the power to break out early if the callback function returns true.
func (k Keeper) IterateAllRemainingBounties(ctx sdk.Context, callback func(remainingBounty types.RemainingBounty) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RemainingBountyStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var remainingBounty types.RemainingBounty
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &remainingBounty)
		if callback(remainingBounty) {
			break
		}
	}
}
