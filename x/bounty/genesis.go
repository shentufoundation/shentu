package bounty

import (
	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, ak types.AccountKeeper, k keeper.Keeper, data *types.GenesisState) error {
	// validate genesis state
	if err := types.ValidateGenesis(data); err != nil {
		return err
	}

	// initialize programs
	for _, program := range data.Programs {
		if err := k.Programs.Set(ctx, program.ProgramId, *program); err != nil {
			return err
		}
	}

	// initialize findings
	for _, finding := range data.Findings {
		if err := k.Findings.Set(ctx, finding.FindingId, *finding); err != nil {
			return err
		}
		if err := k.ProgramFindings.Set(ctx, collections.Join(finding.ProgramId, finding.FindingId)); err != nil {
			return err
		}
	}

	// initialize theorem ID
	if err := k.TheoremID.Set(ctx, data.StartingTheoremId); err != nil {
		return err
	}

	// initialize theorems
	for _, theorem := range data.Theorems {
		switch theorem.Status {
		case types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD:
			if err := k.ActiveTheoremsQueue.Set(ctx, collections.Join(*theorem.EndTime, theorem.Id), theorem.Id); err != nil {
				return err
			}
		}
		if err := k.Theorems.Set(ctx, theorem.Id, *theorem); err != nil {
			return err
		}
	}

	// initialize grants
	for _, grant := range data.Grants {
		addr, err := ak.AddressCodec().StringToBytes(grant.Grantor)
		if err != nil {
			return err
		}
		if err := k.Grants.Set(ctx, collections.Join(grant.TheoremId, sdk.AccAddress(addr)), *grant); err != nil {
			return err
		}
	}

	// initialize deposits
	for _, deposit := range data.Deposits {
		addr, err := ak.AddressCodec().StringToBytes(deposit.Depositor)
		if err != nil {
			return err
		}
		if err := k.Deposits.Set(ctx, collections.Join(deposit.ProofId, sdk.AccAddress(addr)), *deposit); err != nil {
			return err
		}
	}

	// initialize rewards
	for _, reward := range data.Rewards {
		addr, err := ak.AddressCodec().StringToBytes(reward.Address)
		if err != nil {
			return err
		}
		if err := k.Rewards.Set(ctx, addr, *reward); err != nil {
			return err
		}
	}

	// initialize params
	if err := k.Params.Set(ctx, *data.Params); err != nil {
		return err
	}

	// initialize proofs
	for _, proof := range data.Proofs {
		switch proof.Status {
		case types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD:
			if err := k.ActiveProofsQueue.Set(ctx, collections.Join(*proof.SubmitTime, proof.Id), *proof); err != nil {
				return err
			}
		}
		if err := k.Proofs.Set(ctx, proof.Id, *proof); err != nil {
			return err
		}

		if err := k.ProofsByTheorem.Set(ctx, collections.Join(proof.TheoremId, proof.Id), []byte{}); err != nil {
			return err
		}
	}

	return nil
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var (
		programs []*types.Program
		findings []*types.Finding
		theorems []*types.Theorem
		proofs   []*types.Proof
		grants   []*types.Grant
		rewards  []*types.Reward
		deposits []*types.Deposit
	)

	err := k.Programs.Walk(ctx, nil, func(_ string, value types.Program) (stop bool, err error) {
		programs = append(programs, &value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	err = k.Findings.Walk(ctx, nil, func(_ string, value types.Finding) (stop bool, err error) {
		findings = append(findings, &value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	err = k.Theorems.Walk(ctx, nil, func(_ uint64, value types.Theorem) (stop bool, err error) {
		theorems = append(theorems, &value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	err = k.Proofs.Walk(ctx, nil, func(_ string, value types.Proof) (stop bool, err error) {
		proofs = append(proofs, &value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	err = k.Grants.Walk(ctx, nil, func(_ collections.Pair[uint64, sdk.AccAddress], value types.Grant) (stop bool, err error) {
		grants = append(grants, &value)
		return false, nil
	})
	if err != nil {
		return nil
	}
	if err != nil {
		panic(err)
	}

	err = k.Rewards.Walk(ctx, nil, func(key sdk.AccAddress, value types.Reward) (stop bool, err error) {
		rewards = append(rewards, &value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	err = k.Deposits.Walk(ctx, nil, func(_ collections.Pair[string, sdk.AccAddress], value types.Deposit) (stop bool, err error) {
		deposits = append(deposits, &value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	startingTheoremId, err := k.TheoremID.Peek(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Programs:          programs,
		Findings:          findings,
		StartingTheoremId: startingTheoremId,
		Theorems:          theorems,
		Proofs:            proofs,
		Grants:            grants,
		Rewards:           rewards,
		Deposits:          deposits,
		Params:            &params,
	}
}
