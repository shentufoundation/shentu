package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

var _ typesv1.QueryServer = queryServer{}

type queryServer struct{ k *Keeper }

func NewQueryServer(k *Keeper) typesv1.QueryServer {
	return queryServer{k: k}
}

func (q queryServer) Constitution(ctx context.Context, _ *govtypesv1.QueryConstitutionRequest) (*govtypesv1.QueryConstitutionResponse, error) {
	constitution, err := q.k.Constitution.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &govtypesv1.QueryConstitutionResponse{Constitution: constitution}, nil
}

// Proposal returns proposal details based on ProposalID
func (q queryServer) Proposal(ctx context.Context, req *govtypesv1.QueryProposalRequest) (*govtypesv1.QueryProposalResponse, error) {
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

	return &govtypesv1.QueryProposalResponse{Proposal: &proposal}, nil
}

// Proposals implements the Query/Proposals gRPC method
func (q queryServer) Proposals(ctx context.Context, req *govtypesv1.QueryProposalsRequest) (*govtypesv1.QueryProposalsResponse, error) {
	filteredProposals, pageRes, err := query.CollectionFilteredPaginate(ctx, q.k.Proposals, req.Pagination, func(key uint64, p govtypesv1.Proposal) (include bool, err error) {
		matchVoter, matchDepositor, matchStatus := true, true, true

		// match status (if supplied/valid)
		if govtypesv1.ValidProposalStatus(req.ProposalStatus) {
			matchStatus = p.Status == req.ProposalStatus
		}

		// match voter address (if supplied)
		if len(req.Voter) > 0 {
			voter, err := q.k.authKeeper.AddressCodec().StringToBytes(req.Voter)
			if err != nil {
				return false, err
			}

			has, err := q.k.Votes.Has(ctx, collections.Join(p.Id, sdk.AccAddress(voter)))
			// if no error, vote found, matchVoter = true
			matchVoter = err == nil && has
		}

		// match depositor (if supplied)
		if len(req.Depositor) > 0 {
			depositor, err := q.k.authKeeper.AddressCodec().StringToBytes(req.Depositor)
			if err != nil {
				return false, err
			}
			has, err := q.k.Deposits.Has(ctx, collections.Join(p.Id, sdk.AccAddress(depositor)))
			// if no error, deposit found, matchDepositor = true
			matchDepositor = err == nil && has
		}

		// if all match, append to results
		if matchVoter && matchDepositor && matchStatus {
			return true, nil
		}
		// continue to next item, do not include because we're appending results above.
		return false, nil
	}, func(_ uint64, value govtypesv1.Proposal) (*govtypesv1.Proposal, error) {
		return &value, nil
	})

	if err != nil && !errors.IsOf(err, collections.ErrInvalidIterator) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &govtypesv1.QueryProposalsResponse{Proposals: filteredProposals, Pagination: pageRes}, nil
}

// Vote returns Voted information based on proposalID, voterAddr
func (q queryServer) Vote(ctx context.Context, req *govtypesv1.QueryVoteRequest) (*govtypesv1.QueryVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	if req.Voter == "" {
		return nil, status.Error(codes.InvalidArgument, "empty voter address")
	}

	voter, err := q.k.authKeeper.AddressCodec().StringToBytes(req.Voter)
	if err != nil {
		return nil, err
	}
	vote, err := q.k.Votes.Get(ctx, collections.Join(req.ProposalId, sdk.AccAddress(voter)))
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.InvalidArgument,
				"voter: %v not found for proposal: %v", req.Voter, req.ProposalId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &govtypesv1.QueryVoteResponse{Vote: &vote}, nil
}

// Votes returns single proposal's votes
func (q queryServer) Votes(ctx context.Context, req *govtypesv1.QueryVotesRequest) (*govtypesv1.QueryVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	votes, pageRes, err := query.CollectionPaginate(ctx, q.k.Votes, req.Pagination, func(_ collections.Pair[uint64, sdk.AccAddress], value govtypesv1.Vote) (vote *govtypesv1.Vote, err error) {
		return &value, nil
	}, query.WithCollectionPaginationPairPrefix[uint64, sdk.AccAddress](req.ProposalId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &govtypesv1.QueryVotesResponse{Votes: votes, Pagination: pageRes}, nil
}

// Deposit queries single deposit information based on proposalID, depositAddr.
func (q queryServer) Deposit(ctx context.Context, req *govtypesv1.QueryDepositRequest) (*govtypesv1.QueryDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	if req.Depositor == "" {
		return nil, status.Error(codes.InvalidArgument, "empty depositor address")
	}

	depositor, err := q.k.authKeeper.AddressCodec().StringToBytes(req.Depositor)
	if err != nil {
		return nil, err
	}
	deposit, err := q.k.Deposits.Get(ctx, collections.Join(req.ProposalId, sdk.AccAddress(depositor)))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &govtypesv1.QueryDepositResponse{Deposit: &deposit}, nil
}

// Deposits returns single proposal's all deposits
func (q queryServer) Deposits(ctx context.Context, req *govtypesv1.QueryDepositsRequest) (*govtypesv1.QueryDepositsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	var deposits []*govtypesv1.Deposit
	deposits, pageRes, err := query.CollectionPaginate(ctx, q.k.Deposits, req.Pagination, func(_ collections.Pair[uint64, sdk.AccAddress], deposit govtypesv1.Deposit) (*govtypesv1.Deposit, error) {
		return &deposit, nil
	}, query.WithCollectionPaginationPairPrefix[uint64, sdk.AccAddress](req.ProposalId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &govtypesv1.QueryDepositsResponse{Deposits: deposits, Pagination: pageRes}, nil
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
// func (q queryServer) Params(ctx context.Context, request *govtypesv1.QueryParamsRequest) (*typesgovtypesv1.QueryParamsResponse, error) {
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
