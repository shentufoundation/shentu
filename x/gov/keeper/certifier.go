package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// ValidateCertifierUpdateSoloMessage re-exports typesv1.ValidateCertifierUpdateSoloMessage
// so existing keeper callers keep compiling. The canonical definition
// lives in types/v1 so genesis validation (types-level) and the keeper
// share one implementation.
func ValidateCertifierUpdateSoloMessage(messages []sdk.Msg) error {
	return typesv1.ValidateCertifierUpdateSoloMessage(messages)
}

func isCertifierUpdateProposalMsg(msg sdk.Msg) (bool, error) {
	return typesv1.IsCertifierUpdateProposalMsg(msg)
}

// SubmitProposal shadows the embedded cosmos-sdk Keeper.SubmitProposal to
// enforce the cert-update solo-message guard for programmatic callers.
// MsgServer.SubmitProposal also runs this check explicitly.
func (k Keeper) SubmitProposal(ctx context.Context, messages []sdk.Msg, metadata, title, summary string, proposer sdk.AccAddress, expedited bool) (govtypesv1.Proposal, error) {
	if err := ValidateCertifierUpdateSoloMessage(messages); err != nil {
		return govtypesv1.Proposal{}, err
	}
	return k.Keeper.SubmitProposal(ctx, messages, metadata, title, summary, proposer, expedited)
}

// AddCertifierVote add a certifier vote
func (k Keeper) AddCertifierVote(ctx context.Context, proposalID uint64, voterAddr sdk.AccAddress, options govtypesv1.WeightedVoteOptions, metadata string) error {
	isCertifier, err := k.IsCertifier(ctx, voterAddr)
	if err != nil {
		return err
	}
	if !isCertifier {
		return errors.Wrapf(govtypes.ErrInvalidVote, "%s is not a certified identity", voterAddr)
	}

	// Certifier round is a single-ballot head-count model: one certifier
	// contributes exactly one head to exactly one option. Reject weighted
	// ballots so MsgVoteWeighted can't be used to split a certifier's
	// vote into fractional heads that would otherwise skew the tally
	// against the existing yes/no model.
	if len(options) != 1 {
		return errors.Wrap(govtypes.ErrInvalidVote, "certifier votes must specify exactly one option")
	}
	weight, err := math.LegacyNewDecFromStr(options[0].Weight)
	if err != nil {
		return errors.Wrapf(govtypes.ErrInvalidVote, "invalid weight %q", options[0].Weight)
	}
	if !weight.Equal(math.LegacyOneDec()) {
		return errors.Wrap(govtypes.ErrInvalidVote, "certifier votes must have weight 1")
	}
	// Restrict to yes/no. Abstain and NoWithVeto are not part of the
	// certifier head-count model, and allowing them lets a certifier
	// contribute to QuorumMet without expressing a yes/no preference —
	// which would let a single abstain flip an expedited security
	// proposal from the (retry-as-regular) no-quorum path to the
	// (conclusive-reject) quorum-met-but-pass-false path.
	opt := options[0].Option
	if opt != govtypesv1.OptionYes && opt != govtypesv1.OptionNo {
		return errors.Wrap(govtypes.ErrInvalidVote, "certifier votes must be yes or no")
	}

	vote := govtypesv1.NewVote(proposalID, voterAddr, options, metadata)
	err = k.Votes.Set(ctx, collections.Join(proposalID, voterAddr), vote)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCertVote,
			sdk.NewAttribute(govtypes.AttributeKeyOption, options.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(govtypes.AttributeKeyVoter, voterAddr.String()),
		),
	)
	return nil
}

// IsCertifier checks if the input address is a certifier.
func (k Keeper) IsCertifier(ctx context.Context, addr sdk.AccAddress) (bool, error) {
	return k.certKeeper.IsCertifier(ctx, addr)
}

// CertifierVoteIsRequired reports whether a proposal must pass the
// certifier round before it can terminate. Only CertifierUpdate
// proposals need certifier approval; software upgrades and every other
// message type flow through the normal validator stake vote.
//
// A proposal qualifies only when its single message is a cert-update.
// A legacy bundled proposal from v6 ([MsgUpdateCertifier, MsgSend]-style)
// must fall through to the stake round — otherwise the non-cert messages
// would execute on cert-round passage alone, bypassing validator stake.
// v7 submission + genesis already reject bundles via
// ValidateCertifierUpdateSoloMessage; this is the runtime backstop for
// any bundle that survived a v6→v7 upgrade.
func (k Keeper) CertifierVoteIsRequired(ctx context.Context, proposalID uint64) (bool, error) {
	proposal, err := k.Proposals.Get(ctx, proposalID)
	if err != nil {
		return false, err
	}
	proposalMsgs, err := proposal.GetMsgs()
	if err != nil {
		return false, err
	}
	if len(proposalMsgs) != 1 {
		return false, nil
	}
	return isCertifierUpdateProposalMsg(proposalMsgs[0])
}
