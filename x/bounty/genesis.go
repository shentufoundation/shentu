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

	findingIDMap := make(map[uint64][]uint64)
	for _, finding := range data.Findings {
		k.SetFinding(ctx, finding)

		findingList, ok := findingIDMap[finding.ProgramId]
		if !ok {
			findingList = []uint64{finding.FindingId}
		} else {
			findingList = append(findingList, finding.FindingId)
		}
		findingIDMap[finding.ProgramId] = findingList
	}

	for _, program := range data.Programs {
		k.SetProgram(ctx, program)
		findingIDList, ok := findingIDMap[program.ProgramId]
		if ok {
			k.SetPidFindingIDList(ctx, program.ProgramId, findingIDList)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var programs []types.Program
	var findings []types.Finding

	maxFindingID := k.GetNextFindingID(ctx)
	maxProgramID := k.GetNextProgramID(ctx)
	for programID := uint64(1); programID < maxProgramID; programID++ {
		program, ok := k.GetProgram(ctx, programID)
		if ok {
			programs = append(programs, program)

			findingIDs, err := k.GetPidFindingIDList(ctx, program.ProgramId)
			if err == nil {
				for _, fid := range findingIDs {
					finding, ok := k.GetFinding(ctx, fid)
					if ok {
						findings = append(findings, finding)
					}
				}
			}
		}
	}

	return &types.GenesisState{
		StartingFindingId: maxFindingID,
		StartingProgramId: maxProgramID,
		Programs:          programs,
		Findings:          findings,
	}
}
