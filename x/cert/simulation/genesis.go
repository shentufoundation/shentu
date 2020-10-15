package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand
	gs := types.GenesisState{}

	numOfCertifiers := r.Intn(10)
	for i := 0; i < numOfCertifiers; i++ {
		acc, _ := simulation.RandomAcc(r, simState.Accounts)
		certifier := types.NewCertifier(acc.Address, "", nil, "")
		gs.Certifiers = append(gs.Certifiers, certifier)
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}
