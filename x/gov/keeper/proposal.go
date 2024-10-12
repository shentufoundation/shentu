package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// ActivateVotingPeriod activates the voting period of a proposal
func (k Keeper) ActivateVotingPeriod(ctx context.Context, proposal govtypesv1.Proposal) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	startTime := sdkCtx.BlockHeader().Time
	proposal.VotingStartTime = &startTime
	var votingPeriod *time.Duration
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	if proposal.Expedited {
		votingPeriod = params.ExpeditedVotingPeriod
	} else {
		votingPeriod = params.VotingPeriod
	}
	endTime := proposal.VotingStartTime.Add(*votingPeriod)
	proposal.VotingEndTime = &endTime
	proposal.Status = govtypesv1.StatusVotingPeriod
	err = k.SetProposal(ctx, proposal)
	if err != nil {
		return err
	}

	err = k.InactiveProposalsQueue.Remove(ctx, collections.Join(*proposal.DepositEndTime, proposal.Id))
	if err != nil {
		return err
	}
	return k.ActiveProposalsQueue.Set(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id), proposal.Id)
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

// assertMetadataLength returns an error if given metadata length
// is greater than a pre-defined maxMetadataLen.
func (k Keeper) assertMetadataLength(metadata string) error {
	if metadata != "" && uint64(len(metadata)) > k.config.MaxMetadataLen {
		return govtypes.ErrMetadataTooLong.Wrapf("got metadata with length %d", len(metadata))
	}
	return nil
}
