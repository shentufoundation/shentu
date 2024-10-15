package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
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

	required, err := k.CertifierVoteIsRequired(ctx, proposalID)
	if err != nil {
		return err
	}
	if required {
		voted, err := k.GetCertifierVoted(ctx, proposalID)
		if err != nil {
			return err
		}
		if !voted {
			return k.AddCertifierVote(ctx, proposalID, voterAddr, options)
		}
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

// DeleteVotes deletes all the votes from a given proposalID.
func (keeper Keeper) DeleteVotes(ctx context.Context, proposalID uint64) error {
	rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](proposalID)
	err := keeper.Votes.Clear(ctx, rng)
	if err != nil {
		return err
	}
	return nil
}
