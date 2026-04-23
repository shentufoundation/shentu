package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	govkeeper "github.com/shentufoundation/shentu/v2/x/gov/keeper"
)

func TestCertifierVoteIsRequiredForMsgUpdateCertifier(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	governanceAuthority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	// SubmitProposal enforces the certifier-only proposer rule for
	// cert-update proposals; register addrs[0] so submission succeeds.
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(addrs[0], "")))

	proposer, err := app.AccountKeeper.AddressCodec().StringToBytes(addrs[0].String())
	require.NoError(t, err)

	msg := certtypes.NewMsgUpdateCertifier(governanceAuthority, addrs[1], "governance certifier update", certtypes.Add)
	proposal, err := app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{msg}, "", "certifier update", "summary", proposer, false)
	require.NoError(t, err)

	required, err := app.GovKeeper.CertifierVoteIsRequired(ctx, proposal.Id)
	require.NoError(t, err)
	require.True(t, required)
}

// TestCertifierVoteIsRequiredForLegacyCertifierUpdateProposal simulates a
// legacy CertifierUpdateProposal that was submitted before the chain upgrade
// and is still in the store. The tally helper must still recognize it.
func TestCertifierVoteIsRequiredForLegacyCertifierUpdateProposal(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	governanceAuthority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	//nolint:staticcheck // testing deprecated legacy proposal path
	content := &certtypes.CertifierUpdateProposal{
		Title:       "legacy certifier update",
		Description: "summary",
		Certifier:   addrs[1].String(),
		Proposer:    addrs[0].String(),
		AddOrRemove: certtypes.Add,
	}
	legacyMsg, err := govtypesv1.NewLegacyContent(content, governanceAuthority.String())
	require.NoError(t, err)

	msgAny, err := codectypes.NewAnyWithValue(legacyMsg)
	require.NoError(t, err)

	// Directly inject a proposal into the store to simulate a pre-upgrade
	// legacy proposal that the handler route no longer accepts.
	now := time.Now()
	depositEnd := now.Add(48 * time.Hour)
	proposal := govtypesv1.Proposal{
		Id:             999,
		Messages:       []*codectypes.Any{msgAny},
		Status:         govtypesv1.StatusVotingPeriod,
		SubmitTime:     &now,
		DepositEndTime: &depositEnd,
		TotalDeposit:   sdk.NewCoins(),
		Title:          "legacy certifier update",
		Summary:        "summary",
		Proposer:       addrs[0].String(),
	}
	err = app.GovKeeper.SetProposal(ctx, proposal)
	require.NoError(t, err)

	required, err := app.GovKeeper.CertifierVoteIsRequired(ctx, 999)
	require.NoError(t, err)
	require.True(t, required)
}

// Bundling a MsgUpdateCertifier with any other message must be
// rejected at submission. Otherwise the bundled message would ride the
// certifier head-count tally and bypass validator stake voting.
func TestValidateCertifierUpdateSoloMessage_RejectsBundle(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()
	proposer, err := app.AccountKeeper.AddressCodec().StringToBytes(addrs[0].String())
	require.NoError(t, err)

	certMsg := certtypes.NewMsgUpdateCertifier(authority, addrs[1], "add", certtypes.Add)
	sendMsg := banktypes.NewMsgSend(authority, addrs[0], sdk.NewCoins(sdk.NewCoin("uctk", math.NewInt(1))))

	_, err = app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{certMsg, sendMsg}, "", "mixed", "summary", proposer, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "certifier-update message must contain exactly one message")

	// A lone MsgSend must still succeed — only cert-update triggers the guard.
	_, err = app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{sendMsg}, "", "send", "summary", proposer, false)
	require.NoError(t, err)
}

// AddCertifierVote must reject ballots that would break the
// head-count invariant (exactly one yes/no option at weight 1).
// Weighted, abstain, or no-with-veto ballots would otherwise let a
// certifier contribute a partial or non-directional head and skew the
// tally.
func TestAddCertifierVote_RejectsInvalidBallots(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 1, math.NewInt(10000))
	certifier := addrs[0]
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(certifier, "")))

	const proposalID uint64 = 1

	weighted := govtypesv1.WeightedVoteOptions{
		{Option: govtypesv1.OptionYes, Weight: "0.5"},
		{Option: govtypesv1.OptionNo, Weight: "0.5"},
	}
	err := app.GovKeeper.AddCertifierVote(ctx, proposalID, certifier, weighted, "")
	require.Error(t, err)

	abstain := govtypesv1.WeightedVoteOptions{{Option: govtypesv1.OptionAbstain, Weight: "1"}}
	err = app.GovKeeper.AddCertifierVote(ctx, proposalID, certifier, abstain, "")
	require.Error(t, err)

	veto := govtypesv1.WeightedVoteOptions{{Option: govtypesv1.OptionNoWithVeto, Weight: "1"}}
	err = app.GovKeeper.AddCertifierVote(ctx, proposalID, certifier, veto, "")
	require.Error(t, err)

	halfYes := govtypesv1.WeightedVoteOptions{{Option: govtypesv1.OptionYes, Weight: "0.5"}}
	err = app.GovKeeper.AddCertifierVote(ctx, proposalID, certifier, halfYes, "")
	require.Error(t, err)

	// Canonical single-option weight-1 yes vote must succeed.
	yes := govtypesv1.WeightedVoteOptions{{Option: govtypesv1.OptionYes, Weight: "1"}}
	err = app.GovKeeper.AddCertifierVote(ctx, proposalID, certifier, yes, "")
	require.NoError(t, err)
}

// A legacy bundled cert-update proposal (e.g. surviving a v6→v7
// upgrade when submission-time bundle rejection wasn't yet in place)
// must not be treated as cert-only. Classifying such a bundle as
// cert-only would let its non-cert messages execute on certifier-round
// passage alone, bypassing the validator stake tally.
func TestCertifierVoteIsRequired_RejectsBundledProposal(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	certMsg := certtypes.NewMsgUpdateCertifier(authority, addrs[1], "add", certtypes.Add)
	sendMsg := banktypes.NewMsgSend(authority, addrs[0], sdk.NewCoins(sdk.NewCoin("uctk", math.NewInt(1))))

	// Construct and persist the bundle without going through SubmitProposal,
	// which would reject it at the solo-message guard.
	now := time.Now()
	end := now.Add(48 * time.Hour)
	bundle, err := govtypesv1.NewProposal(
		[]sdk.Msg{certMsg, sendMsg},
		55555,
		now, end,
		"legacy bundle", "title", "summary", addrs[0], false,
	)
	require.NoError(t, err)
	bundle.Status = govtypesv1.StatusVotingPeriod
	bundle.VotingStartTime = &now
	bundle.VotingEndTime = &end
	require.NoError(t, app.GovKeeper.SetProposal(ctx, bundle))

	required, err := app.GovKeeper.CertifierVoteIsRequired(ctx, bundle.Id)
	require.NoError(t, err)
	require.False(t, required, "bundled cert-update must fall through to stake round")
}

// The exported ValidateCertifierUpdateSoloMessage must be callable by
// the genesis import path. This exercises the predicate without the
// SubmitProposal wrapper.
func TestValidateCertifierUpdateSoloMessage_Predicate(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	certMsg := certtypes.NewMsgUpdateCertifier(authority, addrs[1], "add", certtypes.Add)
	sendMsg := banktypes.NewMsgSend(authority, addrs[0], sdk.NewCoins(sdk.NewCoin("uctk", math.NewInt(1))))

	require.NoError(t, govkeeper.ValidateCertifierUpdateSoloMessage([]sdk.Msg{certMsg}))
	require.NoError(t, govkeeper.ValidateCertifierUpdateSoloMessage([]sdk.Msg{sendMsg}))
	require.Error(t, govkeeper.ValidateCertifierUpdateSoloMessage([]sdk.Msg{certMsg, sendMsg}))
}
