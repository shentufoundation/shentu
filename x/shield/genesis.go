package shield

import (
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/keeper"
	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// InitGenesis initialize store values with genesis states.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	k.SetPoolParams(ctx, data.PoolParams)
	k.SetClaimProposalParams(ctx, data.ClaimProposalParams)

	adminAddr := sdk.AccAddress{}
	var err error
	if len(strings.TrimSpace(data.ShieldAdmin)) != 0 {
		adminAddr, err = sdk.AccAddressFromBech32(data.ShieldAdmin)
		if err != nil {
			panic(err)
		}
	}

	k.SetAdmin(ctx, adminAddr)
	k.SetTotalCollateral(ctx, data.TotalCollateral)
	k.SetTotalWithdrawing(ctx, data.TotalWithdrawing)
	k.SetTotalShield(ctx, data.TotalShield)
	k.SetTotalClaimed(ctx, data.TotalClaimed)
	k.SetServiceFees(ctx, data.ServiceFees)
	k.SetRemainingServiceFees(ctx, data.RemainingServiceFees)
	k.SetGlobalStakingPool(ctx, data.GlobalStakingPool)
	k.SetShieldStakingRate(ctx, data.ShieldStakingRate)
	for _, pool := range data.Pools {
		k.SetPool(ctx, pool)
	}
	k.SetNextPoolID(ctx, data.NextPoolId)
	k.SetNextPurchaseID(ctx, data.NextPurchaseId)

	for _, purchase := range data.Purchases {
		purchaserAddr, err := sdk.AccAddressFromBech32(purchase.Purchaser)
		if err != nil {
			panic(err)
		}
		k.SetPurchase(ctx, purchase.PoolId, purchaserAddr, purchase)
	}

	for _, provider := range data.Providers {
		providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
		if err != nil {
			panic(err)
		}
		k.SetProvider(ctx, providerAddr, provider)
	}
	for _, withdraw := range data.Withdraws {
		k.InsertWithdrawQueue(ctx, withdraw)
	}
	for _, pRPair := range data.ProposalIDReimbursementPairs {
		k.SetReimbursement(ctx, pRPair.ProposalId, pRPair.Reimbursement)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values to a genesis file,
// which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
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
	providers := k.GetAllProviders(ctx)
	withdraws := k.GetAllWithdraws(ctx)
	stakingPurchaseRate := k.GetShieldStakingRate(ctx)
	globalStakingPool := k.GetGlobalStakingPool(ctx)
	stakingPurchases := k.GetAllPurchase(ctx)
	reimbursements := k.GetAllProposalIDReimbursementPairs(ctx)

	return types.NewGenesisState(shieldAdmin, nextPoolID, nextPurchaseID, poolParams, claimProposalParams,
		totalCollateral, totalWithdrawing, totalShield, totalClaimed, serviceFees, remainingServiceFees, pools,
		providers, withdraws, stakingPurchaseRate, globalStakingPool, stakingPurchases, reimbursements)
}
