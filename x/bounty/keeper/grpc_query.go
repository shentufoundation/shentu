package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

var _ types.QueryServer = Keeper{}

// Programs implements the Query/Programs gRPC method
func (k Keeper) Programs(c context.Context, req *types.QueryProgramsRequest) (*types.QueryProgramsResponse, error) {
	var programs types.Programs

	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(c))
	programStore := prefix.NewStore(kvStore, types.ProgramKey)

	pageRes, err := query.FilteredPaginate(programStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var p types.Program
		if err := k.cdc.Unmarshal(value, &p); err != nil {
			return false, status.Error(codes.Internal, err.Error())
		}

		if accumulate {
			programs = append(programs, p)
		}

		return true, nil

	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProgramsResponse{
		Programs:   programs,
		Pagination: pageRes,
	}, nil
}

// Program returns program details based on ProgramId
func (k Keeper) Program(c context.Context, req *types.QueryProgramRequest) (*types.QueryProgramResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if len(req.ProgramId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "program-id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)
	program, found := k.GetProgram(ctx, req.ProgramId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "program %s doesn't exist", req.ProgramId)
	}

	return &types.QueryProgramResponse{Program: program}, nil
}

func (k Keeper) Findings(c context.Context, req *types.QueryFindingsRequest) (*types.QueryFindingsResponse, error) {
	var queryFindings types.Findings

	if len(req.ProgramId) == 0 && len(req.SubmitterAddress) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(c))
	programStore := prefix.NewStore(kvStore, types.FindingKey)

	pageRes, err := query.FilteredPaginate(programStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var finding types.Finding
		if err := k.cdc.Unmarshal(value, &finding); err != nil {
			return false, status.Error(codes.Internal, err.Error())
		}

		switch {
		case len(req.ProgramId) != 0 && len(req.SubmitterAddress) != 0:
			if finding.ProgramId == req.ProgramId && finding.SubmitterAddress == req.SubmitterAddress {
				queryFindings = append(queryFindings, finding)
			}
		case len(req.ProgramId) != 0:
			if finding.ProgramId == req.ProgramId {
				queryFindings = append(queryFindings, finding)

			}
		case len(req.SubmitterAddress) != 0:
			if finding.SubmitterAddress == req.SubmitterAddress {
				queryFindings = append(queryFindings, finding)
			}
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFindingsResponse{
		Findings:   queryFindings,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Finding(c context.Context, req *types.QueryFindingRequest) (*types.QueryFindingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if len(req.FindingId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "finding-id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)
	finding, found := k.GetFinding(ctx, req.FindingId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "finding %s doesn't exist", req.FindingId)
	}

	return &types.QueryFindingResponse{Finding: finding}, nil
}

func (k Keeper) FindingFingerprint(c context.Context, req *types.QueryFindingFingerprintRequest) (*types.QueryFindingFingerprintResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if len(req.FindingId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "finding-id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)
	finding, found := k.GetFinding(ctx, req.FindingId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "finding %s doesn't exist", req.FindingId)
	}

	findingFingerPrintHash := k.GetFindingFingerprintHash(&finding)
	return &types.QueryFindingFingerprintResponse{Fingerprint: findingFingerPrintHash}, nil
}

func (k Keeper) ProgramFingerprint(c context.Context, req *types.QueryProgramFingerprintRequest) (*types.QueryProgramFingerprintResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if len(req.ProgramId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "program-id can not be 0")
	}
	ctx := sdk.UnwrapSDKContext(c)
	program, found := k.GetProgram(ctx, req.ProgramId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "program %s doesn't exist", req.ProgramId)
	}
	programFingerPrintHash := k.GetProgramFingerprintHash(&program)
	return &types.QueryProgramFingerprintResponse{Fingerprint: programFingerPrintHash}, nil
}

func (k Keeper) AllTheorems(c context.Context, req *types.QueryTheoremsRequest) (*types.QueryTheoremsResponse, error) {
	panic("implement me")
}

func (k Keeper) Theorem(c context.Context, req *types.QueryTheoremRequest) (*types.QueryTheoremResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.TheoremId == 0 {
		return nil, status.Error(codes.InvalidArgument, "theorem id can not be 0")
	}

	theorem, err := k.Theorems.Get(c, req.TheoremId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "theorem %d doesn't exist", req.TheoremId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTheoremResponse{Theorem: theorem}, nil
}

func (k Keeper) Proof(c context.Context, req *types.QueryProofRequest) (*types.QueryProofResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProofId == "" {
		return nil, status.Error(codes.InvalidArgument, "proof id can not be empty")
	}

	proof, err := k.Proofs.Get(c, req.ProofId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "proof %s doesn't exist", req.ProofId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProofResponse{Proof: proof}, nil
}
