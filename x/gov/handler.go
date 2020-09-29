package gov

import (
	"encoding/hex"
	"fmt"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/gov/internal/keeper"
	"github.com/certikfoundation/shentu/x/gov/internal/types"
	"github.com/certikfoundation/shentu/x/shield"
)

// NewHandler handles all "gov" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case gov.MsgDeposit:
			return handleMsgDeposit(ctx, k, msg)

		case gov.MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, k, msg)

		case gov.MsgVote:
			return handleMsgVote(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized gov message type: %T", msg)
		}
	}
}

func handleMsgDeposit(ctx sdk.Context, k keeper.Keeper, msg gov.MsgDeposit) (*sdk.Result, error) {
	votingStarted, err := k.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor.String()),
			sdk.NewAttribute(AttributeTxHash, hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))),
		),
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govtypes.EventTypeProposalDeposit,
				sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalID)),
			),
		)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSubmitProposal(ctx sdk.Context, k keeper.Keeper, msg gov.MsgSubmitProposal) (*sdk.Result, error) {
	var initialDepositAmount = msg.InitialDeposit.AmountOf(common.MicroCTKDenom)
	var depositParams = k.GetDepositParams(ctx)
	var minimalInitialDepositAmount = depositParams.MinInitialDeposit.AmountOf(common.MicroCTKDenom)
	// Check if delegator proposal reach the bar, current bar is 0 ctk.
	if initialDepositAmount.LT(minimalInitialDepositAmount) && !k.IsCouncilMember(ctx, msg.Proposer) {
		return &sdk.Result{}, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"insufficient initial deposits amount: %v, minimum: %v",
			initialDepositAmount,
			minimalInitialDepositAmount,
		)
	}

	err := validateProposalByType(ctx, k, msg)
	if err != nil {
		return &sdk.Result{}, err
	}

	proposal, err := k.SubmitProposal(ctx, msg.Content, msg.Proposer)
	if err != nil {
		return nil, err
	}

	// Skip deposit period for proposals of council members.
	isVotingPeriodActivated := k.ActivateCouncilProposalVotingPeriod(ctx, proposal)
	if !isVotingPeriodActivated {
		// Non council members can add deposit to their newly submitted proposals.
		isVotingPeriodActivated, err = k.AddDeposit(ctx, proposal.ProposalID, msg.Proposer, msg.InitialDeposit)
		if err != nil {
			return nil, err
		}
	}

	if err := updateAfterSubmitProposal(ctx, k, proposal); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proposer.String()),
		),
	)

	submitEvent := sdk.NewEvent(
		govtypes.EventTypeSubmitProposal,
		sdk.NewAttribute(govtypes.AttributeKeyProposalType, msg.Content.ProposalType()),
		sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
	)
	if isVotingPeriodActivated {
		submitEvent = submitEvent.AppendAttributes(
			sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalID)),
		)
	}
	ctx.EventManager().EmitEvent(submitEvent)
	return &sdk.Result{
		Data:   gov.GetProposalIDBytes(proposal.ProposalID),
		Events: ctx.EventManager().Events(),
	}, nil
}

func updateAfterSubmitProposal(ctx sdk.Context, k keeper.Keeper, proposal types.Proposal) error {
	if proposal.ProposalType() == shield.ProposalTypeShieldClaim {
		c := proposal.Content.(shield.ClaimProposal)
		if err := k.ShieldKeeper.ClaimLock(ctx, c.PoolID, c.Loss, c.PurchaseTxHash, proposal.ProposalID); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func validateProposalByType(ctx sdk.Context, k keeper.Keeper, msg gov.MsgSubmitProposal) error {
	switch c := msg.Content.(type) {
	case cert.CertifierUpdateProposal:
		if !k.IsCertifier(ctx, msg.Proposer) {
			return sdkerrors.Wrapf(cert.ErrUnqualifiedCertifier,
				"%s is not a certifier", msg.Proposer,
			)
		}
		isCertifier := k.IsCertifier(ctx, c.Certifier)
		if c.AddOrRemove == cert.Add && isCertifier {
			return cert.ErrCertifierAlreadyExists
		}
		if c.AddOrRemove == cert.Remove && !isCertifier {
			return sdkerrors.Wrapf(cert.ErrUnqualifiedCertifier,
				"%s is not a certifier", msg.Proposer,
			)
		}
		if c.Alias != "" && k.CertKeeper.HasCertifierAlias(ctx, c.Alias) {
			return cert.ErrRepeatedAlias
		}

	case upgrade.SoftwareUpgradeProposal:
		return k.UpgradeKeeper.ValidatePlan(ctx, c.Plan)

	case shield.ClaimProposal:
		// check initial deposit >= max(<loss>*ClaimDepositRate, MinimumClaimDeposit)
		initialDepositAmount := msg.InitialDeposit.AmountOf(common.MicroCTKDenom).ToDec()
		lossAmount := c.Loss.AmountOf(common.MicroCTKDenom).ToDec()
		claimProposalParams := k.ShieldKeeper.GetClaimProposalParams(ctx)
		depositRate := claimProposalParams.DepositRate
		minDeposit := claimProposalParams.MinDeposit.AmountOf(common.MicroCTKDenom).ToDec()
		if initialDepositAmount.LT(lossAmount.Mul(depositRate)) && initialDepositAmount.LT(minDeposit) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInsufficientFunds,
				"insufficient initial deposits amount: %v, minimum: max(%v, %v)",
				initialDepositAmount, lossAmount.Mul(depositRate), minDeposit,
			)
		}

		// check shield >= loss
		purchase, err := k.ShieldKeeper.GetPurchase(ctx, c.PurchaseTxHash)
		if err != nil {
			return err
		}
		if !purchase.Shield.IsAllGTE(c.Loss) {
			return fmt.Errorf("insufficient shield: %s, loss: %s", purchase.Shield, c.Loss)
		}

		// check the purchase is not expired
		if purchase.ClaimPeriodEndTime.Before(ctx.BlockTime()) {
			return fmt.Errorf("after claim period end time: %s", purchase.ClaimPeriodEndTime)
		}
		return nil

	default:
		return nil
	}
	return nil
}

func handleMsgVote(ctx sdk.Context, k keeper.Keeper, msg gov.MsgVote) (*sdk.Result, error) {
	err := k.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
