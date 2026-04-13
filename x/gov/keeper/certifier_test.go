package keeper_test

import (
	"testing"

	"time"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/stretchr/testify/require"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
)

func TestCertifierVoteIsRequiredForMsgUpdateCertifier(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	governanceAuthority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	proposer, err := app.AccountKeeper.AddressCodec().StringToBytes(addrs[0].String())
	require.NoError(t, err)

	msg := certtypes.NewMsgUpdateCertifier(governanceAuthority, addrs[1], "governance certifier update", certtypes.Add, addrs[0])
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
