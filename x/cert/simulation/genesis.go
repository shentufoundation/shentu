package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	gs := &types.GenesisState{}

	for _, acc := range simState.Accounts {
		if simState.Rand.Intn(100) < 10 {
			certifier := types.NewCertifier(acc.Address, "", nil, "")
			gs.Certifiers = append(gs.Certifiers, certifier)
		}
	}

	gs.NextCertificateId = 1

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}
