package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	v046 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v046"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
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

// isValidator checks if the input address is a validator.
func (k Keeper) isValidator(ctx sdk.Context, addr sdk.AccAddress) bool {
	isValidator := false
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		if validator.GetOperator().Equals(addr) {
			isValidator = true
			return true
		}
		return false
	})
	return isValidator
}

// IsCertifier checks if the input address is a certifier.
func (k Keeper) IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.CertKeeper.IsCertifier(ctx, addr)
}

// IsCouncilMember checks if the address is either a validator or a certifier.
func (k Keeper) IsCouncilMember(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.isValidator(ctx, addr) || k.IsCertifier(ctx, addr)
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

// SubmitProposal creates a new proposal with given content.
func (k Keeper) SubmitProposal(ctx sdk.Context, messages []sdk.Msg, metadata string) (govtypesv1.Proposal, error) {
	err := k.assertMetadataLength(metadata)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	// Will hold a comma-separated string of all Msg type URLs.
	msgsStr := ""

	// Loop through all messages and confirm that each has a handler and the gov module account
	// as the only signer
	for _, msg := range messages {
		msgsStr += fmt.Sprintf(",%s", sdk.MsgTypeURL(msg))

		// perform a basic validation of the message
		if err := msg.ValidateBasic(); err != nil {
			return govtypesv1.Proposal{}, sdkerrors.Wrap(govtypes.ErrInvalidProposalMsg, err.Error())
		}

		signers := msg.GetSigners()
		if len(signers) != 1 {
			return govtypesv1.Proposal{}, govtypes.ErrInvalidSigner
		}

		// assert that the governance module account is the only signer of the messages
		if !signers[0].Equals(k.GetGovernanceAccount(ctx).GetAddress()) {
			return govtypesv1.Proposal{}, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, signers[0].String())
		}

		// use the msg service router to see that there is a valid route for that message.
		handler := k.router.Handler(msg)
		if handler == nil {
			return govtypesv1.Proposal{}, sdkerrors.Wrap(govtypes.ErrUnroutableProposalMsg, sdk.MsgTypeURL(msg))
		}

		// Only if it's a MsgExecLegacyContent do we try to execute the
		// proposal in a cached context.
		// For other Msgs, we do not verify the proposal messages any further.
		// They may fail upon execution.
		// ref: https://github.com/cosmos/cosmos-sdk/pull/10868#discussion_r784872842
		if msg, ok := msg.(*govtypesv1.MsgExecLegacyContent); ok {
			cacheCtx, _ := ctx.CacheContext()
			if _, err := handler(cacheCtx, msg); err != nil {
				return govtypesv1.Proposal{}, sdkerrors.Wrap(govtypes.ErrNoProposalHandlerExists, err.Error())
			}
		}

	}

	proposalID, err := k.GetProposalID(ctx)
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := k.GetDepositParams(ctx).MaxDepositPeriod

	proposal, err := govtypesv1.NewProposal(messages, proposalID, metadata, submitTime, submitTime.Add(*depositPeriod))
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	k.SetProposal(ctx, proposal)
	k.InsertInactiveProposalQueue(ctx, proposalID, *proposal.DepositEndTime)
	k.SetProposalID(ctx, proposalID+1)

	// called right after a proposal is submitted
	k.AfterProposalSubmission(ctx, proposalID)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeSubmitProposal,
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(govtypes.AttributeKeyProposalMessages, msgsStr),
		),
	)

	return proposal, nil
}

func (k Keeper) HasSecurityVoting(p govtypesv1.Proposal) bool {
	legacyProposal, err := v046.ConvertToLegacyProposal(p)
	if err != nil {
		return false
	}
	switch legacyProposal.GetContent().(type) {
	case *upgradetypes.SoftwareUpgradeProposal, *certtypes.CertifierUpdateProposal, *shieldtypes.ShieldClaimProposal:
		return true
	default:
		return false
	}
}

// ActivateVotingPeriodCustom switches proposals to voting period for customization.
func (k Keeper) ActivateVotingPeriodCustom(ctx sdk.Context, proposal govtypesv1.Proposal, addr sdk.AccAddress) bool {
	legacyProposal, err := v046.ConvertToLegacyProposal(proposal)
	if err != nil {
		return false
	}
	if !k.IsCertifier(ctx, addr) && legacyProposal.ProposalType() != shieldtypes.ProposalTypeShieldClaim {
		return false
	}
	if k.IsCertifier(ctx, addr) && k.HasSecurityVoting(proposal) {
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
