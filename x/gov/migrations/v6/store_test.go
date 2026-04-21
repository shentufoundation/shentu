package v6_test

import (
	"encoding/binary"
	"testing"
	"time"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
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

// TestMigrate6to7_DeletesEveryCertVotedEntry confirms the migration's
// one and only job: sweep every key under CertVotesKeyPrefix. The
// migration is intentionally scan-free (see package doc), so this test
// also asserts it succeeds regardless of what lives in Proposals/Votes.
func TestMigrate6to7_DeletesEveryCertVotedEntry(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	ids := []uint64{101, 202, 303}
	for _, id := range ids {
		writeLegacyCertVotedEntry(t, ctx, app, id)
		require.True(t, certVotedHas(ctx, app, id))
	}

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	require.NoError(t, m.Migrate6to7(ctx))

	for _, id := range ids {
		require.False(t, certVotedHas(ctx, app, id),
			"cert_voted entry for proposal %d must be swept", id)
	}
}

// TestMigrate6to7_NoEntriesIsNoop confirms the migration is safe to
// run against a store with no legacy entries (the common case on any
// chain that never accumulated cert_voted state).
func TestMigrate6to7_NoEntriesIsNoop(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	require.NoError(t, m.Migrate6to7(ctx))
}

// TestMigrate6to7_TolerantOfUndecodableProposals is a regression test:
// a live proposal whose message type was removed from the v7 interface
// registry must not block the upgrade. The simplified migration never
// calls proposal.GetMsgs(), so this is true by construction, but the
// test encodes the contract so a future "add a scan" change has to
// explicitly break it.
func TestMigrate6to7_TolerantOfUndecodableProposals(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false).WithBlockTime(time.Unix(100, 0))

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 1, math.NewInt(10000))
	now := time.Unix(100, 0)
	endTime := now.Add(48 * time.Hour)

	bogus := &codectypes.Any{
		TypeUrl: "/shentu.deprecated.v1.MsgNoLongerExists",
		Value:   []byte{0x08, 0x01},
	}
	proposal := govtypesv1.Proposal{
		Id:              777001,
		Messages:        []*codectypes.Any{bogus},
		Status:          govtypesv1.StatusVotingPeriod,
		SubmitTime:      &now,
		DepositEndTime:  &endTime,
		VotingStartTime: &now,
		VotingEndTime:   &endTime,
		Title:           "legacy undecodable",
		Summary:         "deprecated msg type",
		Proposer:        addrs[0].String(),
	}
	require.NoError(t, app.GovKeeper.Proposals.Set(ctx, proposal.Id, proposal))
	writeLegacyCertVotedEntry(t, ctx, app, proposal.Id)

	m := govkeeper.NewMigrator(app.GovKeeper, nil)
	require.NoError(t, m.Migrate6to7(ctx))
	require.False(t, certVotedHas(ctx, app, proposal.Id))
}
