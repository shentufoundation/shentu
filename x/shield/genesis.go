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
	k.SetNextPurchaseID(ctx, data.NextPurchaseID)
	k.SetPoolParams(ctx, data.PoolParams)
	k.SetClaimProposalParams(ctx, data.ClaimProposalParams)
	for _, pool := range data.Pools {
		k.SetPool(ctx, pool)
	}
	for _, collateral := range data.Collaterals {
		pool, found := k.GetPool(ctx, collateral.PoolID)
		if !found {
			panic(types.ErrNoPoolFound)
		}
		k.SetCollateral(ctx, pool, collateral.Provider, collateral)
	}
	for _, purchaseList := range data.PurchaseLists {
		k.SetPurchaseList(ctx, purchaseList)
		for _, entry := range purchaseList.Entries {
			k.InsertPurchaseQueue(ctx, purchaseList, entry.DeleteTime)
		}
	}
	for _, provider := range data.Providers {
		k.SetProvider(ctx, provider.Address, provider)
	}
	for _, withdraw := range data.Withdraws {
		k.InsertWithdrawQueue(ctx, withdraw)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	shieldAdmin := k.GetAdmin(ctx)
	nextPoolID := k.GetNextPoolID(ctx)
	nextPurchaseID := k.GetNextPurchaseID(ctx)
	poolParams := k.GetPoolParams(ctx)
	claimProposalParams := k.GetClaimProposalParams(ctx)
	pools := k.GetAllPools(ctx)
	collaterals := k.GetAllCollaterals(ctx)
	providers := k.GetAllProviders(ctx)
	purchaseLists := k.GetAllPurchaseLists(ctx)
	withdraws := k.GetAllWithdraws(ctx)

	return types.NewGenesisState(shieldAdmin, nextPoolID, nextPurchaseID, poolParams, claimProposalParams, pools, collaterals, providers, purchaseLists, withdraws)
}
