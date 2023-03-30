package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	var poolParams types.LockedPoolParams
	simState.AppParams.GetOrGenerate(
		simState.Cdc, string(types.ParamsStoreKeyPoolParams), &poolParams, simState.Rand,
		func(r *rand.Rand) {
			poolParams = GenPoolParams(r)
		})

	var taskParams types.TaskParams
	simState.AppParams.GetOrGenerate(
		simState.Cdc, string(types.ParamsStoreKeyTaskParams), &taskParams, simState.Rand,
		func(r *rand.Rand) {
			taskParams = GenTaskParams(r)
		})

	gs := types.NewGenesisState(
		nil,
		nil,
		poolParams,
		taskParams,
		nil,
		nil,
		nil,
	)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&gs)
}

// GenPoolParams returns a randomized LockedPoolParams object.
func GenPoolParams(r *rand.Rand) types.LockedPoolParams {
	return types.LockedPoolParams{
		LockedInBlocks:    r.Int63n(10),
		MinimumCollateral: r.Int63n(10000),
	}
}

// GenTaskParams returns a randomized TaskParams object.
func GenTaskParams(r *rand.Rand) types.TaskParams {
	return types.TaskParams{
		ExpirationDuration: time.Duration(48) * time.Hour,
		AggregationWindow:  r.Int63n(40),
		AggregationResult:  sdk.NewInt(r.Int63n(100)),
		ThresholdScore:     sdk.NewInt(r.Int63n(100)),
		Epsilon1:           sdk.NewInt(r.Int63n(10) + 1),
		Epsilon2:           sdk.NewInt(r.Int63n(10) + 90),
		ShortcutQuorum:     sdk.NewDecWithPrec(int64((r.Float64()+r.Float64())/2*10000)+1, 4),
	}
}
