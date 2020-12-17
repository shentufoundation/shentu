package keeper

import (
	"context"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the shield MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateOperator(goCtx context.Context, msg *types.MsgCreateOperator) (*types.MsgCreateOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	proposerAddr, err := sdk.AccAddressFromBech32(msg.Proposer)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.CreateOperator(ctx, operatorAddr, msg.Collateral, proposerAddr, msg.Name); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateOperator,
			sdk.NewAttribute("operator", msg.Address),
			sdk.NewAttribute("operator_name", msg.Name),
			sdk.NewAttribute("collateral", msg.Collateral.String()),
		),
	})

	return &types.MsgCreateOperatorResponse{}, nil
}

func (k msgServer) RemoveOperator(goCtx context.Context, msg *types.MsgRemoveOperator) (*types.MsgRemoveOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.RemoveOperator(ctx, addr); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveOperator,
			sdk.NewAttribute("operator", msg.Address),
		),
	})

	return &types.MsgRemoveOperatorResponse{}, nil
}

func (k msgServer) AddCollateral(goCtx context.Context, msg *types.MsgAddCollateral) (*types.MsgAddCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.AddCollateral(ctx, addr, msg.CollateralIncrement); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddCollateral,
			sdk.NewAttribute("operator", msg.Address),
			sdk.NewAttribute("collateral_increment", msg.CollateralIncrement.String()),
		),
	})

	return &types.MsgAddCollateralResponse{}, nil
}

func (k msgServer) ReduceCollateral(goCtx context.Context, msg *types.MsgReduceCollateral) (*types.MsgReduceCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.ReduceCollateral(ctx, addr, msg.CollateralDecrement); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReduceCollateral,
			sdk.NewAttribute("operator", msg.Address),
			sdk.NewAttribute("collateral_decrement", msg.CollateralDecrement.String()),
		),
	})

	return &types.MsgReduceCollateralResponse{}, nil
}

func (k msgServer) WithdrawReward(goCtx context.Context, msg *types.MsgWithdrawReward) (*types.MsgWithdrawRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	reward, err := k.Keeper.WithdrawAllReward(ctx, addr)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawReward,
			sdk.NewAttribute("operator", msg.Address),
			sdk.NewAttribute("reward", reward.String()),
		),
	})

	return &types.MsgWithdrawRewardResponse{}, nil
}

func (k msgServer) CreateTask(goCtx context.Context, msg *types.MsgCreateTask) (*types.MsgCreateTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	taskParams := k.Keeper.GetTaskParams(ctx)
	var windowSize int64
	if msg.Wait == 0 {
		windowSize = taskParams.AggregationWindow
	} else {
		windowSize = msg.Wait
	}
	var expiration time.Time
	if msg.ValidDuration.Microseconds() == 0 {
		expiration = ctx.BlockTime().Add(taskParams.ExpirationDuration)
	} else {
		expiration = ctx.BlockTime().Add(msg.ValidDuration)
	}

	if err := k.Keeper.CreateTask(ctx, msg.Contract, msg.Function, msg.Bounty, msg.Description,
		expiration, creatorAddr, windowSize); err != nil {
		return nil, err
	}

	createTaskEvent := sdk.NewEvent(
		types.EventTypeCreateTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("bounty", msg.Bounty.String()),
		sdk.NewAttribute("description", msg.Description),
		sdk.NewAttribute("expiration", expiration.String()),
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("windowSize", strconv.FormatInt(windowSize, 10)),
		sdk.NewAttribute("closingHeight", strconv.FormatInt(ctx.BlockHeight()+windowSize, 10)),
	)
	ctx.EventManager().EmitEvent(createTaskEvent)

	return &types.MsgCreateTaskResponse{}, nil
}

func (k msgServer) TaskResponse(goCtx context.Context, msg *types.MsgTaskResponse) (*types.MsgTaskResponseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.AccAddressFromBech32(msg.Operator)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.RespondToTask(ctx, msg.Contract, msg.Function, msg.Score, operatorAddr); err != nil {
		return nil, err
	}

	respondToTaskEvent := sdk.NewEvent(
		types.EventTypeRespondToTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("score", strconv.FormatInt(msg.Score, 10)),
		sdk.NewAttribute("operator", msg.Operator),
	)
	ctx.EventManager().EmitEvent(respondToTaskEvent)

	return &types.MsgTaskResponseResponse{}, nil
}

func (k msgServer) InquiryTask(goCtx context.Context, msg *types.MsgInquiryTask) (*types.MsgInquiryTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	task, err := k.Keeper.GetTask(ctx, msg.Contract, msg.Function)
	if err != nil {
		return nil, err
	}

	InquiryTaskEvent := sdk.NewEvent(
		types.EventTypeInquireTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("txhash", msg.TxHash),
		sdk.NewAttribute("inquirer", msg.Inquirer),
		sdk.NewAttribute("result", strconv.FormatUint(task.Result.Uint64(), 10)),
		sdk.NewAttribute("expiration", task.Expiration.String()),
	)
	ctx.EventManager().EmitEvent(InquiryTaskEvent)

	return &types.MsgInquiryTaskResponse{}, nil
}

func (k msgServer) DeleteTask(goCtx context.Context, msg *types.MsgDeleteTask) (*types.MsgDeleteTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	deleterAddr, err := sdk.AccAddressFromBech32(msg.Deleter)
	if err != nil {
		return nil, err
	}

	if err := k.RemoveTask(ctx, msg.Contract, msg.Function, msg.Force, deleterAddr); err != nil {
		return nil, err
	}

	DeleteTaskEvent := sdk.NewEvent(
		types.EventTypeDeleteTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("creator", msg.Deleter),
		sdk.NewAttribute("expired", strconv.FormatBool(msg.Force)),
	)
	ctx.EventManager().EmitEvent(DeleteTaskEvent)

	return &types.MsgDeleteTaskResponse{}, nil
}
