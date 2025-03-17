package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"

	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ types.QueryServer = queryServer{}

type queryServer struct{ k *Keeper }

func NewQueryServer(k Keeper) types.QueryServer {
	return queryServer{k: &k}
}

// Programs implements the Query/Programs gRPC method
func (q queryServer) Programs(c context.Context, req *types.QueryProgramsRequest) (*types.QueryProgramsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	programs, pageRes, err := query.CollectionFilteredPaginate(c, q.k.Programs, req.Pagination, func(key string, p types.Program) (include bool, err error) {
		return true, nil
	}, func(_ string, value types.Program) (*types.Program, error) {
		return &value, nil
	})

	if err != nil && !errors.IsOf(err, collections.ErrInvalidIterator) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProgramsResponse{
		Programs:   programs,
		Pagination: pageRes,
	}, nil
}

// Program returns program details based on ProgramId
func (q queryServer) Program(c context.Context, req *types.QueryProgramRequest) (*types.QueryProgramResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if len(req.ProgramId) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "program-id can not less than 0")
	}

	program, err := q.k.Programs.Get(c, req.ProgramId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "program %s doesn't exist", req.ProgramId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProgramResponse{Program: &program}, nil
}

func (q queryServer) Findings(c context.Context, req *types.QueryFindingsRequest) (*types.QueryFindingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if len(req.ProgramId) == 0 && len(req.SubmitterAddress) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	findings, pageRes, err := query.CollectionFilteredPaginate(c, q.k.Findings, req.Pagination, func(key string, f types.Finding) (include bool, err error) {
		switch {
		case len(req.ProgramId) != 0 && len(req.SubmitterAddress) != 0:
			return f.ProgramId == req.ProgramId && f.SubmitterAddress == req.SubmitterAddress, nil
		case len(req.ProgramId) != 0:
			return f.ProgramId == req.ProgramId, nil
		case len(req.SubmitterAddress) != 0:
			return f.SubmitterAddress == req.SubmitterAddress, nil
		default:
			return true, nil
		}
	}, func(_ string, value types.Finding) (*types.Finding, error) {
		return &value, nil
	})

	if err != nil && !errors.IsOf(err, collections.ErrInvalidIterator) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFindingsResponse{
		Findings:   findings,
		Pagination: pageRes,
	}, nil
}

func (q queryServer) Finding(c context.Context, req *types.QueryFindingRequest) (*types.QueryFindingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if len(req.FindingId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "finding-id can not be 0")
	}

	finding, err := q.k.Findings.Get(c, req.FindingId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "finding %s doesn't exist", req.FindingId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFindingResponse{Finding: &finding}, nil
}

func (q queryServer) FindingFingerprint(c context.Context, req *types.QueryFindingFingerprintRequest) (*types.QueryFindingFingerprintResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if len(req.FindingId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "finding-id can not be 0")
	}

	finding, err := q.k.Findings.Get(c, req.FindingId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "finding %s doesn't exist", req.FindingId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	findingFingerPrintHash := q.k.GetFindingFingerprintHash(&finding)
	return &types.QueryFindingFingerprintResponse{Fingerprint: findingFingerPrintHash}, nil
}

func (q queryServer) ProgramFingerprint(c context.Context, req *types.QueryProgramFingerprintRequest) (*types.QueryProgramFingerprintResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if len(req.ProgramId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "program-id can not be 0")
	}

	program, err := q.k.Programs.Get(c, req.ProgramId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "program %s doesn't exist", req.ProgramId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	programFingerPrintHash := q.k.GetProgramFingerprintHash(&program)
	return &types.QueryProgramFingerprintResponse{Fingerprint: programFingerPrintHash}, nil
}

func (q queryServer) Theorems(c context.Context, req *types.QueryTheoremsRequest) (*types.QueryTheoremsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	filteredTheorems, pageRes, err := query.CollectionFilteredPaginate(c, q.k.Theorems, req.Pagination, func(key uint64, t types.Theorem) (include bool, err error) {
		return true, nil
	}, func(_ uint64, value types.Theorem) (*types.Theorem, error) {
		return &value, nil
	})

	if err != nil && !errors.IsOf(err, collections.ErrInvalidIterator) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTheoremsResponse{Theorems: filteredTheorems, Pagination: pageRes}, nil
}

func (q queryServer) Theorem(c context.Context, req *types.QueryTheoremRequest) (*types.QueryTheoremResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.TheoremId == 0 {
		return nil, status.Error(codes.InvalidArgument, "theorem id can not be 0")
	}

	theorem, err := q.k.Theorems.Get(c, req.TheoremId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "theorem %d doesn't exist", req.TheoremId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTheoremResponse{Theorem: &theorem}, nil
}

func (q queryServer) Proof(c context.Context, req *types.QueryProofRequest) (*types.QueryProofResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProofId == "" {
		return nil, status.Error(codes.InvalidArgument, "proof id can not be empty")
	}

	proof, err := q.k.Proofs.Get(c, req.ProofId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "proof %s doesn't exist", req.ProofId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProofResponse{Proof: &proof}, nil
}

func (q queryServer) Proofs(c context.Context, req *types.QueryProofsRequest) (*types.QueryProofsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.TheoremId == 0 {
		return nil, status.Error(codes.InvalidArgument, "theorem id can not be 0")
	}

	var (
		proofs  []*types.Proof
		pageRes *query.PageResponse
		err     error
	)

	proofs, pageRes, err = query.CollectionPaginate(c, q.k.ProofsByTheorem,
		req.Pagination, func(key collections.Pair[uint64, string], _ []byte) (*types.Proof, error) {
			proof, err := q.k.Proofs.Get(c, key.K2())
			if err != nil {
				return nil, err
			}
			return &proof, nil
		}, query.WithCollectionPaginationPairPrefix[uint64, string](req.TheoremId),
	)
	if err != nil && !errors.IsOf(err, collections.ErrInvalidIterator) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProofsResponse{
		Proofs:     proofs,
		Pagination: pageRes,
	}, nil
}

func (q queryServer) Reward(c context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	addr, err := q.k.authKeeper.AddressCodec().StringToBytes(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	reward, err := q.k.Rewards.Get(c, addr)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "reward for address %s doesn't exist", req.Address)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRewardsResponse{Rewards: reward.Reward}, nil
}

func (q queryServer) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	params, err := q.k.Params.Get(c)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParamsResponse{Params: &params}, nil
}
