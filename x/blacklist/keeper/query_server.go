package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/shentufoundation/shentu/v2/x/blacklist/types"
)

type queryServer struct {
	k Keeper
}

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the blacklist QueryServer
// interface for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k: k}
}

// Entries implements the Query/Entries RPC method.
func (qs queryServer) Entries(ctx context.Context, req *types.QueryEntriesRequest) (*types.QueryEntriesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	entries, pageRes, err := query.CollectionPaginate(
		ctx,
		qs.k.Entries,
		req.Pagination,
		func(_ string, entry types.Entry) (types.Entry, error) {
			return entry, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if entries == nil {
		entries = []types.Entry{}
	}

	return &types.QueryEntriesResponse{Entries: entries, Pagination: pageRes}, nil
}

// Entry implements the Query/Entry RPC method.
func (qs queryServer) Entry(ctx context.Context, req *types.QueryEntryRequest) (*types.QueryEntryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	entry, err := qs.k.GetEntry(ctx, req.Address)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "address %s is not blacklisted", req.Address)
	}

	return &types.QueryEntryResponse{Entry: entry}, nil
}
