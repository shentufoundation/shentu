package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/cert/types"
)

var _ types.QueryServer = Keeper{}

// Certifier queries a certifier given its address or alias.
func (q Keeper) Certifier(c context.Context, req *types.QueryCertifierRequest) (*types.QueryCertifierResponse, error) {
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
			panic(err)
		}

		certifier, err = q.GetCertifier(ctx, certifierAddr)
		if err != nil {
			return nil, err
		}
	}

	return &types.QueryCertifierResponse{Certifier: certifier}, nil
}

// Certifiers queries all certifiers.
func (q Keeper) Certifiers(c context.Context, req *types.QueryCertifiersRequest) (*types.QueryCertifiersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryCertifiersResponse{Certifiers: q.GetAllCertifiers(ctx)}, nil
}

// Validator queries the certifier of a certified validator.
func (q Keeper) Validator(c context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	pk, ok := req.Pubkey.GetCachedValue().(crypto.PubKey)
	if !ok {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into cryto.PubKey %T", req.Pubkey))
	}

	certifier, err := q.GetValidatorCertifier(ctx, pk)
	if err != nil {
		return nil, err
	}

	return &types.QueryValidatorResponse{Certifier: certifier.String()}, nil
}

// Validators returns all validators' public keys.
func (q Keeper) Validators(c context.Context, req *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryValidatorsResponse{Pubkeys: q.GetAllValidatorPubkeys(ctx)}, nil
}

func (q Keeper) Platform(c context.Context, req *types.QueryPlatformRequest) (*types.QueryPlatformResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	pk, ok := req.Pubkey.GetCachedValue().(crypto.PubKey)
	if !ok {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into cryto.PubKey %T", req.Pubkey))
	}

	platform, ok := q.GetPlatform(ctx, pk)
	if !ok {
		return nil, nil
	}

	return &types.QueryPlatformResponse{Platform: platform}, nil
}

func (q Keeper) Certificate(c context.Context, req *types.QueryCertificateRequest) (*types.QueryCertificateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	certificate, err := q.GetCertificateByID(ctx, types.CertificateID(req.CertificateId))
	if err != nil {
		return nil, err
	}

	reqContent := certificate.RequestContent()

	return &types.QueryCertificateResponse{
		CertificateId:      certificate.ID().String(),
		CertificateType:    certificate.Type().String(),
		RequestContent:     &reqContent,
		CertificateContent: certificate.FormattedCertificateContent(),
		Description:        certificate.Description(),
		Certifier:          certificate.Certifier().String(),
		TxHash:             certificate.TxHash(),
	}, nil
}

func (q Keeper) Certificates(c context.Context, req *types.QueryCertificatesRequest) (*types.QueryCertificatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	certifierAddr, err := sdk.AccAddressFromBech32(req.Certifier)
	if err != nil {
		panic(err)
	}

	// TODO Use new cosmos pagniation mechanism
	params := types.QueryCertificatesParams{
		Page:        int(req.Pagination.Offset),
		Limit:       int(req.Pagination.Limit),
		Certifier:   certifierAddr,
		ContentType: req.ContentType,
		Content:     req.Content,
	}

	total, certificates, err := q.GetCertificatesFiltered(ctx, params)
	if err != nil {
		return nil, err
	}

	var results []types.QueryCertificateResponse
	for _, certificate := range certificates {
		reqContent := certificate.RequestContent()

		results = append(results,
			types.QueryCertificateResponse{
				CertificateId: certificate.ID().String(), CertificateType: certificate.Type().String(),
				RequestContent:     &reqContent,
				CertificateContent: certificate.FormattedCertificateContent(),
				Description:        certificate.Description(),
				Certifier:          certificate.Certifier().String(),
				TxHash:             certificate.TxHash(),
			})
	}

	return &types.QueryCertificatesResponse{Total: total, Certificates: results}, nil
}
