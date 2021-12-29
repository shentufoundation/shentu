package keeper

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/certikfoundation/shentu/v2/x/gov/types"
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
)

// AddVote Adds a vote on a specific proposal.
<<<<<<< HEAD
func (k Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options govTypes.WeightedVoteOptions) error {
=======
func (k Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypes.WeightedVoteOptions) error {
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		return sdkerrors.Wrapf(govtypes.ErrUnknownProposal, "%v", proposalID)
	}
	if proposal.Status != types.StatusCertifierVotingPeriod &&
		proposal.Status != types.StatusValidatorVotingPeriod {
		return sdkerrors.Wrapf(govtypes.ErrInactiveProposal, "%v", proposalID)
	}

	for _, option := range options {
<<<<<<< HEAD
		if !govTypes.ValidWeightedVoteOption(option) {
			return sdkerrors.Wrapf(govTypes.ErrInvalidVote, "%s", option)
		}

		if proposal.Status == types.StatusCertifierVotingPeriod {
			if !(option.Option == govTypes.OptionYes ||
				option.Option == govTypes.OptionNo) {
				return sdkerrors.Wrapf(govTypes.ErrInvalidVote,
=======
		if !govtypes.ValidWeightedVoteOption(option) {
			return sdkerrors.Wrapf(govtypes.ErrInvalidVote, "%s", option)
		}

		if proposal.Status == types.StatusCertifierVotingPeriod {
			if !(option.Option == govtypes.OptionYes ||
				option.Option == govtypes.OptionNo) {
				return sdkerrors.Wrapf(govtypes.ErrInvalidVote,
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a
					"'%s' is not valid option in certifier voting; must be 'yes' or 'no'", option.Option)
			}
		}
	}

	if proposal.Status == types.StatusCertifierVotingPeriod && !k.IsCertifier(ctx, voterAddr) {
		return sdkerrors.Wrapf(govtypes.ErrInvalidVote, "'%s' is not a certifier.", voterAddr)
	}

	if proposal.GetContent().ProposalType() == shieldtypes.ProposalTypeShieldClaim &&
		proposal.Status == types.StatusValidatorVotingPeriod &&
		!k.IsCertifiedIdentity(ctx, voterAddr) {
		return sdkerrors.Wrapf(govtypes.ErrInvalidVote, "'%s' is not a certified identity", voterAddr)
	}

	txhash := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
<<<<<<< HEAD
	vote := types.NewVote(proposalID, voterAddr, options, txhash)
=======
	vote := govtypes.NewVote(proposalID, voterAddr, options)
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a
	k.setVote(ctx, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
<<<<<<< HEAD
			govTypes.EventTypeProposalVote,
			sdk.NewAttribute(govTypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(govTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
=======
			govtypes.EventTypeProposalVote,
			sdk.NewAttribute(govtypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a
			sdk.NewAttribute(types.AttributeKeyVoter, voterAddr.String()),
			sdk.NewAttribute(types.AttributeTxHash, txhash),
		),
	)

	return nil
}

// GetAllVotes returns all the votes from the store.
func (k Keeper) GetAllVotes(ctx sdk.Context) (votes govtypes.Votes) {
	k.IterateAllVotes(ctx, func(vote govtypes.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVotes returns all votes on a given proposal.
func (k Keeper) GetVotes(ctx sdk.Context, proposalID uint64) (votes govtypes.Votes) {
	k.IterateVotes(ctx, proposalID, func(vote govtypes.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVotesPaginated performs paginated query of votes on a given proposal.
func (k Keeper) GetVotesPaginated(ctx sdk.Context, proposalID uint64, page, limit uint) (votes govtypes.Votes) {
	k.IterateVotesPaginated(ctx, proposalID, page, limit, func(vote govtypes.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVote gets the vote from an address on a specific proposal.
func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (vote govtypes.Vote, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(govtypes.VoteKey(proposalID, voterAddr))
	if bz == nil {
		return vote, false
	}

	k.cdc.MustUnmarshal(bz, &vote)
	return vote, true
}

// setVote set a vote.
func (k Keeper) setVote(ctx sdk.Context, vote govtypes.Vote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&vote)
	addr, err := sdk.AccAddressFromBech32(vote.Voter)
	if err != nil {
		panic(err)
	}
	store.Set(govtypes.VoteKey(vote.ProposalId, addr), bz)
}

// SetVote set a vote.
func (k Keeper) SetVote(ctx sdk.Context, vote govtypes.Vote) {
	k.setVote(ctx, vote)
}

// GetVotesIterator returns an iterator to go over all votes on a given proposal.
func (k Keeper) GetVotesIterator(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, govtypes.VotesKey(proposalID))
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

// IterateAllVotes iterates over the all the stored votes and performs a callback function.
func (k Keeper) IterateAllVotes(ctx sdk.Context, cb func(vote govtypes.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, govtypes.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
<<<<<<< HEAD
		var vote types.Vote
=======
		var vote govtypes.Vote
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a
		k.cdc.MustUnmarshal(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// IterateVotes iterates over the all votes on a given proposal and performs a callback function.
func (k Keeper) IterateVotes(ctx sdk.Context, proposalID uint64, cb func(vote govtypes.Vote) (stop bool)) {
	iterator := k.GetVotesIterator(ctx, proposalID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
<<<<<<< HEAD
		var vote types.Vote
=======
		var vote govtypes.Vote
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a
		k.cdc.MustUnmarshal(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// IterateVotesPaginated iterates over votes on a given proposal
// based on pagination parameters and performs a callback function.
func (k Keeper) IterateVotesPaginated(ctx sdk.Context, proposalID uint64, page, limit uint, cb func(vote govtypes.Vote) (stop bool)) {
	iterator := k.GetVotesIteratorPaginated(ctx, proposalID, page, limit)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
<<<<<<< HEAD
		var vote types.Vote
=======
		var vote govtypes.Vote
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a
		k.cdc.MustUnmarshal(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}
