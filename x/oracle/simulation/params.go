package simulation

import (
	"encoding/json"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals on the simulation.
func ParamChanges(_ *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamsStoreKeyPoolParams),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenPoolParams(r))
				return string(bz)
			},
		),

		simulation.NewSimParamChange(types.ModuleName, string(types.ParamsStoreKeyTaskParams),
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenTaskParams(r))
				return string(bz)
			},
		),
	}
}
