package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
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

// SubmitProposal implements the MsgServer.SubmitProposal method.
func (k msgServer) SubmitProposal(goCtx context.Context, msg *govtypesv1.MsgSubmitProposal) (*govtypesv1.MsgSubmitProposalResponse, error) {
	if msg.Title == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proposal title cannot be empty")
	}
	if msg.Summary == "" {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proposal summary cannot be empty")
	}

	proposer, err := k.authKeeper.AddressCodec().StringToBytes(msg.GetProposer())
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid proposer address: %s", err)
	}

	// check that either metadata or Msgs length is non nil.
	if len(msg.Messages) == 0 && len(msg.Metadata) == 0 {
		return nil, errors.Wrap(govtypes.ErrNoProposalMsgs, "either metadata or Msgs length must be non-nil")
	}

	// verify that if present, the metadata title and summary equals the proposal title and summary
	if len(msg.Metadata) != 0 {
		proposalMetadata := govtypes.ProposalMetadata{}
		if err := json.Unmarshal([]byte(msg.Metadata), &proposalMetadata); err == nil {
			if proposalMetadata.Title != msg.Title {
				return nil, errors.Wrapf(govtypes.ErrInvalidProposalContent, "metadata title '%s' must equal proposal title '%s'", proposalMetadata.Title, msg.Title)
			}

			if proposalMetadata.Summary != msg.Summary {
				return nil, errors.Wrapf(govtypes.ErrInvalidProposalContent, "metadata summary '%s' must equal proposal summary '%s'", proposalMetadata.Summary, msg.Summary)
			}
		}

		// if we can't unmarshal the metadata, this means the client didn't use the recommended metadata format
		// nothing can be done here, and this is still a valid case, so we ignore the error
	}

	proposalMsgs, err := msg.GetMsgs()
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	initialDeposit := msg.GetInitialDeposit()

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get governance parameters: %w", err)
	}

	if err := k.validateInitialDeposit(ctx, params, initialDeposit, msg.Expedited); err != nil {
		return nil, err
	}

	if err := k.validateDepositDenom(ctx, params, initialDeposit); err != nil {
		return nil, err
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, proposalMsgs, msg.Metadata, msg.Title, msg.Summary, proposer, msg.Expedited)
	if err != nil {
		return nil, err
	}

	bytes, err := proposal.Marshal()
	if err != nil {
		return nil, err
	}

	// ref: https://github.com/cosmos/cosmos-sdk/issues/9683
	ctx.GasMeter().ConsumeGas(
		3*ctx.KVGasConfig().WriteCostPerByte*uint64(len(bytes)),
		"submit proposal",
	)

	votingStarted, err := k.Keeper.AddDeposit(ctx, proposal.Id, proposer, msg.GetInitialDeposit())
	if err != nil {
		return nil, err
	}

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

// CancelProposal implements the MsgServer.CancelProposal method.
func (k msgServer) CancelProposal(goCtx context.Context, msg *govtypesv1.MsgCancelProposal) (*govtypesv1.MsgCancelProposalResponse, error) {
	_, err := k.authKeeper.AddressCodec().StringToBytes(msg.Proposer)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid proposer address: %s", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.Keeper.CancelProposal(ctx, msg.ProposalId, msg.Proposer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeCancelProposal,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proposer),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprint(msg.ProposalId)),
		),
	)

	return &govtypesv1.MsgCancelProposalResponse{
		ProposalId:     msg.ProposalId,
		CanceledTime:   ctx.BlockTime(),
		CanceledHeight: uint64(ctx.BlockHeight()),
	}, nil
}

// ExecLegacyContent implements the MsgServer.ExecLegacyContent method.
func (k msgServer) ExecLegacyContent(goCtx context.Context, msg *govtypesv1.MsgExecLegacyContent) (*govtypesv1.MsgExecLegacyContentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	govAcct := k.GetGovernanceAccount(ctx).GetAddress().String()
	if govAcct != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", govAcct, msg.Authority)
	}

	content, err := govtypesv1.LegacyContentFromMessage(msg)
	if err != nil {
		return nil, errors.Wrapf(govtypes.ErrInvalidProposalContent, "%+v", err)
	}

	// Ensure that the content has a respective handler
	if !k.Keeper.legacyRouter.HasRoute(content.ProposalRoute()) {
		return nil, errors.Wrap(govtypes.ErrNoProposalHandlerExists, content.ProposalRoute())
	}

	handler := k.Keeper.legacyRouter.GetRoute(content.ProposalRoute())
	if err := handler(ctx, content); err != nil {
		return nil, errors.Wrapf(govtypes.ErrInvalidProposalContent, "failed to run legacy handler %s, %+v", content.ProposalRoute(), err)
	}

	return &govtypesv1.MsgExecLegacyContentResponse{}, nil
}

func (k msgServer) Vote(goCtx context.Context, msg *govtypesv1.MsgVote) (*govtypesv1.MsgVoteResponse, error) {
	accAddr, err := k.authKeeper.AddressCodec().StringToBytes(msg.Voter)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid voter address: %s", err)
	}

	if !govtypesv1.ValidVoteOption(msg.Option) {
		return nil, errors.Wrap(govtypes.ErrInvalidVote, msg.Option.String())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err = k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, govtypesv1.NewNonSplitVoteOption(msg.Option), msg.Metadata)
	if err != nil {
		return nil, err
	}

	return &govtypesv1.MsgVoteResponse{}, nil
}

func (k msgServer) VoteWeighted(goCtx context.Context, msg *govtypesv1.MsgVoteWeighted) (*govtypesv1.MsgVoteWeightedResponse, error) {
	accAddr, accErr := k.authKeeper.AddressCodec().StringToBytes(msg.Voter)
	if accErr != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid voter address: %s", accErr)
	}

	if len(msg.Options) == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, govtypesv1.WeightedVoteOptions(msg.Options).String())
	}

	totalWeight := math.LegacyNewDec(0)
	usedOptions := make(map[govtypesv1.VoteOption]bool)
	for _, option := range msg.Options {
		if !option.IsValid() {
			return nil, errors.Wrap(govtypes.ErrInvalidVote, option.String())
		}
		weight, err := math.LegacyNewDecFromStr(option.Weight)
		if err != nil {
			return nil, errors.Wrapf(govtypes.ErrInvalidVote, "invalid weight: %s", err)
		}
		totalWeight = totalWeight.Add(weight)
		if usedOptions[option.Option] {
			return nil, errors.Wrap(govtypes.ErrInvalidVote, "duplicated vote option")
		}
		usedOptions[option.Option] = true
	}

	if totalWeight.GT(math.LegacyNewDec(1)) {
		return nil, errors.Wrap(govtypes.ErrInvalidVote, "total weight overflow 1.00")
	}

	if totalWeight.LT(math.LegacyNewDec(1)) {
		return nil, errors.Wrap(govtypes.ErrInvalidVote, "total weight lower than 1.00")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, msg.Options, msg.Metadata)
	if err != nil {
		return nil, err
	}

	return &govtypesv1.MsgVoteWeightedResponse{}, nil
}

func (k msgServer) Deposit(goCtx context.Context, msg *govtypesv1.MsgDeposit) (*govtypesv1.MsgDepositResponse, error) {
	accAddr, err := k.authKeeper.AddressCodec().StringToBytes(msg.Depositor)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid depositor address: %s", err)
	}

	if err := validateDeposit(msg.Amount); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	votingStarted, err := k.Keeper.AddDeposit(ctx, msg.ProposalId, accAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

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

func (k msgServer) UpdateParams(goCtx context.Context, msg *govtypesv1.MsgUpdateParams) (*govtypesv1.MsgUpdateParamsResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	if err := msg.Params.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &govtypesv1.MsgUpdateParamsResponse{}, nil
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
	content := msg.GetContent()
	if content == nil {
		return nil, errors.Wrap(govtypes.ErrInvalidProposalContent, "missing content")
	}
	if !govtypesv1beta1.IsValidProposalType(content.ProposalType()) {
		return nil, errors.Wrap(govtypes.ErrInvalidProposalType, content.ProposalType())
	}
	if err := content.ValidateBasic(); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(goCtx)
	err := validateProposalByType(sdkCtx, k.keeper, msg)
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
		msg.GetContent().GetTitle(),
		msg.GetContent().GetDescription(),
		false,
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
		hasAlias, err := k.certKeeper.HasCertifierAlias(ctx, c.Alias)
		if err != nil {
			return err
		}
		if c.Alias != "" && hasAlias {
			return certtypes.ErrRepeatedAlias
		}
	default:
		return nil
	}
	return nil
}

// validateDeposit validates the deposit amount, do not use for initial deposit.
func validateDeposit(amount sdk.Coins) error {
	if !amount.IsValid() || !amount.IsAllPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap(amount.String())
	}

	return nil
}
