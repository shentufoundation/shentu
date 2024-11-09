package keeper

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

var _ types.QueryServer = Keeper{}

// Operator queries an operator based on its address.
func (k Keeper) Operator(c context.Context, req *types.QueryOperatorRequest) (*types.QueryOperatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return nil, err
	}

	return &types.QueryOperatorResponse{Operator: operator}, nil
}

// Operators queries all operators.
func (k Keeper) Operators(c context.Context, req *types.QueryOperatorsRequest) (*types.QueryOperatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryOperatorsResponse{Operators: k.GetAllOperators(ctx)}, nil
}

// Withdraws queries all withdraws.
func (k Keeper) Withdraws(c context.Context, req *types.QueryWithdrawsRequest) (*types.QueryWithdrawsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryWithdrawsResponse{Withdraws: k.GetAllWithdraws(ctx)}, nil
}

// Task queries a task given its contract and function.
func (k Keeper) Task(c context.Context, req *types.QueryTaskRequest) (*types.QueryTaskResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	task, err := k.GetTask(ctx, types.NewTaskID(req.Contract, req.Function))
	if err != nil {
		return nil, err
	}

	if smartContractTask, ok := task.(*types.Task); ok {
		return &types.QueryTaskResponse{Task: *smartContractTask}, nil
	}
	return nil, types.ErrFailedToCastTask
}

// TxTask queries a tx task given its tx hash.
func (k Keeper) TxTask(c context.Context, req *types.QueryTxTaskRequest) (*types.QueryTxTaskResponse, error) {
	if req == nil || len(req.AtxHash) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	taskID, err := types.NewTxTaskID(req.AtxHash)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	task, err := k.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if txTask, ok := task.(*types.TxTask); ok {
		return &types.QueryTxTaskResponse{Task: *txTask}, nil
	}
	return nil, types.ErrFailedToCastTask
}

// Response queries a response based on its task contract, task function,
// and operator address.
func (k Keeper) Response(c context.Context, req *types.QueryResponseRequest) (*types.QueryResponseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	task, err := k.GetTask(ctx, types.NewTaskID(req.Contract, req.Function))
	if err != nil {
		return nil, err
	}

	for _, response := range task.GetResponses() {
		if response.Operator == req.OperatorAddress {
			return &types.QueryResponseResponse{Response: response}, nil
		}
	}
	return &types.QueryResponseResponse{}, fmt.Errorf("there is no response from this operator")
}

// TxResponse queries a tx response based on its tx hash,
// and operator address.
func (k Keeper) TxResponse(c context.Context, req *types.QueryTxResponseRequest) (*types.QueryTxResponseResponse, error) {
	if req == nil || len(req.AtxHash) == 0 || len(req.OperatorAddress) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	taskID, err := types.NewTxTaskID(req.AtxHash)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	task, err := k.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	for _, response := range task.GetResponses() {
		if response.Operator == req.OperatorAddress {
			return &types.QueryTxResponseResponse{Response: response}, nil
		}
	}
	return &types.QueryTxResponseResponse{}, fmt.Errorf("there is no response from this operator")
}

// Params queries all params
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	taskParams := k.GetTaskParams(ctx)
	poolParams := k.GetLockedPoolParams(ctx)
	return &types.QueryParamsResponse{TaskParams: taskParams, PoolParams: poolParams}, nil
}
