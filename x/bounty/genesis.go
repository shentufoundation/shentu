package bounty

import (
	"cosmossdk.io/collections"

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

	if err := k.TheoremID.Set(ctx, data.StartingTheoremId); err != nil {
		panic(err)
	}

	if err := k.Params.Set(ctx, *data.Params); err != nil {
		panic(err)
	}

	for _, theorem := range data.Theorems {
		switch theorem.Status {
		case types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD:
			err := k.ActiveTheoremsQueue.Set(ctx, collections.Join(*theorem.ProofEndTime, theorem.Id), theorem.Id)
			if err != nil {
				panic(err)
			}
		}
		err := k.SetTheorem(ctx, *theorem)
		if err != nil {
			panic(err)
		}
	}

	//for _, proof := range data.Proofs {
	//	switch proof.Status {
	//	case types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD:
	//		k.HashLockedProofsQueue.Set(ctx, collections.Join(*proof.SubmitTime))
	//	}
	//}

	//for _, grant := range data.Grants {
	//
	//}

}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	programs := k.GetAllPrograms(ctx)
	findings := k.GetAllFindings(ctx)

	return &types.GenesisState{
		Programs: programs,
		Findings: findings,
	}
}
