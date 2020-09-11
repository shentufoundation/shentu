package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

const (
	lockedInBlocksKey     = "locked_in_blocks"
	minimumCollateralKey  = "minimum_collateral"
	expirationDurationKey = "expiration_duration"
	aggregationWindowKey  = "aggregation_window"
	aggregationResultKey  = "aggregation_result"
	thresholdScoreKey     = "threshold_score"
	epsilon1Key           = "epsilon_1"
	epsilon2Key           = "epsilon_2"
)

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	var lockedInBlocks int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, lockedInBlocksKey, &lockedInBlocks, simState.Rand,
		func(r *rand.Rand) {
			lockedInBlocks = r.Int63n(60)
		})

	var minimumCollateral int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, minimumCollateralKey, &minimumCollateral, simState.Rand,
		func(r *rand.Rand) {
			minimumCollateral = r.Int63n(100000)
		})

	var expirationDuration time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, expirationDurationKey, &expirationDuration, simState.Rand,
		func(r *rand.Rand) {
			expirationDuration = time.Duration(r.Int63n(1000 * 1000 * 1000 * 60 * 60 * 48))
		})

	var aggregationWindow int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, aggregationWindowKey, &aggregationWindow, simState.Rand,
		func(r *rand.Rand) {
			aggregationWindow = r.Int63n(40)
		})

	var aggregationResult sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, aggregationResultKey, &aggregationResult, simState.Rand,
		func(r *rand.Rand) {
			aggregationResult = sdk.NewInt(r.Int63n(3))
		})

	var thresholdScore sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, thresholdScoreKey, &thresholdScore, simState.Rand,
		func(r *rand.Rand) {
			thresholdScore = sdk.NewInt(r.Int63n(257))
		})

	var epsilon1 sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, epsilon1Key, &epsilon1, simState.Rand,
		func(r *rand.Rand) {
			epsilon1 = sdk.NewInt(r.Int63n(3))
		})

	var epsilon2 sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, epsilon2Key, &epsilon2, simState.Rand,
		func(r *rand.Rand) {
			epsilon2 = sdk.NewInt(r.Int63n(201))
		})

	gs := types.NewGenesisState(
		nil,
		nil,
		types.NewLockedPoolParams(lockedInBlocks, minimumCollateral),
		types.NewTaskParams(expirationDuration, aggregationWindow, aggregationResult, thresholdScore, epsilon1, epsilon2),
		nil,
		nil,
	)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}
