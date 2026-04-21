package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/common"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
)

// submitAndActivate submits a proposal and flips it into VotingPeriod
// so AddVote's "inVotingPeriod" guard passes. It returns the proposal ID.
func submitAndActivate(t *testing.T, app *shentuapp.ShentuApp, ctx sdk.Context, msgs []sdk.Msg, proposer sdk.AccAddress) uint64 {
	t.Helper()
	proposerBz, err := app.AccountKeeper.AddressCodec().StringToBytes(proposer.String())
	require.NoError(t, err)
	proposal, err := app.GovKeeper.SubmitProposal(ctx, msgs, "", "title", "summary", proposerBz, false)
	require.NoError(t, err)
	require.NoError(t, app.GovKeeper.ActivateVotingPeriod(ctx, proposal))
	return proposal.Id
}

// AddVote on a solo cert-update proposal must route to AddCertifierVote,
// which rejects a non-certifier voter. If the router ever misclassified
// cert-update as a stake proposal, a non-certifier's ballot would be
// stored as a stake vote — that's the exact silent failure this test
// guards against.
func TestAddVote_CertUpdateRoutedToCertifierVote_NonCertifierRejected(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	proposer := addrs[0]
	nonCertifier := addrs[1]
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	// Cert-update submission requires a certifier proposer.
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(proposer, "")))

	certMsg := certtypes.NewMsgUpdateCertifier(authority, addrs[0], "add", certtypes.Add)
	proposalID := submitAndActivate(t, app, ctx, []sdk.Msg{certMsg}, proposer)

	err := app.GovKeeper.AddVote(ctx, proposalID, nonCertifier, govtypesv1.WeightedVoteOptions{
		{Option: govtypesv1.OptionYes, Weight: "1"},
	}, "")
	require.Error(t, err, "non-certifier must be rejected by the cert-vote path")
	require.ErrorIs(t, err, govtypes.ErrInvalidVote)

	// Verify no stake-style ballot was persisted by the fallback.
	has, err := app.GovKeeper.Votes.Has(ctx, collections.Join(proposalID, nonCertifier))
	require.NoError(t, err)
	require.False(t, has)
}

// AddVote on a solo cert-update proposal, from a registered certifier,
// must succeed and persist exactly one single-option weight-1 ballot
// (the invariant SecurityTally relies on). Without correct routing the
// vote would land as a multi-weighted stake ballot and silently
// corrupt the certifier head-count.
func TestAddVote_CertUpdateRoutedToCertifierVote_CertifierAccepted(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	proposer := addrs[0]
	certifier := addrs[1]
	// Proposer must also be a certifier so cert-update submission passes.
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(proposer, "")))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(certifier, "")))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	certMsg := certtypes.NewMsgUpdateCertifier(authority, addrs[0], "add", certtypes.Add)
	proposalID := submitAndActivate(t, app, ctx, []sdk.Msg{certMsg}, proposer)

	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, certifier, govtypesv1.WeightedVoteOptions{
		{Option: govtypesv1.OptionYes, Weight: "1"},
	}, ""))

	vote, err := app.GovKeeper.Votes.Get(ctx, collections.Join(proposalID, certifier))
	require.NoError(t, err)
	require.Len(t, vote.Options, 1, "certifier ballot must be single-option")
	weight, err := math.LegacyNewDecFromStr(vote.Options[0].Weight)
	require.NoError(t, err)
	require.True(t, weight.Equal(math.LegacyOneDec()), "certifier ballot must be weight 1")
	require.Equal(t, govtypesv1.OptionYes, vote.Options[0].Option)
}

// AddVote on a non-cert proposal (MsgSend) must land on the stake-vote
// path regardless of the voter's certifier status. Routing a stake
// proposal to the certifier path would reject every non-certifier
// voter — the entire validator electorate — breaking governance.
func TestAddVote_NonCertProposalRoutedToStakeVote(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	proposer := addrs[0]
	voter := addrs[1]
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	sendMsg := banktypes.NewMsgSend(authority, proposer, sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, math.NewInt(1))))
	proposalID := submitAndActivate(t, app, ctx, []sdk.Msg{sendMsg}, proposer)

	// Voter is NOT a certifier. On the stake path this must still succeed —
	// the cert-voter guard must not apply.
	options := govtypesv1.WeightedVoteOptions{
		{Option: govtypesv1.OptionYes, Weight: "0.6"},
		{Option: govtypesv1.OptionNo, Weight: "0.4"},
	}
	require.NoError(t, app.GovKeeper.AddVote(ctx, proposalID, voter, options, ""))

	vote, err := app.GovKeeper.Votes.Get(ctx, collections.Join(proposalID, voter))
	require.NoError(t, err)
	require.Len(t, vote.Options, 2, "stake ballot may carry multiple weighted options")
}
