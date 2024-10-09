package keeper

import (
	"context"

	"cosmossdk.io/math"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
)

func (k Keeper) ActivateVotingPeriod(ctx sdk.Context, proposal govtypesv1.Proposal) {
	startTime := ctx.BlockHeader().Time
	proposal.VotingStartTime = &startTime
	votingPeriod := k.GetParams(ctx).VotingPeriod
	oldVotingEndTime := proposal.VotingEndTime
	endTime := proposal.VotingStartTime.Add(*votingPeriod)
	proposal.VotingEndTime = &endTime
	oldDepositEndTime := proposal.DepositEndTime

	// Default case: for plain text proposals, community pool spend proposals;
	// and second round of software upgrade, certifier update and shield claim
	// proposals.
	if k.GetCertifierVoted(ctx, proposal.Id) {
		k.RemoveFromActiveProposalQueue(ctx, proposal.Id, *oldVotingEndTime)
	} else {
		proposal.DepositEndTime = &endTime
	}
	proposal.Status = govtypesv1.StatusVotingPeriod

	k.SetProposal(ctx, proposal)
	k.RemoveFromInactiveProposalQueue(ctx, proposal.Id, *oldDepositEndTime)
	k.InsertActiveProposalQueue(ctx, proposal.Id, *proposal.VotingEndTime)
}

// IsCertifier checks if the input address is a certifier.
func (k Keeper) IsCertifier(ctx context.Context, addr sdk.AccAddress) (bool, error) {
	return k.CertKeeper.IsCertifier(ctx, addr)
}

// TotalBondedByCertifiedIdentities calculates the amount of total bonded stakes by certified identities.
func (k Keeper) TotalBondedByCertifiedIdentities(ctx sdk.Context) math.Int {
	bonded := sdk.ZeroInt()
	for _, identity := range k.CertKeeper.GetCertifiedIdentities(ctx) {
		k.stakingKeeper.IterateDelegations(ctx, identity, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			val, found := k.stakingKeeper.GetValidator(ctx, delegation.GetValidatorAddr())
			if !found {
				return false
			}
			bonded = bonded.Add(delegation.GetShares().Quo(val.GetDelegatorShares()).MulInt(val.GetBondedTokens()).TruncateInt())
			return false
		})
	}
	return bonded
}

func (k Keeper) CertifierVoteIsRequired(ctx context.Context, proposalID uint64) bool {
	proposal, err := k.Proposals.Get(ctx, proposalID)
	proposalMsgs, err := proposal.GetMsgs()
	if err != nil {
		return false
	}

	for _, proposalmsg := range proposalMsgs {
		// upgrade msg need certifier vote
		if sdk.MsgTypeURL(proposalmsg) == sdk.MsgTypeURL(&upgradetypes.MsgSoftwareUpgrade{}) {
			return true
		}

		if legacyMsg, ok := proposalmsg.(*govtypesv1.MsgExecLegacyContent); ok {
			// check that the content struct can be unmarshalled
			content, err := govtypesv1.LegacyContentFromMessage(legacyMsg)
			if err != nil {
				return false
			}
			switch content.(type) {
			// nolint
			case *upgradetypes.SoftwareUpgradeProposal, *certtypes.CertifierUpdateProposal:
				return true
			default:
				return false
			}
		}
	}

	return false
}

// assertMetadataLength returns an error if given metadata length
// is greater than a pre-defined maxMetadataLen.
func (k Keeper) assertMetadataLength(metadata string) error {
	if metadata != "" && uint64(len(metadata)) > k.config.MaxMetadataLen {
		return govtypes.ErrMetadataTooLong.Wrapf("got metadata with length %d", len(metadata))
	}
	return nil
}
