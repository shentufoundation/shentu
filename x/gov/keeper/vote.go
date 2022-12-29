package keeper

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// AddVote Adds a vote on a specific proposal.
func (k Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypes.WeightedVoteOptions) error {
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		return sdkerrors.Wrapf(govtypes.ErrUnknownProposal, "%d", proposalID)
	}

	if proposal.Status != govtypes.StatusVotingPeriod {
		return sdkerrors.Wrapf(govtypes.ErrInactiveProposal, "%d", proposalID)
	}

	for _, option := range options {
		if !govtypes.ValidWeightedVoteOption(option) {
			return sdkerrors.Wrapf(govtypes.ErrInvalidVote, "%s", option)
		}
	}

	// Add cert vote
	if k.HasSecurityVoting(proposal) && !k.IsCertifierVoted(ctx, proposalID) {
		return k.AddCertifierVoted(ctx, proposalID, voterAddr, options)
	}

	if proposal.GetContent().ProposalType() == shieldtypes.ProposalTypeShieldClaim &&
		proposal.Status == govtypes.StatusVotingPeriod &&
		!k.IsCertifiedIdentity(ctx, voterAddr) {
		return sdkerrors.Wrapf(govtypes.ErrInvalidVote, "'%s' is not a certified identity", voterAddr)
	}

	txhash := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
	vote := govtypes.NewVote(proposalID, voterAddr, options)
	k.SetVote(ctx, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeProposalVote,
			sdk.NewAttribute(govtypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyVoter, voterAddr.String()),
			sdk.NewAttribute(types.AttributeTxHash, txhash),
		),
	)

	return nil
}

// GetVotesIteratorPaginated returns an iterator to go over
// votes on a given proposal based on pagination parameters.
func (k Keeper) GetVotesIteratorPaginated(ctx sdk.Context, proposalID uint64, page, limit uint) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIteratorPaginated(store, govtypes.VotesKey(proposalID), page, limit)
}

// deleteVote delete a vote for a proposal.
func (k Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(govtypes.VoteKey(proposalID, voterAddr))
}

// DeleteAllVotes deletes all votes for a proposal.
func (k Keeper) DeleteAllVotes(ctx sdk.Context, proposalID uint64) {
	k.IterateVotes(ctx, proposalID, func(vote govtypes.Vote) bool {
		addr, err := sdk.AccAddressFromBech32(vote.Voter)
		if err != nil {
			panic(err)
		}
		k.deleteVote(ctx, proposalID, addr)
		return false
	})
}

// AddCertifierVoted add a certifier vote
// The only voting options are "yes" and "no".
func (k Keeper) AddCertifierVoted(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypes.WeightedVoteOptions) error {
	if !k.IsCertifier(ctx, voterAddr) {
		return sdkerrors.Wrapf(govtypes.ErrInvalidVote, "%s is not a certified identity", voterAddr)
	}

	for _, option := range options {
		if !(option.Option == govtypes.OptionYes ||
			option.Option == govtypes.OptionNo) {
			return sdkerrors.Wrapf(govtypes.ErrInvalidVote,
				"'%s' is not valid option in certifier voting; must be 'yes' or 'no'", option.Option)
		}
	}

	txhash := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
	vote := govtypes.NewVote(proposalID, voterAddr, options)
	k.SetVote(ctx, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetCertVote,
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeTxHash, txhash),
		),
	)
	return nil
}

func (k Keeper) SetCertifierVoted(ctx sdk.Context, proposalID uint64) {
	k.SetCertVote(ctx, proposalID)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetCertVote,
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)
}

// SetCertVote sets a cert vote to the gov store
func (k Keeper) SetCertVote(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CertVotesKey(proposalID), govtypes.GetProposalIDBytes(proposalID))
}

// IsCertifierVoted determine cert vote for custom proposal types have finished
func (k Keeper) IsCertifierVoted(ctx sdk.Context, proposalID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CertVotesKey(proposalID))
	return bz != nil
}
