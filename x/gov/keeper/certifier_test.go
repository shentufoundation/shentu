package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
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

func TestCertifierVoteIsRequiredForLegacyCertifierUpdateProposal(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(10000))
	governanceAuthority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	proposer, err := app.AccountKeeper.AddressCodec().StringToBytes(addrs[0].String())
	require.NoError(t, err)

	content := certtypes.NewCertifierUpdateProposal("legacy certifier update", "summary", addrs[1], addrs[0], certtypes.Add)
	legacyMsg, err := govtypesv1.NewLegacyContent(content, governanceAuthority.String())
	require.NoError(t, err)

	proposal, err := app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{legacyMsg}, "", "legacy certifier update", "summary", proposer, false)
	require.NoError(t, err)

	required, err := app.GovKeeper.CertifierVoteIsRequired(ctx, proposal.Id)
	require.NoError(t, err)
	require.True(t, required)
}
