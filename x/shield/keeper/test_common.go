package keeper

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// RandomValidator returns a random validator given access to the keeper and ctx.
func RandomValidator(r *rand.Rand, k Keeper, ctx sdk.Context) (staking.Validator, bool) {
	vals := k.sk.GetAllValidators(ctx)
	if len(vals) == 0 {
		return staking.Validator{}, false
	}

	i := r.Intn(len(vals))
	return vals[i], true
}

// RandomDelegation returns a random delegation info given access to the keeper and ctx.
func RandomDelegation(r *rand.Rand, k Keeper, ctx sdk.Context) (sdk.AccAddress, sdk.Int, bool) {
	val, ok := RandomValidator(r, k, ctx)
	if !ok {
		return nil, sdk.Int{}, false
	}

	dels := k.sk.GetValidatorDelegations(ctx, val.OperatorAddress)

	i := r.Intn(len(dels))
	return dels[i].DelegatorAddress, val.TokensFromShares(dels[i].Shares).TruncateInt(), true
}

// RandomPoolInfo returns info of a random pool given access to the keeper and ctx.
func RandomPoolInfo(r *rand.Rand, k Keeper, ctx sdk.Context) (uint64, string, bool) {
	pools := k.GetAllPools(ctx)
	if len(pools) == 0 {
		return 0, "", false
	}
	i := r.Intn(len(pools))
	return pools[i].ID, pools[i].Sponsor, true
}

// RandomPurchaseList returns a random purchase given access to the keeper and ctx.
func RandomPurchaseList(r *rand.Rand, k Keeper, ctx sdk.Context) (types.PurchaseList, bool) {
	purchaseLists := k.GetAllPurchaseLists(ctx)
	if len(purchaseLists) == 0 {
		return types.PurchaseList{}, false
	}
	i := r.Intn(len(purchaseLists))
	return purchaseLists[i], true
}

// RandomProvider returns a random provider of collaterals.
func RandomProvider(r *rand.Rand, k Keeper, ctx sdk.Context) (types.Provider, bool) {
	providers := k.GetAllProviders(ctx)
	if len(providers) == 0 {
		return types.Provider{}, false
	}

	i := r.Intn(len(providers))
	return providers[i], true
}

// RandomMaturedProposalIDReimbursementPair returns a random proposal ID - reimbursement pair for a matured reimbursement.
func RandomMaturedProposalIDReimbursementPair(r *rand.Rand, k Keeper, ctx sdk.Context) (types.ProposalIDReimbursementPair, bool) {
	prPairs := k.GetAllProposalIDReimbursementPairs(ctx)
	if len(prPairs) == 0 {
		return types.ProposalIDReimbursementPair{}, false
	}
	var maturedPRPairs []types.ProposalIDReimbursementPair
	for _, prPair := range prPairs {
		if prPair.Reimbursement.PayoutTime.Before(ctx.BlockTime()) {
			maturedPRPairs = append(maturedPRPairs, prPair)
		}
	}
	if len(maturedPRPairs) == 0 {
		return types.ProposalIDReimbursementPair{}, false
	}

	i := r.Intn(len(maturedPRPairs))
	return maturedPRPairs[i], true
}
