package keeper

import (
	"context"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Hosts(ctx context.Context, request *types.QueryHostsRequest) (*types.QueryHostsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Host(ctx context.Context, request *types.QueryHostRequest) (*types.QueryHostResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Programs(ctx context.Context, request *types.QueryProgramsRequest) (*types.QueryProgramsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Program(ctx context.Context, request *types.QueryProgramRequest) (*types.QueryProgramResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Findings(ctx context.Context, requests *types.QueryFindingsRequests) (*types.QueryFindingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Finding(ctx context.Context, request *types.QueryFindingRequest) (*types.QueryFindingResponse, error) {
	//TODO implement me
	panic("implement me")
}
