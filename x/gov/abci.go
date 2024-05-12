package gov

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v046 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v046"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

func removeInactiveProposals(ctx sdk.Context, k keeper.Keeper) {
	logger := k.Logger(ctx)

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
			"min_deposit", sdk.NewCoins(k.GetDepositParams(ctx).MinDeposit...).String(),
			"total_deposit", sdk.NewCoins(proposal.TotalDeposit...).String(),
		)

		return false
	})
}

func updateVeto(ctx sdk.Context, k keeper.Keeper, proposal govtypesv1.Proposal) {
	legacyProposal, err := v046.ConvertToLegacyProposal(proposal)
	if err != nil {
		return
	}
	if legacyProposal.ProposalType() == shieldtypes.ProposalTypeShieldClaim {
		c := legacyProposal.GetContent().(*shieldtypes.ShieldClaimProposal)
		k.ShieldKeeper.ClaimEnd(ctx, c.ProposalId, c.PoolId, c.Loss)
	}
}

func updateAbstain(ctx sdk.Context, k keeper.Keeper, proposal govtypesv1.Proposal) {
	legacyProposal, err := v046.ConvertToLegacyProposal(proposal)
	if err != nil {
		return
	}
	if legacyProposal.ProposalType() == shieldtypes.ProposalTypeShieldClaim {
		c := legacyProposal.GetContent().(*shieldtypes.ShieldClaimProposal)
		proposer, err := sdk.AccAddressFromBech32(c.Proposer)
		if err != nil {
			panic(err)
		}
		k.ShieldKeeper.RestoreShield(ctx, c.PoolId, proposer, c.PurchaseId, c.Loss)
		k.ShieldKeeper.ClaimEnd(ctx, c.ProposalId, c.PoolId, c.Loss)
	}
}

func processActiveProposal(ctx sdk.Context, k keeper.Keeper, proposal govtypesv1.Proposal) bool {
	var (
		tagValue, logMsg string
		pass, veto       bool
		tallyResults     govtypesv1.TallyResult
	)
	logger := k.Logger(ctx)

	if k.HasSecurityVoting(proposal) && !k.GetCertifierVoted(ctx, proposal.Id) {
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
		pass, veto, tallyResults = keeper.Tally(ctx, k, proposal)
	}

	if veto {
		k.DeleteAndBurnDeposits(ctx, proposal.Id)
		updateVeto(ctx, k, proposal)
	} else {
		k.RefundAndDeleteDeposits(ctx, proposal.Id)
		if !pass {
			updateAbstain(ctx, k, proposal)
		}
	}

	if pass {
		//handler := k.Router().GetRoute(proposal.ProposalRoute())
		//cacheCtx, writeCache := ctx.CacheContext()
		//// The proposal handler may execute state mutating logic depending on the
		//// proposal content. If the handler fails, no state mutation is written and
		//// the error message is logged.
		//err := handler(cacheCtx, proposal.GetContent())
		//if err == nil {
		//	proposal.Status = govtypesv1.StatusPassed
		//	tagValue = govtypes.AttributeValueProposalPassed
		//	logMsg = "passed"
		//
		//	// write state to the underlying multi-store
		//	writeCache()
		//} else {
		//	proposal.Status = govtypesv1.StatusFailed
		//	tagValue = govtypes.AttributeValueProposalFailed
		//	logMsg = fmt.Sprintf("passed, but failed on execution: %s", err)
		//}
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
	if !k.HasSecurityVoting(proposal) {
		return false
	}
	// Only process proposals in the security voting period.
	if k.HasSecurityVoting(proposal) && k.GetCertifierVoted(ctx, proposal.Id) {
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
		//handler := k.Router().GetRoute(proposal.ProposalRoute())
		//cacheCtx, writeCache := ctx.CacheContext()
		//
		//// The proposal handler may execute state mutating logic depending on the
		//// proposal content. If the handler fails, no state mutation is written and
		//// the error message is logged.
		//err := handler(cacheCtx, proposal.GetContent())
		//if err == nil {
		//	proposal.Status = govtypesv1.StatusPassed
		//	tagValue = govtypes.AttributeValueProposalPassed
		//	logMsg = "passed"
		//
		//	// write state to the underlying multi-store
		//	writeCache()
		//} else {
		//proposal.Status = govtypesv1.StatusFailed
		//tagValue = govtypes.AttributeValueProposalFailed
		//logMsg = fmt.Sprintf("passed, but failed on execution: %s", err)
		//}

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
