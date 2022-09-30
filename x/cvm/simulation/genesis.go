package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/shentufoundation/shentu/v2/x/cvm/types"
)

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	gs := types.GenesisState{}

	gs.GasRate = 1

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&gs)
}
