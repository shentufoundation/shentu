package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

func (k Keeper) ActivateVotingPeriod(ctx sdk.Context, proposal govtypesv1.Proposal) {
	startTime := ctx.BlockHeader().Time
	proposal.VotingStartTime = &startTime
	votingPeriod := k.GetVotingParams(ctx).VotingPeriod
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
func (k Keeper) IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.CertKeeper.IsCertifier(ctx, addr)
}

// IsCertifiedIdentity checks if the input address is a certified identity.
func (k Keeper) IsCertifiedIdentity(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.CertKeeper.IsCertified(ctx, addr.String(), "identity")
}

// TotalBondedByCertifiedIdentities calculates the amount of total bonded stakes by certified identities.
func (k Keeper) TotalBondedByCertifiedIdentities(ctx sdk.Context) sdk.Int {
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

func (k Keeper) CertifierVoteIsRequired(proposal govtypesv1.Proposal) bool {
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
			case *upgradetypes.SoftwareUpgradeProposal, *certtypes.CertifierUpdateProposal, *shieldtypes.ShieldClaimProposal:
				return true
			default:
				return false
			}
		}
	}

	return false
}

// ActivateVotingPeriodCustom switches proposals to voting period for customization.
func (k Keeper) ActivateVotingPeriodCustom(ctx sdk.Context, c govtypesv1beta1.Content, proposal govtypesv1.Proposal, addr sdk.AccAddress) bool {
	if !k.IsCertifier(ctx, addr) && c.ProposalType() != shieldtypes.ProposalTypeShieldClaim {
		return false
	}
	if k.IsCertifier(ctx, addr) && k.CertifierVoteIsRequired(proposal) {
		k.SetCertifierVoted(ctx, proposal.Id)
	}
	k.ActivateVotingPeriod(ctx, proposal)
	return true
}

// assertMetadataLength returns an error if given metadata length
// is greater than a pre-defined maxMetadataLen.
func (k Keeper) assertMetadataLength(metadata string) error {
	if metadata != "" && uint64(len(metadata)) > k.config.MaxMetadataLen {
		return govtypes.ErrMetadataTooLong.Wrapf("got metadata with length %d", len(metadata))
	}
	return nil
}
