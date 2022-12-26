package keeper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) govtypes.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ govtypes.MsgServer = msgServer{}

func (k msgServer) SubmitProposal(goCtx context.Context, msg *govtypes.MsgSubmitProposal) (*govtypes.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := validateProposalByType(ctx, k.Keeper, msg)
	if err != nil {
		return nil, err
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, msg.GetContent())
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, govtypes.ModuleName, "proposal")

	votingStarted, err := k.Keeper.AddDeposit(ctx, proposal.ProposalId, msg.GetProposer(), msg.GetInitialDeposit())
	if err != nil {
		return nil, err
	}

	if err := updateAfterSubmitProposal(ctx, k.Keeper, proposal); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.GetProposer().String()),
		),
	)

	submitEvent := sdk.NewEvent(govtypes.EventTypeSubmitProposal, sdk.NewAttribute(govtypes.AttributeKeyProposalType, msg.GetContent().ProposalType()))
	if votingStarted {
		submitEvent = submitEvent.AppendAttributes(
			sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalId)),
		)
	}

	ctx.EventManager().EmitEvent(submitEvent)
	return &govtypes.MsgSubmitProposalResponse{
		ProposalId: proposal.ProposalId,
	}, nil
}

func (k msgServer) Vote(goCtx context.Context, msg *govtypes.MsgVote) (*govtypes.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, err := sdk.AccAddressFromBech32(msg.Voter)
	if err != nil {
		return nil, err
	}

	// Custom proposal type, need cert to vote first
	isNeedCertVote, err := k.IsNeedCertVote(ctx, msg.ProposalId)
	if err != nil {
		return nil, err
	}
	if isNeedCertVote && !k.IsCertifierVoted(ctx, msg.ProposalId) {
		return k.AddCertVote(ctx, msg.ProposalId, accAddr)
	}

	err = k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, govtypes.NewNonSplitVoteOption(msg.Option))
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounterWithLabels(
		[]string{govtypes.ModuleName, "vote"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter),
		),
	)

	return &govtypes.MsgVoteResponse{}, nil
}

func (k msgServer) VoteWeighted(goCtx context.Context, msg *govtypes.MsgVoteWeighted) (*govtypes.MsgVoteWeightedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, accErr := sdk.AccAddressFromBech32(msg.Voter)
	if accErr != nil {
		return nil, accErr
	}

	// Custom proposal type, need cert to vote first
	isNeedCertVote, err := k.IsNeedCertVote(ctx, msg.ProposalId)
	if err != nil {
		return nil, err
	}
	if isNeedCertVote {
		return nil, sdkerrors.Wrap(govtypes.ErrInvalidVote, "need cert vote fist")
	}

	err = k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, msg.Options)
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounterWithLabels(
		[]string{govtypes.ModuleName, "vote"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter),
		),
	)

	return &govtypes.MsgVoteWeightedResponse{}, nil
}

func (k msgServer) Deposit(goCtx context.Context, msg *govtypes.MsgDeposit) (*govtypes.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	votingStarted, err := k.Keeper.AddDeposit(ctx, msg.ProposalId, accAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounterWithLabels(
		[]string{govtypes.ModuleName, "deposit"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor),
		),
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govtypes.EventTypeProposalDeposit,
				sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalId)),
			),
		)
	}

	return &govtypes.MsgDepositResponse{}, nil
}

func validateProposalByType(ctx sdk.Context, k Keeper, msg *govtypes.MsgSubmitProposal) error {
	switch c := msg.GetContent().(type) {
	case *certtypes.CertifierUpdateProposal:
		if c.Alias != "" && k.CertKeeper.HasCertifierAlias(ctx, c.Alias) {
			return certtypes.ErrRepeatedAlias
		}

	case shieldtypes.ShieldClaimProposal:
		// check initial deposit >= max(<loss>*ClaimDepositRate, MinimumClaimDeposit)
		denom := k.BondDenom(ctx)
		initialDepositAmount := msg.InitialDeposit.AmountOf(denom).ToDec()
		lossAmount := c.Loss.AmountOf(denom)
		lossAmountDec := lossAmount.ToDec()
		claimProposalParams := k.ShieldKeeper.GetClaimProposalParams(ctx)
		depositRate := claimProposalParams.DepositRate
		minDeposit := claimProposalParams.MinDeposit.AmountOf(denom).ToDec()
		if initialDepositAmount.LT(lossAmountDec.Mul(depositRate)) || initialDepositAmount.LT(minDeposit) {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInsufficientFunds,
				"insufficient initial deposits amount: %v, minimum: max(%v, %v)",
				initialDepositAmount, lossAmountDec.Mul(depositRate), minDeposit,
			)
		}

		// check shield >= loss
		proposerAddr, err := sdk.AccAddressFromBech32(c.Proposer)
		if err != nil {
			return err
		}
		purchaseList, found := k.ShieldKeeper.GetPurchaseList(ctx, c.PoolId, proposerAddr)
		if !found {
			return shieldtypes.ErrPurchaseNotFound
		}
		purchase, found := k.ShieldKeeper.GetPurchase(purchaseList, c.PurchaseId)
		if !found {
			return shieldtypes.ErrPurchaseNotFound
		}
		if !purchase.Shield.GTE(lossAmount) {
			return fmt.Errorf("insufficient shield: %s, loss: %s", purchase.Shield, c.Loss)
		}

		// check the purchaseList is not expired
		if purchase.ProtectionEndTime.Before(ctx.BlockTime()) {
			return fmt.Errorf("after protection end time: %s", purchase.ProtectionEndTime)
		}
		return nil

	default:
		return nil
	}
	return nil
}

func updateAfterSubmitProposal(ctx sdk.Context, k Keeper, proposal govtypes.Proposal) error {
	if proposal.ProposalType() == shieldtypes.ProposalTypeShieldClaim {
		c := proposal.GetContent().(*shieldtypes.ShieldClaimProposal)
		lockPeriod := k.GetVotingParams(ctx).VotingPeriod * 2
		proposerAddr, err := sdk.AccAddressFromBech32(c.Proposer)
		if err != nil {
			return err
		}
		return k.ShieldKeeper.SecureCollaterals(ctx, c.PoolId, proposerAddr, c.PurchaseId, c.Loss, lockPeriod)
	}
	return nil
}
