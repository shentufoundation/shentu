package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

func (k Keeper) ActivateVotingPeriod(ctx sdk.Context, proposal govtypes.Proposal) {
	proposal.VotingStartTime = ctx.BlockHeader().Time
	votingPeriod := k.GetVotingParams(ctx).VotingPeriod
	oldVotingEndTime := proposal.VotingEndTime
	proposal.VotingEndTime = proposal.VotingStartTime.Add(votingPeriod)
	oldDepositEndTime := proposal.DepositEndTime

	// Default case: for plain text proposals, community pool spend proposals;
	// and second round of software upgrade, certifier update and shield claim
	// proposals.
	if k.IsCertifierVoted(ctx, proposal.ProposalId) {
		k.RemoveFromActiveProposalQueue(ctx, proposal.ProposalId, oldVotingEndTime)
	} else {
		proposal.DepositEndTime = ctx.BlockHeader().Time
	}
	proposal.Status = govtypes.StatusVotingPeriod

	k.SetProposal(ctx, proposal)
	k.RemoveFromInactiveProposalQueue(ctx, proposal.ProposalId, oldDepositEndTime)
	k.InsertActiveProposalQueue(ctx, proposal.ProposalId, proposal.VotingEndTime)

}

// DeleteProposalByProposalID deletes a proposal from store.
func (k Keeper) DeleteProposalByProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		panic(fmt.Sprintf("couldn't find proposal with id#%d", proposalID))
	}
	k.RemoveFromInactiveProposalQueue(ctx, proposalID, proposal.DepositEndTime)
	k.RemoveFromActiveProposalQueue(ctx, proposalID, proposal.VotingEndTime)
	store.Delete(ProposalKey(proposalID))
}

// ProposalKey gets a specific proposal from the store.
func ProposalKey(proposalID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, proposalID)
	return append(govtypes.ProposalsKeyPrefix, bz...)
}

// isValidator checks if the input address is a validator.
func (k Keeper) isValidator(ctx sdk.Context, addr sdk.AccAddress) bool {
	isValidator := false
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		if validator.GetOperator().Equals(addr) {
			isValidator = true
			return true
		}
		return false
	})
	return isValidator
}

// IsCertifier checks if the input address is a certifier.
func (k Keeper) IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.CertKeeper.IsCertifier(ctx, addr)
}

// IsCouncilMember checks if the address is either a validator or a certifier.
func (k Keeper) IsCouncilMember(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.isValidator(ctx, addr) || k.IsCertifier(ctx, addr)
}

// IsCertifiedIdentity checks if the input address is a certified identity.
func (k Keeper) IsCertifiedIdentity(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.CertKeeper.IsCertified(ctx, addr.String(), "identity")
}

// TotalBondedByCertifiedIdentities calculates the amount of total bonded stakes by certified identities.
func (k Keeper) TotalBondedByCertifiedIdentities(ctx sdk.Context) sdk.Int {
	bonded := sdk.ZeroInt()
	for _, identity := range k.CertKeeper.GetCertifiedIdentities(ctx) {
		k.stakingKeeper.IterateDelegations(ctx, identity, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			val, found := k.stakingKeeper.GetValidator(ctx, delegation.GetValidatorAddr())
			if !found {
				return false
			}
			bonded = bonded.Add(delegation.GetShares().Quo(val.GetDelegatorShares()).MulInt(val.GetBondedTokens()).TruncateInt())
			return false
		})
	}
	return bonded
}

// SubmitProposal creates a new proposal with given content.
func (k Keeper) SubmitProposal(ctx sdk.Context, content govtypes.Content) (govtypes.Proposal, error) {
	if !k.router.HasRoute(content.ProposalRoute()) {
		return govtypes.Proposal{}, sdkerrors.Wrap(govtypes.ErrNoProposalHandlerExists, content.ProposalRoute())
	}

	proposalID, err := k.GetProposalID(ctx)
	if err != nil {
		return govtypes.Proposal{}, err
	}

	if c, ok := content.(*shieldtypes.ShieldClaimProposal); ok {
		c.ProposalId = proposalID
	}

	// Execute the proposal content in a cache-wrapped context to validate the
	// actual parameter changes before the proposal proceeds through the
	// governance process. State is not persisted.
	cacheCtx, _ := ctx.CacheContext()
	handler := k.router.GetRoute(content.ProposalRoute())
	if err := handler(cacheCtx, content); err != nil {
		return govtypes.Proposal{}, sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, err.Error())
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := k.GetDepositParams(ctx).MaxDepositPeriod

	proposal, err := govtypes.NewProposal(content, proposalID, submitTime, submitTime.Add(depositPeriod))
	if err != nil {
		return govtypes.Proposal{}, err
	}

	k.SetProposal(ctx, proposal)
	k.InsertInactiveProposalQueue(ctx, proposalID, proposal.DepositEndTime)
	k.SetProposalID(ctx, proposalID+1)

	// called right after a proposal is submitted
	k.AfterProposalSubmission(ctx, proposalID)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeSubmitProposal,
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return proposal, nil
}

// GetProposalsFiltered returns proposals filtered.
func (k Keeper) GetProposalsFiltered(ctx sdk.Context, params govtypes.QueryProposalsParams) []govtypes.Proposal {
	proposals := k.GetProposals(ctx)
	filteredProposals := make([]govtypes.Proposal, 0, len(proposals))

	for _, p := range proposals {
		matchVoter, matchDepositor, matchStatus := true, true, true

		// match status (if supplied/valid)
		if govtypes.ValidProposalStatus(params.ProposalStatus) {
			matchStatus = p.Status == params.ProposalStatus
		}

		// match voter address (if supplied)
		if len(params.Voter) > 0 {
			_, matchVoter = k.GetVote(ctx, p.ProposalId, params.Voter)
		}

		// match depositor (if supplied)
		if len(params.Depositor) > 0 {
			_, matchDepositor = k.GetDeposit(ctx, p.ProposalId, params.Depositor)
		}

		if matchVoter && matchDepositor && matchStatus {
			filteredProposals = append(filteredProposals, p)
		}
	}

	start, end := client.Paginate(len(filteredProposals), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredProposals = []govtypes.Proposal{}
	} else {
		filteredProposals = filteredProposals[start:end]
	}

	return filteredProposals
}

func (k Keeper) HasSecurityVoting(p govtypes.Proposal) bool {
	switch p.GetContent().(type) {
	case *upgradetypes.SoftwareUpgradeProposal, *certtypes.CertifierUpdateProposal, *shieldtypes.ShieldClaimProposal:
		return true
	default:
		return false
	}
}
