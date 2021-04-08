package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	certtypes "github.com/certikfoundation/shentu/x/cert/types"
	"github.com/certikfoundation/shentu/x/gov/types"
	shieldtypes "github.com/certikfoundation/shentu/x/shield/types"
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

func (k msgServer) SubmitProposal(goCtx context.Context, msg *govtypes.MsgSubmitProposal) (*govtypes.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var initialDepositAmount = msg.InitialDeposit.AmountOf(k.stakingKeeper.BondDenom(ctx))
	var depositParams = k.GetDepositParams(ctx)
	var minimalInitialDepositAmount = depositParams.MinInitialDeposit.AmountOf(k.stakingKeeper.BondDenom(ctx))
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

func updateAfterSubmitProposal(ctx sdk.Context, k Keeper, proposal types.Proposal) error {
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor),
			//sdk.NewAttribute(types.AttributeTxHash, hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))),
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
