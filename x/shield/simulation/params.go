package simulation

import (
	"encoding/json"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals on the simulation.
func ParamChanges(_ *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(v1beta1.ParamStoreKeyPoolParams),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenPoolParams(r))
				return string(bz)
			},
		),

		simulation.NewSimParamChange(types.ModuleName, string(v1beta1.ParamStoreKeyClaimProposalParams),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenClaimProposalParams(r))
				return string(bz)
			},
		),

		simulation.NewSimParamChange(types.ModuleName, string(v1beta1.ParamStoreKeyStakingShieldRate),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenShieldStakingRateParam(r))
				return string(bz)
			},
		),

		simulation.NewSimParamChange(types.ModuleName, string(v1beta1.ParamStoreKeyBlockRewardParams),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenBlockRewardParams(r))
				return string(bz)
			},
		),
	}
}
