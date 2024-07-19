package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) govtypesv1.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ govtypesv1.MsgServer = msgServer{}

func (k msgServer) SubmitProposal(goCtx context.Context, msg *govtypesv1.MsgSubmitProposal) (*govtypesv1.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposalMsgs, err := msg.GetMsgs()
	if err != nil {
		return nil, err
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, proposalMsgs, msg.Metadata)
	if err != nil {
		return nil, err
	}

	// Skip deposit period for proposals from certifier members or shield claim proposals.
	proposer, _ := sdk.AccAddressFromBech32(msg.GetProposer())
	votingStarted := false

	// update shield SecureCollaterals
	if len(proposalMsgs) == 1 {
		if legacyMsg, ok := proposalMsgs[0].(*govtypesv1.MsgExecLegacyContent); ok {
			// check that the content struct can be unmarshalled
			content, err := govtypesv1.LegacyContentFromMessage(legacyMsg)
			if err != nil {
				return nil, err
			}
			votingStarted = k.ActivateVotingPeriodCustom(ctx, content, proposal, proposer)

			if err = updateAfterSubmitProposal(ctx, k.Keeper, content); err != nil {
				return nil, err
			}
		}
	}

	if !votingStarted {
		votingStarted, err = k.Keeper.AddDeposit(ctx, proposal.Id, proposer, msg.GetInitialDeposit())
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, govtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.GetProposer()),
		),
	)

	if votingStarted {
		submitEvent := sdk.NewEvent(
			govtypes.EventTypeSubmitProposal,
			sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.Id)),
		)
		ctx.EventManager().EmitEvent(submitEvent)
	}

	return &govtypesv1.MsgSubmitProposalResponse{
		ProposalId: proposal.Id,
	}, nil
}

func (k msgServer) ExecLegacyContent(goCtx context.Context, msg *govtypesv1.MsgExecLegacyContent) (*govtypesv1.MsgExecLegacyContentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	govAcct := k.GetGovernanceAccount(ctx).GetAddress().String()
	if govAcct != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", govAcct, msg.Authority)
	}

	content, err := govtypesv1.LegacyContentFromMessage(msg)
	if err != nil {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "%+v", err)
	}

	// Ensure that the content has a respective handler
	if !k.Keeper.legacyRouter.HasRoute(content.ProposalRoute()) {
		return nil, sdkerrors.Wrap(govtypes.ErrNoProposalHandlerExists, content.ProposalRoute())
	}

	handler := k.Keeper.legacyRouter.GetRoute(content.ProposalRoute())
	if err := handler(ctx, content); err != nil {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "failed to run legacy handler %s, %+v", content.ProposalRoute(), err)
	}

	return &govtypesv1.MsgExecLegacyContentResponse{}, nil
}

func (k msgServer) Vote(goCtx context.Context, msg *govtypesv1.MsgVote) (*govtypesv1.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, accErr := sdk.AccAddressFromBech32(msg.Voter)
	if accErr != nil {
		return nil, accErr
	}
	options := govtypesv1.NewNonSplitVoteOption(msg.Option)

	err := k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, options, msg.Metadata)
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

	return &govtypesv1.MsgVoteResponse{}, nil
}

func (k msgServer) VoteWeighted(goCtx context.Context, msg *govtypesv1.MsgVoteWeighted) (*govtypesv1.MsgVoteWeightedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, accErr := sdk.AccAddressFromBech32(msg.Voter)
	if accErr != nil {
		return nil, accErr
	}
	err := k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, msg.Options, msg.Metadata)
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

	return &govtypesv1.MsgVoteWeightedResponse{}, nil
}

func (k msgServer) Deposit(goCtx context.Context, msg *govtypesv1.MsgDeposit) (*govtypesv1.MsgDepositResponse, error) {
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

	return &govtypesv1.MsgDepositResponse{}, nil
}

func updateAfterSubmitProposal(ctx sdk.Context, k Keeper, content govtypesv1beta1.Content) error {
	if content.ProposalType() == shieldtypes.ProposalTypeShieldClaim {
		c := content.(*shieldtypes.ShieldClaimProposal)
		lockPeriod := *k.GetVotingParams(ctx).VotingPeriod * 2
		proposerAddr, err := sdk.AccAddressFromBech32(c.Proposer)
		if err != nil {
			return err
		}
		return k.ShieldKeeper.SecureCollaterals(ctx, c.PoolId, proposerAddr, c.PurchaseId, c.Loss, lockPeriod)
	}
	return nil
}

type legacyMsgServer struct {
	govAcct string
	server  govtypesv1.MsgServer
	keeper  Keeper
}

// NewLegacyMsgServerImpl returns an implementation of the v1beta1 legacy MsgServer interface. It wraps around
// the current MsgServer
func NewLegacyMsgServerImpl(govAcct string, v1Server govtypesv1.MsgServer, k Keeper) govtypesv1beta1.MsgServer {
	return &legacyMsgServer{govAcct: govAcct, server: v1Server, keeper: k}
}

var _ govtypesv1beta1.MsgServer = legacyMsgServer{}

func (k legacyMsgServer) SubmitProposal(goCtx context.Context, msg *govtypesv1beta1.MsgSubmitProposal) (*govtypesv1beta1.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := validateProposalByType(ctx, k.keeper, msg)
	if err != nil {
		return nil, err
	}

	contentMsg, err := govtypesv1.NewLegacyContent(msg.GetContent(), k.govAcct)
	if err != nil {
		return nil, fmt.Errorf("error converting legacy content into proposal message: %w", err)
	}

	proposal, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{contentMsg},
		msg.InitialDeposit,
		msg.Proposer,
		"",
	)
	if err != nil {
		return nil, err
	}

	resp, err := k.server.SubmitProposal(goCtx, proposal)
	if err != nil {
		return nil, err
	}

	return &govtypesv1beta1.MsgSubmitProposalResponse{ProposalId: resp.ProposalId}, nil
}

func (k legacyMsgServer) Vote(goCtx context.Context, msg *govtypesv1beta1.MsgVote) (*govtypesv1beta1.MsgVoteResponse, error) {
	_, err := k.server.Vote(goCtx, &govtypesv1.MsgVote{
		ProposalId: msg.ProposalId,
		Voter:      msg.Voter,
		Option:     govtypesv1.VoteOption(msg.Option),
	})
	if err != nil {
		return nil, err
	}
	return &govtypesv1beta1.MsgVoteResponse{}, nil
}

func (k legacyMsgServer) VoteWeighted(goCtx context.Context, msg *govtypesv1beta1.MsgVoteWeighted) (*govtypesv1beta1.MsgVoteWeightedResponse, error) {
	opts := make([]*govtypesv1.WeightedVoteOption, len(msg.Options))
	for idx, opt := range msg.Options {
		opts[idx] = &govtypesv1.WeightedVoteOption{
			Option: govtypesv1.VoteOption(opt.Option),
			Weight: opt.Weight.String(),
		}
	}

	_, err := k.server.VoteWeighted(goCtx, &govtypesv1.MsgVoteWeighted{
		ProposalId: msg.ProposalId,
		Voter:      msg.Voter,
		Options:    opts,
	})
	if err != nil {
		return nil, err
	}
	return &govtypesv1beta1.MsgVoteWeightedResponse{}, nil
}

func (k legacyMsgServer) Deposit(goCtx context.Context, msg *govtypesv1beta1.MsgDeposit) (*govtypesv1beta1.MsgDepositResponse, error) {
	_, err := k.server.Deposit(goCtx, &govtypesv1.MsgDeposit{
		ProposalId: msg.ProposalId,
		Depositor:  msg.Depositor,
		Amount:     msg.Amount,
	})
	if err != nil {
		return nil, err
	}
	return &govtypesv1beta1.MsgDepositResponse{}, nil
}

func validateProposalByType(ctx sdk.Context, k Keeper, msg *govtypesv1beta1.MsgSubmitProposal) error {
	switch c := msg.GetContent().(type) {
	case *certtypes.CertifierUpdateProposal:
		if c.Alias != "" && k.CertKeeper.HasCertifierAlias(ctx, c.Alias) {
			return certtypes.ErrRepeatedAlias
		}

	case *shieldtypes.ShieldClaimProposal:
		// check initial deposit >= max(<loss>*ClaimDepositRate, MinimumClaimDeposit)
		denom := k.stakingKeeper.BondDenom(ctx)

		initialDepositAmount := sdk.NewDecFromInt(msg.InitialDeposit.AmountOf(denom))
		lossAmount := c.Loss.AmountOf(denom)
		lossAmountDec := sdk.NewDecFromInt(lossAmount)
		claimProposalParams := k.ShieldKeeper.GetClaimProposalParams(ctx)
		depositRate := claimProposalParams.DepositRate
		minDeposit := sdk.NewDecFromInt(claimProposalParams.MinDeposit.AmountOf(denom))
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
