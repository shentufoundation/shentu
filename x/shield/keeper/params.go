package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
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

// GetShieldStakingRate returns shield to staked rate.
func (k Keeper) GetShieldStakingRate(ctx sdk.Context) (rate sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyStakingShieldRate, &rate)
	return
}

// SetShieldStakingRate sets shield to staked rate.
func (k Keeper) SetShieldStakingRate(ctx sdk.Context, rate sdk.Dec) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyStakingShieldRate, &rate)
}

// GetDistributionParams returns distribution parameters.
func (k Keeper) GetDistributionParams(ctx sdk.Context) (distrParams types.DistributionParams) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyDistribution, &distrParams)
	return
}

// SetDistributionParams sets distribution parameters.
func (k Keeper) SetDistributionParams(ctx sdk.Context, distrParams types.D) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyDistribution, &rate)
}
