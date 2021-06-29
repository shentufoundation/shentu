package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/nft"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"

	"github.com/certikfoundation/shentu/x/nft/keeper"
	"github.com/certikfoundation/shentu/x/nft/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	baseGenesis := nfttypes.GenesisState{
		Collections: data.Collections,
	}
	nft.InitGenesis(ctx, k.Keeper, baseGenesis)

	for _, a := range data.Admin {
		addr, err := sdk.AccAddressFromBech32(a.Address)
		if err != nil {
			panic("Error while translating NFT Admin address: " + err.Error())
		}
		k.SetAdmin(ctx, addr)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	collections := k.GetCollections(ctx)
	admin := k.GetAdmins(ctx)

	return &types.GenesisState{
		Collections: collections,
		Admin:       admin,
	}
}
