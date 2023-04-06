package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// SetCreatorLeftBounty This function saves the LeftBounty struct associated to its address in the store.
func (k Keeper) SetCreatorLeftBounty(ctx sdk.Context, leftBounty types.LeftBounty) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(&leftBounty)
	addr := sdk.MustAccAddressFromBech32(leftBounty.Address)
	store.Set(types.LeftBountyStoreKey(addr), bz)
}

// GetCreatorLeftBounty This function retrieves the LeftBounty of a given address from the store.
func (k Keeper) GetCreatorLeftBounty(ctx sdk.Context, address sdk.AccAddress) (types.LeftBounty, error) {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.LeftBountyStoreKey(address))
	if opBz != nil {
		var leftBounty types.LeftBounty
		k.cdc.MustUnmarshalLengthPrefixed(opBz, &leftBounty)
		return leftBounty, nil
	}
	return types.LeftBounty{}, types.ErrNoLeftBountyFound
}

// DeleteCreatorLeftBounty This function deletes the store key associated with a Creator's Left Bounty within the context of the given Keeper.
func (k Keeper) DeleteCreatorLeftBounty(ctx sdk.Context, address sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.LeftBountyStoreKey(address))
	return nil
}

func (k Keeper) GetAllLeftBounties(ctx sdk.Context) []types.LeftBounty {
	var leftBounties []types.LeftBounty
	k.IterateAllLeftBounties(ctx, func(leftBounty types.LeftBounty) bool {
		leftBounties = append(leftBounties, leftBounty)
		return false
	})
	return leftBounties
}

// IterateAllLeftBounties This function enables iteration over all LeftBounties stored in the store while giving the power to break out early if the callback function returns true.
func (k Keeper) IterateAllLeftBounties(ctx sdk.Context, callback func(leftBounty types.LeftBounty) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LeftBountyStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var leftBounty types.LeftBounty
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &leftBounty)
		if callback(leftBounty) {
			break
		}
	}
}
