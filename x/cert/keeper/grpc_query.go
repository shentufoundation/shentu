package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

var _ types.QueryServer = Querier{}

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

// Certifier queries a certifier given its address or alias.
func (q Querier) Certifier(c context.Context, req *types.QueryCertifierRequest) (*types.QueryCertifierResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	var certifier types.Certifier
	var err error
	if req.Alias != "" {
		// query by alias
		certifier, err = q.GetCertifierByAlias(ctx, req.Alias)
		if err != nil {
			return nil, err
		}
	} else {
		// query by address
		certifierAddr, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			return nil, err
		}

		certifier, err = q.GetCertifier(ctx, certifierAddr)
		if err != nil {
			return nil, err
		}
	}

	return &types.QueryCertifierResponse{Certifier: certifier}, nil
}

// Certifiers queries all certifiers.
func (q Querier) Certifiers(c context.Context, req *types.QueryCertifiersRequest) (*types.QueryCertifiersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryCertifiersResponse{Certifiers: q.GetAllCertifiers(ctx)}, nil
}
