package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/v2/x/gov/types"
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
)

// Proposal

// GetProposal get Proposal from store by ProposalID.
func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (types.Proposal, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(ProposalKey(proposalID))
	if bz == nil {
		return types.Proposal{}, false
	}

	var proposal types.Proposal
	k.MustUnmarshalProposal(bz, &proposal)

	return proposal, true
}

// SetProposal sets a proposal to store.
func (k Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)
	bz := k.MustMarshalProposal(proposal)
	store.Set(ProposalKey(proposal.ProposalId), bz)
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

// isCertifier checks if the input address is a certifier.
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
func (k Keeper) SubmitProposal(ctx sdk.Context, content govtypes.Content, addr sdk.AccAddress) (types.Proposal, error) {
	if !k.router.HasRoute(content.ProposalRoute()) {
		return types.Proposal{}, sdkerrors.Wrap(govtypes.ErrNoProposalHandlerExists, content.ProposalRoute())
	}

	proposalID, err := k.GetProposalID(ctx)
	if err != nil {
		return types.Proposal{}, err
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
		return types.Proposal{}, sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, err.Error())
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := k.GetDepositParams(ctx).MaxDepositPeriod

	var proposal types.Proposal
	proposal, err = types.NewProposal(content, proposalID, addr, k.IsCouncilMember(ctx, addr), submitTime, submitTime.Add(depositPeriod))
	if err != nil {
		return types.Proposal{}, err
	}
	k.SetProposal(ctx, proposal)
	k.InsertInactiveProposalQueue(ctx, proposalID, proposal.DepositEndTime)
	k.SetProposalID(ctx, proposalID+1)
	return proposal, nil
}

// IterateProposals iterates over the all the proposals and performs a callback function.
func (k Keeper) IterateProposals(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, govtypes.ProposalsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proposal types.Proposal
		k.MustUnmarshalProposal(iterator.Value(), &proposal)

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
			k.RemoveFromActiveProposalQueue(ctx, proposal.ProposalId, oldVotingEndTime)
		} else {
			proposal.DepositEndTime = ctx.BlockHeader().Time
		}
		proposal.Status = types.StatusValidatorVotingPeriod
	}

	k.SetProposal(ctx, proposal)
	k.RemoveFromInactiveProposalQueue(ctx, proposal.ProposalId, oldDepositEndTime)
	k.InsertActiveProposalQueue(ctx, proposal.ProposalId, proposal.VotingEndTime)
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
		filteredProposals = []types.Proposal{}
	} else {
		filteredProposals = filteredProposals[start:end]
	}

	return filteredProposals
}

func (keeper Keeper) MarshalProposal(proposal types.Proposal) ([]byte, error) {
	bz, err := keeper.cdc.Marshal(&proposal)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func (keeper Keeper) UnmarshalProposal(bz []byte, proposal *types.Proposal) error {
	err := keeper.cdc.Unmarshal(bz, proposal)
	if err != nil {
		return err
	}
	return nil
}

func (keeper Keeper) MustMarshalProposal(proposal types.Proposal) []byte {
	bz, err := keeper.MarshalProposal(proposal)
	if err != nil {
		panic(err)
	}
	return bz
}

func (keeper Keeper) MustUnmarshalProposal(bz []byte, proposal *types.Proposal) {
	err := keeper.UnmarshalProposal(bz, proposal)
	if err != nil {
		panic(err)
	}
}
