package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

var _ typesv1.CustomQueryServer = customQueryServer{}

type customQueryServer struct{ k *Keeper }

func NewCustomQueryServer(k *Keeper) typesv1.CustomQueryServer {
	return customQueryServer{k: k}
}

// CertVoted returns certifier voting
func (cq customQueryServer) CertVoted(c context.Context, req *typesv1.QueryCertVotedRequest) (*typesv1.QueryCertVotedResponse, error) {
	//ctx := sdk.UnwrapSDKContext(c)
	voted, err := cq.k.GetCertifierVoted(c, req.ProposalId)
	if err != nil {
		return nil, err
	}
	return &typesv1.QueryCertVotedResponse{CertVoted: voted}, nil
}

func (cq customQueryServer) CustomParams(c context.Context, req *govtypesv1.QueryParamsRequest) (*typesv1.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	response := &typesv1.QueryParamsResponse{}

	switch req.ParamsType {
	case typesv1.ParamCustom:
		customParams, err := cq.k.GetCustomParams(ctx)
		if err != nil {
			return nil, err
		}
		response.CustomParams = &customParams

	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"%s is not a valid parameter type", req.ParamsType)
	}

	return response, nil
}
