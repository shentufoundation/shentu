package shield

import (
	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/shield/keeper"
	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

// InitGenesis initialize store values with genesis states.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	err := k.SetRemainingServiceFees(ctx, data.RemainingServiceFees)
	if err != nil {
		panic(err)
	}
	for _, provider := range data.Providers {
		providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
		if err != nil {
			panic(err)
		}
		err = k.SetProvider(ctx, providerAddr, provider)
		if err != nil {
			panic(err)
		}
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values to a genesis file,
// which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	remainingServiceFees, err := k.GetRemainingServiceFees(ctx)
	if err != nil {
		panic(err)
	}
	providers := k.GetAllProviders(ctx)
	return types.NewGenesisState(remainingServiceFees, providers)
}
