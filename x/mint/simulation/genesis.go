package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	cosmosMint "github.com/cosmos/cosmos-sdk/x/mint"
	mintSim "github.com/cosmos/cosmos-sdk/x/mint/simulation"

	"github.com/certikfoundation/shentu/common"
)

// RandomizedGenState generates a random GenesisState for mint.
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintSim.Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = mintSim.GenInflation(r) },
	)

	// params
	var inflationRateChange sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintSim.InflationRateChange, &inflationRateChange, simState.Rand,
		func(r *rand.Rand) { inflationRateChange = mintSim.GenInflationRateChange(r) },
	)

	var inflationMax sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintSim.InflationMax, &inflationMax, simState.Rand,
		func(r *rand.Rand) { inflationMax = mintSim.GenInflationMax(r) },
	)

	var inflationMin sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintSim.InflationMin, &inflationMin, simState.Rand,
		func(r *rand.Rand) { inflationMin = mintSim.GenInflationMin(r) },
	)

	var goalBonded sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintSim.GoalBonded, &goalBonded, simState.Rand,
		func(r *rand.Rand) { goalBonded = mintSim.GenGoalBonded(r) },
	)

	mintDenom := common.MicroCTKDenom
	blocksPerYear := common.BlocksPerYear
	params := cosmosMint.NewParams(mintDenom, inflationRateChange, inflationMax, inflationMin, goalBonded, blocksPerYear)

	mintGenesis := cosmosMint.NewGenesisState(cosmosMint.InitialMinter(inflation), params)

	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, mintGenesis))
	simState.GenState[cosmosMint.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
