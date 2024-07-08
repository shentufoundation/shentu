package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

var _ typesv1.QueryServer = Keeper{}

// CertVoted returns certifier voting
func (k Keeper) CertVoted(c context.Context, req *typesv1.QueryCertVotedRequest) (*typesv1.QueryCertVotedResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	voted := k.GetCertifierVoted(ctx, req.ProposalId)

	return &typesv1.QueryCertVotedResponse{CertVoted: voted}, nil
}

// Params queries all params
func (k Keeper) Params(c context.Context, req *govtypesv1.QueryParamsRequest) (*typesv1.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	response := &typesv1.QueryParamsResponse{}

	switch req.ParamsType {
	case govtypesv1.ParamDeposit:
		depositParams := govtypesv1.NewDepositParams(params.MinDeposit, params.MaxDepositPeriod)
		response.DepositParams = &depositParams

	case govtypesv1.ParamVoting:
		votingParams := govtypesv1.NewVotingParams(params.VotingPeriod)
		response.VotingParams = &votingParams

	case govtypesv1.ParamTallying:
		tallyParams := govtypesv1.NewTallyParams(params.Quorum, params.Threshold, params.VetoThreshold)
		response.TallyParams = &tallyParams

	case typesv1.ParamCustom:
		customParams := k.GetCustomParams(ctx)
		response.CustomParams = &customParams

	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"%s is not a valid parameter type", req.ParamsType)
	}

	response.Params = &params

	return response, nil
}
