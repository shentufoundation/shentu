package oracle

import (
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

// NewHandler returns a handler for oracle type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateOperator:
			return handleMsgCreateOperator(ctx, k, msg)
		case types.MsgRemoveOperator:
			return handleMsgRemoveOperator(ctx, k, msg)
		case types.MsgAddCollateral:
			return handleMsgAddCollateral(ctx, k, msg)
		case types.MsgReduceCollateral:
			return handleMsgReduceCollateral(ctx, k, msg)
		case types.MsgWithdrawReward:
			return handleMsgWithdrawReward(ctx, k, msg)
		case types.MsgCreateTask:
			return handleMsgCreateTask(ctx, k, msg)
		case types.MsgTaskResponse:
			return handleMsgTaskResponse(ctx, k, msg)
		case types.MsgInquiryTask:
			return handleMsgInquiryTask(ctx, k, msg)
		case types.MsgDeleteTask:
			return handleMsgDeleteTask(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgCreateOperator(ctx sdk.Context, k Keeper, msg types.MsgCreateOperator) (*sdk.Result, error) {
	if err := k.CreateOperator(ctx, msg.Address, msg.Collateral, msg.Proposer, msg.Name); err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}

func handleMsgRemoveOperator(ctx sdk.Context, k Keeper, msg types.MsgRemoveOperator) (*sdk.Result, error) {
	if err := k.RemoveOperator(ctx, msg.Address); err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}

func handleMsgAddCollateral(ctx sdk.Context, k Keeper, msg types.MsgAddCollateral) (*sdk.Result, error) {
	if err := k.AddCollateral(ctx, msg.Address, msg.CollateralIncrement); err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}

func handleMsgReduceCollateral(ctx sdk.Context, k Keeper, msg types.MsgReduceCollateral) (*sdk.Result, error) {
	if err := k.ReduceCollateral(ctx, msg.Address, msg.CollateralDecrement); err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}

func handleMsgWithdrawReward(ctx sdk.Context, k Keeper, msg types.MsgWithdrawReward) (*sdk.Result, error) {
	if err := k.WithdrawAllReward(ctx, msg.Address); err != nil {
		return nil, err
	}

	WithdrawEvent := sdk.NewEvent(
		types.EventTypeWithdraw,
		sdk.NewAttribute("address", msg.Address.String()),
	)
	ctx.EventManager().EmitEvent(WithdrawEvent)

	return &sdk.Result{}, nil
}

func handleMsgCreateTask(ctx sdk.Context, k Keeper, msg types.MsgCreateTask) (*sdk.Result, error) {
	taskParams := k.GetTaskParams(ctx)
	var windowSize int64
	if msg.Wait == 0 {
		windowSize = taskParams.AggregationWindow
	} else {
		windowSize = msg.Wait
	}
	var expiration time.Time
	if msg.ValidDuration.Microseconds() == 0 {
		expiration = msg.Now.Add(taskParams.ExpirationDuration)
	} else {
		expiration = msg.Now.Add(msg.ValidDuration)
	}

	if err := k.CreateTask(ctx, msg.Contract, msg.Function, msg.Bounty, msg.Description,
		expiration, msg.Creator, windowSize); err != nil {
		return nil, err
	}

	createTaskEvent := sdk.NewEvent(
		types.EventTypeCreateTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("bounty", msg.Bounty.String()),
		sdk.NewAttribute("description", msg.Description),
		sdk.NewAttribute("expiration", expiration.String()),
		sdk.NewAttribute("creator", msg.Creator.String()),
		sdk.NewAttribute("windowSize", strconv.FormatInt(windowSize, 10)),
		sdk.NewAttribute("closingHeight", strconv.FormatInt(ctx.BlockHeight()+windowSize, 10)),
	)
	ctx.EventManager().EmitEvent(createTaskEvent)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgTaskResponse(ctx sdk.Context, k Keeper, msg types.MsgTaskResponse) (*sdk.Result, error) {
	if err := k.RespondToTask(ctx, msg.Contract, msg.Function, msg.Score, msg.Operator); err != nil {
		return nil, err
	}
	respondToTaskEvent := sdk.NewEvent(
		types.EventTypeRespondToTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("score", strconv.FormatInt(msg.Score, 10)),
		sdk.NewAttribute("operator", msg.Operator.String()),
	)
	ctx.EventManager().EmitEvent(respondToTaskEvent)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgInquiryTask(ctx sdk.Context, k Keeper, msg types.MsgInquiryTask) (*sdk.Result, error) {
	task, err := k.GetTask(ctx, msg.Contract, msg.Function)
	if err != nil {
		return nil, err
	}
	InquiryTaskEvent := sdk.NewEvent(
		types.EventTypeInquireTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("txhash", msg.TxHash),
		sdk.NewAttribute("inquirer", msg.Inquirer.String()),
		sdk.NewAttribute("result", strconv.FormatUint(task.Result.Uint64(), 10)),
		sdk.NewAttribute("expiration", task.Expiration.String()),
	)
	ctx.EventManager().EmitEvent(InquiryTaskEvent)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgDeleteTask(ctx sdk.Context, k Keeper, msg types.MsgDeleteTask) (*sdk.Result, error) {
	if err := k.RemoveTask(ctx, msg.Contract, msg.Function, msg.Force, msg.Deleter); err != nil {
		return nil, err
	}
	DeleteTaskEvent := sdk.NewEvent(
		types.EventTypeDeleteTask,
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("function", msg.Function),
		sdk.NewAttribute("creator", msg.Deleter.String()),
		sdk.NewAttribute("expired", strconv.FormatBool(msg.Force)),
	)
	ctx.EventManager().EmitEvent(DeleteTaskEvent)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
