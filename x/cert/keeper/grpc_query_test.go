package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"

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

func issueCert(t *testing.T, app *shentuapp.ShentuApp, ctx sdk.Context, certifier sdk.AccAddress, certType, content string) types.Certificate {
	t.Helper()
	cert, err := types.NewCertificate(certType, content, "", "", "", certifier)
	require.NoError(t, err)
	id, err := app.CertKeeper.IssueCertificate(ctx, cert)
	require.NoError(t, err)
	issued, err := app.CertKeeper.GetCertificateByID(ctx, id)
	require.NoError(t, err)
	return issued
}

// ---------------------------------------------------------------------------
// Certifier
// ---------------------------------------------------------------------------

func TestQueryCertifier_Found(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))

	expected := types.NewCertifier(addrs[0], "test certifier")
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
		require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addr, "")))
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

func TestQueryCertifiers_NilRequest(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.Certifiers(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Certificate
// ---------------------------------------------------------------------------

func TestQueryCertificate_Found(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))

	cert, err := types.NewCertificate(types.GeneralCertificateTypeName, "query-test", "", "", "", addrs[0])
	require.NoError(t, err)
	id, err := app.CertKeeper.IssueCertificate(ctx, cert)
	require.NoError(t, err)

	resp, err := querier.Certificate(sdk.WrapSDKContext(ctx), &types.QueryCertificateRequest{CertificateId: id})
	require.NoError(t, err)
	require.Equal(t, id, resp.Certificate.CertificateId)
}

func TestQueryCertificate_NotFound(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))

	_, err := querier.Certificate(sdk.WrapSDKContext(ctx), &types.QueryCertificateRequest{CertificateId: 99999})
	require.Error(t, err)
	require.True(t, errorsmod.IsOf(err, types.ErrCertificateNotExists))
}

func TestQueryCertificate_NilRequest(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.Certificate(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Certificates
// ---------------------------------------------------------------------------

func TestQueryCertificates_All(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))

	for i := 0; i < 5; i++ {
		issueCert(t, app, ctx, addrs[0], types.GeneralCertificateTypeName, fmt.Sprintf("cert-all-%d", i))
	}

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Certificates, 5)
	require.NotNil(t, resp.Pagination)
	require.Equal(t, uint64(5), resp.Pagination.Total)
}

func TestQueryCertificates_ByCertifier(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[1], "")))

	for i := 0; i < 3; i++ {
		issueCert(t, app, ctx, addrs[0], types.GeneralCertificateTypeName, fmt.Sprintf("addr0-%d", i))
	}
	issueCert(t, app, ctx, addrs[1], types.GeneralCertificateTypeName, "addr1-0")

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{Certifier: addrs[0].String()})
	require.NoError(t, err)
	require.Equal(t, uint64(3), resp.Pagination.Total)
	require.Len(t, resp.Certificates, 3)
}

func TestQueryCertificates_ByType(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))

	issueCert(t, app, ctx, addrs[0], types.AuditingCertificateTypeName, "audit-content")
	issueCert(t, app, ctx, addrs[0], types.GeneralCertificateTypeName, "general-content")

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{CertificateType: types.CertificateTypeAuditing})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.Pagination.Total)
	require.Len(t, resp.Certificates, 1)
	require.Equal(t, types.CertificateTypeAuditing, types.TranslateCertificateType(resp.Certificates[0]))
}

func TestQueryCertificates_ByContent(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[1], "")))

	issueCert(t, app, ctx, addrs[0], types.GeneralCertificateTypeName, "shared-content")
	issueCert(t, app, ctx, addrs[1], types.OpenMathCertificateTypeName, "shared-content")
	issueCert(t, app, ctx, addrs[0], types.GeneralCertificateTypeName, "other-content")

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{Content: "shared-content"})
	require.NoError(t, err)
	require.Equal(t, uint64(2), resp.Pagination.Total)
	require.Len(t, resp.Certificates, 2)
}

func TestQueryCertificates_ByOpenMathType(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))

	for _, prover := range []sdk.AccAddress{addrs[1], addrs[2]} {
		issueCert(t, app, ctx, addrs[0], types.OpenMathCertificateTypeName, prover.String())
	}
	issueCert(t, app, ctx, addrs[0], types.GeneralCertificateTypeName, "some-content")

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{CertificateType: types.CertificateTypeOpenMath})
	require.NoError(t, err)
	require.Equal(t, uint64(2), resp.Pagination.Total)

	resp, err = querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{CertificateType: types.CertificateTypeGeneral})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.Pagination.Total)

	resp, err = querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{})
	require.NoError(t, err)
	require.Equal(t, uint64(3), resp.Pagination.Total)
}

func TestQueryCertificates_ByCertifierAndType(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[1], "")))

	for _, content := range []string{"a", "b"} {
		issueCert(t, app, ctx, addrs[0], types.GeneralCertificateTypeName, content)
	}
	issueCert(t, app, ctx, addrs[0], types.OpenMathCertificateTypeName, addrs[2].String())
	issueCert(t, app, ctx, addrs[1], types.GeneralCertificateTypeName, "c")

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{
		Certifier:       addrs[0].String(),
		CertificateType: types.CertificateTypeGeneral,
	})
	require.NoError(t, err)
	require.Equal(t, uint64(2), resp.Pagination.Total)

	resp, err = querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{
		Certifier:       addrs[0].String(),
		CertificateType: types.CertificateTypeOpenMath,
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.Pagination.Total)

	resp, err = querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{
		Certifier:       addrs[1].String(),
		CertificateType: types.CertificateTypeOpenMath,
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.Pagination.Total)
}

func TestQueryCertificates_Pagination(t *testing.T) {
	app, ctx, querier := setupQuerier(t)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))

	for i := 0; i < 5; i++ {
		issueCert(t, app, ctx, addrs[0], types.GeneralCertificateTypeName, fmt.Sprintf("paged-%d", i))
	}

	firstPage, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{
		Pagination: &querytypes.PageRequest{Limit: 2, CountTotal: true},
	})
	require.NoError(t, err)
	require.Len(t, firstPage.Certificates, 2)
	require.Equal(t, uint64(5), firstPage.Pagination.Total)
	require.NotEmpty(t, firstPage.Pagination.NextKey)

	secondPage, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{
		Pagination: &querytypes.PageRequest{Key: firstPage.Pagination.NextKey, Limit: 2},
	})
	require.NoError(t, err)
	require.Len(t, secondPage.Certificates, 2)
	require.NotEmpty(t, secondPage.Pagination.NextKey)
	for _, cert := range secondPage.Certificates {
		require.NotZero(t, cert.CertificateId)
	}

	lastPage, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{
		Pagination: &querytypes.PageRequest{Key: secondPage.Pagination.NextKey, Limit: 2},
	})
	require.NoError(t, err)
	require.Len(t, lastPage.Certificates, 1)
	require.Empty(t, lastPage.Pagination.NextKey)
}

func TestQueryCertificates_InvalidType(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{CertificateType: types.CertificateType(99)})
	require.Error(t, err)
}

func TestQueryCertificates_Empty(t *testing.T) {
	_, ctx, querier := setupQuerier(t)

	resp, err := querier.Certificates(sdk.WrapSDKContext(ctx), &types.QueryCertificatesRequest{})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.Pagination.Total)
	require.Empty(t, resp.Certificates)
}

func TestQueryCertificates_NilRequest(t *testing.T) {
	_, ctx, querier := setupQuerier(t)
	_, err := querier.Certificates(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)
}
