package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/exported"

	"github.com/certikfoundation/shentu/x/gov/internal/types"
	"github.com/certikfoundation/shentu/x/shield"
)

// Proposal

// GetProposal get Proposal from store by ProposalID.
func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (proposal types.Proposal, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(ProposalKey(proposalID))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposal)
	return proposal, true
}

// SetProposal sets a proposal to store.
func (k Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(proposal)
	store.Set(ProposalKey(proposal.ProposalID), bz)
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
	return append(govTypes.ProposalsKeyPrefix, bz...)
}

// isValidator checks if the input address is a validator.
func (k Keeper) isValidator(ctx sdk.Context, addr sdk.AccAddress) bool {
	isValidator := false
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator staking.ValidatorI) (stop bool) {
		if validator.GetOperator().Equals(addr) {
			isValidator = true
			return true
		}
		return false
	})
	return isValidator
}

// isCertifier checks if the input address is a certifier.
func (k Keeper) IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.CertKeeper.IsCertifier(ctx, addr)
}

// IsCouncilMember checks if the address is either a validator or a certifier.
func (k Keeper) IsCouncilMember(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.isValidator(ctx, addr) || k.IsCertifier(ctx, addr)
}

// SubmitProposal creates a new proposal with given content.
func (k Keeper) SubmitProposal(ctx sdk.Context, content govTypes.Content, addr sdk.AccAddress) (types.Proposal, error) {
	proposalID, err := k.GetProposalID(ctx)
	if err != nil {
		return types.Proposal{}, err
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := k.GetDepositParams(ctx).MaxDepositPeriod

	var proposal types.Proposal
	if content.ProposalType() == shield.ProposalTypeShieldClaim {
		c := content.(shield.ClaimProposal)
		c.ProposalID = proposalID
		proposal = types.NewProposal(c, proposalID, addr, k.IsCouncilMember(ctx, addr), submitTime, submitTime.Add(depositPeriod))
	} else {
		proposal = types.NewProposal(content, proposalID, addr, k.IsCouncilMember(ctx, addr), submitTime, submitTime.Add(depositPeriod))
	}

	k.SetProposal(ctx, proposal)
	k.InsertInactiveProposalQueue(ctx, proposalID, proposal.DepositEndTime)
	k.SetProposalID(ctx, proposalID+1)

	return proposal, nil
}

// IterateProposals iterates over the all the proposals and performs a callback function.
func (k Keeper) IterateProposals(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, govTypes.ProposalsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proposal types.Proposal
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &proposal)

		if cb(proposal) {
			break
		}
	}
}

// GetProposals returns all the proposals from store.
func (k Keeper) GetProposals(ctx sdk.Context) (proposals types.Proposals) {
	k.IterateProposals(ctx, func(proposal types.Proposal) bool {
		proposals = append(proposals, proposal)
		return false
	})
	return
}

// ActivateVotingPeriod switches proposals from deposit period to voting period.
func (k Keeper) ActivateVotingPeriod(ctx sdk.Context, proposal types.Proposal) {
	proposal.VotingStartTime = ctx.BlockHeader().Time
	votingPeriod := k.GetVotingParams(ctx).VotingPeriod
	oldVotingEndTime := proposal.VotingEndTime
	proposal.VotingEndTime = proposal.VotingStartTime.Add(votingPeriod)
	oldDepositEndTime := proposal.DepositEndTime

	if proposal.HasSecurityVoting() && (proposal.Status != types.StatusCertifierVotingPeriod) {
		// Special case: just for software upgrade, certifier update and shield claim proposals.
		proposal.Status = types.StatusCertifierVotingPeriod
	} else {
		// Default case: for plain text proposals, community pool spend proposals;
		// and second round of software upgrade, certifier update and shield claim
		// proposals.
		if proposal.Status == types.StatusCertifierVotingPeriod {
			k.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, oldVotingEndTime)
		} else {
			proposal.DepositEndTime = ctx.BlockHeader().Time
		}
		proposal.Status = types.StatusValidatorVotingPeriod
	}

	k.SetProposal(ctx, proposal)
	k.RemoveFromInactiveProposalQueue(ctx, proposal.ProposalID, oldDepositEndTime)
	k.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
}

// ActivateCouncilProposalVotingPeriod only switches proposals of council members.
func (k Keeper) ActivateCouncilProposalVotingPeriod(ctx sdk.Context, proposal types.Proposal) bool {
	if proposal.IsProposerCouncilMember {
		k.ActivateVotingPeriod(ctx, proposal)
		return true
	}
	return false
}

// GetProposalsFiltered returns proposals filtered.
func (k Keeper) GetProposalsFiltered(ctx sdk.Context, params types.QueryProposalsParams) []types.Proposal {
	proposals := k.GetProposals(ctx)
	filteredProposals := make([]types.Proposal, 0, len(proposals))

	for _, p := range proposals {
		matchVoter, matchDepositor, matchStatus := true, true, true

		// match status (if supplied/valid)
		if types.ValidProposalStatus(params.ProposalStatus) {
			matchStatus = p.Status == params.ProposalStatus
		}

		// match voter address (if supplied)
		if len(params.Voter) > 0 {
			_, matchVoter = k.GetVote(ctx, p.ProposalID, params.Voter)
		}

		// match depositor (if supplied)
		if len(params.Depositor) > 0 {
			_, matchDepositor = k.GetDeposit(ctx, p.ProposalID, params.Depositor)
		}

		if matchVoter && matchDepositor && matchStatus {
			filteredProposals = append(filteredProposals, p)
		}
	}

	start, end := client.Paginate(len(filteredProposals), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredProposals = []types.Proposal{}
	} else {
		filteredProposals = filteredProposals[start:end]
	}

	return filteredProposals
}
