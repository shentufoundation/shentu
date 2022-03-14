package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

// SetPoolParams sets parameters subspace for shield pool parameters.
func (k Keeper) SetPoolParams(ctx sdk.Context, poolParams v1beta1.PoolParams) {
	k.paramSpace.Set(ctx, v1beta1.ParamStoreKeyPoolParams, &poolParams)
}

// GetPoolParams returns shield pool parameters.
func (k Keeper) GetPoolParams(ctx sdk.Context) v1beta1.PoolParams {
	var poolParams v1beta1.PoolParams
	k.paramSpace.Get(ctx, v1beta1.ParamStoreKeyPoolParams, &poolParams)
	return poolParams
}

// SetClaimProposalParams sets parameters subspace for shield claim proposal parameters.
func (k Keeper) SetClaimProposalParams(ctx sdk.Context, claimProposalParams v1beta1.ClaimProposalParams) {
	k.paramSpace.Set(ctx, v1beta1.ParamStoreKeyClaimProposalParams, &claimProposalParams)
}

// GetClaimProposalParams returns shield claim proposal parameters.
func (k Keeper) GetClaimProposalParams(ctx sdk.Context) v1beta1.ClaimProposalParams {
	var claimProposalParams v1beta1.ClaimProposalParams
	k.paramSpace.Get(ctx, v1beta1.ParamStoreKeyClaimProposalParams, &claimProposalParams)
	return claimProposalParams
}

// SetBlockRewardParams sets parameters subspace for shield block reward parameters.
func (k Keeper) SetBlockRewardParams(ctx sdk.Context, blockRewardParams v1beta1.BlockRewardParams) {
	k.paramSpace.Set(ctx, v1beta1.ParamStoreKeyBlockRewardParams, &blockRewardParams)
}

// GetBlockRewardParams returns shield block reward parameters.
func (k Keeper) GetBlockRewardParams(ctx sdk.Context) v1beta1.BlockRewardParams {
	var blockRewardParams v1beta1.BlockRewardParams
	k.paramSpace.Get(ctx, v1beta1.ParamStoreKeyBlockRewardParams, &blockRewardParams)
	return blockRewardParams
}
