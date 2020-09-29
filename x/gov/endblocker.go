package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

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
		_ = k.ShieldKeeper.ClaimUnlock(ctx, c.PoolID, c.Loss, proposal.ProposalID)
	}
}

func updateAbstain(ctx sdk.Context, k keeper.Keeper, proposal types.Proposal) {
	if proposal.ProposalType() == shield.ProposalTypeShieldClaim {
		c := proposal.Content.(shield.ClaimProposal)
		_ = k.ShieldKeeper.ClaimUnlock(ctx, c.PoolID, c.Loss, proposal.ProposalID)
		_ = k.ShieldKeeper.RestoreShield(ctx, c.PoolID, c.Loss, c.PurchaseTxHash)
	}
}

// EndBlocker is called every block, removes inactive proposals, tallies active proposals and deletes/refunds deposits.
// TODO refactor into smaller functions
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// delete inactive proposal from store and its deposits
	removeInactiveProposals(ctx, k)

	// fetch active proposals whose voting periods have ended (are passed the block time)
	k.IterateActiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal types.Proposal) bool {
		var (
			tagValue     string
			pass, veto   bool
			tallyResults govTypes.TallyResult
		)

		if proposal.Status == types.StatusCertifierVotingPeriod {
			var endVoting bool
			pass, endVoting, tallyResults = keeper.SecurityTally(ctx, k, proposal)
			if !endVoting {
				// Skip the rest of this iteration, because the proposal needs
				// to go through the validator voting period now.
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

			// The proposal handler may execute state mutating logic depending
			// on the proposal content. If the handler fails, no state mutation
			// is written and the error message is logged.
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
	})
}
