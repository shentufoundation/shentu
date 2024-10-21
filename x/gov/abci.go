package gov

import (
	"errors"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
)

// EndBlocker called every block, process inflation, update validator set.
// proposals and deletes/refunds deposits.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) error {
	logger := ctx.Logger().With("module", "x/"+govtypes.ModuleName)

	// delete dead proposals from store and returns theirs deposits.
	// A proposal is dead when it's inactive and didn't get enough deposit on time to get into voting phase.
	if err := removeInactiveProposals(ctx, &k, logger); err != nil {
		return err
	}

	// fetch active proposals whose voting periods have ended (are passed the block time)
	if err := processActiveProposal(ctx, &k, logger); err != nil {
		return err
	}

	//// Iterate over all active proposals, regardless of end time, so that
	//// security voting can end as soon as a passing threshold is met.
	if err := processSecurityVote(ctx, &k, logger); err != nil {
		return err
	}

	return nil
}

func removeInactiveProposals(ctx sdk.Context, k *keeper.Keeper, logger log.Logger) error {
	rng := collections.NewPrefixUntilPairRange[time.Time, uint64](ctx.BlockTime())
	err := k.InactiveProposalsQueue.Walk(ctx, rng, func(key collections.Pair[time.Time, uint64], _ uint64) (bool, error) {
		proposal, err := k.Proposals.Get(ctx, key.K2())
		if err != nil {
			// if the proposal has an encoding error, this means it cannot be processed by x/gov
			// this could be due to some types missing their registration
			// instead of returning an error (i.e, halting the chain), we fail the proposal
			if errors.Is(err, collections.ErrEncoding) {
				proposal.Id = key.K2()
				if err := failUnsupportedProposal(logger, ctx, k, proposal, err.Error(), false); err != nil {
					return false, err
				}

				if err = k.DeleteProposal(ctx, proposal.Id); err != nil {
					return false, err
				}

				return false, nil
			}

			return false, err
		}

		if err = k.DeleteProposal(ctx, proposal.Id); err != nil {
			return false, err
		}

		params, err := k.Params.Get(ctx)
		if err != nil {
			return false, err
		}
		if !params.BurnProposalDepositPrevote {
			err = k.RefundAndDeleteDeposits(ctx, proposal.Id) // refund deposit if proposal got removed without getting 100% of the proposal
		} else {
			err = k.DeleteAndBurnDeposits(ctx, proposal.Id) // burn the deposit if proposal got removed without getting 100% of the proposal
		}
		if err != nil {
			return false, err
		}

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
			"expedited", proposal.Expedited,
			"title", proposal.Title,
			"min_deposit", sdk.NewCoins(proposal.GetMinDepositFromParams(params)...).String(),
			"total_deposit", sdk.NewCoins(proposal.TotalDeposit...).String(),
		)

		return false, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func processActiveProposal(ctx sdk.Context, k *keeper.Keeper, logger log.Logger) error {
	var (
		tagValue, logMsg     string
		passes, burnDeposits bool
		tallyResults         govtypesv1.TallyResult
	)
	rng := collections.NewPrefixUntilPairRange[time.Time, uint64](ctx.BlockTime())
	err := k.ActiveProposalsQueue.Walk(ctx, rng, func(key collections.Pair[time.Time, uint64], _ uint64) (bool, error) {
		proposal, err := k.Proposals.Get(ctx, key.K2())
		if err != nil {
			// if the proposal has an encoding error, this means it cannot be processed by x/gov
			// this could be due to some types missing their registration
			// instead of returning an error (i.e, halting the chain), we fail the proposal
			if errors.Is(err, collections.ErrEncoding) {
				proposal.Id = key.K2()
				if err := failUnsupportedProposal(logger, ctx, k, proposal, err.Error(), true); err != nil {
					return false, err
				}

				if err = k.ActiveProposalsQueue.Remove(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id)); err != nil {
					return false, err
				}

				return false, nil
			}

			return false, err
		}

		certifierVoteIsRequired, err := k.CertifierVoteIsRequired(ctx, proposal.Id)
		if err != nil {
			return false, err
		}
		if certifierVoteIsRequired {
			certifierVoted, err := k.GetCertifierVoted(ctx, proposal.Id)
			if err != nil {
				return false, err
			}
			if !certifierVoted {
				var endVoting bool
				_, endVoting, tallyResults = keeper.SecurityTally(ctx, *k, proposal)
				if !endVoting {
					// Skip the rest of this iteration, because the proposal needs to go
					// through the validator voting period now.
					err = k.SetCertifierVoted(ctx, proposal.Id)
					if err != nil {
						return false, err
					}
					err = k.DeleteVotes(ctx, proposal.Id)
					if err != nil {
						return false, err
					}
				}
				return false, nil
			}
		}

		//var tagValue, logMsg string
		passes, burnDeposits, tallyResults, err = k.Tally(ctx, proposal)
		if err != nil {
			return false, err
		}

		// If an expedited proposal fails, we do not want to update
		// the deposit at this point since the proposal is converted to regular.
		// As a result, the deposits are either deleted or refunded in all cases
		// EXCEPT when an expedited proposal fails.
		if !(proposal.Expedited && !passes) {
			if burnDeposits {
				err = k.DeleteAndBurnDeposits(ctx, proposal.Id)
			} else {
				err = k.RefundAndDeleteDeposits(ctx, proposal.Id)
			}
			if err != nil {
				return false, err
			}
		}

		if err = k.ActiveProposalsQueue.Remove(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id)); err != nil {
			return false, err
		}

		switch {
		case passes:
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
			if err != nil {
				proposal.Status = govtypesv1.StatusFailed
				proposal.FailedReason = err.Error()
				tagValue = govtypes.AttributeValueProposalFailed
				logMsg = fmt.Sprintf("passed proposal (%v) failed to execute; msgs: %s", proposal, err)

				break
			}

			// execute all messages
			for idx, msg = range messages {
				handler := k.Router().Handler(msg)
				var res *sdk.Result
				res, err = safeExecuteHandler(cacheCtx, msg, handler)
				if err != nil {
					break
				}

				events = append(events, res.GetEvents()...)
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
				proposal.FailedReason = err.Error()
				tagValue = govtypes.AttributeValueProposalFailed
				logMsg = fmt.Sprintf("passed, but msg %d (%s) failed on execution: %s", idx, sdk.MsgTypeURL(msg), err)
			}
		case proposal.Expedited:
			// When expedited proposal fails, it is converted
			// to a regular proposal. As a result, the voting period is extended, and,
			// once the regular voting period expires again, the tally is repeated
			// according to the regular proposal rules.
			proposal.Expedited = false
			params, err := k.Params.Get(ctx)
			if err != nil {
				return false, err
			}
			endTime := proposal.VotingStartTime.Add(*params.VotingPeriod)
			proposal.VotingEndTime = &endTime

			err = k.ActiveProposalsQueue.Set(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id), proposal.Id)
			if err != nil {
				return false, err
			}

			tagValue = govtypes.AttributeValueExpeditedProposalRejected
			logMsg = "expedited proposal converted to regular"
		default:
			proposal.Status = govtypesv1.StatusRejected
			proposal.FailedReason = "proposal did not get enough votes to pass"
			tagValue = govtypes.AttributeValueProposalRejected
			logMsg = "rejected"
		}

		proposal.FinalTallyResult = &tallyResults

		err = k.SetProposal(ctx, proposal)
		if err != nil {
			return false, err
		}

		logger.Info(
			"proposal tallied",
			"proposal", proposal.Id,
			"status", proposal.Status.String(),
			"expedited", proposal.Expedited,
			"title", proposal.Title,
			"results", logMsg,
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govtypes.EventTypeActiveProposal,
				sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.Id)),
				sdk.NewAttribute(govtypes.AttributeKeyProposalResult, tagValue),
				sdk.NewAttribute(govtypes.AttributeKeyProposalLog, logMsg),
			),
		)

		return false, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func processSecurityVote(ctx sdk.Context, k *keeper.Keeper, logger log.Logger) error {
	var (
		tagValue, logMsg string
		passes           bool
		tallyResults     govtypesv1.TallyResult
	)
	// Iterate over all active proposals, regardless of end time, so that
	// security voting can end as soon as a passing threshold is met.
	rng := collections.NewPrefixUntilPairRange[time.Time, uint64](time.Unix(common.MaxTimestamp, 0))
	err := k.ActiveProposalsQueue.Walk(ctx, rng, func(key collections.Pair[time.Time, uint64], _ uint64) (bool, error) {
		proposal, err := k.Proposals.Get(ctx, key.K2())
		if err != nil {
			// if the proposal has an encoding error, this means it cannot be processed by x/gov
			// this could be due to some types missing their registration
			// instead of returning an error (i.e, halting the chain), we fail the proposal
			if errors.Is(err, collections.ErrEncoding) {
				proposal.Id = key.K2()
				if err := failUnsupportedProposal(logger, ctx, k, proposal, err.Error(), true); err != nil {
					return false, err
				}

				if err = k.ActiveProposalsQueue.Remove(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id)); err != nil {
					return false, err
				}

				return false, nil
			}

			return false, err
		}

		certifierVoteIsRequired, err := k.CertifierVoteIsRequired(ctx, proposal.Id)
		if err != nil {
			return false, err
		}
		// Only process security proposals
		if !certifierVoteIsRequired {
			return false, nil
		}
		// Only process proposals in the security voting period.
		certifierVoted, err := k.GetCertifierVoted(ctx, proposal.Id)
		if err != nil {
			return false, err
		}
		if certifierVoted {
			return false, nil
		}

		var endVoting bool
		passes, endVoting, tallyResults = keeper.SecurityTally(ctx, *k, proposal)
		if !passes {
			// Do nothing, because the proposal still has time before the voting period ends.
			return false, nil
		}
		//Else: the proposal passed the certifier voting period.
		if endVoting {
			var (
				//idx    int
				events sdk.Events
				msg    sdk.Msg
			)

			cacheCtx, writeCache := ctx.CacheContext()
			messages, err := proposal.GetMsgs()
			if err != nil {
				proposal.Status = govtypesv1.StatusFailed
				proposal.FailedReason = err.Error()
				tagValue = govtypes.AttributeValueProposalFailed
				logMsg = fmt.Sprintf("passed proposal (%v) failed to execute; msgs: %s", proposal, err)
			} else {
				for _, msg = range messages {
					handler := k.Router().Handler(msg)

					var res *sdk.Result
					res, err = handler(cacheCtx, msg)
					if err != nil {
						break
					}

					events = append(events, res.GetEvents()...)
				}

				proposal.Status = govtypesv1.StatusPassed
				tagValue = govtypes.AttributeValueProposalPassed
				logMsg = "passed"

				// write state to the underlying multi-store
				writeCache()

				// propagate the msg events to the current context
				ctx.EventManager().EmitEvents(events)
			}

			proposal.FinalTallyResult = &tallyResults
			err = k.SetProposal(ctx, proposal)
			if err != nil {
				return false, err
			}
			if err = k.ActiveProposalsQueue.Remove(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id)); err != nil {
				return false, err
			}

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
			err = k.SetCertifierVoted(ctx, proposal.Id)
			if err != nil {
				return false, err
			}
			err = k.DeleteVotes(ctx, proposal.Id)
			if err != nil {
				return false, err
			}
		}

		return true, nil
	})

	return err
}

// executes handle(msg) and recovers from panic.
func safeExecuteHandler(ctx sdk.Context, msg sdk.Msg, handler baseapp.MsgServiceHandler,
) (res *sdk.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("handling x/gov proposal msg [%s] PANICKED: %v", msg, r)
		}
	}()
	res, err = handler(ctx, msg)
	return
}

// failUnsupportedProposal fails a proposal that cannot be processed by gov
func failUnsupportedProposal(
	logger log.Logger,
	ctx sdk.Context,
	keeper *keeper.Keeper,
	proposal govtypesv1.Proposal,
	errMsg string,
	active bool,
) error {
	proposal.Status = govtypesv1.StatusFailed
	proposal.FailedReason = fmt.Sprintf("proposal failed because it cannot be processed by gov: %s", errMsg)
	proposal.Messages = nil // clear out the messages

	if err := keeper.SetProposal(ctx, proposal); err != nil {
		return err
	}

	if err := keeper.RefundAndDeleteDeposits(ctx, proposal.Id); err != nil {
		return err
	}

	eventType := govtypes.EventTypeInactiveProposal
	if active {
		eventType = govtypes.EventTypeActiveProposal
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.Id)),
			sdk.NewAttribute(govtypes.AttributeKeyProposalResult, govtypes.AttributeValueProposalFailed),
		),
	)

	logger.Info(
		"proposal failed to decode; deleted",
		"proposal", proposal.Id,
		"expedited", proposal.Expedited,
		"title", proposal.Title,
		"results", errMsg,
	)

	return nil
}
