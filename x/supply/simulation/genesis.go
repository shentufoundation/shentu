package simulation

// DONTCOVER

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	cosmosSupply "github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/certikfoundation/shentu/common"
)

// RandomizedGenState generates a random GenesisState for supply.
func RandomizedGenState(simState *module.SimulationState) {
	numAccs := int64(len(simState.Accounts))
	totalSupply := sdk.NewInt(simState.InitialStake * (numAccs + simState.NumBonded))
	supplyGenesis := cosmosSupply.NewGenesisState(sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, totalSupply)))

	fmt.Printf("Generated supply parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, supplyGenesis))
	simState.GenState[cosmosSupply.ModuleName] = simState.Cdc.MustMarshalJSON(supplyGenesis)
}
