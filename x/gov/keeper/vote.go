package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
	"github.com/shentufoundation/shentu/v2/x/gov/types/v1"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// AddVote Adds a vote on a specific proposal.
func (k Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypesv1.WeightedVoteOptions, metadata string) error {
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		return sdkerrors.Wrapf(govtypes.ErrUnknownProposal, "%d", proposalID)
	}
	if proposal.Status != govtypesv1.StatusVotingPeriod {
		return sdkerrors.Wrapf(govtypes.ErrInactiveProposal, "%d", proposalID)
	}
	err := k.assertMetadataLength(metadata)
	if err != nil {
		return err
	}

	for _, option := range options {
		if !govtypesv1.ValidWeightedVoteOption(*option) {
			return sdkerrors.Wrapf(govtypes.ErrInvalidVote, "%s", option)
		}
	}

	// Add certifier vote
	if k.CertifierVoteIsRequired(proposal) && !k.GetCertifierVoted(ctx, proposalID) {
		return k.AddCertifierVote(ctx, proposalID, voterAddr, options)
	}

	// update certifier vote
	proposalMsgs, err := proposal.GetMsgs()
	if err != nil {
		return err
	}
	for _, proposalmsg := range proposalMsgs {
		if legacyMsg, ok := proposalmsg.(*govtypesv1.MsgExecLegacyContent); ok {
			// check that the content struct can be unmarshalled
			content, err := govtypesv1.LegacyContentFromMessage(legacyMsg)
			if err != nil {
				return err
			}
			if content.ProposalType() == shieldtypes.ProposalTypeShieldClaim &&
				proposal.Status == govtypesv1.StatusVotingPeriod &&
				!k.IsCertifiedIdentity(ctx, voterAddr) {
				return sdkerrors.Wrapf(govtypes.ErrInvalidVote, "'%s' is not a certified identity", voterAddr)
			}
		}
	}

	txhash := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
	vote := govtypesv1.NewVote(proposalID, voterAddr, options, metadata)
	k.SetVote(ctx, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeProposalVote,
			sdk.NewAttribute(govtypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(govtypes.AttributeKeyVoter, voterAddr.String()),
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
	k.IterateVotes(ctx, proposalID, func(vote govtypesv1.Vote) bool {
		addr, err := sdk.AccAddressFromBech32(vote.Voter)
		if err != nil {
			panic(err)
		}
		k.deleteVote(ctx, proposalID, addr)
		return false
	})
}

// AddCertifierVote add a certifier vote
func (k Keeper) AddCertifierVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypesv1.WeightedVoteOptions) error {
	if !k.IsCertifier(ctx, voterAddr) {
		return sdkerrors.Wrapf(govtypes.ErrInvalidVote, "%s is not a certified identity", voterAddr)
	}

	txhash := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
	vote := govtypesv1.NewVote(proposalID, voterAddr, options, "")
	k.SetVote(ctx, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCertVote,
			sdk.NewAttribute(govtypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(govtypes.AttributeKeyVoter, voterAddr.String()),
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
	store.Set(v1.CertVotesKey(proposalID), govtypes.GetProposalIDBytes(proposalID))
}

// GetCertifierVoted determine cert vote for custom proposal types have finished
func (k Keeper) GetCertifierVoted(ctx sdk.Context, proposalID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(v1.CertVotesKey(proposalID))
	return bz != nil
}
