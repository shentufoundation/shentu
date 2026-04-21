package v6_test

import (
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	govkeeper "github.com/shentufoundation/shentu/v2/x/gov/keeper"
	v6 "github.com/shentufoundation/shentu/v2/x/gov/migrations/v6"
)

func proposalIDBytes(id uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return b
}

func writeLegacyCertVotedEntry(t *testing.T, ctx sdk.Context, app *shentuapp.ShentuApp, proposalID uint64) {
	t.Helper()
	kv := app.GetKey(govtypes.StoreKey)
	store := ctx.KVStore(kv)
	certVotes := prefix.NewStore(store, v6.CertVotesKeyPrefix)
	certVotes.Set(proposalIDBytes(proposalID), proposalIDBytes(proposalID))
}

func certVotedHas(ctx sdk.Context, app *shentuapp.ShentuApp, proposalID uint64) bool {
	kv := app.GetKey(govtypes.StoreKey)
	store := ctx.KVStore(kv)
	certVotes := prefix.NewStore(store, v6.CertVotesKeyPrefix)
	return certVotes.Has(proposalIDBytes(proposalID))
}

// requirePanicContains recovers from a panic and asserts the panic
// message contains each of the given substrings. The migration's panic
// message concatenates a list of blockers with newlines; substring
// checks are more stable than comparing the full string.
func requirePanicContains(t *testing.T, fn func(), substrings ...string) {
	t.Helper()
	defer func() {
		r := recover()
		require.NotNil(t, r, "expected panic")
		msg := fmt.Sprint(r)
		for _, s := range substrings {
			require.Contains(t, msg, s)
		}
	}()
	fn()
}

// submitActiveCertUpdate submits a solo cert-update proposal and
// activates its voting period. Returns the proposal.
func submitActiveCertUpdate(t *testing.T, app *shentuapp.ShentuApp, ctx sdk.Context, proposer sdk.AccAddress, certifier sdk.AccAddress) govtypesv1.Proposal {
	t.Helper()
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()
	proposerBytes, err := app.AccountKeeper.AddressCodec().StringToBytes(proposer.String())
	require.NoError(t, err)
	msg := certtypes.NewMsgUpdateCertifier(authority, certifier, "add", certtypes.Add)
	proposal, err := app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{msg}, "", "cert update", "summary", proposerBytes, false)
	require.NoError(t, err)
	require.NoError(t, app.GovKeeper.ActivateVotingPeriod(ctx, proposal))
	return proposal
}

// The migration must remove every cert_voted entry when the matching
// proposal is no longer in flight. Orphan entries (no matching
// proposal, or one already out of VotingPeriod) are the common case
// after the binary switches off the two-round model.
func TestMigrate6to7_DeletesOrphanedEntries(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	writeLegacyCertVotedEntry(t, ctx, app, 101)
	writeLegacyCertVotedEntry(t, ctx, app, 202)
	require.True(t, certVotedHas(ctx, app, 101))
	require.True(t, certVotedHas(ctx, app, 202))

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	require.NoError(t, m.Migrate6to7(ctx))

	require.False(t, certVotedHas(ctx, app, 101))
	require.False(t, certVotedHas(ctx, app, 202))
}

// The migration must panic when a cert_voted entry still belongs to an
// active CertifierUpdate proposal: running the new code against such a
// proposal would re-tally its stake votes as certifier head-counts and
// silently corrupt the outcome.
func TestMigrate6to7_PanicsOnInFlightProposal(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false).WithBlockTime(time.Unix(100, 0))

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	proposal := submitActiveCertUpdate(t, app, ctx, addrs[0], addrs[1])
	writeLegacyCertVotedEntry(t, ctx, app, proposal.Id)

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	requirePanicContains(t, func() {
		_ = m.Migrate6to7(ctx)
	},
		"gov v6→v7 migration blocked",
		fmt.Sprintf("proposal %d", proposal.Id),
		"cert_voted=true",
	)

	// The entry must still be present — the migration aborts before deletion.
	require.True(t, certVotedHas(ctx, app, proposal.Id))
}

// A cert_voted entry for a proposal whose status has moved past
// VotingPeriod (passed/rejected/failed) is not considered in-flight and
// must be swept.
func TestMigrate6to7_IgnoresCompletedProposals(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false).WithBlockTime(time.Unix(100, 0))

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()
	proposer, err := app.AccountKeeper.AddressCodec().StringToBytes(addrs[0].String())
	require.NoError(t, err)

	certMsg := certtypes.NewMsgUpdateCertifier(authority, addrs[1], "add", certtypes.Add)
	proposal, err := app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{certMsg}, "", "cert update", "summary", proposer, false)
	require.NoError(t, err)

	proposal.Status = govtypesv1.StatusRejected
	require.NoError(t, app.GovKeeper.SetProposal(ctx, proposal))

	writeLegacyCertVotedEntry(t, ctx, app, proposal.Id)

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	require.NoError(t, m.Migrate6to7(ctx))
	require.False(t, certVotedHas(ctx, app, proposal.Id))
}

// A bundled cert-update proposal (one of v6's legal shapes) must block
// the upgrade. Under v7, CertifierVoteIsRequired returns false for
// bundles and they fall through to the stake round — but that's still
// an operator-visible behavior change, so make them surface at
// migration time rather than silently rerouting.
func TestMigrate6to7_PanicsOnLegacyBundledProposal(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false).WithBlockTime(time.Unix(100, 0))

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	certMsg := certtypes.NewMsgUpdateCertifier(authority, addrs[1], "add", certtypes.Add)
	sendMsg := banktypes.NewMsgSend(authority, addrs[0], sdk.NewCoins(sdk.NewCoin("uctk", math.NewInt(1))))

	// v7 rejects bundles at submission; directly inject the proposal to
	// simulate a v6-submitted bundle that survived the upgrade.
	now := time.Unix(100, 0)
	votingEnd := now.Add(48 * time.Hour)
	bundle, err := govtypesv1.NewProposal(
		[]sdk.Msg{certMsg, sendMsg},
		12345,
		now, votingEnd,
		"legacy bundle", "title", "summary", addrs[0], false,
	)
	require.NoError(t, err)
	bundle.Status = govtypesv1.StatusVotingPeriod
	bundle.VotingStartTime = &now
	bundle.VotingEndTime = &votingEnd
	require.NoError(t, app.GovKeeper.SetProposal(ctx, bundle))

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	requirePanicContains(t, func() {
		_ = m.Migrate6to7(ctx)
	},
		"bundles cert-update",
		fmt.Sprintf("proposal %d", bundle.Id),
	)
}

// A bundled cert-update proposal that is still in the deposit period
// at upgrade time must also block. Skipping it would let
// ActivateVotingPeriod admit it post-upgrade; v7's
// CertifierVoteIsRequired returns false for bundles, so the proposal
// would silently become a plain stake-only vote with no certifier
// round at all — and on stake pass, MsgUpdateCertifier would execute
// without any certifier approval.
func TestMigrate6to7_PanicsOnLegacyBundledDepositPeriodProposal(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false).WithBlockTime(time.Unix(100, 0))

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	certMsg := certtypes.NewMsgUpdateCertifier(authority, addrs[1], "add", certtypes.Add)
	sendMsg := banktypes.NewMsgSend(authority, addrs[0], sdk.NewCoins(sdk.NewCoin("uctk", math.NewInt(1))))

	now := time.Unix(100, 0)
	depositEnd := now.Add(48 * time.Hour)
	bundle, err := govtypesv1.NewProposal(
		[]sdk.Msg{certMsg, sendMsg},
		54321,
		now, depositEnd,
		"legacy bundle in deposit", "title", "summary", addrs[0], false,
	)
	require.NoError(t, err)
	bundle.Status = govtypesv1.StatusDepositPeriod
	bundle.DepositEndTime = &depositEnd
	require.NoError(t, app.GovKeeper.SetProposal(ctx, bundle))

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	requirePanicContains(t, func() {
		_ = m.Migrate6to7(ctx)
	},
		"bundles cert-update",
		"deposit period",
		fmt.Sprintf("proposal %d", bundle.Id),
	)
}

// A legacy weighted certifier ballot stored under v6 must block the
// upgrade. v6's tally counted any single-option weighted ballot as 1
// head; v7's tally drops anything not weight=1, so leaving such a
// ballot in place would silently flip the outcome.
func TestMigrate6to7_PanicsOnLegacyWeightedCertifierVote(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false).WithBlockTime(time.Unix(100, 0))

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	certifier := addrs[0]
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(certifier, "")))

	proposal := submitActiveCertUpdate(t, app, ctx, addrs[0], addrs[1])

	// Inject a weighted ballot directly into the Votes collection to
	// simulate a v6 MsgVoteWeighted that AddCertifierVote now rejects.
	weighted := govtypesv1.NewVote(proposal.Id, certifier, govtypesv1.WeightedVoteOptions{
		{Option: govtypesv1.OptionYes, Weight: "0.5"},
	}, "")
	require.NoError(t, app.GovKeeper.Votes.Set(ctx, collections.Join(proposal.Id, certifier), weighted))

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	requirePanicContains(t, func() {
		_ = m.Migrate6to7(ctx)
	},
		"weighted ballot",
		`weight="0.5"`,
		fmt.Sprintf("proposal %d", proposal.Id),
	)
}

// A legacy multi-option certifier ballot (also formerly possible via
// MsgVoteWeighted) must block the upgrade for the same reason: v7's
// SecurityTally skips any ballot that isn't single-option.
func TestMigrate6to7_PanicsOnLegacyMultiOptionCertifierVote(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false).WithBlockTime(time.Unix(100, 0))

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	certifier := addrs[0]
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(certifier, "")))

	proposal := submitActiveCertUpdate(t, app, ctx, addrs[0], addrs[1])

	multi := govtypesv1.NewVote(proposal.Id, certifier, govtypesv1.WeightedVoteOptions{
		{Option: govtypesv1.OptionYes, Weight: "0.5"},
		{Option: govtypesv1.OptionNo, Weight: "0.5"},
	}, "")
	require.NoError(t, app.GovKeeper.Votes.Set(ctx, collections.Join(proposal.Id, certifier), multi))

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	requirePanicContains(t, func() {
		_ = m.Migrate6to7(ctx)
	},
		"2-option ballot",
		fmt.Sprintf("proposal %d", proposal.Id),
	)
}
