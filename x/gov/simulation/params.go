package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	sim "github.com/cosmos/cosmos-sdk/types/simulation"

	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

const (
	keyVotingParams  = "votingparams"
	keyDepositParams = "depositparams"
	keyTallyParams   = "tallyparams"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []sim.ParamChange {
	votingPeriod := time.Duration(sim.RandIntBetween(r, 1, 2*60*60*24*2)) * time.Second
	depositPeriod := time.Duration(sim.RandIntBetween(r, 1, 2*60*60*24*2)) * time.Second

	return []sim.ParamChange{
		simulation.NewSimParamChange(govTypes.ModuleName, keyVotingParams,
			func(r *rand.Rand) string {
				return fmt.Sprintf(`{"voting_period": "%d"}`, votingPeriod)
			},
		),
		simulation.NewSimParamChange(govTypes.ModuleName, keyDepositParams,
			func(r *rand.Rand) string {
				return fmt.Sprintf(`{"max_deposit_period": "%d"}`, depositPeriod)
			},
		),
		simulation.NewSimParamChange(govTypes.ModuleName, keyTallyParams,
			func(r *rand.Rand) string {
				bz, _ := json.Marshal(GenerateTallyParams(r))
				return string(bz)
			},
		),
	}
}
