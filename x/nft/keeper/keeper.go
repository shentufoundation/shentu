package keeper

import (
	certkeeper "github.com/certikfoundation/shentu/x/cert/keeper"
	customtypes "github.com/certikfoundation/shentu/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irisnet/irismod/modules/nft/keeper"
)

type Keeper struct {
	keeper.Keeper
	certKeeper certkeeper.Keeper
	storeKey sdk.StoreKey
	cdc      codec.Marshaler
}

func (k Keeper) DeleteAdmin(ctx sdk.Context, addr sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(customtypes.AdminKey(addr))
	return nil
}

func (k Keeper) SetAdmin(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	newAdmin := customtypes.Admin{
		Address: addr.String(),
	}
	bz := k.cdc.MustMarshalBinaryBare(&newAdmin)
	store.Set(customtypes.AdminKey(addr), bz)
}

func (k Keeper) GetAdmin(ctx sdk.Context, addr sdk.AccAddress) (customtypes.Admin, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(customtypes.AdminKey(addr))
	if bz == nil {
		return customtypes.Admin{}, sdkerrors.Wrapf(customtypes.ErrAdminNotFound, "not found NFT: %s", addr.String())
	}
	var admin customtypes.Admin
	k.cdc.MustUnmarshalBinaryBare(bz, &admin)
	return admin, nil
}

func (k Keeper) GetAdmins(ctx sdk.Context) []customtypes.Admin {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, customtypes.AdminKeyPrefix)
	var res []customtypes.Admin
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var admin customtypes.Admin
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &admin)
		res = append(res, admin)
	}

	return res
}