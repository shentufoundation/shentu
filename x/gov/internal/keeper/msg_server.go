package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/gov/internal/types"
	"github.com/certikfoundation/shentu/x/shield"
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

	var initialDepositAmount = msg.InitialDeposit.AmountOf(common.MicroCTKDenom)
	var depositParams = k.GetDepositParams(ctx)
	var minimalInitialDepositAmount = depositParams.MinInitialDeposit.AmountOf(common.MicroCTKDenom)
	// Check if delegator proposal reach the bar, current bar is 0 ctk.
	if initialDepositAmount.LT(minimalInitialDepositAmount) && !k.IsCouncilMember(ctx, msg.GetProposer()) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"insufficient initial deposits amount: %v, minimum: %v",
			initialDepositAmount,
			minimalInitialDepositAmount,
		)
	}

	err := validateProposalByType(ctx, k.Keeper, msg)
	if err != nil {
		return nil, err
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, msg.GetContent(), msg.GetProposer())
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounter(1, govtypes.ModuleName, "proposal")

	// Skip deposit period for proposals of council members.
	isVotingPeriodActivated := k.ActivateCouncilProposalVotingPeriod(ctx, proposal)
	if !isVotingPeriodActivated {
		// Non council members can add deposit to their newly submitted proposals.
		isVotingPeriodActivated, err = k.AddDeposit(ctx, proposal.ProposalId, msg.GetProposer(), msg.GetInitialDeposit())
		if err != nil {
			return nil, err
		}
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

	submitEvent := sdk.NewEvent(
		govtypes.EventTypeSubmitProposal,
		sdk.NewAttribute(govtypes.AttributeKeyProposalType, msg.GetContent().ProposalType()),
		sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalId)),
	)
	if isVotingPeriodActivated {
		submitEvent = submitEvent.AppendAttributes(
			sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalId)),
		)
	}

	ctx.EventManager().EmitEvent(submitEvent)
	return &govtypes.MsgSubmitProposalResponse{
		ProposalId: proposal.ProposalId,
	}, nil
}

func validateProposalByType(ctx sdk.Context, k Keeper, msg *govtypes.MsgSubmitProposal) error {
	switch c := msg.GetContent().(type) {
	case cert.CertifierUpdateProposal:
		if c.Alias != "" && k.CertKeeper.HasCertifierAlias(ctx, c.Alias) {
			return cert.ErrRepeatedAlias
		}

	case *upgrade.SoftwareUpgradeProposal:
		return k.UpgradeKeeper.ValidatePlan(ctx, c.Plan)

	case shield.ClaimProposal:
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
		purchaseList, found := k.ShieldKeeper.GetPurchaseList(ctx, c.PoolID, c.Proposer)
		if !found {
			return shield.ErrPurchaseNotFound
		}
		purchase, found := shield.GetPurchase(purchaseList, c.PurchaseID)
		if !found {
			return shield.ErrPurchaseNotFound
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

func updateAfterSubmitProposal(ctx sdk.Context, k Keeper, proposal types.Proposal) error {
	if proposal.ProposalType() == shield.ProposalTypeShieldClaim {
		c := proposal.GetContent().(shield.ClaimProposal)
		lockPeriod := k.GetVotingParams(ctx).VotingPeriod * 2
		return k.ShieldKeeper.SecureCollaterals(ctx, c.PoolID, c.Proposer, c.PurchaseID, c.Loss, lockPeriod)
	}
	return nil
}

func (k msgServer) Vote(goCtx context.Context, msg *govtypes.MsgVote) (*govtypes.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, accErr := sdk.AccAddressFromBech32(msg.Voter)
	if accErr != nil {
		return nil, accErr
	}
	err := k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, msg.Option)
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
			sdk.NewAttribute(types.AttributeTxHash, hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))),
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
