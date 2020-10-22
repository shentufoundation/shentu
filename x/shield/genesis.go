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
	k.SetTotalWithdrawing(ctx, data.TotalWithdrawing)
	k.SetTotalShield(ctx, data.TotalShield)
	k.SetTotalLocked(ctx, data.TotalLocked)
	k.SetServiceFees(ctx, data.ServiceFees)
	k.SetRemainingServiceFees(ctx, data.RemainingServiceFees)
	for _, pool := range data.Pools {
		k.SetPool(ctx, pool)
	}
	k.SetNextPoolID(ctx, data.NextPoolID)
	k.SetNextPurchaseID(ctx, data.NextPurchaseID)
	for _, purchaseList := range data.PurchaseLists {
		k.SetPurchaseList(ctx, purchaseList)
		for _, entry := range purchaseList.Entries {
			k.InsertExpiringPurchaseQueue(ctx, purchaseList, entry.ProtectionEndTime)
		}
	}
	for _, provider := range data.Providers {
		k.SetProvider(ctx, provider.Address, provider)
	}
	for _, withdraw := range data.Withdraws {
		k.InsertWithdrawQueue(ctx, withdraw)
	}
	k.SetLastUpdateTime(ctx, data.LastUpdateTime)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values to a genesis file,
// which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	poolParams := k.GetPoolParams(ctx)
	claimProposalParams := k.GetClaimProposalParams(ctx)
	shieldAdmin := k.GetAdmin(ctx)
	totalCollateral := k.GetTotalCollateral(ctx)
	totalWithdrawing := k.GetTotalWithdrawing(ctx)
	totalShield := k.GetTotalShield(ctx)
	totalLocked := k.GetTotalLocked(ctx)
	serviceFees := k.GetServiceFees(ctx)
	remainingServiceFees := k.GetRemainingServiceFees(ctx)
	pools := k.GetAllPools(ctx)
	nextPoolID := k.GetNextPoolID(ctx)
	nextPurchaseID := k.GetNextPurchaseID(ctx)
	purchaseLists := k.GetAllPurchaseLists(ctx)
	providers := k.GetAllProviders(ctx)
	withdraws := k.GetAllWithdraws(ctx)
	lastUpdateTime, _ := k.GetLastUpdateTime(ctx)

	return types.NewGenesisState(shieldAdmin, nextPoolID, nextPurchaseID, poolParams, claimProposalParams, totalCollateral, totalWithdrawing, totalShield, totalLocked, serviceFees, remainingServiceFees, pools, providers, purchaseLists, withdraws, lastUpdateTime)
}
