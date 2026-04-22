package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// securityTallyFixture wires up N registered certifiers, a pending
// MsgUpdateCertifier proposal, and returns the proposal so tests can
// write single-option weight-1 ballots directly to the Votes
// collection (bypassing AddVote, which reroutes via CertifierVoteIsRequired).
func securityTallyFixture(t *testing.T, nCertifiers int) (
	*shentuapp.ShentuApp, sdk.Context, govtypesv1.Proposal, []sdk.AccAddress,
) {
	t.Helper()
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, nCertifiers+1, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()
	proposerAddr := addrs[0]
	certifiers := addrs[1:]

	// Register the proposer as a certifier too — SubmitProposal enforces
	// the certifier-only proposer rule for cert-update proposals. The
	// fixture deliberately does NOT include the proposer in the voting
	// electorate (`certifiers` slice) so tally math stays unaffected.
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(proposerAddr, "")))
	for _, c := range certifiers {
		require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(c, "")))
	}

	proposer, err := app.AccountKeeper.AddressCodec().StringToBytes(proposerAddr.String())
	require.NoError(t, err)

	msg := certtypes.NewMsgUpdateCertifier(authority, addrs[0], "add", certtypes.Add)
	proposal, err := app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{msg}, "", "cert update", "summary", proposer, false)
	require.NoError(t, err)

	return app, ctx, proposal, certifiers
}

// setVote writes a ballot directly to the Votes collection, bypassing
// AddVote's routing so the SecurityTally caller can test the raw
// tally math without the certifier-role guard.
func setVote(
	t *testing.T,
	app *shentuapp.ShentuApp,
	ctx sdk.Context,
	proposalID uint64,
	voter sdk.AccAddress,
	opt govtypesv1.VoteOption,
	weight string,
) {
	t.Helper()
	vote := govtypesv1.NewVote(proposalID, voter, govtypesv1.WeightedVoteOptions{
		{Option: opt, Weight: weight},
	}, "")
	require.NoError(t, app.GovKeeper.Votes.Set(ctx, collections.Join(proposalID, voter), vote))
}

// Three certifiers all vote yes → quorum met, yes > threshold (66.7%),
// so SecurityTally passes. Confirms the happy path with default
// certifier tally params from genesis.
func TestSecurityTally_Pass(t *testing.T) {
	app, ctx, proposal, certifiers := securityTallyFixture(t, 3)

	for _, c := range certifiers {
		setVote(t, app, ctx, proposal.Id, c, govtypesv1.OptionYes, "1")
	}

	result, err := app.GovKeeper.SecurityTally(ctx, proposal)
	require.NoError(t, err)
	require.True(t, result.Pass, "3/3 yes must pass")
	require.Equal(t, math.NewInt(3).String(), result.Tally.YesCount)
	require.Equal(t, math.NewInt(0).String(), result.Tally.NoCount)
}

// Three certifiers vote 2 no / 1 yes → quorum met but yes below the
// 66.7% threshold, so the proposal fails. This is the "quorum met,
// threshold missed" branch that's easy to regress if passAndVetoSecurityResult
// is reworked.
func TestSecurityTally_RejectBelowThreshold(t *testing.T) {
	app, ctx, proposal, certifiers := securityTallyFixture(t, 3)

	setVote(t, app, ctx, proposal.Id, certifiers[0], govtypesv1.OptionYes, "1")
	setVote(t, app, ctx, proposal.Id, certifiers[1], govtypesv1.OptionNo, "1")
	setVote(t, app, ctx, proposal.Id, certifiers[2], govtypesv1.OptionNo, "1")

	result, err := app.GovKeeper.SecurityTally(ctx, proposal)
	require.NoError(t, err)
	require.False(t, result.Pass, "1/3 yes is below 66.7% threshold")
	require.Equal(t, math.NewInt(1).String(), result.Tally.YesCount)
	require.Equal(t, math.NewInt(2).String(), result.Tally.NoCount)
}

// One certifier out of three votes → participation 1/3 = 33.3%,
// below the default 33.4% quorum. SecurityTally must report Pass=false
// without consulting the threshold.
func TestSecurityTally_NoQuorum(t *testing.T) {
	app, ctx, proposal, certifiers := securityTallyFixture(t, 3)

	setVote(t, app, ctx, proposal.Id, certifiers[0], govtypesv1.OptionYes, "1")

	result, err := app.GovKeeper.SecurityTally(ctx, proposal)
	require.NoError(t, err)
	require.False(t, result.Pass, "below quorum must not pass")
}

// Clearing CertifierUpdateSecurityVoteTally from CustomParams leaves
// SecurityTally without the knobs it needs. It must surface the error
// to the caller (who treats it as a rejection) instead of silently
// defaulting.
func TestSecurityTally_MissingCustomParamsReturnsError(t *testing.T) {
	app, ctx, proposal, certifiers := securityTallyFixture(t, 1)

	setVote(t, app, ctx, proposal.Id, certifiers[0], govtypesv1.OptionYes, "1")

	require.NoError(t, app.GovKeeper.SetCustomParams(ctx, typesv1.CustomParams{
		CertifierUpdateSecurityVoteTally: nil,
	}))

	_, err := app.GovKeeper.SecurityTally(ctx, proposal)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CertifierUpdateSecurityVoteTally")
}

// AddCertifierVote rejects weighted/multi-option ballots at
// submission, but the store might already hold such ballots from a
// pre-guard binary. SecurityTally's defensive filter must skip them
// silently — neither crash nor count them.
func TestSecurityTally_SkipsLegacyWeightedAndMultiOptionBallots(t *testing.T) {
	app, ctx, proposal, certifiers := securityTallyFixture(t, 4)

	// One valid yes — it's the only ballot that should count.
	setVote(t, app, ctx, proposal.Id, certifiers[0], govtypesv1.OptionYes, "1")

	// Weight < 1 — must be skipped.
	setVote(t, app, ctx, proposal.Id, certifiers[1], govtypesv1.OptionYes, "0.5")

	// Multi-option ballot written directly — must be skipped.
	multi := govtypesv1.NewVote(proposal.Id, certifiers[2], govtypesv1.WeightedVoteOptions{
		{Option: govtypesv1.OptionYes, Weight: "0.5"},
		{Option: govtypesv1.OptionNo, Weight: "0.5"},
	}, "")
	require.NoError(t, app.GovKeeper.Votes.Set(ctx, collections.Join(proposal.Id, certifiers[2]), multi))

	// Unparseable weight — must be skipped, not panic.
	setVote(t, app, ctx, proposal.Id, certifiers[3], govtypesv1.OptionYes, "not-a-decimal")

	result, err := app.GovKeeper.SecurityTally(ctx, proposal)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(1).String(), result.Tally.YesCount,
		"only the weight-1 single-option ballot should count")
	require.Equal(t, math.NewInt(0).String(), result.Tally.NoCount)
}
