package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govsim "github.com/cosmos/cosmos-sdk/x/gov/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {
	startingProposalID := uint64(simState.Rand.Intn(100))

	var minDeposit sdk.Coins
	simState.AppParams.GetOrGenerate(
		simState.Cdc, govsim.DepositParamsMinDeposit, &minDeposit, simState.Rand,
		func(r *rand.Rand) { minDeposit = govsim.GenDepositParamsMinDeposit(r) },
	)

	var depositPeriod time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, govsim.DepositParamsDepositPeriod, &depositPeriod, simState.Rand,
		func(r *rand.Rand) { depositPeriod = govsim.GenDepositParamsDepositPeriod(r) },
	)

	var minInitialDepositRatio sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, govsim.DepositMinInitialRatio, &minInitialDepositRatio, simState.Rand,
		func(r *rand.Rand) { minInitialDepositRatio = govsim.GenDepositMinInitialDepositRatio(r) },
	)

	var votingPeriod time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, govsim.VotingParamsVotingPeriod, &votingPeriod, simState.Rand,
		func(r *rand.Rand) { votingPeriod = govsim.GenVotingParamsVotingPeriod(r) },
	)

	var quorum sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, govsim.TallyParamsQuorum, &quorum, simState.Rand,
		func(r *rand.Rand) { quorum = govsim.GenTallyParamsQuorum(r) },
	)

	var threshold sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, govsim.TallyParamsThreshold, &threshold, simState.Rand,
		func(r *rand.Rand) { threshold = govsim.GenTallyParamsThreshold(r) },
	)

	var veto sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, govsim.TallyParamsVeto, &veto, simState.Rand,
		func(r *rand.Rand) { veto = govsim.GenTallyParamsVeto(r) },
	)

	tally := &govtypesv1.TallyParams{
		Quorum:        quorum.String(),
		Threshold:     threshold.String(),
		VetoThreshold: veto.String(),
	}

	customParam := typesv1.CustomParams{
		CertifierUpdateSecurityVoteTally: tally,
		CertifierUpdateStakeVoteTally:    tally,
	}

	govGenesis := typesv1.NewGenesisState(
		startingProposalID,
		govtypesv1.NewParams(minDeposit, depositPeriod, votingPeriod, quorum.String(), threshold.String(), veto.String(), minInitialDepositRatio.String(), simState.Rand.Intn(2) == 0, simState.Rand.Intn(2) == 0, simState.Rand.Intn(2) == 0),
		customParam,
	)

	bz, err := json.MarshalIndent(&govGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated governance parameters:\n%s\n", bz)
	simState.GenState[govtypes.ModuleName] = simState.Cdc.MustMarshalJSON(govGenesis)
}

