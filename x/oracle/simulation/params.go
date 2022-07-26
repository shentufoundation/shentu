package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/v2/x/oracle/types"
)

const (
	subkeyLockedInBlocks     = "locked_in_blocks"
	subkeyMinimumCollateral  = "minimum_collateral"
	subkeyExpirationDuration = "expiration_duration"
	subkeyAggregationWindow  = "aggregation_window"
	subkeyAggregationResult  = "aggregation_result"
	subkeyThresholdScore     = "threshold_score"
	subkeyEpsilon1           = "epsilon1"
	subkeyEpsilon2           = "epsilon2"
)

// ParamChanges defines the parameters that can be modified by param change proposals on the simulation.
func ParamChanges(_ *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamsStoreKeyPoolParams),
			func(r *rand.Rand) string {
				pp := GenPoolParams(r)
				changes := []struct {
					key   string
					value int64
				}{
					{subkeyLockedInBlocks, pp.LockedInBlocks},
					{subkeyMinimumCollateral, pp.MinimumCollateral},
				}

				pc := make(map[string]string)
				numChanges := len(changes)
				for i := 0; i < numChanges; i++ {
					c := changes[i]
					pc[c.key] = strconv.FormatInt(c.value, 10)
				}
				bz, _ := json.Marshal(pc)
				return string(bz)
			},
		),

		simulation.NewSimParamChange(types.ModuleName, string(types.ParamsStoreKeyTaskParams),
			func(r *rand.Rand) string {
				tp := GenTaskParams(r)
				changes := []struct {
					key   string
					value int64
				}{
					{subkeyAggregationWindow, tp.AggregationWindow},
					{subkeyAggregationResult, tp.AggregationResult.Int64()},
					{subkeyThresholdScore, tp.ThresholdScore.Int64()},
					{subkeyEpsilon1, tp.Epsilon1.Int64()},
					{subkeyEpsilon2, tp.Epsilon2.Int64()},
				}

				pc := make(map[string]string)
				numChanges := len(changes)
				for i := 0; i < numChanges; i++ {
					c := changes[i]
					pc[c.key] = strconv.FormatInt(c.value, 10)
				}
				pc[subkeyExpirationDuration] = fmt.Sprintf("%d", tp.ExpirationDuration)
				bz, _ := json.Marshal(pc)
				return string(bz)
			},
		),
	}
}
