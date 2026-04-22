package gov_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	govmodule "github.com/shentufoundation/shentu/v2/x/gov"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// CertifierUpdate proposals are decided entirely by the certifier head-count
// tally: with no certifier votes, expiry must reject the proposal and burn
// its deposit.
func TestEndBlocker_CertifierUpdateRejectedOnNoCertifierVotes(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false).WithBlockTime(time.Unix(100, 0))

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	govAuthority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()
	// SubmitProposal rejects cert-update from non-certifiers; register
	// the proposer so submission reaches the end-blocker path under test.
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(addrs[0], "")))
	proposer, err := app.AccountKeeper.AddressCodec().StringToBytes(addrs[0].String())
	require.NoError(t, err)

	msg := certtypes.NewMsgUpdateCertifier(govAuthority, addrs[1], "add certifier", certtypes.Add)
	proposal, err := app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{msg}, "", "certifier update", "summary", proposer, false)
	require.NoError(t, err)
	require.NoError(t, app.GovKeeper.ActivateVotingPeriod(ctx, proposal))

	proposal, err = app.GovKeeper.Proposals.Get(ctx, proposal.Id)
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(proposal.VotingEndTime.Add(time.Second))
	require.NoError(t, govmodule.EndBlocker(ctx, app.GovKeeper))

	proposal, err = app.GovKeeper.Proposals.Get(ctx, proposal.Id)
	require.NoError(t, err)
	require.Equal(t, govtypesv1.StatusRejected, proposal.Status)
	require.Equal(t, "proposal did not pass the certifier voting period", proposal.FailedReason)

	has, err := app.GovKeeper.ActiveProposalsQueue.Has(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id))
	require.NoError(t, err)
	require.False(t, has)
}

// Export after init must produce a genesis that re-imports cleanly and
// round-trips to the same logical state. A failure here usually means a
// field that used to be persisted (e.g. CertVotedProposalIds) is still
// being read or written somewhere.
func TestGenesisRoundTrip(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	exported, err := govmodule.ExportGenesis(ctx, app.GovKeeper)
	require.NoError(t, err)
	require.NotNil(t, exported)

	// Re-import into a fresh app and export again; the two exports must
	// be byte-for-byte equivalent through the wire format.
	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	govmodule.InitGenesis(ctx2, app2.GovKeeper, app2.AccountKeeper, app2.BankKeeper, exported)

	exported2, err := govmodule.ExportGenesis(ctx2, app2.GovKeeper)
	require.NoError(t, err)

	bz1, err := app.AppCodec().MarshalJSON(exported)
	require.NoError(t, err)
	bz2, err := app2.AppCodec().MarshalJSON(exported2)
	require.NoError(t, err)
	require.JSONEq(t, string(bz1), string(bz2))
}

// InitGenesis must refuse a bundled cert-update proposal. Persisting
// one via genesis would let it ride the certifier head-count tally on
// the first block and bypass the validator stake tally — the exact
// exploit ValidateCertifierUpdateSoloMessage exists to block at submit
// time.
func TestInitGenesis_RejectsBundledCertUpdateProposal(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	certMsg := certtypes.NewMsgUpdateCertifier(authority, addrs[1], "add", certtypes.Add)
	sendMsg := banktypes.NewMsgSend(authority, addrs[0], sdk.NewCoins(sdk.NewCoin("uctk", math.NewInt(1))))

	exported, err := govmodule.ExportGenesis(ctx, app.GovKeeper)
	require.NoError(t, err)

	bundle, err := govtypesv1.NewProposal(
		[]sdk.Msg{certMsg, sendMsg},
		42,
		time.Unix(1, 0), time.Unix(2, 0),
		"bundled cert update", "title", "summary", addrs[0], false,
	)
	require.NoError(t, err)
	bundle.Status = govtypesv1.StatusDepositPeriod
	exported.Proposals = append(exported.Proposals, &bundle)

	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	require.Panics(t, func() {
		govmodule.InitGenesis(ctx2, app2.GovKeeper, app2.AccountKeeper, app2.BankKeeper, exported)
	})
}

// Default genesis must no longer carry a StakeVoteTally — that field
// backed the retired validator stake round for cert-update proposals.
// The default must contain only the certifier security-round tally.
func TestDefaultGenesis_NoStakeTally(t *testing.T) {
	gs := typesv1.DefaultGenesisState()
	require.NotNil(t, gs.CustomParams)
	require.NotNil(t, gs.CustomParams.CertifierUpdateSecurityVoteTally)
}

// SoftwareUpgrade no longer passes through the certifier round; it flows
// through the normal validator stake vote. CertifierVoteIsRequired must
// report false for it, so the cert-round code path is never triggered.
func TestSoftwareUpgradeSkipsCertifierRound(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false).WithBlockTime(time.Unix(100, 0))

	_, _, proposer := testdata.KeyTestPubAddr()
	govAuthority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress().String()

	upgradeMsg := &upgradetypes.MsgSoftwareUpgrade{
		Authority: govAuthority,
		Plan: upgradetypes.Plan{
			Name:   "v2.18.0",
			Height: 200,
		},
	}

	proposal, err := app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{upgradeMsg}, "", "upgrade", "summary", proposer, false)
	require.NoError(t, err)

	required, err := app.GovKeeper.CertifierVoteIsRequired(ctx, proposal.Id)
	require.NoError(t, err)
	require.False(t, required)
}
