package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/types/simulation"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// RandomizedGenState creates a randomGenesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand
	gs := typesv1.GenesisState{}
	gs.StartingProposalId = uint64(simState.Rand.Intn(100))

	a := GenerateADepositParams(r)
	gs.DepositParams = &a
	b := GenerateAVotingParams(r)
	gs.VotingParams = &b
	tallyParams := GenerateTallyParams(r)
	gs.TallyParams = &tallyParams
	gs.CustomParams = &typesv1.CustomParams{
		CertifierUpdateSecurityVoteTally: &tallyParams,
		CertifierUpdateStakeVoteTally:    &tallyParams,
	}

	// For the shield module, locking period should be shorter than unbonding period.
	stakingGenStatebz := simState.GenState[stakingtypes.ModuleName]
	var stakingGenState stakingtypes.GenesisState
	simState.Cdc.MustUnmarshalJSON(stakingGenStatebz, &stakingGenState)
	ubdTime := stakingGenState.Params.UnbondingTime
	if *gs.VotingParams.VotingPeriod*2 >= ubdTime {
		votingPeriod := time.Duration(sim.RandIntBetween(r, int(ubdTime)/10, int(ubdTime)/2))
		gs.VotingParams.VotingPeriod = &votingPeriod
	}

	simState.GenState[govTypes.ModuleName] = simState.Cdc.MustMarshalJSON(&gs)
}

// GenerateADepositParams returns a DepositParams object with all of its fields randomized.
func GenerateADepositParams(r *rand.Rand) govtypesv1.DepositParams {
	minDeposit := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(sim.RandIntBetween(r, 1, 10))))
	maxDepositPeriod := sim.RandIntBetween(r, 1, 2*60*60*24*2)
	return govtypesv1.NewDepositParams(minDeposit, time.Duration(maxDepositPeriod)*time.Second)
}

// GenerateAVotingParams returns a VotingParams object with all of its fields randomized.
func GenerateAVotingParams(r *rand.Rand) govtypesv1.VotingParams {
	votingPeriod := sim.RandIntBetween(r, 1, 2*60*60*24*2)
	return govtypesv1.NewVotingParams(time.Duration(votingPeriod) * time.Second)
}

// GenerateTallyParams returns a TallyParams object with all of its fields randomized.
func GenerateTallyParams(r *rand.Rand) govtypesv1.TallyParams {
	quorum := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 334, 500)), 3)
	threshold := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 450, 550)), 3)
	veto := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 250, 334)), 3)
	return govtypesv1.NewTallyParams(quorum, threshold, veto)
}
