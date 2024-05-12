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

	switch req.ParamsType {
	case govtypesv1.ParamDeposit:
		depositParmas := k.GetDepositParams(ctx)
		return &typesv1.QueryParamsResponse{DepositParams: &depositParmas}, nil

	case govtypesv1.ParamVoting:
		votingParmas := k.GetVotingParams(ctx)
		return &typesv1.QueryParamsResponse{VotingParams: &votingParmas}, nil

	case govtypesv1.ParamTallying:
		tallyParams := k.GetTallyParams(ctx)
		return &typesv1.QueryParamsResponse{TallyParams: &tallyParams}, nil

	case typesv1.ParamCustom:
		customParams := k.GetCustomParams(ctx)
		return &typesv1.QueryParamsResponse{CustomParams: &customParams}, nil

	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"%s is not a valid parameter type", req.ParamsType)
	}
}
