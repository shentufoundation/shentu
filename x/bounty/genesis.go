package bounty

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	k.SetNextProgramID(ctx, data.StartingProgramId)
	k.SetNextFindingID(ctx, data.StartingFindingId)
	// TODO Complete InitGenesis
}

// TODO Implement ExportGenesis
