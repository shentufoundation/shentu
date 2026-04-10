package keeper_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func setupMsgServer(t *testing.T) (*shentuapp.ShentuApp, sdk.Context, types.MsgServer) {
	t.Helper()
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	msgServer := keeper.NewMsgServerImpl(app.CertKeeper)
	return app, ctx, msgServer
}

// ---------------------------------------------------------------------------
// MsgUpdateCertifier
// ---------------------------------------------------------------------------

func TestMsgServerUpdateCertifier_Add(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	msg := types.NewMsgUpdateCertifier(authority, addrs[0], "first", types.Add, addrs[1])
	_, err := msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	certifier, err := app.CertKeeper.GetCertifier(ctx, addrs[0])
	require.NoError(t, err)
	require.Equal(t, addrs[0].String(), certifier.Address)
	require.Equal(t, addrs[1].String(), certifier.Proposer)
}

func TestMsgServerUpdateCertifier_AddDuplicate(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	msg := types.NewMsgUpdateCertifier(authority, addrs[0], "", types.Add, addrs[1])
	_, err := msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	_, err = msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), msg)
	require.True(t, errorsmod.IsOf(err, types.ErrCertifierAlreadyExists))
}

func TestMsgServerUpdateCertifier_Remove(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	// Add two certifiers so one can be removed.
	_, err := msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), types.NewMsgUpdateCertifier(authority, addrs[0], "", types.Add, nil))
	require.NoError(t, err)
	_, err = msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), types.NewMsgUpdateCertifier(authority, addrs[1], "", types.Add, nil))
	require.NoError(t, err)

	_, err = msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), types.NewMsgUpdateCertifier(authority, addrs[0], "", types.Remove, nil))
	require.NoError(t, err)

	_, err = app.CertKeeper.GetCertifier(ctx, addrs[0])
	require.True(t, errorsmod.IsOf(err, types.ErrCertifierNotExists))
}

func TestMsgServerUpdateCertifier_RemoveLastFails(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	_, err := msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), types.NewMsgUpdateCertifier(authority, addrs[0], "", types.Add, nil))
	require.NoError(t, err)

	_, err = msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), types.NewMsgUpdateCertifier(authority, addrs[0], "", types.Remove, nil))
	require.True(t, errorsmod.IsOf(err, types.ErrOnlyOneCertifier))
}

func TestMsgServerUpdateCertifier_UnauthorizedAuthority(t *testing.T) {
	_, ctx, msgServer := setupMsgServer(t)
	fakeAuthority := sdk.AccAddress([]byte("fake_authority_addr_"))

	msg := types.NewMsgUpdateCertifier(fakeAuthority, fakeAuthority, "", types.Add, nil)
	_, err := msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), msg)
	require.True(t, errorsmod.IsOf(err, sdkerrors.ErrUnauthorized))
}

// ---------------------------------------------------------------------------
// MsgIssueCertificate
// ---------------------------------------------------------------------------

func TestMsgServerIssueCertificate_Success(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	content := types.AssembleContent("auditing", "some-content")
	msg := types.NewMsgIssueCertificate(content, "", "", "test cert", addrs[0])
	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	require.True(t, app.CertKeeper.IsCertified(ctx, "some-content", "auditing"))
}

func TestMsgServerIssueCertificate_NonCertifierFails(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	// addrs[0] is NOT a certifier

	content := types.AssembleContent("general", "some-content")
	msg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.True(t, errorsmod.IsOf(err, types.ErrUnqualifiedCertifier))
}

func TestMsgServerIssueCertificate_MultipleCertTypes(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	for _, certType := range []string{"general", "auditing", "proof", "identity"} {
		content := types.AssembleContent(certType, "content-"+certType)
		msg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
		_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err, "failed to issue %s certificate", certType)

		require.True(t, app.CertKeeper.IsCertified(ctx, "content-"+certType, certType))
	}
}

// ---------------------------------------------------------------------------
// MsgRevokeCertificate
// ---------------------------------------------------------------------------

func TestMsgServerRevokeCertificate_Success(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	content := types.AssembleContent("general", "revoke-me")
	issueMsg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), issueMsg)
	require.NoError(t, err)

	certs := app.CertKeeper.GetAllCertificates(ctx)
	require.Len(t, certs, 1)

	revokeMsg := types.NewMsgRevokeCertificate(addrs[0], certs[0].CertificateId, "revoking")
	_, err = msgServer.RevokeCertificate(sdk.WrapSDKContext(ctx), revokeMsg)
	require.NoError(t, err)

	require.False(t, app.CertKeeper.IsContentCertified(ctx, "revoke-me"))
}

func TestMsgServerRevokeCertificate_NonCertifierFails(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	content := types.AssembleContent("general", "revoke-test")
	issueMsg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), issueMsg)
	require.NoError(t, err)

	certs := app.CertKeeper.GetAllCertificates(ctx)
	require.Len(t, certs, 1)

	// addrs[1] is not a certifier
	revokeMsg := types.NewMsgRevokeCertificate(addrs[1], certs[0].CertificateId, "")
	_, err = msgServer.RevokeCertificate(sdk.WrapSDKContext(ctx), revokeMsg)
	require.True(t, errorsmod.IsOf(err, types.ErrUnqualifiedCertifier))
}

func TestMsgServerRevokeCertificate_NotFoundFails(t *testing.T) {
	app, ctx, _ := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	// GetCertificateByID with a non-existent ID does not return an error
	// (gogoproto MustUnmarshal on nil returns an empty struct). Verify that
	// HasCertificateByID correctly reports absence.
	has, err := app.CertKeeper.HasCertificateByID(ctx, 99999)
	require.NoError(t, err)
	require.False(t, has)
}
