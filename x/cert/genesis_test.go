package cert_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	cert "github.com/shentufoundation/shentu/v2/x/cert"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func TestGenesisRoundTrip(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))

	// Seed certifiers and certificates.
	for _, addr := range addrs {
		require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addr, "genesis test")))
	}
	for i, content := range []string{"content-a", "content-b", "content-c"} {
		c, err := types.NewCertificate(types.GeneralCertificateTypeName, content, "", "", "", addrs[i%len(addrs)])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, c)
		require.NoError(t, err)
	}

	// Export genesis.
	exported := cert.ExportGenesis(ctx, app.CertKeeper)
	require.Len(t, exported.Certifiers, 3)
	require.Len(t, exported.Certificates, 3)
	require.True(t, exported.NextCertificateId > 0)

	// Set up a fresh app and import.
	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	cert.InitGenesis(ctx2, app2.CertKeeper, *exported)

	// Re-export and compare.
	reExported := cert.ExportGenesis(ctx2, app2.CertKeeper)
	require.Equal(t, len(exported.Certifiers), len(reExported.Certifiers))
	require.Equal(t, len(exported.Certificates), len(reExported.Certificates))
	require.Equal(t, exported.NextCertificateId, reExported.NextCertificateId)

	// Verify state is queryable after import.
	for _, addr := range addrs {
		_, err := app2.CertKeeper.GetCertifier(ctx2, addr)
		require.NoError(t, err)
	}
	require.True(t, app2.CertKeeper.IsContentCertified(ctx2, "content-a"))
	require.True(t, app2.CertKeeper.IsContentCertified(ctx2, "content-b"))
	require.True(t, app2.CertKeeper.IsContentCertified(ctx2, "content-c"))
}

func TestGenesisRoundTrip_OpenMath(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))

	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))

	// Issue openmath certs for two provers.
	for _, prover := range []string{addrs[1].String(), addrs[2].String()} {
		c, err := types.NewCertificate(types.OpenMathCertificateTypeName, prover, "", "", "prover cert", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, c)
		require.NoError(t, err)
	}

	require.True(t, app.CertKeeper.IsCertified(ctx, addrs[1].String(), types.OpenMathCertificateTypeName))

	// Export and re-import.
	exported := cert.ExportGenesis(ctx, app.CertKeeper)
	require.Len(t, exported.Certificates, 2)

	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	cert.InitGenesis(ctx2, app2.CertKeeper, *exported)

	// OpenMath certs survive the round trip.
	require.True(t, app2.CertKeeper.IsCertified(ctx2, addrs[1].String(), types.OpenMathCertificateTypeName))
	require.True(t, app2.CertKeeper.IsCertified(ctx2, addrs[2].String(), types.OpenMathCertificateTypeName))
	require.False(t, app2.CertKeeper.IsCertified(ctx2, addrs[0].String(), types.OpenMathCertificateTypeName))
}

func TestGenesisDefault(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	cert.InitDefaultGenesis(ctx, app.CertKeeper)

	exported := cert.ExportGenesis(ctx, app.CertKeeper)
	require.Empty(t, exported.Certifiers)
	require.Empty(t, exported.Certificates)
}

func TestGenesisRoundTrip_MultipleCertTypes(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))

	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))

	// Issue one of each type.
	certTypes := types.IssueableCertificateTypeNames()
	for i, ct := range certTypes {
		c, err := types.NewCertificate(ct, "content-"+ct, "", "", "", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, c)
		require.NoError(t, err, "failed to issue cert type %s (idx %d)", ct, i)
	}

	exported := cert.ExportGenesis(ctx, app.CertKeeper)
	require.Len(t, exported.Certificates, len(certTypes))

	// Re-import.
	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	cert.InitGenesis(ctx2, app2.CertKeeper, *exported)

	reExported := cert.ExportGenesis(ctx2, app2.CertKeeper)
	require.Equal(t, len(exported.Certificates), len(reExported.Certificates))

	// Verify each type survives.
	for _, ct := range certTypes {
		require.True(t, app2.CertKeeper.IsCertified(ctx2, "content-"+ct, ct), "type %s not found after reimport", ct)
	}
}

func TestGenesisRoundTrip_PreservesNextID(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))

	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], "")))

	// Issue 5 certificates to advance the next ID.
	for i := 0; i < 5; i++ {
		c, err := types.NewCertificate(types.GeneralCertificateTypeName, "id-test-"+string(rune('A'+i)), "", "", "", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, c)
		require.NoError(t, err)
	}

	exported := cert.ExportGenesis(ctx, app.CertKeeper)
	require.True(t, exported.NextCertificateId >= 6)

	// Re-import and issue one more — its ID should continue from the exported next ID.
	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	cert.InitGenesis(ctx2, app2.CertKeeper, *exported)

	c, err := types.NewCertificate(types.GeneralCertificateTypeName, "after-reimport", "", "", "", addrs[0])
	require.NoError(t, err)
	newID, err := app2.CertKeeper.IssueCertificate(ctx2, c)
	require.NoError(t, err)
	require.Equal(t, exported.NextCertificateId, newID)
}

func TestGenesisRoundTrip_MultipleCertifiers(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 5, math.NewInt(10000))

	for i, addr := range addrs {
		require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addr, "certifier-"+string(rune('A'+i)))))
	}

	exported := cert.ExportGenesis(ctx, app.CertKeeper)
	require.Len(t, exported.Certifiers, 5)

	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	cert.InitGenesis(ctx2, app2.CertKeeper, *exported)

	for _, addr := range addrs {
		_, err := app2.CertKeeper.GetCertifier(ctx2, addr)
		require.NoError(t, err)
	}
}
