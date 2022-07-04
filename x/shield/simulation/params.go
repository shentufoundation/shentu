package simulation

import (
	"encoding/json"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// ParamChanges defines the parameters that can be modified by param change proposals on the simulation.
func ParamChanges(_ *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
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

		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyStakingShieldRate),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenShieldStakingRateParam(r))
				return string(bz)
			},
		),
	}
}
