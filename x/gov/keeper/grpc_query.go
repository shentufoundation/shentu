package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

var _ typesv1.QueryServer = queryServer{}

type queryServer struct{ k *Keeper }

func NewQueryServer(k *Keeper) typesv1.QueryServer {
	return queryServer{k: k}
}

func (q queryServer) Constitution(ctx context.Context, request *govtypesv1.QueryConstitutionRequest) (*govtypesv1.QueryConstitutionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q queryServer) Proposal(ctx context.Context, request *govtypesv1.QueryProposalRequest) (*govtypesv1.QueryProposalResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q queryServer) Proposals(ctx context.Context, request *govtypesv1.QueryProposalsRequest) (*govtypesv1.QueryProposalsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q queryServer) Vote(ctx context.Context, request *govtypesv1.QueryVoteRequest) (*govtypesv1.QueryVoteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q queryServer) Votes(ctx context.Context, request *govtypesv1.QueryVotesRequest) (*govtypesv1.QueryVotesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q queryServer) Deposit(ctx context.Context, request *govtypesv1.QueryDepositRequest) (*govtypesv1.QueryDepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q queryServer) Deposits(ctx context.Context, request *govtypesv1.QueryDepositsRequest) (*govtypesv1.QueryDepositsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q queryServer) TallyResult(ctx context.Context, req *govtypesv1.QueryTallyResultRequest) (*govtypesv1.QueryTallyResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	proposal, err := q.k.Proposals.Get(ctx, req.ProposalId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "proposal %d doesn't exist", req.ProposalId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	var tallyResult govtypesv1.TallyResult

	switch {
	case proposal.Status == govtypesv1.StatusDepositPeriod:
		tallyResult = govtypesv1.EmptyTallyResult()

	case proposal.Status == govtypesv1.StatusPassed || proposal.Status == govtypesv1.StatusRejected || proposal.Status == govtypesv1.StatusFailed:
		tallyResult = *proposal.FinalTallyResult

	default:
		// proposal is in voting period
		var err error
		_, _, tallyResult, err = q.k.Tally(ctx, proposal)
		if err != nil {
			return nil, err
		}
	}

	return &govtypesv1.QueryTallyResultResponse{Tally: &tallyResult}, nil
}

// CertVoted returns certifier voting
func (q queryServer) CertVoted(c context.Context, req *typesv1.QueryCertVotedRequest) (*typesv1.QueryCertVotedResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	voted, err := q.k.GetCertifierVoted(ctx, req.ProposalId)
	if err != nil {
		return nil, err
	}
	return &typesv1.QueryCertVotedResponse{CertVoted: voted}, nil
}

// Params queries all params
// func (q queryServer) Params(ctx context.Context, request *govtypesv1.QueryParamsRequest) (*typesv1.QueryParamsResponse, error) {
func (q queryServer) Params(c context.Context, req *govtypesv1.QueryParamsRequest) (*typesv1.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	params, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
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
		customParams, err := q.k.GetCustomParams(ctx)
		if err != nil {
			return nil, err
		}
		response.CustomParams = &customParams

	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"%s is not a valid parameter type", req.ParamsType)
	}

	response.Params = &params

	return response, nil
}
