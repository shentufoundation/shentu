package keeper

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

var _ types.QueryServer = Keeper{}

// Operator queries an operator based on its address.
func (q Keeper) Operator(c context.Context, req *types.QueryOperatorRequest) (*types.QueryOperatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	operator, err := q.GetOperator(ctx, address)
	if err != nil {
		return nil, err
	}

	return &types.QueryOperatorResponse{Operator: operator}, nil
}

// Operators queries all operators.
func (q Keeper) Operators(c context.Context, req *types.QueryOperatorsRequest) (*types.QueryOperatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryOperatorsResponse{Operators: q.GetAllOperators(ctx)}, nil
}

// Withdraws queries all withdraws.
func (q Keeper) Withdraws(c context.Context, req *types.QueryWithdrawsRequest) (*types.QueryWithdrawsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryWithdrawsResponse{Withdraws: q.GetAllWithdraws(ctx)}, nil
}

// Task queries a task given its contract and function.
func (q Keeper) Task(c context.Context, req *types.QueryTaskRequest) (*types.QueryTaskResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	task, err := q.GetTask(ctx, req.Contract, req.Function)
	if err != nil {
		return nil, err
	}

	return &types.QueryTaskResponse{Task: task}, nil
}

// Response queries a response based on its task contract, task function,
// and operator address.
func (q Keeper) Response(c context.Context, req *types.QueryResponseRequest) (*types.QueryResponseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	task, err := q.GetTask(ctx, req.Contract, req.Function)
	if err != nil {
		return nil, err
	}

	for _, response := range task.Responses {
		if response.Operator == req.OperatorAddress {
			return &types.QueryResponseResponse{Response: response}, nil
		}
	}
	return &types.QueryResponseResponse{}, fmt.Errorf("there is no response from this operator")
}
