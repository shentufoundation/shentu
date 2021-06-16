package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	nftkeeper "github.com/irisnet/irismod/modules/nft/keeper"

	"github.com/certikfoundation/shentu/x/nft/types"
)

type Keeper struct {
	nftkeeper.Keeper
	certKeeper types.CertKeeper
	storeKey   sdk.StoreKey
	cdc        codec.Marshaler
}

// NewKeeper creates a new instance of the NFT Keeper
func NewKeeper(cdc codec.Marshaler, certKeeper types.CertKeeper, storeKey sdk.StoreKey) Keeper {
	baseKeeper := nftkeeper.NewKeeper(cdc, storeKey)
	return Keeper{
		Keeper:     baseKeeper,
		certKeeper: certKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) DeleteAdmin(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AdminKey(addr))
}

func (k Keeper) SetAdmin(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	newAdmin := types.Admin{
		Address: addr.String(),
	}
	bz := k.cdc.MustMarshalBinaryBare(&newAdmin)
	store.Set(types.AdminKey(addr), bz)
}

func (k Keeper) GetAdmin(ctx sdk.Context, addr sdk.AccAddress) (types.Admin, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AdminKey(addr))
	if bz == nil {
		return types.Admin{}, sdkerrors.Wrapf(types.ErrAdminNotFound, "not found NFT: %s", addr.String())
	}
	var admin types.Admin
	k.cdc.MustUnmarshalBinaryBare(bz, &admin)
	return admin, nil
}

func (k Keeper) GetAdmins(ctx sdk.Context) []types.Admin {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.AdminKeyPrefix)
	var res []types.Admin
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var admin types.Admin
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &admin)
		res = append(res, admin)
	}

	return res
}

func (k Keeper) CheckAdmin(ctx sdk.Context, addr string) bool {
	admin, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return false
	}
	_, err = k.GetAdmin(ctx, admin)
	return err == nil
}
