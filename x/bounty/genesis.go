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

	for _, program := range data.Programs {
		k.SetProgram(ctx, program)
	}

	findingIdMap := make(map[uint64][]uint64)
	for _, finding := range data.Findings {
		k.SetFinding(ctx, finding)

		findingList, ok := findingIdMap[finding.ProgramId]
		if !ok {
			findingList = []uint64{finding.FindingId}
		} else {
			findingList = append(findingList, finding.FindingId)
		}
		findingIdMap[finding.ProgramId] = findingList
	}

	for programId, findingIdList := range findingIdMap {
		k.SetPidFindingIDList(ctx, programId, findingIdList)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var programs []types.Program
	var findings []types.Finding

	maxFindingId := k.GetNextFindingID(ctx)
	maxProgramId := k.GetNextProgramID(ctx)
	for programId := uint64(1); programId < maxProgramId; programId++ {
		program, ok := k.GetProgram(ctx, programId)
		if ok {
			programs = append(programs, program)

			findingIds, err := k.GetPidFindingIDList(ctx, program.ProgramId)
			if err == nil {
				for _, fid := range findingIds {
					finding, ok := k.GetFinding(ctx, fid)
					if ok {
						findings = append(findings, finding)
					}
				}
			}
		}
	}

	return &types.GenesisState{
		StartingFindingId: maxFindingId,
		StartingProgramId: maxProgramId,
		Programs:          programs,
		Findings:          findings,
	}
}
