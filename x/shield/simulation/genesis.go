package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	// TODO set random genesis state
	gs := types.DefaultGenesisState()
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}

// GenPoolParams returns a randomized PoolParams object.
func GenPoolParams(r *rand.Rand) types.PoolParams {
	// TODO set random params
	return types.DefaultPoolParams()
}

// GenClaimProposalParams returns a randomized ClaimProposalParams object.
func GenClaimProposalParams(r *rand.Rand) types.ClaimProposalParams {
	// TODO set random params
	return types.DefaultClaimProposalParams()
}
