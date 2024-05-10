package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

var _ types.QueryServer = Keeper{}

// CertVoted returns certifier voting
func (k Keeper) CertVoted(c context.Context, req *types.QueryCertVotedRequest) (*types.QueryCertVotedResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	voted := k.GetCertifierVoted(ctx, req.ProposalId)

	return &types.QueryCertVotedResponse{CertVoted: voted}, nil
}

// Params queries all params
func (k Keeper) Params(c context.Context, req *govtypesv1.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	switch req.ParamsType {
	case govtypesv1.ParamDeposit:
		depositParmas := k.GetDepositParams(ctx)
		return &types.QueryParamsResponse{DepositParams: depositParmas}, nil

	case govtypesv1.ParamVoting:
		votingParmas := k.GetVotingParams(ctx)
		return &types.QueryParamsResponse{VotingParams: votingParmas}, nil

	case govtypesv1.ParamTallying:
		tallyParams := k.GetTallyParams(ctx)
		return &types.QueryParamsResponse{TallyParams: tallyParams}, nil

	case types.ParamCustom:
		customParams := k.GetCustomParams(ctx)
		return &types.QueryParamsResponse{CustomParams: customParams}, nil

	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"%s is not a valid parameter type", req.ParamsType)
	}
}
