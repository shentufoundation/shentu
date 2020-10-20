package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// SetPoolParams sets parameters subspace for shield pool parameters.
func (k Keeper) SetPoolParams(ctx sdk.Context, poolParams types.PoolParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyPoolParams, &poolParams)
}

// GetPoolParams returns shield pool parameters.
func (k Keeper) GetPoolParams(ctx sdk.Context) types.PoolParams {
	var poolParams types.PoolParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyPoolParams, &poolParams)
	return poolParams
}

// SetClaimProposalParams sets parameters subspace for shield claim proposal parameters.
func (k Keeper) SetClaimProposalParams(ctx sdk.Context, claimProposalParams types.ClaimProposalParams) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyClaimProposalParams, &claimProposalParams)
}

// GetClaimProposalParams returns shield claim proposal parameters.
func (k Keeper) GetClaimProposalParams(ctx sdk.Context) types.ClaimProposalParams {
	var claimProposalParams types.ClaimProposalParams
	k.paramSpace.Get(ctx, types.ParamStoreKeyClaimProposalParams, &claimProposalParams)
	return claimProposalParams
}

// GetPurchaseDeletionPeriod returns time duration from purchase protection end time to deletion time.
func (k Keeper) GetPurchaseDeletionPeriod(ctx sdk.Context) time.Duration {
	paramProtectionPeriodMs := k.GetPoolParams(ctx).ProtectionPeriod.Milliseconds()
	paramClaimPeriodMs := k.GetClaimProposalParams(ctx).ClaimPeriod.Milliseconds()
	paramVotingPeriodMs := (k.GetVotingParams(ctx).VotingPeriod * 2).Milliseconds()
	return time.Duration(paramClaimPeriodMs-paramProtectionPeriodMs+paramVotingPeriodMs) * time.Millisecond
}
