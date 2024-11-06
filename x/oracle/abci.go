package oracle

import (
	"context"
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func BeginBlocker(ctx context.Context, k keeper.Keeper) error {
	return k.FinalizeMatureWithdraws(ctx)
}

func EndBlocker(ctx context.Context, k keeper.Keeper) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	closingTaskIDs := k.GetInvalidTaskIDs(ctx)
	tasks, err := k.GetShortcutTasks(ctx)
	if err != nil {
		return err
	}
	tasks = append(tasks, closingTaskIDs...)
	for _, taskID := range tasks {
		err := k.Aggregate(ctx, taskID.Tid)
		if err != nil {
			continue
		}
		task, err := k.GetTask(ctx, taskID.Tid)
		if err != nil {
			continue
		}

		distributeErr := k.DistributeBounty(ctx, task)
		task, _ = k.GetTask(ctx, task.GetID())
		_ = k.RefundBounty(ctx, task)

		if distributeErr != nil {
			continue
		}

		switch task := task.(type) {
		case *types.Task:
			EmitEventsForTask(sdkCtx, task)
		case *types.TxTask:
			EmitEventsForTxTask(sdkCtx, task)
		}
	}
	err = k.DeleteClosingTaskIDs(ctx)
	if err != nil {
		return err
	}
	err = k.DeleteShortcutTasks(ctx)
	if err != nil {
		return err
	}
	err = k.DeleteExpiredTasks(ctx)
	if err != nil {
		return err
	}
	return nil
}

func EmitEventsForTask(ctx sdk.Context, task *types.Task) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAggTask,
			sdk.NewAttribute("contract", task.Contract),
			sdk.NewAttribute("function", task.Function),
			sdk.NewAttribute("begin_block_height", strconv.FormatInt(task.BeginBlock, 10)),
			sdk.NewAttribute("bounty", task.Bounty.String()),
			sdk.NewAttribute("description", task.Description),
			sdk.NewAttribute("expiration", task.Expiration.String()),
			sdk.NewAttribute("creator", task.Creator),
			sdk.NewAttribute("responses", task.Responses.String()),
			sdk.NewAttribute("result", task.Result.String()),
			sdk.NewAttribute("end_block_height", strconv.FormatInt(task.ExpireHeight, 10)),
			sdk.NewAttribute("status", task.Status.String()),
		),
	)
}

func EmitEventsForTxTask(ctx sdk.Context, task *types.TxTask) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAggTxTask,
			sdk.NewAttribute("atx_hash", hex.EncodeToString(task.AtxHash)),
			sdk.NewAttribute("score", strconv.FormatInt(task.Score, 10)),
			sdk.NewAttribute("status", task.Status.String()),
			sdk.NewAttribute("creator", task.Creator),
			sdk.NewAttribute("responses", task.Responses.String()),
			sdk.NewAttribute("expiration", task.Expiration.String()),
			sdk.NewAttribute("bounty", task.Bounty.String()),
		),
	)
}
