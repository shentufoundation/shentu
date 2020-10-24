package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/gov/internal/types"
)

// RandomizedGenState creates a randomGenesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand
	gs := types.GenesisState{}
	gs.StartingProposalID = uint64(simState.Rand.Intn(100))

	gs.DepositParams = GenerateADepositParams(r)
	gs.VotingParams = GenerateAVotingParams(r)
	gs.TallyParams = GenerateTallyParams(r)

	// For the shield module, locking period should be shorter than unbonding period.
	stakingGenStatebz := simState.GenState[staking.ModuleName]
	var stakingGenState staking.GenesisState
	staking.ModuleCdc.MustUnmarshalJSON(stakingGenStatebz, &stakingGenState)
	ubdTime := stakingGenState.Params.UnbondingTime
	if 2*gs.VotingParams.VotingPeriod >= ubdTime {
		gs.VotingParams.VotingPeriod = time.Duration(simulation.RandIntBetween(r, int(ubdTime)/10, int(ubdTime)/2))
	}

	simState.GenState[govTypes.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}

// GenerateADepositParams returns a DepositParams object with all of its fields randomized.
func GenerateADepositParams(r *rand.Rand) types.DepositParams {
	minInitialDeposit := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simulation.RandIntBetween(r, 1, 1e2))))
	minDeposit := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simulation.RandIntBetween(r, 1, 1e3))))
	maxDepositPeriod := simulation.RandIntBetween(r, 1, 2*60*60*24*2)
	return types.NewDepositParams(minInitialDeposit, minDeposit, time.Duration(maxDepositPeriod)*time.Second)
}

// GenerateAVotingParams returns a VotingParams object with all of its fields randomized.
func GenerateAVotingParams(r *rand.Rand) govTypes.VotingParams {
	votingPeriod := simulation.RandIntBetween(r, 1, 2*60*60*24*2)
	return govTypes.NewVotingParams(time.Duration(votingPeriod) * time.Second)
}

//GenerateTallyParams returns a TallyParams object with all of its fields randomized.
func GenerateTallyParams(r *rand.Rand) types.TallyParams {
	return types.TallyParams{
		DefaultTally:                     GenerateATallyParams(r),
		CertifierUpdateSecurityVoteTally: GenerateATallyParams(r),
		CertifierUpdateStakeVoteTally:    GenerateATallyParams(r),
	}
}

// GenerateATallyParams returns a TallyParams object with all of its fields randomized.
func GenerateATallyParams(r *rand.Rand) govTypes.TallyParams {
	quorum := sdk.NewDecWithPrec(int64(simulation.RandIntBetween(r, 334, 500)), 3)
	threshold := sdk.NewDecWithPrec(int64(simulation.RandIntBetween(r, 450, 550)), 3)
	veto := sdk.NewDecWithPrec(int64(simulation.RandIntBetween(r, 250, 334)), 3)
	return govTypes.NewTallyParams(quorum, threshold, veto)
}
