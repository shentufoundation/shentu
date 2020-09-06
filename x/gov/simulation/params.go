package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

const (
	keyVotingParams  = "votingparams"
	keyDepositParams = "depositparams"
	keyTallyParams   = "tallyparams"
	subkeyQuorum     = "quorum"
	subkeyThreshold  = "threshold"
	subkeyVeto       = "veto"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	votingPeriod := time.Duration(simulation.RandIntBetween(r, 1, 2*60*60*24*2)) * time.Second
	depositPeriod := time.Duration(simulation.RandIntBetween(r, 1, 2*60*60*24*2)) * time.Second
	tallyParams := GenerateATallyParams(r)

	return []simulation.ParamChange{
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
				changes := []struct {
					key   string
					value sdk.Dec
				}{
					{subkeyQuorum, tallyParams.Quorum},
					{subkeyThreshold, tallyParams.Threshold},
					{subkeyVeto, tallyParams.Veto},
				}

				pc := make(map[string]string)
				numChanges := simulation.RandIntBetween(r, 1, len(changes))
				for i := 0; i < numChanges; i++ {
					c := changes[r.Intn(len(changes))]

					_, ok := pc[c.key]
					for ok {
						c := changes[r.Intn(len(changes))]
						_, ok = pc[c.key]
					}

					pc[c.key] = c.value.String()
				}

				bz, _ := json.Marshal(pc)
				return string(bz)
			},
		),
	}
}
