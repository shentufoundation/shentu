// Package simulation glues the shentu gov module into the cosmos-sdk
// simulation framework. Most of the knobs (tally params, deposit
// periods, etc.) are identical to the upstream gov module, so this
// package reuses upstream's Gen* helpers directly. The shentu-specific
// piece is the shentu GenesisState wrapper: upstream's genesis only
// carries the v1.Params, while shentu's also carries CustomParams —
// specifically CertifierUpdateSecurityVoteTally for the cert-update
// security round. Under the v7 single-round model there is no separate
// stake round, so CertifierUpdateStakeVoteTally is intentionally not
// populated.
package simulation

import (
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govsim "github.com/cosmos/cosmos-sdk/x/gov/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// RandomizedGenState generates a random GenesisState for the shentu gov
// module. It reuses upstream's Gen* helpers for every param the two
// modules share, then wraps the result in shentu's GenesisState with a
// single CertifierUpdateSecurityVoteTally derived from the same
// randomized quorum/threshold/veto values.
func RandomizedGenState(simState *module.SimulationState) {
	startingProposalID := uint64(simState.Rand.Intn(100))

	var minDeposit sdk.Coins
	simState.AppParams.GetOrGenerate(govsim.MinDeposit, &minDeposit, simState.Rand,
		func(r *rand.Rand) { minDeposit = govsim.GenMinDeposit(r, simState.BondDenom) })

	var expeditedMinDeposit sdk.Coins
	simState.AppParams.GetOrGenerate(govsim.ExpeditedMinDeposit, &expeditedMinDeposit, simState.Rand,
		func(r *rand.Rand) { expeditedMinDeposit = govsim.GenExpeditedMinDeposit(r, simState.BondDenom) })

	var depositPeriod time.Duration
	simState.AppParams.GetOrGenerate(govsim.DepositPeriod, &depositPeriod, simState.Rand,
		func(r *rand.Rand) { depositPeriod = govsim.GenDepositPeriod(r) })

	var minInitialDepositRatio sdkmath.LegacyDec
	simState.AppParams.GetOrGenerate(govsim.MinInitialRatio, &minInitialDepositRatio, simState.Rand,
		func(r *rand.Rand) { minInitialDepositRatio = govsim.GenDepositMinInitialDepositRatio(r) })

	var proposalCancelRate sdkmath.LegacyDec
	simState.AppParams.GetOrGenerate(govsim.ProposalCancelRate, &proposalCancelRate, simState.Rand,
		func(r *rand.Rand) { proposalCancelRate = govsim.GenProposalCancelRate(r) })

	var votingPeriod time.Duration
	simState.AppParams.GetOrGenerate(govsim.VotingPeriod, &votingPeriod, simState.Rand,
		func(r *rand.Rand) { votingPeriod = govsim.GenVotingPeriod(r) })

	var expeditedVotingPeriod time.Duration
	simState.AppParams.GetOrGenerate(govsim.ExpeditedVotingPeriod, &expeditedVotingPeriod, simState.Rand,
		func(r *rand.Rand) { expeditedVotingPeriod = govsim.GenExpeditedVotingPeriod(r) })

	var quorum sdkmath.LegacyDec
	simState.AppParams.GetOrGenerate(govsim.Quorum, &quorum, simState.Rand,
		func(r *rand.Rand) { quorum = govsim.GenQuorum(r) })

	var threshold sdkmath.LegacyDec
	simState.AppParams.GetOrGenerate(govsim.Threshold, &threshold, simState.Rand,
		func(r *rand.Rand) { threshold = govsim.GenThreshold(r) })

	var expeditedThreshold sdkmath.LegacyDec
	simState.AppParams.GetOrGenerate(govsim.ExpeditedThreshold, &expeditedThreshold, simState.Rand,
		func(r *rand.Rand) { expeditedThreshold = govsim.GenExpeditedThreshold(r) })

	var veto sdkmath.LegacyDec
	simState.AppParams.GetOrGenerate(govsim.Veto, &veto, simState.Rand,
		func(r *rand.Rand) { veto = govsim.GenVeto(r) })

	var minDepositRatio sdkmath.LegacyDec
	simState.AppParams.GetOrGenerate(govsim.MinDepositRatio, &minDepositRatio, simState.Rand,
		func(r *rand.Rand) { minDepositRatio = govsim.GenMinDepositRatio(r) })

	params := govtypesv1.NewParams(
		minDeposit,
		expeditedMinDeposit,
		depositPeriod,
		votingPeriod,
		expeditedVotingPeriod,
		quorum.String(),
		threshold.String(),
		expeditedThreshold.String(),
		veto.String(),
		minInitialDepositRatio.String(),
		proposalCancelRate.String(),
		"",
		simState.Rand.Intn(2) == 0,
		simState.Rand.Intn(2) == 0,
		simState.Rand.Intn(2) == 0,
		minDepositRatio.String(),
	)

	// Certifier security round reuses the same quorum/threshold/veto
	// values as the main stake round — it's a separate electorate but
	// the same dec shape. Only CertifierUpdateSecurityVoteTally is set;
	// the stake round for cert-update was removed in v7.
	certifierTally := govtypesv1.NewTallyParams(quorum.String(), threshold.String(), veto.String())
	customParams := typesv1.CustomParams{
		CertifierUpdateSecurityVoteTally: &certifierTally,
	}

	govGenesis := typesv1.NewGenesisState(startingProposalID, params, customParams)
	simState.GenState[govtypes.ModuleName] = simState.Cdc.MustMarshalJSON(govGenesis)
}
