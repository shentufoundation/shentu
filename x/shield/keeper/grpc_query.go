package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

var _ v1beta1.QueryServer = Keeper{}

// Pool queries a pool based on the ID or sponsor.
func (q Keeper) Pool(c context.Context, req *v1beta1.QueryPoolRequest) (*v1beta1.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// query by ID
	pool, found := q.GetPool(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pool under ID %d doesn't exist", req.PoolId)
	}

	return &v1beta1.QueryPoolResponse{Pool: pool}, nil
}

// Pool queries a pool based on the ID or sponsor.
func (q Keeper) Sponsor(c context.Context, req *v1beta1.QuerySponsorRequest) (*v1beta1.QuerySponsorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// query by ID
	pool, found := q.GetPoolsBySponsor(ctx, req.Sponsor)
	if !found {
		return nil, status.Errorf(codes.NotFound, "there is no pool with sponsor %s", req.Sponsor)
	}

	return &v1beta1.QuerySponsorResponse{Pools: pool}, nil
}

// Pools queries all pools.
func (q Keeper) Pools(c context.Context, req *v1beta1.QueryPoolsRequest) (*v1beta1.QueryPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &v1beta1.QueryPoolsResponse{Pools: q.GetAllPools(ctx)}, nil
}

// Provider queries a provider given the address.
func (q Keeper) Provider(c context.Context, req *v1beta1.QueryProviderRequest) (*v1beta1.QueryProviderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	provider, found := q.GetProvider(ctx, address)
	if !found {
		return nil, types.ErrProviderNotFound
	}

	return &v1beta1.QueryProviderResponse{Provider: provider}, nil
}

// Providers queries all providers.
func (q Keeper) Providers(c context.Context, req *v1beta1.QueryProvidersRequest) (*v1beta1.QueryProvidersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &v1beta1.QueryProvidersResponse{Providers: q.GetAllProviders(ctx)}, nil
}

// PoolParams queries pool parameters.
func (q Keeper) PoolParams(c context.Context, req *v1beta1.QueryPoolParamsRequest) (*v1beta1.QueryPoolParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &v1beta1.QueryPoolParamsResponse{Params: q.GetPoolParams(ctx)}, nil
}

// ClaimParams queries claim proposal parameters.
func (q Keeper) ClaimParams(c context.Context, req *v1beta1.QueryClaimParamsRequest) (*v1beta1.QueryClaimParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &v1beta1.QueryClaimParamsResponse{Params: q.GetClaimProposalParams(ctx)}, nil
}

// BlockRewardParams queries block reward parameters.
func (q Keeper) BlockRewardParams(c context.Context, req *v1beta1.QueryBlockRewardParamsRequest) (*v1beta1.QueryBlockRewardParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &v1beta1.QueryBlockRewardParamsResponse{Params: q.GetBlockRewardParams(ctx)}, nil
}

// ShieldStatus queries the global status of the shield module.
func (q Keeper) ShieldStatus(c context.Context, req *v1beta1.QueryShieldStatusRequest) (*v1beta1.QueryShieldStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &v1beta1.QueryShieldStatusResponse{
		TotalCollateral:         q.GetTotalCollateral(ctx),
		TotalShield:             q.GetTotalShield(ctx),
		TotalWithdrawing:        q.GetTotalWithdrawing(ctx),
		GlobalShieldStakingPool: q.GetGlobalStakingPool(ctx),
	}, nil
}

// ShieldStaking queries staked-for-shield for pool-purchaser pair.
func (q Keeper) Purchase(c context.Context, req *v1beta1.QueryPurchaseRequest) (*v1beta1.QueryPurchaseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	purchaser, err := sdk.AccAddressFromBech32(req.Purchaser)
	if err != nil {
		return nil, err
	}
	shieldStaking, found := q.GetPurchase(ctx, req.PoolId, purchaser)
	if !found {
		return nil, types.ErrPurchaseNotFound
	}

	return &v1beta1.QueryPurchaseResponse{Purchase: shieldStaking}, nil
}

// Reserve queries all shield reserve amount.
func (q Keeper) Reserve(c context.Context, req *v1beta1.QueryReserveRequest) (*v1beta1.QueryReserveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	return &v1beta1.QueryReserveResponse{Reserve: q.GetReserve(ctx)}, nil
}

// PendingPayouts queries all pending payouts.
func (q Keeper) PendingPayouts(c context.Context, req *v1beta1.QueryPendingPayoutsRequest) (*v1beta1.QueryPendingPayoutsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	return &v1beta1.QueryPendingPayoutsResponse{PendingPayouts: q.GetAllPendingPayouts(ctx)}, nil
}

// PoolPurchases queries for all purchases for a specific pool.
func (k Keeper) PoolPurchases(c context.Context, req *v1beta1.QueryPoolPurchasesRequest) (*v1beta1.QueryPurchasesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	res := k.GetPoolPurchases(ctx, req.PoolId)
	return &v1beta1.QueryPurchasesResponse{Purchases: res}, nil
}

// Purchases queries for all purchases.
func (k Keeper) Purchases(c context.Context, req *v1beta1.QueryAllPurchasesRequest) (*v1beta1.QueryPurchasesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	res := k.GetAllPurchase(ctx)
	return &v1beta1.QueryPurchasesResponse{Purchases: res}, nil
}

// Purchaser queries for information on a purchaser.
func (k Keeper) Purchaser(c context.Context, req *v1beta1.QueryPurchaserRequest) (*v1beta1.QueryPurchaserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	purchaser, err := sdk.AccAddressFromBech32(req.Purchaser)
	if err != nil {
		return nil, err
	}
	res := k.GetPurchaserPurchases(ctx, purchaser)

	shield := sdk.ZeroInt()
	deposit := sdk.ZeroInt()
	for _, p := range res {
		shield = shield.Add(p.Shield)
		deposit = deposit.Add(p.Amount)
	}
	return &v1beta1.QueryPurchaserResponse{
		Purchases:    res,
		TotalShield:  shield,
		TotalDeposit: deposit,
	}, nil
}
