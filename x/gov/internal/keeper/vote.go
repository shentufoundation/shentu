package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/gov/internal/types"
	"github.com/certikfoundation/shentu/x/shield"
)

// AddVote Adds a vote on a specific proposal.
func (k Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, option govTypes.VoteOption) error {
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		return sdkerrors.Wrapf(govTypes.ErrUnknownProposal, "%v", proposalID)
	}
	if proposal.Status != types.StatusCertifierVotingPeriod &&
		proposal.Status != types.StatusValidatorVotingPeriod {
		return sdkerrors.Wrapf(govTypes.ErrInactiveProposal, "%v", proposalID)
	}

	if !govTypes.ValidVoteOption(option) {
		return sdkerrors.Wrapf(govTypes.ErrInvalidVote, "%s", option)
	}

	if proposal.Status == types.StatusCertifierVotingPeriod {
		if !(option == govTypes.OptionYes ||
			option == govTypes.OptionNo) {
			return sdkerrors.Wrapf(govTypes.ErrInvalidVote,
				"'%s' is not valid option in certifier voting; must be 'yes' or 'no'", option)
		}
	}

	if proposal.Status == types.StatusCertifierVotingPeriod && !k.IsCertifier(ctx, voterAddr) {
		return sdkerrors.Wrapf(govTypes.ErrInvalidVote, "'%s' is not a certifier.", voterAddr)
	}

	if proposal.Content.ProposalType() == shield.ProposalTypeShieldClaim &&
		proposal.Status == types.StatusValidatorVotingPeriod &&
		!k.IsCertifiedIdentity(ctx, voterAddr) {
		return sdkerrors.Wrapf(govTypes.ErrInvalidVote, "'%s' is not a certified identity", voterAddr)
	}

	txhash := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
	vote := types.NewVote(proposalID, voterAddr, option, txhash)
	k.setVote(ctx, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govTypes.EventTypeProposalVote,
			sdk.NewAttribute(govTypes.AttributeKeyOption, option.String()),
			sdk.NewAttribute(govTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyVoter, voterAddr.String()),
			sdk.NewAttribute(types.AttributeTxHash, txhash),
		),
	)

	return nil
}

// GetAllVotes returns all the votes from the store.
func (k Keeper) GetAllVotes(ctx sdk.Context) (votes types.Votes) {
	k.IterateAllVotes(ctx, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVotes returns all votes on a given proposal.
func (k Keeper) GetVotes(ctx sdk.Context, proposalID uint64) (votes types.Votes) {
	k.IterateVotes(ctx, proposalID, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVotesPaginated performs paginated query of votes on a given proposal.
func (k Keeper) GetVotesPaginated(ctx sdk.Context, proposalID uint64, page, limit uint) (votes types.Votes) {
	k.IterateVotesPaginated(ctx, proposalID, page, limit, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVote gets the vote from an address on a specific proposal.
func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (vote types.Vote, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(govTypes.VoteKey(proposalID, voterAddr))
	if bz == nil {
		return vote, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &vote)
	return vote, true
}

// setVote set a vote.
func (k Keeper) setVote(ctx sdk.Context, vote types.Vote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(govTypes.VoteKey(vote.ProposalID, vote.Voter), bz)
}

// SetVote set a vote.
func (k Keeper) SetVote(ctx sdk.Context, vote types.Vote) {
	k.setVote(ctx, vote)
}

// GetVotesIterator returns an iterator to go over all votes on a given proposal.
func (k Keeper) GetVotesIterator(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, govTypes.VotesKey(proposalID))
}

// GetVotesIteratorPaginated returns an iterator to go over
// votes on a given proposal based on pagination parameters.
func (k Keeper) GetVotesIteratorPaginated(ctx sdk.Context, proposalID uint64, page, limit uint) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIteratorPaginated(store, govTypes.VotesKey(proposalID), page, limit)
}

// deleteVote delete a vote for a proposal.
func (k Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(govTypes.VoteKey(proposalID, voterAddr))
}

// DeleteAllVotes deletes all votes for a proposal.
func (k Keeper) DeleteAllVotes(ctx sdk.Context, proposalID uint64) {
	k.IterateVotes(ctx, proposalID, func(vote types.Vote) bool {
		k.deleteVote(ctx, proposalID, vote.Voter)
		return false
	})
}

// IterateAllVotes iterates over the all the stored votes and performs a callback function.
func (k Keeper) IterateAllVotes(ctx sdk.Context, cb func(vote types.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, govTypes.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// IterateVotes iterates over the all votes on a given proposal and performs a callback function.
func (k Keeper) IterateVotes(ctx sdk.Context, proposalID uint64, cb func(vote types.Vote) (stop bool)) {
	iterator := k.GetVotesIterator(ctx, proposalID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// IterateVotesPaginated iterates over votes on a given proposal
// based on pagination parameters and performs a callback function.
func (k Keeper) IterateVotesPaginated(ctx sdk.Context, proposalID uint64, page, limit uint, cb func(vote types.Vote) (stop bool)) {
	iterator := k.GetVotesIteratorPaginated(ctx, proposalID, page, limit)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}
