package keeper

import (
	certtypes "github.com/certikfoundation/shentu/x/cert/types"
	customtypes "github.com/certikfoundation/shentu/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/nft/keeper"

	certkeeper "github.com/certikfoundation/shentu/x/cert/keeper"
)

type Keeper struct {
	keeper.Keeper
	certKeeper certkeeper.Keeper
	storeKey sdk.StoreKey
	cdc      codec.Marshaler
}

func (k Keeper) CreateNFTAdmin(ctx sdk.Context, addr sdk.AccAddress) error {
	certifiers := k.certKeeper.GetAllCertifiers(ctx)
	for _, c := range certifiers {
		certaddr, err := sdk.AccAddressFromBech32(c.Address)
		if err != nil {
			return certtypes.ErrCertifierNotExists
		}
		if addr.Equals(certaddr) {
			k.SetAdmin(ctx, addr)
			return nil
		}
	}
	return certtypes.ErrCertifierNotExists
}

func (k Keeper) SetAdmin(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	newAdmin := customtypes.NFTAdmin{
		Address: addr.String(),
	}
	bz := k.cdc.MustMarshalBinaryBare(&newAdmin)
	store.Set(customtypes.AdminKey(addr), bz)
}