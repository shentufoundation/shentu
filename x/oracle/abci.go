package oracle

import (
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.FinalizeMatureWithdraws(ctx)
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	closingTaskIDs := k.GetInvalidTaskIDs(ctx)
	toAggTaskIDs := append(k.GetShortcutTasks(ctx), closingTaskIDs...)
	for _, taskID := range toAggTaskIDs {
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
		remainingBounty := k.HandleRemainingBounty(ctx, task)

		if distributeErr != nil {
			continue
		}

		switch task := task.(type) {
		case *types.Task:
			EmitEventsForTask(ctx, task, remainingBounty)
		case *types.TxTask:
			EmitEventsForTxTask(ctx, task, remainingBounty)
		}
	}
	k.DeleteClosingTaskIDs(ctx)
	k.DeleteShortcutTasks(ctx)
	k.DeleteExpiredTasks(ctx)
}

func EmitEventsForTask(ctx sdk.Context, task *types.Task, remainingBounty string) {
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
			sdk.NewAttribute("remaining_bounty", remainingBounty),
		),
	)
}

func EmitEventsForTxTask(ctx sdk.Context, task *types.TxTask, remainingBounty string) {
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
			sdk.NewAttribute("remaining_bounty", remainingBounty),
		),
	)
}
