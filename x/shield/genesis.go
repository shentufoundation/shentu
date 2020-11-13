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
	k.SetTotalClaimed(ctx, data.TotalClaimed)
	k.SetServiceFees(ctx, data.ServiceFees)
	k.SetRemainingServiceFees(ctx, data.RemainingServiceFees)
	k.SetGlobalShieldStakingPool(ctx, data.GlobalStakingPool)
	k.SetShieldStakingRate(ctx, data.ShieldStakingRate)
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
	for _, purchase := range data.StakeForShields {
		k.SetStakeForShield(ctx, purchase.PoolID, purchase.Purchaser, purchase)
	}
	for _, originalStaking := range data.OriginalStakings {
		k.SetOriginalStaking(ctx, originalStaking.PurchaseID, originalStaking.Amount)
	}
	for _, provider := range data.Providers {
		k.SetProvider(ctx, provider.Address, provider)
	}
	for _, withdraw := range data.Withdraws {
		k.InsertWithdrawQueue(ctx, withdraw)
	}
	k.SetLastUpdateTime(ctx, data.LastUpdateTime)
	for _, pRPair := range data.ProposalIDReimbursementPairs {
		k.SetReimbursement(ctx, pRPair.ProposalID, pRPair.Reimbursement)
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
	totalWithdrawing := k.GetTotalWithdrawing(ctx)
	totalShield := k.GetTotalShield(ctx)
	totalClaimed := k.GetTotalClaimed(ctx)
	serviceFees := k.GetServiceFees(ctx)
	remainingServiceFees := k.GetRemainingServiceFees(ctx)
	pools := k.GetAllPools(ctx)
	nextPoolID := k.GetNextPoolID(ctx)
	nextPurchaseID := k.GetNextPurchaseID(ctx)
	purchaseLists := k.GetAllPurchaseLists(ctx)
	providers := k.GetAllProviders(ctx)
	withdraws := k.GetAllWithdraws(ctx)
	lastUpdateTime, _ := k.GetLastUpdateTime(ctx)
	stakingPurchaseRate := k.GetShieldStakingRate(ctx)
	globalStakingPool := k.GetGlobalShieldStakingPool(ctx)
	stakingPurchases := k.GetAllStakeForShields(ctx)
	originalStaking := k.GetAllOriginalStakings(ctx)
	reimbursements := k.GetAllProposalIDReimbursementPairs(ctx)

	return types.NewGenesisState(shieldAdmin, nextPoolID, nextPurchaseID, poolParams, claimProposalParams,
		totalCollateral, totalWithdrawing, totalShield, totalClaimed, serviceFees, remainingServiceFees, pools,
		providers, purchaseLists, withdraws, lastUpdateTime, stakingPurchaseRate, globalStakingPool, stakingPurchases, originalStaking, reimbursements)
}
