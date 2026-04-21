package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/common"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
)

// msgServerFixture wires a MsgServer, a funded proposer (uctk >=
// MinDeposit * MinDepositRatio), and the gov module authority. Every
// msg-server SubmitProposal test needs these three handles plus the
// deposit amount to submit alongside the proposal.
func msgServerFixture(t *testing.T) (
	*shentuapp.ShentuApp, sdk.Context, govtypesv1.MsgServer,
	sdk.AccAddress, sdk.AccAddress, sdk.Coins,
) {
	t.Helper()
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	proposer := addrs[0]

	// Fund the proposer with enough uctk to satisfy MinDepositRatio
	// (default 0.01 of MinDeposit). MintCoins + SendCoinsFromModuleToAccount
	// is the standard shentu test pattern for uctk balances.
	deposit := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, math.NewInt(1_000_000_000)))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, deposit))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, proposer, deposit))

	msgSrvr := keeper.NewMsgServerImpl(app.GovKeeper)
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()
	return app, ctx, msgSrvr, authority, proposer, deposit
}

// MsgServer.SubmitProposal must reject a proposal that bundles a
// MsgUpdateCertifier with any other message. This is the front-line
// enforcement point for external clients — if it's ever bypassed, a
// bundle would reach the store and let its non-cert messages ride the
// certifier head-count tally.
func TestMsgServerSubmitProposal_RejectsBundledCertUpdate(t *testing.T) {
	app, ctx, msgSrvr, authority, proposer, deposit := msgServerFixture(t)

	certMsg := certtypes.NewMsgUpdateCertifier(authority, proposer, "add", certtypes.Add)
	sendMsg := banktypes.NewMsgSend(
		authority, proposer,
		sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, math.NewInt(1))),
	)

	msg, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{certMsg, sendMsg},
		deposit,
		proposer.String(),
		"",
		"bundle",
		"summary",
		false,
	)
	require.NoError(t, err)

	_, err = msgSrvr.SubmitProposal(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "certifier-update message must contain exactly one message")

	// No proposal may have been persisted.
	count, err := proposalCount(app, ctx)
	require.NoError(t, err)
	require.Zero(t, count)
}

// A solo MsgUpdateCertifier from a registered certifier must go through
// the msg-server entry point cleanly — the solo-message guard only trips
// on bundles, and the certifier-only proposer guard is satisfied.
func TestMsgServerSubmitProposal_AcceptsSoloCertUpdate(t *testing.T) {
	app, ctx, msgSrvr, authority, proposer, deposit := msgServerFixture(t)

	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(proposer, "")))

	certMsg := certtypes.NewMsgUpdateCertifier(authority, proposer, "add", certtypes.Add)

	msg, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{certMsg},
		deposit,
		proposer.String(),
		"",
		"solo cert",
		"summary",
		false,
	)
	require.NoError(t, err)

	resp, err := msgSrvr.SubmitProposal(ctx, msg)
	require.NoError(t, err)
	require.NotZero(t, resp.ProposalId)
}

// A plain MsgSend proposal must not be affected by the cert-update
// guard — it should reach the embedded cosmos-sdk keeper and persist
// normally. Guards against over-eager rejection logic.
func TestMsgServerSubmitProposal_AcceptsPlainProposal(t *testing.T) {
	_, ctx, msgSrvr, authority, proposer, deposit := msgServerFixture(t)

	sendMsg := banktypes.NewMsgSend(
		authority, proposer,
		sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, math.NewInt(1))),
	)

	msg, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{sendMsg},
		deposit,
		proposer.String(),
		"",
		"send",
		"summary",
		false,
	)
	require.NoError(t, err)

	resp, err := msgSrvr.SubmitProposal(ctx, msg)
	require.NoError(t, err)
	require.NotZero(t, resp.ProposalId)
}

// A cert-update proposal from a non-certifier must be rejected with
// ErrInvalidProposer. Without this guard, any address with a
// deposit-sized balance could flood the cert round with spam proposals.
func TestMsgServerSubmitProposal_RejectsCertUpdateFromNonCertifier(t *testing.T) {
	app, ctx, msgSrvr, authority, proposer, deposit := msgServerFixture(t)

	// proposer is NOT registered as a certifier.
	certMsg := certtypes.NewMsgUpdateCertifier(authority, proposer, "add", certtypes.Add)

	msg, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{certMsg},
		deposit,
		proposer.String(),
		"",
		"solo cert",
		"summary",
		false,
	)
	require.NoError(t, err)

	_, err = msgSrvr.SubmitProposal(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, govtypes.ErrInvalidProposer)

	// No proposal may have been persisted.
	count, err := proposalCount(app, ctx)
	require.NoError(t, err)
	require.Zero(t, count)
}

// Regression: empty title must still be rejected. Makes sure the
// solo-message guard wasn't reordered in front of basic validation.
func TestMsgServerSubmitProposal_RejectsEmptyTitle(t *testing.T) {
	_, ctx, msgSrvr, authority, proposer, deposit := msgServerFixture(t)

	sendMsg := banktypes.NewMsgSend(
		authority, proposer,
		sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, math.NewInt(1))),
	)

	msg, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{sendMsg},
		deposit,
		proposer.String(),
		"",
		"",
		"summary",
		false,
	)
	require.NoError(t, err)

	_, err = msgSrvr.SubmitProposal(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "title")
}

func proposalCount(app *shentuapp.ShentuApp, ctx sdk.Context) (int, error) {
	n := 0
	err := app.GovKeeper.Proposals.Walk(ctx, nil, func(_ uint64, _ govtypesv1.Proposal) (bool, error) {
		n++
		return false, nil
	})
	return n, err
}
