package bounty

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	for _, finding := range data.Findings {
		k.SetFinding(ctx, finding)
		if err := k.AppendFidToFidList(ctx, finding.ProgramId, finding.FindingId); err != nil {
			panic(err)
		}
	}

	for _, program := range data.Programs {
		k.SetProgram(ctx, program)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	programs := k.GetAllPrograms(ctx)
	findings := k.GetAllFindings(ctx)

	return &types.GenesisState{
		Programs: programs,
		Findings: findings,
	}
}
