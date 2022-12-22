package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/shentufoundation/shentu/v2/x/gov/types"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// AddVote adds a vote on a specific proposal
func (keeper Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypes.WeightedVoteOptions) error {
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return sdkerrors.Wrapf(govtypes.ErrUnknownProposal, "%d", proposalID)
	}
	if proposal.Status != govtypes.StatusVotingPeriod {
		return sdkerrors.Wrapf(govtypes.ErrInactiveProposal, "%d", proposalID)
	}

	for _, option := range options {
		if !govtypes.ValidWeightedVoteOption(option) {
			return sdkerrors.Wrap(govtypes.ErrInvalidVote, option.String())
		}
	}

	vote := govtypes.NewVote(proposalID, voterAddr, options)
	keeper.SetVote(ctx, vote)

	// called after a vote on a proposal is cast
	keeper.AfterProposalVote(ctx, proposalID, voterAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeProposalVote,
			sdk.NewAttribute(govtypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

func (keeper Keeper) IsNeedCertVote(ctx sdk.Context, proposalID uint64) (bool, error) {
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return false, sdkerrors.Wrapf(govtypes.ErrUnknownProposal, "%d", proposalID)
	}
	switch proposal.GetContent().(type) {
	case *upgradetypes.SoftwareUpgradeProposal, *certtypes.CertifierUpdateProposal, shieldtypes.ShieldClaimProposal:
		return true, nil
	default:
		return false, nil
	}
}

func (keeper Keeper) AddCertVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (*govtypes.MsgVoteResponse, error) {
	if keeper.IsCertifier(ctx, voterAddr) {
		keeper.SetCertifierVote(ctx, proposalID)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeSetCertVote,
				sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			),
		)
		return &govtypes.MsgVoteResponse{}, nil
	}
	return nil, sdkerrors.Wrapf(govtypes.ErrInvalidVote, "%s is not a certified identity", voterAddr)
}

// SetCertifierVote sets a Certifier to the gov store
func (keeper Keeper) SetCertifierVote(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Set(types.CertVotesKey(proposalID), govtypes.GetProposalIDBytes(proposalID))
}

// IsCertifierVote determine cert vote for custom proposal types
func (keeper Keeper) isCertifierVoted(ctx sdk.Context, proposalID uint64) bool {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.CertVotesKey(proposalID))
	if bz == nil {
		return false
	}
	return true
}
