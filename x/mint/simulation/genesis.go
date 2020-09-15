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
		func(r *rand.Rand) { inflation = sdk.NewDecWithPrec(7, 2) },
	)

	// params
	var inflationRateChange sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintSim.InflationRateChange, &inflationRateChange, simState.Rand,
		func(r *rand.Rand) { inflationRateChange = sdk.NewDecWithPrec(10, 2) },
	)

	var inflationMax sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintSim.InflationMax, &inflationMax, simState.Rand,
		func(r *rand.Rand) { inflationMax = sdk.NewDecWithPrec(14, 2) },
	)

	var inflationMin sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintSim.InflationMin, &inflationMin, simState.Rand,
		func(r *rand.Rand) { inflationMin = sdk.NewDecWithPrec(4, 2) },
	)

	var goalBonded sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintSim.GoalBonded, &goalBonded, simState.Rand,
		func(r *rand.Rand) { goalBonded = sdk.NewDecWithPrec(67, 2) },
	)

	mintDenom := common.MicroCTKDenom
	blocksPerYear := common.BlocksPerYear
	params := cosmosMint.NewParams(mintDenom, inflationRateChange, inflationMax, inflationMin, goalBonded, blocksPerYear)

	mintGenesis := cosmosMint.NewGenesisState(cosmosMint.InitialMinter(inflation), params)

	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, mintGenesis))
	simState.GenState[cosmosMint.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
