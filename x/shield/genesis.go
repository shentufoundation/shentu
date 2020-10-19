package shield

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// InitGenesis initialize store values with genesis states.
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	k.SetPoolParams(ctx, data.PoolParams)
	k.SetClaimProposalParams(ctx, data.ClaimProposalParams)
	k.SetAdmin(ctx, data.ShieldAdmin)
	k.SetTotalCollateral(ctx, data.TotalCollateral)
	k.SetTotalShield(ctx, data.TotalShield)
	k.SetTotalLocked(ctx, data.TotalLocked)
	k.SetServiceFees(ctx, data.ServiceFees)
	for _, pool := range data.Pools {
		k.SetPool(ctx, pool)
	}
	k.SetNextPoolID(ctx, data.NextPoolID)
	k.SetNextPurchaseID(ctx, data.NextPurchaseID)
	for _, purchaseList := range data.PurchaseLists {
		k.SetPurchaseList(ctx, purchaseList)
		for _, entry := range purchaseList.Entries {
			k.InsertPurchaseQueue(ctx, purchaseList, entry.ProtectionEndTime.Add(k.GetPurchaseDeletionPeriod(ctx)))
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

// ExportGenesis writes the current store values to a genesis file,
// which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	poolParams := k.GetPoolParams(ctx)
	claimProposalParams := k.GetClaimProposalParams(ctx)
	shieldAdmin := k.GetAdmin(ctx)
	totalCollateral := k.GetTotalCollateral(ctx)
	totalShield := k.GetTotalShield(ctx)
	totalLocked := k.GetTotalLocked(ctx)
	serviceFees := k.GetServiceFees(ctx)
	pools := k.GetAllPools(ctx)
	nextPoolID := k.GetNextPoolID(ctx)
	nextPurchaseID := k.GetNextPurchaseID(ctx)
	purchaseLists := k.GetAllPurchaseLists(ctx)
	providers := k.GetAllProviders(ctx)
	withdraws := k.GetAllWithdraws(ctx)

	return types.NewGenesisState(shieldAdmin, nextPoolID, nextPurchaseID, poolParams, claimProposalParams, totalCollateral, totalShield, totalLocked, serviceFees, pools, providers, purchaseLists, withdraws)
}
