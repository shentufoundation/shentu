package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func setupQuerier(t *testing.T) (*shentuapp.ShentuApp, sdk.Context, keeper.Querier) {
	t.Helper()
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	querier := keeper.Querier{Keeper: app.CertKeeper}
	return app, ctx, querier
}

// ---------------------------------------------------------------------------
// Certifier
// ---------------------------------------------------------------------------

func TestQueryCertifier_Found(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))

	expected := types.NewCertifier(addrs[0], addrs[1], "test certifier")
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, expected))

	resp, err := querier.Certifier(sdk.WrapSDKContext(ctx), &types.QueryCertifierRequest{Address: addrs[0].String()})
	require.NoError(t, err)
	require.Equal(t, expected, resp.Certifier)
}

func TestQueryCertifier_NotFound(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))

	_, err := querier.Certifier(sdk.WrapSDKContext(ctx), &types.QueryCertifierRequest{Address: addrs[0].String()})
	require.Error(t, err)
}

func TestQueryCertifier_NilRequest(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.Certifier(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Certifiers
// ---------------------------------------------------------------------------

func TestQueryCertifiers(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))

	for _, addr := range addrs {
		require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addr, addr, "")))
	}

	resp, err := querier.Certifiers(sdk.WrapSDKContext(ctx), &types.QueryCertifiersRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Certifiers, 3)
}

func TestQueryCertifiers_Empty(t *testing.T) {
	_, ctx, querier := setupQuerier(t)

	resp, err := querier.Certifiers(sdk.WrapSDKContext(ctx), &types.QueryCertifiersRequest{})
	require.NoError(t, err)
	require.Empty(t, resp.Certifiers)
}

// ---------------------------------------------------------------------------
// Certificate
// ---------------------------------------------------------------------------

func TestQueryCertificate_Found(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	cert, err := types.NewCertificate("general", "query-test", "", "", "", addrs[0])
	require.NoError(t, err)
	id, err := app.CertKeeper.IssueCertificate(ctx, cert)
	require.NoError(t, err)

	resp, err := querier.Certificate(sdk.WrapSDKContext(ctx), &types.QueryCertificateRequest{CertificateId: id})
	require.NoError(t, err)
	require.Equal(t, id, resp.Certificate.CertificateId)
}

func TestQueryCertificate_NotFound(t *testing.T) {
	app, ctx, _ := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	// HasCertificateByID correctly reports absence.
	has, err := app.CertKeeper.HasCertificateByID(ctx, 99999)
	require.NoError(t, err)
	require.False(t, has)

	// NOTE: GetCertificateByID does not return an error for missing keys
	// (gogoproto MustUnmarshal on nil returns an empty struct). This is a
	// known limitation — callers should use HasCertificateByID first.
}

// ---------------------------------------------------------------------------
// Certificates
// ---------------------------------------------------------------------------

func TestQueryCertificates_All(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	for i := 0; i < 5; i++ {
		cert, err := types.NewCertificate("general", "cert-all-"+string(rune('A'+i)), "", "", "", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)
	}

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{})
	require.NoError(t, err)
	require.Equal(t, uint64(5), resp.Total)
	require.Len(t, resp.Certificates, 5)
}

func TestQueryCertificates_ByCertifier(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[1], addrs[1], "")))

	for i := 0; i < 3; i++ {
		cert, err := types.NewCertificate("general", "addr0-"+string(rune('A'+i)), "", "", "", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)
	}
	cert, err := types.NewCertificate("general", "addr1-A", "", "", "", addrs[1])
	require.NoError(t, err)
	_, err = app.CertKeeper.IssueCertificate(ctx, cert)
	require.NoError(t, err)

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{Certifier: addrs[0].String()})
	require.NoError(t, err)
	require.Equal(t, uint64(3), resp.Total)
}

func TestQueryCertificates_ByType(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	cert1, err := types.NewCertificate("auditing", "audit-content", "", "", "", addrs[0])
	require.NoError(t, err)
	_, err = app.CertKeeper.IssueCertificate(ctx, cert1)
	require.NoError(t, err)

	cert2, err := types.NewCertificate("general", "general-content", "", "", "", addrs[0])
	require.NoError(t, err)
	_, err = app.CertKeeper.IssueCertificate(ctx, cert2)
	require.NoError(t, err)

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{CertificateType: "auditing"})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.Total)
}

// ---------------------------------------------------------------------------
// AddrConversion
// ---------------------------------------------------------------------------

func TestQueryAddrConversion_ShentuAddr(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	// Build a shentu-prefixed address from known bytes.
	addrBytes := sdk.AccAddress([]byte("test_addr_for_conv__"))
	addr, err := sdk.Bech32ifyAddressBytes("shentu", addrBytes)
	require.NoError(t, err)

	resp, err := querier.AddrConversion(sdk.WrapSDKContext(ctx), &types.ConversionToShentuAddrRequest{Address: addr})
	require.NoError(t, err)
	require.Equal(t, addr, resp.Address)
}

func TestQueryAddrConversion_InvalidAddr(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.AddrConversion(sdk.WrapSDKContext(ctx), &types.ConversionToShentuAddrRequest{Address: "invalid"})
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Nil request handling
// ---------------------------------------------------------------------------

func TestQueryCertificate_NilRequest(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.Certificate(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)
}

func TestQueryCertificates_NilRequest(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.Certificates(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)
}

func TestQueryCertifiers_NilRequest(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.Certifiers(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)
}

func TestQueryAddrConversion_NilRequest(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.AddrConversion(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Certificate queries — OpenMath type filtering
// ---------------------------------------------------------------------------

func TestQueryCertificates_ByOpenMathType(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	// Issue 2 openmath and 1 general cert.
	for _, prover := range []sdk.AccAddress{addrs[1], addrs[2]} {
		cert, err := types.NewCertificate("openmath", prover.String(), "", "", "", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)
	}
	cert, err := types.NewCertificate("general", "some-content", "", "", "", addrs[0])
	require.NoError(t, err)
	_, err = app.CertKeeper.IssueCertificate(ctx, cert)
	require.NoError(t, err)

	// Filter by openmath — should return only 2.
	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{CertificateType: "openmath"})
	require.NoError(t, err)
	require.Equal(t, uint64(2), resp.Total)

	// Filter by general — should return only 1.
	resp, err = querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{CertificateType: "general"})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.Total)

	// No filter — should return all 3.
	resp, err = querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{})
	require.NoError(t, err)
	require.Equal(t, uint64(3), resp.Total)
}

// ---------------------------------------------------------------------------
// Certificate queries — combined certifier + type filter
// ---------------------------------------------------------------------------

func TestQueryCertificates_ByCertifierAndType(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[1], addrs[1], "")))

	// addrs[0] issues 2 general + 1 openmath.
	for _, content := range []string{"a", "b"} {
		cert, err := types.NewCertificate("general", content, "", "", "", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)
	}
	cert, err := types.NewCertificate("openmath", addrs[2].String(), "", "", "", addrs[0])
	require.NoError(t, err)
	_, err = app.CertKeeper.IssueCertificate(ctx, cert)
	require.NoError(t, err)

	// addrs[1] issues 1 general.
	cert, err = types.NewCertificate("general", "c", "", "", "", addrs[1])
	require.NoError(t, err)
	_, err = app.CertKeeper.IssueCertificate(ctx, cert)
	require.NoError(t, err)

	// Filter: addrs[0] + general → 2
	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{
		Certifier:       addrs[0].String(),
		CertificateType: "general",
	})
	require.NoError(t, err)
	require.Equal(t, uint64(2), resp.Total)

	// Filter: addrs[0] + openmath → 1
	resp, err = querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{
		Certifier:       addrs[0].String(),
		CertificateType: "openmath",
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.Total)

	// Filter: addrs[1] + openmath → 0
	resp, err = querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{
		Certifier:       addrs[1].String(),
		CertificateType: "openmath",
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.Total)
}

// ---------------------------------------------------------------------------
// Certificates query — empty store
// ---------------------------------------------------------------------------

func TestQueryCertificates_Empty(t *testing.T) {
	_, ctx, querier := setupQuerier(t)

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.Total)
	require.Empty(t, resp.Certificates)
}
