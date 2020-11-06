package gov

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/gov/internal/keeper"
	"github.com/certikfoundation/shentu/x/gov/internal/types"
	"github.com/certikfoundation/shentu/x/shield"
)

func removeInactiveProposals(ctx sdk.Context, k keeper.Keeper) {
	k.IterateInactiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal types.Proposal) bool {
		k.DeleteProposalByProposalID(ctx, proposal.ProposalID)
		k.RefundDepositsByProposalID(ctx, proposal.ProposalID)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govTypes.EventTypeInactiveProposal,
				sdk.NewAttribute(govTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
				sdk.NewAttribute(govTypes.AttributeKeyProposalResult, govTypes.AttributeValueProposalDropped),
			),
		)

		// TODO log reason of proposal deletion
		return false
	})
}

func updateVeto(ctx sdk.Context, k keeper.Keeper, proposal types.Proposal) {
	if proposal.ProposalType() == shield.ProposalTypeShieldClaim {
		c := proposal.Content.(shield.ClaimProposal)
		k.ShieldKeeper.ClaimEnd(ctx, c.ProposalID, c.PoolID, c.Loss)
	}
}

func updateAbstain(ctx sdk.Context, k keeper.Keeper, proposal types.Proposal) {
	if proposal.ProposalType() == shield.ProposalTypeShieldClaim {
		c := proposal.Content.(shield.ClaimProposal)
		k.ShieldKeeper.RestoreShield(ctx, c.PoolID, proposal.ProposerAddress, c.PurchaseID, c.Loss)
		k.ShieldKeeper.ClaimEnd(ctx, c.ProposalID, c.PoolID, c.Loss)
	}
}

func processActiveProposal(ctx sdk.Context, k keeper.Keeper, proposal types.Proposal) bool {
	var (
		tagValue     string
		pass, veto   bool
		tallyResults govTypes.TallyResult
	)

	if proposal.Status == types.StatusCertifierVotingPeriod {
		var endVoting bool
		pass, endVoting, tallyResults = keeper.SecurityTally(ctx, k, proposal)
		if !endVoting {
			// Skip the rest of this iteration, because the proposal needs to go
			// through the validator voting period now.
			k.DeleteAllVotes(ctx, proposal.ProposalID)
			k.ActivateVotingPeriod(ctx, proposal)
			return false
		}
	} else {
		pass, veto, tallyResults = keeper.Tally(ctx, k, proposal)
	}

	if veto {
		k.DeleteDepositsByProposalID(ctx, proposal.ProposalID)
		updateVeto(ctx, k, proposal)
	} else {
		k.RefundDepositsByProposalID(ctx, proposal.ProposalID)
		if !pass {
			updateAbstain(ctx, k, proposal)
		}
	}

	if pass {
		handler := k.Router().GetRoute(proposal.ProposalRoute())
		cacheCtx, writeCache := ctx.CacheContext()

		// The proposal handler may execute state mutating logic depending on the
		// proposal content. If the handler fails, no state mutation is written and
		// the error message is logged.
		err := handler(cacheCtx, proposal.Content)
		if err == nil {
			proposal.Status = types.StatusPassed
			tagValue = govTypes.AttributeValueProposalPassed

			// write state to the underlying multi-store
			writeCache()
		} else {
			proposal.Status = types.StatusFailed
			tagValue = govTypes.AttributeValueProposalFailed
		}
	} else {
		proposal.Status = types.StatusRejected
		tagValue = govTypes.AttributeValueProposalRejected
	}

	proposal.FinalTallyResult = tallyResults

	k.SetProposal(ctx, proposal)
	k.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)

	// TODO log tallying result

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govTypes.EventTypeActiveProposal,
			sdk.NewAttribute(govTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
			sdk.NewAttribute(govTypes.AttributeKeyProposalResult, tagValue),
		),
	)
	return false
}

func processSecurityVote(ctx sdk.Context, k keeper.Keeper, proposal types.Proposal) bool {
	var (
		tagValue     string
		pass         bool
		tallyResults govTypes.TallyResult
	)

	// Only process proposals in the security voting period.
	if proposal.Status != types.StatusCertifierVotingPeriod {
		return false
	}

	var endVoting bool
	pass, endVoting, tallyResults = keeper.SecurityTally(ctx, k, proposal)
	if !pass {
		// Do nothing, because the proposal still has time before the voting period
		// ends.
		return false
	}
	// Else: the proposal passed the certifier voting period.

	if endVoting {
		handler := k.Router().GetRoute(proposal.ProposalRoute())
		cacheCtx, writeCache := ctx.CacheContext()

		// The proposal handler may execute state mutating logic depending on the
		// proposal content. If the handler fails, no state mutation is written and
		// the error message is logged.
		err := handler(cacheCtx, proposal.Content)
		if err == nil {
			proposal.Status = types.StatusPassed
			tagValue = govTypes.AttributeValueProposalPassed

			// write state to the underlying multi-store
			writeCache()
		} else {
			proposal.Status = types.StatusFailed
			tagValue = govTypes.AttributeValueProposalFailed
		}

		proposal.FinalTallyResult = tallyResults

		k.SetProposal(ctx, proposal)
		k.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)

		// TODO log tallying result

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govTypes.EventTypeActiveProposal,
				sdk.NewAttribute(govTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
				sdk.NewAttribute(govTypes.AttributeKeyProposalResult, tagValue),
			),
		)
	} else {
		// Activate validator voting period
		k.DeleteAllVotes(ctx, proposal.ProposalID)
		k.ActivateVotingPeriod(ctx, proposal)
	}
	return false
}

// EndBlocker is called every block, removes inactive proposals, tallies active
// proposals and deletes/refunds deposits.
// TODO refactor into smaller functions
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// delete inactive proposal from store and its deposits
	removeInactiveProposals(ctx, k)

	// fetch active proposals whose voting periods have ended (are passed the
	// block time)
	k.IterateActiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal types.Proposal) bool {
		return processActiveProposal(ctx, k, proposal)
	})

	// Iterate over all active proposals, regardless of end time, so that
	// security voting can end as soon as a passing threshold is met.
	k.IterateActiveProposalsQueue(ctx, time.Unix(cert.MaxTimestamp, 0), func(proposal types.Proposal) bool {
		return processSecurityVote(ctx, k, proposal)
	})
}
