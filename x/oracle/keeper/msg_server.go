package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
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
			types.TypeMsgCreateOperator,
			sdk.NewAttribute("operator", msg.Address),
			sdk.NewAttribute("operator_name", msg.Name),
			sdk.NewAttribute("collateral", msg.Collateral.String()),
		),
	})

	return &types.MsgCreateOperatorResponse{}, nil
}

func (k msgServer) RemoveOperator(goCtx context.Context, msg *types.MsgRemoveOperator) (*types.MsgRemoveOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.RemoveOperator(ctx, msg.Address, msg.Proposer); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgRemoveOperator,
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
			types.TypeMsgAddCollateral,
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
			types.TypeMsgReduceCollateral,
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
			types.TypeMsgWithdrawReward,
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

	smartContractTask := types.NewTask(
		msg.Contract, msg.Function, ctx.BlockHeight(),
		msg.Bounty, msg.Description, expiration,
		creatorAddr, ctx.BlockHeight()+windowSize, windowSize)
	if err := k.Keeper.CreateTask(ctx, creatorAddr, &smartContractTask); err != nil {
		return nil, err
	}

	createTaskEvent := sdk.NewEvent(
		types.TypeMsgCreateTask,
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

	if err := k.Keeper.RespondToTask(ctx, types.NewTaskID(msg.Contract, msg.Function), msg.Score, operatorAddr); err != nil {
		return nil, err
	}

	respondToTaskEvent := sdk.NewEvent(
		types.TypeMsgRespondToTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("score", strconv.FormatInt(msg.Score, 10)),
		sdk.NewAttribute("operator", msg.Operator),
	)
	ctx.EventManager().EmitEvent(respondToTaskEvent)

	return &types.MsgTaskResponseResponse{}, nil
}

func (k msgServer) DeleteTask(goCtx context.Context, msg *types.MsgDeleteTask) (*types.MsgDeleteTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	deleterAddr, _ := sdk.AccAddressFromBech32(msg.From)

	if err := k.RemoveTask(ctx, types.NewTaskID(msg.Contract, msg.Function), msg.Force, deleterAddr); err != nil {
		return nil, err
	}

	DeleteTaskEvent := sdk.NewEvent(
		types.TypeMsgDeleteTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("creator", msg.From),
		sdk.NewAttribute("expired", strconv.FormatBool(msg.Force)),
	)
	ctx.EventManager().EmitEvent(DeleteTaskEvent)

	return &types.MsgDeleteTaskResponse{}, nil
}

func (k msgServer) CreateTxTask(goCtx context.Context, msg *types.MsgCreateTxTask) (*types.MsgCreateTxTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creatorAddr, _ := sdk.AccAddressFromBech32(msg.Creator)
	if msg.ValidTime.Before(ctx.BlockTime()) {
		return nil, types.ErrOverdueValidTime
	}

	hashByte := sha256.Sum256(msg.TxBytes)
	hash := base64.StdEncoding.EncodeToString(hashByte[:])

	txTask, err := k.BuildTxTask(ctx, hashByte[:], msg.Creator, msg.Bounty, msg.ValidTime)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.CreateTask(ctx, creatorAddr, txTask); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgCreateTxTask,
			sdk.NewAttribute("tx_hash", hash),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("chain_id", msg.ChainId),
			sdk.NewAttribute("bounty", msg.Bounty.String()),
			sdk.NewAttribute("valid_time", msg.ValidTime.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})

	return &types.MsgCreateTxTaskResponse{
		TxHash: hashByte[:],
	}, nil
}

func (k msgServer) TxTaskResponse(goCtx context.Context, msg *types.MsgTxTaskResponse) (*types.MsgTxTaskResponseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, _ := sdk.AccAddressFromBech32(msg.Operator)
	if err := k.HandleNoneTxTaskForResponse(ctx, msg.TxHash); err != nil {
		return nil, err
	}

	if err := k.RespondToTask(ctx, msg.TxHash, msg.Score, operatorAddr); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgRespondToTxTask,
			sdk.NewAttribute("operator", msg.Operator),
			sdk.NewAttribute("score", strconv.FormatInt(msg.Score, 10)),
			sdk.NewAttribute("txHash", base64.StdEncoding.EncodeToString(msg.TxHash)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Operator),
		),
	})

	return &types.MsgTxTaskResponseResponse{}, nil
}

func (k msgServer) DeleteTxTask(goCtx context.Context, msg *types.MsgDeleteTxTask) (*types.MsgDeleteTxTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	deleterAddr, _ := sdk.AccAddressFromBech32(msg.From)
	if err := k.RemoveTask(ctx, msg.TxHash, true, deleterAddr); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgDeleteTxTask,
			sdk.NewAttribute("deleter", msg.From),
			sdk.NewAttribute("txHash", base64.StdEncoding.EncodeToString(msg.TxHash)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgDeleteTxTaskResponse{}, nil
}
