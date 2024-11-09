package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	"github.com/shentufoundation/shentu/v2/x/gov/types"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

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

func (k Keeper) SetCertifierVoted(ctx context.Context, proposalID uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.SetCertVote(ctx, proposalID)
	if err != nil {
		return err
	}

	sdkCtx.EventManager().EmitEvent(
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

// IsCertifier checks if the input address is a certifier.
func (k Keeper) IsCertifier(ctx context.Context, addr sdk.AccAddress) (bool, error) {
	return k.certKeeper.IsCertifier(ctx, addr)
}

func (k Keeper) CertifierVoteIsRequired(ctx context.Context, proposalID uint64) (bool, error) {
	proposal, err := k.Proposals.Get(ctx, proposalID)
	proposalMsgs, err := proposal.GetMsgs()
	if err != nil {
		return false, err
	}

	for _, proposalmsg := range proposalMsgs {
		// upgrade msg need certifier vote
		if sdk.MsgTypeURL(proposalmsg) == sdk.MsgTypeURL(&upgradetypes.MsgSoftwareUpgrade{}) {
			return true, nil
		}

		if legacyMsg, ok := proposalmsg.(*govtypesv1.MsgExecLegacyContent); ok {
			// check that the content struct can be unmarshalled
			content, err := govtypesv1.LegacyContentFromMessage(legacyMsg)
			if err != nil {
				return false, err
			}
			switch content.(type) {
			// nolint
			case *upgradetypes.SoftwareUpgradeProposal, *certtypes.CertifierUpdateProposal:
				return true, nil
			default:
				return false, nil
			}
		}
	}

	return false, nil
}
