package bounty

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data *types.GenesisState) {
	k.SetNextProgramID(ctx, data.StartingProgramId)

	for _, program := range data.Programs {
		k.SetProgram(ctx, program)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	startingProgramID, _ := k.GetNextProgramID(ctx)
	programs := k.GetPrograms(ctx)

	return &types.GenesisState{
		StartingProgramId: startingProgramID,
		Programs:          programs,
	}
}
