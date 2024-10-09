package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// AddVote Adds a vote on a specific proposal.
func (k Keeper) AddVote(ctx context.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypesv1.WeightedVoteOptions, metadata string) error {
	// Check if proposal is in voting period.
	inVotingPeriod, err := k.VotingPeriodProposals.Has(ctx, proposalID)
	if err != nil {
		return err
	}

	if !inVotingPeriod {
		return errors.Wrapf(govtypes.ErrInactiveProposal, "%d", proposalID)
	}

	err = k.assertMetadataLength(metadata)
	if err != nil {
		return err
	}

	for _, option := range options {
		if !govtypesv1.ValidWeightedVoteOption(*option) {
			return errors.Wrapf(govtypes.ErrInvalidVote, "%s", option)
		}
	}

	// Add certifier vote
	if k.CertifierVoteIsRequired(ctx, proposalID) && !k.GetCertifierVoted(ctx, proposalID) {
		return k.AddCertifierVote(ctx, proposalID, voterAddr, options)
	}

	vote := govtypesv1.NewVote(proposalID, voterAddr, options, metadata)
	err = k.Votes.Set(ctx, collections.Join(proposalID, voterAddr), vote)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeProposalVote,
			sdk.NewAttribute(govtypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(govtypes.AttributeKeyVoter, voterAddr.String()),
		),
	)

	return nil
}

// deleteVotes deletes all the votes from a given proposalID.
func (keeper Keeper) deleteVotes(ctx context.Context, proposalID uint64) error {
	rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](proposalID)
	err := keeper.Votes.Clear(ctx, rng)
	if err != nil {
		return err
	}

	return nil
}

// AddCertifierVote add a certifier vote
func (k Keeper) AddCertifierVote(ctx context.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypesv1.WeightedVoteOptions) error {
	if !k.IsCertifier(ctx, voterAddr) {
		return errors.Wrapf(govtypes.ErrInvalidVote, "%s is not a certified identity", voterAddr)
	}

	txhash := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))
	vote := govtypesv1.NewVote(proposalID, voterAddr, options, "")
	err := k.Votes.Set(ctx, collections.Join(proposalID, voterAddr), vote)
	if err != nil {
		return err
	}

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
func (k Keeper) SetCertVote(ctx context.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(typesv1.CertVotesKey(proposalID), govtypes.GetProposalIDBytes(proposalID))
}

// GetCertifierVoted determine cert vote for custom proposal types have finished
func (k Keeper) GetCertifierVoted(ctx context.Context, proposalID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(typesv1.CertVotesKey(proposalID))
	return bz != nil
}
