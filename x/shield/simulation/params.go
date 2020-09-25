package simulation

import (
	"encoding/json"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals on the simulation.
func ParamChanges(_ *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyPoolParams),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenPoolParams(r))
				return string(bz)
			},
		),

		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyClaimProposalParams),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenClaimProposalParams(r))
				return string(bz)
			},
		),
	}
}
