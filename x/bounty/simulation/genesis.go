package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	gs := types.DefaultGenesisState()
	gs.StartingProgramId = uint64(simState.Rand.Int63n(100))
	gs.StartingFindingId = uint64(simState.Rand.Int63n(100))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}
