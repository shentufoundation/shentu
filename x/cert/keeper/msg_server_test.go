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

	for _, certType := range []string{"general", "auditing", "proof", "identity", "openmath"} {
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
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	revokeMsg := types.NewMsgRevokeCertificate(addrs[0], 99999, "")
	_, err := msgServer.RevokeCertificate(sdk.WrapSDKContext(ctx), revokeMsg)
	require.Error(t, err)
	require.True(t, errorsmod.IsOf(err, types.ErrCertificateNotExists))
}

// ---------------------------------------------------------------------------
// OpenMath certificate
// ---------------------------------------------------------------------------

func TestMsgServerOpenMath_IssueAndQuery(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	// Issue an openmath certificate for a prover (addrs[1]).
	content := types.AssembleContent("openmath", addrs[1].String())
	require.NotNil(t, content)
	msg := types.NewMsgIssueCertificate(content, "", "", "openmath prover cert", addrs[0])
	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	// The prover's address should be certified under the openmath type.
	require.True(t, app.CertKeeper.IsCertified(ctx, addrs[1].String(), "openmath"))
	// But not under other types.
	require.False(t, app.CertKeeper.IsCertified(ctx, addrs[1].String(), "general"))
	// Content-level check should also find it.
	require.True(t, app.CertKeeper.IsContentCertified(ctx, addrs[1].String()))
}

func TestMsgServerOpenMath_MultipleProvers(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	// Certify addrs[1] and addrs[2] as openmath provers.
	for _, prover := range []sdk.AccAddress{addrs[1], addrs[2]} {
		content := types.AssembleContent("openmath", prover.String())
		msg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
		_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)
	}

	require.True(t, app.CertKeeper.IsCertified(ctx, addrs[1].String(), "openmath"))
	require.True(t, app.CertKeeper.IsCertified(ctx, addrs[2].String(), "openmath"))
	// addrs[3] was not certified.
	require.False(t, app.CertKeeper.IsCertified(ctx, addrs[3].String(), "openmath"))

	// Query by type should return exactly 2 openmath certificates.
	params := types.NewQueryCertificatesParams(1, 100, nil, types.CertificateTypeOpenMath)
	certs, pagination, err := app.CertKeeper.GetCertificatesFiltered(ctx, params)
	require.NoError(t, err)
	require.Equal(t, uint64(2), pagination.Total)
	require.Len(t, certs, 2)
}

func TestMsgServerOpenMath_RevokeProverCert(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	// Issue openmath cert.
	content := types.AssembleContent("openmath", addrs[1].String())
	msg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.True(t, app.CertKeeper.IsCertified(ctx, addrs[1].String(), "openmath"))

	// Revoke it.
	certs := app.CertKeeper.GetAllCertificates(ctx)
	require.Len(t, certs, 1)
	revokeMsg := types.NewMsgRevokeCertificate(addrs[0], certs[0].CertificateId, "revoke prover")
	_, err = msgServer.RevokeCertificate(sdk.WrapSDKContext(ctx), revokeMsg)
	require.NoError(t, err)

	// After revocation the prover is no longer certified.
	require.False(t, app.CertKeeper.IsCertified(ctx, addrs[1].String(), "openmath"))
	require.False(t, app.CertKeeper.IsContentCertified(ctx, addrs[1].String()))
}

func TestMsgServerOpenMath_NonCertifierCannotIssue(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	// addrs[0] is NOT a certifier.

	content := types.AssembleContent("openmath", addrs[1].String())
	msg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
	require.True(t, errorsmod.IsOf(err, types.ErrUnqualifiedCertifier))
}

// ---------------------------------------------------------------------------
// Cross-certifier operations
// ---------------------------------------------------------------------------

func TestMsgServerRevokeCertificate_CrossCertifierRevoke(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	// Both are certifiers.
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[1], addrs[1], "")))

	// addrs[0] issues a certificate.
	content := types.AssembleContent("general", "cross-revoke-content")
	msg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	certs := app.CertKeeper.GetAllCertificates(ctx)
	require.Len(t, certs, 1)

	// addrs[1] (a different certifier) revokes it — should succeed.
	revokeMsg := types.NewMsgRevokeCertificate(addrs[1], certs[0].CertificateId, "cross-revoke")
	_, err = msgServer.RevokeCertificate(sdk.WrapSDKContext(ctx), revokeMsg)
	require.NoError(t, err)

	require.False(t, app.CertKeeper.IsContentCertified(ctx, "cross-revoke-content"))
}

func TestMsgServerIssueCertificate_DuplicateContent(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	// Issue the same content twice — both should succeed (different certificate IDs).
	for i := 0; i < 2; i++ {
		content := types.AssembleContent("general", "duplicate-content")
		msg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
		_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)
	}

	certs := app.CertKeeper.GetAllCertificates(ctx)
	require.Len(t, certs, 2)
	require.NotEqual(t, certs[0].CertificateId, certs[1].CertificateId)
}

func TestMsgServerIssueCertificate_WithCompilationContent(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	content := types.AssembleContent("compilation", "source-hash")
	msg := types.NewMsgIssueCertificate(content, "solc-0.8.0", "0xdeadbeef", "compilation cert", addrs[0])
	_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	certs := app.CertKeeper.GetAllCertificates(ctx)
	require.Len(t, certs, 1)
	require.Equal(t, "solc-0.8.0", certs[0].CompilationContent.Compiler)
	require.Equal(t, "0xdeadbeef", certs[0].CompilationContent.BytecodeHash)
}

func TestMsgServerUpdateCertifier_RemoveNonExistent(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()

	// Add one certifier so there's at least one.
	_, err := msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), types.NewMsgUpdateCertifier(authority, addrs[0], "", types.Add, nil))
	require.NoError(t, err)

	// Try to remove addrs[1] which was never added.
	_, err = msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), types.NewMsgUpdateCertifier(authority, addrs[1], "", types.Remove, nil))
	require.Error(t, err)
	require.True(t, errorsmod.IsOf(err, types.ErrCertifierNotExists))
}

// ---------------------------------------------------------------------------
// OpenMath certificate — additional edge cases
// ---------------------------------------------------------------------------

func TestMsgServerOpenMath_DoesNotInterfereWithOtherTypes(t *testing.T) {
	app, ctx, msgServer := setupMsgServer(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	// Issue both an openmath and a general cert with the same content string.
	contentStr := addrs[1].String()
	for _, certType := range []string{"openmath", "general"} {
		content := types.AssembleContent(certType, contentStr)
		msg := types.NewMsgIssueCertificate(content, "", "", "", addrs[0])
		_, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)
	}

	// Each type should be independently queryable.
	require.True(t, app.CertKeeper.IsCertified(ctx, contentStr, "openmath"))
	require.True(t, app.CertKeeper.IsCertified(ctx, contentStr, "general"))

	// Revoke only the openmath cert.
	params := types.NewQueryCertificatesParams(1, 100, nil, types.CertificateTypeOpenMath)
	certs, _, err := app.CertKeeper.GetCertificatesFiltered(ctx, params)
	require.NoError(t, err)
	require.Len(t, certs, 1)

	revokeMsg := types.NewMsgRevokeCertificate(addrs[0], certs[0].CertificateId, "")
	_, err = msgServer.RevokeCertificate(sdk.WrapSDKContext(ctx), revokeMsg)
	require.NoError(t, err)

	// openmath is gone, general remains.
	require.False(t, app.CertKeeper.IsCertified(ctx, contentStr, "openmath"))
	require.True(t, app.CertKeeper.IsCertified(ctx, contentStr, "general"))
}
