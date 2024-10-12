package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

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
	voted, err := k.GetCertifierVoted(ctx, proposalID)
	if err != nil {
		return err
	}
	required, err := k.CertifierVoteIsRequired(ctx, proposalID)
	if err != nil {
		return err
	}
	if required && !voted {
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
	isCertifier, err := k.IsCertifier(ctx, voterAddr)
	if err != nil {
		return err
	}
	if !isCertifier {
		return errors.Wrapf(govtypes.ErrInvalidVote, "%s is not a certified identity", voterAddr)
	}

	vote := govtypesv1.NewVote(proposalID, voterAddr, options, "")
	err = k.Votes.Set(ctx, collections.Join(proposalID, voterAddr), vote)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCertVote,
			sdk.NewAttribute(govtypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(govtypes.AttributeKeyVoter, voterAddr.String()),
		),
	)
	return nil
}

func (k Keeper) SetCertifierVoted(ctx sdk.Context, proposalID uint64) error {
	err := k.SetCertVote(ctx, proposalID)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetCertVote,
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)
	return nil
}

// SetCertVote sets a cert vote to the gov store
func (k Keeper) SetCertVote(ctx context.Context, proposalID uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(typesv1.CertVotesKey(proposalID), typesv1.GetProposalIDBytes(proposalID))
}

// GetCertifierVoted determine cert vote for custom proposal types have finished
func (k Keeper) GetCertifierVoted(ctx context.Context, proposalID uint64) (bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	return store.Has(typesv1.CertVotesKey(proposalID))
}
