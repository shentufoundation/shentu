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

	task, err := q.GetTask(ctx, types.NewTaskID(req.Contract, req.Function))
	if err != nil {
		return nil, err
	}

	if smartContractTask, ok := task.(*types.Task); ok {
		return &types.QueryTaskResponse{Task: *smartContractTask}, nil
	}
	return nil, types.ErrFailedToCastTask
}

// TxTask queries a tx task given its tx hash.
func (q Keeper) TxTask(c context.Context, req *types.QueryTxTaskRequest) (*types.QueryTxTaskResponse, error) {
	if req == nil || len(req.AtxHash) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	taskID, err := types.NewTxTaskID(req.AtxHash)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	task, err := q.GetTask(ctx, taskID)
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
func (q Keeper) Response(c context.Context, req *types.QueryResponseRequest) (*types.QueryResponseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	task, err := q.GetTask(ctx, types.NewTaskID(req.Contract, req.Function))
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
func (q Keeper) TxResponse(c context.Context, req *types.QueryTxResponseRequest) (*types.QueryTxResponseResponse, error) {
	if req == nil || len(req.AtxHash) == 0 || len(req.OperatorAddress) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	taskID, err := types.NewTxTaskID(req.AtxHash)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	task, err := q.GetTask(ctx, taskID)
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

// RemainingBounty This function retrieves the amount of bounty remaining for the task creator
func (q Keeper) RemainingBounty(c context.Context, req *types.QueryRemainingBountyRequest) (*types.QueryRemainingBountyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	remainingBounty, err := q.GetRemainingBounty(ctx, address)
	if err != nil {
		return nil, err
	}
	return &types.QueryRemainingBountyResponse{
		Bounty: remainingBounty,
	}, nil

}
