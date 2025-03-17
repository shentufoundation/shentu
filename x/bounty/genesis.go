package bounty

import (
	"cosmossdk.io/collections"

	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	for _, program := range data.Programs {
		err := k.Programs.Set(ctx, program.ProgramId, program)
		if err != nil {
			panic(err)
		}
	}

	for _, finding := range data.Findings {
		err := k.Findings.Set(ctx, finding.FindingId, finding)
		if err != nil {
			panic(err)
		}
		if err = k.ProgramFindings.Set(ctx, collections.Join(finding.ProgramId, finding.FindingId)); err != nil {
			panic(err)
		}
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
			err := k.ActiveTheoremsQueue.Set(ctx, collections.Join(*theorem.EndTime, theorem.Id), theorem.Id)
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
	var programs types.Programs
	var findings types.Findings

	err := k.Programs.Walk(ctx, nil, func(_ string, value types.Program) (stop bool, err error) {
		programs = append(programs, value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	err = k.Findings.Walk(ctx, nil, func(_ string, value types.Finding) (stop bool, err error) {
		findings = append(findings, value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Programs: programs,
		Findings: findings,
	}
}
