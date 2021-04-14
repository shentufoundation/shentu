package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

var _ types.QueryServer = Keeper{}

// Pool queries a pool based on the ID or sponsor.
func (q Keeper) Pool(c context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// query by ID
	pool, found := q.GetPool(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pool under ID %d doesn't exist", req.PoolId)
	}

	return &types.QueryPoolResponse{Pool: pool}, nil
}

// Pool queries a pool based on the ID or sponsor.
func (q Keeper) Sponsor(c context.Context, req *types.QuerySponsorRequest) (*types.QuerySponsorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// query by ID
	pool, found := q.GetPoolsBySponsor(ctx, req.Sponsor)
	if !found {
		return nil, status.Errorf(codes.NotFound, "there is no pool with sponsor %s", req.Sponsor)
	}

	return &types.QuerySponsorResponse{Pools: pool}, nil
}

// Pools queries all pools.
func (q Keeper) Pools(c context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryPoolsResponse{Pools: q.GetAllPools(ctx)}, nil
}

// PurchaseList queries a purchase list given a pool-purchase pair.
func (q Keeper) PurchaseList(c context.Context, req *types.QueryPurchaseListRequest) (*types.QueryPurchaseListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	purchaser, err := sdk.AccAddressFromBech32(req.Purchaser)
	if err != nil {
		return nil, err
	}

	purchaseList, found := q.GetPurchaseList(ctx, req.PoolId, purchaser)
	if !found {
		return nil, types.ErrPurchaseNotFound
	}

	return &types.QueryPurchaseListResponse{PurchaseList: purchaseList}, nil
}

// PurchaserPurchaseLists queries purchase lists for a given pool.
func (q Keeper) PoolPurchaseLists(c context.Context, req *types.QueryPoolPurchaseListsRequest) (*types.QueryPurchaseListsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	purchaseLists := q.GetPoolPurchaseLists(ctx, req.PoolId)

	return &types.QueryPurchaseListsResponse{PurchaseLists: purchaseLists}, nil
}

// PurchaseLists queries purchase lists purchaser.
func (q Keeper) PurchaseLists(c context.Context, req *types.QueryPurchaseListsRequest) (*types.QueryPurchaseListsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	var purchaseLists []types.PurchaseList
	purchaser, err := sdk.AccAddressFromBech32(req.Purchaser)
	if err != nil {
		return nil, err
	}
	purchaseLists = q.GetPurchaserPurchases(ctx, purchaser)

	return &types.QueryPurchaseListsResponse{PurchaseLists: purchaseLists}, nil
}

// Purchases queries all purchases.
func (q Keeper) Purchases(c context.Context, req *types.QueryPurchasesRequest) (*types.QueryPurchasesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryPurchasesResponse{Purchases: q.GetAllPurchases(ctx)}, nil
}

// Provider queries a provider given the address.
func (q Keeper) Provider(c context.Context, req *types.QueryProviderRequest) (*types.QueryProviderResponse, error) {
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

	return &types.QueryProviderResponse{Provider: provider}, nil
}

// Providers queries all providers.
func (q Keeper) Providers(c context.Context, req *types.QueryProvidersRequest) (*types.QueryProvidersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryProvidersResponse{Providers: q.GetAllProviders(ctx)}, nil
}

// PoolParams queries pool parameters.
func (q Keeper) PoolParams(c context.Context, req *types.QueryPoolParamsRequest) (*types.QueryPoolParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryPoolParamsResponse{Params: q.GetPoolParams(ctx)}, nil
}

// ClaimParams queries claim proposal parameters.
func (q Keeper) ClaimParams(c context.Context, req *types.QueryClaimParamsRequest) (*types.QueryClaimParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryClaimParamsResponse{Params: q.GetClaimProposalParams(ctx)}, nil
}

// ShieldStatus queries the global status of the shield module.
func (q Keeper) ShieldStatus(c context.Context, req *types.QueryShieldStatusRequest) (*types.QueryShieldStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryShieldStatusResponse{
		TotalCollateral:         q.GetTotalCollateral(ctx),
		TotalShield:             q.GetTotalShield(ctx),
		TotalWithdrawing:        q.GetTotalWithdrawing(ctx),
		CurrentServiceFees:      q.GetServiceFees(ctx),
		RemainingServiceFees:    q.GetRemainingServiceFees(ctx),
		GlobalShieldStakingPool: q.GetGlobalShieldStakingPool(ctx),
	}, nil
}

// ShieldStaking queries staked-for-shield for pool-purchaser pair.
func (q Keeper) ShieldStaking(c context.Context, req *types.QueryShieldStakingRequest) (*types.QueryShieldStakingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	purchaser, err := sdk.AccAddressFromBech32(req.Purchaser)
	if err != nil {
		return nil, err
	}
	shieldStaking, found := q.GetStakeForShield(ctx, req.PoolId, purchaser)
	if !found {
		return nil, types.ErrPurchaseNotFound
	}

	return &types.QueryShieldStakingResponse{ShieldStaking: shieldStaking}, nil
}

// ShieldStakingRate queries the shield staking rate for shield.
func (q Keeper) ShieldStakingRate(c context.Context, req *types.QueryShieldStakingRateRequest) (*types.QueryShieldStakingRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryShieldStakingRateResponse{Rate: q.GetShieldStakingRate(ctx)}, nil
}

// Reimbursement queries a reimbursement by proposal ID.
func (q Keeper) Reimbursement(c context.Context, req *types.QueryReimbursementRequest) (*types.QueryReimbursementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	reimbursement, err := q.GetReimbursement(ctx, req.ProposalId)
	if err != nil {
		return nil, err
	}

	return &types.QueryReimbursementResponse{Reimbursement: reimbursement}, nil
}

// Reimbursements queries all reimbursements.
func (q Keeper) Reimbursements(c context.Context, req *types.QueryReimbursementsRequest) (*types.QueryReimbursementsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryReimbursementsResponse{Pairs: q.GetAllProposalIDReimbursementPairs(ctx)}, nil
}
