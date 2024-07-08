package gov

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
)

func removeInactiveProposals(ctx sdk.Context, k keeper.Keeper) {
	logger := k.Logger(ctx)

	// delete dead proposals from store and returns theirs deposits. A proposal is dead when it's inactive and didn't get enough deposit on time to get into voting phase.
	k.IterateInactiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal govtypesv1.Proposal) bool {
		k.DeleteProposal(ctx, proposal.Id)
		k.RefundAndDeleteDeposits(ctx, proposal.Id)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govtypes.EventTypeInactiveProposal,
				sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.Id)),
				sdk.NewAttribute(govtypes.AttributeKeyProposalResult, govtypes.AttributeValueProposalDropped),
			),
		)

		logger.Info(
			"proposal did not meet minimum deposit; deleted",
			"proposal", proposal.Id,
			"min_deposit", sdk.NewCoins(k.GetParams(ctx).MinDeposit...).String(),
			"total_deposit", sdk.NewCoins(proposal.TotalDeposit...).String(),
		)

		return false
	})
}

// fetch active proposals whose voting periods have ended (are passed the block time)
func processActiveProposal(ctx sdk.Context, k keeper.Keeper, proposal govtypesv1.Proposal) bool {
	var (
		tagValue, logMsg string
		pass, veto       bool
		tallyResults     govtypesv1.TallyResult
	)
	logger := k.Logger(ctx)

	if k.CertifierVoteIsRequired(proposal) && !k.GetCertifierVoted(ctx, proposal.Id) {
		var endVoting bool
		pass, endVoting, tallyResults = keeper.SecurityTally(ctx, k, proposal)
		if !endVoting {
			// Skip the rest of this iteration, because the proposal needs to go
			// through the validator voting period now.
			k.SetCertifierVoted(ctx, proposal.Id)
			k.DeleteAllVotes(ctx, proposal.Id)
			k.ActivateVotingPeriod(ctx, proposal)
			return false
		}
	} else {
		pass, veto, tallyResults = k.Tally(ctx, proposal)
	}

	if veto {
		k.DeleteAndBurnDeposits(ctx, proposal.Id)
	} else {
		k.RefundAndDeleteDeposits(ctx, proposal.Id)
	}

	if pass {
		var (
			idx    int
			events sdk.Events
			msg    sdk.Msg
		)

		// attempt to execute all messages within the passed proposal
		// Messages may mutate state thus we use a cached context. If one of
		// the handlers fails, no state mutation is written and the error
		// message is logged.
		cacheCtx, writeCache := ctx.CacheContext()
		messages, err := proposal.GetMsgs()
		if err == nil {
			for idx, msg = range messages {
				handler := k.Router().Handler(msg)

				var res *sdk.Result
				res, err = handler(cacheCtx, msg)
				if err != nil {
					break
				}

				events = append(events, res.GetEvents()...)
			}
		}

		// `err == nil` when all handlers passed.
		// Or else, `idx` and `err` are populated with the msg index and error.
		if err == nil {
			proposal.Status = govtypesv1.StatusPassed
			tagValue = govtypes.AttributeValueProposalPassed
			logMsg = "passed"

			// write state to the underlying multi-store
			writeCache()

			// propagate the msg events to the current context
			ctx.EventManager().EmitEvents(events)
		} else {
			proposal.Status = govtypesv1.StatusFailed
			tagValue = govtypes.AttributeValueProposalFailed
			logMsg = fmt.Sprintf("passed, but msg %d (%s) failed on execution: %s", idx, sdk.MsgTypeURL(msg), err)
		}
	} else {
		proposal.Status = govtypesv1.StatusRejected
		tagValue = govtypes.AttributeValueProposalRejected
		logMsg = "rejected"
	}

	proposal.FinalTallyResult = &tallyResults

	k.SetProposal(ctx, proposal)
	k.RemoveFromActiveProposalQueue(ctx, proposal.Id, *proposal.VotingEndTime)

	logger.Info(
		"proposal tallied",
		"proposal", proposal.Id,
		"result", logMsg,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeActiveProposal,
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.Id)),
			sdk.NewAttribute(govtypes.AttributeKeyProposalResult, tagValue),
		),
	)
	return false
}

func processSecurityVote(ctx sdk.Context, k keeper.Keeper, proposal govtypesv1.Proposal) bool {
	var (
		tagValue, logMsg string
		pass             bool
		tallyResults     govtypesv1.TallyResult
	)
	logger := k.Logger(ctx)

	// Only process security proposals
	if !k.CertifierVoteIsRequired(proposal) {
		return false
	}
	// Only process proposals in the security voting period.
	if k.CertifierVoteIsRequired(proposal) && k.GetCertifierVoted(ctx, proposal.Id) {
		return false
	}

	var endVoting bool
	pass, endVoting, tallyResults = keeper.SecurityTally(ctx, k, proposal)
	if !pass {
		// Do nothing, because the proposal still has time before the voting period ends.
		return false
	}
	// Else: the proposal passed the certifier voting period.

	if endVoting {
		var (
			idx    int
			events sdk.Events
			msg    sdk.Msg
		)

		cacheCtx, writeCache := ctx.CacheContext()
		messages, err := proposal.GetMsgs()
		if err == nil {
			for idx, msg = range messages {
				handler := k.Router().Handler(msg)

				var res *sdk.Result
				res, err = handler(cacheCtx, msg)
				if err != nil {
					break
				}

				events = append(events, res.GetEvents()...)
			}
		}

		if err == nil {
			proposal.Status = govtypesv1.StatusPassed
			tagValue = govtypes.AttributeValueProposalPassed
			logMsg = "passed"

			// write state to the underlying multi-store
			writeCache()

			// propagate the msg events to the current context
			ctx.EventManager().EmitEvents(events)
		} else {
			proposal.Status = govtypesv1.StatusFailed
			tagValue = govtypes.AttributeValueProposalFailed
			logMsg = fmt.Sprintf("passed, but msg %d (%s) failed on execution: %s", idx, sdk.MsgTypeURL(msg), err)
		}

		proposal.FinalTallyResult = &tallyResults

		k.SetProposal(ctx, proposal)
		k.RemoveFromActiveProposalQueue(ctx, proposal.Id, *proposal.VotingEndTime)

		logger.Info(
			"proposal tallied",
			"proposal", proposal.Id,
			"result", logMsg,
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govtypes.EventTypeActiveProposal,
				sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.Id)),
				sdk.NewAttribute(govtypes.AttributeKeyProposalResult, tagValue),
			),
		)
	} else {
		// Activate validator voting period
		k.SetCertifierVoted(ctx, proposal.Id)
		k.DeleteAllVotes(ctx, proposal.Id)
		k.ActivateVotingPeriod(ctx, proposal)
	}
	return false
}

// EndBlocker is called every block, removes inactive proposals, tallies active
// proposals and deletes/refunds deposits.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// delete inactive proposal from store and its deposits
	removeInactiveProposals(ctx, k)

	// fetch active proposals whose voting periods have ended (are passed the
	// block time)
	k.IterateActiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal govtypesv1.Proposal) bool {
		return processActiveProposal(ctx, k, proposal)
	})

	// Iterate over all active proposals, regardless of end time, so that
	// security voting can end as soon as a passing threshold is met.
	k.IterateActiveProposalsQueue(ctx, time.Unix(common.MaxTimestamp, 0), func(proposal govtypesv1.Proposal) bool {
		return processSecurityVote(ctx, k, proposal)
	})
}
