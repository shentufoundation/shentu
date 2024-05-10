package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	v046 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v046"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

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

	// todo remove?
	//err := validateProposalByType(ctx, k.Keeper, msg)
	//if err != nil {
	//	return nil, err
	//}

	proposalMsgs, err := msg.GetMsgs()
	if err != nil {
		return nil, err
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, proposalMsgs, msg.Metadata)
	if err != nil {
		return nil, err
	}

	// Skip deposit period for proposals from certifier memebers or shield claim proposals.
	proposer, _ := sdk.AccAddressFromBech32(msg.GetProposer())
	votingStarted := k.ActivateVotingPeriodCustom(ctx, proposal, proposer)
	if !votingStarted {
		votingStarted, err = k.Keeper.AddDeposit(ctx, proposal.Id, proposer, msg.GetInitialDeposit())
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
	//TODO implement me
	panic("implement me")
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

// todo remove to msg check
func validateProposalByType(ctx sdk.Context, k Keeper, msg *v1beta1.MsgSubmitProposal) error {
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

func updateAfterSubmitProposal(ctx sdk.Context, k Keeper, proposal govtypesv1.Proposal) error {
	legacyProposal, err := v046.ConvertToLegacyProposal(proposal)
	if err != nil {
		return err
	}
	if legacyProposal.ProposalType() == shieldtypes.ProposalTypeShieldClaim {
		c := legacyProposal.GetContent().(*shieldtypes.ShieldClaimProposal)
		lockPeriod := time.Duration(*k.GetVotingParams(ctx).VotingPeriod) * 2
		proposerAddr, err := sdk.AccAddressFromBech32(c.Proposer)
		if err != nil {
			return err
		}
		return k.ShieldKeeper.SecureCollaterals(ctx, c.PoolId, proposerAddr, c.PurchaseId, c.Loss, lockPeriod)
	}
	return nil
}
