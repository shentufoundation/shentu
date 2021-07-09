package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	qtypes "github.com/cosmos/cosmos-sdk/types/query"

	"github.com/certikfoundation/shentu/x/nft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Admin(c context.Context, req *types.QueryAdminRequest) (*types.QueryAdminResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	admin, err := k.GetAdmin(ctx, addr)
	return &types.QueryAdminResponse{Admin: admin}, err
}

func (k Keeper) Admins(c context.Context, _ *types.QueryAdminsRequest) (*types.QueryAdminsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	admins := k.GetAdmins(ctx)
	return &types.QueryAdminsResponse{Admins: admins}, nil
}

func (k Keeper) Certificate(c context.Context, req *types.QueryCertificateRequest) (*types.QueryCertificateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	certificate, err := k.GetCertificate(ctx, req.DenomId, req.TokenId)
	if err != nil {
		return nil, err
	}

	return &types.QueryCertificateResponse{Certificate: certificate}, nil
}

func (k Keeper) Certificates(c context.Context, req *types.QueryCertificatesRequest) (*types.QueryCertificatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	var certifierAddr sdk.AccAddress
	var err error
	if req.Certifier != "" {
		certifierAddr, err = sdk.AccAddressFromBech32(req.Certifier)
		if err != nil {
			return nil, err
		}
	}

	page, limit, err := qtypes.ParsePagination(req.Pagination)
	if err != nil {
		return nil, err
	}
	params := types.QueryCertificatesParams{
		Page:      page,
		Limit:     limit,
		Certifier: certifierAddr,
		DenomID:   req.DenomId,
	}

	total, certificates, err := k.GetCertificatesFiltered(ctx, params)
	if err != nil {
		return nil, err
	}

	results := make([]types.QueryCertificateResponse, total)
	for i, certificate := range certificates {
		results[i] = types.QueryCertificateResponse{
			Certificate: certificate,
		}
	}

	return &types.QueryCertificatesResponse{Total: total, Certificates: results}, nil
}
