package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
	"strconv"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SubmitProposal(goCtx context.Context, msg *types.MsgSubmitProposal) (*types.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := validateProposalByType(ctx, k.Keeper, msg)
	if err != nil {
		return nil, err
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, msg.GetContent())
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, types.ModuleName, "proposal")

	votingStarted, err := k.Keeper.AddDeposit(ctx, proposal.ProposalId, msg.GetProposer(), msg.GetInitialDeposit())
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.GetProposer().String()),
		),
	)

	submitEvent := sdk.NewEvent(types.EventTypeSubmitProposal, sdk.NewAttribute(types.AttributeKeyProposalType, msg.GetContent().ProposalType()))
	if votingStarted {
		submitEvent = submitEvent.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalId)),
		)
	}

	ctx.EventManager().EmitEvent(submitEvent)
	return &types.MsgSubmitProposalResponse{
		ProposalId: proposal.ProposalId,
	}, nil
}

func (k msgServer) Vote(goCtx context.Context, msg *types.MsgVote) (*types.MsgVoteResponse, error) {
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

	err = k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, types.NewNonSplitVoteOption(msg.Option))
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, "vote"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter),
		),
	)

	return &types.MsgVoteResponse{}, nil
}

func (k msgServer) VoteWeighted(goCtx context.Context, msg *types.MsgVoteWeighted) (*types.MsgVoteWeightedResponse, error) {
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
		return nil, sdkerrors.Wrap(types.ErrInvalidVote, "need cert vote fist")
	}

	err = k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, msg.Options)
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, "vote"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter),
		),
	)

	return &types.MsgVoteWeightedResponse{}, nil
}

func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
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
		[]string{types.ModuleName, "deposit"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor),
		),
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposalDeposit,
				sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalId)),
			),
		)
	}

	return &types.MsgDepositResponse{}, nil
}

func validateProposalByType(ctx sdk.Context, k Keeper, msg *types.MsgSubmitProposal) error {
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
