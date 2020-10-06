package shield

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// InitGenesis initialize store values with genesis states.
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	k.SetAdmin(ctx, data.ShieldAdmin)
	k.SetNextPoolID(ctx, data.NextPoolID)
	k.SetPoolParams(ctx, data.PoolParams)
	k.SetClaimProposalParams(ctx, data.ClaimProposalParams)
	for _, pool := range data.Pools {
		k.SetPool(ctx, pool)
	}
	for _, collateral := range data.Collaterals {
		pool, err := k.GetPool(ctx, collateral.PoolID)
		if err != nil {
			panic(err)
		}
		k.SetCollateral(ctx, pool, collateral.Provider, collateral)
	}
	for _, purchase := range data.Purchases {
		k.SetPurchase(ctx, purchase.TxHash, purchase)
	}
	for _, provider := range data.Providers {
		k.SetProvider(ctx, provider.Address, provider)
		k.UpdateDelegationAmount(ctx, provider.Address)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	shieldAdmin := k.GetAdmin(ctx)
	nextPoolID := k.GetNextPoolID(ctx)
	poolParams := k.GetPoolParams(ctx)
	claimProposalParams := k.GetClaimProposalParams(ctx)
	pools := k.GetAllPools(ctx)
	collaterals := k.GetAllCollaterals(ctx)
	providers := k.GetAllProviders(ctx)
	purchases := k.GetAllPurchases(ctx)

	return types.NewGenesisState(shieldAdmin, nextPoolID, poolParams, claimProposalParams, pools, collaterals, providers, purchases)
}
