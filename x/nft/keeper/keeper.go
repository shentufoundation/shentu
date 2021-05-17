package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irismod/modules/nft/keeper"

	certkeeper "github.com/certikfoundation/shentu/x/cert/keeper"
)

type Keeper struct {
	keeper.Keeper
	certKeeper certkeeper.Keeper
}

func (k Keeper) IssueNFTAdmin(ctx sdk.Context, addr sdk.AccAddress) bool {
	certifiers := k.certKeeper.GetAllCertifiers(ctx)

}