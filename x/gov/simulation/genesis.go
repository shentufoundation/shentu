package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/types/simulation"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

// RandomizedGenState creates a randomGenesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand
	gs := types.GenesisState{}
	gs.StartingProposalId = uint64(simState.Rand.Intn(100))

	gs.DepositParams = GenerateADepositParams(r)
	gs.VotingParams = GenerateAVotingParams(r)
	tallyParams := GenerateTallyParams(r)
	gs.TallyParams = tallyParams
	gs.CustomParams = types.CustomParams{
		CertifierUpdateSecurityVoteTally: &tallyParams,
		CertifierUpdateStakeVoteTally:    &tallyParams,
	}

	// For the shield module, locking period should be shorter than unbonding period.
	stakingGenStatebz := simState.GenState[stakingtypes.ModuleName]
	var stakingGenState stakingtypes.GenesisState
	simState.Cdc.MustUnmarshalJSON(stakingGenStatebz, &stakingGenState)
	ubdTime := stakingGenState.Params.UnbondingTime
	if 2*gs.VotingParams.VotingPeriod >= ubdTime {
		gs.VotingParams.VotingPeriod = time.Duration(sim.RandIntBetween(r, int(ubdTime)/10, int(ubdTime)/2))
	}

	simState.GenState[govTypes.ModuleName] = simState.Cdc.MustMarshalJSON(&gs)
}

// GenerateADepositParams returns a DepositParams object with all of its fields randomized.
func GenerateADepositParams(r *rand.Rand) govTypes.DepositParams {
	//minInitialDeposit := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(sim.RandIntBetween(r, 1, 1e2))))
	minDeposit := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(sim.RandIntBetween(r, 1, 1e3))))
	maxDepositPeriod := sim.RandIntBetween(r, 1, 2*60*60*24*2)
	return govTypes.NewDepositParams(minDeposit, time.Duration(maxDepositPeriod)*time.Second)
}

// GenerateAVotingParams returns a VotingParams object with all of its fields randomized.
func GenerateAVotingParams(r *rand.Rand) govTypes.VotingParams {
	votingPeriod := sim.RandIntBetween(r, 1, 2*60*60*24*2)
	return govTypes.NewVotingParams(time.Duration(votingPeriod) * time.Second)
}

// GenerateTallyParams returns a TallyParams object with all of its fields randomized.
func GenerateTallyParams(r *rand.Rand) govTypes.TallyParams {
	quorum := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 334, 500)), 3)
	threshold := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 450, 550)), 3)
	veto := sdk.NewDecWithPrec(int64(sim.RandIntBetween(r, 250, 334)), 3)
	return govTypes.NewTallyParams(quorum, threshold, veto)
}
